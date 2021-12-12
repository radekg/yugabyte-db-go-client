package implementation

import (
	"github.com/radekg/yugabyte-db-go-client/configs"
	clientErrors "github.com/radekg/yugabyte-db-go-client/errors"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// CheckExists returns table schema if table exists or an error.
func (c *defaultYBCliClient) CheckExists(opConfig *configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error) {
	payload := &ybApi.GetTableSchemaRequestPB{
		Table: &ybApi.TableIdentifierPB{},
	}
	if opConfig.Name != "" {
		payload.Table.Namespace = &ybApi.NamespaceIdentifierPB{
			Name: utils.PString(opConfig.Keyspace),
		}
		payload.Table.TableName = utils.PString(opConfig.Name)
	} else {
		payload.Table.TableId = []byte(opConfig.UUID)
	}
	responsePayload := &ybApi.GetTableSchemaResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, clientErrors.NewMasterError(responsePayload.Error)
}
