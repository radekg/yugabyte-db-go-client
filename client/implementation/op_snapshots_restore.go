package implementation

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Restore snapshot.
func (c *defaultYBCliClient) SnapshotsRestore(opConfig *configs.OpSnapshotRestoreConfig) (*ybApi.RestoreSnapshotResponsePB, error) {

	ybDbID, err := ybdbid.TryParseFromString(opConfig.SnapshotID)
	if err != nil {
		c.logger.Error("given snapshot id is not valid",
			"original-value", opConfig.SnapshotID,
			"reason", err)
		return nil, err
	}

	payload := &ybApi.RestoreSnapshotRequestPB{
		SnapshotId: ybDbID.Bytes(),
	}

	if opConfig.RestoreAt > 0 {
		payload.RestoreHt = &opConfig.RestoreAt
	}
	if opConfig.RestoreRelative > 0 {
		serverClock, err := c.ServerClock()
		if err != nil {
			return nil, err
		}
		if serverClock.HybridTime == nil {
			return nil, fmt.Errorf("no hybrid time in server clock response")
		}
		newHybridTime := *serverClock.HybridTime - utils.ClockTimestampToHTTimestamp(uint64(opConfig.RestoreRelative.Microseconds()))
		payload.RestoreHt = &newHybridTime
	}

	responsePayload := &ybApi.RestoreSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
