package jupiter

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

const (
	// Name is the name of the bingx provider.
	Name = "jupiter_api"

	// URL is the base URL of the bingx API. This includes the base and quote
	// currency pairs that need to be inserted into the URL. This URL should be utilized
	// by Non-US users.
	URL = "https://price.jup.ag/v6/price?ids=%s&vsToken=%s"
)

// DefaultAPIConfig is the default configuration for the bingx API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           false,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         750 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}
