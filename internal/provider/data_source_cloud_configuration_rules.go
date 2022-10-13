package provider

import (
	//"bytes"
	"context"
	//"crypto/sha1"
	//"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	//"wiz.io/hashicorp/terraform-provider-wiz/internal"
	//"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func dataSourceWizCloudConfigurationRules() *schema.Resource {
	return &schema.Resource{
		Description: "Query cloud configuration rules.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal identifier for the data.",
			},
			"first": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "How many results to return",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Free text search on CSPM name or resource ID.",
			},
			"scope_account_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Find CSPM rules applied on cloud account IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_provider": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("Find CSPM rules related to cloud provider.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.CloudProvider,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							vendor.CloudProvider,
							false,
						),
					),
				},
			},
			"service_type": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("Find CSPM rules related to the service.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.CloudConfigurationRuleServiceType,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							vendor.CloudConfigurationRuleServiceType,
							false,
						),
					),
				},
			},
			"subject_entity_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "FreeFind rules by their entity type subject.",
			},
			"severity": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("CSPM Rule severity.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.Severity,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							vendor.Severity,
							false,
						),
					),
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "CSPM Rule enabled status.",
			},
			"has_auto_remediation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Rule has auto remediation.",
			},
			"has_remediation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Rule has remediation.",
			},
			"framework_category": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search rules by any of securityFramework | securitySubCategory | securityCategory.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_native_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search rules by target native type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_by": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search rules by user.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_opa_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Search by opaPolicy presence.",
			},
			"project": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search by project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"matcher_type": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("Search rules by target native type.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.CloudConfigurationRuleMatcherTypeFilter,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							vendor.CloudConfigurationRuleMatcherTypeFilter,
							false,
						),
					),
				},
			},
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "GetSearch by IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"function_as_control": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Search by function as control.",
			},
			"risk_equals_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"risk_equals_all": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_configuration_rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The returned cloud configuration rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Wiz UUID.",
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"short_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Rule enabled status.",
						},
						"severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule severity will outcome to finding severity. This filed initial value is set as the severity of the CSPM rule.",
						},
						"external_references": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"target_native_types": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The identifier types of the objects targeted by this rule, as seen on the cloud provider service. e.g. 'ec2'.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"supports_nrt": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Rule enabled status.",
						},
						"subject_entity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule subject entity type, as represented on Wiz Security Graph.",
						},
						"cloud_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud provider this rule is relevant to.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service this rule is relevant to.",
						},
						"scope_accounts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Scope of subscription IDs for automatically asses with this rule on. If set to empty array rule will run on all environment",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"security_sub_category_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"builtin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indication whether the rule is built-in or custom.",
						},
						"opa_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OPA rego policy that this rule runs. Undefined for built-in code based configuration rules.",
						},
						"function_as_control": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Make this rule also function as a control which means findings by this control will also trigger Issues.",
						},
						"control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "In case this rule also functions as a control, this property will contain its details.",
						},
						"graph_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vertex id of this rule on the graph.",
						},
						"has_auto_remediation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remediation_instructions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iac_matcher_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "OPA rego policies that this rule runs (Cloud / IaC rules).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		ReadContext: dataSourceWizCloudConfigurationRuleRead,
	}
}

func dataSourceWizCloudConfigurationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizCloudConfigurationRuleRead called...")

	return diags
}
