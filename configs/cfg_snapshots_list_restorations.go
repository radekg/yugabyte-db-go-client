package configs

import "github.com/spf13/pflag"

// OpSnapshotListRestorationsConfig represents a command specific config.
type OpSnapshotListRestorationsConfig struct {
	flagBase

	SnapshotID    string
	RestorationID string
}

// NewOpSnapshotListRestorationsConfig returns an instance of the command specific config.
func NewOpSnapshotListRestorationsConfig() *OpSnapshotListRestorationsConfig {
	return &OpSnapshotListRestorationsConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotListRestorationsConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.StringVar(&c.RestorationID, "restoration-id", "", "Restoration identifier")
	}
	return c.flagSet
}
