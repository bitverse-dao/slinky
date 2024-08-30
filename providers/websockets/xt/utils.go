package xt

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (

	// Name is the name of the OKX provider.
	Name = "xt_ws"

	// URL_PROD is the public OKX Websocket URL.
	URL_PROD = "wss://stream.xt.com/public"

	WriteInterval = 3000 * time.Millisecond

	// DefaultPingInterval is the default ping interval for the bitstamp websocket.
	DefaultPingInterval = 25 * time.Second
)

// DefaultWebSocketConfig is the default configuration for the OKX Websocket.
var DefaultWebSocketConfig = config.WebSocketConfig{
	Name:                          Name,
	Enabled:                       true,
	MaxBufferSize:                 config.DefaultMaxBufferSize,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: URL_PROD}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   config.DefaultReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  DefaultPingInterval,
	WriteInterval:                 WriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
}
