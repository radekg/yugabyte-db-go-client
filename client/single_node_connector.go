package client

import (
	"crypto/tls"
	"net"

	"github.com/hashicorp/go-hclog"
	"github.com/radekg/yugabyte-db-go-client/configs"
	"github.com/radekg/yugabyte-db-go-client/metrics"
)

// Connector connects a single node client.
type Connector interface {
	Connect(cfg *configs.YBSingleNodeClientConfig) (YBConnectedClient, error)
	// Allows configuring the logger used by the client.
	// Uses go-hclog. Users can provide integrate with any logging
	// framework using https://pkg.go.dev/github.com/hashicorp/go-hclog#InterceptLogger.
	WithLogger(logger hclog.Logger) Connector
	// Allows providing custom implementation of the metrics callback.
	WithMetricsCallback(callback metrics.Callback) Connector
}

type defaultClientConnector struct {
	logger          hclog.Logger
	metricsCallback metrics.Callback
}

// NewDefaultConnector returns a new instance of the default connector.
func NewDefaultConnector() Connector {
	return &defaultClientConnector{
		logger:          hclog.Default(),
		metricsCallback: metrics.Noop(),
	}
}

// WithLogger configures the logger for the connector and resulting client.
func (dcc *defaultClientConnector) WithLogger(logger hclog.Logger) Connector {
	if logger != nil {
		dcc.logger = logger
	}
	return dcc
}

// WithMetricsCallback configures the metrics callback for the connector and resulting client.
func (dcc *defaultClientConnector) WithMetricsCallback(callback metrics.Callback) Connector {
	if callback != nil {
		dcc.metricsCallback = callback
	}
	return dcc
}

// Connect connects to the master server without TLS.
func (dcc *defaultClientConnector) Connect(cfg *configs.YBSingleNodeClientConfig) (YBConnectedClient, error) {
	if cfg.TLSConfig != nil {
		return dcc.connectTLS(cfg)
	}
	return dcc.connect(cfg)
}

func (dcc *defaultClientConnector) connect(cfg *configs.YBSingleNodeClientConfig) (YBConnectedClient, error) {
	dcc.logger.Debug("connecting non-TLS client")
	conn, err := net.Dial("tcp", cfg.MasterHostPort)
	if err != nil {
		return nil, err
	}
	client := &defaultSingleNodeClient{
		originalConfig: cfg,
		chanConnected:  make(chan struct{}, 1),
		chanConnectErr: make(chan error, 1),
		closeFunc: func() error {
			return conn.Close()
		},
		conn:        conn,
		logger:      dcc.logger,
		svcRegistry: NewDefaultServiceRegistry(),
	}
	return client.
		withLogger(dcc.logger).
		withMetricsCallback(dcc.metricsCallback).
		afterConnect(), nil
}

func (dcc *defaultClientConnector) connectTLS(cfg *configs.YBSingleNodeClientConfig) (YBConnectedClient, error) {
	dcc.logger.Debug("connecting TLS client")
	conn, err := tls.Dial("tcp", cfg.MasterHostPort, cfg.TLSConfig)
	if err != nil {
		return nil, err
	}
	client := &defaultSingleNodeClient{
		originalConfig: cfg,
		chanConnected:  make(chan struct{}, 1),
		chanConnectErr: make(chan error, 1),
		closeFunc: func() error {
			return conn.Close()
		},
		conn:        conn,
		svcRegistry: NewDefaultServiceRegistry(),
	}
	return client.
		withLogger(dcc.logger).
		withMetricsCallback(dcc.metricsCallback).
		afterConnect(), nil
}
