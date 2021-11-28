package cli

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// DescribeTable returns info on a table in this database.
func (c *defaultYBCliClient) DescribeTable(opConfig *configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error) {

	if opConfig.UUID != "" {
		// we can short circuit everything below:
		return c.getTableSchemaByUUID([]byte(opConfig.UUID))
	}

	parsedKeyspace := parseKeyspace(opConfig.Keyspace)
	payloadListTables := &ybApi.ListTablesRequestPB{
		Namespace: parsedKeyspace.toProtoKeyspace(),
	}
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
			if *namespace.Name == parsedKeyspace.Keyspace && *tableInfo.Name == opConfig.Name {
				return c.getTableSchemaByUUID(tableInfo.Id)
			}
		}
	}

	return nil, fmt.Errorf("not found")
}
