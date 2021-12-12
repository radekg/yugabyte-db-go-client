package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// MasterError is an error representation of the MasterErrorPB.
type MasterError struct {
	Code   *ybApi.MasterErrorPB_Code
	Status *ybApi.AppStatusPB
}

func (e *MasterError) statusToString() string {
	if e.Status == nil {
		return "status: <unknown>"
	}
	code := int32(999)
	status := ybApi.AppStatusPB_ErrorCode_name[code]
	if e.Status.Code != nil {
		if v, ok := ybApi.AppStatusPB_ErrorCode_name[int32(*e.Status.Code)]; ok {
			status = v
			code = int32(*e.Status.Code)
		}
	}
	errString := fmt.Sprintf("status: %d (%s)", code, status)
	if e.Status.Message != nil {
		errString = fmt.Sprintf("%s\n\tmessage: %s", errString, *e.Status.Message)
	}
	if e.Status.SourceFile != nil {
		errString = fmt.Sprintf("%s\n\tsource: %s", errString, *e.Status.SourceFile)
		if e.Status.SourceLine != nil {
			errString = fmt.Sprintf("%s@%d", errString, *e.Status.SourceLine)
		}
	}
	return errString
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
	return fmt.Sprintf("rpc error: code: %d (%s), %s",
		code, codeName, e.statusToString())
}

// ToError converts MasterErrorPB into an error.
func ToError(genericError *ybApi.MasterErrorPB) error {
	if genericError == nil {
		return nil
	}
	return &MasterError{
		Code:   genericError.Code,
		Status: genericError.Status,
	}
}
