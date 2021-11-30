package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Restore snapshot.
func (c *defaultYBCliClient) SnapshotsRestore(opConfig *configs.OpSnapshotRestoreConfig) (*ybApi.RestoreSnapshotResponsePB, error) {

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

	payload := &ybApi.RestoreSnapshotRequestPB{
		SnapshotId: protoSnapshotID,
	}
	if opConfig.RestoreHt > 0 {
		payload.RestoreHt = &opConfig.RestoreHt
	}

	responsePayload := &ybApi.RestoreSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
