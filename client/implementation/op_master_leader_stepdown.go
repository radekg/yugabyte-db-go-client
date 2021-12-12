package implementation

import (
	"fmt"
	"time"

	"github.com/radekg/yugabyte-db-go-client/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// MasterTabletID is the master TableID
// TODO: find out how to get this via the API.
const MasterTabletID = "00000000000000000000000000000000"

// MasterLeaderStepDown attempts a master leader step down procedure.
func (c *defaultYBCliClient) MasterLeaderStepDown(opConfig *configs.OpMMasterLeaderStepdownConfig) (*ybApi.GetMasterRegistrationResponsePB, error) {

	masterRegistration, err := c.GetMasterRegistration()
	if err != nil {
		return nil, err
	}

	payload := &ybApi.LeaderStepDownRequestPB{
		DestUuid: masterRegistration.InstanceId.PermanentUuid,
		TabletId: []byte(MasterTabletID),
	}

	if opConfig.NewLeaderID != "" {
		payload.NewLeaderUuid = []byte(opConfig.NewLeaderID)
	}

	responsePayload := &ybApi.LeaderStepDownResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if responsePayload.Error != nil {
		return nil, clientErrors.NewTabletServerError(responsePayload.Error)
	}

	// TODO: this does not belong here...
	continuousErrors := 0
	var registration *ybApi.GetMasterRegistrationResponsePB
topExit:
	for {
		masterRegistration, err := c.GetMasterRegistration()
		if err != nil {
			continuousErrors = continuousErrors + 1
			if continuousErrors >= 30 {
				return nil, fmt.Errorf("electing a new leader failed for 30 continuous attempts")
			}
		}
		switch *masterRegistration.Role {
		case ybApi.RaftPeerPB_LEADER:
			registration = masterRegistration
			break topExit
		case ybApi.RaftPeerPB_FOLLOWER:
			registration = masterRegistration
			break topExit
		default:
			continuousErrors = 0
		}
		time.Sleep(time.Second)
	}
	return registration, nil
}
