package configs

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

// OpSnapshotRestoreConfig represents a command specific config.
type OpSnapshotRestoreConfig struct {
	flagBase

	SnapshotID      string
	RestoreAt       uint64
	RestoreRelative time.Duration
}

// NewOpSnapshotRestoreConfig returns an instance of the command specific config.
func NewOpSnapshotRestoreConfig() *OpSnapshotRestoreConfig {
	return &OpSnapshotRestoreConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotRestoreConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.Uint64Var(&c.RestoreAt, "restore-at", 0, "Absolute Timing Option: Max HybridTime, in Micros")
		c.flagSet.DurationVar(&c.RestoreRelative, "restore-relative", 0, "Relative restore time in the past to fetched server clock time, takes precedence when specified")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotRestoreConfig) Validate() error {
	if c.SnapshotID == "" {
		return fmt.Errorf("--snapshot-id required")
	}
	if c.RestoreAt > 0 && c.RestoreRelative > 0 {
		return fmt.Errorf("--restore-at or --restore-relative: choose one")
	}
	return nil
}
