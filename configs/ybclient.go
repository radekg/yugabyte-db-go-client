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
