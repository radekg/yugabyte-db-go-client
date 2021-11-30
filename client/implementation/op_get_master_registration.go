package implementation

import (
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// GetMasterRegistration retrieves the master registration information or error if the request failed.
func (c *defaultYBCliClient) GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error) {
	payload := &ybApi.GetMasterRegistrationRequestPB{}
	responsePayload := &ybApi.GetMasterRegistrationResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
