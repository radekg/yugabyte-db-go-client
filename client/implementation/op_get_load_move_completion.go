package implementation

import (
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// GetLoadMoveCompletion gets the completion percentage of tablet load move from blacklisted servers.
func (c *defaultYBCliClient) GetLoadMoveCompletion() (*ybApi.GetLoadMovePercentResponsePB, error) {
	payload := &ybApi.GetLoadMovePercentRequestPB{}
	responsePayload := &ybApi.GetLoadMovePercentResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
