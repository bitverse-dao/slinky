package bitmart

import (
	"fmt"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
)

// parseTickerResponseMessage parses a ticker response message. The format of the message is defined
// in the messages.go file. This message contains the latest price data for a set of instruments.
func (h *WebSocketHandler) parseTickerResponseMessage(
	resp TickersResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// The channel must be the index tickers channel.
	if Channel(resp.Table) != TickersChannel {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("invalid channel %s", resp.Table)
	}

	// Iterate through all tickers and add them to the response.
	for _, instrument := range resp.Data {
		ticker, ok := h.cache.FromOffChainTicker(instrument.Symbol)
		if !ok {
			h.logger.Debug("ticker not found for instrument ID", zap.String("instrument_id", instrument.Symbol))
			continue
		}

		// Convert the price to a big.Float.
		price, err := math.Float64StringToBigFloat(instrument.LastPrice)
		if err != nil {
			wErr := fmt.Errorf("failed to convert price to big.Float: %w", err)
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
			}
			continue
		}

		resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	}

	return types.NewPriceResponse(resolved, unresolved), nil
}
