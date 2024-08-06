package provider

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

// ReadSAMLGroupMappings represents the structure of a SAML group mappings read operation.
// It includes a SAMLGroupMappings object.
type ReadSAMLGroupMappings struct {
	SAMLGroupMappings SAMLGroupMappings `json:"samlIdentityProviderGroupMappings"`
}

// SAMLGroupMappings represents the structure of SAML group mappings.
// It includes PageInfo and a list of Nodes.
type SAMLGroupMappings struct {
	PageInfo wiz.PageInfo            `json:"pageInfo"`
	Nodes    []*wiz.SAMLGroupMapping `json:"nodes,omitempty"`
}

// SAMLGroupMappingsImport represents the structure of a SAML group mapping import.
// It includes the SAML IdP ID, provider group ID, project IDs, and role.
type SAMLGroupMappingsImport struct {
	SamlIdpID       string
	ProviderGroupID string
	ProjectIDs      []string
	Role            string
}

// UpdateSAMLGroupMappingPayload struct
type UpdateSAMLGroupMappingPayload struct {
	SAMLGroupMapping wiz.SAMLGroupMapping `json:"samlGroupMapping,omitempty"`
}

// DeleteSAMLGroupMappingInput struct
type DeleteSAMLGroupMappingInput struct {
	ID    string                 `json:"id"`
	Patch DeleteSAMLGroupMapping `json:"patch"`
}

// DeleteSAMLGroupMapping struct
type DeleteSAMLGroupMapping struct {
	Delete []string `json:"delete"`
}

