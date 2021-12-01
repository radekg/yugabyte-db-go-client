package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultYBCliClient) SnapshotsDeleteSchedule(opConfig *configs.OpSnapshotDeleteScheduleConfig) (*ybApi.DeleteSnapshotScheduleResponsePB, error) {

	ybDbID, err := ybdbid.TryParseFromString(opConfig.ScheduleID)
	if err != nil {
		c.logger.Error("given schedule id is not valid",
			"original-value", opConfig.ScheduleID,
			"reason", err)
		return nil, err
	}

	payload := &ybApi.DeleteSnapshotScheduleRequestPB{
		SnapshotScheduleId: ybDbID.Bytes(),
	}
	responsePayload := &ybApi.DeleteSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
