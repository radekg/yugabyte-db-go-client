package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/metrics"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var recvChunkSize = 4 * 1024

// YBConnectedClient represents a connected client.
type YBConnectedClient interface {
	ClientID() string
	// Close closes the connected client.
	Close() error
	// Execute executes the payload against the service
	// and populates the response with the response data.
	Execute(payload, response protoreflect.ProtoMessage) error
	// Retrieves the master registration information or error if the request failed.
	GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error)
	// Returns a channel which closed when the client is connected.
	OnConnected() <-chan struct{}
	// Returns a channel which will return an error if connect fails.
	OnConnectError() <-chan error
}

type defaultSingleNodeClient struct {
	id              string
	originalConfig  *configs.YBSingleNodeClientConfig
	callCounter     int
	chanConnected   chan struct{}
	chanConnectErr  chan error
	closeFunc       func() error
	conn            net.Conn
	logger          hclog.Logger
	metricsCallback metrics.Callback
	svcRegistry     ServiceRegistry
}

// Close closes a connected client.
func (c *defaultSingleNodeClient) ClientID() string {
	return c.id
}

// Close closes a connected client.
func (c *defaultSingleNodeClient) Close() error {
	return c.closeFunc()
}

// Execute executes the payload against the service
// and populates the response with the response data.
func (c *defaultSingleNodeClient) Execute(payload, response protoreflect.ProtoMessage) error {
	return c.executeOp(payload, response)
}

// GetMasterRegistration retrieves the master registration information or error if the request failed.
func (c *defaultSingleNodeClient) GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error) {
	payload := &ybApi.GetMasterRegistrationRequestPB{}
	responsePayload := &ybApi.GetMasterRegistrationResponsePB{}
	if err := c.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, errors.NewMasterError(responsePayload.Error)
}

// OnConnected returns a channel which closed when the client is connected.
func (c *defaultSingleNodeClient) OnConnected() <-chan struct{} {
	return c.chanConnected
}

// OnConnectError returns a channel which will return an error if connect fails.
func (c *defaultSingleNodeClient) OnConnectError() <-chan error {
	return c.chanConnectErr
}

func (c *defaultSingleNodeClient) WithMetricsCallback(callback metrics.Callback) YBConnectedClient {
	c.metricsCallback = callback
	return c
}

/// Private interface

func (c *defaultSingleNodeClient) afterConnect() *defaultSingleNodeClient {
	loadServiceDefinitions(c.svcRegistry)
	go func() {
		c.logger.Debug("sending connection header")
		header := append([]byte("YB"), 1)
		n, err := c.conn.Write(header)
		if err != nil {
			c.chanConnectErr <- &errors.ProtocolConnectionHeaderWriteError{
				Cause: err,
			}
			close(c.chanConnected)
			return
		}
		if n != len(header) {
			c.chanConnectErr <- &errors.ProtocolConnectionHeaderWriteIncompleteError{
				Header:   header,
				Expected: len(header),
				Written:  n,
			}
			close(c.chanConnected)
			return
		}
		c.logger.Debug("client connected")
		close(c.chanConnected)
	}()
	return c
}

func (c *defaultSingleNodeClient) callID() int {
	currentID := c.callCounter
	c.callCounter = c.callCounter + 1
	return currentID
}

