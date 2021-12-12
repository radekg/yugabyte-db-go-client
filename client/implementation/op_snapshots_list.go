package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultYBCliClient) SnapshotsList(opConfig *configs.OpSnapshotListConfig) (*ybApi.ListSnapshotsResponsePB, error) {
	payload := &ybApi.ListSnapshotsRequestPB{
		ListDeletedSnapshots: &opConfig.ListDeletedSnapshots,
		PrepareForBackup:     &opConfig.PrepareForBackup,
	}
	if len(opConfig.SnapshotID) > 0 {

		ybDbID, err := ybdbid.TryParseFromString(opConfig.SnapshotID)
		if err != nil {
			c.logger.Error("given snapshot id is not valid",
				"original-value", opConfig.SnapshotID,
				"reason", err)
			return nil, err
		}

		payload.SnapshotId = ybDbID.Bytes()
	}
	responsePayload := &ybApi.ListSnapshotsResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
