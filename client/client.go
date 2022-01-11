package client

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"google.golang.org/protobuf/reflect/protoreflect"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// YBClient is a high-level YugabyteDB client.
type YBClient interface {
	// Close closes the connected client.
	Close() error
	// Connects the client.
	Connect() error
	// Execute executes the payload against the service
	// and populates the response with the response data.
	Execute(payload, response protoreflect.ProtoMessage) error
}

var (
	errConnected         = fmt.Errorf(clientErrors.ErrorMessageConnected)
	errConnecting        = fmt.Errorf(clientErrors.ErrorMessageConnecting)
	errLeaderWaitTimeout = fmt.Errorf(clientErrors.ErrorMessageLeaderWaitTimeout)
	errNoClient          = fmt.Errorf(clientErrors.ErrorMessageNoClient)
	errNotConnected      = fmt.Errorf(clientErrors.ErrorMessageNotConnected)
	errNotReconnected    = fmt.Errorf(clientErrors.ErrorMessageReconnectFailed)
)

type defaultYBClient struct {
	config          *configs.YBClientConfig
	connectedClient YBConnectedClient
	isConnecting    bool
	isConnected     bool
	lock            *sync.Mutex
	logger          hclog.Logger
}

// NewYBClient constructs a new instance of the high-level YugabyteDB client.
func NewYBClient(config *configs.YBClientConfig, logger hclog.Logger) YBClient {
	return &defaultYBClient{
		config: config.WithDefaults(),
		lock:   &sync.Mutex{},
		logger: logger,
	}
}

func (c *defaultYBClient) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.connectedClient == nil {
		return errNoClient
	}
	closeError := c.closeUnsafe()
	c.isConnected = false
	c.connectedClient = nil
	return closeError
}

func (c *defaultYBClient) Connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.isConnecting {
		return errConnecting
	}
	if c.isConnected || c.connectedClient != nil {
		return errConnected
	}
	return c.connectUnsafe()
}

func (c *defaultYBClient) Execute(payload, response protoreflect.ProtoMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.isConnected {
		return errNotConnected
	}

	if c.connectedClient == nil {
		return errNoClient
	}

	currentAttempt := int32(1)

	for {

		executeErr := c.connectedClient.Execute(payload, response)

		// the response might have an error in it, check if this is a response returning ybApi.MasterErrorPB
		if tResponse, ok := response.(clientErrors.AbstractMasterErrorResponse); ok {
			// was there an error in that response?
			if masterError := clientErrors.NewMasterError(tResponse.GetError()); masterError != nil {
				// we know it was not nil because masterError was not nil
				responseError := tResponse.GetError()
				if responseError.Code != nil {
					responseErrorCode := *responseError.Code
					if int32(responseErrorCode.Number()) == int32(ybApi.MasterErrorPB_NOT_THE_LEADER.Number()) {
						c.logger.Warn("execute: response with NOT_THE_LEADER master status code, reconnect", "reason", masterError)
						executeErr = &clientErrors.RequiresReconnectError{
							Cause: masterError,
						}
					}
				}
			}
		}

		if executeErr == nil {
			return nil
		}

		if c.config.MaxExecuteRetries <= configs.NoExecuteRetry {
			reportErr := executeErr
			if tReconnectError, ok := executeErr.(*clientErrors.RequiresReconnectError); ok {
				reportErr = tReconnectError.Cause
			}
			c.logger.Error("execute: retry disabled, not retrying", "reason", reportErr)
			return reportErr
		}

		if currentAttempt > c.config.MaxExecuteRetries {
			reportErr := executeErr
			if tReconnectError, ok := executeErr.(*clientErrors.RequiresReconnectError); ok {
				reportErr = tReconnectError.Cause
			}
			c.logger.Error("execute: failed for a maximum number of allowed attempts, giving up", "reason", reportErr)
			return reportErr
		}

		if errors.Is(executeErr, syscall.EPIPE) {
			// broken pipe qualifies for immediate retry:
			executeErr = &clientErrors.RequiresReconnectError{
				Cause: executeErr,
			}
		} else if _, ok := executeErr.(*clientErrors.SendReceiveError); ok {
			// the client was connected but is no longer able to
			// communicate with the server, this qualifies
			// for reconnect
			executeErr = &clientErrors.RequiresReconnectError{
				Cause: executeErr,
			}
		} else if _, ok := executeErr.(*clientErrors.UnprocessableResponseError); ok {
			// complete payload has been read from the server
			// but payload could not be deserialized as protobuf,
			// this qualifies for immediate retry:
			currentAttempt = currentAttempt + 1
			<-time.After(c.config.RetryInterval)
			continue
		}

		if tReconnectError, ok := executeErr.(*clientErrors.RequiresReconnectError); ok {

			if c.config.MaxReconnectAttempts <= configs.NoReconnectAttempts {
				c.logger.Error("execute: not reconnecting after error, max reconnect attempts not set",
					"reason", tReconnectError.Cause)
				return tReconnectError.Cause
			}

			c.logger.Debug("execute: attempting reconnect due to an error",
				"reason", tReconnectError.Cause)

			// reconnect:
			currentReconnectAttempt := int32(1)
			reconnected := false
			for {

				reconnectErr := c.reconnect()

				if reconnectErr == nil {
					reconnected = true
					break
				}

				if currentReconnectAttempt == c.config.MaxReconnectAttempts {
					break
				}

				c.logger.Error("execute: failed reconnect",
					"attempt", currentReconnectAttempt,
					"max-attempts", c.config.MaxReconnectAttempts,
					"reason", reconnectErr)

				currentReconnectAttempt = currentReconnectAttempt + 1
				<-time.After(c.config.ReconnectRetryInterval)

			}

			if !reconnected {
				c.logger.Error("execute: failed reconnect consecutive maximum reconnect attempts",
					"max-attempts", c.config.MaxReconnectAttempts,
					"reason", tReconnectError.Cause)
				return fmt.Errorf("%s: %s", errNotReconnected.Error(), tReconnectError.Cause.Error())
			}

			// retry:
			<-time.After(c.config.RetryInterval)
			currentAttempt = currentAttempt + 1
			continue

		} // reconnect handling / end

		// in case of any other error, no recovery:
		return executeErr

	}

}

