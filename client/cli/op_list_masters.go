package cli

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// ListMasters returns a list of masters or an error if call failed.
func (c *defaultYBCliClient) ListMasters() (*ybApi.ListMastersResponsePB, error) {
	payload := &ybApi.ListMastersRequestPB{}
	responsePayload := &ybApi.ListMastersResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}
