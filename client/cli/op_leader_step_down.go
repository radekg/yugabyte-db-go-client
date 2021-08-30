package cli

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// ListMasters returns a list of masters or an error if call failed.
func (c *defaultYBCliClient) LeaderStepDown(opConfig *configs.OpLeaderStepDownConfig) (*ybApi.LeaderStepDownResponsePB, error) {

	payload := &ybApi.LeaderStepDownRequestPB{
		DestUuid:                  []byte(opConfig.DestUUID),
		TabletId:                  []byte(opConfig.TabletID),
		DisableGracefulTransition: utils.PBool(opConfig.DisableGracefulTansition),
		NewLeaderUuid:             []byte(opConfig.NewLeaderUUID),
	}
	responsePayload := &ybApi.LeaderStepDownResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}

	return responsePayload, nil
	//return responsePayload, nil
}
