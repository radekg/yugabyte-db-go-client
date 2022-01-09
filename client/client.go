package client

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/configs"
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
	errAlreadyConnected = fmt.Errorf("client: already connected")
	errConnecting       = fmt.Errorf("client: connecting")
	errNoClient         = fmt.Errorf("client: no client")
)

type defaultYBClient struct {
	config          *configs.YBClientConfig
	connectedClient YBConnectedClient
	isConnecting    bool
	lock            *sync.Mutex
	logger          hclog.Logger
}

// NewYBClient constructs a new instance of the high-level YugabyteDB client.
func NewYBClient(config *configs.YBClientConfig, logger hclog.Logger) YBClient {
	return &defaultYBClient{
		config: config,
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

	closeError := c.connectedClient.Close()
	c.connectedClient = nil
	return closeError
}

func (c *defaultYBClient) Connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.isConnecting {
		return errConnecting
	}

	if c.connectedClient != nil {
		return errAlreadyConnected
	}

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
				return fmt.Errorf("no reachable leaders")
			}
		case connectedClient := <-chanConnectedClient:
			c.connectedClient = connectedClient
			c.isConnecting = false
			return nil
		case <-time.After(c.config.OpTimeout):
			c.isConnecting = false
			return fmt.Errorf("failed to connect to a leader master within timeout")
		}
	}

}

func (c *defaultYBClient) Execute(payload, response protoreflect.ProtoMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.connectedClient == nil {
		return errNoClient
	}
	return c.connectedClient.Execute(payload, response)
}
