package implementation

import (
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Gets the server clock value.
// Returned server time is represented in microseconds.
func (c *defaultYBCliClient) ServerClock() (*ybApi.ServerClockResponsePB, error) {
	payload := &ybApi.ServerClockRequestPB{}
	responsePayload := &ybApi.ServerClockResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
