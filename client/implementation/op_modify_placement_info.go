package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	"github.com/radekg/yugabyte-db-go-client/utils/ybdbid"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Modifies the placement information (cloud, region, and zone) for a deployment.
func (c *defaultYBCliClient) ModifyPlacementInfo(opConfig *configs.OpModifyPlacementInfoConfig) (*ybApi.ChangeMasterClusterConfigResponsePB, error) {

	// Get current config:
	payloadCurrent := &ybApi.GetMasterClusterConfigRequestPB{}
	responseCurrent := &ybApi.GetMasterClusterConfigResponsePB{}
	if err := c.connectedClient.Execute(payloadCurrent, responseCurrent); err != nil {
		return nil, err
	}

	sysClusterConfigEntry := responseCurrent.ClusterConfig

	placementInfo := &ybApi.PlacementInfoPB{}
	placementInfo.NumReplicas = utils.PInt32(int32(opConfig.ReplicationFactor))
	placementInfo.PlacementBlocks = []*ybApi.PlacementBlockPB{}

	if opConfig.PlacementUUID != "" {
		placementYbDbID, err := ybdbid.TryParseFromString(opConfig.PlacementUUID)
		if err != nil {
			return nil, err
		}
		placementInfo.PlacementUuid = placementYbDbID.Bytes()
	} else if sysClusterConfigEntry.ReplicationInfo != nil &&
		sysClusterConfigEntry.ReplicationInfo.LiveReplicas != nil &&
		sysClusterConfigEntry.ReplicationInfo.LiveReplicas.PlacementUuid != nil {
		placementInfo.PlacementUuid = sysClusterConfigEntry.ReplicationInfo.LiveReplicas.PlacementUuid
	}

	placementToMinReplicas := map[string]int32{}
	for _, z := range opConfig.PlacementInfos {
		if _, ok := placementToMinReplicas[z]; !ok {
			placementToMinReplicas[z] = 0
		}
		placementToMinReplicas[z] = placementToMinReplicas[z] + 1
	}

	for placement, minReplicas := range placementToMinReplicas {
		cloudInfo, err := zoneInfoToCloudPB(placement)
		if err != nil {
			return nil, err
		}
		placementInfo.PlacementBlocks = append(placementInfo.PlacementBlocks, &ybApi.PlacementBlockPB{
			CloudInfo:      cloudInfo,
			MinNumReplicas: utils.PInt32(minReplicas),
		})
	}

	sysClusterConfigEntry.ReplicationInfo.LiveReplicas = placementInfo

	payload := &ybApi.ChangeMasterClusterConfigRequestPB{
		ClusterConfig: sysClusterConfigEntry,
	}

	// change configuration:
	responsePayload := &ybApi.ChangeMasterClusterConfigResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
