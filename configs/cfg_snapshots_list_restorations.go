package configs

import "github.com/spf13/pflag"

// OpSnapshotListRestorationsConfig represents a command specific config.
type OpSnapshotListRestorationsConfig struct {
	flagBase

	SnapshotID    string
	RestorationID string
	Base64Encoded bool
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
		c.flagSet.BoolVar(&c.Base64Encoded, "base64-encoded", false, "If true, accepts the --snapshot-id as base64 encoded string")
	}
	return c.flagSet
}
