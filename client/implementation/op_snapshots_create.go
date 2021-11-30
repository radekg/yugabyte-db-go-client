package implementation

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Create a snapshot.
func (c *defaultYBCliClient) SnapshotsCreate(opConfig *configs.OpSnapshotCreateConfig) (*ybApi.CreateSnapshotResponsePB, error) {

	if len(opConfig.ScheduleID) > 0 {
		// short circuit
		payload := &ybApi.CreateSnapshotRequestPB{
			ScheduleId: opConfig.ScheduleID,
		}
		responsePayload := &ybApi.CreateSnapshotResponsePB{}
		if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
			return nil, err
		}
		return responsePayload, nil
	}

	parsedKeyspace := parseKeyspace(opConfig.Keyspace)

	switch parsedKeyspace.YQLDatabaseType {
	case "ycql":
		return createSnapshotYCQL(c, opConfig, parsedKeyspace)
	case "ysql":
		return createSnapshotYSQL(c, opConfig, parsedKeyspace)
	default:
		return nil, fmt.Errorf("unsupported snapshot keyspace type: %s", parsedKeyspace.YQLDatabaseType)
	}
}

func createSnapshotYCQL(c *defaultYBCliClient, opConfig *configs.OpSnapshotCreateConfig, ns *parsedKeyspace) (*ybApi.CreateSnapshotResponsePB, error) {
	tableIdentifiers := []*ybApi.TableIdentifierPB{}
	for _, tableUUID := range opConfig.TableUUIDs {
		tableIdentifiers = append(tableIdentifiers, &ybApi.TableIdentifierPB{
			TableId: []byte(tableUUID),
		})
	}

	if len(opConfig.TableNames) > 0 {
		mappedIDs, err := c.lookupTableIDsByNames(opConfig.Keyspace, opConfig.TableNames)
		if err != nil {
			return nil, err
		}
		for _, id := range mappedIDs {
			tableIdentifiers = append(tableIdentifiers, &ybApi.TableIdentifierPB{
				TableId: id,
			})
		}
	}

	if len(tableIdentifiers) == 0 {
		tables, err := c.ListTables(&configs.OpListTablesConfig{
			Keyspace: opConfig.Keyspace,
		})
		if err != nil {
			return nil, err
		}
		for _, tableInfo := range tables.Tables {
			tableIdentifiers = append(tableIdentifiers, &ybApi.TableIdentifierPB{
				TableId: tableInfo.Id,
			})
		}
	}

	payload := &ybApi.CreateSnapshotRequestPB{
		Tables: tableIdentifiers,
		AddIndexes: func() *bool {
			v := true // https://github.com/yugabyte/yugabyte-db/blob/d4d5688147734d1a36bbe58430f35ba4db2770f1/ent/src/yb/tools/yb-admin_cli_ent.cc#L119
			return &v
		}(),
		TransactionAware: func() *bool {
			v := true // https://github.com/yugabyte/yugabyte-db/blob/d4d5688147734d1a36bbe58430f35ba4db2770f1/ent/src/yb/tools/yb-admin_client_ent.cc#L247
			return &v
		}(),
	}

	responsePayload := &ybApi.CreateSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}

func createSnapshotYSQL(c *defaultYBCliClient, opConfig *configs.OpSnapshotCreateConfig, ns *parsedKeyspace) (*ybApi.CreateSnapshotResponsePB, error) {
	tablesPayload, err := c.ListTables(&configs.OpListTablesConfig{
		Keyspace:            opConfig.Keyspace,
		ExcludeSystemTables: true,                                  // https://github.com/yugabyte/yugabyte-db/blob/d4d5688147734d1a36bbe58430f35ba4db2770f1/ent/src/yb/tools/yb-admin_client_ent.cc#L262
		RelationType:        []string{"user_table", "index_table"}, // https://github.com/yugabyte/yugabyte-db/blob/d4d5688147734d1a36bbe58430f35ba4db2770f1/ent/src/yb/tools/yb-admin_client_ent.cc#L263-L264
	})
	if err != nil {
		return nil, err
	}
	tableIdentifiers := []*ybApi.TableIdentifierPB{}
	for _, tableInfo := range tablesPayload.Tables {
		tableIdentifiers = append(tableIdentifiers, &ybApi.TableIdentifierPB{
			TableId: tableInfo.Id,
		})
	}

	payload := &ybApi.CreateSnapshotRequestPB{
		Tables: tableIdentifiers,
		AddIndexes: func() *bool {
			v := false
			return &v
		}(),
	}

	responsePayload := &ybApi.CreateSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
