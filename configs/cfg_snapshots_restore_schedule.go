package configs

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

// OpSnapshotRestoreScheduleConfig represents a command specific config.
type OpSnapshotRestoreScheduleConfig struct {
	flagBase

	ScheduleID      string
	RestoreAt       uint64
	RestoreRelative time.Duration
}

// NewOpSnapshotRestoreScheduleConfig returns an instance of the command specific config.
func NewOpSnapshotRestoreScheduleConfig() *OpSnapshotRestoreScheduleConfig {
	return &OpSnapshotRestoreScheduleConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotRestoreScheduleConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.ScheduleID, "schedule-id", "", "Schedule identifier")
		c.flagSet.Uint64Var(&c.RestoreAt, "restore-at", 0, "Absolute Timing Option: Max HybridTime, in Micros")
		c.flagSet.DurationVar(&c.RestoreRelative, "restore-relative", 0, "Relative restore time in the past to fetched server clock time, takes precedence when specified")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotRestoreScheduleConfig) Validate() error {
	if c.ScheduleID == "" {
		return fmt.Errorf("--schedule-id required")
	}
	if c.RestoreAt > 0 && c.RestoreRelative > 0 {
		return fmt.Errorf("--restore-at or --restore-relative: choose one")
	}
	return nil
}