package configs

import "github.com/spf13/pflag"

// OpListTabletServersConfig represents a command specific config.
type OpListTabletServersConfig struct {
	flagBase

	PrimaryOnly bool
}

// NewOpListTableServersConfig returns an instance of the command specific config.
func NewOpListTableServersConfig() *OpListTabletServersConfig {
	return &OpListTabletServersConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpListTabletServersConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.BoolVar(&c.PrimaryOnly, "primary-only", false, "Primary only")
	}
	return c.flagSet
}
