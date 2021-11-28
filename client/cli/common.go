package cli

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

func (c *defaultYBCliClient) getTableSchemaByUUID(tableID []byte) (*ybApi.GetTableSchemaResponsePB, error) {
	payload := &ybApi.GetTableSchemaRequestPB{
		Table: &ybApi.TableIdentifierPB{
			TableId: tableID,
		},
	}
	responsePayload := &ybApi.GetTableSchemaResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

func (c *defaultYBCliClient) getTabletsForTableByUUID(tableID []byte, opConfig *configs.OpGetTableLocationsConfig) (*ybApi.GetTableLocationsResponsePB, error) {
	payload := &ybApi.GetTableLocationsRequestPB{
		Table: &ybApi.TableIdentifierPB{
			TableId: tableID,
		},
		PartitionKeyStart:     opConfig.PartitionKeyStart,
		PartitionKeyEnd:       opConfig.PartitionKeyEnd,
		MaxReturnedLocations:  utils.PUint32(opConfig.MaxReturnedLocations),
		RequireTabletsRunning: utils.PBool(opConfig.RequireTabletsRunning),
	}
	responsePayload := &ybApi.GetTableLocationsResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}
