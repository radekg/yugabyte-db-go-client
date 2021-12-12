package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// ConsensusError is an error representation of the ConsensusErrorPB.
type ConsensusError struct {
	genericError
	Code *ybApi.ConsensusErrorPB_Code
}

func (e *ConsensusError) Error() string {
	code := int32(0)
	codeName := ybApi.ConsensusErrorPB_Code_name[code]
	if e.Code != nil {
		if v, ok := ybApi.ConsensusErrorPB_Code_name[int32(*e.Code)]; ok {
			codeName = v
			code = int32(*e.Code)
		}
	}
	return fmt.Sprintf("consensus rpc error: code: %d (%s), %s",
		code, codeName, e.statusToString())
}

// NewConsensusError converts ConsensusErrorPB into an error.
func NewConsensusError(input *ybApi.ConsensusErrorPB) error {
	if input == nil {
		return nil
	}
	return &ConsensusError{
		Code: input.Code,
		genericError: genericError{
			Status: input.Status,
		},
	}
}
