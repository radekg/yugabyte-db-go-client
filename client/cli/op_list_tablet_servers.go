package cli

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// ListTabletServers returns a list of tablet servers or an error if call failed.
func (c *defaultYBCliClient) ListTabletServers(opConfig *configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error) {
	payload := &ybApi.ListTabletServersRequestPB{
		PrimaryOnly: utils.PBool(opConfig.PrimaryOnly),
	}
	responsePayload := &ybApi.ListTabletServersResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
