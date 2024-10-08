package yymm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/yymm"
	apihandlermocks "github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	mmclient "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestYYMMMultiMarketMapFetcher(t *testing.T) {
	yymmMainnetMMFetcher := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)
	yymmResearchMMFetcher := apihandlermocks.NewAPIFetcher[mmclient.Chain, *mmtypes.MarketMapResponse](t)

	fetcher := yymm.NewYYMMResearchMarketMapFetcher(yymmMainnetMMFetcher, yymmResearchMMFetcher, zap.NewExample(), false)

	t.Run("test that if the mainnet api-price fetcher response is unresolved, we return it", func(t *testing.T) {
		ctx := context.Background()
		yymmMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{
			UnResolved: map[mmclient.Chain]providertypes.UnresolvedResult{
				yymm.YYMMChain: {
					ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("error"), providertypes.ErrorAPIGeneral),
				},
			},
		}, nil).Once()
		yymmResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{yymm.YYMMChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test that if the yymm-research response is unresolved, we return that", func(t *testing.T) {
		ctx := context.Background()
		yymmMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				yymm.YYMMChain: providertypes.NewResult(&mmtypes.MarketMapResponse{}, time.Now()),
			},
		}, nil).Once()
		yymmResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{
			UnResolved: map[mmclient.Chain]providertypes.UnresolvedResult{
				yymm.YYMMChain: {},
			},
		}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{yymm.YYMMChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test if both responses are resolved, the tickers are appended to each other + validation fails", func(t *testing.T) {
		ctx := context.Background()
		yymmMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				yymm.YYMMChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"BTC/USD": {},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()
		yymmResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				yymm.YYMMChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"ETH/USD": {},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{yymm.YYMMChain})
		require.Len(t, response.UnResolved, 1)
	})

	t.Run("test that if both responses are resolved, the responses are aggregated + validation passes", func(t *testing.T) {
		ctx := context.Background()
		yymmMainnetMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				yymm.YYMMChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"BTC/USD": {
								Ticker: mmtypes.Ticker{
									CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USD"),
									Decimals:         8,
									MinProviderCount: 1,
									Enabled:          true,
								},
								ProviderConfigs: []mmtypes.ProviderConfig{
									{
										Name:           "yymm",
										OffChainTicker: "BTC/USD",
									},
								},
							},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()
		yymmResearchMMFetcher.On("Fetch", ctx, []mmclient.Chain{yymm.YYMMChain}).Return(mmclient.MarketMapResponse{
			Resolved: map[mmclient.Chain]providertypes.ResolvedResult[*mmtypes.MarketMapResponse]{
				yymm.YYMMChain: providertypes.NewResult(&mmtypes.MarketMapResponse{
					MarketMap: mmtypes.MarketMap{
						Markets: map[string]mmtypes.Market{
							"ETH/USD": {
								Ticker: mmtypes.Ticker{
									CurrencyPair:     slinkytypes.NewCurrencyPair("ETH", "USD"),
									Decimals:         8,
									MinProviderCount: 1,
								},
								ProviderConfigs: []mmtypes.ProviderConfig{
									{
										Name:           "yymm",
										OffChainTicker: "BTC/USD",
									},
								},
							},
						},
					},
				}, time.Now()),
			},
		}, nil).Once()

		response := fetcher.Fetch(ctx, []mmclient.Chain{yymm.YYMMChain})
		require.Len(t, response.Resolved, 1)

		marketMap := response.Resolved[yymm.YYMMChain].Value.MarketMap

		require.Len(t, marketMap.Markets, 2)
		require.Contains(t, marketMap.Markets, "BTC/USD")
		require.Contains(t, marketMap.Markets, "ETH/USD")
	})
}
