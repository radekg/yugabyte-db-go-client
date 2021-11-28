package cli

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// DescribeTable returns info on a table in this database.
func (c *defaultYBCliClient) GetTabletsForTable(opConfig *configs.OpGetTableLocationsConfig) (*ybApi.GetTableLocationsResponsePB, error) {

	if opConfig.UUID != "" {
		// we can short circuit everything below:
		return c.getTabletsForTableByUUID([]byte(opConfig.UUID), opConfig)
	}

	payloadListTables := &ybApi.ListTablesRequestPB{}
	responseListTablesPayload := &ybApi.ListTablesResponsePB{}
	if err := c.connectedClient.Execute(payloadListTables, responseListTablesPayload); err != nil {
		return nil, err
	}
	if err := responseListTablesPayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}

	for _, tableInfo := range responseListTablesPayload.Tables {
		if tableInfo.Namespace != nil {
			namespace := *tableInfo.Namespace
			if *namespace.Name == opConfig.Keyspace && *tableInfo.Name == opConfig.Name {
				return c.getTabletsForTableByUUID(tableInfo.Id, opConfig)
			}
		}
	}

	return nil, fmt.Errorf("not found")
}
