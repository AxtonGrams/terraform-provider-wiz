package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizSAMLIdP() *schema.Resource {
	return &schema.Resource{
		Description: "Configure SAML Providers and associated resources (group mappings).",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Internal identifier for the Saml Provider",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "IdP name to display in Wiz.",
				Required:    true,
			},
			"issuer_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If undefined, this will default to the login_url value. Set to the same value as login_url if unsure what value to use.",
			},
			"login_url": {
				Type:        schema.TypeString,
				Description: "IdP Login URL",
				Required:    true,
			},
			"logout_url": {
				Type:        schema.TypeString,
				Description: "IdP Logout URL",
				Optional:    true,
			},
			"use_provider_managed_roles": {
				Type:        schema.TypeBool,
				Description: "When set to true, roles will be provided by the SSO provider. Manage the roles via Wiz portal otherwise.",
				Optional:    true,
				Default:     false,
			},
			"allow_manual_role_override": {
				Type:        schema.TypeBool,
				Description: "When set to true, allow overriding the mapped SSO role for specific users. Must be set `true` if `use_provided_roles` is false.",
				Optional:    true,
				Default:     true,
				RequiredWith: []string{
					"use_provider_managed_roles",
				},
			},
			"certificate": {
				Type:        schema.TypeString,
				Description: "PEM certificate from IdP",
				Required:    true,
			},
			"domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of domains the IdP handles.",
				Deprecated:  "This field is no longer supported by Wiz. If defined, this will result in change detection on every run.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_mapping": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Group mappings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Description: "Description",
							Optional:    true,
						},
						"provider_group_id": {
							Type:        schema.TypeString,
							Description: "Provider group ID",
							Required:    true,
						},
						"role": {
							Type:        schema.TypeString,
							Description: "Wiz Role name",
							Required:    true,
						},
						"projects": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Project mapping",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"merge_groups_mapping_by_role": {
				Type:        schema.TypeBool,
				Description: "Manage group mapping by role?",
				Optional:    true,
			},
		},
		CreateContext: resourceWizSAMLIdPCreate,
		ReadContext:   resourceWizSAMLIdPRead,
		UpdateContext: resourceWizSAMLIdPUpdate,
		DeleteContext: resourceWizSAMLIdPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func getGroupMappingVar(ctx context.Context, d *schema.ResourceData) []*wiz.SAMLGroupMappingCreateInput {
	groupMapping := d.Get("group_mapping").(*schema.Set).List()
	var myGroupMappings []*wiz.SAMLGroupMappingCreateInput
	for _, a := range groupMapping {
		tflog.Debug(ctx, fmt.Sprintf("groupMapping: %t %s", a, utils.PrettyPrint(a)))
		localGroupMapping := &wiz.SAMLGroupMappingCreateInput{}
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "description":
				localGroupMapping.Description = c.(string)
			case "role":
				localGroupMapping.Role = c.(string)
			case "provider_group_id":
				localGroupMapping.ProviderGroupID = c.(string)
			case "projects":
				for _, f := range c.([]interface{}) {
					tflog.Trace(ctx, fmt.Sprintf("f: %t %s", f, f))
					localGroupMapping.Projects = append(localGroupMapping.Projects, f.(string))
				}
			}
		}
		myGroupMappings = append(myGroupMappings, localGroupMapping)
	}
	tflog.Debug(ctx, fmt.Sprintf("myGroupMappings: %s", utils.PrettyPrint(myGroupMappings)))
	return myGroupMappings
}

// CreateSAMLIdentityProvider struct
type CreateSAMLIdentityProvider struct {
	CreateSAMLIdentityProvider wiz.CreateSAMLIdentityProviderPayload `json:"createSAMLIdentityProvider"`
}

func resourceWizSAMLIdPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLIdPCreate called...")

	// define the graphql query
	query := `mutation CreateSAMLIdentityProvider ($input: CreateSAMLIdentityProviderInput!) {
	  createSAMLIdentityProvider(
	    input: $input
	  ) {
	    samlIdentityProvider {
	      id
	    }
	  }
	}`

	// populate the graphql variables
	vars := &wiz.CreateSAMLIdentityProviderInput{}
	vars.Name = d.Get("name").(string)
	vars.IssuerURL = d.Get("issuer_url").(string)
	vars.LoginURL = d.Get("login_url").(string)
	vars.LogoutURL = d.Get("logout_url").(string)
	vars.UseProviderManagedRoles = d.Get("use_provider_managed_roles").(bool)
	vars.AllowManualRoleOverride = utils.ConvertBoolToPointer(d.Get("allow_manual_role_override").(bool))
	vars.Certificate = d.Get("certificate").(string)
	vars.MergeGroupsMappingByRole = utils.ConvertBoolToPointer(d.Get("merge_groups_mapping_by_role").(bool))
	vars.Domains = utils.ConvertListToString(d.Get("domains").([]interface{}))
	vars.GroupMapping = getGroupMappingVar(ctx, d)

	// process the request
	data := &CreateSAMLIdentityProvider{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_idp", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateSAMLIdentityProvider.SAMLIdentityProvider.ID)

	return resourceWizSAMLIdPRead(ctx, d, m)
}

