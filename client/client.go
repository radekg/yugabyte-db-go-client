package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var recvChunkSize = 4 * 1024

// Connect connects to the master server without TLS.
func Connect(cfg *configs.YBClientConfig, logger hclog.Logger) (YBConnectedClient, error) {
	if logger == nil {
		logger = hclog.Default().Named("default-client-log")
	}
	if cfg.TLSConfig != nil {
		return connectTLS(cfg, logger)
	}
	return connect(cfg, logger)
}

func connect(cfg *configs.YBClientConfig, logger hclog.Logger) (YBConnectedClient, error) {
	logger.Debug("connecting non-TLS client")
	conn, err := net.Dial("tcp", cfg.MasterHostPort)
	if err != nil {
		return nil, err
	}
	client := &ybDefaultConnectedClient{
		originalConfig: cfg,
		chanConnected:  make(chan struct{}, 1),
		chanConnectErr: make(chan error, 1),
		closeFunc: func() error {
			return conn.Close()
		},
		conn:   conn,
		logger: logger,
	}
	return client.afterConnect(), nil
}

func connectTLS(cfg *configs.YBClientConfig, logger hclog.Logger) (YBConnectedClient, error) {
	logger.Debug("connecting TLS client")
	conn, err := tls.Dial("tcp", cfg.MasterHostPort, cfg.TLSConfig)
	if err != nil {
		return nil, err
	}
	client := &ybDefaultConnectedClient{
		originalConfig: cfg,
		chanConnected:  make(chan struct{}, 1),
		chanConnectErr: make(chan error, 1),
		closeFunc: func() error {
			return conn.Close()
		},
		conn:   conn,
		logger: logger,
	}
	return client.afterConnect(), nil
}

// Connected client

// YBConnectedClient represents a connected client.
type YBConnectedClient interface {
	Close() error

	GetLoadMoveCompletion() (*ybApi.GetLoadMovePercentResponsePB, error)
	GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error)
	GetUniverseConfig() (*ybApi.GetMasterClusterConfigResponsePB, error)
	IsLoadBalanced(*configs.OpIsLoadBalancedConfig) (*ybApi.IsLoadBalancedResponsePB, error)
	ListMasters() (*ybApi.ListMastersResponsePB, error)
	ListTabletServers(*configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error)
	Ping() (*ybApi.PingResponsePB, error)

	OnConnected() <-chan struct{}
	OnConnectError() <-chan error
}

type ybDefaultConnectedClient struct {
	originalConfig *configs.YBClientConfig
	callCounter    int
	chanConnected  chan struct{}
	chanConnectErr chan error
	closeFunc      func() error
	conn           net.Conn
	logger         hclog.Logger
}

// Close closes a connected client.
func (c *ybDefaultConnectedClient) Close() error {
	return c.closeFunc()
}

