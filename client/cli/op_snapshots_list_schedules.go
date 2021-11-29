package cli

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// List snapshots.
func (c *defaultYBCliClient) SnapshotsListSchedules(opConfig *configs.OpSnapshotListSchedulesConfig) (*ybApi.ListSnapshotSchedulesResponsePB, error) {
	payload := &ybApi.ListSnapshotSchedulesRequestPB{}
	if len(opConfig.ScheduleID) > 0 {
		payload.SnapshotScheduleId = opConfig.ScheduleID
	}
	responsePayload := &ybApi.ListSnapshotSchedulesResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
