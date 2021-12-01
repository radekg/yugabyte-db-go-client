package implementation

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

func (c *defaultYBCliClient) defaultServerClockResolver() (uint64, error) {
	serverClock, err := c.ServerClock()
	if err != nil {
		return 0, err
	}
	if serverClock.HybridTime == nil {
		return 0, fmt.Errorf("no hybrid time in server clock response")
	}
	return *serverClock.HybridTime, nil
}

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
