package implementation

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// SnapshotExportData contains snapshot file info.
type SnapshotExportData struct {
	FilePath string `json:"file_path"`
	Size     int64  `json:"size"`
}

// Export snapshot.
func (c *defaultYBCliClient) SnapshotsExport(opConfig *configs.OpSnapshotExportConfig) (*SnapshotExportData, error) {

	ybDbID, err := ybdbid.TryParseFromString(opConfig.SnapshotID)
	if err != nil {
		c.logger.Error("given snapshot id is not valid",
			"original-value", opConfig.SnapshotID,
			"reason", err)
		return nil, err
	}

	payload := &ybApi.ListSnapshotsRequestPB{
		PrepareForBackup: func() *bool {
			v := true
			return &v
		}(),
		SnapshotId: ybDbID.Bytes(),
	}
	responsePayload := &ybApi.ListSnapshotsResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	if len(responsePayload.Snapshots) > 1 {
		c.logger.Warn("too many snapshots returned, expected at  most 1",
			"snapshot-id", ybDbID.String(),
			"found", len(responsePayload.Snapshots))
	}

	var snapshotExportEntry *ybApi.SnapshotInfoPB
	for _, snapshotInfoEntry := range responsePayload.Snapshots {
		if bytes.Compare(ybDbID.Bytes(), snapshotInfoEntry.Id) == 0 {
			snapshotExportEntry = snapshotInfoEntry
			break
		}
	}

	if snapshotExportEntry == nil {
		return nil, fmt.Errorf("Snapshot '%s' not found", ybDbID.String())
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
