package configs

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

// OpSnapshotCreateScheduleConfig represents a command specific config.
type OpSnapshotCreateScheduleConfig struct {
	flagBase

	Keyspace              string
	IntervalSecs          time.Duration
	RetendionDurationSecs time.Duration
	DeleteAfter           string
}

// NewOpSnapshotCreateScheduleConfig returns an instance of the command specific config.
func NewOpSnapshotCreateScheduleConfig() *OpSnapshotCreateScheduleConfig {
	return &OpSnapshotCreateScheduleConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotCreateScheduleConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "Keyspace for the tables in this create request")
		c.flagSet.DurationVar(&c.IntervalSecs, "interval", time.Second*0, "Interval for taking snapshot in seconds")
		c.flagSet.DurationVar(&c.RetendionDurationSecs, "retention-duration", time.Second*0, "How long store snapshots in seconds")
		c.flagSet.StringVar(&c.DeleteAfter, "delete-after", "", "How long until schedule is removed in seconds, hybrid time will be calculated by fetching server hybrid time and adding this value")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotCreateScheduleConfig) Validate() error {
	if c.Keyspace == "" {
		return fmt.Errorf("--keyspace is required")
	}
	return nil
}
