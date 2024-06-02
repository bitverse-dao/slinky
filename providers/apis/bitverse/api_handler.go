package bitverse

import (
	"fmt"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"net/http"
	"time"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for Kraken.
// for more information about the Kraken API, refer to the following link:
type APIHandler struct {
	// market is the config for the Kraken API.
	market types.ProviderMarketMap
	// api is the config for the Kraken API.
	api config.APIConfig
}

// NewAPIHandler returns a new Kraken PriceAPIDataHandler.
func NewAPIHandler(
	market types.ProviderMarketMap,
	api config.APIConfig,
) (types.PriceAPIDataHandler, error) {
	if err := market.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market config for %s: %w", Name, err)
	}

	if market.Name != Name {
		return nil, fmt.Errorf("expected market config name %s, got %s", Name, market.Name)
	}

	if api.Name != Name {
		return nil, fmt.Errorf("expected api config name %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", Name)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config for %s: %w", Name, err)
	}

	return &APIHandler{
		market: market,
		api:    api,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the Bitverse API for the
// given tickers.
func (h *APIHandler) CreateURL(
	tickers []mmtypes.Ticker,
) (string, error) {
	if len(tickers) != 1 {
		return "", fmt.Errorf("expected 1 ticker, got %d", len(tickers))
	}

	// Ensure that the base and quote currencies are supported by the Coinbase API and
	// are configured for the handler.
	ticker := tickers[0]
	market, ok := h.market.TickerConfigs[ticker]
	if !ok {
		return "", fmt.Errorf("unknown ticker %s", ticker.String())
	}

	return fmt.Sprintf(h.api.URL, market.OffChainTicker), nil
}

// ParseResponse parses the response from the Kraken API and returns a GetResponse. Each
// of the tickers supplied will get a response or an error.
func (h *APIHandler) ParseResponse(
	tickers []mmtypes.Ticker,
	resp *http.Response,
) types.PriceResponse {
	if len(tickers) != 1 {
		return types.NewPriceResponseWithErr(tickers,
			providertypes.NewErrorWithCode(fmt.Errorf("expected 1 ticker, got %d", len(tickers)), providertypes.ErrorInvalidResponse),
		)
	}

	// Check if this ticker is supported by the Coinbase market config.
	ticker := tickers[0]
	_, ok := h.market.TickerConfigs[ticker]
	if !ok {
		return types.NewPriceResponseWithErr(tickers,
			providertypes.NewErrorWithCode(fmt.Errorf("unknown ticker %s", ticker.String()), providertypes.ErrorUnknownPair),
		)
	}

	// Parse the response into a ResponseBody.
	result, err := Decode(resp)
	if err != nil {
		return types.NewPriceResponseWithErr(tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	if result.Code != 200 {
		err := fmt.Errorf(
			"bitverse API call fail, code: %d", result.Code,
		)
		return types.NewPriceResponseWithErr(tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		)
	}

	price, err := math.Float64StringToBigInt(result.Tickers.LastPrice, ticker.Decimals)
	if err != nil {
		return types.NewPriceResponseWithErr(tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		)
	}

	return types.NewPriceResponse(
		types.ResolvedPrices{
			ticker: types.NewPriceResult(price, time.Now()),
		},
		nil,
	)
}
