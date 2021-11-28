package cli

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
// Keyspace handling

type parsedKeyspace struct {
	YQLDatabaseType string
	Keyspace        string
}

func (pk *parsedKeyspace) toProtoKeyspace() *ybApi.NamespaceIdentifierPB {
	if yqlDatabaseType, ok := mapYQLDatabaseType(pk.YQLDatabaseType); ok {
		return &ybApi.NamespaceIdentifierPB{
			Name:         &pk.Keyspace,
			DatabaseType: pYQLDatabase(yqlDatabaseType),
		}
	}
	return &ybApi.NamespaceIdentifierPB{
		Name:         &pk.Keyspace,
		DatabaseType: pYQLDatabase(ybApi.YQLDatabase_YQL_DATABASE_CQL),
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
	case "index":
		return ybApi.RelationType_INDEX_TABLE_RELATION, true
	default:
		return -1, false
	}
}

func pYQLDatabase(input ybApi.YQLDatabase) *ybApi.YQLDatabase {
	return &input
}
