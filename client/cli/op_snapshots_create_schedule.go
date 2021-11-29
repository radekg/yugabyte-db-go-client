package cli

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Create a snapshot.
func (c *defaultYBCliClient) SnapshotsCreateSchedule(opConfig *configs.OpSnapshotCreateScheduleConfig) (*ybApi.CreateSnapshotScheduleResponsePB, error) {

	payload := &ybApi.CreateSnapshotScheduleRequestPB{
		Options: &ybApi.SnapshotScheduleOptionsPB{
			Filter: &ybApi.SnapshotScheduleFilterPB{
				Filter: &ybApi.SnapshotScheduleFilterPB_Tables{
					Tables: &ybApi.TableIdentifiersPB{
						Tables: []*ybApi.TableIdentifierPB{
							{
								Namespace: parseKeyspace(opConfig.Keyspace).toProtoKeyspace(),
							},
						},
					},
				},
			},
		},
	}

	if opConfig.IntervalSecs > 0 {
		payload.Options.IntervalSec = func() *uint64 {
			v := uint64(opConfig.IntervalSecs.Seconds())
			return &v
		}()
	}
	if opConfig.RetendionDurationSecs > 0 {
		payload.Options.RetentionDurationSec = func() *uint64 {
			v := uint64(opConfig.RetendionDurationSecs.Seconds())
			return &v
		}()
	}
	if opConfig.DeleteTime > 0 {
		payload.Options.DeleteTime = &opConfig.DeleteTime
	}
	if opConfig.DeleteAfter > 0 {
		payloadHybridTime := &ybApi.ServerClockRequestPB{}
		responseHybridTimePayload := &ybApi.ServerClockResponsePB{}
		if err := c.connectedClient.Execute(payloadHybridTime, responseHybridTimePayload); err != nil {
			return nil, err
		}
		newHybridTime := *responseHybridTimePayload.HybridTime + uint64(opConfig.DeleteAfter.Seconds())
		payload.Options.DeleteTime = &newHybridTime
	}

	responsePayload := &ybApi.CreateSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
