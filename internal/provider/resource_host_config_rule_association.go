package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizHostConfigRuleAssociations() *schema.Resource {
	return &schema.Resource{
		Description: "Manage associations between host configuration rules and security sub-categories. Associations defined outside this resouce will remain untouched through the lifecycle of this resource. Wiz managed host configuration rules cannot be associated to Wiz managed security sub-categories. This resource does not support imports; it can, however, overlay existing resources to bring them under management.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Internal identifier for the association.",
				Computed:    true,
			},
			"details": {
				Type:        schema.TypeString,
				Description: "Details of the association. This information is not used to manage resources but can serve as notes or documentation for the associations.",
				Optional:    true,
				Default:     "undefined",
			},
			"host_config_rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of host configuration rule IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_sub_category_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of security sub-category IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		CreateContext: resourceWizHostConfigRuleAssociationsCreate,
		ReadContext:   resourceWizHostConfigRuleAssociationsRead,
		UpdateContext: resourceWizHostConfigRuleAssociationsUpdate,
		DeleteContext: resourceWizHostConfigRuleAssociationsDelete,
	}
}

func validateHostConfigRulesExist(ctx context.Context, m interface{}, hostConfigRuleIDs []string) (diags diag.Diagnostics) {
	tflog.Info(ctx, "validateHostConfigRulesExist called...")

	query := `query hostConfigurationRule (
	  $id: ID!
	) {
	  hostConfigurationRule(
	    id: $id
	  ) {
	    id
	    securitySubCategories {
	      id
	    }
	  }
	}`

	for _, b := range hostConfigRuleIDs {
		vars := &internal.QueryVariables{}
		vars.ID = b

		// process the request
		data := &ReadHostConfigurationRulePayload{}
		requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "host_config_rule", "read")

		// handle any errors
		if len(requestDiags) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Host config rule not found: %s", b))
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when creating wiz_host_config_rule_associations",
				Detail:   fmt.Sprintf("Host config rule not found: %s", b),
			})
		}
	}
	return diags
}

// UpdateHostConfigurationRules struct
type UpdateHostConfigurationRules struct {
	UpdateHostConfigurationRules wiz.UpdateHostConfigurationRulesPayload `json:"updateHostConfigurationRules"`
}

func resourceWizHostConfigRuleAssociationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizHostConfigRuleAssociationsCreate called...")

	// validate each host configuration rule and security sub-category exists
	hostConfigRuleDiags := validateHostConfigRulesExist(ctx, m, utils.ConvertListToString(d.Get("host_config_rule_ids").([]interface{})))
	diags = append(diags, hostConfigRuleDiags...)

	securitySubCategoryDiags := validateSecuritySubCategoriesExist(ctx, m, utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{})))
	diags = append(diags, securitySubCategoryDiags...)

	if len(diags) > 0 {
		return diags
	}

	// generate an id for this resource
	uuid := uuid.New().String()

	// Set the id
	d.SetId(uuid)

	// define the graphql query
	mutation := `mutation UpdateHostConfigurationRules(
	  $input: UpdateHostConfigurationRulesInput!
	) {
	  updateHostConfigurationRules(
	    input: $input
	  ) {
	    successCount
	    failCount
	    errors {
	      reason
	      rule {
	        id
	      }
	    }
	  }
	}`

	// populate the graphql variables
	mvars := &wiz.UpdateHostConfigurationRulesInput{}
	mvars.IDs = utils.ConvertListToString(d.Get("host_config_rule_ids").([]interface{}))
	mvars.SecuritySubCategoriesToAdd = utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{}))

	// print the input variables
	tflog.Debug(ctx, fmt.Sprintf("UpdateHostConfigurationRulesInput: %s", utils.PrettyPrint(mvars)))

	// process the request
	mdata := &UpdateHostConfigurationRules{}
	mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "host_config_rule_association", "create")
	diags = append(diags, mrequestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// error handling
	if mdata.UpdateHostConfigurationRules.FailCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Error encountered during operation: %s", utils.PrettyPrint(mdata.UpdateHostConfigurationRules.Errors)))
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error during UpdateHostConfigurationRules: %d", mdata.UpdateHostConfigurationRules.FailCount),
			Detail:   fmt.Sprintf("Details: %s", utils.PrettyPrint(mdata.UpdateHostConfigurationRules.Errors)),
		})
	}

	return resourceWizHostConfigRuleAssociationsRead(ctx, d, m)
}

func resourceWizHostConfigRuleAssociationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizHostConfigRuleAssociationsRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// set each parameter. since the list of ids triggers a new resource, we can pass through
	err := d.Set("details", d.Get("details").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// read current sub-categories for each host config rule
	// if a sub-category is missing from a host config rule, taint the resource by removing the host config rule from the state
	tflog.Debug(ctx, fmt.Sprintf("Host config rule IDs: %T %s", d.Get("host_config_rule_ids"), utils.PrettyPrint(d.Get("host_config_rule_ids"))))

	// define the graphql query
	query := `query hostConfigurationRule (
	  $id: ID!
	){
	  hostConfigurationRule(
	    id: $id
	  ) {
	    id
	    securitySubCategories {
	      id
	    }
	  }
	}`

	// declare a variable to store the host config rule ids that have the desired security sub-categories
	var cleanHostConfigRules = make([]string, 0, 0)

	// iterate over each host config rule
	tflog.Debug(ctx, fmt.Sprintf("host_config_rule_ids for read: %s", d.Get("host_config_rule_ids").([]interface{})))
	for _, b := range d.Get("host_config_rule_ids").([]interface{}) {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, b))

		// populate the graphql variables
		vars := &internal.QueryVariables{}
		vars.ID = b.(string)

		// process the request
		data := &ReadHostConfigurationRulePayload{}
		requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "host_config_rule_association", "read")
		diags = append(diags, requestDiags...)
		if len(diags) > 0 {
			tflog.Error(ctx, "Error from API call, resource not found.")
			if data.HostConfigurationRule.ID == "" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error reading host configuration rule",
					Detail:   fmt.Sprintf("Host configuration rule ID: %s", b.(string)),
				})
				return diags
			}
		}

		// store the security sub-category ids for the host config rule in sscids
		var sscids []string
		for e, f := range data.HostConfigurationRule.SecuritySubCategories {
			tflog.Debug(ctx, fmt.Sprintf("e: %T %d", e, e))
			tflog.Debug(ctx, fmt.Sprintf("f: %T %s", f, utils.PrettyPrint(f)))
			sscids = append(sscids, f.ID)
		}

		// compare the security sub-categories read with those defined in the hcl
		missing := utils.Missing(sscids, utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{})))
		tflog.Debug(ctx, fmt.Sprintf("Missing security sub-categories for control %s: %s", b.(string), missing))
		if len(missing) == 0 {
			cleanHostConfigRules = append(cleanHostConfigRules, b.(string))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Clean host config rules from read operation: %T %s", cleanHostConfigRules, utils.PrettyPrint(cleanHostConfigRules)))

	err = d.Set("host_config_rule_ids", cleanHostConfigRules)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceWizHostConfigRuleAssociationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizHostConfigRuleAssociationsUpdate called...")

	err := d.Set("details", d.Get("details").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return resourceWizHostConfigRuleAssociationsRead(ctx, d, m)
}

func resourceWizHostConfigRuleAssociationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizHostConfigRuleAssociationsDelete called...")

	// define the graphql query
	mutation := `mutation UpdateHostConfigurationRules(
	  $input: UpdateHostConfigurationRulesInput!
	) {
	  updateHostConfigurationRules(
	    input: $input
	  ) {
	    successCount
	    failCount
	    errors {
	      reason
	      rule {
	        id
	      }
	    }
	  }
	}`

	// populate the graphql variables
	mvars := &wiz.UpdateHostConfigurationRulesInput{}
	mvars.IDs = utils.ConvertListToString(d.Get("host_config_rule_ids").([]interface{}))
	mvars.SecuritySubCategoriesToRemove = utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{}))

	// print the input variables
	tflog.Debug(ctx, fmt.Sprintf("UpdateHostConfigurationRulesInput IDs: %s", d.Get("host_config_rule_ids").([]interface{})))
	tflog.Debug(ctx, fmt.Sprintf("UpdateHostConfigurationRulesInput: %s", utils.PrettyPrint(mvars)))

	// process the request
	mdata := &UpdateHostConfigurationRules{}
	mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "host_config_rule_association", "delete")
	diags = append(diags, mrequestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// error handling
	if mdata.UpdateHostConfigurationRules.FailCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Error encountered during operation: %s", utils.PrettyPrint(mdata.UpdateHostConfigurationRules.Errors)))
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error during DeleteHostConfigRuleAssociations: %d", mdata.UpdateHostConfigurationRules.FailCount),
			Detail:   fmt.Sprintf("Details: %s", utils.PrettyPrint(mdata.UpdateHostConfigurationRules.Errors)),
		})
	}

	return diags
}
