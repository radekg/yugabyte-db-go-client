package configs

import (
	"fmt"

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
