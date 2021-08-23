package configs

import (
	"crypto/tls"
)

// YBClientConfig is the client configuration.
type YBClientConfig struct {
	MasterHostPort string
	TLSConfig      *tls.Config
	OpTimeout      uint32
}

// NewYBClientConfigFromCliConfig constructs YB client config from the cli config.
func NewYBClientConfigFromCliConfig(input *CliConfig) *YBClientConfig {
	return &YBClientConfig{
		MasterHostPort: input.MasterHostPort,
		TLSConfig:      input.TLSConfig(),
		OpTimeout:      uint32(input.OpTimeout.Milliseconds()),
	}
}
