package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpPingConfig represents a command specific config.
type OpPingConfig struct {
	flagBase

	Host string
	Port int
}

// NewOpPingConfig returns an instance of the command specific config.
func NewOpPingConfig() *OpPingConfig {
	return &OpPingConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpPingConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Host, "host", "", "Host to ping")
		c.flagSet.IntVar(&c.Port, "port", 0, "Port to ping")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpPingConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("--host is required")
	}
	if c.Port < 1 {
		return fmt.Errorf("--port is required")
	}
	return nil
}
