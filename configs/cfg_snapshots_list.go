package configs

import "github.com/spf13/pflag"

// OpSnapshotListConfig represents a command specific config.
type OpSnapshotListConfig struct {
	flagBase

	SnapshotID           string
	ListDeletedSnapshots bool
	PrepareForBackup     bool
}

// NewOpSnapshotListConfig returns an instance of the command specific config.
func NewOpSnapshotListConfig() *OpSnapshotListConfig {
	return &OpSnapshotListConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotListConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.BoolVar(&c.ListDeletedSnapshots, "list-deleted-snapshots", false, "List deleted snapshots")
		c.flagSet.BoolVar(&c.PrepareForBackup, "prepare-for-backup", false, "Prepare for backup")
	}
	return c.flagSet
}
