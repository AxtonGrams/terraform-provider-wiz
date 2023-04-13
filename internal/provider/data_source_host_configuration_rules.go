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
	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func dataSourceWizHostConfigurationRules() *schema.Resource {
	return &schema.Resource{
		Description: "Query cloud configuration rules.",
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
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Free text search on id, name, externalId.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Host Configuration Rule enabled status.",
			},
			"framework_category": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search rules by any of securityFramework | securitySubCategory | securityCategory.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_platform": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search by target platforms.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_configuration_rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The returned cloud configuration rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Wiz UUID.",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An external id for the rule.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the rule.",
						},
						"short_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A short name that identifies the rule.",
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Rule enabled status.",
						},
						"security_sub_category_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"builtin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indication whether the rule is built-in or custom.",
						},
						"direct_oval": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Direct OVAL definition assessed on hosts during disk scanning.",
						},
						"target_platform_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The platforms the rule is targeting. e.g Ubuntu, RedHat, NGINX.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		ReadContext: dataSourceWizHostConfigurationRuleRead,
	}
}

// ReadHostConfigurationRules struct
type ReadHostConfigurationRules struct {
	HostConfigurationRules wiz.HostConfigurationRuleConnection `json:"hostConfigurationRules"`
}

func dataSourceWizHostConfigurationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizHostConfigurationRuleRead called...")

	// generate the id for this resource
	// id must be deterministic, so the id is based on a hash of the search parameters
	var identifier bytes.Buffer

	a, b := d.GetOk("first")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("search")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("enabled")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("framework_category")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("target_platform")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}

	h := sha1.New()
	h.Write([]byte(identifier.String()))
	hashID := hex.EncodeToString(h.Sum(nil))

	// Set the id
	d.SetId(hashID)

	// define the graphql query
	query := `query hostConfigurationRules(
	  $filterBy: HostConfigurationRuleFilters
	  $first: Int
	  $after: String
	  $orderBy: HostConfigurationRuleOrder
	) {
	  hostConfigurationRules(
	    filterBy: $filterBy
	    first: $first
	    after: $after
	    orderBy: $orderBy
	  ) {
	    nodes {
	      id
	      externalId
	      name
	      shortName
	      description
	      enabled
	      securitySubCategories {
	        id
	      }
	      builtin
	      directOVAL
	      targetPlatforms {
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
	err = d.Set("search", d.Get("search").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	a, b = d.GetOk("enabled")
	if b {
		err = d.Set("enabled", a.(bool))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	err = d.Set("framework_category", d.Get("framework_category").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("target_platform", d.Get("target_platform").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.First = d.Get("first").(int)
	filterBy := &wiz.HostConfigurationRuleFilters{}
	a, b = d.GetOk("search")
	if b {
		filterBy.Search = a.(string)
	}
	a, b = d.GetOk("enabled")
	if b {
		filterBy.Enabled = utils.ConvertBoolToPointer(a.(bool))
	}
	a, b = d.GetOk("framework_category")
	if b {
		filterBy.FrameworkCategory = utils.ConvertListToString(a.([]interface{}))
	}
	a, b = d.GetOk("target_platform")
	if b {
		filterBy.TargetPlatforms = utils.ConvertListToString(a.([]interface{}))
	}

	vars.FilterBy = filterBy

	// process the request
	data := &ReadHostConfigurationRules{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "host_config_rules", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	hostConfigurationRules := flattenHostConfigurationRules(ctx, &data.HostConfigurationRules.Nodes)
	if err := d.Set("host_configuration_rules", hostConfigurationRules); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func flattenHostConfigurationRules(ctx context.Context, nodes *[]*wiz.HostConfigurationRule) []interface{} {
	tflog.Info(ctx, "flattenHostConfigurationRules called...")
	tflog.Debug(ctx, fmt.Sprintf("HostConfigurationRules: %s", utils.PrettyPrint(nodes)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range *nodes {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		ruleMap := make(map[string]interface{})
		ruleMap["id"] = b.ID
		ruleMap["name"] = b.Name
		ruleMap["short_name"] = b.ShortName
		ruleMap["builtin"] = b.Builtin
		ruleMap["description"] = b.Description
		ruleMap["direct_oval"] = b.DirectOVAL
		ruleMap["external_id"] = b.ExternalID
		ruleMap["security_sub_category_ids"] = flattenSecuritySubCategoryIDs(ctx, &b.SecuritySubCategories)
		ruleMap["target_platform_ids"] = flattenTargetPlatformIDs(ctx, b.TargetPlatforms)

		output = append(output, ruleMap)
	}

	tflog.Debug(ctx, fmt.Sprintf("flattenCloudConfigurationRules output: %s", utils.PrettyPrint(output)))

	return output
}

func flattenTargetPlatformIDs(ctx context.Context, plats []wiz.Technology) []interface{} {
	tflog.Info(ctx, "flattenTargetPlatformIDs called...")
	tflog.Debug(ctx, fmt.Sprintf("TargetPlatforms: %s", utils.PrettyPrint(plats)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0, 0)
	for _, b := range plats {
		tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}

	// sort the return slice to avoid unwanted diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].(string) < output[j].(string)
	})

	tflog.Debug(ctx, fmt.Sprintf("flattenTargetPlatformIDs output: %s", utils.PrettyPrint(output)))

	return output
}
