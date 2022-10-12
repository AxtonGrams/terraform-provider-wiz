package provider

import (
	"context"
	//	"encoding/json"
	//	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	//	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	//	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	//	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func dataSourceWizCloudAccounts() *schema.Resource {
	return &schema.Resource{
		Description: "Query cloud accounts (subscriptions).",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Internal identifier for the Control",
			},
			"filters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Get specific Cloud Accounts by their IDs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"search": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Free text search on cloud account name or tags or external-id.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Query cloud accounts of a specific linked project, given its id.",
						},
						"cloud_provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Query cloud accounts of specific cloud provider.",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.CloudProvider,
									false,
								),
							),
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Query cloud accounts by it's status.",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.CloudAccountStatus,
									false,
								),
							),
						},
						"connector_id": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Query cloud accounts by specific connector ID.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"connector_issue_id": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Query cloud accounts by specific connector issue ID.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"assigned_to_project": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "QueryQuery cloud accounts by project assignment state.",
						},
						"has_multiple_connector_source": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "QueryQuery cloud accounts by project assignment state.",
						},
					},
				},
			},
			"cloud_accounts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal Wiz ID.",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External subscription id from cloud provider (subscriptionId in security graph).",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name for this account.",
						},
						"cloud_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_connectors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connectors detected this cloud account.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud Account connectivity status as affected by configured connectors.",
						},
					},
				},
			},
		},
		ReadContext: dataSourceWizCloudAccountsRead,
	}
}

func dataSourceWizCloudAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizOrganizationsRead called...")

	return diags
}
