package implementation

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// SnapshotExportData contains snapshot file info.
type SnapshotExportData struct {
	FilePath string `json:"file_path"`
	Size     int64  `json:"size"`
}

// Export snapshot.
func (c *defaultYBCliClient) SnapshotsExport(opConfig *configs.OpSnapshotExportConfig) (*SnapshotExportData, error) {

	givenSnapshotID, err := utils.DecodeAsYugabyteID(opConfig.SnapshotID, opConfig.Base64Encoded)
	if err != nil {
		c.logger.Error("failed fetching normalized snapshot id",
			"given-value", opConfig.SnapshotID,
			"reason", err)
		return nil, err
	}

	protoSnapshotID, err := utils.StringUUIDToProtoYugabyteID(givenSnapshotID)
	if err != nil {
		return nil, err
	}

	payload := &ybApi.ListSnapshotsRequestPB{
		PrepareForBackup: func() *bool {
			v := true
			return &v
		}(),
		SnapshotId: protoSnapshotID,
	}
	responsePayload := &ybApi.ListSnapshotsResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	if len(responsePayload.Snapshots) > 1 {
		c.logger.Warn("too many snapshots returned, expected at  most 1",
			"snapshot-id", givenSnapshotID,
			"found", len(responsePayload.Snapshots))
	}

	var snapshotExportEntry *ybApi.SnapshotInfoPB
	for _, snapshotEntry := range responsePayload.Snapshots {
		stringID, err := utils.ProtoYugabyteIDToString(snapshotEntry.Id)
		if err != nil {
			c.logger.Warn("skipping snapshot, could not parse Id value as string UUID")
			continue
		}
		if stringID == givenSnapshotID {
			snapshotExportEntry = snapshotEntry
			break
		}
	}

	if snapshotExportEntry == nil {
		return nil, fmt.Errorf("Snapshot '%s' not found", givenSnapshotID)
	}

	bys, err := utils.SerializeProto(snapshotExportEntry)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(opConfig.FilePath, bys, 0644); err != nil {
		return nil, err
	}

	statResult, err := os.Stat(opConfig.FilePath)
	if err != nil {
		return nil, err
	}

	return &SnapshotExportData{
		FilePath: opConfig.FilePath,
		Size:     statResult.Size(),
	}, nil
}
