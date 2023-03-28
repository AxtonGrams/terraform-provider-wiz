package provider

import (
	"context"
	//"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	//"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func resourceWizIntegrationAwsSNS() *schema.Resource {
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
				Description: "The project this action is scoped to.",
			},
			"is_accessible_to_all_projects": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, any scoped and non scoped project users will be able to use this action, false by default.",
			},
			"aws_sns_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is EMAIL, define these paramemters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_arn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The SNS Topic Arn.",
						},
						"access_method": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Required: true,
										Type:     schema.TypeString,
										Description: fmt.Sprintf(
											"The access method this integration should use. \n    - Allowed values: %s",
											utils.SliceOfStringToMDUList(
												vendor.AwsSNSIntegrationAccessMethodType,
											),
										),
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.StringInSlice(
												vendor.IntegrationType,
												false,
											),
										),
									},
									"access_connector_id": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "Required if and only if accessMethod is ASSUME_CONNECTOR_ROLE, this should be a valid existing AWS connector ID from which the role ARN will be taken.",
									},
									"customer_role_arn": {
										Optional:    true,
										Type:        schema.TypeString,
										Description: "Required if and only if accessMethod is ASSUME_SPECIFIED_ROLE, this is the role that should be assumed, the ExternalID of the role must be your Wiz Tenant ID (a GUID).",
									},
								},
							},
						},
					},
				},
			},
		},
		CreateContext: resourceWizIntegrationAwsSNSCreate,
		ReadContext:   resourceWizIntegrationAwsSNSRead,
		UpdateContext: resourceWizIntegrationAwsSNSUpdate,
		DeleteContext: resourceWizIntegrationAwsSNSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func getAwsSNSIntegrationVar(ctx context.Context, d *schema.ResourceData) vendor.CreateAwsSNSIntegrationParamsInput {
	params := d.Get("aws_sns_params").(*schema.Set).List()
	var myParams vendor.CreateAwsSNSIntegrationParamsInput
	for _, y := range params {
		for a, b := range y.(map[string]interface{}) {
			tflog.Debug(ctx, fmt.Sprintf("a: %T %s", a, utils.PrettyPrint(a)))
			tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		}
	}
	return myParams
}

func resourceWizIntegrationAwsSNSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsSNSCreate called...")

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

	vars := &vendor.CreateIntegrationInput{}
	vars.Name = d.Get("name").(string)
	vars.Type = "AWS_SNS"
	vars.ProjectID = d.Get("project_id").(string)
	vars.IsAccessibleToAllProjects = utils.ConvertBoolToPointer(d.Get("is_accessible_to_all_projects").(bool))
	vars.Params.AwsSNS = getAwsSNSIntegrationVar(ctx, d)

	// process the request
	data := &CreateIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration_aws_sns", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateIntegration.Integration.ID)

	return resourceWizIntegrationAwsSNSRead(ctx, d, m)
}

func resourceWizIntegrationAwsSNSRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsSNSRead called...")

	return diags
}

func resourceWizIntegrationAwsSNSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsSNSUpdate called...")

	return diags
}
