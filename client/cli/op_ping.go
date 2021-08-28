package cli

import ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"

// Ping pings a certain YB server.
func (c *defaultYBCliClient) Ping() (*ybApi.PingResponsePB, error) {
	payload := &ybApi.PingRequestPB{}
	responsePayload := &ybApi.PingResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
