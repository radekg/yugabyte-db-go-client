package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// MasterError is an error representation of the MasterErrorPB.
type MasterError struct {
	genericError
	Code *ybApi.MasterErrorPB_Code
}

func (e *MasterError) Error() string {
	code := int32(1)
	codeName := ybApi.MasterErrorPB_Code_name[code]
	if e.Code != nil {
		if v, ok := ybApi.MasterErrorPB_Code_name[int32(*e.Code)]; ok {
			codeName = v
			code = int32(*e.Code)
		}
	}
	return fmt.Sprintf("master rpc error: code: %d (%s), %s",
		code, codeName, e.statusToString())
}

// NewMasterError converts MasterErrorPB into an error.
func NewMasterError(input *ybApi.MasterErrorPB) error {
	if input == nil {
		return nil
	}
	return &MasterError{
		Code: input.Code,
		genericError: genericError{
			Status: input.Status,
		},
	}
}
