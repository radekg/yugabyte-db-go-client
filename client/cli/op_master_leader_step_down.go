package cli

import (
	"fmt"
	"time"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// MasterTabletID is the master TableID
// TODO: find out how to get this via the API.
const MasterTabletID = "00000000000000000000000000000000"

// MasterLeaderStepDown attempts a master leader step down procedure.
func (c *defaultYBCliClient) MasterLeaderStepDown() (*ybApi.GetMasterRegistrationResponsePB, error) {

	masterRegistration, err := c.GetMasterRegistration()
	if err != nil {
		return nil, err
	}

	payload := &ybApi.LeaderStepDownRequestPB{
		DestUuid: masterRegistration.InstanceId.PermanentUuid,
		TabletId: []byte(MasterTabletID),
	}
	responsePayload := &ybApi.LeaderStepDownResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
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
		default:
			continuousErrors = 0
		}
		time.Sleep(time.Second)
	}
	return registration, nil
}
