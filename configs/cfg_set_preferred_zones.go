package configs

import (
	"github.com/spf13/pflag"
)

// OpSetPreferredZonesConfig represents a command specific config.
type OpSetPreferredZonesConfig struct {
	flagBase

	ZonesInfos []string
}

// NewOpSetPreferredZonesConfig returns an instance of the command specific config.
func NewOpSetPreferredZonesConfig() *OpSetPreferredZonesConfig {
	return &OpSetPreferredZonesConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *OpSetPreferredZonesConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringSliceVar(&c.ZonesInfos, "zone-info", []string{}, "Specifies the cloud, region, and zone. Default value is cloud1.datacenter1.rack1")
	}
	return c.flagSet
}
