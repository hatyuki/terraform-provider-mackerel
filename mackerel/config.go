package mackerel

import (
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type Config struct {
	ApiKey string
	RawURL string
}

func (c *Config) NewClient() (*mackerel.Client, error) {
	var verbose bool

	if logging.IsDebugOrHigher() {
		verbose = true
	}

	return mackerel.NewClientWithOptions(c.ApiKey, c.RawURL, verbose)
}
