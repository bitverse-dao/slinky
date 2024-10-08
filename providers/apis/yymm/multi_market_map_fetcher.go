package yymm

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/cmd/constants/marketmaps"
	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/coinmarketcap"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	mmclient "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	_         mmclient.MarketMapFetcher = &MultiMarketMapRestAPIFetcher{}
	YYMMChain                           = mmclient.Chain{
		ChainID: ChainID,
	}
)

// NewYYMMResearchMarketMapFetcher returns a MultiMarketMapFetcher composed of yymm mainnet + research
// apiDataHandlers.
func DefaultYYMMResearchMarketMapFetcher(
	rh apihandlers.RequestHandler,
	metrics metrics.APIMetrics,
	api config.APIConfig,
	logger *zap.Logger,
) (*MultiMarketMapRestAPIFetcher, error) {
	if rh == nil {
		return nil, fmt.Errorf("request handler is nil")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics is nil")
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, err
	}

	if len(api.Endpoints) != 2 {
		return nil, fmt.Errorf("expected two endpoint, got %d", len(api.Endpoints))
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	// make a yymm research api-handler
	researchAPIDataHandler, err := NewResearchAPIHandler(logger, api)
	if err != nil {
		return nil, err
	}

	mainnetAPIDataHandler := &APIHandler{
		logger: logger,
		api:    api,
	}

	mainnetFetcher, err := apihandlers.NewRestAPIFetcher(
		rh,
		mainnetAPIDataHandler,
		metrics,
		api,
		logger,
	)
	if err != nil {
		return nil, err
	}

	researchFetcher, err := apihandlers.NewRestAPIFetcher(
		rh,
		researchAPIDataHandler,
		metrics,
		api,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return NewYYMMResearchMarketMapFetcher(
		mainnetFetcher,
		researchFetcher,
		logger,
		api.Name == ResearchCMCAPIHandlerName,
	), nil
}

// MultiMarketMapRestAPIFetcher is an implementation of a RestAPIFetcher that wraps
// two underlying Fetchers for fetching the market-map according to yymm mainnet and
// the additional markets that can be added according to the yymm research json.
type MultiMarketMapRestAPIFetcher struct {
	// yymm mainnet fetcher is the api-fetcher for the yymm mainnet market-map
	yymmMainnetFetcher mmclient.MarketMapFetcher

	// yymm research fetcher is the api-fetcher for the yymm research market-map
	yymmResearchFetcher mmclient.MarketMapFetcher

	// logger is the logger for the fetcher
	logger *zap.Logger

	// isCMCOnly is a flag that indicates whether the fetcher should only return CoinMarketCap markets.
	isCMCOnly bool
}

// NewYYMMResearchMarketMapFetcher returns an aggregated market-map among the yymm mainnet and the yymm research json.
func NewYYMMResearchMarketMapFetcher(
	mainnetFetcher, researchFetcher mmclient.MarketMapFetcher,
	logger *zap.Logger,
	isCMCOnly bool,
) *MultiMarketMapRestAPIFetcher {
	return &MultiMarketMapRestAPIFetcher{
		yymmMainnetFetcher:  mainnetFetcher,
		yymmResearchFetcher: researchFetcher,
		logger:              logger.With(zap.String("module", "yymm-research-market-map-fetcher")),
		isCMCOnly:           isCMCOnly,
	}
}

// Fetch fetches the market map from the underlying fetchers and combines the results. If any of the underlying
// fetchers fetch for a chain that is different from the chain that the fetcher is initialized with, those responses
// will be ignored.
func (f *MultiMarketMapRestAPIFetcher) Fetch(ctx context.Context, chains []mmclient.Chain) mmclient.MarketMapResponse {
	// call the underlying fetchers + await their responses
	// channel to aggregate responses
	yymmMainnetResponseChan := make(chan mmclient.MarketMapResponse, 1) // buffer so that sends / receives are non-blocking
	yymmResearchResponseChan := make(chan mmclient.MarketMapResponse, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	// fetch yymm mainnet
	go func() {
		defer wg.Done()
		yymmMainnetResponseChan <- f.yymmMainnetFetcher.Fetch(ctx, chains)
		f.logger.Debug("fetched valid market-map from yymm mainnet")
	}()

	// fetch yymm research
	go func() {
		defer wg.Done()
		yymmResearchResponseChan <- f.yymmResearchFetcher.Fetch(ctx, chains)
		f.logger.Debug("fetched valid market-map from yymm research")
	}()

	// wait for both fetchers to finish
	wg.Wait()

	yymmMainnetMarketMapResponse := <-yymmMainnetResponseChan
	yymmResearchMarketMapResponse := <-yymmResearchResponseChan

	// if the yymm mainnet market-map response failed, return the yymm mainnet failed response
	if _, ok := yymmMainnetMarketMapResponse.UnResolved[YYMMChain]; ok {
		f.logger.Error("yymm mainnet market-map fetch failed", zap.Any("response", yymmMainnetMarketMapResponse))
		return yymmMainnetMarketMapResponse
	}

	// if the yymm research market-map response failed, return the yymm research failed response
	if _, ok := yymmResearchMarketMapResponse.UnResolved[YYMMChain]; ok {
		f.logger.Error("yymm research market-map fetch failed", zap.Any("response", yymmResearchMarketMapResponse))
		return yymmResearchMarketMapResponse
	}

	// otherwise, add all markets from yymm research
	yymmMainnetMarketMap := yymmMainnetMarketMapResponse.Resolved[YYMMChain].Value.MarketMap

	resolved, ok := yymmResearchMarketMapResponse.Resolved[YYMMChain]
	if ok {
		for ticker, market := range resolved.Value.MarketMap.Markets {
			// if the market is not already in the yymm mainnet market-map, add it
			if _, ok := yymmMainnetMarketMap.Markets[ticker]; !ok {
				f.logger.Debug("adding market from yymm research", zap.String("ticker", ticker))
				yymmMainnetMarketMap.Markets[ticker] = market
			}
		}
	}

	// if the fetcher is only for CoinMarketCap markets, filter out all non-CMC markets
	if f.isCMCOnly {
		for ticker, market := range yymmMainnetMarketMap.Markets {
			market.Ticker.MinProviderCount = 1
			yymmMainnetMarketMap.Markets[ticker] = market

			var (
				seenCMC     = false
				cmcProvider mmtypes.ProviderConfig
			)

			for _, provider := range market.ProviderConfigs {
				if provider.Name == coinmarketcap.Name {
					seenCMC = true
					cmcProvider = provider
				}
			}

			// if we saw a CMC provider, add it to the market
			if seenCMC {
				market.ProviderConfigs = []mmtypes.ProviderConfig{cmcProvider}
				yymmMainnetMarketMap.Markets[ticker] = market
				continue
			}

			// If we did not see a CMC provider, we can attempt to add it using the CMC marketmap
			cmcMarket, ok := marketmaps.CoinMarketCapMarketMap.Markets[ticker]
			if !ok {
				f.logger.Info("did not find CMC market for ticker", zap.String("ticker", ticker))
				delete(yymmMainnetMarketMap.Markets, ticker)
				continue
			}

			// add the CMC provider to the market
			market.ProviderConfigs = cmcMarket.ProviderConfigs
			yymmMainnetMarketMap.Markets[ticker] = market
		}
	}

	// validate the combined market-map
	if err := yymmMainnetMarketMap.ValidateBasic(); err != nil {
		f.logger.Error("combined market-map failed validation", zap.Error(err))

		return mmclient.NewMarketMapResponseWithErr(
			chains,
			providertypes.NewErrorWithCode(
				fmt.Errorf("combined market-map failed validation: %w", err),
				providertypes.ErrorUnknown,
			),
		)
	}

	yymmMainnetMarketMapResponse.Resolved[YYMMChain].Value.MarketMap = yymmMainnetMarketMap

	return yymmMainnetMarketMapResponse
}
