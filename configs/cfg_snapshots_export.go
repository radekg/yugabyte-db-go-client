package configs

import "github.com/spf13/pflag"

// OpSnapshotExportConfig represents a command specific config.
type OpSnapshotExportConfig struct {
	flagBase

	SnapshotID    string
	Base64Encoded bool
	FilePath      string
}

// NewOpSnapshotExportConfig returns an instance of the command specific config.
func NewOpSnapshotExportConfig() *OpSnapshotExportConfig {
	return &OpSnapshotExportConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotExportConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.BoolVar(&c.Base64Encoded, "base64-encoded", false, "If true, accepts the --snapshot-id as base64 encoded string")
		c.flagSet.StringVar(&c.FilePath, "file-path", "", "Absolute path to the snapshot export file, parent directories must exist")
	}
	return c.flagSet
}
