package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizSecurityFramework() *schema.Resource {
	return &schema.Resource{
		Description: "Configure Security Frameworks and associated resources (Categories and Subcategories). Support for extended fields has not been implemented due to issues with the API. This includes: category.external_id, category.sub_category.resolution_recommendation, and category.sub_category.external_id.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Internal identifier for the Security Framework",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the security framework.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the security framework.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the security framework.",
				Optional:    true,
				Default:     true,
			},
			"category": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Security framework category.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Internal identifier for the security category. Specify an existing identifier to use an existing category. If not provided, a new category will be created.",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Name fo the security category.",
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the security category.",
						},
						"sub_category": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Security subcategory.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "Internal identifier for the security subcategory. Specify an existing identifier to use an existing subcategory. If not provided, a new subcategory will be created.",
										Computed:    true,
									},
									"title": {
										Type:        schema.TypeString,
										Description: "Title of the security subcategory.",
										Required:    true,
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the security subcategory.",
									},
								},
							},
						},
					},
				},
			},
		},
		CreateContext: resourceWizSecurityFrameworkCreate,
		ReadContext:   resourceWizSecurityFrameworkRead,
		UpdateContext: resourceWizSecurityFrameworkUpdate,
		DeleteContext: resourceWizSecurityFrameworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func getSecurityCategories(ctx context.Context, data *schema.ResourceData) []wiz.SecurityCategoryInput {
	tflog.Debug(ctx, "getSecurityCategories called...")

	securityCategories := data.Get("category").(*schema.Set)
	var mySecurityCategories []wiz.SecurityCategoryInput

	for a, b := range securityCategories.List() {
		tflog.Debug(ctx, fmt.Sprintf("a: %d", a))
		tflog.Debug(ctx, fmt.Sprintf("b: %t %s", b, utils.PrettyPrint(b)))

		localSecurityCategory := wiz.SecurityCategoryInput{}

		for c, d := range b.(map[string]interface{}) {
			tflog.Debug(ctx, fmt.Sprintf("c: %s", utils.PrettyPrint(c)))
			tflog.Debug(ctx, fmt.Sprintf("d: %t %s", d, utils.PrettyPrint(d)))

			//localSecuritySubCategories := []wiz.SecuritySubCategoryInput{}

			switch c {
			case "name":
				localSecurityCategory.Name = d.(string)
			case "description":
				localSecurityCategory.Description = d.(string)
			case "id":
				localSecurityCategory.ID = d.(string)
			case "sub_category":
				localSecurityCategory.SubCategories = getSecuritySubCategories(ctx, d.(*schema.Set))
			}
		}
		mySecurityCategories = append(mySecurityCategories, localSecurityCategory)
	}
	tflog.Debug(ctx, fmt.Sprintf("getSecurityCategories: %s", utils.PrettyPrint(mySecurityCategories)))

	return mySecurityCategories
}

func getSecuritySubCategories(ctx context.Context, set *schema.Set) []wiz.SecuritySubCategoryInput {
	tflog.Debug(ctx, "getSecuritySubCategories called...")

	var mySecuritySubCategories []wiz.SecuritySubCategoryInput

	for a, b := range set.List() {
		tflog.Debug(ctx, fmt.Sprintf("a: %d", a))
		tflog.Debug(ctx, fmt.Sprintf("b: %t %s", b, utils.PrettyPrint(b)))

		localSubCategory := wiz.SecuritySubCategoryInput{}

		for c, d := range b.(map[string]interface{}) {
			tflog.Debug(ctx, fmt.Sprintf("c: %s", utils.PrettyPrint(c)))
			tflog.Debug(ctx, fmt.Sprintf("d: %t %s", d, utils.PrettyPrint(d)))

			switch c {
			case "title":
				localSubCategory.Title = d.(string)
			case "description":
				localSubCategory.Description = d.(string)
			case "id":
				localSubCategory.ID = d.(string)
			}
		}
		mySecuritySubCategories = append(mySecuritySubCategories, localSubCategory)
	}
	tflog.Debug(ctx, fmt.Sprintf("getSecuritySubCategories: %s", utils.PrettyPrint(mySecuritySubCategories)))

	return mySecuritySubCategories
}

// CreateSecurityFramework struct
type CreateSecurityFramework struct {
	CreateSecurityFramework wiz.CreateSecurityFrameworkPayload `json:"createSecurityFramework"`
}

func resourceWizSecurityFrameworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecurityFrameworkCreate called...")

	// define the graphql query
	query := `mutation CreateSecurityFramework(
    $input: CreateSecurityFrameworkInput!
	) {
	    createSecurityFramework(
	        input: $input
	    ) {
	        framework {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.CreateSecurityFrameworkInput{}
	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	vars.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	vars.Categories = getSecurityCategories(ctx, d)

	// process the request
	data := &CreateSecurityFramework{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "security_framework", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateSecurityFramework.Framework.ID)

	return resourceWizSecurityFrameworkRead(ctx, d, m)
}

func flattenSecurityCategories(ctx context.Context, securityFrameworks wiz.SecurityFramework) []interface{} {
	tflog.Info(ctx, "flattenSecurityCategories called...")
	tflog.Debug(ctx, fmt.Sprintf("flattenSecurityCategories input: %T %s", securityFrameworks, utils.PrettyPrint(securityFrameworks)))

	var output = make([]interface{}, 0, 0)

	for a, b := range securityFrameworks.Categories {
		tflog.Debug(ctx, fmt.Sprintf("a: %d", a))
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		var securityCategory = make(map[string]interface{})

		securityCategory["description"] = b.Description
		securityCategory["id"] = b.ID
		securityCategory["name"] = b.Name
		securityCategory["sub_category"] = flattenSecuritySubCategories(ctx, b.SubCategories)
		output = append(output, securityCategory)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenSecurityCategories output: %+v", output))

	return output
}

func flattenSecuritySubCategories(ctx context.Context, securitySubCategories []wiz.SecuritySubCategory) []interface{} {
	tflog.Info(ctx, "flattenSecuritySubCategories called...")
	tflog.Debug(ctx, fmt.Sprintf("flattenSecuritySubCategories input: %T %s", securitySubCategories, utils.PrettyPrint(securitySubCategories)))

	var output = make([]interface{}, 0, 0)

	for a, b := range securitySubCategories {
		tflog.Debug(ctx, fmt.Sprintf("a: %d", a))
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		var securitySubCategory = make(map[string]interface{})
		securitySubCategory["description"] = b.Description
		securitySubCategory["title"] = b.Title
		securitySubCategory["id"] = b.ID
		output = append(output, securitySubCategory)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenSecuritySubCategories output: %+v", output))

	return output
}

// ReadSecurityFrameworkPayload struct -- updates
type ReadSecurityFrameworkPayload struct {
	SecurityFramework wiz.SecurityFramework `json:"securityFramework"`
}

func resourceWizSecurityFrameworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecurityFrameworkRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query securityFramework  (
	    $id: ID!
	){
	    securityFramework(
	        id: $id
	    ) {
	        id
	        name
	        description
	        enabled
	        categories {
	            id
	            name
	            description
	            subCategories {
	                id
	                title
	                description
	            }
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: oops! an internal error has occurred. for reference purposes, this is your request id
	data := &ReadSecurityFrameworkPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "security_framework", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.SecurityFramework.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.SecurityFramework.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("description", data.SecurityFramework.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("enabled", data.SecurityFramework.Enabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	securityCategories := flattenSecurityCategories(ctx, data.SecurityFramework)
	if err := d.Set("category", securityCategories); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateSecurityFramework struct
type UpdateSecurityFramework struct {
	UpdateSecurityFramework wiz.UpdateSecurityFrameworkPayload `json:"updateSecurityFramework"`
}

func resourceWizSecurityFrameworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecurityFrameworkUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateSecurityFramework(
	    $input: UpdateSecurityFrameworkInput!
	) {
	    updateSecurityFramework(
	        input: $input
	    ) {
	        framework {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateSecurityFrameworkInput{}
	vars.ID = d.Id()

	// description must be passed with every update
	vars.Patch.Description = d.Get("description").(string)

	// these can optionally be included in the patch
	if d.HasChange("enabled") {
		vars.Patch.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	}
	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}

	// if security catetories are altered, we must send the all security categories
	if d.HasChange("category") {
		vars.Patch.Categories = getSecurityCategories(ctx, d)
	}

	// process the request
	data := &UpdateSecurityFramework{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "security_framework", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizSecurityFrameworkRead(ctx, d, m)
}

// DeleteSecurityFramework struct
type DeleteSecurityFramework struct {
	DeleteSecurityFramework wiz.DeleteSecurityFrameworkPayload `json:"deleteSecurityFramework"`
}

func resourceWizSecurityFrameworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSecurityFrameworkDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteSecurityFramework (
	    $input: DeleteSecurityFrameworkInput!
	) {
	    deleteSecurityFramework(
	        input: $input
	    ) {
	        _stub
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteSecurityFrameworkInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateSecurityFramework{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "security_framework", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
