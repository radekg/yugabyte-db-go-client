package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

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
