package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// OpGetTableLocationsConfig represents a command specific config.
type OpGetTableLocationsConfig struct {
	flagBase

	Keyspace string
	Name     string
	UUID     string

	PartitionKeyStart     []byte
	PartitionKeyEnd       []byte
	MaxReturnedLocations  uint32
	RequireTabletsRunning bool
}

// NewOpGetTableLocationsConfig returns an instance of the command specific config.
func NewOpGetTableLocationsConfig() *OpGetTableLocationsConfig {
	return &OpGetTableLocationsConfig{
		MaxReturnedLocations: uint32(10),
	}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpGetTableLocationsConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "Keyspace to check in")
		c.flagSet.StringVar(&c.Name, "name", "", "Table name to check for")
		c.flagSet.StringVar(&c.UUID, "uuid", "", "Table identifier to check for")
		c.flagSet.BytesBase64Var(&c.PartitionKeyStart, "partition-key-start", []byte{}, "Partition key range start")
		c.flagSet.BytesBase64Var(&c.PartitionKeyEnd, "partition-key-end", []byte{}, "Partition key range end")
		c.flagSet.Uint32Var(&c.MaxReturnedLocations, "max-returned-locations", 10, "Maximum number of returned locations")
		c.flagSet.BoolVar(&c.RequireTabletsRunning, "require-tablet-running", false, "Require tablet running")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpGetTableLocationsConfig) Validate() error {
	if c.Name != "" && c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required when --name is given")
	}
	if c.Name == "" && c.UUID == "" {
		return fmt.Errorf("--name or --uuid is required")
	}
	return nil
}

// ==

// OpGetTableSchemaConfig represents a command specific config.
type OpGetTableSchemaConfig struct {
	flagBase

	Keyspace string
	Name     string
	UUID     string
}

// NewOpGetTableSchemaConfig returns an instance of the command specific config.
func NewOpGetTableSchemaConfig() *OpGetTableSchemaConfig {
	return &OpGetTableSchemaConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpGetTableSchemaConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "Keyspace to check in")
		c.flagSet.StringVar(&c.Name, "name", "", "Table name to check for")
		c.flagSet.StringVar(&c.UUID, "uuid", "", "Table identifier to check for")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpGetTableSchemaConfig) Validate() error {
	if c.Name != "" && c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required when --name is given")
	}
	if c.Name == "" && c.UUID == "" {
		return fmt.Errorf("--name or --uuid is required")
	}
	return nil
}

// ==

var (
	supportedNamespaceType = []string{"cql", "pgsql", "redis"}
	supportedRelationType  = []string{"system_table", "user_table", "index"}
)

// OpListTablesConfig represents a command specific config.
type OpListTablesConfig struct {
	flagBase

	NameFilter          string
	NamespaceName       string
	NamespaceType       string
	ExcludeSystemTables bool
	IncludeNotRunning   bool
	RelationType        []string
}

// NewOpListTablesConfig returns an instance of the command specific config.
func NewOpListTablesConfig() *OpListTablesConfig {
	return &OpListTablesConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpListTablesConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.NameFilter, "name-filter", "", "When used, only returns tables that satisfy a substring match on name_filter")
		c.flagSet.StringVar(&c.NamespaceName, "keyspace", "", "The namespace name to fetch info")
		c.flagSet.StringVar(&c.NamespaceType, "namespace-type", "", fmt.Sprintf("Database type: %s", strings.Join(supportedNamespaceType, ", ")))
		c.flagSet.BoolVar(&c.ExcludeSystemTables, "exclude-system-tables", false, "Exclude system tables")
		c.flagSet.BoolVar(&c.IncludeNotRunning, "include-not-running", false, "Include not running")
		c.flagSet.StringSliceVar(&c.RelationType, "relation-type", supportedRelationType, fmt.Sprintf("Filter tables based on RelationType: %s", strings.Join(supportedRelationType, ", ")))
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpListTablesConfig) Validate() error {
	if c.NamespaceType != "" {
		var found bool
		for _, opt := range supportedNamespaceType {
			if opt == c.NamespaceType {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unsupported value '%s' for --namespace-type", c.NamespaceType)
		}
	}
	for _, relation := range c.RelationType {
		var found bool
		for _, opt := range supportedRelationType {
			if opt == relation {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unsupported value '%s' for --relation-type", relation)
		}
	}
	return nil
}

// ==

// OpListTabletServersConfig represents a command specific config.
type OpListTabletServersConfig struct {
	flagBase

	PrimaryOnly bool
}

// NewOpListTableServersConfig returns an instance of the command specific config.
func NewOpListTableServersConfig() *OpListTabletServersConfig {
	return &OpListTabletServersConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpListTabletServersConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.BoolVar(&c.PrimaryOnly, "primary-only", false, "Primary only")
	}
	return c.flagSet
}

// ==

// OpIsLoadBalancedConfig represents a command specific config.
type OpIsLoadBalancedConfig struct {
	flagBase

	ExpectedNumServers int
}

// NewOpIsLoadBalancedConfig returns an instance of the command specific config.
func NewOpIsLoadBalancedConfig() *OpIsLoadBalancedConfig {
	return &OpIsLoadBalancedConfig{
		ExpectedNumServers: 0,
	}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpIsLoadBalancedConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.IntVar(&c.ExpectedNumServers, "expected-num-servers", 0, "How many servers to expect")
	}
	return c.flagSet
}

