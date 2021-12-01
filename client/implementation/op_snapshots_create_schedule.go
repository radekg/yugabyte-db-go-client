package implementation

import (
	"fmt"

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
		serverClock, err := c.ServerClock()
		if err != nil {
			return nil, err
		}
		if serverClock.HybridTime == nil {
			return nil, fmt.Errorf("no hybrid time in server clock response")
		}
		newHybridTime := *serverClock.HybridTime + uint64(opConfig.DeleteAfter.Microseconds())
		payload.Options.DeleteTime = &newHybridTime
	}

	responsePayload := &ybApi.CreateSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
