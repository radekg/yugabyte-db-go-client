package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultYBCliClient) SnapshotsListSchedules(opConfig *configs.OpSnapshotListSchedulesConfig) (*ybApi.ListSnapshotSchedulesResponsePB, error) {
	payload := &ybApi.ListSnapshotSchedulesRequestPB{}
	if len(opConfig.ScheduleID) > 0 {

		decodedScheduleID, err := utils.DecodeAsYugabyteID(opConfig.ScheduleID, opConfig.Base64Encoded)
		if err != nil {
			c.logger.Error("failed fetching normalized schedule id",
				"given-value", opConfig.ScheduleID,
				"reason", err)
			return nil, err
		}

		protoID, err := utils.StringUUIDToProtoYugabyteID(decodedScheduleID)
		if err != nil {
			return nil, err
		}

		payload.SnapshotScheduleId = protoID
	}
	responsePayload := &ybApi.ListSnapshotSchedulesResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
