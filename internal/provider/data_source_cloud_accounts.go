package provider

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
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
				Description: "How many results to return, maximum is `500` is per page.",
			},
			"max_pages": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "How many pages to return. 0 means all pages.",
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
				Description: "Free text search on cloud account name or tags or external-id. Specify list of empty string to return all cloud accounts.",
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
						wiz.CloudProvider,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.CloudProvider,
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
						wiz.CloudAccountStatus,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.CloudAccountStatus,
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
	CloudAccounts wiz.CloudAccountConnection `json:"cloudAccounts"`
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
	maxPages, b := d.GetOk("max_pages")
	if b {
		identifier.WriteString(utils.PrettyPrint(maxPages))
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
	filterBy := &wiz.CloudAccountFilters{}
	a, b = d.GetOk("ids")
	if b {
		filterBy.ID = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("search")
	if b {
		filterBy.Search = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("project_id")
	if b {
		filterBy.ProjectID = a.(string)
	}
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
	a, b = d.GetOk("assigned_to_project")
	if b {
		filterBy.AssignedToProject = utils.ConvertBoolToPointer(a.(bool))
	}
	a, b = d.GetOk("has_multiple_connector_sources")
	if b {
		filterBy.HasMultipleConnectorSources = utils.ConvertBoolToPointer(a.(bool))
	}
	vars.FilterBy = filterBy

	// process the request
	data := &ReadCloudAccounts{}
	requestDiags, allData := client.ProcessPagedRequest(ctx, m, vars, data, query, "cloud_accounts", "read", maxPages.(int))

	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	cloudAccounts := flattenCloudAccounts(ctx, allData)
	if err := d.Set("cloud_accounts", cloudAccounts); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func flattenCloudAccounts(ctx context.Context, cloudAccounts []interface{}) []interface{} {
	tflog.Info(ctx, "flattenCloudAccounts called...")
	tflog.Debug(ctx, fmt.Sprintf("cloudAccounts: %s", utils.PrettyPrint(cloudAccounts)))

	// walk the slice and construct the map
	var output = make([]interface{}, 0)
	for _, cldacc := range cloudAccounts {
		ca := cldacc.(*ReadCloudAccounts)
		for _, c := range ca.CloudAccounts.Nodes {
			tflog.Debug(ctx, fmt.Sprintf("c: %T %s", c, utils.PrettyPrint(c)))
			accountMap := make(map[string]interface{})
			accountMap["id"] = c.ID
			accountMap["external_id"] = c.ExternalID
			accountMap["name"] = c.Name
			accountMap["cloud_provider"] = c.CloudProvider
			accountMap["status"] = c.Status
			accountMap["linked_project_ids"] = flattenProjectIDs(ctx, &c.LinkedProjects)
			accountMap["source_connector_ids"] = flattenSourceConnectorIDs(ctx, &c.SourceConnectors)
			output = append(output, accountMap)

		}
		// sort the return slice to avoid unwanted diffs
		sort.Slice(output, func(i, j int) bool {
			return output[i].(map[string]interface{})["id"].(string) < output[j].(map[string]interface{})["id"].(string)
		})
		tflog.Debug(ctx, fmt.Sprintf("flattenCloudAccounts output: %s", utils.PrettyPrint(output)))

	}
	return output
}

func flattenProjectIDs(ctx context.Context, projects *[]*wiz.Project) []interface{} {
	tflog.Info(ctx, "flattenProjectIDs called...")
	tflog.Debug(ctx, fmt.Sprintf("Projects: %s", utils.PrettyPrint(projects)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0)
	for _, b := range *projects {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	// sort the return slice to avoid unwanted diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].(string) < output[j].(string)
	})

	tflog.Debug(ctx, fmt.Sprintf("flattenProjectIDs output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenSourceConnectorIDs(ctx context.Context, connectors *[]wiz.Connector) []interface{} {
	tflog.Info(ctx, "flattenSourceConnectorIDs called...")
	tflog.Debug(ctx, fmt.Sprintf("Projects: %s", utils.PrettyPrint(connectors)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0)
	for _, b := range *connectors {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	// sort the return slice to avoid unwanted diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].(string) < output[j].(string)
	})

	tflog.Debug(ctx, fmt.Sprintf("flattenSourceConnectorIDs output: %s", utils.PrettyPrint(output)))

	return output
}
