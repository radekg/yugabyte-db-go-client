package implementation

import ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"

// GetYsqlCatalogVersion gets the current YSQL schema catalog version.
func (c *defaultYBCliClient) GetYsqlCatalogVersion() (*ybApi.GetYsqlCatalogConfigResponsePB, error) {
	payload := &ybApi.GetYsqlCatalogConfigRequestPB{}
	responsePayload := &ybApi.GetYsqlCatalogConfigResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}
