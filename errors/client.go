package errors

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

const (
	// ErrorMessageConnected is an error message.
	ErrorMessageConnected = "client: connected"
	// ErrorMessageConnecting is an error message.
	ErrorMessageConnecting = "client: connecting"
	// ErrorMessageLeaderWaitTimeout is an error message.
	ErrorMessageLeaderWaitTimeout = "client: leader wait timed out"
	// ErrorMessageNoClient is an error message.
	ErrorMessageNoClient = "client: no client"
	// ErrorMessageNoLeader is an error message.
	ErrorMessageNoLeader = "client: no leader"
	// ErrorMessageNotConnected is an error message.
	ErrorMessageNotConnected = "client: not connected"
	// ErrorMessagePayloadError is an error message.
	ErrorMessagePayloadError = "client: payload error"
	// ErrorMessageProtocolConnectionHeader is an error message.
	ErrorMessageProtocolConnectionHeader = "client: protocol connection header error"
	// ErrorMessageProtoServiceError is an error message.
	ErrorMessageProtoServiceError = "client: proto service error"
	// ErrorMessageReconnectFailed is an error message.
	ErrorMessageReconnectFailed = "client: reconnect failed"
	// ErrorMessageReconnectRequired is an error message.
	ErrorMessageReconnectRequired = "client: reconnect required"
	// ErrorMessageSendReceiveFailed is an error message.
	ErrorMessageSendReceiveFailed = "client: send/receive failed"
	// ErrorMessageUnprocessableResponse is an error message.
	ErrorMessageUnprocessableResponse = "client: unprocessable response"
)

// AbstractMasterErrorResponse isn't an error. It represents an RPC response
// returning an instance of the MasterErrorPB error.
// This type is used to check if the client needs to reconnect and retry a call
// in case of a call not being issued against a leader master.
type AbstractMasterErrorResponse interface {
	GetError() *ybApi.MasterErrorPB
}

// NoLeaderError represents a client without a leader error.
type NoLeaderError struct{}

func (e *NoLeaderError) Error() string {
	return ErrorMessageNoLeader
}

// PayloadWriteError happens when the client cannot serialize the header
// or the payload. This is a non-recoverable error.
type PayloadWriteError struct {
	Cause   error
	Header  *ybApi.RequestHeader
	Payload protoreflect.ProtoMessage
}

func (e *PayloadWriteError) Error() string {
	return fmt.Sprintf("%s: %s", ErrorMessagePayloadError, e.Cause.Error())
}

// ProtocolConnectionHeaderWriteError is an error returned when the initial
// connect header could not be written.
type ProtocolConnectionHeaderWriteError struct {
	Cause error
}

func (e *ProtocolConnectionHeaderWriteError) Error() string {
	return fmt.Sprintf("%s: %s", ErrorMessageProtocolConnectionHeader, e.Cause.Error())
}

// ProtocolConnectionHeaderWriteIncompleteError is an error returned when the initial
// connect header could not be fully written.
type ProtocolConnectionHeaderWriteIncompleteError struct {
	Header   []byte
	Expected int
	Written  int
}

func (e *ProtocolConnectionHeaderWriteIncompleteError) Error() string {
	return fmt.Sprintf("%s: written %d bytes vs expected %d bytes", ErrorMessageProtocolConnectionHeader, e.Written, e.Expected)
}

// ProtoServiceError happens when the service registry cannot identify
// a service for a protobuf type. This is a non-recoverable error.
type ProtoServiceError struct {
	ProtoType protoreflect.FullName
}

func (e *ProtoServiceError) Error() string {
	return fmt.Sprintf("%s: %s", ErrorMessageProtoServiceError, e.ProtoType)
}

// RequiresReconnectError is an error indicating a need to reconnect.
type RequiresReconnectError struct {
	Cause error
}

func (e *RequiresReconnectError) Error() string {
	return fmt.Sprintf("%s: no service for type '%s'", ErrorMessageReconnectRequired, e.Cause.Error())
}

// SendReceiveError is returned when the client is unable to
// send the payload or receive from the server.
type SendReceiveError struct {
	Cause error
}

func (e *SendReceiveError) Error() string {
	return fmt.Sprintf("%s: %s", ErrorMessageSendReceiveFailed, e.Cause.Error())
}

// UnprocessableResponseError represents a client error where a fully read response
// cannot be deserialized as a protobuf message.
// This error usually implies that a retry is required.
type UnprocessableResponseError struct {
	Cause           error
	ConsumedPayload []byte
}

func (e *UnprocessableResponseError) Error() string {
	return fmt.Sprintf("%s: %s", ErrorMessageUnprocessableResponse, e.Cause.Error())
}
