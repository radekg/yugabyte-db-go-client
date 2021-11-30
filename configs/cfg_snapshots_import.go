package configs

import "github.com/spf13/pflag"

// OpSnapshotImportConfig represents a command specific config.
type OpSnapshotImportConfig struct {
	flagBase

	FilePath  string
	Keyspace  string
	TableName []string
}

// NewOpSnapshotImportConfig returns an instance of the command specific config.
func NewOpSnapshotImportConfig() *OpSnapshotImportConfig {
	return &OpSnapshotImportConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSnapshotImportConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.FilePath, "file-path", "", "Absolute path to the snapshot export file, parent directories must exist")
		c.flagSet.StringVar(&c.Keyspace, "keyspace", "", "The name of the database or keyspace; YCQL only")
		c.flagSet.StringSliceVar(&c.TableName, "table-name", []string{}, "The name of the table; YCQL only")
	}
	return c.flagSet
}
