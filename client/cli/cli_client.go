package cli

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/client/base"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// YBCliClient is a client implementing the CLI functionality.
type YBCliClient interface {
	Close() error
	GetLoadMoveCompletion() (*ybApi.GetLoadMovePercentResponsePB, error)
	GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error)
	GetTableSchema(*configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error)
	GetUniverseConfig() (*ybApi.GetMasterClusterConfigResponsePB, error)
	IsLoadBalanced(*configs.OpIsLoadBalancedConfig) (*ybApi.IsLoadBalancedResponsePB, error)
	IsTabletServerReady() (*ybApi.IsTabletServerReadyResponsePB, error)
	ListMasters() (*ybApi.ListMastersResponsePB, error)
	ListTabletServers(*configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error)
	Ping() (*ybApi.PingResponsePB, error)
	OnConnected() <-chan struct{}
	OnConnectError() <-chan error
}

type defaultYBCliClient struct {
	connectedClient base.YBConnectedClient
}

// NewYBConnectedClient returns a configured instance of the default CLI client.
func NewYBConnectedClient(cfg *configs.YBClientConfig, logger hclog.Logger) (YBCliClient, error) {
	connectedClient, err := base.Connect(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &defaultYBCliClient{
		connectedClient: connectedClient,
	}, nil
}

func (c *defaultYBCliClient) Close() error {
	return c.connectedClient.Close()
}

// GetLoadMoveCompletion gets the completion percentage of tablet load move from blacklisted servers.
func (c *defaultYBCliClient) GetLoadMoveCompletion() (*ybApi.GetLoadMovePercentResponsePB, error) {
	payload := &ybApi.GetLoadMovePercentRequestPB{}
	responsePayload := &ybApi.GetLoadMovePercentResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

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

// GetUniverseConfig get the placement info and blacklist info of the universe.
func (c *defaultYBCliClient) GetUniverseConfig() (*ybApi.GetMasterClusterConfigResponsePB, error) {
	payload := &ybApi.GetMasterClusterConfigRequestPB{}
	responsePayload := &ybApi.GetMasterClusterConfigResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// IsLoadBalanced returns a list of masters or an error if call failed.
func (c *defaultYBCliClient) IsLoadBalanced(opConfig *configs.OpIsLoadBalancedConfig) (*ybApi.IsLoadBalancedResponsePB, error) {
	payload := &ybApi.IsLoadBalancedRequestPB{
		ExpectedNumServers: func() *int32 {
			if opConfig.ExpectedNumServers > 0 {
				return utils.PInt32(int32(opConfig.ExpectedNumServers))
			}
			return nil
		}(),
	}
	responsePayload := &ybApi.IsLoadBalancedResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// IsTabletServerReady checks if a given tablet server is ready or returns an error.
func (c *defaultYBCliClient) IsTabletServerReady() (*ybApi.IsTabletServerReadyResponsePB, error) {
	payload := &ybApi.IsTabletServerReadyRequestPB{}
	responsePayload := &ybApi.IsTabletServerReadyResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}

// ListMasters returns a list of masters or an error if call failed.
func (c *defaultYBCliClient) ListMasters() (*ybApi.ListMastersResponsePB, error) {
	payload := &ybApi.ListMastersRequestPB{}
	responsePayload := &ybApi.ListMastersResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// ListTabletServers returns a list of tablet servers or an error if call failed.
func (c *defaultYBCliClient) ListTabletServers(opConfig *configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error) {
	payload := &ybApi.ListTabletServersRequestPB{
		PrimaryOnly: utils.PBool(opConfig.PrimaryOnly),
	}
	responsePayload := &ybApi.ListTabletServersResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	if err := responsePayload.GetError(); err != nil {
		return nil, fmt.Errorf(err.String())
	}
	return responsePayload, nil
}

// Ping pings a certain YB server.
func (c *defaultYBCliClient) Ping() (*ybApi.PingResponsePB, error) {
	payload := &ybApi.PingRequestPB{}
	responsePayload := &ybApi.PingResponsePB{}
	if err := c.connectedClient.Execute(payload, responsePayload); err != nil {
		return nil, err
	}
	return responsePayload, nil
}

// OnConnected returns a channel which closed when the client is connected.
func (c *defaultYBCliClient) OnConnected() <-chan struct{} {
	return c.connectedClient.OnConnected()
}

// OnConnectError returns a channel which will return an error if connect fails.
func (c *defaultYBCliClient) OnConnectError() <-chan error {
	return c.connectedClient.OnConnectError()
}
