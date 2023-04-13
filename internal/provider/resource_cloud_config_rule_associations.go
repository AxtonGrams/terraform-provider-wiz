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

func resourceWizCloudConfigRuleAssociations() *schema.Resource {
	return &schema.Resource{
		Description: "Manage associations between cloud configuration rules and security sub-categories. Associations defined outside this resouce will remain untouched through the lifecycle of this resource. Wiz managed cloud configuration rules cannot be associated to Wiz managed security sub-categories. This resource does not support imports; it can, however, overlay existing resources to bring them under management.",
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
			"cloud_config_rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of cloud configuration rule IDs.",
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
		CreateContext: resourceWizCloudConfigRuleAssociationsCreate,
		ReadContext:   resourceWizCloudConfigRuleAssociationsRead,
		UpdateContext: resourceWizCloudConfigRuleAssociationsUpdate,
		DeleteContext: resourceWizCloudConfigRuleAssociationsDelete,
	}
}

func validateCloudConfigRulesExist(ctx context.Context, m interface{}, cloudConfigRuleIDs []string) (diags diag.Diagnostics) {
	tflog.Info(ctx, "validateCloudConfigRulesExist called...")

	// define the graphql query
	query := `query cloudConfigurationRule (
	  $id: ID!
	){
	  cloudConfigurationRule(
	    id: $id
	  ) {
	    id
	    securitySubCategories {
	      id
	    }
	  }
	}`

	for _, b := range cloudConfigRuleIDs {
		vars := &internal.QueryVariables{}
		vars.ID = b

		// process the request
		data := &ReadCloudConfigurationRulePayload{}
		requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_config_rule", "read")

		// handle any errors
		if len(requestDiags) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Cloud config rule not found: %s", b))
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when creating wiz_cloud_config_rule_associations",
				Detail:   fmt.Sprintf("Cloud config rule not found: %s", b),
			})
		}
	}
	return diags
}

// UpdateCloudConfigurationRules struct
type UpdateCloudConfigurationRules struct {
	UpdateCloudConfigurationRules wiz.UpdateCloudConfigurationRulesPayload `json:"updateCloudConfigurationRules"`
}

func resourceWizCloudConfigRuleAssociationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigRuleAssociationsCreate called...")

	// validate each cloud_config_rule and security sub-category exists
	cloudConfigRuleDiags := validateCloudConfigRulesExist(ctx, m, utils.ConvertListToString(d.Get("cloud_config_rule_ids").([]interface{})))
	diags = append(diags, cloudConfigRuleDiags...)

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
	mutation := `mutation UpdateCloudConfigurationRulesInput(
	  $input: UpdateCloudConfigurationRulesInput!
	) {
	  updateCloudConfigurationRules(
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
	mvars := &wiz.UpdateCloudConfigurationRulesInput{}
	mvars.IDs = utils.ConvertListToString(d.Get("cloud_config_rule_ids").([]interface{}))
	mvars.SecuritySubCategoriesToAdd = utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{}))

	// print the input variables
	tflog.Debug(ctx, fmt.Sprintf("UpdateCloudConfigRulesInput: %s", utils.PrettyPrint(mvars)))

	// process the request
	mdata := &UpdateCloudConfigurationRules{}
	mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "cloud_config_rule_association", "create")
	diags = append(diags, mrequestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// error handling
	if mdata.UpdateCloudConfigurationRules.FailCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Error encountered during operation: %s", utils.PrettyPrint(mdata.UpdateCloudConfigurationRules.Errors)))
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error during UpdateCloudConfigurationRules: %d", mdata.UpdateCloudConfigurationRules.FailCount),
			Detail:   fmt.Sprintf("Details: %s", utils.PrettyPrint(mdata.UpdateCloudConfigurationRules.Errors)),
		})
	}

	return resourceWizCloudConfigRuleAssociationsRead(ctx, d, m)
}

func resourceWizCloudConfigRuleAssociationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigRuleAssociationsRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// set each parameter. since the list of ids triggers a new resource, we can pass through
	err := d.Set("details", d.Get("details").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// read current sub-categories for each cloud_config_rule
	// if a sub-category is missing from a cloud_config_rule, taint the resource by removing the cloud_config_rule from the state
	tflog.Debug(ctx, fmt.Sprintf("Cloud Config Rule IDs: %T %s", d.Get("cloud_config_rule_ids"), utils.PrettyPrint(d.Get("cloud_config_rule_ids"))))

	// define the graphql query
	query := `query cloudConfigurationRule(
	  $id: ID!
	){
	  cloudConfigurationRule(
	    id: $id
	  ) {
	    id
	    securitySubCategories {
	      id
	    }
	  }
	}`

	// declare a variable to store the cloud_config_rule ids that have the desired security sub-categories
	var cleanCloudConfigRules = make([]string, 0, 0)

	// iterate over each cloud_config_rule
	tflog.Debug(ctx, fmt.Sprintf("cloud_config_rule_ids for read: %s", d.Get("cloud_config_rule_ids").([]interface{})))
	for _, b := range d.Get("cloud_config_rule_ids").([]interface{}) {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, b))

		// populate the graphql variables
		vars := &internal.QueryVariables{}
		vars.ID = b.(string)

		// process the request
		data := &ReadCloudConfigurationRulePayload{}
		requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_config_rule_association", "read")
		diags = append(diags, requestDiags...)
		if len(diags) > 0 {
			tflog.Error(ctx, "Error from API call, resource not found.")
			if data.CloudConfigurationRule.ID == "" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error reading cloud_config_rule",
					Detail:   fmt.Sprintf("Cloud Config Rule ID: %s", b.(string)),
				})
				return diags
			}
		}

		// store the security sub-category ids for the cloud_config_rule in sscids
		var sscids []string
		for e, f := range data.CloudConfigurationRule.SecuritySubCategories {
			tflog.Debug(ctx, fmt.Sprintf("e: %T %d", e, e))
			tflog.Debug(ctx, fmt.Sprintf("f: %T %s", f, utils.PrettyPrint(f)))
			sscids = append(sscids, f.ID)
		}

		// compare the security sub-categories read with those defined in the hcl
		missing := utils.Missing(sscids, utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{})))
		tflog.Debug(ctx, fmt.Sprintf("Missing security sub-categories for cloud_config_rule %s: %s", b.(string), missing))
		if len(missing) == 0 {
			cleanCloudConfigRules = append(cleanCloudConfigRules, b.(string))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Clean cloud_config_rules from read operation: %T %s", cleanCloudConfigRules, utils.PrettyPrint(cleanCloudConfigRules)))

	err = d.Set("cloud_config_rule_ids", cleanCloudConfigRules)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceWizCloudConfigRuleAssociationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigRuleAssociationsUpdate called...")

	err := d.Set("details", d.Get("details").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return resourceWizCloudConfigRuleAssociationsRead(ctx, d, m)
}

func resourceWizCloudConfigRuleAssociationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCloudConfigRuleAssociationsDelete called...")

	// define the graphql query
	mutation := `mutation UpdateCloudConfigurationRulesInput(
	  $input: UpdateCloudConfigurationRulesInput!
	) {
	  updateCloudConfigurationRules(
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
	mvars := &wiz.UpdateCloudConfigurationRulesInput{}
	mvars.IDs = utils.ConvertListToString(d.Get("cloud_config_rule_ids").([]interface{}))
	mvars.SecuritySubCategoriesToRemove = utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{}))

	// print the input variables
	tflog.Debug(ctx, fmt.Sprintf("DeleteCloudConfigurationRulesInput: %s", utils.PrettyPrint(mvars)))

	// process the request
	mdata := &UpdateControls{}
	mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "cloud_config_rule_association", "delete")
	diags = append(diags, mrequestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// error handling
	if mdata.UpdateControls.FailCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Error encountered during operation: %s", utils.PrettyPrint(mdata.UpdateControls.Errors)))
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error during DeleteCloudConfigurationRuleAssociations: %d", mdata.UpdateControls.FailCount),
			Detail:   fmt.Sprintf("Details: %s", utils.PrettyPrint(mdata.UpdateControls.Errors)),
		})
	}

	return diags
}
