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
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func resourceWizSecuritySubCategoryAssociation() *schema.Resource {
	return &schema.Resource{
		Description: "Manage associations between security sub-categories and policies. This resource can only be used with custom security sub-categories. Wiz managed or custom policies can be referenced. When the association is removed from state, all associations managed by this resource will be removed. Associations managed outside this resouce declaration will remain untouched through the lifecycle of this resource.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Internal identifier for the association.",
				Computed:    true,
			},
			"details": {
				Type:        schema.TypeString,
				Description: "Details of the association. This information is not used to manage resources, but can serve as notes for the associations.",
				Optional:    true,
			},
			"cloud_config_rule_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of cloud config rule IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ExactlyOneOf: []string{
					"cloud_config_rule_ids",
					"control_ids",
					"host_config_rule_ids",
				},
			},
			"control_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of control IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ExactlyOneOf: []string{
					"cloud_config_rule_ids",
					"control_ids",
					"host_config_rule_ids",
				},
			},
			"host_config_rule_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of host config rule IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ExactlyOneOf: []string{
					"cloud_config_rule_ids",
					"control_ids",
					"host_config_rule_ids",
				},
			},
			"security_sub_category_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security sub-category ID.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.IsUUID,
				),
			},
		},
		CreateContext: resourceWizSecuritySubCategoryAssociationCreate,
		ReadContext:   resourceWizSecuritySubCategoryAssociationRead,
		UpdateContext: resourceWizSecuritySubCategoryAssociationUpdate,
		DeleteContext: resourceWizSecuritySubCategoryAssociationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadSecurityFrameworkPayload struct -- updates
type ReadSecuritySubCategoryPayload struct {
	SecuritySubCategory vendor.SecuritySubCategory `json:"securitySubCategory"`
}

// UpdateControls struct
type UpdateControls struct {
	UpdateControls vendor.UpdateControlsPayload `json:"updateControls"`
}

func resourceWizSecuritySubCategoryAssociationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecuritySubCategoryAssociationCreate called...")

	// determine the policy type
	var associationType string
	var checkAssociationType bool

	_, checkAssociationType = d.GetOk("host_config_rule_ids")
	if checkAssociationType {
		associationType = "host_config_rule"
	}
	_, checkAssociationType = d.GetOk("control_ids")
	if checkAssociationType {
		associationType = "control"
	}
	_, checkAssociationType = d.GetOk("cloud_config_rule_ids")
	if checkAssociationType {
		associationType = "cloud_config_rule"
	}
	tflog.Debug(ctx, fmt.Sprintf("Association Type: %s", associationType))

	// read security sub-category associations for the policy type
	switch associationType {
	case "control":
		tflog.Debug(ctx, fmt.Sprintf("Control IDs: %T %s", d.Get("control_ids"), utils.PrettyPrint(d.Get("control_ids"))))

		// define the graphql query
		query := `query Control (
		  $id: ID!
		){
		  control(
		    id: $id
		  ) {
		    id
		    name
		    description
		    enabled
		    securitySubCategories {
		      id
		      title
		    }
		  }
		}`

		for a, b := range d.Get("control_ids").([]interface{}) {
			tflog.Debug(ctx, fmt.Sprintf("a: %T %d", a, a))
			tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, b))

			// populate the graphql variables
			vars := &internal.QueryVariables{}
			vars.ID = b.(string)

			// process the request
			data := &ReadControlPayload{}
			requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "security_sub_category_association", "read")
			diags = append(diags, requestDiags...)
			if len(diags) > 0 {
				tflog.Error(ctx, "Error from API call, resource not found.")
				if data.Control.ID == "" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Control not found",
						Detail:   fmt.Sprintf("Control ID: %s", b.(string)),
					})

					return diags
				}
			}

			// Retrieve the security sub-category ids
			var sscids []string
			for e, f := range data.Control.SecuritySubCategories {
				tflog.Debug(ctx, fmt.Sprintf("e: %T %d", e, e))
				tflog.Debug(ctx, fmt.Sprintf("f: %T %s", f, utils.PrettyPrint(f)))
				sscids = append(sscids, f.ID)
			}

			tflog.Debug(ctx, fmt.Sprintf("Existing security sub-category associations for control %s %s", b, utils.PrettyPrint(sscids)))

			// determine which security sub-categories are defined in the tf resource but not in wiz
			missing := utils.Missing(
				sscids,
				[]string{d.Get("security_sub_category_id").(string)},
			)
			tflog.Debug(ctx, fmt.Sprintf("Security sub-categories missing from control %s: %s", b.(string), missing))

			// define the graphql query
			mutation := `mutation UpdateControls(
		          $input: UpdateControlsInput!
		        ) {
		          updateControls(
		            input: $input
		          ) {
		            successCount
			    failCount
			    errors {
			      reason
			      control {
			        id
			      }
			    }
		          }
		        }`

			// populate the graphql variables
			mvars := &vendor.UpdateControlsInput{}
			mvars.IDS = []string{b.(string)}
			mvars.SecuritySubCategoriesToAdd = missing

			// print the input variables
			tflog.Debug(ctx, fmt.Sprintf("Updates: %s", utils.PrettyPrint(mvars)))

			// process the request
			mdata := &vendor.UpdateControlsPayload{}
			mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "security_sub_category_association", "update")
			diags = append(diags, mrequestDiags...)
			if len(diags) > 0 {
				return diags
			}
		}
	}

	// set the id
	d.SetId(fmt.Sprintf("%s-%s", associationType, d.Get("security_sub_category_id").(string)))

	//return resourceWizControlRead(ctx, d, m)
	return diags
}

func resourceWizSecuritySubCategoryAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecuritySubCategoryAssociationRead called...")

	// read current security sub-categories for the control
	// remove all sub-categories not defined in the tf resource
	// populate the resource data

	return diags
}

func resourceWizSecuritySubCategoryAssociationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecuritySubCategoryAssociationUpdate called...")

	// read current security sub-categories for the control
	// determine which security sub-categories are defined in the tf resource but not in wiz
	// compute a superset of security sub-categories
	// issue the update

	return resourceWizSecuritySubCategoryAssociationRead(ctx, d, m)
}

func resourceWizSecuritySubCategoryAssociationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecuritySubCategoryAssociationDelete called...")

	// read current security sub-categories for the control
	// strip all security sub-categories defined in the tf resource from the list
	// issue the update

	return diags
}
