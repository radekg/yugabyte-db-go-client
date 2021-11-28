package cli

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultYBCliClient) SnapshotsDelete(opConfig *configs.OpSnapshotDeleteConfig) (*ybApi.DeleteSnapshotResponsePB, error) {
	payload := &ybApi.DeleteSnapshotRequestPB{
		SnapshotId: []byte(opConfig.SnapshotID),
	}
	responsePayload := &ybApi.DeleteSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
