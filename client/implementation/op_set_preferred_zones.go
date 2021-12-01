package implementation

import (
	"fmt"
	"strings"

	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// Sets the preferred availability zones (AZs) and regions.
func (c *defaultYBCliClient) SetPreferredZones(opConfig *configs.OpSetPreferredZonesConfig) (*ybApi.SetPreferredZonesResponsePB, error) {
	zones := []*ybApi.CloudInfoPB{}
	for _, z := range opConfig.ZonesInfos {
		zz, err := zoneInfoToCloudPB(z)
		if err != nil {
			return nil, err
		}
		zones = append(zones, zz)
	}

	payload := &ybApi.SetPreferredZonesRequestPB{
		PreferredZones: zones,
	}
	responsePayload := &ybApi.SetPreferredZonesResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}

func zoneInfoToCloudPB(input string) (*ybApi.CloudInfoPB, error) {
	parts := strings.Split(input, ".")
	if len(parts) > 3 {
		return nil, fmt.Errorf("invalid zone info: %s", input)
	}
	if len(parts) == 3 {
		return &ybApi.CloudInfoPB{
			PlacementCloud:  &parts[0],
			PlacementRegion: &parts[1],
			PlacementZone:   &parts[2],
		}, nil
	}
	if len(parts) == 2 {
		return &ybApi.CloudInfoPB{
			PlacementCloud:  &parts[0],
			PlacementRegion: &parts[1],
		}, nil
	}
	if len(parts) == 1 {
		return &ybApi.CloudInfoPB{
			PlacementCloud: &parts[0],
		}, nil
	}
	return nil, fmt.Errorf("empty zone info")
}

func defaultZoneInfo() *ybApi.CloudInfoPB {
	return &ybApi.CloudInfoPB{}
}
