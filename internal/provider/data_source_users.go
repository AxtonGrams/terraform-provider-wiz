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

// ReadUsers struct
type ReadUsers struct {
	Users wiz.UserConnection `json:"users"`
}

func dataSourceWizUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details for Wiz users.",
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
				Description: "How many matches to return, maximum is `100` is per page.",
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
				Description: "Free text search. Specify empty string to return all users.",
			},
			"authentication_source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MODERN",
				Description: fmt.Sprintf(
					"Authentication Source.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.AuthenticationSource,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.AuthenticationSource,
						false,
					),
				),
			},
			"users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The returned wiz users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Internal Wiz ID.",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User email address.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User email name.",
						},
						"is_suspended": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If user is suspended.",
						},
						"identity_provider_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identity Provider Type.",
						},
						"identity_provider": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The identity provider details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The identity provider name.",
									},
								},
							},
						},
						"effective_role": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The effective role details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The role internal identifier.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The role name.",
									},
									"scopes": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Permission Scopes.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"roles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of roles to filter by.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: dataSourceWizUsersRead,
	}
}

func dataSourceWizUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "dataSourceWizUsersRead called...")

	// //generate the id for this resource
	// //id must be deterministic, so the id is based on a hash of the search parameters
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
	a, b = d.GetOk("roles")
	if b {
		identifier.WriteString(utils.PrettyPrint(a))
	}
	a, b = d.GetOk("authentication_source")
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
	query := `query users(
	  $first: Int
	  $filterBy: UserFilters
	  $after: String
	){
	  users(
	    first: $first,
	    filterBy: $filterBy,
	    after: $after
	  ) {
	      nodes {
	        id
	        name
	        email	      
	        identityProvider {
	          name
	        }
	        identityProviderType
	        effectiveRole {
	          id
	          name
	          scopes
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
	err = d.Set("roles", d.Get("roles").([]interface{}))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("authentication_source", d.Get("authentication_source").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.First = d.Get("first").(int)
	filterBy := &wiz.UserFilters{}
	a, b = d.GetOk("search")
	if b {
		filterBy.Search = a.(string)
	}
	a, b = d.GetOk("authentication_source")
	if b {
		filterBy.AuthenticationSource = a.(string)
	}
	a, b = d.GetOk("roles")
	if b {
		filterBy.Roles = utils.ConvertListToString(a.([]interface{}))
	}

	vars.FilterBy = filterBy

	// process the request
	data := &ReadUsers{}
	requestDiags, allData := client.ProcessPagedRequest(ctx, m, vars, data, query, "users", "read", maxPages.(int))
	tflog.Debug(ctx, fmt.Sprintf("allData: %s", utils.PrettyPrint(allData)))

	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	users := flattenUsers(ctx, allData)
	if err := d.Set("users", users); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func flattenUsers(ctx context.Context, users []interface{}) interface{} {
	tflog.Info(ctx, "flattenUsers called...")
	tflog.Debug(ctx, fmt.Sprintf("Users: %s", utils.PrettyPrint(users)))

	// walk the slice and construct the list
	var output = make([]interface{}, 0)
	for _, u := range users {
		readUsers := u.(*ReadUsers)
		for _, c := range readUsers.Users.Nodes {
			tflog.Debug(ctx, fmt.Sprintf("c: %T %s", c, utils.PrettyPrint(c)))
			userMap := make(map[string]interface{})
			userMap["id"] = c.ID
			userMap["email"] = c.Email
			userMap["name"] = c.Name
			userMap["is_suspended"] = c.IsSuspended
			userMap["identity_provider_type"] = c.IdentityProviderType

			idpMap := make(map[string]interface{})
			idpMap["name"] = c.IdentityProvider.Name
			userMap["identity_provider"] = []interface{}{idpMap}

			roleMap := make(map[string]interface{})
			roleMap["id"] = c.EffectiveRole.ID
			roleMap["name"] = c.EffectiveRole.Name
			roleMap["scopes"] = c.EffectiveRole.Scopes
			userMap["effective_role"] = []interface{}{roleMap}

			output = append(output, userMap)
		}

	}
	return output
}
