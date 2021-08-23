package configs

import (
	"crypto/tls"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

// CliConfig represents the shared CLI config.
type CliConfig struct {
	sync.Mutex
	flagBase

	MasterHostPort    string
	OpTimeout         time.Duration
	TLSCaCertFilePath string
	TLSCertFilePath   string
	TLSKeyFilePath    string
}

// NewCliConfig returns a new instance of the configuration.
func NewCliConfig() *CliConfig {
	return &CliConfig{}
}

// FlagSet returns an instance of the flag set for the configuration.
func (c *CliConfig) FlagSet() *pflag.FlagSet {
	if c.initFlagSet() {
		c.flagSet.StringVar(&c.MasterHostPort, "master", "127.0.0.1:7100", "Master host port")
		c.flagSet.DurationVar(&c.OpTimeout, "operation-timeout", time.Duration(time.Second*5), "Operation timeout")
		c.flagSet.StringVar(&c.TLSCaCertFilePath, "tls-ca-cert-file-path", "", "TLS CA certificate file path")
		c.flagSet.StringVar(&c.TLSCertFilePath, "tls-cert-file-path", "", "TLS certificate file path")
		c.flagSet.StringVar(&c.TLSKeyFilePath, "tls-key-file-path", "", "TLS key file path")
	}
	return c.flagSet
}

// TLSConfig returns TLS config is TLS is configured.
func (c *CliConfig) TLSConfig() *tls.Config {
	// TODO: implement
	return nil
}

// Validate validates the correctness of the configuration.
func (c *CliConfig) Validate() error {
	if c.MasterHostPort == "" {
		return fmt.Errorf("--master is required")
	}
	if c.TLSCertFilePath != "" && c.TLSKeyFilePath == "" {
		return fmt.Errorf("both --tls-cert-file-path and --tls-key-file-path are required")
	}
	if c.TLSKeyFilePath != "" && c.TLSCertFilePath == "" {
		return fmt.Errorf("both --tls-cert-file-path and --tls-key-file-path are required")
	}
	if c.TLSCaCertFilePath != "" {
		if c.TLSKeyFilePath == "" || c.TLSCertFilePath == "" {
			return fmt.Errorf("both --tls-cert-file-path and --tls-key-file-path are required when --tls-ca-cert-file-path")
		}
	}
	if c.OpTimeout.Milliseconds() < 0 {
		return fmt.Errorf("--operation-timeout must be greater than 0")
	}
	if c.OpTimeout.Milliseconds() > math.MaxUint32 {
		return fmt.Errorf("--operation-timeout is too large, cannot be greater than %d milliseconds", math.MaxUint32)
	}
	return nil
}
