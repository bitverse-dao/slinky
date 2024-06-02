package bitverse

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

// NOTE: All documentation for this file can be located on the Kraken docs.
// API documentation: https://docs.kraken.com/rest/. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Kraken API provider.
	Name = "bitverse_api"

	// URL is the base URL of the Kraken API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://market.bitverse-dev.bitverse.zone/api/v1/market/ticker?symbol=%s"
)

var (
	// DefaultAPIConfig is the default configuration for the Kraken API.
	DefaultAPIConfig = config.APIConfig{
		Name:             Name,
		Atomic:           true,
		Enabled:          true,
		Timeout:          500 * time.Millisecond,
		Interval:         400 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       1,
		URL:              URL,
	}

	// DefaultMarketConfig is the default market configuration for bitverse.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.APE_USD: {
			Name:           Name,
			OffChainTicker: "APE-USD",
		},
		constants.APE_USDC: {
			Name:           Name,
			OffChainTicker: "APE-USDC",
		},
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APE-USDT",
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
		constants.ATOM_USDC: {
			Name:           Name,
			OffChainTicker: "ATOM-USDC",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAX-USD",
		},
		constants.AVAX_USDC: {
			Name:           Name,
			OffChainTicker: "AVAX-USDC",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX-USDT",
		},
		constants.BCH_USD: {
			Name:           Name,
			OffChainTicker: "BCH-USD",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "BTC-USD",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC-USDT",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BTC-USDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC-USDT",
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
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "TIA-USDT",
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
		constants.DYDX_USD: {
			Name:           Name,
			OffChainTicker: "DYDX-USD",
		},
		constants.DYDX_USDC: {
			Name:           Name,
			OffChainTicker: "DYDX-USDC",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDX-USDT",
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
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETH-USDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETH-USDT",
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
		constants.OSMOSIS_USDT: {
			Name:           Name,
			OffChainTicker: "OSMO-USDT",
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
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOL-USDT",
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
		constants.USDC_USD: {
			Name:           Name,
			OffChainTicker: "USDC-USD",
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

// BitverseTicker is our representation of ticker information returned in Binance response.
// It implements interface `Ticker` in util.go.
type BitverseTicker struct {
	Pair      string `json:"symbol" validate:"required"`
	AskPrice  string `json:"indexPrice" validate:"required,positive-float-string"`
	BidPrice  string `json:"oraclePrice" validate:"required,positive-float-string"`
	LastPrice string `json:"lastPrice" validate:"required,positive-float-string"`
}

// ResponseBody returns a list of tickers for the response.  If there is an error, it will be included,
// and all Tickers will be undefined.
// ResponseBody defines the overall Huobi response.
type ResponseBody struct {
	Code    uint32         `json:"code" validate:"required"`
	Tickers BitverseTicker `json:"data" validate:"required"`
}

// Decode decodes the given http response into a TickerResult.
func Decode(resp *http.Response) (ResponseBody, error) {
	// Parse the response into a ResponseBody.
	var result ResponseBody
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
