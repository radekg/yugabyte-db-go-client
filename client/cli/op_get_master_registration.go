package cli

import (
	"fmt"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// GetMasterRegistration retrieves the master registration information or error if the request failed.
func (c *defaultYBCliClient) GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error) {
	payload := &ybApi.GetMasterRegistrationRequestPB{}
	responsePayload := &ybApi.GetMasterRegistrationResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}
