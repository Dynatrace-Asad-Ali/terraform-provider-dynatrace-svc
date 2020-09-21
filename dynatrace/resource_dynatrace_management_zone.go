package dynatrace

import (
	"context"

	"github.com/antihax/optional"
	dynatraceConfigV1 "github.com/dynatrace-ace/dynatrace-go-api-client/api/v1/config/dynatrace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDynatraceManagementZones() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDynatraceManagementZoneCreate,
		ReadContext:   resourceDynatraceManagementZoneRead,
		UpdateContext: resourceDynatraceManagementZoneUpdate,
		DeleteContext: resourceDynatraceManagementZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the management zone.",
				Required:    true,
			},
			"rule": &schema.Schema{
				Type:        schema.TypeList,
				Description: "A list of rules for management zone usage. Each rule is evaluated independently of all other rules.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The type of Dynatrace entities the management zone can be applied to.",
							Required:    true,
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "The rule is enabled (true) or disabled (false).",
							Required:    true,
						},
						"propagation_types": &schema.Schema{
							Type:        schema.TypeSet,
							Description: "How to apply the management zone to underlying entities.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"condition": &schema.Schema{
							Type:        schema.TypeList,
							Description: "A list of matching rules for the management zone. The management zone applies only if all conditions are fulfilled.",
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:        schema.TypeList,
										Description: "The key to identify the data we're matching.",
										Required:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"attribute": &schema.Schema{
													Type:        schema.TypeString,
													Description: "The attribute to be used for comparision.",
													Required:    true,
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Description: "Defines the actual set of fields depending on the value.",
													Optional:    true,
												},
											},
										},
									},
									"comparison_info": &schema.Schema{
										Type:        schema.TypeList,
										Description: "Defines how the matching is actually performed: what and how are we comparing.",
										Required:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": &schema.Schema{
													Type:        schema.TypeString,
													Description: "Operator of the comparison. You can reverse it by setting negate to true. Possible values depend on the type of the comparison. Find the list of actual models in the description of the type field and check the description of the model you need.",
													Required:    true,
												},
												"value": {
													Type:        schema.TypeMap,
													Description: "The value to compare to.",
													Optional:    true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"negate": &schema.Schema{
													Type:        schema.TypeBool,
													Description: "Reverses the comparison operator. For example it turns the begins with into does not begin with.",
													Required:    true,
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Description: "Defines the actual set of fields depending on the value.",
													Required:    true,
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

func resourceDynatraceManagementZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConf := m.(*ProviderConfiguration)
	dynatraceConfigClientV1 := providerConf.DynatraceConfigClientV1
	authConfigV1 := providerConf.AuthConfigV1

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	mz := dynatraceConfigV1.ManagementZone{
		Name:  d.Get("name").(string),
		Rules: expandManagementZoneRules(d.Get("rule").([]interface{})),
	}

	mzBody := dynatraceConfigV1.CreateManagementZoneOpts{
		ManagementZone: optional.NewInterface(mz),
	}

	managementZone, _, err := dynatraceConfigClientV1.ManagementZonesApi.CreateManagementZone(authConfigV1, &mzBody)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create dynatrace client",
			Detail:   "Bad Request or unable to connect to environment/authenticate API token",
		})
		return diags
	}

	d.SetId(managementZone.Id)

	resourceDynatraceManagementZoneRead(ctx, d, m)

	return diags
}

func resourceDynatraceManagementZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConf := m.(*ProviderConfiguration)
	dynatraceConfigClientV1 := providerConf.DynatraceConfigClientV1
	authConfigV1 := providerConf.AuthConfigV1

	var diags diag.Diagnostics

	managementZoneID := d.Id()

	managementZone, _, err := dynatraceConfigClientV1.ManagementZonesApi.GetSingleManagementZoneConfig(authConfigV1, managementZoneID, &dynatraceConfigV1.GetSingleManagementZoneConfigOpts{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create dynatrace client",
			Detail:   "Bad Request or unable to connect to environment/authenticate API token",
		})
		return diags
	}

	managementZoneRules := flattenManagementZoneRulesData(&managementZone.Rules)
	if err := d.Set("rule", managementZoneRules); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", &managementZone.Name)

	return diags
}

func resourceDynatraceManagementZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConf := m.(*ProviderConfiguration)
	dynatraceConfigClientV1 := providerConf.DynatraceConfigClientV1
	authConfigV1 := providerConf.AuthConfigV1

	var diags diag.Diagnostics

	managementZoneID := d.Id()

	if d.HasChange("name") || d.HasChange("rule") {

		mz := dynatraceConfigV1.ManagementZone{
			Name:  d.Get("name").(string),
			Rules: expandManagementZoneRules(d.Get("rule").([]interface{})),
		}

		mzBody := dynatraceConfigV1.CreateOrUpdateManagementZoneOpts{
			ManagementZone: optional.NewInterface(mz),
		}

		_, _, err := dynatraceConfigClientV1.ManagementZonesApi.CreateOrUpdateManagementZone(authConfigV1, managementZoneID, &mzBody)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create dynatrace client",
				Detail:   "Bad Request or unable to connect to environment/authenticate API token",
			})
			return diags
		}
	}

	return resourceDynatraceManagementZoneRead(ctx, d, m)
}

func resourceDynatraceManagementZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerConf := m.(*ProviderConfiguration)
	dynatraceConfigClientV1 := providerConf.DynatraceConfigClientV1
	authConfigV1 := providerConf.AuthConfigV1

	var diags diag.Diagnostics

	managementZoneID := d.Id()

	_, err := dynatraceConfigClientV1.ManagementZonesApi.DeleteManagementZone(authConfigV1, managementZoneID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create dynatrace client",
			Detail:   "Unable to connect to environment and/or authenticate API token",
		})
		return diags
	}

	d.SetId("")

	return diags
}

