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
func NewYBClientConfigFromCliConfig(hostPort string, input *CliConfig) (*YBClientConfig, error) {
	tlsConfig, err := input.TLSConfig()
	if err != nil {
		return nil, err
	}
	return &YBClientConfig{
		MasterHostPort: hostPort,
		TLSConfig:      tlsConfig,
		OpTimeout:      uint32(input.OpTimeout.Milliseconds()),
	}, nil
}
