package mackerel

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	defaultBaseURL = "https://api.mackerelio.com/"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MACKEREL_APIKEY", nil),
				Description: `The API key of the organization to which targeted hosts and services belong.`,
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MACKEREL_BASE_URL", defaultBaseURL),
				Description: ``,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"mackerel_dashboard": resourceMackerelDashboard(),
			//"mackerel_expression_monitor": resourceMackerelExpressionMonitor(),
			//"mackerel_external_monitor":   resourceMackerelExternalMonitor(),
			//"mackerel_host_monitor":       resourceMackerelHostMonitor(),
			//"mackerel_service_monitor":    resourceMackerelServiceMonitor(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		ApiKey: d.Get("api_key").(string),
		RawURL: d.Get("base_url").(string),
	}

	return config.NewClient()
}