// GetLoadMoveCompletion gets the completion percentage of tablet load move from blacklisted servers.
func (c *ybDefaultConnectedClient) GetLoadMoveCompletion() (*ybApi.GetLoadMovePercentResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("GetLoadMoveCompletion"),
		},
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}
	payload := &ybApi.GetLoadMovePercentRequestPB{}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	buffer, err := c.recv() // TODO: can move this to readResponseInto
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.GetLoadMovePercentResponsePB{}
	readResponseErr := c.readResponseInto(buffer, responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// GetMasterRegistration retrieves the master registration information or error if the request failed.
func (c *ybDefaultConnectedClient) GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("GetMasterRegistration"),
		},
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}
	payload := &ybApi.GetMasterRegistrationRequestPB{}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	buffer, err := c.recv() // TODO: can move this to readResponseInto
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.GetMasterRegistrationResponsePB{}
	readResponseErr := c.readResponseInto(buffer, responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// GetUniverseConfig get the placement info and blacklist info of the universe.
func (c *ybDefaultConnectedClient) GetUniverseConfig() (*ybApi.GetMasterClusterConfigResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("GetMasterClusterConfig"),
		},
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}
	payload := &ybApi.GetMasterClusterConfigRequestPB{}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	buffer, err := c.recv() // TODO: can move this to readResponseInto
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.GetMasterClusterConfigResponsePB{}
	readResponseErr := c.readResponseInto(buffer, responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// IsLoadBalanced returns a list of masters or an error if call failed.
func (c *ybDefaultConnectedClient) IsLoadBalanced(opConfig *configs.OpIsLoadBalancedConfig) (*ybApi.IsLoadBalancedResponsePB, error) {

	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("IsLoadBalanced"),
		},
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}

	payload := &ybApi.IsLoadBalancedRequestPB{
		ExpectedNumServers: func() *int32 {
			if opConfig.ExpectedNumServers > 0 {
				return utils.PInt32(int32(opConfig.ExpectedNumServers))
			}
			return nil
		}(),
	}

	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	buffer, err := c.recv()
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.IsLoadBalancedResponsePB{}
	readResponseErr := c.readResponseInto(buffer, responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// ListMasters returns a list of masters or an error if call failed.
func (c *ybDefaultConnectedClient) ListMasters() (*ybApi.ListMastersResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("ListMasters"),
		},
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}
	payload := &ybApi.ListMastersRequestPB{}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	buffer, err := c.recv()
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.ListMastersResponsePB{}
	readResponseErr := c.readResponseInto(buffer, responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// ListTabletServers returns a list of tablet servers or an error if call failed.
func (c *ybDefaultConnectedClient) ListTabletServers(opConfig *configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("ListTabletServers"),
		},
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}
	payload := &ybApi.ListTabletServersRequestPB{
		PrimaryOnly: utils.PBool(opConfig.PrimaryOnly),
	}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	buffer, err := c.recv()
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.ListTabletServersResponsePB{}
	readResponseErr := c.readResponseInto(buffer, responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// Ping pings a certain YB server.
func (c *ybDefaultConnectedClient) Ping() (*ybApi.PingResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.server.GenericService"),
			MethodName:  utils.PString("Ping"),
		},
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}
	payload := &ybApi.PingRequestPB{}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	buffer, err := c.recv()
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.PingResponsePB{}
	readResponseErr := c.readResponseInto(buffer, responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	return responsePayload, nil
}

// OnConnected returns a channel which closed when the client is connected.
func (c *ybDefaultConnectedClient) OnConnected() <-chan struct{} {
	return c.chanConnected
}

// OnConnectError returns a channel which will return an error if connect fails.
func (c *ybDefaultConnectedClient) OnConnectError() <-chan error {
	return c.chanConnectErr
}

/// Private interface

func (c *ybDefaultConnectedClient) afterConnect() *ybDefaultConnectedClient {
	go func() {
		c.logger.Debug("sending connection header")
		header := append([]byte("YB"), 1)
		n, err := c.conn.Write(header)
		if err != nil {
			c.chanConnectErr <- err
			close(c.chanConnected)
			return
		}
		if n != len(header) {
			c.chanConnectErr <- fmt.Errorf("header not written: %d vs expected %d", n, len(header))
			close(c.chanConnected)
			return
		}
		c.logger.Debug("client connected")
		close(c.chanConnected)
	}()
	return c
}

func (c *ybDefaultConnectedClient) callID() int {
	currentID := c.callCounter
	c.callCounter = c.callCounter + 1
	return currentID
}

func (c *ybDefaultConnectedClient) recv() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	for {
		chunk := make([]byte, recvChunkSize)
		n, err := c.conn.Read(chunk)
		if err != nil {
			return buf, err
		}
		// we read an EOF, finished reading
		// previous iteration was fitting
		// all in one recvChunkSize
		if n == 4 && chunk[0] == 0 && chunk[1] == 0 && chunk[2] == 0 && chunk[3] == 0 {
			break
		}
		// otherwise, read what we've got:
		buf.Write(chunk[0:n])
		// if we have read exactly recvChunkSize, we continue reading
		if n == recvChunkSize {
			continue
		}
		// we can't read more so this implies we read less
		break
	}
	return buf, nil
}

