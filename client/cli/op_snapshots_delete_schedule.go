package cli

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Delete snapshot.
func (c *defaultYBCliClient) SnapshotsDeleteSchedule(opConfig *configs.OpSnapshotDeleteScheduleConfig) (*ybApi.DeleteSnapshotScheduleResponsePB, error) {
	payload := &ybApi.DeleteSnapshotScheduleRequestPB{
		SnapshotScheduleId: opConfig.ScheduleID,
	}
	responsePayload := &ybApi.DeleteSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
