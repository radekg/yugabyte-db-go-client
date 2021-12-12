package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// PerOpError is an error representation of the PerOpErrorPB.
type PerOpError struct {
	genericError
	ID *ybApi.OpIdPB
}

func (e *PerOpError) Error() string {
	errString := "per op rpc error"
	if e.ID != nil {
		if e.ID.Index != nil {
			errString = fmt.Sprintf("%s, index: %d", errString, *e.ID.Index)
		}
		if e.ID.Term != nil {
			errString = fmt.Sprintf("%s, term: %d", errString, *e.ID.Term)
		}
	}
	if e.Status != nil {
		errString = fmt.Sprintf("%s, %s", errString, e.statusToString())
	}
	return errString
}

// NewPerOpError converts PerOpErrorPB into an error.
func NewPerOpError(input *ybApi.PerOpErrorPB) error {
	if input == nil {
		return nil
	}
	return &PerOpError{
		ID: input.Id,
		genericError: genericError{
			Status: input.Status,
		},
	}
}
