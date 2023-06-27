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
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizAutomationRuleJiraTransitionTicket() *schema.Resource {
	return &schema.Resource{
		Description: "Automation Rules define associations between actions and findings.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Wiz internal identifier.",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date/time at which the automation rule was created.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the automation rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the automation rule",
			},
			"trigger_source": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Trigger source.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.AutomationRuleTriggerSource,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.AutomationRuleTriggerSource,
						false,
					),
				),
			},
			"trigger_type": {
				Type:     schema.TypeList,
				Required: true,
				Description: fmt.Sprintf(
					"Trigger type. Must be set to `CREATED` for wiz_automation_rule_jira_transition_ticket.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.AutomationRuleTriggerType,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.AutomationRuleTriggerType,
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
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled?",
				Default:     true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Wiz internal ID for a project.",
			},
			"action_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Wiz internal ID for the action.",
			},
			"integration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Wiz identifier for the Integration to leverage for this action. Must be resource type integration_jira.",
			},
			"jira_project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue project",
			},
			"jira_transition_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue transition ID or Name",
			},
			"jira_advanced_fields": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
			"jira_comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue Jira comment",
			},
			"jira_comment_on_transition": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not to send comment during follow-up call, if this is disabled comment will be sent as update field",
			},
			"jira_attach_evidence_csv": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Upload issues report as attachment Only relevant in CONTROL-triggered Actions.",
			},
		},
		CreateContext: resourceWizAutomationRuleJiraTransitionTicketCreate,
		ReadContext:   resourceWizAutomationRuleJiraTransitionTicketRead,
		UpdateContext: resourceWizAutomationRuleJiraTransitionTicketUpdate,
		DeleteContext: resourceWizAutomationRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizAutomationRuleJiraTransitionTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleJiraTransitionTicketCreate called...")

	// define the graphql query
	query := `mutation CreateAutomationRule (
	  $input: CreateAutomationRuleInput!
	) {
	  createAutomationRule(
	    input: $input
	  ) {
	    automationRule {
	      id
	    }
	  }
	}`

	// populate the graphql variables
	vars := &wiz.CreateAutomationRuleInput{}
	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	vars.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	vars.Filters = json.RawMessage(d.Get("filters").(string))
	vars.ProjectID = d.Get("project_id").(string)
	vars.TriggerType = utils.ConvertListToString(d.Get("trigger_type").([]interface{}))
	vars.TriggerSource = d.Get("trigger_source").(string)

	// populate the actions parameter
	jiraTransitionTicketParams := &wiz.JiraActionTransitionTicketTemplateParamsInput{
		Project:             d.Get("jira_project").(string),
		TransitionID:        d.Get("jira_transition_id").(string),
		AdvancedFields:      json.RawMessage(d.Get("jira_advanced_fields").(string)),
		Comment:             d.Get("jira_comment").(string),
		CommentOnTransition: utils.ConvertBoolToPointer(d.Get("jira_comment_on_transition").(bool)),
		AttachEvidenceCSV:   utils.ConvertBoolToPointer(d.Get("jira_attach_evidence_csv").(bool)),
	}
	actionTemplateParams := wiz.ActionTemplateParamsInput{
		JiraTransitionTicket: jiraTransitionTicketParams,
	}
	actions := []wiz.AutomationRuleActionInput{}
	action := wiz.AutomationRuleActionInput{
		IntegrationID:        d.Get("integration_id").(string),
		ActionTemplateParams: actionTemplateParams,
		ActionTemplateType:   "JIRA_TRANSITION_TICKET",
	}
	actions = append(actions, action)
	vars.Actions = actions

	// process the request
	data := &CreateAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule_jira_transition_ticket", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id and computed values
	d.SetId(data.CreateAutomationRule.AutomationRule.ID)

	return resourceWizAutomationRuleJiraTransitionTicketRead(ctx, d, m)
}

func resourceWizAutomationRuleJiraTransitionTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleJiraTransitionTicketRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query automationRule (
	  $id: ID!
	){
	  automationRule(
	    id: $id
	  ){
	    id
	    name
	    description
	    createdAt
	    triggerSource
	    triggerType
	    filters
	    enabled
	    project {
	      id
	    }
	    actions {
	      id
	      actionTemplateType
	      integration {
	        id
	      }
	      actionTemplateParams {
	        ... on JiraActionTransitionTicketTemplateParams {
			project
			transitionId
			advancedFields
			comment
			commentOnTransition
			attachEvidenceCSV 
	        }
	      }
	    }
	  }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	automationRuleActions := make([]*wiz.AutomationRuleAction, 0)
	automationRuleAction := &wiz.AutomationRuleAction{
		ActionTemplateParams: &wiz.JiraActionTransitionTicketTemplateParams{},
	}
	automationRuleActions = append(automationRuleActions, automationRuleAction)
	data := &ReadAutomationRulePayload{
		AutomationRule: wiz.AutomationRule{
			Actions: automationRuleActions,
		},
	}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule_jira_transition_ticket", "read")
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
	err = d.Set("description", data.AutomationRule.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("enabled", data.AutomationRule.Enabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("trigger_type", data.AutomationRule.TriggerType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("trigger_source", data.AutomationRule.TriggerSource)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("filters", string(data.AutomationRule.Filters))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", data.AutomationRule.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("created_at", data.AutomationRule.CreatedAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", data.AutomationRule.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("action_id", data.AutomationRule.Actions[0].ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("integration_id", data.AutomationRule.Actions[0].Integration.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_project", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionTransitionTicketTemplateParams).Project)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_transition_id", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionTransitionTicketTemplateParams).TransitionID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_advanced_fields", string(data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionTransitionTicketTemplateParams).AdvancedFields))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_comment", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionTransitionTicketTemplateParams).Comment)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_comment_on_transition", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionTransitionTicketTemplateParams).CommentOnTransition)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_attach_evidence_csv", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionTransitionTicketTemplateParams).AttachEvidenceCSV)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceWizAutomationRuleJiraTransitionTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleJiraTransitionTicketUpdate called...")

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
	vars := &wiz.UpdateAutomationRuleInput{}
	vars.ID = d.Id()
	vars.Patch.Name = d.Get("name").(string)
	vars.Patch.Description = d.Get("description").(string)
	vars.Patch.TriggerSource = d.Get("trigger_source").(string)
	triggerTypes := make([]string, 0, 0)
	for _, j := range d.Get("trigger_type").([]interface{}) {
		triggerTypes = append(triggerTypes, j.(string))
	}
	vars.Patch.TriggerType = triggerTypes
	vars.Patch.Filters = json.RawMessage(d.Get("filters").(string))
	vars.Patch.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))

	actions := []wiz.AutomationRuleActionInput{}
	jiraTransitionTicket := &wiz.JiraActionTransitionTicketTemplateParamsInput{
		Project:             d.Get("jira_project").(string),
		TransitionID:        d.Get("jira_transition_id").(string),
		AdvancedFields:      json.RawMessage(d.Get("jira_advanced_fields").(string)),
		Comment:             d.Get("jira_comment").(string),
		CommentOnTransition: utils.ConvertBoolToPointer(d.Get("jira_comment_on_transition").(bool)),
		AttachEvidenceCSV:   utils.ConvertBoolToPointer(d.Get("jira_attach_evidence_csv").(bool)),
	}
	actionTemplateParams := wiz.ActionTemplateParamsInput{
		JiraTransitionTicket: jiraTransitionTicket,
	}
	action := wiz.AutomationRuleActionInput{
		IntegrationID:        d.Get("integration_id").(string),
		ActionTemplateType:   "JIRA_TRANSITION_TICKET",
		ActionTemplateParams: actionTemplateParams,
	}
	actions = append(actions, action)
	vars.Patch.Actions = actions

	// process the request
	data := &UpdateAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule_jira_transition_ticket", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizAutomationRuleJiraTransitionTicketRead(ctx, d, m)
}
