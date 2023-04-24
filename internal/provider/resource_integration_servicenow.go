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
				Description: "ServiceNow URL. (default: none, environment variable: WIZ_INTEGRATION_SERVICENOW_URL)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_SERVICENOW_URL",
					nil,
				),
			},
			"servicenow_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email of a ServiceNow user with permissions to create tickets. (default: none, environment variable: WIZ_INTEGRATION_SERVICENOW_USERNAME)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_SERVICENOW_USERNAME",
					nil,
				),
			},
			"servicenow_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "ServiceNow password. (default: none, environment variable: WIZ_INTEGRATION_SERVICENOW_PASSWORD)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_SERVICENOW_PASSWORD",
					nil,
				),
			},
			"servicenow_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ServiceNow OAuth Client ID. (default: none, environment variable: WIZ_INTEGRATION_SERVICENOW_CLIENT_ID)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_SERVICENOW_CLIENT_ID",
					nil,
				),
			},
			"servicenow_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "ServiceNow OAuth Client Secret. (default: none, environment variable: WIZ_INTEGRATION_SERVICENOW_CLIENT_SECRET)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_SERVICENOW_CLIENT_SECRET",
					nil,
				),
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
	        authorizationType: authorization {
	          type: __typename
	        }
	        authorization {
	          ... on ServiceNowIntegrationBasicAuthorization {
	            password
	            username
	          }
	          ... on ServiceNowIntegrationOAuthAuthorization {
	            password
	            username
	            clientId
	            clientSecret
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
	data := &ReadIntegrationPayload{}
	params := &wiz.ServiceNowIntegrationParams{}
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
	err = d.Set("servicenow_url", params.URL)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("servicenow_username", params.Authorization.(map[string]interface{})["username"])
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("servicenow_password", d.Get("servicenow_password").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// determine credentials type and populate resource data
	tflog.Debug(ctx, fmt.Sprintf("params.AuthorizationType.Type %s (%T)", params.AuthorizationType.Type, params.Authorization))

	switch params.AuthorizationType.Type {
	case "ServiceNowIntegrationOAuthAuthorization":
		err = d.Set("servicenow_client_id", params.Authorization.(map[string]interface{})["servicenow_client_id"])
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		err = d.Set("servicenow_client_secret", d.Get("servicenow_client_secret").(string))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
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
	vars.Patch.Params.ServiceNow = &wiz.UpdateServiceNowIntegrationParamsInput{}
	vars.Patch.Params.ServiceNow.URL = d.Get("servicenow_url").(string)
	vars.Patch.Params.ServiceNow.Authorization.ClientID = d.Get("servicenow_client_id").(string)
	vars.Patch.Params.ServiceNow.Authorization.ClientSecret = d.Get("servicenow_client_secret").(string)
	vars.Patch.Params.ServiceNow.Authorization.Username = d.Get("servicenow_username").(string)
	vars.Patch.Params.ServiceNow.Authorization.Password = d.Get("servicenow_password").(string)

	// process the request
	data := &UpdateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_servicenow", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
