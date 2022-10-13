package provider

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func dataSourceWizCloudAccounts() *schema.Resource {
	return &schema.Resource{
		Description: "Query cloud accounts (subscriptions).",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal identifier for the data.",
			},
			"first": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "How many results to return",
			},
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Get specific Cloud Accounts by their IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"search": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Free text search on cloud account name or tags or external-id.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query cloud accounts of a specific linked project, given its id.",
			},
			"cloud_provider": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("Query cloud accounts of specific cloud provider.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.CloudProvider,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							vendor.CloudProvider,
							false,
						),
					),
				},
			},
			"status": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("Query cloud accounts by status.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						vendor.CloudAccountStatus,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							vendor.CloudAccountStatus,
							false,
						),
					),
				},
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
			"has_multiple_connector_sources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "QueryQuery cloud accounts by project assignment state.",
			},
			"cloud_accounts": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The returned cloud accounts.",
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
						"source_connector_ids": {
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
						"linked_project_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Projects list this cloud account is assigned to/",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		ReadContext: dataSourceWizCloudAccountsRead,
	}
}

// ReadCloudAccounts struct
type ReadCloudAccounts struct {
	CloudAccounts vendor.CloudAccountConnection `json:"cloudAccounts"`
}

func dataSourceWizCloudAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizCloudAccountsRead called...")

	// generate the id for this resource
	// id must be deterministic, so the id is based on a hash of the search parameters
	var identifier bytes.Buffer

	a, b := d.GetOk("first")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("ids")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("search")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("project_id")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("cloud_provider")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("status")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("connector_id")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("connector_issue_id")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("assigned_to_project")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("has_multiple_connector_sources")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}

	h := sha1.New()
	h.Write([]byte(identifier.String()))
	hashID := hex.EncodeToString(h.Sum(nil))

	// Set the id
	d.SetId(hashID)

	// define the graphql query
	query := `query cloudAccounts(
	  $filterBy: CloudAccountFilters
	  $first: Int
	  $after: String
	) {
	  cloudAccounts(
	    filterBy: $filterBy
	    first: $first
	    after: $after
	  ) {
	    nodes {
	      id
	      externalId
	      name
	      cloudProvider
	      sourceConnectors {
	        id
	      }
	      status
	      linkedProjects {
	        id
	      }
	    }
	    pageInfo {
	      hasNextPage
	      endCursor
	    }
	    totalCount
	  }
	}`

	// set the resource parameters
	err := d.Set("first", d.Get("first").(int))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("ids", d.Get("ids").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("search", d.Get("search").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", d.Get("project_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("cloud_provider", d.Get("cloud_provider").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("status", d.Get("status").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("connector_id", d.Get("connector_id").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("connector_issue_id", d.Get("connector_issue_id").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	a, b = d.GetOk("assigned_to_project")
	if b {
		err = d.Set("assigned_to_project", a.(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	a, b = d.GetOk("has_multiple_connector_sources")
	if b {
		err = d.Set("has_multiple_connector_sources", d.Get("has_multiple_connector_sources").(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.First = d.Get("first").(int)
	filterBy := &vendor.CloudAccountFilters{}
	a, b = d.GetOk("ids")
	if b {
		filterBy.ID = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("search")
	if b {
		filterBy.Search = utils.ConvertListToString(a.([]interface{}))
	}
	filterBy.ProjectID = d.Get("project_id").(string)
	a, b = d.GetOk("cloud_provider")
	if b {
		filterBy.CloudProvider = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("status")
	if b {
		filterBy.Status = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("connector_id")
	if b {
		filterBy.ConnectorID = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("connector_issue_id")
	if b {
		filterBy.ConnectorIssueID = utils.ConvertListToString(a.([]interface{}))
	}
	filterBy.AssignedToProject = utils.ConvertBoolToPointer(d.Get("assigned_to_project").(bool))
	filterBy.HasMultipleConnectorSources = utils.ConvertBoolToPointer(d.Get("has_multiple_connector_sources").(bool))
	vars.FilterBy = filterBy

	// process the request
	data := &ReadCloudAccounts{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cloud_accounts", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	cloudAccounts := flattenCloudAccounts(ctx, &data.CloudAccounts.Nodes)
	if err := d.Set("cloud_accounts", cloudAccounts); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	tflog.Debug(ctx, "Finished")

	return diags
}

func flattenCloudAccounts(ctx context.Context, nodes *[]*vendor.CloudAccount) []interface{} {
	tflog.Info(ctx, "flattenCloudAccounts called...")
	tflog.Debug(ctx, fmt.Sprintf("Nodes: %s", utils.PrettyPrint(nodes)))

	// walk the slice and construct the map
	var output = make([]interface{}, 0, 0)
	for _, b := range *nodes {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		accountMap := make(map[string]interface{})
		accountMap["id"] = b.ID
		accountMap["external_id"] = b.ExternalID
		accountMap["name"] = b.Name
		accountMap["cloud_provider"] = b.CloudProvider
		accountMap["status"] = b.Status
		accountMap["linked_project_ids"] = flattenProjectIDs(ctx, &b.LinkedProjects)
		accountMap["source_connector_ids"] = flattenSourceConnectorIDs(ctx, &b.SourceConnectors)
		output = append(output, accountMap)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenCloudAccounts output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenProjectIDs(ctx context.Context, projects *[]*vendor.Project) []interface{} {
	tflog.Info(ctx, "flattenProjectIDs called...")
	tflog.Debug(ctx, fmt.Sprintf("Projects: %s", utils.PrettyPrint(projects)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *projects {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenProjectIDs output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenSourceConnectorIDs(ctx context.Context, connectors *[]vendor.Connector) []interface{} {
	tflog.Info(ctx, "flattenSourceConnectorIDs called...")
	tflog.Debug(ctx, fmt.Sprintf("Projects: %s", utils.PrettyPrint(connectors)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *connectors {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenSourceConnectorIDs output: %s", utils.PrettyPrint(output)))

	return output
}
