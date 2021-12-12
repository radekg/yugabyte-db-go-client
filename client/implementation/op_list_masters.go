package implementation

import (
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// ListMasters returns a list of masters or an error if call failed.
func (c *defaultYBCliClient) ListMasters() (*ybApi.ListMastersResponsePB, error) {
	payload := &ybApi.ListMastersRequestPB{}
	responsePayload := &ybApi.ListMastersResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
