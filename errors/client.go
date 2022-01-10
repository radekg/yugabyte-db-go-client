package errors

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

var (
	errNoLeader = fmt.Errorf("client: no leader")
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
	return errNoLeader.Error()
}

// RequiresReconnectError is an error indicating a need to reconnect.
type RequiresReconnectError struct {
	Cause error
}

func (e *RequiresReconnectError) Error() string {
	return multierror.Append(fmt.Errorf("client: requires reconnect"), e.Cause).Error()
}

// UnprocessableResponseError represents a client error where a fully read response
// cannot be deserialized as a protobuf message.
// This error usually implies that a retry is required.
type UnprocessableResponseError struct {
	Cause           error
	ConsumedPayload []byte
}

func (e *UnprocessableResponseError) Error() string {
	return multierror.Append(fmt.Errorf("client: unprocessable response"), e.Cause).Error()
}
