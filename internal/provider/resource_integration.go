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
	//"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func resourceWizIntegration() *schema.Resource {
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Type of integration action. The following are implemented: AWS_SNS, JIRA, PAGER_DUTY, SERVICE_NOW, WEBHOOK.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.IntegrationType,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						vendor.IntegrationType,
						false,
					),
				),
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
				ConflictsWith: []string{
					"webhook_params",
					"pagerduty_params",
					"servicenow_params",
				},
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
										Required: true,
										Type:     schema.TypeString,
									},
									"customer_role_arn": {
										Required: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
			"webhook_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is EMAIL, define these paramemters.",
				ConflictsWith: []string{
					"aws_sns_params",
					"pagerduty_params",
					"servicenow_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL of the webhook.",
						},
						"is_on_prem": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"authorization": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"password": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"token": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"headers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"tls_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_insecurity_tls": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Setting this to true will ignore any TLS validation errors on the server side certificate Warning: should only be used to validate that the action works regardless of TLS validation, if for example your server is presenting self signed or expired TLS certificate.",
									},
									"server_ca": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A PEM of the certificate authority that your server presents (if you use self signed, or custom CA).",
									},
									"client_certificate_and_private_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A PEM of the client certificate as well as the certificate private key.",
									},
								},
							},
						},
					},
				},
			},
			"pagerduty_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is EMAIL, define these paramemters.",
				ConflictsWith: []string{
					"aws_sns_params",
					"webhook_params",
					"servicenow_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"integration_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"servicenow_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is EMAIL, define these paramemters.",
				ConflictsWith: []string{
					"aws_sns_params",
					"webhook_params",
					"pagerduty_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authorization": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Email of a ServiceNow user with permissions to create tickets.",
									},
									"password": {
										Type:     schema.TypeString,
										Required: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"client_secret": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
		CreateContext: resourceWizIntegrationCreate,
		ReadContext:   resourceWizIntegrationRead,
		UpdateContext: resourceWizIntegrationUpdate,
		DeleteContext: resourceWizIntegrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationCreate called...")

	return diags
}

func resourceWizIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationRead called...")

	return diags
}

func resourceWizIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationUpdate called...")

	return diags
}

func resourceWizIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationDelete called...")

	return diags
}
