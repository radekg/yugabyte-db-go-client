package implementation

import (
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// IsTabletServerReady checks if a given tablet server is ready or returns an error.
func (c *defaultYBCliClient) IsTabletServerReady() (*ybApi.IsTabletServerReadyResponsePB, error) {
	payload := &ybApi.IsTabletServerReadyRequestPB{}
	responsePayload := &ybApi.IsTabletServerReadyResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewTabletServerError(responsePayload.Error)
}
