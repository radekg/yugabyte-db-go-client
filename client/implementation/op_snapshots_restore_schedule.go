package implementation

import (
	"fmt"
	"time"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils/relativetime"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Restore schedule.
func (c *defaultYBCliClient) SnapshotsRestoreSchedule(opConfig *configs.OpSnapshotRestoreScheduleConfig) (*ybApi.RestoreSnapshotResponsePB, error) {

	restoreAt, err := relativetime.RelativeOrFixedPastWithFallback(opConfig.RestoreAt,
		opConfig.RestoreRelative,
		c.defaultServerClockResolver)
	if err != nil {
		return nil, fmt.Errorf("could not establish restore at time")
	}

	c.logger.Trace("calculated restore-at",
		"restore-at", restoreAt)

	suitableSnapshotID, err := c.suitableSnapshotID(opConfig.ScheduleID, restoreAt)
	if err != nil {
		return nil, err
	}

	suitableYbDbID, err := ybdbid.TryParseFromBytes(suitableSnapshotID)
	if err != nil {
		c.logger.Error("suitable snapshot id cann't be parsed as YugabyteDB ID",
			"bytes", suitableSnapshotID,
			"reason", err)
		return nil, err
	}

	c.logger.Trace("found suitable snapshot id",
		"snapshot-id", suitableYbDbID.String())

	// wait for the snapshot to be complete:
loop:
	for {
		snapshotsResponse, err := c.SnapshotsList(&configs.OpSnapshotListConfig{
			SnapshotID: suitableYbDbID.String(),
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
			"snapshot-id", suitableYbDbID.String(),
			"snapshot", snapshotsResponse.Snapshots[0].Entry)

		if snapshotsResponse.Snapshots[0].Entry == nil {
			return nil, fmt.Errorf("snapshot without an entry, snapshot ID %s", suitableYbDbID.String())
		}
		if snapshotsResponse.Snapshots[0].Entry.State == nil {
			return nil, fmt.Errorf("snapshot entry without a state, snapshot ID %s", suitableYbDbID.String())
		}

		c.logger.Trace("loaded snapshot for suitable snapshot id",
			"snapshot-id", suitableYbDbID.String(),
			"state", snapshotsResponse.Snapshots[0].Entry.State)

		switch *snapshotsResponse.Snapshots[0].Entry.State {
		case ybApi.SysSnapshotEntryPB_COMPLETE:
			break loop
		default:
			return nil, fmt.Errorf("snapshot is not suitable for restore at %d", restoreAt)
		}
	}

	restoreResponse, err := c.SnapshotsRestore(&configs.OpSnapshotRestoreConfig{
		SnapshotID: suitableYbDbID.String(),
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
		for _, candidateSnapshot := range schedules.Schedules[0].Snapshots {

			candidateSnapshotYbDbID, err := ybdbid.TryParseFromBytes(candidateSnapshot.Id)
			if err != nil {
				c.logger.Error("skipping candidate snapshot with invalid id", "bytes", candidateSnapshot.Id)
				continue
			}
			snapshotHt := candidateSnapshot.Entry.SnapshotHybridTime
			if snapshotHt == nil {
				c.logger.Error("Snapshot without hybrid time", "snapshot", candidateSnapshotYbDbID.String())
				continue
			}
			if *snapshotHt > lastSnapshotTime {
				lastSnapshotTime = *snapshotHt
			}

			// is it suitable...
			if c.snapshotSuitableForRestoreAt(candidateSnapshot.Entry, restoreAt) {
				c.logger.Info("candidate snaphost ACCEPTED for restore",
					"snapshot-id", candidateSnapshotYbDbID.String())
				return candidateSnapshot.Id, nil
			}

			c.logger.Info("candidate snapshot REJECTED for restore",
				"snapshot-id", candidateSnapshotYbDbID.String())

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
