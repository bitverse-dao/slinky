package manager

import (
	"fmt"
	"math/big"
	"sync"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type (
	// ProviderManager is a stateful manager that is responsible for updating the
	// oracle and its providers with the latest market map changes.
	ProviderManager struct {
		mut    sync.Mutex
		logger *zap.Logger

		// providers is a map of all of the providers that the oracle is using.
		providers map[string]ProviderState

		// -------------------Oracle Configuration Fields-------------------//
		//
		// cfg is the oracle configuration.
		cfg config.OracleConfig
		// marketMap is the market map that the oracle is using.
		marketMap mmtypes.MarketMap

		// -------------------Provider Constructor Fields-------------------//
		//
		// apiQueryHandler factory is a factory function that creates API query handlers.
		apiQueryHandlerFactory types.PriceAPIQueryHandlerFactory
		// webSocketQueryHandlerFactory is a factory function that creates websocket query handlers.
		webSocketQueryHandlerFactory types.PriceWebSocketQueryHandlerFactory
		// wsMetrics is the web socket metrics.
		wsMetrics wsmetrics.WebSocketMetrics
		// apiMetrics is the API metrics.
		apiMetrics apimetrics.APIMetrics
		// providerMetrics is the provider metrics.
		providerMetrics providermetrics.ProviderMetrics
	}

	ProviderState struct {
		// provider is the price provider implementation.
		provider *types.PriceProvider
		// market is the market map for the provider.
		market types.ProviderMarketMap
		// enabled is a flag that indicates whether the provider is enabled.
		enabled bool
	}

	// Option is a functional option for the market map state.
	Option func(*ProviderManager)
)

// NewProviderManager returns a new provider manager.
func NewProviderManager(
	cfg config.OracleConfig,
	opts ...Option,
) (*ProviderManager, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	manager := &ProviderManager{
		cfg:       cfg,
		providers: make(map[string]ProviderState),
		logger:    zap.NewNop(),
	}

	for _, opt := range opts {
		opt(manager)
	}

	if err := manager.Init(); err != nil {
		return nil, err
	}

	return manager, nil
}

// Init initializes the all of the providers that are configured via the oracle config. Specifically,
// this will:
//
// 1. This will initialize the provider,
// 2. Create the provider specific market map, if configured with a marketmap.
// 3. Enable the provider if the provider is included in the oracle config and marketmap.
func (m *ProviderManager) Init() error {
	m.mut.Lock()
	defer m.mut.Unlock()

	for _, providerCfg := range m.cfg.Providers {
		// Initialize the provider.
		state, err := m.CreateProviderState(providerCfg)
		if err != nil {
			return fmt.Errorf("failed to create provider state: %w", err)
		}

		// Add the provider to the manager.
		m.providers[providerCfg.Name] = state
	}

	return nil
}

// CreateProviderState creates a provider state for the given provider.
func (m *ProviderManager) CreateProviderState(
	cfg config.ProviderConfig,
) (ProviderState, error) {
	// Create the provider state.
	state := ProviderState{}

	// Create the provider market map.
	market, err := types.ProviderMarketMapFromMarketMap(cfg.Name, m.marketMap)
	if err != nil {
		return ProviderState{}, fmt.Errorf("failed to create %s's provider market map: %w", cfg.Name, err)
	}

	state.market = market
	if len(market.GetTickers()) != 0 {
		state.enabled = true
	}

	// Create the provider.
	var (
		provider *types.PriceProvider
	)
	switch {
	case cfg.API.Enabled:
		if m.apiQueryHandlerFactory == nil {
			return ProviderState{}, fmt.Errorf("api query handler factory is not set")
		}

		queryHandler, err := m.apiQueryHandlerFactory(m.logger, cfg, m.apiMetrics, market)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's API query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[mmtypes.Ticker, *big.Int](cfg.Name),
			base.WithLogger[mmtypes.Ticker, *big.Int](m.logger),
			base.WithAPIQueryHandler(queryHandler),
			base.WithAPIConfig[mmtypes.Ticker, *big.Int](cfg.API),
			base.WithIDs[mmtypes.Ticker, *big.Int](market.GetTickers()),
			base.WithMetrics[mmtypes.Ticker, *big.Int](m.providerMetrics),
		)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	case cfg.WebSocket.Enabled:
		if m.webSocketQueryHandlerFactory == nil {
			return ProviderState{}, fmt.Errorf("web socket query handler factory is not set")
		}

		queryHandler, err := m.webSocketQueryHandlerFactory(m.logger, cfg, m.wsMetrics, market)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's web socket query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[mmtypes.Ticker, *big.Int](cfg.Name),
			base.WithLogger[mmtypes.Ticker, *big.Int](m.logger),
			base.WithWebSocketQueryHandler(queryHandler),
			base.WithWebSocketConfig[mmtypes.Ticker, *big.Int](cfg.WebSocket),
			base.WithIDs[mmtypes.Ticker, *big.Int](market.GetTickers()),
			base.WithMetrics[mmtypes.Ticker, *big.Int](m.providerMetrics),
		)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	default:
		return ProviderState{}, fmt.Errorf("provider %s has no enabled query handlers", cfg.Name)
	}

	state.provider = provider
	return state, nil
}

// GetProviders returns all of the providers that are configured on the manager. Specifically, this
// will return all of the providers that are enabled.
func (m *ProviderManager) GetProviders() []*types.PriceProvider {
	m.mut.Lock()
	defer m.mut.Unlock()

	providers := make([]*types.PriceProvider, 0, len(m.providers))
	for _, state := range m.providers {
		if state.enabled {
			providers = append(providers, state.provider)
		}
	}

	return providers
}