func resourceWizSAMLGroupMapping() *schema.Resource {
	return &schema.Resource{
		Description: "Configure SAML Group Role Mapping. When using SSO to authenticate with Wiz, you can map group memberships in SAML assertions to Wiz roles across specific scopes.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique tf-internal identifier for the saml group mapping",
				Computed:    true,
			},
			"saml_idp_id": {
				Type:        schema.TypeString,
				Description: "Identifier for the Saml Provider",
				Required:    true,
				ForceNew:    true,
			},
			"provider_group_id": {
				Type:        schema.TypeString,
				Description: "Provider group ID",
				Required:    true,
				ForceNew:    true,
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
		CreateContext: resourceSAMLGroupMappingCreate,
		ReadContext:   resourceSAMLGroupMappingRead,
		UpdateContext: resourceSAMLGroupMappingUpdate,
		DeleteContext: resourceSAMLGroupMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// schema for import id: mapping|<saml_idp_id>|<provider_group_id>|<project_ids>|<role>

				mappingToImport, err := extractIDsFromSamlIdpGroupMappingImportID(d.Id())
				if err != nil {
					return nil, err
				}

				err = d.Set("saml_idp_id", mappingToImport.SamlIdpID)
				if err != nil {
					return nil, err
				}

				err = d.Set("provider_group_id", mappingToImport.ProviderGroupID)
				if err != nil {
					return nil, err
				}

				err = d.Set("role", mappingToImport.Role)
				if err != nil {
					return nil, err
				}

				err = d.Set("projects", mappingToImport.ProjectIDs)
				if err != nil {
					return nil, err
				}

				d.SetId(uuid.NewString())

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceSAMLGroupMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLGroupMappingCreate called...")

	samlIdpID := d.Get("saml_idp_id").(string)
	providerGroupID := d.Get("provider_group_id").(string)
	role := d.Get("role").(string)
	projectIDs := utils.ConvertListToString(d.Get("projects").([]interface{}))

	// verify the mapping doesn't already exist
	matchingNode, diags := querySAMLGroupMappings(ctx, m, samlIdpID, providerGroupID, role, projectIDs)
	if len(diags) != 0 {
		return diags
	}

	if matchingNode != nil {
		return diag.Errorf("saml group mapping for group: %s and role: %s to project(s): %s already exists for saml idp provider: %s and should be imported instead",
			providerGroupID, role, strings.Join(projectIDs, ", "), samlIdpID)
	}

	// define the graphql query
	query := `mutation SetSAMLGroupMapping ($input: ModifySAMLGroupMappingInput!) {
	  modifySAMLIdentityProviderGroupMappings(input: $input) {
            _stub
          }
	}`
	// populate the graphql variables
	vars := &wiz.UpdateSAMLGroupMappingInput{}
	vars.ID = samlIdpID
	vars.Patch.Upsert.ProviderGroupID = providerGroupID
	vars.Patch.Upsert.Role = role
	vars.Patch.Upsert.Projects = projectIDs

	// process the request
	data := &UpdateSAMLGroupMappingPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_group_mapping", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(uuid.NewString())

	return resourceSAMLGroupMappingRead(ctx, d, m)
}

func extractIDsFromSamlIdpGroupMappingImportID(id string) (SAMLGroupMappingsImport, error) {
	parts := strings.Split(id, "|")
	if len(parts) != 5 {
		return SAMLGroupMappingsImport{}, errors.New("invalid ID format")
	}

	// if user species the mapping to be global we return an empty slice
	var projectIDs []string
	if parts[3] != "global" {
		for _, projectID := range strings.Split(parts[3], ",") {
			projectIDs = append(projectIDs, strings.TrimSpace(projectID))
		}
	}

	return SAMLGroupMappingsImport{
		SamlIdpID:       parts[1],
		ProviderGroupID: parts[2],
		ProjectIDs:      projectIDs,
		Role:            parts[4],
	}, nil
}

func extractProjectIDs(projects []wiz.Project) []string {
	projectIDs := make([]string, len(projects))
	for i, project := range projects {
		projectIDs[i] = project.ID
	}

	return projectIDs
}

func resourceSAMLGroupMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLGroupMappingRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}
	samlIdpID := d.Get("saml_idp_id").(string)
	providerGroupID := d.Get("provider_group_id").(string)
	role := d.Get("role").(string)
	projectIDs := utils.ConvertListToString(d.Get("projects").([]interface{}))

	matchingNode, diags := querySAMLGroupMappings(ctx, m, samlIdpID, providerGroupID, role, projectIDs)
	if len(diags) > 0 {
		return diags
	}

	// If no matching node was found, return error
	if matchingNode == nil {
		return diag.Errorf("saml group mapping for group: %s not found for saml idp provider: %s", providerGroupID, samlIdpID)
	}

	// set the resource parameters
	err := d.Set("saml_idp_id", samlIdpID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("provider_group_id", matchingNode.ProviderGroupID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("role", matchingNode.Role.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	projectIDs = extractProjectIDs(matchingNode.Projects)
	err = d.Set("projects", projectIDs)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceSAMLGroupMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLGroupMappingUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation SetSAMLGroupMapping ($input: ModifySAMLGroupMappingInput!) {
	  modifySAMLIdentityProviderGroupMappings(input: $input) {
            _stub
          }
	}`

	samlIdpID := d.Get("saml_idp_id").(string)
	providerGroupID := d.Get("provider_group_id").(string)
	role := d.Get("role").(string)
	projects := utils.ConvertListToString(d.Get("projects").([]interface{}))

	// populate the graphql variables
	vars := &wiz.UpdateSAMLGroupMappingInput{}
	vars.ID = samlIdpID
	vars.Patch.Upsert.ProviderGroupID = providerGroupID
	vars.Patch.Upsert.Role = role
	vars.Patch.Upsert.Projects = projects

	// process the request
	data := &UpdateSAMLGroupMappingPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_group_mapping", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceSAMLGroupMappingRead(ctx, d, m)
}

func resourceSAMLGroupMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizSAMLGroupMappingDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation SetSAMLGroupMapping ($input: ModifySAMLGroupMappingInput!) {
	  modifySAMLIdentityProviderGroupMappings(input: $input) {
            _stub
          }
	}`

	samlIdpID := d.Get("saml_idp_id").(string)
	providerGroupID := d.Get("provider_group_id").(string)

	// populate the graphql variables
	vars := &DeleteSAMLGroupMappingInput{}
	vars.ID = samlIdpID
	vars.Patch.Delete = []string{providerGroupID}

	// process the request
	data := &UpdateSAMLGroupMappingPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_group_mapping", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}

func querySAMLGroupMappings(ctx context.Context, m interface{}, samlIdpID string, providerGroupID string, roleId string, projectIDs []string) (*wiz.SAMLGroupMapping, diag.Diagnostics) {
	// define the graphql query
	query := `query samlIdentityProviderGroupMappings ($id: ID!, $first: Int! $after: String){
	    samlIdentityProviderGroupMappings (
	        id: 	$id,
			first: 	$first
			after: 	$after
	    ) {
			pageInfo {
				  hasNextPage
				  endCursor
			}
	        nodes {
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
	vars.ID = samlIdpID
	vars.First = 100

	// Call ProcessPagedRequest
	diags, allData := client.ProcessPagedRequest(ctx, m, vars, &ReadSAMLGroupMappings{}, query, "saml_idp", "read", 0)
	if diags.HasError() {
		return nil, diags
	}

	var matchingNode *wiz.SAMLGroupMapping
	// Process the data...
	for _, data := range allData {
		typedData, ok := data.(*ReadSAMLGroupMappings)
		if !ok {
			return nil, diag.Errorf("data is not of type *ReadSAMLGroupMappings")
		}
		nodes := typedData.SAMLGroupMappings.Nodes
		for _, node := range nodes {
			nodeProjectIDs := extractProjectIDs(node.Projects)
			// If we find a match, store the node and break the loop
			if node.ProviderGroupID == providerGroupID && node.Role.ID == roleId && slices.Equal(projectIDs, nodeProjectIDs) {
				matchingNode = node
				break
			}
		}
	}

	return matchingNode, nil
}
