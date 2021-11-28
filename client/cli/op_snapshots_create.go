package cli

import (
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
		// load all tables:
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
		Tables:           tableIdentifiers,
		TransactionAware: &opConfig.TransactionAware,
		AddIndexes:       &opConfig.AddIndexes,
		Imported:         &opConfig.Imported,
	}

	responsePayload := &ybApi.CreateSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
