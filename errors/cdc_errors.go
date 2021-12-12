package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// CDCError is an error representation of the CDCErrorPB.
type CDCError struct {
	genericError
	Code *ybApi.CDCErrorPB_Code
}

func (e *CDCError) Error() string {
	code := int32(1)
	codeName := ybApi.CDCErrorPB_Code_name[code]
	if e.Code != nil {
		if v, ok := ybApi.CDCErrorPB_Code_name[int32(*e.Code)]; ok {
			codeName = v
			code = int32(*e.Code)
		}
	}
	return fmt.Sprintf("CDC rpc error: code: %d (%s), %s",
		code, codeName, e.statusToString())
}

// NewCDCError converts MasterErrorPB into an error.
func NewCDCError(input *ybApi.CDCErrorPB) error {
	if input == nil {
		return nil
	}
	return &CDCError{
		Code: input.Code,
		genericError: genericError{
			Status: input.Status,
		},
	}
}
