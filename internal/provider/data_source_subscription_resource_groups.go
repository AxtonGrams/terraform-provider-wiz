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
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

// ReadSubscriptionResourceGroups struct
type ReadSubscriptionResourceGroups struct {
	SubscriptionResourceGroups wiz.GraphSearchResultConnection `json:"graphsearch"`
}

func dataSourceWizSubscriptionResourceGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches the resource groups that are part of the subscription.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal identifier for the data.",
			},
			"first": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     50,
				Description: "How many matches to return.",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Wiz subscription ID to search by.",
			},
			"relationship_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "CONTAINS",
				Description: fmt.Sprintf("Relationship type, will default to `CONTAINS` if not specified.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.GraphRelationshipType,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.GraphRelationshipType,
						false,
					),
				),
			},
			"resource_groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The returned subscription resource groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Internal Wiz ID of Resource Group.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the Resource Group.",
						},
					},
				},
			},
		},
		ReadContext: dataSourceWizSubscriptionResourceGroupsRead,
	}
}

func dataSourceWizSubscriptionResourceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizSubscriptionResourceGroupsRead called...")
	var identifier bytes.Buffer

	a, b := d.GetOk("first")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("id")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("subscription_id")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("relationship_type")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	h := sha1.New()
	h.Write([]byte(identifier.String()))
	hashID := hex.EncodeToString(h.Sum(nil))

	// Set the id
	d.SetId(hashID)

	// We set quick to false to ensure ordering, no noticable performance trade-off detected
	query := `query ResourceGroupQuery($query: GraphEntityQueryInput, $quick: Boolean = false, $first:Int) {
	  graphSearch(first: $first, query: $query, quick: $quick) {
	    nodes {
	      entities {
	        id
	        name
	      }
	    }
	   }
	  }`

	// set the resource parameters
	err := d.Set("id", d.Get("id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("subscription_id", d.Get("subscription_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("relationship_type", d.Get("relationship_type").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.First = d.Get("first").(int)

	// declare main graph query
	resourceGroupQuery := &wiz.GraphEntityQueryInput{}
	resourceGroupQuery.Type = []string{"RESOURCE_GROUP"}

	// set the relationship type
	relationshipDirectedType := &wiz.GraphDirectedRelationshipTypeInput{}
	a, b = d.GetOk("relationship_type")
	if b {
		relationshipDirectedType.Type = a.(string)
	}
	// reverse needs be to be true to fetch resource group subscription relationships from the edges
	reverse := true
	relationshipDirectedType.Reverse = &reverse

	var directedRelationshipQueryInput = []wiz.GraphDirectedRelationshipTypeInput{}
	directedRelationshipQueryInput = append(directedRelationshipQueryInput, *relationshipDirectedType)

	entityInput := &wiz.GraphEntityQueryInput{}
	entityInput.Type = []string{"SUBSCRIPTION"}

	// set the where predicate for the query to the subscription id
	a, b = d.GetOk("subscription_id")
	if b {
		wherePredicate := map[string]interface{}{
			"_vertexID": map[string]interface{}{
				"EQUALS": a.(string),
			},
		}
		entityInput.Where = wherePredicate
	}

	var relationshipQueryInput = &wiz.GraphRelationshipQueryInput{
		With: *entityInput,
		Type: directedRelationshipQueryInput,
	}
	var relationships = []*wiz.GraphRelationshipQueryInput{}
	relationships = append(relationships, relationshipQueryInput)
	resourceGroupQuery.Relationships = relationships

	vars.Query = resourceGroupQuery

	// process the request
	data := &ReadSubscriptionResourceGroups{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "subscriptionResourceGroups", "read")

	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	resourceGroups := flattenResourceGroups(ctx, &data.SubscriptionResourceGroups.Nodes)

	if err := d.Set("resource_groups", resourceGroups); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags

}

func flattenResourceGroups(ctx context.Context, resgroups *[]*wiz.GraphSearchResult) []interface{} {
	tflog.Info(ctx, "flattenResourceGroups called...")
	tflog.Debug(ctx, fmt.Sprintf("resourceGroups: %s", utils.PrettyPrint(resgroups)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *resgroups {
		for _, e := range b.Entities {
			resourceGroups := make(map[string]interface{})
			resourceGroups["id"] = e.ID
			resourceGroups["name"] = e.Name
			tflog.Debug(ctx, fmt.Sprintf("resourceGroups output: id: %s, name: %s", utils.PrettyPrint(e.ID), utils.PrettyPrint(e.Name)))
			output = append(output, resourceGroups)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenResourceGroups output: %s", utils.PrettyPrint(output)))
	return output
}
