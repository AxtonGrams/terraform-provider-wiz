package provider

import (
	"context"
	"encoding/json"
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

func resourceWizAutomationRule() *schema.Resource {
	return &schema.Resource{
		Description: "Automation Rules define associations between actions and findings.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Wiz internal identifier.",
				Computed:    true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
				Default:     "",
			},
			"trigger_source": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Trigger source.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.AutomationRuleTriggerSource,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						vendor.AutomationRuleTriggerSource,
						false,
					),
				),
			},
			"trigger_type": {
				Type:     schema.TypeList,
				Required: true,
				Description: fmt.Sprintf(
					"Trigger type.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.AutomationRuleTriggerType,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							vendor.AutomationRuleTriggerType,
							false,
						),
					),
				},
			},
			"filters": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
				Description: "Value should be wrapped in jsonencode() to avoid diff detection. This is required even though the API states it is not required.  Validate is performed by the UI.",
			},
			"action_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AutomationActions to execute once an automation rule event is triggered and passes the filters",
			},
			"override_action_params": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "{}",
				Description: "Optional parameters that can override the default automationaction parameters that have been defined when the automationaction was created.  Value should be wrapped in jsonencode() to avoid diff detection.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled?",
				Default:     true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
		CreateContext: resourceWizAutomationRuleCreate,
		ReadContext:   resourceWizAutomationRuleRead,
		UpdateContext: resourceWizAutomationRuleUpdate,
		DeleteContext: resourceWizAutomationRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAutomationRule struct
type CreateAutomationRule struct {
	CreateAutomationRule vendor.CreateAutomationRulePayload `json:"createAutomationRule"`
}

func resourceWizAutomationRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleCreate called...")

	// define the graphql query
	query := `mutation CreateAutomationRule (
	    $input: CreateAutomationRuleInput!
	) {
	    createAutomationRule (
	        input: $input
	    ) {
	        automationRule {
	            id
	            createdAt
	        }
	    }
	}`

	// populate the graphql variables
	vars := &vendor.CreateAutomationRuleInput{}
	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	vars.ActionID = d.Get("action_id").(string)
	vars.TriggerSource = d.Get("trigger_source").(string)
	vars.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	vars.ProjectID = d.Get("project_id").(string)
	vars.OverrideActionParams = json.RawMessage(d.Get("override_action_params").(string))
	vars.Filters = json.RawMessage(d.Get("filters").(string))
	vars.TriggerType = utils.ConvertListToString(d.Get("trigger_type").([]interface{}))

	// process the request
	data := &CreateAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id and computed values
	d.SetId(data.CreateAutomationRule.AutomationRule.ID)
	d.Set("created_at", data.CreateAutomationRule.AutomationRule.CreatedAt)

	return resourceWizAutomationRuleRead(ctx, d, m)
}

// ReadAutomationRulePayload struct -- updates
type ReadAutomationRulePayload struct {
	AutomationRule vendor.AutomationRule `json:"automationRule"`
}

func resourceWizAutomationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query automationRule (
	    $id: ID!
	){
	    automationRule (
	        id: $id
	    ) {
	        id
	        createdAt
	        name
	        description
	        action {id}
	        triggerSource
	        triggerType
	        enabled
	        filters
	        overrideActionParams
	        project {id}
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: record not found for id
	data := &ReadAutomationRulePayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_action", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.AutomationRule.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.AutomationRule.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("action_id", data.AutomationRule.Action.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("trigger_source", data.AutomationRule.TriggerSource)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("description", data.AutomationRule.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("enabled", data.AutomationRule.Enabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("filters", string(data.AutomationRule.Filters))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("override_action_params", string(data.AutomationRule.OverrideActionParams))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", data.AutomationRule.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("trigger_type", data.AutomationRule.TriggerType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("created_at", data.AutomationRule.CreatedAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateAutomationRule struct
type UpdateAutomationRule struct {
	UpdateAutomationRule vendor.UpdateAutomationRulePayload `json:"updateAutomationRule"`
}

func resourceWizAutomationRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation updateAutomationRule($input: UpdateAutomationRuleInput!) {
	    updateAutomationRule(
	        input: $input
	    ) {
	        automationRule {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &vendor.UpdateAutomationRuleInput{}
	vars.ID = d.Id()
	vars.Patch.Name = d.Get("name").(string)
	vars.Patch.ActionID = d.Get("action_id").(string)
	vars.Patch.TriggerSource = d.Get("trigger_source").(string)

	triggerTypes := make([]string, 0, 0)
	for _, j := range d.Get("trigger_type").([]interface{}) {
		triggerTypes = append(triggerTypes, j.(string))
	}
	vars.Patch.TriggerType = triggerTypes

	vars.Patch.Description = d.Get("description").(string)
	vars.Patch.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	vars.Patch.Filters = json.RawMessage(d.Get("filters").(string))
	vars.Patch.OverrideActionParams = json.RawMessage(d.Get("override_action_params").(string))

	// process the request
	data := &UpdateAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizAutomationRuleRead(ctx, d, m)
}

// DeleteAutomationRule struct
type DeleteAutomationRule struct {
	DeleteAutomationRule vendor.DeleteAutomationRulePayload `json:"deleteAutomationRule"`
}

func resourceWizAutomationRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteAutomationRule (
            $input: DeleteAutomationRuleInput!
        ) {
            deleteAutomationRule (
                input: $input
            ) {
                _stub
            }
        }`

	// populate the graphql variables
	vars := &vendor.DeleteAutomationRuleInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
