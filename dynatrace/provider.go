package dynatrace

import (
	"context"

	dynatraceConfigV1 "github.com/dynatrace-ace/dynatrace-go-api-client/api/v1/config/dynatrace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var dynatraceProvider *schema.Provider

// Provider function for Dynatrace API
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dt_env_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_ENV_URL", "DT_ENV_URL"}, nil),
			},
			"dt_api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_API_TOKEN", "DT_API_TOKEN"}, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_alerting_profiles": resourceDynatraceAlertingProfile(),
			"dynatrace_management_zones":  resourceDynatraceManagementZones(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dynatrace_alerting_profiles": dataSourceDynatraceAlertingProfiles(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

//ProviderConfiguration contains the initialized API clients to communicate with the Datadog API
type ProviderConfiguration struct {
	DynatraceConfigClientV1 *dynatraceConfigV1.APIClient
	AuthConfigV1            context.Context
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	dtEnvURL := d.Get("dt_env_url").(string)
	apiToken := d.Get("dt_api_token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Initialize the official Dynatrace Config V1 API client
	authConfigV1 := context.WithValue(
		context.Background(),
		dynatraceConfigV1.ContextAPIKey,
		dynatraceConfigV1.APIKey{
			Prefix: "Api-token",
			Key:    apiToken,
		},
	)

	configV1 := dynatraceConfigV1.NewConfiguration()
	configV1.BasePath = dtEnvURL + "/api/config/v1"

	dynatraceConfigClientV1 := dynatraceConfigV1.NewAPIClient(configV1)

	return &ProviderConfiguration{
		DynatraceConfigClientV1: dynatraceConfigClientV1,
		AuthConfigV1:            authConfigV1,
	}, diags

}
