package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultYBCliClient) SnapshotsDeleteSchedule(opConfig *configs.OpSnapshotDeleteScheduleConfig) (*ybApi.DeleteSnapshotScheduleResponsePB, error) {

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

	payload := &ybApi.DeleteSnapshotScheduleRequestPB{
		SnapshotScheduleId: protoID,
	}
	responsePayload := &ybApi.DeleteSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
