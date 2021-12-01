package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultYBCliClient) SnapshotsListRestorations(opConfig *configs.OpSnapshotListRestorationsConfig) (*ybApi.ListSnapshotRestorationsResponsePB, error) {

	useSnapshotID := []byte{}
	useRestorationID := []byte{}
	if opConfig.SnapshotID != "" {
		ybDbID, err := ybdbid.TryParseFromString(opConfig.SnapshotID)
		if err != nil {
			c.logger.Error("given snapshot id is not valid",
				"original-value", opConfig.SnapshotID,
				"reason", err)
			return nil, err
		}
		useSnapshotID = ybDbID.Bytes()
	}

	if opConfig.RestorationID != "" {
		ybDbID, err := ybdbid.TryParseFromString(opConfig.RestorationID)
		if err != nil {
			c.logger.Error("given snapshot id is not valid",
				"original-value", opConfig.SnapshotID,
				"reason", err)
			return nil, err
		}
		useRestorationID = ybDbID.Bytes()
	}

	payload := &ybApi.ListSnapshotRestorationsRequestPB{}
	if len(useSnapshotID) > 0 {
		payload.SnapshotId = useSnapshotID
	}
	if len(useRestorationID) > 0 {
		payload.RestorationId = useRestorationID
	}

	responsePayload := &ybApi.ListSnapshotRestorationsResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
