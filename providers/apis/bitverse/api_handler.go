package bitverse

import (
	"fmt"
	"net/http"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for bitverse, which can be used
// by a base provider. The DataHandler fetches data from the spot price bitverse API. It is
// atomic in that it must request data from the bitverse API sequentially for each ticker.
type APIHandler struct {
	// api is the config for the bitverse API.
	api config.APIConfig
}

// NewAPIHandler returns a new bitverse PriceAPIDataHandler.
func NewAPIHandler(
	api config.APIConfig,
) (types.PriceAPIDataHandler, error) {
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
		api: api,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the bitverse API for the
// given tickers. Since the bitverse API only supports fetching spot prices for a single
// ticker at a time, this function will return an error if the ticker slice contains more
// than one ticker.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	if len(tickers) != 1 {
		return "", fmt.Errorf("expected 1 ticker, got %d", len(tickers))
	}
	return fmt.Sprintf(h.api.Endpoints[0].URL, tickers[0].GetOffChainTicker()), nil
}

// ParseResponse parses the spot price HTTP response from the bitverse API and returns
// the resulting price. Note that this can only parse a single ticker at a time.
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	if len(tickers) != 1 {
		return types.NewPriceResponseWithErr(tickers,
			providertypes.NewErrorWithCode(fmt.Errorf("expected 1 ticker, got %d", len(tickers)), providertypes.ErrorInvalidResponse),
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

	// Convert the float64 price into a big.Float.
	ticker := tickers[0]
	price, err := math.Float64StringToBigFloat(result.Tickers.LastPrice)
	if err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		)
	}

	return types.NewPriceResponse(
		types.ResolvedPrices{
			ticker: types.NewPriceResult(price, time.Now().UTC()),
		},
		nil,
	)
}
