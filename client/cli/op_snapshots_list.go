package cli

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultYBCliClient) SnapshotsList(opConfig *configs.OpSnapshotListConfig) (*ybApi.ListSnapshotsResponsePB, error) {
	payload := &ybApi.ListSnapshotsRequestPB{
		ListDeletedSnapshots: &opConfig.ListDeletedSnapshots,
		PrepareForBackup:     &opConfig.PrepareForBackup,
	}
	if len(opConfig.SnapshotID) > 0 {
		payload.SnapshotId = []byte(opConfig.SnapshotID)
	}
	responsePayload := &ybApi.ListSnapshotsResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
