package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultYBCliClient) SnapshotsDelete(opConfig *configs.OpSnapshotDeleteConfig) (*ybApi.DeleteSnapshotResponsePB, error) {

	ybDbID, err := ybdbid.TryParseFromString(opConfig.SnapshotID)
	if err != nil {
		c.logger.Error("given snapshot id is not valid",
			"original-value", opConfig.SnapshotID,
			"reason", err)
		return nil, err
	}

	payload := &ybApi.DeleteSnapshotRequestPB{
		SnapshotId: ybDbID.Bytes(),
	}
	responsePayload := &ybApi.DeleteSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
