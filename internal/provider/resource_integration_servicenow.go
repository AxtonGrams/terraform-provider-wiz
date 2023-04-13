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
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizIntegrationServiceNow() *schema.Resource {
	return &schema.Resource{
		Description: "Integrations are reusable, generic connections between Wiz and third-party platforms like Slack, Google Chat, and Jira that allow data from Wiz to be passed to your preferred tool.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Identifier for this object.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the integration.",
				Required:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Identifies the date and time when the object was created.",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The project this action is scoped to.",
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "All Resources, Restrict this Integration to global roles only",
				Description: fmt.Sprintf(
					"Scoping to a selected Project makes this Integration accessible only to users with global roles or Project-scoped access to the selected Project. Other users will not be able to see it, use it, or view its results. Integrations restricted to global roles cannot be seen or used by users with Project-scoped roles. \n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						internal.IntegrationScope,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						internal.IntegrationScope,
						false,
					),
				),
			},
			"servicenow_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ServiceNow URL.",
			},
			"servicenow_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email of a ServiceNow user with permissions to create tickets",
			},
			"servicenow_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"servicenow_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"servicenow_client_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceWizIntegrationAwsServiceNowCreate,
		ReadContext:   resourceWizIntegrationAwsServiceNowRead,
		UpdateContext: resourceWizIntegrationAwsServiceNowUpdate,
		DeleteContext: resourceWizIntegrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizIntegrationAwsServiceNowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsServiceNowCreate called...")

	// define the graphql query
	query := `mutation CreateIntegration($input: CreateIntegrationInput!) {
	  createIntegration(
	    input: $input
	  ) {
	    integration {
	      id
	    }
	  }
	}`

	vars := &wiz.CreateIntegrationInput{}
	vars.Name = d.Get("name").(string)
	vars.Type = "SERVICE_NOW"
	vars.ProjectID = d.Get("project_id").(string)
	vars.IsAccessibleToAllProjects = convertIntegrationScopeToBool(d.Get("scope").(string))
	vars.Params.ServiceNow = &wiz.CreateServiceNowIntegrationParamsInput{}
	vars.Params.ServiceNow.URL = d.Get("servicenow_url").(string)
	vars.Params.ServiceNow.Authorization.Username = d.Get("servicenow_username").(string)
	vars.Params.ServiceNow.Authorization.Password = d.Get("servicenow_password").(string)
	vars.Params.ServiceNow.Authorization.ClientID = d.Get("servicenow_client_id").(string)
	vars.Params.ServiceNow.Authorization.ClientSecret = d.Get("servicenow_client_secret").(string)

	// process the request
	data := &CreateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_servicenow", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateIntegration.Integration.ID)

	return resourceWizIntegrationAwsServiceNowRead(ctx, d, m)
}

func resourceWizIntegrationAwsServiceNowRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsServiceNowRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query integration (
	  $id: ID!
	) {
	  integration(
	    id: $id
	  ) {
	    id
	    name
	    createdAt
	    updatedAt
	    project {
	      id
	    }
	    type
	    isAccessibleToAllProjects
	    usedByRules {
	      id
	    }
	    paramsType: params {
	      type: __typename
	    }
	    params {
	      ... on ServiceNowIntegrationParams {
	        url
	      }
	    }
	  }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	data := &ReadIntegrationPayload{}
	params := &wiz.AwsSNSIntegrationParams{}
	data.Integration.Params = params
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_servicenow", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.Integration.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.Integration.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("created_at", data.Integration.CreatedAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", data.Integration.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("aws_sns_topic_arn", params.TopicARN)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("aws_sns_access_method", params.AccessMethod)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("aws_sns_connector_id", params.AccessConnector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("aws_sns_customer_role_arn", params.CustomerRoleARN)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceWizIntegrationAwsServiceNowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsServiceNowUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateIntegration(
	  $input: UpdateIntegrationInput!
	) {
	  updateIntegration(input: $input) {
	    integration {
	      id
	    }
	  }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateIntegrationInput{}
	vars.ID = d.Id()
	vars.Patch.Name = d.Get("name").(string)
	vars.Patch.Params.AwsSNS = &wiz.UpdateAwsSNSIntegrationParamsInput{}
	vars.Patch.Params.AwsSNS.TopicARN = d.Get("aws_sns_topic_arn").(string)
	vars.Patch.Params.AwsSNS.AccessMethod.Type = d.Get("aws_sns_access_method").(string)
	vars.Patch.Params.AwsSNS.AccessMethod.AccessConnectorID = d.Get("aws_sns_connector_id").(string)
	vars.Patch.Params.AwsSNS.AccessMethod.CustomerRoleARN = d.Get("aws_sns_customer_role_arn").(string)

	// process the request
	data := &UpdateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_servicenow", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
