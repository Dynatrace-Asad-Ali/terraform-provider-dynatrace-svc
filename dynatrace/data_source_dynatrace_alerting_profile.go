package dynatrace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dynatraceConfigV1 "github.com/dynatrace-ace/dynatrace-go-api-client/api/v1/config/dynatrace"
)

func dataSourceDynatraceAlertingProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDynatraceAlertingProfilesRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"mz_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity_level": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_filters": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include_mode": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_filter": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"context": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"delay_in_minutes": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"event_type_filters": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"predefined_event_filter": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"negate": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"custom_event_filter": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_title_filter": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": &schema.Schema{
													Type:     schema.TypeBool,
													Computed: true,
												},
												"value": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"negate": &schema.Schema{
													Type:     schema.TypeBool,
													Computed: true,
												},
												"case_insensitive": &schema.Schema{
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"custom_description_filter": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": &schema.Schema{
													Type:     schema.TypeBool,
													Computed: true,
												},
												"value": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"negate": &schema.Schema{
													Type:     schema.TypeBool,
													Computed: true,
												},
												"case_insensitive": &schema.Schema{
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDynatraceAlertingProfilesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConf := m.(*ProviderConfiguration)
	dynatraceConfigClientV1 := providerConf.DynatraceConfigClientV1
	authConfigV1 := providerConf.AuthConfigV1

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	alertingProfileID := d.Get("id").(string)

	alertingProfile, _, err := dynatraceConfigClientV1.AlertingProfilesApi.GetAlertingProfile(authConfigV1, alertingProfileID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create dynatrace client",
			Detail:   "Bad Request or unable to connect to environment/authenticate API token",
		})
		return diags
	}

	alertingProfileRules := flattenAlertingProfileRulesData(&alertingProfile.Rules)
	if err := d.Set("rules", alertingProfileRules); err != nil {
		return diag.FromErr(err)
	}

	alertingProfileEventTypeFilters := flattenAlertingProfileEventTypeFiltersData(&alertingProfile.EventTypeFilters)
	if err := d.Set("event_type_filters", alertingProfileEventTypeFilters); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(alertingProfileID)
	d.Set("display_name", &alertingProfile.DisplayName)
	d.Set("mz_id", &alertingProfile.MzId)

	return diags

}

func flattenAlertingProfileRulesData(alertingProfileRules *[]dynatraceConfigV1.AlertingProfileSeverityRule) []interface{} {
	if alertingProfileRules != nil {
		ars := make([]interface{}, len(*alertingProfileRules), len(*alertingProfileRules))

		for i, alertingProfileRules := range *alertingProfileRules {
			ar := make(map[string]interface{})

			ar["severity_level"] = alertingProfileRules.SeverityLevel
			ar["delay_in_minutes"] = alertingProfileRules.DelayInMinutes
			ar["tag_filters"] = flattenAlertingProfileTagFilter(&alertingProfileRules.TagFilter)
			ars[i] = ar
		}

		return ars
	}

	return make([]interface{}, 0)
}

func flattenAlertingProfileTagFilter(alertingProfileTagFilter *dynatraceConfigV1.AlertingProfileTagFilter) []interface{} {
	if alertingProfileTagFilter == nil {
		return []interface{}{alertingProfileTagFilter}
	}
	t := make(map[string]interface{})

	t["include_mode"] = alertingProfileTagFilter.IncludeMode
	t["tag_filter"] = flattenAlertingProfileTagFilters(&alertingProfileTagFilter.TagFilters)

	return []interface{}{t}

}

func flattenAlertingProfileTagFilters(alertingProfileTagFilters *[]dynatraceConfigV1.TagFilter) []interface{} {

	if alertingProfileTagFilters != nil {
		tfs := make([]interface{}, len(*alertingProfileTagFilters), len(*alertingProfileTagFilters))

		for i, alertingProfileTagFilters := range *alertingProfileTagFilters {
			tf := make(map[string]interface{})

			tf["context"] = alertingProfileTagFilters.Context
			tf["key"] = alertingProfileTagFilters.Key
			tf["value"] = alertingProfileTagFilters.Value
			tfs[i] = tf
		}

		return tfs
	}

	return make([]interface{}, 0)
}

func flattenAlertingProfileEventTypeFiltersData(alertingProfileEventTypeFilters *[]dynatraceConfigV1.AlertingEventTypeFilter) []interface{} {
	if alertingProfileEventTypeFilters != nil {

		efs := make([]interface{}, len(*alertingProfileEventTypeFilters), len(*alertingProfileEventTypeFilters))

		for i, alertingProfileEventTypeFilters := range *alertingProfileEventTypeFilters {
			ef := make(map[string]interface{})

			ef["predefined_event_filter"] = flattenPredefinedEventFilter(alertingProfileEventTypeFilters.PredefinedEventFilter)
			ef["custom_event_filter"] = flattenCustomEventFilter(alertingProfileEventTypeFilters.CustomEventFilter)
			efs[i] = ef

		}
		return efs
	}

	return make([]interface{}, 0)
}

func flattenPredefinedEventFilter(alertingProfilePredefinedEventFilters *dynatraceConfigV1.AlertingPredefinedEventFilter) []interface{} {
	if alertingProfilePredefinedEventFilters == nil {
		return []interface{}{alertingProfilePredefinedEventFilters}
	}

	pef := make(map[string]interface{})

	pef["event_type"] = alertingProfilePredefinedEventFilters.EventType
	pef["negate"] = alertingProfilePredefinedEventFilters.Negate

	return []interface{}{pef}

}

func flattenCustomEventFilter(alertingProfileCustomEventFilters *dynatraceConfigV1.AlertingCustomEventFilter) []interface{} {
	if alertingProfileCustomEventFilters == nil {
		return []interface{}{alertingProfileCustomEventFilters}
	}

	cef := make(map[string]interface{})

	cef["custom_title_filter"] = flattenCustomTextFilter(alertingProfileCustomEventFilters.CustomTitleFilter)
	cef["custom_description_filter"] = flattenCustomTextFilter(alertingProfileCustomEventFilters.CustomDescriptionFilter)

	return []interface{}{cef}

}

func flattenCustomTextFilter(alertingProfileCustomTextFilters *dynatraceConfigV1.AlertingCustomTextFilter) []interface{} {
	if alertingProfileCustomTextFilters == nil {
		return []interface{}{alertingProfileCustomTextFilters}
	}

	ctf := make(map[string]interface{})

	ctf["enabled"] = alertingProfileCustomTextFilters.Enabled
	ctf["value"] = alertingProfileCustomTextFilters.Value
	ctf["operator"] = alertingProfileCustomTextFilters.Operator
	ctf["negate"] = alertingProfileCustomTextFilters.Negate
	ctf["case_insensitive"] = alertingProfileCustomTextFilters.CaseInsensitive

	return []interface{}{ctf}
}
