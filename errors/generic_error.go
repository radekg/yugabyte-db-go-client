package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

type genericError struct {
	Status *ybApi.AppStatusPB
}

func (e *genericError) statusToString() string {
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