func (c *defaultSingleNodeClient) recv() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	for {
		chunk := make([]byte, recvChunkSize)
		n, err := c.conn.Read(chunk)
		if err != nil {
			return buf, err
		}
		c.metricsCallback.ClientBytesReceived(n)
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

func (c *defaultSingleNodeClient) send(buf *bytes.Buffer) error {
	nBytesToWrite := buf.Len()
	n, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	c.metricsCallback.ClientBytesSent(n)
	if n != nBytesToWrite {
		return fmt.Errorf("write incomplete: %d bytes vs %d expected", n, nBytesToWrite)
	}
	return nil
}

func (c *defaultSingleNodeClient) readResponseInto(reader *bytes.Buffer, m protoreflect.ProtoMessage) error {

	opLogger := c.logger.With("message", m.ProtoReflect().Type().Descriptor().Name())

	// Read the complete data length:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L71
	var dataLength int32
	if err := binary.Read(reader, binary.BigEndian, &dataLength); err != nil {
		opLogger.Error("failed reading response data length", "reason", err)
		return &errors.ReceiveError{
			Cause: fmt.Errorf("response data length read failed: %s", err.Error()),
		}
	}

	opLogger.Trace("data-length", "value", dataLength)

	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L76
	responseHeaderLength, err := utils.ReadUvarint32(reader)
	if err != nil {
		opLogger.Error("failed reading response header length", "reason", err)
		return &errors.ReceiveError{
			Cause: fmt.Errorf("response header length read failed: %s", err.Error()),
		}
	}

	opLogger.Trace("response-header-length", "value", responseHeaderLength)

	// Now I can read the response header:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L78
	responseHeaderBuf := make([]byte, responseHeaderLength)
	n, err := reader.Read(responseHeaderBuf)
	if err != nil {
		opLogger.Error("failed reading response header", "reason", err)
		return &errors.ReceiveError{
			Cause: fmt.Errorf("response header read failed: %s", err.Error()),
		}
	}

	opLogger.Trace("response-header-read",
		"expected-header-length", responseHeaderLength,
		"read-header-length", n)

	if uint64(n) != responseHeaderLength {
		opLogger.Trace("response header read bytes count != expected count",
			"read-data", string(responseHeaderBuf))
		opLogger.Error("response header read bytes count != expected count",
			"expected-header-length", responseHeaderLength,
			"read-header-length", n)
		return &errors.ReceiveError{
			Cause: fmt.Errorf("response header incomplete: read %d bytes vs expected %d",
				n, responseHeaderLength),
		}
	}

	responseHeader := &ybApi.ResponseHeader{}
	protoErr := utils.DeserializeProto(responseHeaderBuf, responseHeader)
	if protoErr != nil {
		opLogger.Error("failed unmarshalling response header", "reason", protoErr)
		return &errors.ReceiveError{
			Cause: fmt.Errorf("response header unprocessable: %s", protoErr.Error()),
		}
	}

	opLogger = opLogger.With("call-id", *responseHeader.CallId,
		"is-error", *responseHeader.IsError,
		"sidecars-count", len(responseHeader.SidecarOffsets))

	// This here is currently a guess but I believe the corretc mechanism sits here:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L113
	// The encoding/binary.ReadUvarint and encoding/binary.ReadVarint doesn't do what it supposed to do
	// hence the custom code here.
	responsePayloadLength, err := utils.ReadUvarint32(reader)
	if err != nil {
		opLogger.Error("failed reading response payload length", "reason", err)
		return &errors.ReceiveError{
			Cause: fmt.Errorf("response payload length read failed: %s", err.Error()),
		}
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
		return &errors.ReceiveError{
			Cause: fmt.Errorf("response data read failed: %s", err.Error()),
		}
	}

	opLogger.Trace("response-payload-read",
		"expected-payload-length", responsePayloadLength,
		"read-payload-length", n)

	if uint64(n) < responsePayloadLength {
		opLogger.Trace("not all data consumed yet, receiving remainder...",
			"expected-payload-length", responsePayloadLength,
			"consumed-payload-length", n)
		for {
			buf, err := c.recv()
			if err != nil {
				opLogger.Error("response payload read error",
					"expected-payload-length", responsePayloadLength,
					"read-payload-length", n,
					"reason", err)
				return &errors.ReceiveError{
					Cause: fmt.Errorf("response payload read error: read %d bytes vs expected %d, reason: %s",
						n, responsePayloadLength, err.Error()),
				}
			}

			n = n + buf.Len()
			responsePayloadBuf = append(responsePayloadBuf, buf.Bytes()...)

			if uint64(n) == responsePayloadLength {
				opLogger.Trace("consumed all expected data",
					"expected-payload-length", responsePayloadLength,
					"consumed-payload-length", n)
				break
			}

			if uint64(n) > responsePayloadLength {
				opLogger.Error("consumed too much data",
					"expected-payload-length", responsePayloadLength,
					"consumed-payload-length", n)
				return &errors.ReceiveError{
					Cause: fmt.Errorf("response payload too long: read %d bytes vs expected %d",
						n, responsePayloadLength),
				}
			}
		}
	}

	if *responseHeader.IsError {
		errorResponse := &ybApi.ErrorStatusPB{}
		errorUnmarshalErr := utils.DeserializeProto(responsePayloadBuf, errorResponse)
		if errorUnmarshalErr != nil {
			return &errors.UnprocessableResponseError{
				Cause:           errorUnmarshalErr,
				ConsumedPayload: responseHeaderBuf,
			}
		}
		return &errors.ServiceRPCError{
			Cause: errorResponse,
		}
	}

	protoErr2 := utils.DeserializeProto(responsePayloadBuf, m)
	if protoErr2 != nil {
		return &errors.UnprocessableResponseError{
			Cause:           protoErr2,
			ConsumedPayload: responseHeaderBuf,
		}
	}

	return nil
}

func (c *defaultSingleNodeClient) executeOp(payload, result protoreflect.ProtoMessage) error {

	svcInfo := c.svcRegistry.Get(payload)
	if svcInfo == nil {
		c.metricsCallback.ClientError()
		c.metricsCallback.ClientMessageSendFailure()
		return &errors.ProtoServiceError{
			ProtoType: payload.ProtoReflect().Descriptor().FullName(),
		}
	}

	requestHeader := &ybApi.RequestHeader{
		CallId:        utils.PInt32(int32(c.callID())),
		RemoteMethod:  svcInfo.ToRemoteMethodPB(),
		TimeoutMillis: utils.PUint32(c.originalConfig.OpTimeout),
	}

	b := bytes.NewBuffer([]byte{})
	if err := utils.WriteMessages(b, requestHeader, payload); err != nil {
		c.metricsCallback.ClientError()
		c.metricsCallback.ClientMessageSendFailure()
		return &errors.PayloadWriteError{
			Cause:   err,
			Header:  requestHeader,
			Payload: payload,
		}
	}
	if err := c.send(b); err != nil {
		c.metricsCallback.ClientError()
		c.metricsCallback.ClientMessageSendFailure()
		return &errors.SendError{Cause: err}
	}
	buffer, err := c.recv()
	if err != nil {
		c.metricsCallback.ClientError()
		c.metricsCallback.ClientMessageSendFailure()
		return &errors.ReceiveError{Cause: err}
	}
	readResponseErr := c.readResponseInto(buffer, result)
	if readResponseErr != nil {
		c.metricsCallback.ClientError()
		c.metricsCallback.ClientMessageSendFailure()
		return readResponseErr
	}
	c.metricsCallback.ClientMessageSendSuccess()
	return nil
}

func (c *defaultSingleNodeClient) withLogger(logger hclog.Logger) *defaultSingleNodeClient {
	c.logger = logger
	return c
}

func (c *defaultSingleNodeClient) withMetricsCallback(callback metrics.Callback) *defaultSingleNodeClient {
	c.metricsCallback = callback
	return c
}
