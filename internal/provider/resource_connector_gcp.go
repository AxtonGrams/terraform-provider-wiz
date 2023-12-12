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

const (
	auditLogMonitorEnabledKey = "auditLogMonitorEnabled"
	auditLogsConfigKey        = "auditLogsConfig"
	pubSubKey                 = "pub_sub"
	projectIDKey              = "project_id"
	topicIDKey                = "topic_id"
	subscriptionIDKey         = "subscription_id"
	topicNameKey              = "topicName"
	subscriptionIDNameKey     = "subscriptionID"
)

func resourceWizConnectorGcp() *schema.Resource {
	return &schema.Resource{
		Description: "Connectors are used to connect GCP resources to Wiz.",
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
			"is_managed_identity": {
				Type:        schema.TypeString,
				Description: "Is managed identity?",
				Computed:    true,
			},
			"folder_id": {
				Type:        schema.TypeString,
				Description: "The GCP folder ID.",
				Computed:    true,
			},
			"organization_id": {
				Type:        schema.TypeString,
				Description: "The GCP organization ID.",
				Computed:    true,
			},
			"projects": {
				Type:        schema.TypeList,
				Description: "The GCP projects to target with the connector.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excluded_projects": {
				Type:        schema.TypeList,
				Description: "The GCP projects excluded by the connector.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"included_folders": {
				Type:        schema.TypeList,
				Description: "The GCP folders included by the connector.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excluded_folders": {
				Type:        schema.TypeList,
				Description: "The GCP folders excluded by the connector.",
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
			"events_topic_name": {
				Type:        schema.TypeString,
				Description: "If using Wiz Cloud Events, the Topic Name in format `projects/<project_id>/topics/<topic_id>`.",
				Computed:    true,
			},
			"events_pub_sub_subscription_id": {
				Type:        schema.TypeString,
				Description: "If using Wiz Cloud Events, the Pub/Sub Subscription ID.",
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
		CreateContext: resourceWizConnectorGcpCreate,
		ReadContext:   resourceWizConnectorGcpRead,
		UpdateContext: resourceWizConnectorGcpUpdate,
		DeleteContext: resourceWizConnectorGcpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizConnectorGcpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorGcpCreate called...")

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
	vars.Type = "gcp"
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

	return resourceWizConnectorGcpRead(ctx, d, m)
}

func resourceWizConnectorGcpRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorGcpRead called...")

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
	        ... on ConnectorConfigGCP {
	          auditLogMonitorEnabled
	          diskAnalyzerInFlightDisabled
	          includedFolders
	          excludedFolders
	          excludedProjects
	          delegateUser
	          projects
	          organization_id
	          project_id
	          folder_id
	          auditLogMonitorEnabled
	          auditLogsConfig {
	            pub_sub {
	              topicName
	              subscriptionID
	            }
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
	diags = append(diags, requestDiags...)
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

	var connectorConfig wiz.ConnectorConfigGCP
	connectorConfigBytes, err := json.Marshal(data.Connector.Config)
	if err != nil {
		return append(diags, diag.Errorf("unable to marshal ConnectorConfigGCP: %v", err)...)
	}
	if err := json.Unmarshal(connectorConfigBytes, &connectorConfig); err != nil {
		return append(diags, diag.Errorf("unable to unmarshal ConnectorConfigGCP: %v", err)...)
	}

	err = d.Set("projects", utils.ConvertSliceToGenericArray(connectorConfig.Projects))
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
	err = d.Set("organization_id", connectorConfig.OrganizationID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("folder_id", connectorConfig.FolderID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("excluded_folders", utils.ConvertSliceToGenericArray(connectorConfig.ExcludedFolders))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("included_folders", utils.ConvertSliceToGenericArray(connectorConfig.IncludedFolders))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("excluded_projects", utils.ConvertSliceToGenericArray(connectorConfig.ExcludedProjects))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("events_topic_name", connectorConfig.AuditLogsConfig.PubSub.TopicName)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("events_pub_sub_subscription_id", connectorConfig.AuditLogsConfig.PubSub.SubscriptionID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	diagsExtraConfig := updateExtraConfig(ctx, *data, d, diags)
	if diagsExtraConfig != nil {
		return append(diags, diagsExtraConfig...)
	}

	return diags
}

func resourceWizConnectorGcpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorGcpUpdate called...")

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

	return resourceWizConnectorGcpRead(ctx, d, m)
}

func resourceWizConnectorGcpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizConnectorGcpDelete called...")

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

// Wiz API limitations prevent nullifying the `auditLogsConfig/pub_sub` field. For example, disabling `auditLogMonitorEnabled`
// will results in perpetual drift detection of extraConfig as the related pub_sub information will always be in the response once set.
// Furthermore, additional `pub_sub `fields require normalization and removal of unnecessary fields.
func updateExtraConfig(ctx context.Context, data ReadConnectorPayload, d *schema.ResourceData, diags diag.Diagnostics) diag.Diagnostics {
	tflog.Info(ctx, "updateExtraConfig called...")

	var mapExtraConfig map[string]interface{}
	err := json.Unmarshal(data.Connector.ExtraConfig, &mapExtraConfig)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	auditLogMonitoringEnabled, ok := mapExtraConfig[auditLogMonitorEnabledKey].(bool)
	if !ok {
		auditLogMonitoringEnabled = true
	}

	if auditLogMonitoringEnabled {
		if pubSubConfig, ok := mapExtraConfig[auditLogsConfigKey].(map[string]interface{})[pubSubKey].(map[string]interface{}); ok {
			projectID, ok := pubSubConfig[projectIDKey].(string)
			if !ok {
				return addFieldError(diags, projectIDKey, pubSubKey)
			}
			topicID, ok := pubSubConfig[topicIDKey].(string)
			if !ok {
				return addFieldError(diags, topicIDKey, pubSubKey)
			}
			topicName := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)
			subscriptionID, ok := pubSubConfig[subscriptionIDKey].(string)
			if !ok {
				return addFieldError(diags, subscriptionIDKey, pubSubKey)
			}

			pubSubConfig[topicNameKey] = topicName
			pubSubConfig[subscriptionIDNameKey] = subscriptionID
			delete(pubSubConfig, projectIDKey)
			delete(pubSubConfig, topicIDKey)
			delete(pubSubConfig, subscriptionIDKey)

			mapExtraConfig[auditLogsConfigKey].(map[string]interface{})[pubSubKey] = pubSubConfig

		}
	} else {
		delete(mapExtraConfig, auditLogsConfigKey)
	}
	tflog.Debug(ctx, fmt.Sprintf("mapExtraConfig: %s", mapExtraConfig))

	extraConfig, err := json.Marshal(mapExtraConfig)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("extra_config", string(extraConfig))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func addFieldError(diags diag.Diagnostics, fieldName, keyName string) diag.Diagnostics {
	return append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "An issue was encountered while processing the `extraConfig` field.",
		Detail:   fmt.Sprintf("missing or invalid %s field in %s", fieldName, keyName),
	})
}
