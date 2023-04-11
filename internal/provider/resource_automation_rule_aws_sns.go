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

func resourceWizAutomationRuleAwsSns() *schema.Resource {
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
			"action_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"integration_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_sns_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CreateContext: resourceWizAutomationRuleAwsSNSCreate,
		ReadContext:   resourceWizAutomationRuleAwsSNSRead,
		UpdateContext: resourceWizAutomationRuleAwsSNSUpdate,
		DeleteContext: resourceWizAutomationRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAutomationRuleAwsSNS struct
type CreateAutomationRule struct {
	CreateAutomationRule vendor.CreateAutomationRulePayload `json:"createAutomationRule"`
}

func resourceWizAutomationRuleAwsSNSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleCreate called...")

	return resourceWizAutomationRuleAwsSNSRead(ctx, d, m)
}

func resourceWizAutomationRuleAwsSNSRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleRead called...")

	return diags
}

func resourceWizAutomationRuleAwsSNSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleUpdate called...")

	return resourceWizAutomationRuleAwsSNSRead(ctx, d, m)
}
