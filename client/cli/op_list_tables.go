package cli

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// ListTables returns a List all the tables in a database.
func (c *defaultYBCliClient) ListTables(opConfig *configs.OpListTablesConfig) (*ybApi.ListTablesResponsePB, error) {
	payload := &ybApi.ListTablesRequestPB{
		ExcludeSystemTables: &opConfig.ExcludeSystemTables,
		IncludeNotRunning:   &opConfig.IncludeNotRunning,
	}
	if opConfig.NameFilter != "" {
		payload.NameFilter = &opConfig.NameFilter
	}
	if opConfig.NamespaceName != "" || opConfig.NamespaceType != "" {
		payload.Namespace = &ybApi.NamespaceIdentifierPB{
			//DatabaseType: pYQLDatabase(ybApi.YQLDatabase(ybApi.YQLDatabase_YQL_DATABASE_UNKNOWN)),
		}
		if opConfig.NamespaceName != "" {
			payload.Namespace.Name = &opConfig.NamespaceName
		}
		if opConfig.NamespaceType != "" {
			if yqlDatabaseType, ok := mapYQLDatabaseType(opConfig.NamespaceType); ok {
				payload.Namespace.DatabaseType = pYQLDatabase(ybApi.YQLDatabase(yqlDatabaseType))
			}
		}
	} /*
		if len(opConfig.RelationType) > 0 {
			payload.RelationTypeFilter = []ybApi.RelationType{}
			for _, relation := range opConfig.RelationType {
				if relationType, ok := mapRelationTypeFilter(relation); ok {
					payload.RelationTypeFilter = append(payload.RelationTypeFilter, relationType)
				}
			}
		}*/

	responsePayload := &ybApi.ListTablesResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

func mapYQLDatabaseType(input string) (ybApi.YQLDatabase, bool) {
	switch input {
	case "cql":
		return ybApi.YQLDatabase_YQL_DATABASE_CQL, true
	case "pgsql":
		return ybApi.YQLDatabase_YQL_DATABASE_PGSQL, true
	case "redis":
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