func (c *ybDefaultConnectedClient) send(buf *bytes.Buffer) error {
	nBytesToWrite := buf.Len()
	n, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if n != nBytesToWrite {
		return fmt.Errorf("not all bytes written: %d vs expected %d", n, nBytesToWrite)
	}
	return nil
}

func (c *ybDefaultConnectedClient) readResponseInto(reader *bytes.Buffer, m protoreflect.ProtoMessage) error {

	opLogger := c.logger.Named("read-response-into").With("message", m.ProtoReflect().Type().Descriptor().Name())

	// Read the complete data length:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L71
	dataLength, err := utils.ReadInt(reader)
	if err != nil {
		opLogger.Error("failed reading data length", "reason", err)
		return err
	}

	opLogger.Trace("data-length", "value", dataLength)

	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L76
	responseHeaderLength, _, err := utils.ReadVarint(reader)
	if err != nil {
		opLogger.Error("failed reading response header length", "reason", err)
		return err
	}

	opLogger.Trace("response-header-length", "value", responseHeaderLength)

	// Now I can read the response header:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L78
	responseHeaderBuf := make([]byte, responseHeaderLength)
	n, err := reader.Read(responseHeaderBuf)
	if err != nil {
		opLogger.Error("failed reading response header", "reason", err)
		return err
	}

	opLogger.Trace("response-header-read",
		"expected-header-length", responseHeaderLength,
		"read-header-length", n)

	if uint64(n) != responseHeaderLength {
		opLogger.Error("response header read bytes count != expected count",
			"expected-header-length", responseHeaderLength,
			"read-header-length", n)
		return fmt.Errorf("expected to read %d but read %d", responseHeaderLength, n)
	}

	responseHeader := &ybApi.ResponseHeader{}
	protoErr := proto.Unmarshal(responseHeaderBuf, responseHeader)
	if protoErr != nil {
		opLogger.Error("failed unmarshalling response header", "reason", err)
		return err
	}

	opLogger = opLogger.With("call-id", *responseHeader.CallId,
		"is-error", *responseHeader.IsError,
		"sidecars-count", len(responseHeader.SidecarOffsets))

	// This here is currently a guess but I believe the corretc mechanism sits here:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L113
	// The encoding/binary.ReadUvarint and encoding/binary.ReadVarint doesn't do what it supposed to do
	// hence the custom code here.
	responsePayloadLength, _, err := utils.ReadVarint(reader)
	if err != nil {
		opLogger.Error("failed reading response payload length", "reason", err)
		return err
	}

	opLogger.Trace("response-payload-length", "value", responsePayloadLength)

	// if there was no data but the call did not result in an error,
	// return successful no data response:
	if !*responseHeader.IsError && responsePayloadLength == 0 {
		opLogger.Debug("payload was empty but no error, assuming OK")
		return nil
	}

	responsePayloadBuf := make([]byte, responsePayloadLength)
	n, err = reader.Read(responsePayloadBuf)
	if err != nil {
		opLogger.Error("failed reading response payload", "reason", err)
		return err
	}

	opLogger.Trace("response-payload-read",
		"expected-payload-length", responsePayloadLength,
		"read-payload-length", n)

	if uint64(n) != responsePayloadLength {
		opLogger.Error("response payload read bytes count != expected count",
			"expected-payload-length", responsePayloadLength,
			"read-payload-length", n)
		return fmt.Errorf("expected to read %d but read %d", responsePayloadLength, n)
	}

	protoErr2 := proto.Unmarshal(responsePayloadBuf, m)
	if protoErr2 != nil {
		opLogger.Error("failed unmarshalling response payload", "reason", err)
		return err
	}

	return nil
}

func (c *ybDefaultConnectedClient) sendMessages(msgs ...protoreflect.ProtoMessage) error {
	b := bytes.NewBuffer([]byte{})
	if err := utils.WriteMessages(b, msgs...); err != nil {
		return err
	}
	if err := c.send(b); err != nil {
		return err
	}
	return nil
}
