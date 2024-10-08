package bingx

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// Name is the name of the bingx provider.
	Name = "bingx_api"

	// URL is the base URL of the bingx API. This includes the base and quote
	// currency pairs that need to be inserted into the URL. This URL should be utilized
	// by Non-US users.
	URL = "https://open-api.bingx.com/openApi/spot/v1/ticker/price?symbol=%s"
)

// DefaultAPIConfig is the default configuration for the bingx API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         750 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}

type (
	BingxResponse struct {
		Code int         `json:"code"`
		Data []BingxData `json:"data"`
	}

	BingxData struct { //nolint
		Symbol string `json:"symbol"`
		Trades []struct {
			Price string `json:"price"`
		} `json:"trades"`
	}
)