func expandManagementZoneRules(rules []interface{}) []dynatraceConfigV1.ManagementZoneRule {
	if len(rules) < 1 {
		return []dynatraceConfigV1.ManagementZoneRule{}
	}

	mrs := make([]dynatraceConfigV1.ManagementZoneRule, len(rules))

	for i, rule := range rules {

		m := rule.(map[string]interface{})

		mrs[i] = dynatraceConfigV1.ManagementZoneRule{
			Type:             m["type"].(string),
			Enabled:          m["enabled"].(bool),
			PropagationTypes: expandPropagationTypes(m["propagation_types"].(*schema.Set).List()),
			Conditions:       expandManagementZoneConditions(m["condition"].([]interface{})),
		}
	}

	return mrs
}

func expandPropagationTypes(propagationTypes []interface{}) []string {
	pts := make([]string, len(propagationTypes))

	for i, v := range propagationTypes {
		pts[i] = v.(string)
	}

	return pts

}

func expandManagementZoneConditions(conditions []interface{}) []dynatraceConfigV1.EntityRuleEngineCondition {
	if len(conditions) < 1 {
		return []dynatraceConfigV1.EntityRuleEngineCondition{}
	}

	mcs := make([]dynatraceConfigV1.EntityRuleEngineCondition, len(conditions))

	for i, condition := range conditions {

		m := condition.(map[string]interface{})
		mcs[i] = dynatraceConfigV1.EntityRuleEngineCondition{
			Key:            expandConditionKey(m["key"].([]interface{})),
			ComparisonInfo: expandConditionComparisonInfo(m["comparison_info"].([]interface{})),
		}
	}

	return mcs
}

func expandConditionKey(conditionKey []interface{}) dynatraceConfigV1.ConditionKey {
	if len(conditionKey) == 0 || conditionKey[0] == nil {
		return dynatraceConfigV1.ConditionKey{}
	}

	m := conditionKey[0].(map[string]interface{})

	mck := dynatraceConfigV1.ConditionKey{}

	if attribute, ok := m["attribute"]; ok {
		mck.Attribute = attribute.(string)
	}

	if mkType, ok := m["attribute"]; ok {
		mck.Attribute = mkType.(string)
	}

	return mck

}

func expandConditionComparisonInfo(comparisonInfo []interface{}) dynatraceConfigV1.ComparisonBasic {
	if len(comparisonInfo) == 0 || comparisonInfo[0] == nil {
		return dynatraceConfigV1.ComparisonBasic{}
	}

	m := comparisonInfo[0].(map[string]interface{})

	mci := dynatraceConfigV1.ComparisonBasic{}

	if operator, ok := m["operator"]; ok {
		mci.Operator = operator.(string)
	}

	if value, ok := m["value"]; ok {
		mci.Value = value.(map[string]interface{})
	}

	if negate, ok := m["negate"]; ok {
		mci.Negate = negate.(bool)
	}

	if ciType, ok := m["type"]; ok {
		mci.Type = ciType.(string)
	}

	return mci

}

func flattenManagementZoneRulesData(managementZoneRules *[]dynatraceConfigV1.ManagementZoneRule) []interface{} {
	if managementZoneRules != nil {
		mrs := make([]interface{}, len(*managementZoneRules), len(*managementZoneRules))

		for i, managementZoneRules := range *managementZoneRules {
			mr := make(map[string]interface{})

			mr["type"] = managementZoneRules.Type
			mr["enabled"] = managementZoneRules.Enabled
			mr["propagation_types"] = managementZoneRules.PropagationTypes
			mr["condition"] = flattenManagementZoneConditionsData(&managementZoneRules.Conditions)
			mrs[i] = mr

		}
		return mrs
	}

	return make([]interface{}, 0)
}

func flattenManagementZoneConditionsData(managementZoneConditions *[]dynatraceConfigV1.EntityRuleEngineCondition) []interface{} {
	if managementZoneConditions != nil {
		mcs := make([]interface{}, len(*managementZoneConditions), len(*managementZoneConditions))

		for i, managementZoneConditions := range *managementZoneConditions {
			mc := make(map[string]interface{})

			mc["key"] = flattenManagementZoneKey(&managementZoneConditions.Key)
			mc["comparison_info"] = flattenManagementZoneComparisonInfo(&managementZoneConditions.ComparisonInfo)
			mcs[i] = mc
		}

		return mcs
	}

	return make([]interface{}, 0)

}

func flattenManagementZoneKey(managementZoneConditionKey *dynatraceConfigV1.ConditionKey) []interface{} {
	if managementZoneConditionKey == nil {
		return []interface{}{managementZoneConditionKey}
	}

	k := make(map[string]interface{})

	k["attribute"] = managementZoneConditionKey.Attribute
	k["type"] = managementZoneConditionKey.Type

	return []interface{}{k}
}

func flattenManagementZoneComparisonInfo(managementZoneComparisonInfo *dynatraceConfigV1.ComparisonBasic) []interface{} {
	if managementZoneComparisonInfo == nil {
		return []interface{}{managementZoneComparisonInfo}
	}

	c := make(map[string]interface{})

	c["operator"] = managementZoneComparisonInfo.Operator
	c["value"] = managementZoneComparisonInfo.Value
	c["negate"] = managementZoneComparisonInfo.Negate
	c["type"] = managementZoneComparisonInfo.Type

	return []interface{}{c}
}
