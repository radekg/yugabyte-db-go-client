package implementation

import (
	"fmt"
	"time"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Restore schedule.
func (c *defaultYBCliClient) SnapshotsRestoreSchedule(opConfig *configs.OpSnapshotRestoreScheduleConfig) (*ybApi.RestoreSnapshotResponsePB, error) {

	/*
		givenScheduleID, err := utils.DecodeAsYugabyteID(opConfig.ScheduleID, opConfig.Base64Encoded)
		if err != nil {
			c.logger.Error("failed fetching normalized schedule id",
				"given-value", opConfig.ScheduleID,
				"reason", err)
			return nil, err
		}

		protoID, err := utils.StringUUIDToProtoYugabyteID(givenScheduleID)
		if err != nil {
			return nil, err
		}
	*/

	restoreAt, err := getRestoreScheduleAt(c, opConfig)
	if err != nil {
		return nil, fmt.Errorf("could not establish restore at time")
	}

	c.logger.Trace("calculated restore-at",
		"restore-at", restoreAt)

	suitableSnapshotID, err := c.suitableSnapshotID(opConfig.ScheduleID, restoreAt)
	if err != nil {
		return nil, err
	}

	suitableSnapshotIDString, err := utils.ProtoYugabyteIDToString(suitableSnapshotID)
	if err != nil {
		return nil, err
	}

	c.logger.Trace("found suitable snapshot id",
		"snapshot-id", suitableSnapshotIDString)

	// wait for the snapshot to be complete:
loop:
	for {
		snapshotsResponse, err := c.SnapshotsList(&configs.OpSnapshotListConfig{
			SnapshotID: suitableSnapshotIDString,
		})
		if err != nil {
			return nil, err
		}
		if callErr := snapshotsResponse.GetError(); callErr != nil {
			return nil, fmt.Errorf("failed loading suitable snapshot, reason: %+v", callErr)
		}
		if len(snapshotsResponse.Snapshots) != 1 {
			return nil, fmt.Errorf("wrong number of snapshots received: %d", len(snapshotsResponse.Snapshots))
		}

		c.logger.Trace("loaded snapshot for suitable snapshot id",
			"snapshot-id", suitableSnapshotIDString,
			"snapshot", snapshotsResponse.Snapshots[0].Entry)

		if snapshotsResponse.Snapshots[0].Entry == nil {
			return nil, fmt.Errorf("snapshot without an entry, snapshot ID %s", suitableSnapshotIDString)
		}
		if snapshotsResponse.Snapshots[0].Entry.State == nil {
			return nil, fmt.Errorf("snapshot entry without a state, snapshot ID %s", suitableSnapshotIDString)
		}

		c.logger.Trace("loaded snapshot for suitable snapshot id",
			"snapshot-id", suitableSnapshotIDString,
			"state", snapshotsResponse.Snapshots[0].Entry.State)

		switch *snapshotsResponse.Snapshots[0].Entry.State {
		case ybApi.SysSnapshotEntryPB_COMPLETE:
			break loop
		default:
			return nil, fmt.Errorf("snapshot is not suitable for restore at %d", restoreAt)
		}
	}

	restoreResponse, err := c.SnapshotsRestore(&configs.OpSnapshotRestoreConfig{
		SnapshotID: suitableSnapshotIDString,
		RestoreAt:  restoreAt,
	})
	if err != nil {
		return nil, err
	}

	return restoreResponse, nil
}

func (c *defaultYBCliClient) suitableSnapshotID(scheduleID string, restoreAt uint64) ([]byte, error) {
	for {

		schedules, err := c.SnapshotsListSchedules(func() *configs.OpSnapshotListSchedulesConfig {
			listSchedulesConfig := &configs.OpSnapshotListSchedulesConfig{}
			if len(scheduleID) > 0 {
				listSchedulesConfig.ScheduleID = scheduleID
			}
			return listSchedulesConfig
		}())

		if err != nil {
			c.logger.Error("Failed to list snapshot schedules", "reason", err)
			return nil, err
		}

		if len(schedules.Schedules) == 0 {
			return nil, fmt.Errorf("no schedule")
		}

		c.logger.Trace("found requested schedule")

		lastSnapshotTime := uint64(0)

		// only look at first schedule:
		for _, snapshot := range schedules.Schedules[0].Snapshots {
			snapshotIDString, err := utils.ProtoYugabyteIDToString(snapshot.Id)
			if err != nil {
				c.logger.Error("Snapshot without id")
				continue
			}
			snapshotHt := snapshot.Entry.SnapshotHybridTime
			if snapshotHt == nil {
				c.logger.Error("Snapshot without hybrid time", "snapshot", snapshotIDString)
				continue
			}
			if *snapshotHt > lastSnapshotTime {
				lastSnapshotTime = *snapshotHt
			}

			// is it suitable...
			if c.snapshotSuitableForRestoreAt(snapshot.Entry, restoreAt) {
				c.logger.Info("snaphost picked for restore",
					"snapshot-id", snapshotIDString)
				return snapshot.Id, nil
			}

			c.logger.Info("snapshot rejected for restore",
				"snapshot-id", snapshotIDString)

		}

		if lastSnapshotTime > restoreAt {
			return nil, fmt.Errorf("Cannot restore at %d, last snapshot: %d, snapshots: %+v",
				restoreAt, lastSnapshotTime, schedules.Schedules[0].Snapshots)
		}

		// create a snapshot:
		createResponse, err := c.SnapshotsCreate(&configs.OpSnapshotCreateConfig{
			ScheduleID: scheduleID,
		})
		if err != nil {
			return nil, err
		}
		if callErr := createResponse.GetError(); callErr != nil {
			switch *callErr.Code {
			case ybApi.MasterErrorPB_PARALLEL_SNAPSHOT_OPERATION:
				<-time.After(time.Second)
				continue
			default:
				return nil, fmt.Errorf("failed creating snapshot, reason: %v", callErr)
			}
		}

		return createResponse.SnapshotId, nil
	}
}

func (c *defaultYBCliClient) snapshotSuitableForRestoreAt(entry *ybApi.SysSnapshotEntryPB, restoreAt uint64) bool {
	if entry.State == nil || entry.PreviousSnapshotHybridTime == nil || entry.SnapshotHybridTime == nil {
		return false
	}
	if *entry.State == ybApi.SysSnapshotEntryPB_CREATING || *entry.State == ybApi.SysSnapshotEntryPB_COMPLETE {
		return *entry.SnapshotHybridTime >= restoreAt && *entry.PreviousSnapshotHybridTime < restoreAt
	}
	return false
}

func getRestoreScheduleAt(c *defaultYBCliClient, opConfig *configs.OpSnapshotRestoreScheduleConfig) (uint64, error) {

	restoreAt := uint64(0)

	if opConfig.RestoreAt > 0 {
		restoreAt = opConfig.RestoreAt
	}

	if opConfig.RestoreRelative > 0 {
		serverClock, err := c.ServerClock()
		if err != nil {
			return 0, err
		}
		if serverClock.HybridTime == nil {
			return 0, fmt.Errorf("no hybrid time in server clock response")
		}
		restoreAt = *serverClock.HybridTime - utils.ClockTimestampToHTTimestamp(uint64(opConfig.RestoreRelative.Microseconds()))
		return restoreAt, nil
	}

	serverClock, err := c.ServerClock()
	if err != nil {
		return 0, err
	}
	if serverClock.HybridTime == nil {
		return 0, fmt.Errorf("no hybrid time in server clock response")
	}
	return *serverClock.HybridTime, nil

}
