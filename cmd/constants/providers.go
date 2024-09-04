package constants

import (
	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/constants"
	"github.com/skip-mev/connect/v2/oracle/types"
	binanceapi "github.com/skip-mev/connect/v2/providers/apis/binance"
	bingxapi "github.com/skip-mev/connect/v2/providers/apis/bingx"
	bitmartapi "github.com/skip-mev/connect/v2/providers/apis/bitmart"
	bitstampapi "github.com/skip-mev/connect/v2/providers/apis/bitstamp"
	coinbaseapi "github.com/skip-mev/connect/v2/providers/apis/coinbase"
	coinexapi "github.com/skip-mev/connect/v2/providers/apis/coinex"
	"github.com/skip-mev/connect/v2/providers/apis/coingecko"
	"github.com/skip-mev/connect/v2/providers/apis/coinmarketcap"
	"github.com/skip-mev/connect/v2/providers/apis/defi/osmosis"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium"
	"github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	jupiterapi "github.com/skip-mev/connect/v2/providers/apis/jupiter"
	krakenapi "github.com/skip-mev/connect/v2/providers/apis/kraken"
	"github.com/skip-mev/connect/v2/providers/apis/marketmap"
	"github.com/skip-mev/connect/v2/providers/apis/polymarket"
	xtapi "github.com/skip-mev/connect/v2/providers/apis/xt"
	"github.com/skip-mev/connect/v2/providers/apis/yymm"
	"github.com/skip-mev/connect/v2/providers/volatile"
	binancews "github.com/skip-mev/connect/v2/providers/websockets/binance"
	"github.com/skip-mev/connect/v2/providers/websockets/bitfinex"
	"github.com/skip-mev/connect/v2/providers/websockets/bitget"
	"github.com/skip-mev/connect/v2/providers/websockets/bitmart"
	"github.com/skip-mev/connect/v2/providers/websockets/bitstamp"
	"github.com/skip-mev/connect/v2/providers/websockets/bitverse"
	"github.com/skip-mev/connect/v2/providers/websockets/bybit"
	"github.com/skip-mev/connect/v2/providers/websockets/coinbase"
	"github.com/skip-mev/connect/v2/providers/websockets/cryptodotcom"
	"github.com/skip-mev/connect/v2/providers/websockets/gate"
	"github.com/skip-mev/connect/v2/providers/websockets/huobi"
	"github.com/skip-mev/connect/v2/providers/websockets/kraken"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
	"github.com/skip-mev/connect/v2/providers/websockets/mexc"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
	mmtypes "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
)

var (
	Providers = []config.ProviderConfig{
		// DEFI providers
		{
			Name: raydium.Name,
			API:  raydium.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: uniswapv3.ProviderNames[constants.ETHEREUM],
			API:  uniswapv3.DefaultETHAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: uniswapv3.ProviderNames[constants.BASE],
			API:  uniswapv3.DefaultBaseAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: osmosis.Name,
			API:  osmosis.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: pancakeswap.Name,
			API:  pancakeswap.DefaultAPIConfig,
			Type: types.ConfigType,
		},

		// Exchange API providers
		{
			Name: binanceapi.Name,
			API:  binanceapi.DefaultNonUSAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: bitstampapi.Name,
			API:  bitstampapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: coinbaseapi.Name,
			API:  coinbaseapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: coingecko.Name,
			API:  coingecko.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: coinmarketcap.Name,
			API:  coinmarketcap.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: krakenapi.Name,
			API:  krakenapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: bingxapi.Name,
			API:  bingxapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: coinexapi.Name,
			API:  coinexapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: xtapi.Name,
			API:  xtapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: bitmartapi.Name,
			API:  bitmartapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: jupiterapi.Name,
			API:  jupiterapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: volatile.Name,
			API:  volatile.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		// Exchange WebSocket providers
		{
			Name:      binancews.Name,
			WebSocket: binancews.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      bitfinex.Name,
			WebSocket: bitfinex.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      bitstamp.Name,
			WebSocket: bitstamp.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      bybit.Name,
			WebSocket: bybit.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      bitverse.Name,
			WebSocket: bitverse.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      coinbase.Name,
			WebSocket: coinbase.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      cryptodotcom.Name,
			WebSocket: cryptodotcom.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      gate.Name,
			WebSocket: gate.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      huobi.Name,
			WebSocket: huobi.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      kraken.Name,
			WebSocket: kraken.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      kucoin.Name,
			WebSocket: kucoin.DefaultWebSocketConfig,
			API:       kucoin.DefaultAPIConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      mexc.Name,
			WebSocket: mexc.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      okx.Name,
			WebSocket: okx.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},

		{
			Name:      bitget.Name,
			WebSocket: bitget.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},

		{
			Name:      bitmart.Name,
			WebSocket: bitmart.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},

		// Polymarket provider
		{
			Name: polymarket.Name,
			API:  polymarket.DefaultAPIConfig,
			Type: types.ConfigType,
		},

		// MarketMap provider
		{
			Name: marketmap.Name,
			API:  marketmap.DefaultAPIConfig,
			Type: mmtypes.ConfigType,
		},
	}

	AlternativeMarketMapProviders = []config.ProviderConfig{
		{
			Name: dydx.Name,
			API:  dydx.DefaultAPIConfig,
			Type: mmtypes.ConfigType,
		},
		{
			Name: dydx.SwitchOverAPIHandlerName,
			API:  dydx.DefaultSwitchOverAPIConfig,
			Type: mmtypes.ConfigType,
		},
		{
			Name: dydx.ResearchAPIHandlerName,
			API:  dydx.DefaultResearchAPIConfig,
			Type: mmtypes.ConfigType,
		},
		{
			Name: dydx.ResearchCMCAPIHandlerName,
			API:  dydx.DefaultResearchCMCAPIConfig,
			Type: mmtypes.ConfigType,
		},
		// yymm
		{
			Name: yymm.Name,
			API:  yymm.DefaultAPIConfig,
			Type: mmtypes.ConfigType,
		},
		{
			Name: yymm.SwitchOverAPIHandlerName,
			API:  yymm.DefaultSwitchOverAPIConfig,
			Type: mmtypes.ConfigType,
		},
	}

	MarketMapProviderNames = map[string]struct{}{
		dydx.Name:                      {},
		dydx.SwitchOverAPIHandlerName:  {},
		dydx.ResearchAPIHandlerName:    {},
		dydx.ResearchCMCAPIHandlerName: {},
		marketmap.Name:                 {},

		yymm.Name:                     {},
		yymm.SwitchOverAPIHandlerName: {},
	}
)
