package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpSnapshotRestoreScheduleConfig represents a command specific config.
type OpSnapshotRestoreScheduleConfig struct {
	flagBase

	ScheduleID    string
	RestoreTarget string
}

// NewOpSnapshotRestoreScheduleConfig returns an instance of the command specific config.
func NewOpSnapshotRestoreScheduleConfig() *OpSnapshotRestoreScheduleConfig {
	return &OpSnapshotRestoreScheduleConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotRestoreScheduleConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.ScheduleID, "schedule-id", "", "Schedule identifier")
		c.flagSet.StringVar(&c.RestoreTarget, "restore-target", "", "Absolute Timing Option: Max HybridTime, in Micros or duration expression")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotRestoreScheduleConfig) Validate() error {
	if c.ScheduleID == "" {
		return fmt.Errorf("--schedule-id required")
	}
	return nil
}
