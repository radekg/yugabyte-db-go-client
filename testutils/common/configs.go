package common

import (
	"testing"

	dc "github.com/ory/dockertest/v3/docker"
)

const (
	// DefaultYugabyteDBMasterImageName specifies the default Docker image name used in tests.
	DefaultYugabyteDBMasterImageName = "yugabytedb/yugabyte"
	// DefaultYugabyteDBImageVersion specifies the default Docker image version used in tests.
	DefaultYugabyteDBImageVersion = "2.13.0.1-b2"
	// DefaultYugabyteDBContainerUser specifies the default Docker container user.
	DefaultYugabyteDBContainerUser = "yugabyte"
	// DefaultReplicationFactor specifies the default replication factor.
	DefaultReplicationFactor = 1
	// DefaultYugabyteDBFsDataPath specifies the default path of the data directory in the container.
	DefaultYugabyteDBFsDataPath = "/mnt/data"
	// DefaultYugabyteDBMasterRPCPort is the default YugabyteDB RPC port used inside of the container.
	// Currently cannot be changed.
	DefaultYugabyteDBMasterRPCPort = 7100

	// DefaultYugabyteDBTServerRPCPort is the default YugabyteDB RPC port used inside of the container.
	// Currently cannot be changed.
	DefaultYugabyteDBTServerRPCPort = 9100
	// DefaultYugabyteDBTServerYSQLPort is the default YugabyteDB PostgreSQL port used inside of the container.
	// Currently cannot be changed.
	DefaultYugabyteDBTServerYSQLPort = 5433

	// DefaultMasterPrefix is the default master prefix value, if no prefix specified.
	DefaultMasterPrefix = "yb-master"

	// DefaultYugabyteDBEnvVarImageName is an environment variable name used to override the default YugabyteDB docker image name.
	DefaultYugabyteDBEnvVarImageName = "TEST_YUGABYTEDB_IMAGE_NAME"
	// DefaultYugabyteDEnvVarImageVersion is an environment variable name used to override the default YugabyteDB docker image tag.
	DefaultYugabyteDEnvVarImageVersion = "TEST_YUGABYTEDB_IMAGE_VERSION"
)

// MasterInternalAddresses is a list of internal masters addresses used to form the cluster.
type MasterInternalAddresses = []string

// DockerCommand is the YugabyteDB Docker image command.
type DockerCommand = []string

// RPCBindAddress is the bind address used by the test container.
type RPCBindAddress = string

// TestMasterConfiguration is the master configuration for this test.
type TestMasterConfiguration struct {
	YbDBContainerUser string
	YbDBEnv           []string // format: env=value
	YbDBFsDataPath    string
	YbDBDockerImage   string
	YbDBDockerTag     string
	YbDBCmdSupplier   func(MasterInternalAddresses, RPCBindAddress) DockerCommand

	MasterPrefix string

	AdditionalPorts []dc.Port

	LogRegistrationRetryErrors bool
	NoCleanupContainers        bool
	ReplicationFactor          int
}

// ApplyMasterConfigDefaults applies default values to the master configuration.
func ApplyMasterConfigDefaults(t *testing.T, cfg *TestMasterConfiguration) *TestMasterConfiguration {
	if cfg.MasterPrefix == "" {
		cfg.MasterPrefix = DefaultMasterPrefix
	}
	if cfg.ReplicationFactor == 0 {
		cfg.ReplicationFactor = DefaultReplicationFactor
	}
	if cfg.YbDBContainerUser == "" {
		cfg.YbDBContainerUser = DefaultYugabyteDBContainerUser
	}
	if cfg.YbDBDockerImage == "" {
		cfg.YbDBDockerImage = DefaultYugabyteDBMasterImageName
	}
	if cfg.YbDBDockerTag == "" {
		cfg.YbDBDockerTag = DefaultYugabyteDBImageVersion
	}
	if cfg.YbDBFsDataPath == "" {
		cfg.YbDBFsDataPath = DefaultYugabyteDBFsDataPath
	}
	return cfg
}

// TestTServerConfiguration is the TServer configuration for this test.
type TestTServerConfiguration struct {
	YbDBContainerUser string
	YbDBEnv           []string // format: env=value
	YbDBFsDataPath    string
	YbDBDockerImage   string
	YbDBDockerTag     string
	YbDBCmdSupplier   func(MasterInternalAddresses, RPCBindAddress) DockerCommand

	// if empty, a random value will be generated and assigned
	TServerID string

	AdditionalPorts []dc.Port

	LogRegistrationRetryErrors bool
	NoCleanupContainers        bool
}

// ApplyTServerConfigDefaults applies default values to the master configuration.
func ApplyTServerConfigDefaults(t *testing.T, cfg *TestTServerConfiguration) *TestTServerConfiguration {
	if cfg.YbDBContainerUser == "" {
		cfg.YbDBContainerUser = DefaultYugabyteDBContainerUser
	}
	if cfg.YbDBDockerImage == "" {
		cfg.YbDBDockerImage = DefaultYugabyteDBMasterImageName
	}
	if cfg.YbDBDockerTag == "" {
		cfg.YbDBDockerTag = DefaultYugabyteDBImageVersion
	}
	if cfg.YbDBFsDataPath == "" {
		cfg.YbDBFsDataPath = DefaultYugabyteDBFsDataPath
	}
	return cfg
}

// --

// AllocatedAdditionalPort holds references to allocated ports when other ports are requested.
type AllocatedAdditionalPort interface {
	Allocated() string
	Requested() dc.Port
	Use()
}

type defaultAllocatedAdditionalPort struct {
	allocated string
	requested dc.Port
	supplier  RandomPortSupplier
}

// NewDefaultAllocatedAdditionalPort returns an instance of a default AllocatedAdditionalPort implementation.
func NewDefaultAllocatedAdditionalPort(r dc.Port, a string, supplier RandomPortSupplier) AllocatedAdditionalPort {
	return &defaultAllocatedAdditionalPort{
		allocated: a,
		requested: r,
		supplier:  supplier,
	}
}

func (p *defaultAllocatedAdditionalPort) Allocated() string {
	return p.allocated
}

func (p *defaultAllocatedAdditionalPort) Requested() dc.Port {
	return p.requested
}

func (p *defaultAllocatedAdditionalPort) Use() {
	p.supplier.Cleanup()
}

// AllocatedAdditionalPorts is a short hand type storing all allocated additional ports for a single container.
type AllocatedAdditionalPorts = map[dc.Port]AllocatedAdditionalPort

// MultiMasterAllocatedAdditionalPorts is a short hand type storing all allocated additional ports for masters.
// The mapping is:
//   master name => allocated ports
// Master names can be obtainer from masters test context.
type MultiMasterAllocatedAdditionalPorts = map[string]AllocatedAdditionalPorts
