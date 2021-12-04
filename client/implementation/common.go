package implementation

import (
	"fmt"
	"strings"

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

// ==
// Table lookup by name

func (c *defaultYBCliClient) lookupTableIDsByNames(keyspace string, names []string) (map[string][]byte, error) {

	parsedKeyspace := parseKeyspace(keyspace)
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

	results := map[string][]byte{}
	for _, name := range names {
		var found bool
		for _, tableInfo := range responseListTablesPayload.Tables {
			namespace := *tableInfo.Namespace
			if *namespace.Name == parsedKeyspace.Keyspace && *tableInfo.Name == name {
				results[name] = tableInfo.Id
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("table %s.%s not found", keyspace, name)
		}
	}

	return results, nil
}

func (c *defaultYBCliClient) lookupTableByName(keyspace, name string) (*ybApi.GetTableSchemaResponsePB, error) {
	parsedKeyspace := parseKeyspace(keyspace)
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
			if *namespace.Name == parsedKeyspace.Keyspace && *tableInfo.Name == name {
				return c.getTableSchemaByUUID(tableInfo.Id)
			}
		}
	}

	return nil, fmt.Errorf("table %s.%s not found", keyspace, name)
}

// ==
// Keyspace handling

type parsedKeyspace struct {
	YQLDatabaseType string
	Keyspace        string
}

func (pk *parsedKeyspace) toProtoKeyspace() *ybApi.NamespaceIdentifierPB {
	if yqlDatabaseType, ok := mapYQLDatabaseType(pk.YQLDatabaseType); ok {
		return &ybApi.NamespaceIdentifierPB{
			Name:         &pk.Keyspace,
			DatabaseType: PYQLDatabase(yqlDatabaseType),
		}
	}
	return &ybApi.NamespaceIdentifierPB{
		Name:         &pk.Keyspace,
		DatabaseType: PYQLDatabase(ybApi.YQLDatabase_YQL_DATABASE_CQL),
	}
}

func parseKeyspace(input string) *parsedKeyspace {
	if strings.HasPrefix(input, "ycql.") {
		return &parsedKeyspace{
			YQLDatabaseType: "ycql",
			Keyspace:        strings.TrimPrefix(input, "ycql."),
		}
	}
	if strings.HasPrefix(input, "ysql.") {
		return &parsedKeyspace{
			YQLDatabaseType: "ysql",
			Keyspace:        strings.TrimPrefix(input, "ysql."),
		}
	}
	if strings.HasPrefix(input, "yedis.") {
		return &parsedKeyspace{
			YQLDatabaseType: "yedis",
			Keyspace:        strings.TrimPrefix(input, "yedis."),
		}
	}
	return &parsedKeyspace{
		YQLDatabaseType: "ycql",
		Keyspace:        input,
	}
}

func mapYQLDatabaseType(input string) (ybApi.YQLDatabase, bool) {
	switch input {
	case "ycql":
		return ybApi.YQLDatabase_YQL_DATABASE_CQL, true
	case "ysql":
		return ybApi.YQLDatabase_YQL_DATABASE_PGSQL, true
	case "yedis":
		return ybApi.YQLDatabase_YQL_DATABASE_REDIS, true
	default:
		return -1, false
	}
}

func mapRelationTypeFilter(input string) (ybApi.RelationType, bool) {
	switch input {
	case "system_table":
		return ybApi.RelationType_SYSTEM_TABLE_RELATION, true
	case "user_table":
		return ybApi.RelationType_USER_TABLE_RELATION, true
	case "index_table":
		return ybApi.RelationType_INDEX_TABLE_RELATION, true
	default:
		return -1, false
	}
}

// PYQLDatabase returns a pointer to the given input YQLDatabase.
func PYQLDatabase(input ybApi.YQLDatabase) *ybApi.YQLDatabase {
	return &input
}
