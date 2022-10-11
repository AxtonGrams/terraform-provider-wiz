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
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func resourceWizControlAssociations() *schema.Resource {
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
				Description: "Details of the association. This information is not used to manage resources but can serve as notes or documentation for the associations.",
				Optional:    true,
			},
			"control_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of control IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_sub_category_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of security sub-category IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		CreateContext: resourceWizControlAssociationsCreate,
		ReadContext:   resourceWizControlAssociationsRead,
		UpdateContext: resourceWizControlAssociationsUpdate,
		DeleteContext: resourceWizControlAssociationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// UpdateControls struct
type UpdateControls struct {
	UpdateControls vendor.UpdateControlsPayload `json:"updateControls"`
}

// ReadSecuritySubCategoryPayload struct
type ReadSecuritySubCategoryPayload struct {
	SecuritySubCategory vendor.SecuritySubCategory `json:"securitySubCategory"`
}

func resourceWizControlAssociationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsCreate called...")

	// validate each control exists
	queryControl := `query Control (
	  $id: ID!
	) {
	  control(
	    id: $id
	  ) {
	    id
	  }
	}`

	for _, b := range utils.ConvertListToString(d.Get("control_ids").([]interface{})) {
		qcvars := &internal.QueryVariables{}
		qcvars.ID = b
		// process the request
		data := &ReadControlPayload{}
		requestDiags := client.ProcessRequest(ctx, m, qcvars, data, queryControl, "control", "read")
		diags = append(diags, requestDiags...)
		// handle any errors
		if len(diags) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Control not found: %s", b))
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when creating control association",
				Detail:   fmt.Sprintf("Control not found: %s", b),
			})
		}
	}

	// validate each security sub-category exists
	querySecuritySubCategories := `query securitySubCategory  (
	  $id: ID!
	){
	  securitySubCategory(
	    id: $id
	  ) {
	    id
	  }
	}`

	for _, b := range utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{})) {
		qsvars := &internal.QueryVariables{}
		qsvars.ID = b
		// process the request
		data := &ReadSecuritySubCategoryPayload{}
		requestDiags := client.ProcessRequest(ctx, m, qsvars, data, querySecuritySubCategories, "security_sub_category", "read")
		diags = append(diags, requestDiags...)
		// handle any errors
		if len(diags) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Security sub-category not found: %s", b))
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when creating control association",
				Detail:   fmt.Sprintf("Security sub-category not found: %s", b),
			})
		}
	}

	// generate an id for this resource
	uuid := uuid.New().String()

	// Set the id
	d.SetId(uuid)

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
	mvars.IDS = utils.ConvertListToString(d.Get("control_ids").([]interface{}))
	mvars.SecuritySubCategoriesToAdd = utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{}))

	// print the input variables
	tflog.Debug(ctx, fmt.Sprintf("UpdateControlsInput: %s", utils.PrettyPrint(mvars)))

	// process the request
	mdata := &UpdateControls{}
	mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "security_sub_category_association", "update")
	diags = append(diags, mrequestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// error handling
	if mdata.UpdateControls.FailCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Error encountered during operation: %s", utils.PrettyPrint(mdata.UpdateControls.Errors)))
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error during UpdateControls: %d", mdata.UpdateControls.FailCount),
			Detail:   fmt.Sprintf("Details: %s", utils.PrettyPrint(mdata.UpdateControls.Errors)),
		})
	}

	//return resourceWizControlRead(ctx, d, m)
	return diags
}

func resourceWizControlAssociationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// set the common parameters
	err := d.Set("details", d.Get("details").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("security_sub_category_id", d.Get("security_sub_category_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// read current sub-categories for each control
	// if a sub-category is missing from a control, taint the resource
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

	// iterate over each control id
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

		// store the security sub-category ids for the control in sscids
		var sscids []string
		for e, f := range data.Control.SecuritySubCategories {
			tflog.Debug(ctx, fmt.Sprintf("e: %T %d", e, e))
			tflog.Debug(ctx, fmt.Sprintf("f: %T %s", f, utils.PrettyPrint(f)))
			sscids = append(sscids, f.ID)
		}
	}

	return diags
}

func resourceWizControlAssociationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsUpdate called...")

	// read current security sub-categories for the control
	// determine which security sub-categories are defined in the tf resource but not in wiz
	// compute a superset of security sub-categories
	// issue the update

	return resourceWizControlAssociationsRead(ctx, d, m)
}

func resourceWizControlAssociationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsDelete called...")

	// read current security sub-categories for the control
	// strip all security sub-categories defined in the tf resource from the list
	// issue the update

	return diags
}
