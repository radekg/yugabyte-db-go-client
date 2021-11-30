package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Set load balancer state.
func (c *defaultYBCliClient) SetLoadBalancerState(enable bool) (*ybApi.ChangeLoadBalancerStateResponsePB, error) {
	payload := &ybApi.ChangeLoadBalancerStateRequestPB{
		IsEnabled: utils.PBool(enable),
	}
	responsePayload := &ybApi.ChangeLoadBalancerStateResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
