package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

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
	if c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required")
	}
	if c.Name == "" && c.UUID == "" {
		return fmt.Errorf("--name or --uuid is required")
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
