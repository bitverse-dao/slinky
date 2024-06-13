package bitverse

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

// NOTE: All documentation for this file can be located on the Kraken docs.
// API documentation: https://docs.kraken.com/rest/. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the bitverse API provider.
	Name = "bitverse_api"

	// URL is the base URL of the Kraken API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://market.bitverse-dev.bitverse.zone/api/v1/market/ticker?symbol=%s"

	URL_DEV = "https://market.testnet.bitverse.zone/api/v1/market/ticker?symbol=%s"
)

// DefaultAPIConfig is the default configuration for the Kraken API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         600 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}, {URL: URL_DEV}},
}

// Ticker is our representation of ticker information returned in Binance response.
// It implements interface `Ticker` in util.go.
type Ticker struct {
	Pair      string `json:"symbol" validate:"required"`
	AskPrice  string `json:"indexPrice" validate:"required,positive-float-string"`
	BidPrice  string `json:"oraclePrice" validate:"required,positive-float-string"`
	LastPrice string `json:"lastPrice" validate:"required,positive-float-string"`
}

// ResponseBody returns a list of tickers for the response.  If there is an error, it will be included,
// and all Tickers will be undefined.
// ResponseBody defines the overall Huobi response.
type ResponseBody struct {
	Code    uint32 `json:"code" validate:"required"`
	Tickers Ticker `json:"data" validate:"required"`
}

// Decode decodes the given http response into a TickerResult.
func Decode(resp *http.Response) (ResponseBody, error) {
	// Parse the response into a ResponseBody.
	var result ResponseBody
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
