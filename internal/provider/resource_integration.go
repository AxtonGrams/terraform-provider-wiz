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
				Description: "Wiz internal identifier.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the integration.",
				Required:    true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Type of integration action.\n    - Allowed values: %s",
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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"is_on_prem": {
							Type:     schema.TypeBool,
							Optional: true,
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
						"note": {
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
						"note": {
							Type:     schema.TypeString,
							Optional: true,
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
