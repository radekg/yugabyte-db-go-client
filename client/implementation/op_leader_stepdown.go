package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
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
	return responsePayload, clientErrors.NewTabletServerError(responsePayload.Error)
}
