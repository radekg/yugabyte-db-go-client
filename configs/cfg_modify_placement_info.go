package configs

import (
	"fmt"

	"github.com/spf13/pflag"
)

// OpModifyPlacementInfoConfig represents a command specific config.
type OpModifyPlacementInfoConfig struct {
	flagBase

	PlacementInfos    []string
	ReplicationFactor uint32
	PlacementUUID     string
}

// NewOpModifyPlacementInfoConfig returns an instance of the command specific config.
func NewOpModifyPlacementInfoConfig() *OpModifyPlacementInfoConfig {
	return &OpModifyPlacementInfoConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpModifyPlacementInfoConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringSliceVar(&c.PlacementInfos, "placement-info", []string{}, "Placement for cloud.region.zone. Default value is cloud1.datacenter1.rack1.")
		c.flagSet.Uint32Var(&c.ReplicationFactor, "replication-factor", 0, "The number of replicas for each tablet.")
		c.flagSet.StringVar(&c.PlacementUUID, "placement-uuid", "", "The identifier of the primary cluster, which can be any unique string. Optional. If not set, a randomly-generated ID will be used.")
	}
	return c.flagSet
}

// Validate validates the correctness of the configuration.
func (c *OpModifyPlacementInfoConfig) Validate() error {
	if len(c.PlacementInfos) == 0 {
		return fmt.Errorf("at least one --placement-info required")
	}
	if c.ReplicationFactor == 0 {
		return fmt.Errorf("--replication-factor required")
	}
	return nil
}
