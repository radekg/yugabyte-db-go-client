package implementation

import (
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// GetIsLoadBalancerIdle finds out if the load balancer is idle.
func (c *defaultYBCliClient) GetIsLoadBalancerIdle() (*ybApi.IsLoadBalancerIdleResponsePB, error) {
	payload := &ybApi.IsLoadBalancerIdleRequestPB{}
	responsePayload := &ybApi.IsLoadBalancerIdleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