// ==

// OpPingConfig represents a command specific config.
type OpPingConfig struct {
	flagBase

	Host string
	Port int
}

// NewOpPingConfig returns an instance of the command specific config.
func NewOpPingConfig() *OpPingConfig {
	return &OpPingConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpPingConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Host, "host", "", "Host to ping")
		c.flagSet.IntVar(&c.Port, "port", 0, "Port to ping")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpPingConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("--host is required")
	}
	if c.Port < 1 {
		return fmt.Errorf("--port is required")
	}
	return nil
}

// ==

// OpIsServerReadyConfig represents a command specific config.
type OpIsServerReadyConfig struct {
	flagBase

	Host      string
	Port      int
	IsTserver bool
}

// NewOpIsServerReadyConfig returns an instance of the command specific config.
func NewOpIsServerReadyConfig() *OpPingConfig {
	return &OpPingConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpIsServerReadyConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Host, "host", "", "Host to check")
		c.flagSet.IntVar(&c.Port, "port", 0, "Port to check")
		c.flagSet.BoolVar(&c.IsTserver, "is-tserver", false, "Is TServer?")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpIsServerReadyConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("--host is required")
	}
	if c.Port < 1 {
		return fmt.Errorf("--port is required")
	}
	return nil
}

// ==

// OpLeaderStepDownConfig represents a command specific config.
type OpLeaderStepDownConfig struct {
	flagBase

	DestUUID                 string
	DisableGracefulTansition bool
	TabletID                 string
	NewLeaderUUID            string
}

// NewOpLeaderStepDownConfig returns an instance of the command specific config.
func NewOpLeaderStepDownConfig() *OpLeaderStepDownConfig {
	return &OpLeaderStepDownConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpLeaderStepDownConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.DestUUID, "destination-uuid", "", "UUID of server this request is addressed to")
		c.flagSet.BoolVar(&c.DisableGracefulTansition, "disable-graceful-transition", false, "If new_leader_uuid is not specified, the current leader will attempt to gracefully transfer leadership to another peer. Setting this flag disables that behavior")
		c.flagSet.StringVar(&c.NewLeaderUUID, "new-leader-uuid", "", "UUID of the server that should run the election to become the new leader")
		c.flagSet.StringVar(&c.TabletID, "tablet-id", "", "The id of the tablet")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpLeaderStepDownConfig) Validate() error {
	if c.DestUUID == "" || c.TabletID == "" {
		return fmt.Errorf("--destination-uuid and --tablet-id required")
	}
	return nil
}

// ==

// OpSetLoadBalancerEnableConfig represents a command specific config.
type OpSetLoadBalancerEnableConfig struct {
	flagBase

	Enabled  bool
	Disabled bool
}

// NewOpSetLoadBalancerEnableConfig returns an instance of the command specific config.
func NewOpSetLoadBalancerEnableConfig() *OpSetLoadBalancerEnableConfig {
	return &OpSetLoadBalancerEnableConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSetLoadBalancerEnableConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.BoolVar(&c.Enabled, "enabled", false, "Desired state: enabled")
		c.flagSet.BoolVar(&c.Disabled, "disabled", false, "Desired state: disabled")
	}
	return c.flagSet
}

// IsEnabled returns the bool mapping of the state with an ok flag, which is false if the
// given value wss not an expected one.
func (c *OpSetLoadBalancerEnableConfig) IsEnabled() (bool, bool) {
	if c.Enabled {
		return true, true
	}
	if c.Disabled {
		return false, true
	}
	return false, false
}

// Validate validates the correctness of the configuration.
func (c *OpSetLoadBalancerEnableConfig) Validate() error {
	if !c.Enabled && !c.Disabled {
		return fmt.Errorf("--enabled or --disabled required")
	}
	if c.Enabled && c.Disabled {
		return fmt.Errorf("--enabled and --disabled: choose one")
	}
	return nil
}
