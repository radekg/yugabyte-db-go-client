package cli

import (
	"fmt"

	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// GetTableSchema returns table schema if table exists or an error.
func (c *defaultYBCliClient) GetTableSchema(opConfig *configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error) {
	payload := &ybApi.GetTableSchemaRequestPB{
		Table: &ybApi.TableIdentifierPB{
			Namespace: &ybApi.NamespaceIdentifierPB{
				Name: utils.PString(opConfig.Keyspace),
			},
			TableName: func() *string {
				if opConfig.Name == "" {
					return nil
				}
				return utils.PString(opConfig.Name)
			}(),
			TableId: func() []byte {
				if opConfig.UUID == "" {
					return []byte{}
				}
				return []byte(opConfig.UUID)
			}(),
		},
	}
	responsePayload := &ybApi.GetTableSchemaResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}
