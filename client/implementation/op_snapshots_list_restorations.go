package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultYBCliClient) SnapshotsListRestorations(opConfig *configs.OpSnapshotListRestorationsConfig) (*ybApi.ListSnapshotRestorationsResponsePB, error) {

	useSnapshotID := []byte{}
	useRestorationID := []byte{}
	if opConfig.SnapshotID != "" {
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
		useSnapshotID = protoSnapshotID
	}

	if opConfig.RestorationID != "" {
		givenRestorationID, err := utils.DecodeAsYugabyteID(opConfig.RestorationID, opConfig.Base64Encoded)
		if err != nil {
			c.logger.Error("failed fetching normalized restoration id",
				"given-value", opConfig.RestorationID,
				"reason", err)
			return nil, err
		}
		protoRestorationID, err := utils.StringUUIDToProtoYugabyteID(givenRestorationID)
		if err != nil {
			return nil, err
		}
		useRestorationID = protoRestorationID
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
