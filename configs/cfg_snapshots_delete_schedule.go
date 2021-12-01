package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpSnapshotDeleteScheduleConfig represents a command specific config.
type OpSnapshotDeleteScheduleConfig struct {
	flagBase

	ScheduleID    string
	Base64Encoded bool
}

// NewOpSnapshotDeleteScheduleConfig returns an instance of the command specific config.
func NewOpSnapshotDeleteScheduleConfig() *OpSnapshotDeleteScheduleConfig {
	return &OpSnapshotDeleteScheduleConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotDeleteScheduleConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.ScheduleID, "schedule-id", "", "Snapshot schedule identifier")
		c.flagSet.BoolVar(&c.Base64Encoded, "base64-encoded", false, "If true, accepts the --schedule-id as base64 encoded string")
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
