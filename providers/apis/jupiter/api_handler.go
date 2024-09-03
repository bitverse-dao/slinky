package jupiter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for bingx.
// for more information about the bingx API, refer to the following link:
type APIHandler struct {
	// api is the config for the bingx API.
	api config.APIConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewAPIHandler returns a new bingx PriceAPIDataHandler.
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
		api:   api,
		cache: types.NewProviderTickers(),
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the bingx API for the
// given tickers.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	if len(tickers) != 1 {
		return "", fmt.Errorf("expected 1 ticker, got %d", len(tickers))
	}
	ticker := tickers[0].GetOffChainTicker()
	t := strings.Split(ticker, "_")
	if len(t) != 2 {
		return "", fmt.Errorf("ticker incorrect format")
	}
	return fmt.Sprintf(h.api.Endpoints[0].URL, t[0], t[1]), nil
}

// ParseResponse parses the response from the bingx API and returns a GetResponse. Each
// of the tickers supplied will get a response or an error.
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	if len(tickers) != 1 {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("expected 1 ticker, got %d", len(tickers)),
				providertypes.ErrorInvalidResponse,
			),
		)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("data not exist"),
				providertypes.ErrorInvalidResponse,
			))
	}
	ticker := tickers[0]
	t := strings.Split(ticker.GetOffChainTicker(), "_")
	if len(t) != 2 {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("ticker incorrect format"),
				providertypes.ErrorInvalidResponse,
			))
	}

	denom, ok := data[t[0]].(map[string]interface{})
	if !ok {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("%s not exist", t[0]),
				providertypes.ErrorInvalidResponse,
			))
	}
	priceFloat, ok := denom["price"].(float64)
	if !ok {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("price not exist"),
				providertypes.ErrorInvalidResponse,
			))
	}
	priceStr := strconv.FormatFloat(priceFloat, 'f', -1, 64)
	price, err := math.Float64StringToBigFloat(priceStr)
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
