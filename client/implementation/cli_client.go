package implementation

import (
	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/client/base"
	"github.com/radekg/yugabyte-db-go-client/configs"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

// YBCliClient is a client implementing the CLI functionality.
type YBCliClient interface {
	Close() error

	CheckExists(*configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error)
	DescribeTable(*configs.OpGetTableSchemaConfig) (*ybApi.GetTableSchemaResponsePB, error)
	GetIsLoadBalancerIdle() (*ybApi.IsLoadBalancerIdleResponsePB, error)
	GetLoadMoveCompletion() (*ybApi.GetLoadMovePercentResponsePB, error)
	GetMasterRegistration() (*ybApi.GetMasterRegistrationResponsePB, error)
	GetTabletsForTable(*configs.OpGetTableLocationsConfig) (*ybApi.GetTableLocationsResponsePB, error)
	GetUniverseConfig() (*ybApi.GetMasterClusterConfigResponsePB, error)
	IsLoadBalanced(*configs.OpIsLoadBalancedConfig) (*ybApi.IsLoadBalancedResponsePB, error)
	IsTabletServerReady() (*ybApi.IsTabletServerReadyResponsePB, error)
	LeaderStepDown(*configs.OpLeaderStepDownConfig) (*ybApi.LeaderStepDownResponsePB, error)
	ListMasters() (*ybApi.ListMastersResponsePB, error)
	ListTables(*configs.OpListTablesConfig) (*ybApi.ListTablesResponsePB, error)
	ListTabletServers(*configs.OpListTabletServersConfig) (*ybApi.ListTabletServersResponsePB, error)
	MasterLeaderStepDown(*configs.OpMMasterLeaderStepdownConfig) (*ybApi.GetMasterRegistrationResponsePB, error)
	ModifyPlacementInfo(*configs.OpModifyPlacementInfoConfig) (*ybApi.ChangeMasterClusterConfigResponsePB, error)
	Ping() (*ybApi.PingResponsePB, error)
	SetLoadBalancerState(bool) (*ybApi.ChangeLoadBalancerStateResponsePB, error)
	SetPreferredZones(*configs.OpSetPreferredZonesConfig) (*ybApi.SetPreferredZonesResponsePB, error)

	ServerClock() (*ybApi.ServerClockResponsePB, error)

	SnapshotsCreateSchedule(*configs.OpSnapshotCreateScheduleConfig) (*ybApi.CreateSnapshotScheduleResponsePB, error)
	SnapshotsCreate(*configs.OpSnapshotCreateConfig) (*ybApi.CreateSnapshotResponsePB, error)
	SnapshotsDeleteSchedule(*configs.OpSnapshotDeleteScheduleConfig) (*ybApi.DeleteSnapshotScheduleResponsePB, error)
	SnapshotsDelete(*configs.OpSnapshotDeleteConfig) (*ybApi.DeleteSnapshotResponsePB, error)
	SnapshotsExport(*configs.OpSnapshotExportConfig) (*SnapshotExportData, error)
	SnapshotsImport(*configs.OpSnapshotImportConfig) (*ybApi.ImportSnapshotMetaResponsePB, error)
	SnapshotsListSchedules(*configs.OpSnapshotListSchedulesConfig) (*ybApi.ListSnapshotSchedulesResponsePB, error)
	SnapshotsListRestorations(*configs.OpSnapshotListRestorationsConfig) (*ybApi.ListSnapshotRestorationsResponsePB, error)
	SnapshotsList(*configs.OpSnapshotListConfig) (*ybApi.ListSnapshotsResponsePB, error)
	SnapshotsRestoreSchedule(*configs.OpSnapshotRestoreScheduleConfig) (*ybApi.RestoreSnapshotResponsePB, error)
	SnapshotsRestore(*configs.OpSnapshotRestoreConfig) (*ybApi.RestoreSnapshotResponsePB, error)

	YsqlCatalogVersion() (*ybApi.GetYsqlCatalogConfigResponsePB, error)

	OnConnected() <-chan struct{}
	OnConnectError() <-chan error
}

type defaultYBCliClient struct {
	connectedClient base.YBConnectedClient
	logger          hclog.Logger
}

// NewYBConnectedClient returns a configured instance of the default CLI client.
func NewYBConnectedClient(cfg *configs.YBClientConfig, logger hclog.Logger) (YBCliClient, error) {
	connectedClient, err := base.Connect(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &defaultYBCliClient{
		connectedClient: connectedClient,
		logger:          logger,
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
