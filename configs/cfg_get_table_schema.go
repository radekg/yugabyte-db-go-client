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
	if c.Name != "" && c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required when --name is given")
	}
	if c.Name == "" && c.UUID == "" {
		return fmt.Errorf("--name or --uuid is required")
	}
	return nil
}
