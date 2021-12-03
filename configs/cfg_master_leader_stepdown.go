package configs

import (
	"github.com/spf13/pflag"
)

// OpMMasterLeaderStepdownConfig represents a command specific config.
type OpMMasterLeaderStepdownConfig struct {
	flagBase

	NewLeaderID string
}

// NewOpMMasterLeaderStepdownConfig returns an instance of the command specific config.
func NewOpMMasterLeaderStepdownConfig() *OpMMasterLeaderStepdownConfig {
	return &OpMMasterLeaderStepdownConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpMMasterLeaderStepdownConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.NewLeaderID, "new-leader-id", "", "The identifier (ID) of the new YB-Master leader. If not specified, the new leader is automatically elected, optional")
	}
	return c.flagSet
}
