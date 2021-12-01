package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

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
