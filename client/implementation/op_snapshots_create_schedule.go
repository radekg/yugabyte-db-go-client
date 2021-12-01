package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	"github.com/radekg/yugabyte-db-go-client/utils/relativetime"
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

	futureTime, err := relativetime.RelativeOrFixedFuture(opConfig.DeleteTime,
		opConfig.DeleteAfter,
		c.defaultServerClockResolver)
	if err != nil {
		c.logger.Error("failed resolving delete at time", "reason", err)
		return nil, err
	}
	if futureTime > 0 {
		payload.Options.DeleteTime = utils.PUint64(futureTime)
	}

	responsePayload := &ybApi.CreateSnapshotScheduleResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
