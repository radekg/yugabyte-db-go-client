package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpSetLoadBalancerEnableConfig represents a command specific config.
type OpSetLoadBalancerEnableConfig struct {
	flagBase

	Enabled  bool
	Disabled bool
}

// NewOpSetLoadBalancerEnableConfig returns an instance of the command specific config.
func NewOpSetLoadBalancerEnableConfig() *OpSetLoadBalancerEnableConfig {
	return &OpSetLoadBalancerEnableConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSetLoadBalancerEnableConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.BoolVar(&c.Enabled, "enabled", false, "Desired state: enabled")
		c.flagSet.BoolVar(&c.Disabled, "disabled", false, "Desired state: disabled")
	}
	return c.flagSet
}

// IsEnabled returns the bool mapping of the state with an ok flag, which is false if the
// given value wss not an expected one.
func (c *OpSetLoadBalancerEnableConfig) IsEnabled() (bool, bool) {
	if c.Enabled {
		return true, true
	}
	if c.Disabled {
		return false, true
	}
	return false, false
}

// Validate validates the correctness of the configuration.
func (c *OpSetLoadBalancerEnableConfig) Validate() error {
	if !c.Enabled && !c.Disabled {
		return fmt.Errorf("--enabled or --disabled required")
	}
	if c.Enabled && c.Disabled {
		return fmt.Errorf("--enabled and --disabled: choose one")
	}
	return nil
}
