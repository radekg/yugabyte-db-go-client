package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpLeaderStepDownConfig represents a command specific config.
type OpLeaderStepDownConfig struct {
	flagBase

	DestUUID                 string
	DisableGracefulTansition bool
	TabletID                 string
	NewLeaderUUID            string
}

// NewOpLeaderStepDownConfig returns an instance of the command specific config.
func NewOpLeaderStepDownConfig() *OpLeaderStepDownConfig {
	return &OpLeaderStepDownConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpLeaderStepDownConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.DestUUID, "destination-uuid", "", "UUID of server this request is addressed to")
		c.flagSet.BoolVar(&c.DisableGracefulTansition, "disable-graceful-transition", false, "If new_leader_uuid is not specified, the current leader will attempt to gracefully transfer leadership to another peer. Setting this flag disables that behavior")
		c.flagSet.StringVar(&c.NewLeaderUUID, "new-leader-uuid", "", "UUID of the server that should run the election to become the new leader")
		c.flagSet.StringVar(&c.TabletID, "tablet-id", "", "The id of the tablet")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpLeaderStepDownConfig) Validate() error {
	if c.DestUUID == "" || c.TabletID == "" {
		return fmt.Errorf("--destination-uuid and --tablet-id required")
	}
	return nil
}
