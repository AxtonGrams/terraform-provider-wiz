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

func resourceWizControlAssociations() *schema.Resource {
	return &schema.Resource{
		Description: "Manage associations between controls and security sub-categories. Associations defined outside this resouce will remain untouched through the lifecycle of this resource. Wiz managed controls cannot be associated to Wiz managed security sub-categories. This resource does not support imports; it can, however, overlay existing resources to bring them under management.",
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
				ForceNew:    true,
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
	}
}

func validateControlsExist(ctx context.Context, m interface{}, controlIDs []string) (diags diag.Diagnostics) {
	tflog.Info(ctx, "validateControlsExist called...")

	query := `query Control (
	  $id: ID!
	) {
	  control(
	    id: $id
	  ) {
	    id
	  }
	}`

	for _, b := range controlIDs {
		vars := &internal.QueryVariables{}
		vars.ID = b

		// process the request
		data := &ReadControlPayload{}
		requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "read")

		// handle any errors
		if len(requestDiags) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Control not found: %s", b))
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when creating wiz_control_associations",
				Detail:   fmt.Sprintf("Control not found: %s", b),
			})
		}
	}
	return diags
}

func validateSecuritySubCategoriesExist(ctx context.Context, m interface{}, securitySubCategoryIDs []string) (diags diag.Diagnostics) {
	tflog.Info(ctx, "validateSecuritySubCategoriesExist called...")

	query := `query securitySubCategory  (
	  $id: ID!
	){
	  securitySubCategory(
	    id: $id
	  ) {
	    id
	  }
	}`

	for _, b := range securitySubCategoryIDs {
		vars := &internal.QueryVariables{}
		vars.ID = b

		// process the request
		data := &ReadSecuritySubCategoryPayload{}
		requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "security_sub_category", "read")

		// handle any errors
		if len(requestDiags) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Security sub-category not found: %s", b))
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error when creating control association",
				Detail:   fmt.Sprintf("Security sub-category not found: %s", b),
			})
		}
	}
	return diags
}

// UpdateControls struct
type UpdateControls struct {
	UpdateControls wiz.UpdateControlsPayload `json:"updateControls"`
}

// ReadSecuritySubCategoryPayload struct
type ReadSecuritySubCategoryPayload struct {
	SecuritySubCategory wiz.SecuritySubCategory `json:"securitySubCategory"`
}

func resourceWizControlAssociationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsCreate called...")

	// validate each control and security sub-category exists
	controlDiags := validateControlsExist(ctx, m, utils.ConvertListToString(d.Get("control_ids").([]interface{})))
	diags = append(diags, controlDiags...)

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
	mvars := &wiz.UpdateControlsInput{}
	mvars.IDs = utils.ConvertListToString(d.Get("control_ids").([]interface{}))
	mvars.SecuritySubCategoriesToAdd = utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{}))

	// print the input variables
	tflog.Debug(ctx, fmt.Sprintf("UpdateControlsInput: %s", utils.PrettyPrint(mvars)))

	// process the request
	mdata := &UpdateControls{}
	mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "control_association", "create")
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

	return resourceWizControlAssociationsRead(ctx, d, m)
}

func resourceWizControlAssociationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// set each parameter. since the list of ids triggers a new resource, we can pass through
	err := d.Set("details", d.Get("details").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// read current sub-categories for each control
	// if a sub-category is missing from a control, taint the resource by removing the control from the state
	tflog.Debug(ctx, fmt.Sprintf("Control IDs for resource %s: %T %s", d.Get("id"), d.Get("control_ids"), utils.PrettyPrint(d.Get("control_ids"))))

	// define the graphql query
	query := `query Control (
          $id: ID!
        ){
          control(
            id: $id
          ) {
            id
            securitySubCategories {
              id
            }
          }
        }`

	// declare a variable to store the control ids that have the desired security sub-categories
	var cleanControls = make([]string, 0, 0)

	// iterate over each control
	tflog.Debug(ctx, fmt.Sprintf("control_ids for read: %s", d.Get("control_ids").([]interface{})))
	resourceControlIDs := d.Get("control_ids").([]interface{})
	for _, b := range resourceControlIDs {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, b))

		// populate the graphql variables
		vars := &internal.QueryVariables{}
		vars.ID = b.(string)

		// process the request
		data := &ReadControlPayload{}
		requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control_association", "read")
		diags = append(diags, requestDiags...)
		if len(diags) > 0 {
			tflog.Error(ctx, "Error from API call, resource not found.")
			if data.Control.ID == "" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error reading control",
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

		// compare the security sub-categories read with those defined in the hcl
		missing := utils.Missing(sscids, utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{})))
		tflog.Debug(ctx, fmt.Sprintf("Missing security sub-categories for control %s: %s", b.(string), missing))
		if len(missing) == 0 {
			cleanControls = append(cleanControls, b.(string))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Clean controls from read operation: %T %s", cleanControls, utils.PrettyPrint(cleanControls)))

	// terraform presents a diff even when the values and their order are the same, so long as Set() is called
	tflog.Debug(ctx, "Ensuring that we only set control IDs when they have different lengths, or when the orders are different.")
	var newControlIds []string
	resourceControlIDsLen := len(resourceControlIDs)
	cleanControlsLen := len(cleanControls)
	if resourceControlIDsLen != cleanControlsLen {
		newControlIds = cleanControls
	} else {
		for i, resourceControlID := range resourceControlIDs {
			cleanControlID := cleanControls[i]
			if resourceControlID != cleanControlID {
				newControlIds = cleanControls
				break
			}
		}
	}

	if len(newControlIds) > 0 {
		err = d.Set("control_ids", newControlIds)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

func resourceWizControlAssociationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsUpdate called...")

	err := d.Set("details", d.Get("details").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return resourceWizControlAssociationsRead(ctx, d, m)
}

func resourceWizControlAssociationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlAssociationsDelete called...")

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
	mvars := &wiz.UpdateControlsInput{}
	mvars.IDs = utils.ConvertListToString(d.Get("control_ids").([]interface{}))
	mvars.SecuritySubCategoriesToRemove = utils.ConvertListToString(d.Get("security_sub_category_ids").([]interface{}))

	// print the input variables
	tflog.Debug(ctx, fmt.Sprintf("UpdateControlsInput: %s", utils.PrettyPrint(mvars)))

	// process the request
	mdata := &UpdateControls{}
	mrequestDiags := client.ProcessRequest(ctx, m, mvars, mdata, mutation, "control_association", "delete")
	diags = append(diags, mrequestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// error handling
	if mdata.UpdateControls.FailCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Error encountered during operation: %s", utils.PrettyPrint(mdata.UpdateControls.Errors)))
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error during DeleteControlAssociations: %d", mdata.UpdateControls.FailCount),
			Detail:   fmt.Sprintf("Details: %s", utils.PrettyPrint(mdata.UpdateControls.Errors)),
		})
	}

	return diags
}
