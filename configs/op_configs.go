package configs

import (
	"fmt"
	"strings"
	"time"

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
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "Keyspace to describe the table in")
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
	supportedNamespaceType = []string{"ycql", "ysql", "yedis"}
	supportedRelationType  = []string{"system_table", "user_table", "index"}
)

// OpListTablesConfig represents a command specific config.
type OpListTablesConfig struct {
	flagBase

	NameFilter          string
	Keyspace            string
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
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "The namespace name to fetch info")
		c.flagSet.BoolVar(&c.ExcludeSystemTables, "exclude-system-tables", false, "Exclude system tables")
		c.flagSet.BoolVar(&c.IncludeNotRunning, "include-not-running", false, "Include not running")
		c.flagSet.StringSliceVar(&c.RelationType, "relation-type", supportedRelationType, fmt.Sprintf("Filter tables based on RelationType: %s", strings.Join(supportedRelationType, ", ")))
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpListTablesConfig) Validate() error {
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

// ==

// OpSnapshotCreateScheduleConfig represents a command specific config.
type OpSnapshotCreateScheduleConfig struct {
	flagBase

	Keyspace              string
	IntervalSecs          time.Duration
	RetendionDurationSecs time.Duration
	DeleteAfter           time.Duration
	DeleteTime            uint64
}

// NewOpSnapshotCreateScheduleConfig returns an instance of the command specific config.
func NewOpSnapshotCreateScheduleConfig() *OpSnapshotCreateScheduleConfig {
	return &OpSnapshotCreateScheduleConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotCreateScheduleConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "Keyspace for the tables in this create request")
		c.flagSet.DurationVar(&c.IntervalSecs, "interval", time.Second*0, "Interval for taking snapshot in seconds")
		c.flagSet.DurationVar(&c.RetendionDurationSecs, "retention-duration", time.Second*0, "How long store snapshots in seconds")
		c.flagSet.DurationVar(&c.DeleteAfter, "delete-after", time.Second*0, "How long until schedule is removed in seconds, hybrid time will be calculated by fetching server hybrid time and adding this value")
		c.flagSet.Uint64Var(&c.DeleteTime, "delete-at", 0, "Hybrid time when this schedule is deleted")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotCreateScheduleConfig) Validate() error {
	if c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required")
	}
	if c.DeleteAfter.Milliseconds() > 0 && c.DeleteTime > 0 {
		return fmt.Errorf("--delete-after and --delete-at specified, choose one")
	}
	return nil
}

// ==

// OpSnapshotCreateConfig represents a command specific config.
type OpSnapshotCreateConfig struct {
	flagBase

	Keyspace         string
	TableNames       []string
	TableUUIDs       []string
	TransactionAware bool
	AddIndexes       bool
	Imported         bool
	ScheduleID       []byte
}

// NewOpSnapshotCreateConfig returns an instance of the command specific config.
func NewOpSnapshotCreateConfig() *OpSnapshotCreateConfig {
	return &OpSnapshotCreateConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotCreateConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "Keyspace for the tables in this create request")
		c.flagSet.StringSliceVar(&c.TableNames, "name", []string{}, "Table names to create snapshots for")
		c.flagSet.StringSliceVar(&c.TableUUIDs, "uuid", []string{}, "Table IDs to create snapshots for")
		c.flagSet.BoolVar(&c.TransactionAware, "transaction-aware", false, "Transaction aware")
		c.flagSet.BoolVar(&c.AddIndexes, "add-indexes", false, "Add indexes")
		c.flagSet.BoolVar(&c.Imported, "imported", false, "Interpret this snapshot as imported")
		c.flagSet.BytesBase64Var(&c.ScheduleID, "schedule-id", []byte{}, "Create snapshot to this schedule, other fields are ignored")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotCreateConfig) Validate() error {
	if len(c.ScheduleID) > 0 && c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required")
	}
	return nil
}

// ==

// OpSnapshotDeleteScheduleConfig represents a command specific config.
type OpSnapshotDeleteScheduleConfig struct {
	flagBase

	ScheduleID []byte
}

// NewOpSnapshotDeleteScheduleConfig returns an instance of the command specific config.
func NewOpSnapshotDeleteScheduleConfig() *OpSnapshotDeleteScheduleConfig {
	return &OpSnapshotDeleteScheduleConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotDeleteScheduleConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.BytesBase64Var(&c.ScheduleID, "schedule-id", []byte{}, "Snapshot schedule identifier")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotDeleteScheduleConfig) Validate() error {
	if len(c.ScheduleID) == 0 {
		return fmt.Errorf("--schedule-id is required")
	}
	return nil
}

// ==

// OpSnapshotDeleteConfig represents a command specific config.
type OpSnapshotDeleteConfig struct {
	flagBase

	SnapshotID string
}

// NewOpSnapshotDeleteConfig returns an instance of the command specific config.
func NewOpSnapshotDeleteConfig() *OpSnapshotDeleteConfig {
	return &OpSnapshotDeleteConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotDeleteConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotDeleteConfig) Validate() error {
	if c.SnapshotID == "" {
		return fmt.Errorf("--snapshot-id is required")
	}
	return nil
}

// ==

// OpSnapshotListSchedulesConfig represents a command specific config.
type OpSnapshotListSchedulesConfig struct {
	flagBase

	ScheduleID []byte
}

// NewOpSnapshotListSchedulesConfig returns an instance of the command specific config.
func NewOpSnapshotListSchedulesConfig() *OpSnapshotListSchedulesConfig {
	return &OpSnapshotListSchedulesConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotListSchedulesConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.BytesBase64Var(&c.ScheduleID, "schedule-id", []byte{}, "Snapshot schedule identifier")
	}
	return c.flagSet
}

// ==

// OpSnapshotListConfig represents a command specific config.
type OpSnapshotListConfig struct {
	flagBase

	SnapshotID           string
	ListDeletedSnapshots bool
	PrepareForBackup     bool
}

// NewOpSnapshotListConfig returns an instance of the command specific config.
func NewOpSnapshotListConfig() *OpSnapshotListConfig {
	return &OpSnapshotListConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotListConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.BoolVar(&c.ListDeletedSnapshots, "list-deleted-snapshots", false, "List deleted snapshots")
		c.flagSet.BoolVar(&c.PrepareForBackup, "prepare-for-backup", false, "Prepare for backup")
	}
	return c.flagSet
}
