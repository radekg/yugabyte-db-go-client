package implementation

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"google.golang.org/protobuf/proto"
)

// Import snapshot.
func (c *defaultYBCliClient) SnapshotsImport(opConfig *configs.OpSnapshotExportConfig) (*ybApi.ImportSnapshotMetaResponsePB, error) {
	statResult, err := os.Stat(opConfig.FilePath)
	if err != nil {
		return nil, err
	}
	if statResult.IsDir() {
		return nil, fmt.Errorf("path %s points at a directory", opConfig.FilePath)
	}

	rawProtoBytes, err := ioutil.ReadFile(opConfig.FilePath)
	if err != nil {
		return nil, err
	}

	snapshotInfo := &ybApi.SnapshotInfoPB{}
	if err := proto.Unmarshal(rawProtoBytes, snapshotInfo); err != nil {
		return nil, err
	}

	payload := &ybApi.ImportSnapshotMetaRequestPB{
		Snapshot: snapshotInfo,
	}
	responsePayload := &ybApi.ImportSnapshotMetaResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}

	return responsePayload, nil
}
