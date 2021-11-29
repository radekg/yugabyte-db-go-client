package cli

import (
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

	if opConfig.Keyspace != "" {
		parsedKeyspace := parseKeyspace(opConfig.Keyspace)
		payload.Namespace = parsedKeyspace.toProtoKeyspace()
	}

	if len(opConfig.RelationType) > 0 {
		payload.RelationTypeFilter = []ybApi.RelationType{}
		for _, relation := range opConfig.RelationType {
			if relationType, ok := mapRelationTypeFilter(relation); ok {
				payload.RelationTypeFilter = append(payload.RelationTypeFilter, relationType)
			}
		}
	}

	responsePayload := &ybApi.ListTablesResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
