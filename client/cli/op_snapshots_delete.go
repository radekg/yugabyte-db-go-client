package cli

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultYBCliClient) SnapshotsDelete(opConfig *configs.OpSnapshotDeleteConfig) (*ybApi.DeleteSnapshotResponsePB, error) {

	givenSnapshotID, err := utils.SnapshotID(opConfig.SnapshotID, opConfig.Base64Encoded)
	if err != nil {
		c.logger.Error("failed fetching normalized snapshot id",
			"given-value", opConfig.SnapshotID,
			"reason", err)
		return nil, err
	}

	protoSnapshotID, err := utils.StringUUIDToProtoSnapshotID(givenSnapshotID)
	if err != nil {
		return nil, err
	}

	payload := &ybApi.DeleteSnapshotRequestPB{
		SnapshotId: protoSnapshotID,
	}
	responsePayload := &ybApi.DeleteSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
