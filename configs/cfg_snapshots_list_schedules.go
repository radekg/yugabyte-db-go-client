package configs

import "github.com/spf13/pflag"

// OpSnapshotListSchedulesConfig represents a command specific config.
type OpSnapshotListSchedulesConfig struct {
	flagBase

	ScheduleID string
}

// NewOpSnapshotListSchedulesConfig returns an instance of the command specific config.
func NewOpSnapshotListSchedulesConfig() *OpSnapshotListSchedulesConfig {
	return &OpSnapshotListSchedulesConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotListSchedulesConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.ScheduleID, "schedule-id", "", "Snapshot schedule identifier")
	}
	return c.flagSet
}
