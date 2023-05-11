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

// ReadKubernetesClusters struct
type ReadKubernetesClusters struct {
	KubernetesClusters wiz.KubernetesClusterConnection `json:"kubernetesClusters"`
}

func dataSourceWizKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details for Kubernetes clusters.",
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
				Description: "How many matches to return, maximum is `500` per page.",
			},
			"max_pages": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "How many pages to return. 0 means all pages.",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Free text search. Specify empty string to return all kubernetes clusters",
			},
			"external_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The ID(s) to search by. i.e `Azure Subscription ID` or `AWS account number`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kind": {
				Type:     schema.TypeList,
				Optional: true,
				Description: fmt.Sprintf("Query Kubernetes Cluster of specific kind(s) or cloud provider(s).\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.KubernetesClusterKind,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							wiz.KubernetesClusterKind,
							false,
						),
					),
				},
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
			"kubernetes_clusters": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The returned kubernetes clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Internal Wiz ID of Kubernetes Cluster.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the Kubernetes Cluster.",
						},
						"cloud_account": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The cloud account details for the kubernetes cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Internal Wiz ID of Cloud Account.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the cloud account.",
									},
									"cloud_provider": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The cloud provider of the cloud account.",
									},
									"external_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The external ID of the cloud account.",
									},
								},
							},
						},
					},
				},
			},
		},
		ReadContext: dataSourceWizKubernetesClustersRead,
	}
}

func dataSourceWizKubernetesClustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizKubernetesClustersRead called...")
	var identifier bytes.Buffer

	a, b := d.GetOk("first")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("id")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("search")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("kind")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("external_ids")
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

	query := `query ClustersPage($filterBy: KubernetesClusterFilters, $first: Int, $after: String) {
		kubernetesClusters(filterBy: $filterBy, first: $first, after: $after) {
		  nodes {
			id
			externalId
			name
			kind
			cloudAccount {
			  id
			  name
			  cloudProvider
			  externalId
			}
			projects {
			  id
			  name
			  slug
			  riskProfile {
				businessImpact
			  }
			}
		  }
		  pageInfo {
			endCursor
			hasNextPage
		  }
		  totalCount
		}
	  }`

	// set the resource parameters
	err := d.Set("search", d.Get("search").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("id", d.Get("id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("kind", d.Get("kind").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("external_ids", d.Get("external_ids").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.First = d.Get("first").(int)
	filterBy := &wiz.KubernetesClusterFilters{}
	a, b = d.GetOk("search")
	if b {
		filterBy.Search = a.(string)
	}
	a, b = d.GetOk("kind")
	if b {
		filterBy.Kind = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("external_ids")
	if b {
		filterBy.CloudAccount = utils.ConvertListToString(a.([]interface{}))
	}

	vars.FilterBy = filterBy

	// process the request
	data := &ReadKubernetesClusters{}
	requestDiags, allData := client.ProcessPagedRequest(ctx, m, vars, data, query, "kubernetesClusters", "read", maxPages.(int))
	tflog.Debug(ctx, fmt.Sprintf("allData: %s", utils.PrettyPrint(allData)))

	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	clusters := flattenClusters(ctx, allData)

	if err := d.Set("kubernetes_clusters", clusters); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags

}

func flattenClusters(ctx context.Context, clusters []interface{}) []interface{} {
	tflog.Info(ctx, "flattenClusters called...")
	tflog.Debug(ctx, fmt.Sprintf("Clusters: %s", utils.PrettyPrint(clusters)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0)
	for _, b := range clusters {
		readClusters := b.(*ReadKubernetesClusters)
		for _, cluster := range readClusters.KubernetesClusters.Nodes {
			tflog.Debug(ctx, fmt.Sprintf("cluster: %s", utils.PrettyPrint(cluster)))
			rootMap := make(map[string]interface{})
			rootMap["id"] = cluster.ID
			rootMap["name"] = cluster.Name

			clusterMap := make(map[string]interface{})
			clusterMap["cloud_provider"] = cluster.CloudAccount.CloudProvider
			clusterMap["external_id"] = cluster.CloudAccount.ExternalID
			clusterMap["id"] = cluster.CloudAccount.ID
			clusterMap["name"] = cluster.CloudAccount.Name

			cloudAccountMap := make([]interface{}, 0)
			cloudAccountMap = append(cloudAccountMap, clusterMap)
			rootMap["cloud_account"] = cloudAccountMap

			output = append(output, rootMap)

		}
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenClusters output: %s", utils.PrettyPrint(output)))
	return output
}
