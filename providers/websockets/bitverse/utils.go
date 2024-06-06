package bitverse

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// ByBit provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://bybit-exchange.github.io/docs/v5/ws/connect
	// The two production URLs are defined in ProductionURL and TestnetURL.

	// Name is the name of the ByBit provider.
	Name = "bitverse_ws"

	// URLProd is the public ByBit Websocket URL.
	URLProd = "wss://public-stream.testnet.bitverse.zone"

	// URLTest is the public testnet ByBit Websocket URL.
	URLTest = "wss://stream-testnet.bybit.com/v5/public/spot"

	// DefaultPingInterval is the default ping interval for the ByBit websocket.
	DefaultPingInterval = 15 * time.Second
)

var (
	// DefaultWebSocketConfig is the default configuration for the ByBit Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           URLProd,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for ByBit.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.APE_USD: {
			Name:           Name,
			OffChainTicker: "APE-USD",
		},
		constants.APTOS_USD: {
			Name:           Name,
			OffChainTicker: "APT-USD",
		},
		constants.ARBITRUM_USD: {
			Name:           Name,
			OffChainTicker: "ARB-USD",
		},
		constants.ATOM_USD: {
			Name:           Name,
			OffChainTicker: "ATOM-USD",
		},

		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAX-USD",
		},

		constants.BCH_USD: {
			Name:           Name,
			OffChainTicker: "BCH-USD",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "BTC-USD",
		},

		constants.BLUR_USD: {
			Name:           Name,
			OffChainTicker: "BLUR-USD",
		},
		constants.CARDANO_USD: {
			Name:           Name,
			OffChainTicker: "ADA-USD",
		},
		constants.CELESTIA_USD: {
			Name:           Name,
			OffChainTicker: "TIA-USD",
		},
		constants.CELESTIA_USDC: {
			Name:           Name,
			OffChainTicker: "TIA-USDC",
		},

		constants.CHAINLINK_USD: {
			Name:           Name,
			OffChainTicker: "LINK-USD",
		},
		constants.COMPOUND_USD: {
			Name:           Name,
			OffChainTicker: "COMP-USD",
		},
		constants.CURVE_USD: {
			Name:           Name,
			OffChainTicker: "CRV-USD",
		},
		constants.DOGE_USD: {
			Name:           Name,
			OffChainTicker: "DOGE-USD",
		},
		constants.ETC_USD: {
			Name:           Name,
			OffChainTicker: "ETC-USD",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "ETH-USD",
		},

		constants.FILECOIN_USD: {
			Name:           Name,
			OffChainTicker: "FIL-USD",
		},
		constants.LIDO_USD: {
			Name:           Name,
			OffChainTicker: "LDO-USD",
		},
		constants.LITECOIN_USD: {
			Name:           Name,
			OffChainTicker: "LTC-USD",
		},
		constants.MAKER_USD: {
			Name:           Name,
			OffChainTicker: "MKR-USD",
		},
		constants.NEAR_USD: {
			Name:           Name,
			OffChainTicker: "NEAR-USD",
		},
		constants.OPTIMISM_USD: {
			Name:           Name,
			OffChainTicker: "OP-USD",
		},
		constants.OSMOSIS_USD: {
			Name:           Name,
			OffChainTicker: "OSMO-USD",
		},
		constants.OSMOSIS_USDC: {
			Name:           Name,
			OffChainTicker: "OSMO-USDC",
		},

		constants.POLKADOT_USD: {
			Name:           Name,
			OffChainTicker: "DOT-USD",
		},
		constants.POLYGON_USD: {
			Name:           Name,
			OffChainTicker: "MATIC-USD",
		},
		constants.RIPPLE_USD: {
			Name:           Name,
			OffChainTicker: "XRP-USD",
		},
		constants.SEI_USD: {
			Name:           Name,
			OffChainTicker: "SEI-USD",
		},
		constants.SHIBA_USD: {
			Name:           Name,
			OffChainTicker: "SHIB-USD",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "SOL-USD",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOL-USDC",
		},

		constants.STELLAR_USD: {
			Name:           Name,
			OffChainTicker: "XLM-USD",
		},
		constants.SUI_USD: {
			Name:           Name,
			OffChainTicker: "SUI-USD",
		},
		constants.UNISWAP_USD: {
			Name:           Name,
			OffChainTicker: "UNI-USD",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC-USDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDT-USD",
		},
	}
)
