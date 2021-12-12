package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// RemoteBootstrapError is an error representation of the RemoteBootstrapErrorPB.
type RemoteBootstrapError struct {
	genericError
	Code *ybApi.RemoteBootstrapErrorPB_Code
}

func (e *RemoteBootstrapError) Error() string {
	code := int32(1)
	codeName := ybApi.RemoteBootstrapErrorPB_Code_name[code]
	if e.Code != nil {
		if v, ok := ybApi.RemoteBootstrapErrorPB_Code_name[int32(*e.Code)]; ok {
			codeName = v
			code = int32(*e.Code)
		}
	}
	return fmt.Sprintf("remote bootstrap rpc error: code: %d (%s), %s",
		code, codeName, e.statusToString())
}

// NewRemoteBootstrapError converts ConsensusErrorPB into an error.
func NewRemoteBootstrapError(input *ybApi.RemoteBootstrapErrorPB) error {
	if input == nil {
		return nil
	}
	return &RemoteBootstrapError{
		Code: input.Code,
		genericError: genericError{
			Status: input.Status,
		},
	}
}
