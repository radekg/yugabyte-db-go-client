package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpSnapshotExportConfig represents a command specific config.
type OpSnapshotExportConfig struct {
	flagBase

	SnapshotID string
	FilePath   string
}

// NewOpSnapshotExportConfig returns an instance of the command specific config.
func NewOpSnapshotExportConfig() *OpSnapshotExportConfig {
	return &OpSnapshotExportConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotExportConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.SnapshotID, "snapshot-id", "", "Snapshot identifier")
		c.flagSet.StringVar(&c.FilePath, "file-path", "", "Absolute path to the snapshot export file, parent directories must exist")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpSnapshotExportConfig) Validate() error {
	if c.SnapshotID == "" {
		return fmt.Errorf("--snapshot-id is required")
	}
	return nil
}
