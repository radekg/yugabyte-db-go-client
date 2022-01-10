package configs

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"
)

var (
	supportedNamespaceType = []string{"ycql", "ysql", "yedis"}
	supportedRelationType  = []string{"system_table", "user_table", "index_table"}
)

// YBSingleNodeClientConfig is the single node client configuration.
type YBSingleNodeClientConfig struct {
	MasterHostPort string
	TLSConfig      *tls.Config
	OpTimeout      uint32
}

const (
	// DefaultMaxExecuteRetries is the default maximum number of retries for a failed execute.
	DefaultMaxExecuteRetries int32 = 10
	// DefaultMaxReconnectAttempts is the default max reconnect attempts value.
	DefaultMaxReconnectAttempts int32 = 10
	// DefaultOpTimeout is the default operation timeout value.
	DefaultOpTimeout = time.Second * 60
	// DefaultReconnectRetryInterval is the default reconnect retry interval value.
	DefaultReconnectRetryInterval = time.Second
	// DefaultRetryInterval is the default retry interval value.
	DefaultRetryInterval = time.Second

	// NoExecuteRetry is a magic value disabling retry of failed execute.
	NoExecuteRetry int32 = -1
	// NoReconnectAttempts is a magic value disabling reconnect attempts.
	NoReconnectAttempts int32 = -1
)

// YBClientConfig represents the shared CLI config.
type YBClientConfig struct {
	tlsConfig *tls.Config

	MasterHostPort         []string
	OpTimeout              time.Duration
	MaxExecuteRetries      int32
	MaxReconnectAttempts   int32
	ReconnectRetryInterval time.Duration
	RetryInterval          time.Duration
	TLSCaCertFilePath      string
	TLSCertFilePath        string
	TLSKeyFilePath         string
}

// WithDefaults applies defaults to unset values.
func (c *YBClientConfig) WithDefaults() *YBClientConfig {
	if c.MaxExecuteRetries == 0 {
		c.MaxExecuteRetries = DefaultMaxExecuteRetries
	}
	if c.MaxReconnectAttempts == 0 {
		c.MaxReconnectAttempts = DefaultMaxReconnectAttempts
	}
	if c.OpTimeout == 0 {
		c.OpTimeout = DefaultOpTimeout
	}
	if c.ReconnectRetryInterval == 0 {
		c.ReconnectRetryInterval = DefaultReconnectRetryInterval
	}
	if c.RetryInterval == 0 {
		c.RetryInterval = DefaultRetryInterval
	}
	return c
}

// TLSConfig returns TLS config if TLS is configured.
func (c *YBClientConfig) TLSConfig() (*tls.Config, error) {

	if c.tlsConfig != nil {
		return c.tlsConfig, nil
	}

	if c.TLSCertFilePath != "" && c.TLSKeyFilePath != "" {
		cfg := &tls.Config{}
		cfg.RootCAs = x509.NewCertPool()
		if c.TLSCaCertFilePath != "" {
			caCertBytes, err := ioutil.ReadFile(c.TLSCaCertFilePath)
			if err != nil {
				return nil, fmt.Errorf("TLS configuration error, could not read ca cert file, reason: %v", err)
			}
			if ok := cfg.RootCAs.AppendCertsFromPEM(caCertBytes); !ok {
				return nil, fmt.Errorf("TLS configuration error, could append root certificates from file")
			}
		}
		cert, err := tls.LoadX509KeyPair(c.TLSCertFilePath, c.TLSKeyFilePath)
		if err != nil {
			return nil, fmt.Errorf("TLS configuration error, could not load X509 key pair, reason: %v", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
		c.tlsConfig = cfg
		return c.tlsConfig, nil
	}

	return nil, nil
}

// Validate validates the correctness of the configuration.
func (c *YBClientConfig) Validate() error {
	if len(c.MasterHostPort) == 0 {
		return fmt.Errorf("at least one --master is required")
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
	for _, path := range []string{c.TLSCertFilePath, c.TLSKeyFilePath, c.TLSCaCertFilePath} {
		if path != "" {
			fileInfo, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("TLS configuration error, file '%s' does not exist", err)
			}
			if fileInfo.IsDir() {
				return fmt.Errorf("TLS configuration error, path '%s' points at a directory, must be a file", err)
			}
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

func decodePem(certInput []byte) tls.Certificate {
	var cert tls.Certificate
	var certDERBlock *pem.Block
	for {
		certDERBlock, certInput = pem.Decode(certInput)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		}
	}
	return cert
}
