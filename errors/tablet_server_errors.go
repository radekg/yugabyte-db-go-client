package errors

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// TabletServerError is an error representation of the TabletServerErrorPB.
type TabletServerError struct {
	genericError
	Code *ybApi.TabletServerErrorPB_Code
}

func (e *TabletServerError) Error() string {
	code := int32(1)
	codeName := ybApi.TabletServerErrorPB_Code_name[code]
	if e.Code != nil {
		if v, ok := ybApi.TabletServerErrorPB_Code_name[int32(*e.Code)]; ok {
			codeName = v
			code = int32(*e.Code)
		}
	}
	return fmt.Sprintf("tablet server rpc error: code: %d (%s), %s",
		code, codeName, e.statusToString())
}

// NewTabletServerError converts TabletServerErrorPB into an error.
func NewTabletServerError(input *ybApi.TabletServerErrorPB) error {
	if input == nil {
		return nil
	}
	return &TabletServerError{
		Code: input.Code,
		genericError: genericError{
			Status: input.Status,
		},
	}
}
