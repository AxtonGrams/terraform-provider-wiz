package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizConnectorAws() *schema.Resource {
	return &schema.Resource{
		Description: "Connectors are used to connect AWS resources to Wiz.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Wiz internal identifier for the connector.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The connector name.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the connector is enabled.",
				Optional:    true,
				Default:     true,
			},
			"customer_role_arn": {
				Type:        schema.TypeString,
				Description: "The AWS customer role arn for Wiz to assume.",
				Computed:    true,
			},
			"excluded_accounts": {
				Type:        schema.TypeList,
				Description: "The AWS accounts excluded from the connector.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excluded_ous": {
				Type:        schema.TypeList,
				Description: "The AWS OUs excluded from the connector.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"audit_log_monitor_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether audit log monitor is enabled. Note an advanced license is required.",
				Computed:    true,
			},
			"disk_analyzer_inflight_disabled": {
				Type:        schema.TypeBool,
				Description: "If using Outpost, whether disk analyzer inflight scanning is disabled.",
				Computed:    true,
			},
			"skip_organization_scan": {
				Type:        schema.TypeBool,
				Description: "Whether to skip the organization scan (account-scoped only).",
				Computed:    true,
			},
			"external_id_nonce": {
				Type:        schema.TypeString,
				Description: "The AWS external ID / nonce, this will be used for IAM-related dependencies (`sts:ExternalId` conditional trust policies).",
				Computed:    true,
			},
			"opted_in_regions": {
				Type:        schema.TypeList,
				Description: "The AWS regions opted in for the connector.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The AWS region for the connector.",
				Computed:    true,
			},
			"events_cloudtrail_bucket_name": {
				Type:        schema.TypeString,
				Description: "If using Wiz Cloud Events, the CloudTrail bucket name.",
				Computed:    true,
			},
			"events_cloudtrail_bucket_sub_account": {
				Type:        schema.TypeString,
				Description: "If using Wiz Cloud Events and CloudTrail is organizational, the CloudTrail bucket sub account.",
				Computed:    true,
			},
			"events_cloudtrail_organization": {
				Type:        schema.TypeString,
				Description: "If using Wiz Cloud Events and CloudTrail is deployed to AWS organizations, the organizational ID.",
				Computed:    true,
			},
			"auth_params": {
				Type:        schema.TypeString,
				Description: "The authentication parameters. Must be represented in `JSON` format.",
				Required:    true,
				Sensitive:   true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
			"extra_config": {
				// these are JSON fields; the schema does not support overrides, once a field is set, future changes require it to be passed
				Type:        schema.TypeString,
				Description: "Extra configuration for the connector. Must be represented in `JSON` format.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
		},
		// auth_params requires a resource recreation as they cannot be updated.
		// to accommodate for importing resources into state, we can't use `ForceNew` in the schema definition.
		// we use a customdiff and `ForceNewIfChange` for below attributes for only a change condition.
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("auth_params", func(ctx context.Context, old, new, meta any) bool {
				if old.(string) != "" {
					return old.(string) != new.(string)
				}
				return false
			},
			),
		),
		CreateContext: resourceWizConnectorCreate,
		ReadContext:   resourceWizConnectorRead,
		UpdateContext: resourceWizConnectorUpdate,
		DeleteContext: resourceWizConnectorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateConnector struct
type CreateConnector struct {
	CreateConnector wiz.CreateConnectorPayload `json:"createConnector"`
}

func resourceWizConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorCreate called...")

	query := `mutation CreateConnector($input: CreateConnectorInput!) {
				createConnector(input: $input) {
					connector {
						id
					}
				}
			  }
		     `
	// populate the graphql variables
	vars := &wiz.CreateConnectorInput{}
	vars.Name = d.Get("name").(string)
	enabled := d.Get("enabled").(bool)
	vars.Type = "aws"
	vars.Enabled = &enabled

	vars.AuthParams = json.RawMessage(d.Get("auth_params").(string))
	vars.ExtraConfig = json.RawMessage(d.Get("extra_config").(string))

	// process the request
	data := &CreateConnector{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "connector", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateConnector.Connector.ID)

	return resourceWizConnectorRead(ctx, d, m)
}

// ReadConnectorPayload struct
type ReadConnectorPayload struct {
	Connector wiz.Connector `json:"connector"`
}

func resourceWizConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query GetConnector($id: ID!) {
	    connector(id: $id) {
	      id
	      name
	      enabled
	      authParams
	      extraConfig
	      config {
	        ... on ConnectorConfigAWS {
	          region
	          diskAnalyzerInFlightDisabled
	          excludedAccounts
	          excludedOUs
	          externalIdNonce
	          optedInRegions
	          customerRoleARN
	          auditLogMonitorEnabled
	          skipOrganizationScan
	          cloudTrailConfig {
	            bucketName
	            bucketSubAccount
	            trailOrg
	          }
	        }
	      }
	      type {
	        ...ConnectorTypeFrag
	      }
	    }
	  }

	  fragment ConnectorTypeFrag on ConnectorType {
	    id
	    name
	    authorizeUrls
	  }
`
	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	data := &ReadConnectorPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "connector", "read")
	diags = append(diags, requestDiags...,
	)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
		// we are checking if the response does not have an ID and if so we mark the resource as new
		// once the vendor implements consistent and documented error handling we can
		// move to an error handling factory to inspect definitive errors and handle accordingly
		if data.Connector.ID == "" {
			tflog.Info(ctx, "resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	err := d.Set("name", data.Connector.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("id", data.Connector.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("enabled", data.Connector.Enabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	var mapExtraConfig map[string]interface{}
	err = json.Unmarshal(data.Connector.ExtraConfig, &mapExtraConfig)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// When certain fields are deprecated, vendor returns null and/or empty string (i.e. cloudTrailConfig)
	// We need to handle to avoid unwanted diffs, below traverses the map to a maximum depth of 5 levels
	utils.RemoveNullAndEmptyValues(mapExtraConfig, 5)

	extraConfig, err := json.Marshal(mapExtraConfig)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("extra_config", string(extraConfig))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	var connectorConfig wiz.ConnectorConfigAWS
	connectorConfigBytes, err := json.Marshal(data.Connector.Config)
	if err != nil {
		return append(diags, diag.Errorf("unable to marshal ConnectorConfigAWS: %v", err)...)
	}
	if err := json.Unmarshal(connectorConfigBytes, &connectorConfig); err != nil {
		return append(diags, diag.Errorf("unable to unmarshal ConnectorConfigAWS: %v", err)...)
	}

	err = d.Set("events_cloudtrail_bucket_name", connectorConfig.CloudTrailConfig.BucketName)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("events_cloudtrail_bucket_sub_account", connectorConfig.CloudTrailConfig.BucketSubAccount)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("events_cloudtrail_organization", connectorConfig.CloudTrailConfig.TrailOrg)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("external_id_nonce", connectorConfig.ExternalIDNonce)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("customer_role_arn", connectorConfig.CustomerRoleARN)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("disk_analyzer_inflight_disabled", connectorConfig.DiskAnalyzerInFlightDisabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("audit_log_monitor_enabled", connectorConfig.AuditLogMonitorEnabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("skip_organization_scan", connectorConfig.SkipOrganizationScan)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("excluded_accounts", utils.ConvertSliceToGenericArray(connectorConfig.ExcludedAccounts))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("excluded_ous", utils.ConvertSliceToGenericArray(connectorConfig.ExcludedOUs))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("opted_in_regions", utils.ConvertSliceToGenericArray(connectorConfig.OptedInRegions))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("region", connectorConfig.Region)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateConnector struct
type UpdateConnector struct {
	UpdateConnector wiz.UpdateConnectorPayload `json:"updateConnector"`
}

func resourceWizConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateConnector($input: UpdateConnectorInput!) {
	    updateConnector(input: $input) {
	      connector {
	        id
	        name
	        enabled
	        extraConfig
	      }
	    }
	  }`

	// populate the graphql variables
	vars := &wiz.UpdateConnectorInput{}
	vars.ID = d.Id()

	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		vars.Patch.Enabled = &enabled
	}
	if d.HasChange("extra_config") {
		vars.Patch.ExtraConfig = json.RawMessage(d.Get("extra_config").(string))
	}

	// process the request
	data := &UpdateConnector{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "connector", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizConnectorRead(ctx, d, m)
}

// DeleteConnector struct
type DeleteConnector struct {
	DeleteConnector wiz.DeleteConnectorPayload `json:"_stub"`
}

func resourceWizConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteConnector($input: DeleteConnectorInput!) {
		deleteConnector(input: $input) {
		  _stub
		}
	  }
	`
	// populate the graphql variables
	vars := &wiz.DeleteConnectorInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateUser{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "connector", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
