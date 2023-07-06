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

func resourceWizIntegrationJira() *schema.Resource {
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
			"jira_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Jira URL. (default: none, environment variable: WIZ_INTEGRATION_JIRA_URL)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_JIRA_URL",
					nil,
				),
			},
			"jira_server_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Jira server type",
				Default:     "CLOUD",
			},
			"jira_is_on_prem": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether Jira instance is on prem",
				Default:     false,
			},
			"jira_allow_insecure_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Jira integration TLS setting",
			},
			"jira_server_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Jira server CA",
			},
			"jira_client_certificate_and_private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Jira PEM with client certificate and private key",
			},
			"jira_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email of a Jira user with permissions to create tickets. (default: none, environment variable: WIZ_INTEGRATION_JIRA_USERNAME)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_JIRA_USERNAME",
					nil,
				),
			},
			"jira_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Jira password. (default: none, environment variable: WIZ_INTEGRATION_JIRA_PASSWORD)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_JIRA_PASSWORD",
					nil,
				),
			},
			"jira_pat": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Jira personal access token (used for on-prem). (default: none, environment variable: WIZ_INTEGRATION_JIRA_PAT)",
				DefaultFunc: schema.EnvDefaultFunc(
					"WIZ_INTEGRATION_JIRA_PAT",
					nil,
				),
			},
		},
		CreateContext: resourceWizIntegrationJiraCreate,
		ReadContext:   resourceWizIntegrationJiraRead,
		UpdateContext: resourceWizIntegrationJiraUpdate,
		DeleteContext: resourceWizIntegrationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizIntegrationJiraCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationJiraCreate called...")

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
	vars.Type = "JIRA"
	vars.ProjectID = d.Get("project_id").(string)
	vars.IsAccessibleToAllProjects = convertIntegrationScopeToBool(d.Get("scope").(string))
	vars.Params.Jira = &wiz.CreateJiraIntegrationParamsInput{}
	vars.Params.Jira.ServerURL = d.Get("jira_url").(string)
	vars.Params.Jira.ServerType = d.Get("jira_server_type").(string)
	vars.Params.Jira.IsOnPrem = d.Get("jira_is_on_prem").(bool)
	vars.Params.Jira.TLSConfig.AllowInsecureTLS = utils.ConvertBoolToPointer(d.Get("jira_allow_insecure_tls").(bool))
	vars.Params.Jira.TLSConfig.ClientCertificateAndPrivateKey = d.Get("jira_client_certificate_and_private_key").(string)
	vars.Params.Jira.TLSConfig.ServerCA = d.Get("jira_server_ca").(string)
	vars.Params.Jira.Authorization.Username = d.Get("jira_username").(string)
	vars.Params.Jira.Authorization.Password = d.Get("jira_password").(string)
	vars.Params.Jira.Authorization.PersonalAccessToken = d.Get("jira_pat").(string)

	// process the request
	data := &CreateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_jira", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateIntegration.Integration.ID)

	return resourceWizIntegrationJiraRead(ctx, d, m)
}

func resourceWizIntegrationJiraRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationJiraRead called...")

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
	      ... on JiraIntegrationParams {
	        url
			serverType
			onPremConfig {
				isOnPrem
			}
			tlsConfig {
				allowInsecureTLS
				serverCA
				clientCertificateAndPrivateKey
			}
	        authorization {
	          ... on JiraIntegrationBasicAuthorization {
	            password
	            username
	          }
			  ... on JiraIntegrationTokenBearerAuthorization {
				token
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
	params := &wiz.JiraIntegrationParams{}
	data.Integration.Params = params
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_jira", "read")
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
	err = d.Set("jira_url", params.URL)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_server_type", params.ServerType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_is_on_prem", params.OnPremConfig.IsOnPrem)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_allow_insecure_tls", params.TLSConfig.AllowInsecureTLS)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_server_ca", params.TLSConfig.ServerCA)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_client_certificate_and_private_key", params.TLSConfig.ClientCertificateAndPrivateKey)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_username", params.Authorization.(map[string]interface{})["username"])
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_password", d.Get("jira_password").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("jira_pat", d.Get("jira_pat").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceWizIntegrationJiraUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationJiraUpdate called...")

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
	vars.Patch.Params.Jira = &wiz.UpdateJiraIntegrationParamsInput{}
	vars.Patch.Params.Jira.ServerURL = d.Get("jira_url").(string)
	vars.Patch.Params.Jira.ServerType = d.Get("jira_server_type").(string)
	vars.Patch.Params.Jira.IsOnPrem = utils.ConvertBoolToPointer(d.Get("jira_is_on_prem").(bool))
	vars.Patch.Params.Jira.TLSConfig.AllowInsecureTLS = utils.ConvertBoolToPointer(d.Get("jira_allow_insecure_tls").(bool))
	vars.Patch.Params.Jira.TLSConfig.ServerCA = d.Get("jira_server_ca").(string)
	vars.Patch.Params.Jira.TLSConfig.ClientCertificateAndPrivateKey = d.Get("jira_client_certificate_and_private_key").(string)
	vars.Patch.Params.Jira.Authorization.Username = d.Get("jira_username").(string)
	vars.Patch.Params.Jira.Authorization.Password = d.Get("jira_password").(string)
	vars.Patch.Params.Jira.Authorization.PersonalAccessToken = d.Get("jira_pat").(string)

	// process the request
	data := &UpdateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_jira", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
