package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizControl() *schema.Resource {
	return &schema.Resource{
		Description: "A Control consists of a pre-defined Security Graph query and a severity level—if a Control's query returns any results, an Issue is generated for every result. Each Control is assigned to a category in one or more Policy Frameworks.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal identifier for the Control",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Control.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Control.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable the Control. This has a known defect where if set to false, it will be created as true because the API to create Controls does not accept this parameter.",
			},
			"project_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "*",
				Description: "Project scope of the control. Use '*' for all projects.",
			},
			"security_sub_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of security sub-categories IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resolution_recommendation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Guidance on how the user should address an issue that was created by this control.",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The query that the control runs.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					var ov, nv interface{}
					_ = json.Unmarshal([]byte(oldValue), &ov)
					_ = json.Unmarshal([]byte(newValue), &nv)
					return reflect.DeepEqual(ov, nv)
				},
			},
			"scope_query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The query that represents the control's scope.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					var ov, nv interface{}
					_ = json.Unmarshal([]byte(oldValue), &ov)
					_ = json.Unmarshal([]byte(newValue), &nv)
					return reflect.DeepEqual(ov, nv)
				},
			},
			"severity": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Severity that will be set for this control.\n    - Allowed values: %s",
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
			},
		},
		CreateContext: resourceWizControlCreate,
		ReadContext:   resourceWizControlRead,
		UpdateContext: resourceWizControlUpdate,
		DeleteContext: resourceWizControlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateControl struct
type CreateControl struct {
	CreateControl wiz.CreateControlPayload `json:"createControl"`
}

func resourceWizControlCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlCreate called...")

	// define the graphql query
	query := `mutation createControl(
	    $input: CreateControlInput!
	) {
	    createControl(
	        input: $input
	    ) {
	        control {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.CreateControlInput{}
	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	val, hasVal := d.GetOk("resolution_recommendation")
	if hasVal {
		vars.ResolutionRecommendation = val.(string)
	}
	vars.Severity = d.Get("severity").(string)
	vars.ProjectID = d.Get("project_id").(string)
	vars.Query = json.RawMessage(d.Get("query").(string))
	vars.ScopeQuery = json.RawMessage(d.Get("scope_query").(string))
	vars.SecuritySubCategories = utils.ConvertListToString(d.Get("security_sub_categories").([]interface{}))

	// process the request
	data := &CreateControl{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateControl.Control.ID)

	return resourceWizControlRead(ctx, d, m)
}

func flattenControlSecuritySubCategories(ctx context.Context, securitySubCategories []*wiz.SecuritySubCategory) []interface{} {
	tflog.Info(ctx, "flattenControlSecuritySubCategories called...")
	tflog.Debug(ctx, fmt.Sprintf("flattenControlSecuritySubCategories input: %T %s", securitySubCategories, utils.PrettyPrint(securitySubCategories)))

	var output = make([]interface{}, 0, 0)

	for a, b := range securitySubCategories {
		tflog.Debug(ctx, fmt.Sprintf("a: %T %d", a, a))
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenControlSecuritySubCategories output: %+v", output))

	return output
}

// ReadControlPayload struct -- updates
type ReadControlPayload struct {
	Control wiz.Control `json:"control"`
}

func resourceWizControlRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

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
	        query
	        scopeQuery
	        severity
	        securitySubCategories {
	            id
	            title
	        }
	        enabled
	        resolutionRecommendation
	        scopeProject {
	            id
	            name
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: oops! an internal error has occurred. for reference purposes, this is your request id
	data := &ReadControlPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.Control.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.Control.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("description", data.Control.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("severity", data.Control.Severity)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("resolution_recommendation", data.Control.ResolutionRecommendation)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("enabled", data.Control.Enabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	// special handling is required for project since input * is not reflected in read response
	if d.Get("project_id").(string) == "*" {
		err = d.Set("project_id", "*")
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	} else {
		err = d.Set("project_id", data.Control.ScopeProject.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	err = d.Set("security_sub_categories", flattenControlSecuritySubCategories(ctx, data.Control.SecuritySubCategories))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	b, _ := json.Marshal(data.Control.Query)
	err = d.Set("query", string(b))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateControl struct
type UpdateControl struct {
	UpdateControl wiz.UpdateControlPayload `json:"updateControl"`
}

func resourceWizControlUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateControl(
	    $input: UpdateControlInput!
	) {
	    updateControl(
		input: $input
	    ) {
		control {
		    id
		}
	    }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateControlInput{}
	vars.ID = d.Id()

	// these can optionally be included in the patch
	if d.HasChange("enabled") {
		vars.Patch.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	}
	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		vars.Patch.Description = d.Get("description").(string)
	}
	if d.HasChange("severity") {
		vars.Patch.Severity = d.Get("severity").(string)
	}
	if d.HasChange("resolution_recommendation") {
		vars.Patch.ResolutionRecommendation = d.Get("resolution_recommendation").(string)
	}
	if d.HasChange("query") {
		vars.Patch.Query = json.RawMessage(d.Get("query").(string))
	}
	if d.HasChange("scope_query") {
		vars.Patch.ScopeQuery = json.RawMessage(d.Get("scope_query").(string))
	}
	if d.HasChange("security_sub_categories") {
		vars.Patch.SecuritySubCategories = utils.ConvertListToString(d.Get("security_sub_categories").([]interface{}))
	}

	// process the request
	data := &UpdateControl{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizControlRead(ctx, d, m)
}

// DeleteControl struct
type DeleteControl struct {
	DeleteControl wiz.DeleteControlPayload `json:"deleteControl"`
}

func resourceWizControlDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizControlDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteControl (
	    $input: DeleteControlInput!
	) {
	    deleteControl(
		input: $input
	    ) {
		_stub
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteControlInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateControl{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
