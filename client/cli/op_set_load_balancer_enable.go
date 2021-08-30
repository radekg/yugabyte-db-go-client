package cli

import (
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Ping pings a certain YB server.
func (c *defaultYBCliClient) SetLoadBalancerEnable(enable bool) (*ybApi.ChangeLoadBalancerStateResponsePB, error) {
	payload := &ybApi.ChangeLoadBalancerStateRequestPB{
		IsEnabled: utils.PBool(enable),
	}
	responsePayload := &ybApi.ChangeLoadBalancerStateResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
