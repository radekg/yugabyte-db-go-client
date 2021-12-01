package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpSnapshotRestoreConfig represents a command specific config.
type OpSnapshotRestoreConfig struct {
	flagBase

	SnapshotID    string
	RestoreTarget string
}

// NewOpSnapshotRestoreConfig returns an instance of the command specific config.
func NewOpSnapshotRestoreConfig() *OpSnapshotRestoreConfig {
	return &OpSnapshotRestoreConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotRestoreConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.StringVar(&c.RestoreTarget, "restore-target", "", "Absolute Timing Option: Max HybridTime, in Micros or duration expression")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotRestoreConfig) Validate() error {
	if c.SnapshotID == "" {
		return fmt.Errorf("--snapshot-id required")
	}
	return nil
}