func flattenGroupMapping(ctx context.Context, samlGroupMapping []*wiz.SAMLGroupMapping) []interface{} {
	tflog.Info(ctx, "flattenGroupMapping called...")
	var output = make([]interface{}, 0)
	for _, b := range samlGroupMapping {
		tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		var mapping = make(map[string]interface{})
		var projects = make([]interface{}, 0)
		for _, d := range b.Projects {
			tflog.Trace(ctx, fmt.Sprintf("d: %T %s", d, utils.PrettyPrint(d)))
			projects = append(projects, d.ID)
		}
		mapping["projects"] = projects
		mapping["description"] = b.Description
		mapping["provider_group_id"] = b.ProviderGroupID
		mapping["role"] = b.Role.ID
		tflog.Trace(ctx, fmt.Sprintf("projects: %s", projects))
		tflog.Trace(ctx, fmt.Sprintf("mapping: %s", utils.PrettyPrint(mapping)))
		output = append(output, mapping)
	}
	tflog.Debug(ctx, fmt.Sprintf("output: %s", utils.PrettyPrint(output)))
	return output
}

// ReadSAMLIdentityProviderPayload struct -- updates
type ReadSAMLIdentityProviderPayload struct {
	SAMLIdentityProvider wiz.SAMLIdentityProvider `json:"samlIdentityProvider"`
}

func resourceWizSAMLIdPRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLIdPRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query samlIdentityProvider ($id: ID!){
	    samlIdentityProvider (
	        id: $id
	    ) {
	        id
	        name
	        issuerURL
	        loginURL
	        logoutURL
	        useProviderManagedRoles
	        allowManualRoleOverride
	        certificate
	        domains
	        mergeGroupsMappingByRole
	        groupMapping {
	            description
	            providerGroupId
	            role {
	                id
	            }
	            projects {
	                id
	            }
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: oops! an internal error has occurred. for reference purposes, this is your request id
	data := &ReadSAMLIdentityProviderPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_idp", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.SAMLIdentityProvider.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.SAMLIdentityProvider.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("issuer_url", data.SAMLIdentityProvider.IssuerURL)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("login_url", data.SAMLIdentityProvider.LoginURL)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("logout_url", data.SAMLIdentityProvider.LogoutURL)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("certificate", data.SAMLIdentityProvider.Certificate)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("use_provider_managed_roles", data.SAMLIdentityProvider.UseProviderManagedRoles)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("merge_groups_mapping_by_role", data.SAMLIdentityProvider.MergeGroupsMappingByRole)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("allow_manual_role_override", data.SAMLIdentityProvider.AllowManualRoleOverride)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("domains", data.SAMLIdentityProvider.Domains)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	groupMappings := flattenGroupMapping(ctx, data.SAMLIdentityProvider.GroupMapping)
	tflog.Debug(ctx, fmt.Sprintf("groupMappings: %s", utils.PrettyPrint(groupMappings)))
	if err := d.Set("group_mapping", groupMappings); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateSAMLIdentityProvider struct
type UpdateSAMLIdentityProvider struct {
	UpdateSAMLIdentityProvider wiz.UpdateSAMLIdentityProviderPayload `json:"updateSAMLIdentityProvider"`
}

func resourceWizSAMLIdPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLIdPUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateSAMLIdentityProvider($input: UpdateSAMLIdentityProviderInput!) {
	    updateSAMLIdentityProvider(input: $input) {
	        samlIdentityProvider {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateSAMLIdentityProviderInput{}
	vars.ID = d.Id()
	vars.Patch.Name = d.Get("name").(string)
	vars.Patch.IssuerURL = d.Get("issuer_url").(string)
	vars.Patch.LoginURL = d.Get("login_url").(string)
	vars.Patch.LogoutURL = d.Get("logout_url").(string)
	vars.Patch.UseProviderManagedRoles = utils.ConvertBoolToPointer(d.Get("use_provider_managed_roles").(bool))
	vars.Patch.AllowManualRoleOverride = utils.ConvertBoolToPointer(d.Get("allow_manual_role_override").(bool))
	vars.Patch.Certificate = d.Get("certificate").(string)
	vars.Patch.MergeGroupsMappingByRole = utils.ConvertBoolToPointer(d.Get("merge_groups_mapping_by_role").(bool))
	// populate the group mapping
	mappings := d.Get("group_mapping").(*schema.Set).List()
	mappingUpdates := make([]wiz.SAMLGroupMappingUpdateInput, 0)
	for a, b := range mappings {
		var myMap = wiz.SAMLGroupMappingUpdateInput{}
		tflog.Trace(ctx, fmt.Sprintf("a:b: %d %s", a, b))

		for c, d := range b.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("c:d: %s %s", c, d))
			switch c {
			case "description":
				myMap.Description = d.(string)
			case "role":
				myMap.Role = d.(string)
			case "provider_group_id":
				myMap.ProviderGroupID = d.(string)
			case "projects":
				for _, f := range d.([]interface{}) {
					tflog.Trace(ctx, fmt.Sprintf("f: %t %s", f, f))
					myMap.Projects = append(myMap.Projects, f.(string))
				}
			}
		}
		mappingUpdates = append(mappingUpdates, myMap)
	}
	vars.Patch.GroupMapping = mappingUpdates

	// process the request
	data := &UpdateSAMLIdentityProvider{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_idp", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizSAMLIdPRead(ctx, d, m)
}

// DeleteSAMLIdentityProvider struct
type DeleteSAMLIdentityProvider struct {
	DeleteSAMLIdentityProvider wiz.DeleteSAMLIdentityProviderPayload `json:"deleteSAMLIdentityProvider"`
}

func resourceWizSAMLIdPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLIdPDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteSAMLIdentityProvider (
	    $input: DeleteSAMLIdentityProviderInput!
	) {
	    deleteSAMLIdentityProvider(
	        input: $input
	    ) {
	        _stub
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteSAMLIdentityProviderInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateSAMLIdentityProvider{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_idp", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
