package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Connect connects to the master server without TLS.
func Connect(cfg *YBClientConfig) (YBConnectedClient, error) {
	if cfg.TLSConfig != nil {
		return connectTLS(cfg)
	}
	return connect(cfg)
}

func connect(cfg *YBClientConfig) (YBConnectedClient, error) {
	conn, err := net.Dial("tcp", cfg.MasterHostPort)
	if err != nil {
		return nil, err
	}
	client := &ybDefaultConnectedClient{
		chanConnected:  make(chan struct{}, 1),
		chanConnectErr: make(chan error, 1),
		closeFunc: func() error {
			return conn.Close()
		},
		conn: conn,
	}
	return client.doConnect(), nil
}

func connectTLS(cfg *YBClientConfig) (YBConnectedClient, error) {
	conn, err := tls.Dial("tcp", cfg.MasterHostPort, cfg.TLSConfig)
	if err != nil {
		return nil, err
	}
	client := &ybDefaultConnectedClient{
		chanConnected:  make(chan struct{}, 1),
		chanConnectErr: make(chan error, 1),
		closeFunc: func() error {
			return conn.Close()
		},
		conn: conn,
	}
	return client.doConnect(), nil
}

// Client config

// YBClientConfig is the client configuration.
type YBClientConfig struct {
	MasterHostPort string
	TLSConfig      *tls.Config
}

// Connected client

// YBConnectedClient represents a connected client.
type YBConnectedClient interface {
	Close() error

	GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error)
	ListMasters() (*ybApi.ListMastersResponsePB, error)
	ListTabletServers() (*ybApi.ListTabletServersResponsePB, error)

	OnConnected() <-chan struct{}
	OnConnectError() <-chan error
}

type ybDefaultConnectedClient struct {
	chanConnected  chan struct{}
	chanConnectErr chan error
	closeFunc      func() error
	conn           net.Conn
	callCounter    int
}

// Close closes a connected client.
func (c *ybDefaultConnectedClient) Close() error {
	return c.closeFunc()
}

// GetMasterRegistration retrieves the master registration information or error if the request failed.
func (c *ybDefaultConnectedClient) GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("GetMasterRegistration"),
		},
		TimeoutMillis: utils.PUint32(5000), // TODO: must be customizable
	}
	payload := &ybApi.GetMasterRegistrationRequestPB{}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	responseBytes, err := c.recv()
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.GetMasterRegistrationResponsePB{}
	readResponseErr := c.readResponseInto(bytes.NewReader(responseBytes), responsePayload)
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
		TimeoutMillis: utils.PUint32(5000), // TODO: must be customizable
	}
	payload := &ybApi.ListMastersRequestPB{}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	responseBytes, err := c.recv()
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.ListMastersResponsePB{}
	readResponseErr := c.readResponseInto(bytes.NewReader(responseBytes), responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// ListTabletServers returns a list of tablet servers or an error if call failed.
func (c *ybDefaultConnectedClient) ListTabletServers() (*ybApi.ListTabletServersResponsePB, error) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(int32(c.callID())),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("ListTabletServers"),
		},
		TimeoutMillis: utils.PUint32(5000), // TODO: must be customizable
	}
	payload := &ybApi.ListTabletServersRequestPB{PrimaryOnly: utils.PBool(false)}
	if err := c.sendMessages(requestHeader, payload); err != nil {
		return nil, err
	}
	responseBytes, err := c.recv()
	if err != nil {
		return nil, err
	}
	responsePayload := &ybApi.ListTabletServersResponsePB{}
	readResponseErr := c.readResponseInto(bytes.NewReader(responseBytes), responsePayload)
	if readResponseErr != nil {
		return nil, readResponseErr
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
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

func (c *ybDefaultConnectedClient) callID() int {
	currentID := c.callCounter
	c.callCounter = c.callCounter + 1
	return currentID
}

func (c *ybDefaultConnectedClient) doConnect() *ybDefaultConnectedClient {
	go func() {
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
		close(c.chanConnected)
	}()
	return c
}

func (c *ybDefaultConnectedClient) recv() ([]byte, error) {
	buf := make([]byte, 1024*1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return buf, err
	}
	return buf[0:n], nil
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

func (c *ybDefaultConnectedClient) readResponseInto(reader *bytes.Reader, m protoreflect.ProtoMessage) error {
	// Read the complete data length:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L71
	dataLength, err := utils.ReadInt(reader)
	if err != nil {
		return err
	}
	fmt.Println("DEBUG: the response data length is: ", dataLength)

	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L76
	responseHeaderLength, _, err := utils.ReadVarint(reader)
	if err != nil {
		return err
	}

	// Now I can read the response header:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L78
	responseHeaderBuf := make([]byte, responseHeaderLength)
	n, err := reader.Read(responseHeaderBuf)
	if err != nil {
		return err
	}
	if uint64(n) != responseHeaderLength {
		panic(fmt.Errorf("expected to read %d but read %d", responseHeaderLength, n))
	}

	responseHeader := &ybApi.ResponseHeader{}
	protoErr := proto.Unmarshal(responseHeaderBuf, responseHeader)
	if protoErr != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("DEBUG: Response to call id: %d, is error: %v, # of sidecars: %d",
		*responseHeader.CallId,
		*responseHeader.IsError,
		len(responseHeader.SidecarOffsets)))

	// This here is currently a guess but I believe the corretc mechanism sits here:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L113
	// The encoding/binary.ReadUvarint and encoding/binary.ReadVarint doesn't do what it supposed to do
	// hence the custom code here.
	responsePayloadLength, _, err := utils.ReadVarint(reader)
	if err != nil {
		panic(err)
	}

	// if there was no data but the call did not result in an error,
	// return successful no data response:
	if !*responseHeader.IsError && responsePayloadLength == 0 {
		return nil
	}

	responsePayloadBuf := make([]byte, responsePayloadLength)
	n, err = reader.Read(responsePayloadBuf)
	if err != nil {
		return err
	}
	if uint64(n) != responsePayloadLength {
		return fmt.Errorf("expected to read %d but read %d", responsePayloadLength, n)
	}

	protoErr2 := proto.Unmarshal(responsePayloadBuf, m)
	if protoErr2 != nil {
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
