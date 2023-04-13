package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizCloudConfigurationRule() *schema.Resource {
	return &schema.Resource{
		Description: "A Cloud Configuration Rule is a configuration check that applies to a specific cloud resource type—if a resource does not pass a Rule, a Configuration Finding is generated and associated with the resource on the Security Graph.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Wiz internal identifier.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this rule, as appeared in the UI in the portal.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Detailed description for this rule. There is a defect in the API that makes this required; the description field cannot be nullified after one is defined, so we make it required.",
			},
			"target_native_types": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The identifier types of the resources targeted by this rule, as seen on the cloud provider service. e.g. 'ec2'",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"opa_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OPA rego policy that defines this rule.",
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
				Description: fmt.Sprintf(
					"Severity that will be set for findings of this rule.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.Severity,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.Severity,
						false,
					),
				),
				Default: "MEDIUM",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/disable this rule.",
				Default:     true,
			},
			"remediation_instructions": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Steps to mitigate the issue that match this rule. If possible, include sample commands to execute in your cloud provider's console. Markdown formatting is supported.",
			},
			"scope_account_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set the rule scope of cloud account IDs. Select only subscriptions matching to the rule cloud provider. To change scope to 'all relevant resources' set to empty array. This must be the Wiz internal identifier for the account(uuid format).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"function_as_control": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Make this rule function as a Control that creates Issues for new findings. By default only findings are created. If enabled=false, an error will be returned if this is set to true.",
				Default:     false,
			},
			"security_sub_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Associate this rule with security sub-categories to easily monitor your compliance. New Configuration Findings created by this rule will be tagged with the selected sub-categories. There is a defect in the API that makes this required; the security_sub_categories field cannot be nullified after one is defined, so we make it required.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"iac_matchers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "OPA rego policies that this rule runs (Cloud / IaC rules).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							Description: fmt.Sprintf(
								"The type of resource that will be evaluated by the Rego Code.\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.CloudConfigurationRuleMatcherType,
								),
							),
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.CloudConfigurationRuleMatcherType,
									false,
								),
							),
						},
						"rego_code": {
							Type:        schema.TypeString,
							Description: "Write code in the Rego query language. This code will be evaluated against the JSON representation of each resource of the selected Native Type to determine if it passes or fails the rule.",
							Required:    true,
						},
					},
				},
			},
		},
		CreateContext: resourceWizCloudConfigurationRuleCreate,
		ReadContext:   resourceWizCloudConfigurationRuleRead,
		UpdateContext: resourceWizCloudConfigurationRuleUpdate,
		DeleteContext: resourceWizCloudConfigurationRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func getIACMatchers(ctx context.Context, d *schema.ResourceData) []*wiz.CreateCloudConfigurationRuleMatcherInput {
	tflog.Info(ctx, "getIACMatchers called...")

	iacMatchers := d.Get("iac_matchers").(*schema.Set).List()
	var myIacMatchers []*wiz.CreateCloudConfigurationRuleMatcherInput
	for _, a := range iacMatchers {
		tflog.Debug(ctx, fmt.Sprintf("a: %t %s", a, utils.PrettyPrint(a)))
		localIacMatchers := &wiz.CreateCloudConfigurationRuleMatcherInput{}
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "type":
				localIacMatchers.Type = c.(string)
			case "rego_code":
				localIacMatchers.RegoCode = c.(string)
			}
		}
		myIacMatchers = append(myIacMatchers, localIacMatchers)
	}
	tflog.Debug(ctx, fmt.Sprintf("myIacMatchers: %s", utils.PrettyPrint(myIacMatchers)))
	return myIacMatchers
}

// CreateCloudConfigurationRule struct
type CreateCloudConfigurationRule struct {
	CreateCloudConfigurationRule wiz.CreateCloudConfigurationRulePayload `json:"createCloudConfigurationRule"`
}

func resourceWizCloudConfigurationRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigurationRuleCreate called...")

	// define the graphql query
	query := `mutation CreateCloudConfigurationRule(
	    $input: CreateCloudConfigurationRuleInput!
	) {
	    createCloudConfigurationRule(
	        input: $input
	    ) {
	        rule {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.CreateCloudConfigurationRuleInput{}
	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	vars.TargetNativeTypes = utils.ConvertListToString(d.Get("target_native_types").(*schema.Set).List())
	vars.OPAPolicy = d.Get("opa_policy").(string)
	vars.Severity = d.Get("severity").(string)
	vars.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	vars.RemediationInstructions = d.Get("remediation_instructions").(string)
	vars.IACMatchers = getIACMatchers(ctx, d)
	vars.ScopeAccountIDs = utils.ConvertListToString(d.Get("scope_account_ids").(*schema.Set).List())
	vars.FunctionAsControl = utils.ConvertBoolToPointer(d.Get("function_as_control").(bool))

	// process the request
	data := &CreateCloudConfigurationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_configuration_rule", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id and computed values
	d.SetId(data.CreateCloudConfigurationRule.Rule.ID)

	return resourceWizCloudConfigurationRuleRead(ctx, d, m)
}

func flattenIACMatchers(ctx context.Context, iacMatchers []*wiz.CloudConfigurationRuleMatcher) []interface{} {
	tflog.Info(ctx, "flattenIACMatchers called...")

	var output = make([]interface{}, 0, 0)
	for _, b := range iacMatchers {
		tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		var mapping = make(map[string]interface{})
		mapping["type"] = b.Type
		mapping["rego_code"] = b.RegoCode
		output = append(output, mapping)
	}
	tflog.Debug(ctx, fmt.Sprintf("flattenIACMatchers output: %s", utils.PrettyPrint(output)))
	return output
}

func flattenSecuritySubCategoriesID(ctx context.Context, securitySubCategories []*wiz.SecuritySubCategory) []interface{} {
	tflog.Info(ctx, "flattenSecuritySubCategoriesID called...")

	var output = make([]interface{}, 0, 0)
	for _, b := range securitySubCategories {
		tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}
	tflog.Debug(ctx, fmt.Sprintf("flattenSecuritySubCategoriesID output: %s", utils.PrettyPrint(output)))
	return output
}

func flattenScopeAccountIDs(ctx context.Context, scopeAccounts []*wiz.CloudAccount) []interface{} {
	tflog.Info(ctx, "flattenScopeAccountIDs called...")

	var output = make([]interface{}, 0, 0)
	for _, b := range scopeAccounts {
		tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}
	tflog.Debug(ctx, fmt.Sprintf("flattenScopeAccountIDs output: %s", utils.PrettyPrint(output)))
	return output
}

// ReadCloudConfigurationRulePayload struct -- updates
type ReadCloudConfigurationRulePayload struct {
	CloudConfigurationRule wiz.CloudConfigurationRule `json:"cloudConfigurationRule"`
}

func resourceWizCloudConfigurationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigurationRuleRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query cloudConfigurationRule (
	    $id: ID!
	){
	    cloudConfigurationRule(
	        id: $id
	    ) {
	        id
	        name
	        description
	        targetNativeTypes
	        opaPolicy
	        severity
	        enabled
	        remediationInstructions
	        scopeAccounts {
	            id
	        }
	        functionAsControl
	        securitySubCategories {
	            id
	        }
	        iacMatchers {
	            type
	            regoCode
	        }
	        control {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: record not found for id
	data := &ReadCloudConfigurationRulePayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_config_rule", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.CloudConfigurationRule.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.CloudConfigurationRule.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("description", data.CloudConfigurationRule.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("target_native_types", data.CloudConfigurationRule.TargetNativeTypes)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("severity", data.CloudConfigurationRule.Severity)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("enabled", data.CloudConfigurationRule.Enabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("remediation_instructions", data.CloudConfigurationRule.RemediationInstructions)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("function_as_control", data.CloudConfigurationRule.FunctionAsControl)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	scopeAccountIDs := flattenScopeAccountIDs(ctx, data.CloudConfigurationRule.ScopeAccounts)
	if err := d.Set("scope_account_ids", scopeAccountIDs); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	securitySubCategories := flattenSecuritySubCategoriesID(ctx, data.CloudConfigurationRule.SecuritySubCategories)
	if err := d.Set("security_sub_categories", securitySubCategories); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	iacMatchers := flattenIACMatchers(ctx, data.CloudConfigurationRule.IACMatchers)
	if err := d.Set("iac_matchers", iacMatchers); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateCloudConfigurationRule struct
type UpdateCloudConfigurationRule struct {
	UpdateCloudConfigurationRule wiz.UpdateCloudConfigurationRulePayload `json:"updateCloudConfigurationRule"`
}

func resourceWizCloudConfigurationRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigurationRuleUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateCloudConfigurationRule(
	    $input: UpdateCloudConfigurationRuleInput!
	) {
	    updateCloudConfigurationRule(
	        input: $input
	    ) {
	        rule {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateCloudConfigurationRuleInput{}
	vars.ID = d.Id()
	// check if changes were made to required fields
	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		vars.Patch.Description = d.Get("description").(string)
	}
	if d.HasChange("remediation_instructions") {
		vars.Patch.RemediationInstructions = d.Get("remediation_instructions").(string)
	}
	if d.HasChange("target_native_types") {
		targetNativeTypes := make([]string, 0)
		for _, j := range d.Get("target_native_types").(*schema.Set).List() {
			targetNativeTypes = append(targetNativeTypes, j.(string))
		}
		vars.Patch.TargetNativeTypes = targetNativeTypes
	}
	// include all optional fields in the patch in the event they were nullified
	vars.Patch.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	vars.Patch.OPAPolicy = d.Get("opa_policy").(string)
	vars.Patch.Severity = d.Get("severity").(string)
	vars.Patch.FunctionAsControl = utils.ConvertBoolToPointer(d.Get("function_as_control").(bool))
	// flatten scopeAccountIds
	scopeAccountIds := make([]string, 0)
	for _, j := range d.Get("scope_account_ids").(*schema.Set).List() {
		scopeAccountIds = append(scopeAccountIds, j.(string))
	}
	vars.Patch.ScopeAccountIds = scopeAccountIds
	// flatten iacMatchers
	iacMatchers := d.Get("iac_matchers")
	iacMatcherUpdates := make([]*wiz.UpdateCloudConfigurationRuleMatcherInput, 0)
	for _, b := range iacMatchers.(*schema.Set).List() {
		var myMap = &wiz.UpdateCloudConfigurationRuleMatcherInput{}
		tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
		for c, d := range b.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			tflog.Trace(ctx, fmt.Sprintf("d: %T %s", d, d))
			switch c {
			case "type":
				myMap.Type = d.(string)
			case "rego_code":
				myMap.RegoCode = d.(string)
			}
		}
		iacMatcherUpdates = append(iacMatcherUpdates, myMap)
	}
	vars.Patch.IACMatchers = iacMatcherUpdates

	// process the request
	data := &UpdateCloudConfigurationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_configuration_rule", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizCloudConfigurationRuleRead(ctx, d, m)
}

// DeleteCloudConfigurationRule struct
type DeleteCloudConfigurationRule struct {
	DeleteCloudConfigurationRule wiz.DeleteCloudConfigurationRulePayload `json:"deleteCloudConfigurationRule"`
}

func resourceWizCloudConfigurationRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigurationRuleDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteCloudConfigurationRule (
            $input: DeleteCloudConfigurationRuleInput!
        ) {
            deleteCloudConfigurationRule (
                input: $input
            ) {
                _stub
            }
        }`

	// populate the graphql variables
	vars := &wiz.DeleteCloudConfigurationRuleInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateCloudConfigurationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_configuration_rule", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
