package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultYBCliClient) SnapshotsList(opConfig *configs.OpSnapshotListConfig) (*ybApi.ListSnapshotsResponsePB, error) {
	payload := &ybApi.ListSnapshotsRequestPB{
		ListDeletedSnapshots: &opConfig.ListDeletedSnapshots,
		PrepareForBackup:     &opConfig.PrepareForBackup,
	}
	if len(opConfig.SnapshotID) > 0 {

		givenSnapshotID, err := utils.DecodeAsYugabyteID(opConfig.SnapshotID, opConfig.Base64Encoded)
		if err != nil {
			c.logger.Error("failed fetching normalized snapshot id",
				"given-value", opConfig.SnapshotID,
				"reason", err)
			return nil, err
		}

		protoSnapshotID, err := utils.StringUUIDToProtoYugabyteID(givenSnapshotID)
		if err != nil {
			return nil, err
		}

		payload.SnapshotId = protoSnapshotID
	}
	responsePayload := &ybApi.ListSnapshotsResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
