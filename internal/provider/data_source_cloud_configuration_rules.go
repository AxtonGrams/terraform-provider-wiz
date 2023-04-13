package provider

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
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
						wiz.CloudProvider,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.CloudProvider,
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
						wiz.CloudConfigurationRuleServiceType,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.CloudConfigurationRuleServiceType,
							false,
						),
					),
				},
			},
			"subject_entity_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Find rules by their entity type subject.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"severity": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("CSPM Rule severity.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.Severity,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.Severity,
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
						wiz.CloudConfigurationRuleMatcherTypeFilter,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.CloudConfigurationRuleMatcherTypeFilter,
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
						"description": {
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

// ReadCloudConfigurationRules struct
type ReadCloudConfigurationRules struct {
	CloudConfigurationRules wiz.CloudConfigurationRuleConnection `json:"cloudConfigurationRules"`
}

func dataSourceWizCloudConfigurationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizCloudConfigurationRuleRead called...")

	// generate the id for this resource
	// id must be deterministic, so the id is based on a hash of the search parameters
	var identifier bytes.Buffer

	a, b := d.GetOk("first")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("search")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("scope_account_ids")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("service_type")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("subject_entity_type")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("severity")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("enabled")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("has_auto_remediation")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("has_remediation")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("framework_category")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("target_native_type")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("created_by")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("is_opa_policy")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("project")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("matcher_type")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("ids")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("function_as_control")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("risk_equals_any")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("risk_equals_all")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}

	h := sha1.New()
	h.Write([]byte(identifier.String()))
	hashID := hex.EncodeToString(h.Sum(nil))

	// Set the id
	d.SetId(hashID)

	// define the graphql query
	query := `query cloudConfigurationRules(
	  $filterBy: CloudConfigurationRuleFilters
	  $first: Int
	  $after: String
	  $orderBy: CloudConfigurationRuleOrder
	) {
	  cloudConfigurationRules(
	    filterBy: $filterBy
	    first: $first
	    after: $after
	    orderBy: $orderBy
	  ) {
	    nodes {
	      id
	      name
	      shortId
	      description
	      enabled
	      severity
	      externalReferences{
	        id
	        name
	      }
	      targetNativeTypes
	      supportsNRT
	      subjectEntityType
	      cloudProvider
	      serviceType
	      scopeAccounts {
	        id
	      }
	      securitySubCategories {
	        id
	      }
	      builtin
	      opaPolicy
	      functionAsControl
	      control {
	        id
	      }
	      graphId
	      hasAutoRemediation
	      remediationInstructions
	      iacMatchers {
	        id
	      }
	    }
	    pageInfo {
	      hasNextPage
	      endCursor
	    }
	    totalCount
	  }
	}`

	// set the resource parameters
	err := d.Set("first", d.Get("first").(int))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("search", d.Get("search").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("scope_account_ids", d.Get("scope_account_ids").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("cloud_provider", d.Get("cloud_provider").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("service_type", d.Get("service_type").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("subject_entity_type", d.Get("subject_entity_type").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("severity", d.Get("severity").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	a, b = d.GetOk("enabled")
	if b {
		err = d.Set("enabled", a.(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	a, b = d.GetOk("has_auto_remediation")
	if b {
		err = d.Set("has_auto_remediation", a.(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	a, b = d.GetOk("has_remediation")
	if b {
		err = d.Set("has_remediation", a.(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	err = d.Set("framework_category", d.Get("framework_category").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("target_native_type", d.Get("target_native_type").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("created_by", d.Get("created_by").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	a, b = d.GetOk("is_opa_policy")
	if b {
		err = d.Set("is_opa_policy", a.(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	err = d.Set("project", d.Get("project").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("matcher_type", d.Get("matcher_type").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("ids", d.Get("ids").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	a, b = d.GetOk("function_as_control")
	if b {
		err = d.Set("function_as_control", a.(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	err = d.Set("risk_equals_any", d.Get("risk_equals_any").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("risk_equals_all", d.Get("risk_equals_all").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.First = d.Get("first").(int)
	filterBy := &wiz.CloudConfigurationRuleFilters{}
	a, b = d.GetOk("search")
	if b {
		filterBy.Search = a.(string)
	}
	a, b = d.GetOk("scope_account_ids")
	if b {
		filterBy.ScopeAccountIDs = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("cloud_provider")
	if b {
		filterBy.CloudProvider = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("service_type")
	if b {
		filterBy.ServiceType = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("subject_entity_type")
	if b {
		filterBy.SubjectEntityType = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("severity")
	if b {
		filterBy.Severity = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("enabled")
	if b {
		filterBy.Enabled = utils.ConvertBoolToPointer(a.(bool))
	}
	a, b = d.GetOk("has_auto_remediation")
	if b {
		filterBy.HasAutoRemediation = utils.ConvertBoolToPointer(a.(bool))
	}
	a, b = d.GetOk("has_remediation")
	if b {
		filterBy.HasRemediation = utils.ConvertBoolToPointer(a.(bool))
	}
	a, b = d.GetOk("framework_category")
	if b {
		filterBy.FrameworkCategory = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("target_native_type")
	if b {
		filterBy.TargetNativeType = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("created_by")
	if b {
		filterBy.CreatedBy = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("is_opa_policy")
	if b {
		filterBy.IsOPAPolicy = utils.ConvertBoolToPointer(a.(bool))
	}
	a, b = d.GetOk("project")
	if b {
		filterBy.Project = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("matcher_type")
	if b {
		filterBy.MatcherType = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("ids")
	if b {
		filterBy.ID = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("function_as_control")
	if b {
		filterBy.FunctionAsControl = utils.ConvertBoolToPointer(a.(bool))
	}
	a, b = d.GetOk("risk_equals_any")
	if b {
		filterBy.RiskEqualsAny = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("risk_equals_all")
	if b {
		filterBy.RiskEqualsAll = utils.ConvertListToString(a.([]interface{}))
	}

	vars.FilterBy = filterBy

	// process the request
	data := &ReadCloudConfigurationRules{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_config_rules", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	cloudConfigurationRules := flattenCloudConfigurationRules(ctx, &data.CloudConfigurationRules.Nodes)
	if err := d.Set("cloud_configuration_rules", cloudConfigurationRules); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func flattenCloudConfigurationRules(ctx context.Context, nodes *[]*wiz.CloudConfigurationRule) []interface{} {
	tflog.Info(ctx, "flattenCloudConfigurationRules called...")
	tflog.Debug(ctx, fmt.Sprintf("CloudConfigurationRules: %s", utils.PrettyPrint(nodes)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *nodes {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		ruleMap := make(map[string]interface{})
		ruleMap["id"] = b.ID
		ruleMap["name"] = b.Name
		ruleMap["short_id"] = b.ShortID
		ruleMap["description"] = b.Description
		ruleMap["enabled"] = *b.Enabled
		ruleMap["severity"] = b.Severity
		ruleMap["external_references"] = flattenExternalReferences(ctx, &b.ExternalReferences)
		ruleMap["target_native_types"] = b.TargetNativeTypes
		ruleMap["supports_nrt"] = *b.SupportsNRT
		ruleMap["subject_entity_type"] = b.SubjectEntityType
		ruleMap["cloud_provider"] = b.CloudProvider
		ruleMap["service_type"] = b.ServiceType
		ruleMap["scope_accounts"] = flattenScopeAccounts(ctx, &b.ScopeAccounts)
		ruleMap["security_sub_category_ids"] = flattenSecuritySubCategoryIDs(ctx, &b.SecuritySubCategories)
		ruleMap["builtin"] = *b.Builtin
		ruleMap["opa_policy"] = b.OPAPolicy
		ruleMap["function_as_control"] = *b.FunctionAsControl
		if b.Control != nil {
			ruleMap["control_id"] = b.Control.ID
		}
		ruleMap["graph_id"] = b.GraphID
		ruleMap["has_auto_remediation"] = *b.HasAutoRemediation
		ruleMap["remediation_instructions"] = b.RemediationInstructions
		ruleMap["iac_matcher_ids"] = flattenIACMatcherIDs(ctx, &b.IACMatchers)

		output = append(output, ruleMap)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenCloudConfigurationRules output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenExternalReferences(ctx context.Context, refs *[]*wiz.CloudConfigurationRuleExternalReference) []interface{} {
	tflog.Info(ctx, "flattenExternalReferences called...")
	tflog.Debug(ctx, fmt.Sprintf("External References: %s", utils.PrettyPrint(refs)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *refs {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		refMap := make(map[string]interface{})
		refMap["id"] = b.ID
		refMap["name"] = b.Name
		output = append(output, refMap)
	}

	// sort the return slice to avoid unwanted diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].(map[string]interface{})["id"].(string) < output[j].(map[string]interface{})["id"].(string)
	})

	tflog.Debug(ctx, fmt.Sprintf("flattenExternalReferences output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenScopeAccounts(ctx context.Context, accounts *[]*wiz.CloudAccount) []interface{} {
	tflog.Info(ctx, "flattenScopeAccounts called...")
	tflog.Debug(ctx, fmt.Sprintf("ScopeAccounts: %s", utils.PrettyPrint(accounts)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *accounts {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	// sort the return slice to avoid unwanted diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].(string) < output[j].(string)
	})

	tflog.Debug(ctx, fmt.Sprintf("flattenScopeAccounts output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenSecuritySubCategoryIDs(ctx context.Context, subCats *[]*wiz.SecuritySubCategory) []interface{} {
	tflog.Info(ctx, "flattenSecuritySubCategoryIDs called...")
	tflog.Debug(ctx, fmt.Sprintf("SecuritySubCategories: %s", utils.PrettyPrint(subCats)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *subCats {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	// sort the return slice to avoid unwanted diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].(string) < output[j].(string)
	})

	tflog.Debug(ctx, fmt.Sprintf("flattenSecuritySubCategoryIDs output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenIACMatcherIDs(ctx context.Context, matchers *[]*wiz.CloudConfigurationRuleMatcher) []interface{} {
	tflog.Info(ctx, "flattenIACMatchers called...")
	tflog.Debug(ctx, fmt.Sprintf("flattenIACMatchers: %s", utils.PrettyPrint(matchers)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *matchers {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	// sort the return slice to avoid unwanted diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].(string) < output[j].(string)
	})

	tflog.Debug(ctx, fmt.Sprintf("flattenIACMatchers output: %s", utils.PrettyPrint(output)))

	return output
}
