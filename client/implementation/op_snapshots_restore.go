package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	"github.com/radekg/yugabyte-db-go-client/utils/restoretarget"
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

	relativeTime, err := restoretarget.RelativeOrFixedPast(opConfig.RestoreAt,
		opConfig.RestoreRelative,
		c.defaultServerClockResolver)
	if err != nil {
		c.logger.Error("failed resolving restore target time", "reason", err)
		return nil, err
	}
	if relativeTime > 0 {
		payload.RestoreHt = utils.PUint64(relativeTime)
	}

	responsePayload := &ybApi.RestoreSnapshotResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
