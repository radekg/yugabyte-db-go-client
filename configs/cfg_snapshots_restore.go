package configs

import "github.com/spf13/pflag"

// OpSnapshotRestoreConfig represents a command specific config.
type OpSnapshotRestoreConfig struct {
	flagBase

	SnapshotID    string
	Base64Encoded bool
	RestoreHt     uint64
}

// NewOpSnapshotRestoreConfig returns an instance of the command specific config.
func NewOpSnapshotRestoreConfig() *OpSnapshotRestoreConfig {
	return &OpSnapshotRestoreConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotRestoreConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.BoolVar(&c.Base64Encoded, "base64-encoded", false, "If true, accepts the --snapshot-id as base64 encoded string")
		c.flagSet.Uint64Var(&c.RestoreHt, "restore-ht-micros", 0, "Absolute Timing Option: Max HybridTime, in Micros")
	}
	return c.flagSet
}
