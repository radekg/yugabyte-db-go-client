package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpIsServerReadyConfig represents a command specific config.
type OpIsServerReadyConfig struct {
	flagBase

	Host      string
	Port      int
	IsTserver bool
}

// NewOpIsServerReadyConfig returns an instance of the command specific config.
func NewOpIsServerReadyConfig() *OpPingConfig {
	return &OpPingConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpIsServerReadyConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.Host, "host", "", "Host to check")
		c.flagSet.IntVar(&c.Port, "port", 0, "Port to check")
		c.flagSet.BoolVar(&c.IsTserver, "is-tserver", false, "Is TServer?")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpIsServerReadyConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("--host is required")
	}
	if c.Port < 1 {
		return fmt.Errorf("--port is required")
	}
	return nil
}
