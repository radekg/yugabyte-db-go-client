package cli

import (
	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/client/base"
	"github.com/radekg/yugabyte-db-go-client/configs"
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
	LeaderStepDown(*configs.OpLeaderStepDownConfig) (*ybApi.LeaderStepDownResponsePB, error)
	ListMasters() (*ybApi.ListMastersResponsePB, error)
	ListTabletServers(*configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error)
	MasterLeaderStepDown() (*ybApi.GetMasterRegistrationResponsePB, error)
	Ping() (*ybApi.PingResponsePB, error)
	SetLoadBalancerState(bool) (*ybApi.ChangeLoadBalancerStateResponsePB, error)
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

// OnConnected returns a channel which closed when the client is connected.
func (c *defaultYBCliClient) OnConnected() <-chan struct{} {
	return c.connectedClient.OnConnected()
}

// OnConnectError returns a channel which will return an error if connect fails.
func (c *defaultYBCliClient) OnConnectError() <-chan error {
	return c.connectedClient.OnConnectError()
}
