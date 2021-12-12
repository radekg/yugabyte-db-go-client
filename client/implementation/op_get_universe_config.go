package implementation

import (
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// GetUniverseConfig get the placement info and blacklist info of the universe.
func (c *defaultYBCliClient) GetUniverseConfig() (*ybApi.GetMasterClusterConfigResponsePB, error) {
	payload := &ybApi.GetMasterClusterConfigRequestPB{}
	responsePayload := &ybApi.GetMasterClusterConfigResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
