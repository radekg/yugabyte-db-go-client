package common

import "testing"

const (
	// DefaultYugabyteDBMasterImageName specifies the default Docker image name used in tests.
	DefaultYugabyteDBMasterImageName = "yugabytedb/yugabyte"
	// DefaultYugabyteDBImageVersion specifies the default Docker image version used in tests.
	DefaultYugabyteDBImageVersion = "2.11.0.0-b7"
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

// TestMasterConfiguration is the master configuration for this test.
type TestMasterConfiguration struct {
	YbDBContainerUser string
	YbDBEnv           []string // format: env=value
	YbDBFsDataPath    string
	YbDBDockerImage   string
	YbDBDockerTag     string

	MasterPrefix string

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

	// if empty, a random value will be generated and assigned
	TServerID string

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