func (c *defaultYBClient) closeUnsafe() error {
	return c.connectedClient.Close()
}

func (c *defaultYBClient) connectUnsafe() error {

	tlsConfig, err := c.config.TLSConfig()
	if err != nil {
		return err
	}

	c.isConnecting = true

	validConfigs := map[string]*configs.YBSingleNodeClientConfig{}
	for _, hostPort := range c.config.MasterHostPort {
		validConfigs[hostPort] = &configs.YBSingleNodeClientConfig{
			MasterHostPort: hostPort,
			TLSConfig:      tlsConfig,
			OpTimeout:      uint32(c.config.OpTimeout.Milliseconds()),
		}
	}

	chanConnectedClient := make(chan YBConnectedClient, 1)
	chanErrors := make(chan error)
	var done uint64
	max := uint64(len(validConfigs))

	for hostPort, cliConfig := range validConfigs {
		go func(thisHostPort string, thisConfig *configs.YBSingleNodeClientConfig) {
			singleNodeClient, err := Connect(thisConfig, c.logger.Named("client"))
			if err != nil {
				c.logger.Error("failed creating a client",
					"reason", err,
					"host-port", thisHostPort)
				chanErrors <- err
				return
			}

			select {
			case err := <-singleNodeClient.OnConnectError():

				c.logger.Error("connection error",
					"reason", err,
					"host-port", thisHostPort)
				singleNodeClient.Close()
				chanErrors <- err

			case <-singleNodeClient.OnConnected():

				masterRegistration, err := singleNodeClient.GetMasterRegistration()

				if err != nil {
					c.logger.Error("failed querying master registration",
						"reason", err,
						"host-port", thisHostPort)
					singleNodeClient.Close()
					chanErrors <- err
					return
				}

				if masterRegistration == nil {
					c.logger.Trace("master did not send with registration info",
						"host-port", thisHostPort)
					singleNodeClient.Close()
					chanErrors <- fmt.Errorf("master %s did not send registration info", thisHostPort)
					return
				}

				if masterRegistration.Role == nil {
					c.logger.Trace("master did not report its raft peer role",
						"host-port", thisHostPort)
					singleNodeClient.Close()
					chanErrors <- fmt.Errorf("master %s did not report its raft peer role", thisHostPort)
					return
				}

				if *masterRegistration.Role != ybApi.RaftPeerPB_LEADER {
					c.logger.Trace("master not leader",
						"host-port", thisHostPort)
					singleNodeClient.Close()
					chanErrors <- fmt.Errorf("master %s not leader", thisHostPort)
					return
				}

				c.logger.Info("found master leader", "host-port", thisHostPort)
				chanConnectedClient <- singleNodeClient
			}

		}(hostPort, cliConfig)
	}

	for {
		select {
		case <-chanErrors:
			atomic.AddUint64(&done, 1)
			if atomic.LoadUint64(&done) == max {
				c.isConnecting = false
				return &clientErrors.NoLeaderError{}
			}
		case connectedClient := <-chanConnectedClient:
			c.connectedClient = connectedClient
			c.isConnecting = false
			c.isConnected = true
			return nil
		case <-time.After(c.config.OpTimeout):
			c.isConnecting = false
			return errLeaderWaitTimeout
		}
	}

}

func (c *defaultYBClient) reconnect() error {
	// ignore close error
	// if the client isn't connected, it does not matter to us
	//
	c.closeUnsafe()
	return c.connectUnsafe()
}
