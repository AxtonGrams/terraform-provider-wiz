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

func resourceWizAutomationRuleJiraCreateTicket() *schema.Resource {
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
					"Trigger type. Must be set to `CREATED` for wiz_automation_rule_jira_create_ticket.\n    - Allowed values: %s",
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
			"jira_summary": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Wiz Issue: {{control.name}}",
				Description: "Issue summary",
			},
			"jira_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue description",
				Default:     `Description:  {{issue.description}}\nStatus:       {{issue.status}}\nCreated:      {{issue.createdAt}}\nSeverity:     {{issue.severity}}\nProject:      {{#issue.projects}}{{name}}, {{/issue.projects}}\n\n---\nResource:\t            {{issue.entitySnapshot.name}}\nType:\t                {{issue.entitySnapshot.nativeType}}\nCloud Platform:\t        {{issue.entitySnapshot.cloudPlatform}}\nCloud Resource URL:     {{issue.entitySnapshot.cloudProviderURL}}\nSubscription Name (ID): {{issue.entitySnapshot.subscriptionName}} ({{issue.entitySnapshot.subscriptionExternalId}})\nRegion:\t                {{issue.entitySnapshot.region}}\nPlease click the following link to proceed to investigate the issue:\nhttps://{{wizDomain}}/issues#~(issue~'{{issue.id}})\nSource Automation Rule: {{ruleName}}`,
			},
			"jira_issue_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue type",
				Default:     "Vulnerability",
			},
			"jira_assignee": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue assignee",
			},
			"jira_components": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Issue components",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jira_fix_version": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Issue fix versions",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jira_labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Issue labels",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jira_priority": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue priority",
			},
			"jira_project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue project",
			},
			"jira_alternative_description_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issue alternative description field",
			},
			"jira_custom_fields": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom configuration fields as specified in Jira. Make sure you add the fields that are configured as required in Jira Project, otherwise ticket creation will fail. Must be valid JSON.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
			"jira_attach_evidence_csv": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Upload issue evidence CSV as attachment?",
			},
		},
		CreateContext: resourceWizAutomationRuleJiraCreateTicketCreate,
		ReadContext:   resourceWizAutomationRuleJiraCreateTicketRead,
		UpdateContext: resourceWizAutomationRuleJiraCreateTicketUpdate,
		DeleteContext: resourceWizAutomationRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizAutomationRuleJiraCreateTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleJiraCreateTicketCreate called...")

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
	jiraCreateTicketFields := wiz.CreateJiraTicketFieldsInput{
		Summary:                     d.Get("jira_summary").(string),
		Description:                 d.Get("jira_description").(string),
		IssueType:                   d.Get("jira_issue_type").(string),
		Assignee:                    d.Get("jira_assignee").(string),
		Components:                  utils.ConvertListToString(d.Get("jira_components").([]interface{})),
		FixVersion:                  utils.ConvertListToString(d.Get("jira_fix_version").([]interface{})),
		Labels:                      utils.ConvertListToString(d.Get("jira_labels").([]interface{})),
		Priority:                    d.Get("jira_priority").(string),
		Project:                     d.Get("jira_project").(string),
		AlternativeDescriptionField: d.Get("jira_alternative_description_field").(string),
		CustomFields:                json.RawMessage(d.Get("jira_custom_fields").(string)),
		AttachEvidenceCSV:           utils.ConvertBoolToPointer(d.Get("jira_attach_evidence_csv").(bool)),
	}
	jiraCreateTicketParams := &wiz.JiraActionCreateTicketTemplateParamsInput{
		Fields: jiraCreateTicketFields,
	}
	actionTemplateParams := wiz.ActionTemplateParamsInput{
		JiraCreateTicket: jiraCreateTicketParams,
	}
	actions := []wiz.AutomationRuleActionInput{}
	action := wiz.AutomationRuleActionInput{
		IntegrationID:        d.Get("integration_id").(string),
		ActionTemplateParams: actionTemplateParams,
		ActionTemplateType:   "JIRA_CREATE_TICKET",
	}
	actions = append(actions, action)
	vars.Actions = actions

	// process the request
	data := &CreateAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule_jira_create_ticket", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id and computed values
	d.SetId(data.CreateAutomationRule.AutomationRule.ID)

	return resourceWizAutomationRuleJiraCreateTicketRead(ctx, d, m)
}

func resourceWizAutomationRuleJiraCreateTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleJiraCreateTicketRead called...")

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
	        ... on JiraActionCreateTicketTemplateParams {
	          fields {
				summary
				description
				issueType
				assignee
				components
				fixVersion
				labels
				priority
				project
				alternativeDescriptionField
				customFields
				attachEvidenceCSV
	          }
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
		ActionTemplateParams: &wiz.JiraActionCreateTicketTemplateParams{},
	}
	automationRuleActions = append(automationRuleActions, automationRuleAction)
	data := &ReadAutomationRulePayload{
		AutomationRule: wiz.AutomationRule{
			Actions: automationRuleActions,
		},
	}

	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule_jira_create_ticket", "read")
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
	err = d.Set("jira_summary", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.Summary)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_description", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_issue_type", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.IssueType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_assignee", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.Assignee)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_components", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.Components)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_fix_version", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.FixVersion)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_labels", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.Labels)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_priority", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.Priority)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_project", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.Project)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_alternative_description_field", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.AlternativeDescriptionField)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if string(data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.CustomFields) != "null" {
		err = d.Set("jira_custom_fields", string(data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.CustomFields))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	err = d.Set("jira_attach_evidence_csv", data.AutomationRule.Actions[0].ActionTemplateParams.(*wiz.JiraActionCreateTicketTemplateParams).Fields.AttachEvidenceCSV)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceWizAutomationRuleJiraCreateTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleJiraCreateTicketUpdate called...")

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
	jiraFields := wiz.CreateJiraTicketFieldsInput{
		Summary:                     d.Get("jira_summary").(string),
		Description:                 d.Get("jira_description").(string),
		IssueType:                   d.Get("jira_issue_type").(string),
		Assignee:                    d.Get("jira_assignee").(string),
		Components:                  utils.ConvertListToString(d.Get("jira_components").([]interface{})),
		FixVersion:                  utils.ConvertListToString(d.Get("jira_fix_version").([]interface{})),
		Labels:                      utils.ConvertListToString(d.Get("jira_labels").([]interface{})),
		Priority:                    d.Get("jira_priority").(string),
		Project:                     d.Get("jira_project").(string),
		AlternativeDescriptionField: d.Get("jira_alternative_description_field").(string),
		CustomFields:                json.RawMessage(d.Get("jira_custom_fields").(string)),
		AttachEvidenceCSV:           utils.ConvertBoolToPointer(d.Get("jira_attach_evidence_csv").(bool)),
	}
	jiraCreateTicket := &wiz.JiraActionCreateTicketTemplateParamsInput{
		Fields: jiraFields,
	}

	actionTemplateParams := wiz.ActionTemplateParamsInput{
		JiraCreateTicket: jiraCreateTicket,
	}
	action := wiz.AutomationRuleActionInput{
		IntegrationID:        d.Get("integration_id").(string),
		ActionTemplateType:   "JIRA_CREATE_TICKET",
		ActionTemplateParams: actionTemplateParams,
	}
	actions = append(actions, action)
	vars.Patch.Actions = actions

	// process the request
	data := &UpdateAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule_jira_create_ticket", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizAutomationRuleJiraCreateTicketRead(ctx, d, m)
}
