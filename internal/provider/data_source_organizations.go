package provider

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func dataSourceWizOrganizations() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details for Wiz organizations.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the search.  This is a sha1 hash of the search string. Changing the search string on this data source will result in a new data source state entry",
			},
			"organizations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Wiz schema identifier",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization external identifier",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization name",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization path",
						},
						"cloud_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud Provider",
						},
					},
				},
			},
			"search": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization search string. Used to search all organization attributes.",
			},
			"first": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "How many matches to return.",
			},
		},
		ReadContext: dataSourceWizOrganizationsRead,
	}
}

// ReadCloudOrganizations struct
type ReadCloudOrganizations struct {
	CloudOrganizations wiz.CloudOrganizationConnection `json:"cloudOrganizations,omitempty"`
}

func dataSourceWizOrganizationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizOrganizationsRead called...")

	// define the graphql query
	query := `query CloudOrganizations (
	    $filterBy: CloudOrganizationFilters
	    $first: Int
	){
	    cloudOrganizations(
	        filterBy: $filterBy
	        first: $first
	    ) {
	        pageInfo {
	            endCursor
	            hasNextPage
	        }
	        nodes {
	            id
	            externalId
	            name
	            cloudProvider
	            path
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.First = d.Get("first").(int)
	tflog.Debug(ctx, fmt.Sprintf("search strings (%T) %s", d.Get("search"), d.Get("search").(string)))
	filterBy := &wiz.CloudOrganizationFilters{}
	filterBy.Search = append(filterBy.Search, d.Get("search").(string))
	vars.FilterBy = filterBy

	// process the request
	data := &ReadCloudOrganizations{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "organizations", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// compute the id for this resource
	h := sha1.New()
	h.Write([]byte(d.Get("search").(string)))
	hashID := hex.EncodeToString(h.Sum(nil))

	// Set the id
	d.SetId(hashID)

	// Set the resource parameters
	tflog.Debug(ctx, fmt.Sprintf("Organization Read Search: %s", d.Get("search")))
	err := d.Set("search", d.Get("search").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("first", d.Get("first"))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	organizations := flattenOrganizations(ctx, &data.CloudOrganizations.Nodes)
	if err := d.Set("organizations", organizations); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func flattenOrganizations(ctx context.Context, nodes *[]wiz.CloudOrganization) []interface{} {
	tflog.Info(ctx, "flattenOrganizations called...")
	tflog.Debug(ctx, fmt.Sprintf("Nodes: %s", utils.PrettyPrint(nodes)))

	// walk the slices and construct the map
	var output = make([]interface{}, 0, 0)
	for _, b := range *nodes {
		OrganizationMap := make(map[string]interface{})
		OrganizationMap["id"] = b.ID
		OrganizationMap["external_id"] = b.ExternalID
		OrganizationMap["name"] = b.Name
		OrganizationMap["cloud_provider"] = b.CloudProvider
		OrganizationMap["path"] = b.Path
		output = append(output, OrganizationMap)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenOrganizations output: %s", utils.PrettyPrint(output)))

	return output
}
