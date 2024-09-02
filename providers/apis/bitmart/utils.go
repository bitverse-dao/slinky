package bitmart

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// Name is the name of the bingx provider.
	Name = "bitmart_api"

	// URL is the base URL of the bingx API. This includes the base and quote
	// currency pairs that need to be inserted into the URL. This URL should be utilized
	// by Non-US users.
	URL = "https://api-cloud.bitmart.com/spot/quotation/v3/ticker?symbol=%s"
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
	BitmartResponse struct {
		Code int         `json:"code"`
		Data BitmartData `json:"data"`
	}

	BitmartData struct { //nolint
		Symbol string `json:"symbol"`
		Price  string `json:"last"`
	}
)
