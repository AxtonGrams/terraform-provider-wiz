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
// It includes a SAMLIdentityProviderGroupMappingsConnection object.
type ReadSAMLGroupMappings struct {
	SAMLGroupMappings wiz.SAMLIdentityProviderGroupMappingsConnection `json:"samlIdentityProviderGroupMappings"`
}

// UpdateSAMLGroupMappingInput struct
type UpdateSAMLGroupMappingInput struct {
	ID    string                          `json:"id"`
	Patch wiz.ModifySAMLGroupMappingPatch `json:"patch"`
}

// SAMLGroupMappingsImport struct
type SAMLGroupMappingsImport struct {
	SamlIdpID     string
	GroupMappings []wiz.SAMLGroupDetailsInput
}

// UpdateSAMLGroupMappingPayload struct
type UpdateSAMLGroupMappingPayload struct {
	SAMLGroupMapping wiz.SAMLGroupMapping `json:"samlGroupMapping,omitempty"`
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
			"group_mapping": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Group Mapping description",
						},
					},
				},
			},
		},
		CreateContext: resourceSAMLGroupMappingCreate,
		ReadContext:   resourceSAMLGroupMappingRead,
		UpdateContext: resourceSAMLGroupMappingUpdate,
		DeleteContext: resourceSAMLGroupMappingDelete,
		Importer: &schema.ResourceImporter{

			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// schema for import id: <saml_idp_id>|<group_mappings>
				mappingToImport, err := extractIDsFromSamlIdpGroupMappingImportID(d.Id())
				if err != nil {
					return nil, err
				}

				err = d.Set("saml_idp_id", mappingToImport.SamlIdpID)
				if err != nil {
					return nil, err
				}

				var groupMappings []map[string]interface{}
				for _, groupMapping := range mappingToImport.GroupMappings {
					groupMappingMap := map[string]interface{}{
						"provider_group_id": groupMapping.ProviderGroupID,
						"role":              groupMapping.Role,
						"projects":          groupMapping.Projects,
					}
					groupMappings = append(groupMappings, groupMappingMap)
				}

				err = d.Set("group_mappings", groupMappings)
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
	groupMappings := d.Get("group_mapping").(*schema.Set).List()

	var upsertGroupMappings []wiz.SAMLGroupDetailsInput
	for _, item := range groupMappings {
		groupMapping := item.(map[string]interface{})
		providerGroupID := groupMapping["provider_group_id"].(string)
		role := groupMapping["role"].(string)
		projectIDs := utils.ConvertListToString(groupMapping["projects"].([]interface{}))
		description := groupMapping["description"].(string)

		// verify the mapping doesn't already exist
		matchingNodes, diags := querySAMLGroupMappings(ctx, m, samlIdpID, groupMappings)
		if len(diags) != 0 {
			return diags
		}

		for _, matchingNode := range matchingNodes {
			if matchingNode.ProviderGroupID == providerGroupID && matchingNode.Role.ID == role && slices.Equal(projectIDs, extractProjectIDs(matchingNode.Projects)) {
				return diag.Errorf("saml group mapping for group: %s and role: %s to project(s): %s already exists for saml idp provider: %s and should be imported instead",
					providerGroupID, role, strings.Join(projectIDs, ", "), samlIdpID)
			}
		}

		upsertGroupMapping := wiz.SAMLGroupDetailsInput{
			ProviderGroupID: providerGroupID,
			Role:            role,
			Projects:        projectIDs,
			Description:     description,
		}
		upsertGroupMappings = append(upsertGroupMappings, upsertGroupMapping)
	}

	// define the graphql query
	query := `mutation SetSAMLGroupMapping ($input: ModifySAMLGroupMappingInput!) {
	  modifySAMLIdentityProviderGroupMappings(input: $input) {
			_stub
		  }
	}`

	// populate the graphql variables
	vars := &UpdateSAMLGroupMappingInput{}
	vars.ID = samlIdpID
	vars.Patch = wiz.ModifySAMLGroupMappingPatch{
		Upsert: &upsertGroupMappings,
	}

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

	if len(parts) != 3 {
		return SAMLGroupMappingsImport{}, errors.New("invalid ID format")
	}

	groupMappingStrings := strings.Split(parts[2], "#")
	var groupMappings []wiz.SAMLGroupDetailsInput
	for _, groupMappingString := range groupMappingStrings {
		groupMappingParts := strings.Split(groupMappingString, ":")
		if len(groupMappingParts) < 2 {
			return SAMLGroupMappingsImport{}, errors.New("invalid group mapping format")
		}

		providerGroupID := groupMappingParts[0]
		role := groupMappingParts[1]
		var projectIDs []string
		if len(groupMappingParts) > 2 && groupMappingParts[2] != "" {
			projectIDs = strings.Split(groupMappingParts[2], ",")
		}

		groupMapping := wiz.SAMLGroupDetailsInput{
			ProviderGroupID: providerGroupID,
			Role:            role,
			Projects:        projectIDs,
		}
		groupMappings = append(groupMappings, groupMapping)

	}

	return SAMLGroupMappingsImport{
		SamlIdpID:     parts[1],
		GroupMappings: groupMappings,
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
	groupMappings := d.Get("group_mapping").(*schema.Set).List()

	var newGroupMappings []interface{}

	matchingNodes, diags := querySAMLGroupMappings(ctx, m, samlIdpID, groupMappings)
	if len(diags) > 0 {
		return diags
	}

	for _, item := range groupMappings {
		groupMapping := item.(map[string]interface{})
		providerGroupID := groupMapping["provider_group_id"].(string)
		role := groupMapping["role"].(string)
		projectIDs := utils.ConvertListToString(groupMapping["projects"].([]interface{}))

		for _, matchingNode := range matchingNodes {
			if matchingNode.ProviderGroupID == providerGroupID && matchingNode.Role.ID == role && slices.Equal(projectIDs, extractProjectIDs(matchingNode.Projects)) {
				// set the resource parameters
				newGroupMapping := map[string]interface{}{
					"provider_group_id": matchingNode.ProviderGroupID,
					"role":              matchingNode.Role.ID,
					"projects":          extractProjectIDs(matchingNode.Projects),
					"description":       matchingNode.Description,
				}
				newGroupMappings = append(newGroupMappings, newGroupMapping)
			}
		}
	}

	err := d.Set("saml_idp_id", samlIdpID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("group_mapping", newGroupMappings)
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

	samlIdpID := d.Get("saml_idp_id").(string)
	groupMappings := d.Get("group_mapping").(*schema.Set).List()
	var upsertGroupMappings []wiz.SAMLGroupDetailsInput
	for _, item := range groupMappings {
		groupMapping := item.(map[string]interface{})
		providerGroupID := groupMapping["provider_group_id"].(string)
		role := groupMapping["role"].(string)
		projects := utils.ConvertListToString(groupMapping["projects"].([]interface{}))
		description := groupMapping["description"].(string)
		upsertGroupMapping := wiz.SAMLGroupDetailsInput{
			ProviderGroupID: providerGroupID,
			Role:            role,
			Projects:        projects,
			Description:     description,
		}
		upsertGroupMappings = append(upsertGroupMappings, upsertGroupMapping)
	}

	// define the graphql query
	query := `mutation SetSAMLGroupMapping ($input: ModifySAMLGroupMappingInput!) {
	  modifySAMLIdentityProviderGroupMappings(input: $input) {
			_stub
		  }
	}`

	// populate the graphql variables
	vars := &UpdateSAMLGroupMappingInput{}
	vars.ID = samlIdpID
	vars.Patch = wiz.ModifySAMLGroupMappingPatch{
		Upsert: &upsertGroupMappings,
	}

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

	samlIdpID := d.Get("saml_idp_id").(string)
	groupMappings := d.Get("group_mapping").(*schema.Set).List()

	var deleteGroupMappings []string
	for _, item := range groupMappings {
		groupMapping := item.(map[string]interface{})
		providerGroupID := groupMapping["provider_group_id"].(string)
		deleteGroupMappings = append(deleteGroupMappings, providerGroupID)
	}

	// define the graphql query
	query := `mutation SetSAMLGroupMapping ($input: ModifySAMLGroupMappingInput!) {
	  modifySAMLIdentityProviderGroupMappings(input: $input) {
			_stub
		  }
	}`

	// populate the graphql variables
	vars := &UpdateSAMLGroupMappingInput{}
	vars.ID = samlIdpID
	vars.Patch = wiz.ModifySAMLGroupMappingPatch{
		Delete: &deleteGroupMappings,
	}

	// process the request
	data := &UpdateSAMLGroupMappingPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "saml_group_mapping", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}

func querySAMLGroupMappings(ctx context.Context, m interface{}, samlIdpID string, groupMappings []interface{}) ([]*wiz.SAMLGroupMapping, diag.Diagnostics) {
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
			  description
			}
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = samlIdpID
	vars.First = 500

	// Call ProcessPagedRequest
	diags, allData := client.ProcessPagedRequest(ctx, m, vars, &ReadSAMLGroupMappings{}, query, "saml_idp", "read", 0)
	if diags.HasError() {
		return nil, diags
	}

	var matchingNodes []*wiz.SAMLGroupMapping

	// Process the data...
	for _, data := range allData {
		typedData, ok := data.(*ReadSAMLGroupMappings)
		if !ok {
			return nil, diag.Errorf("data is not of type *ReadSAMLGroupMappings")
		}

		nodes := typedData.SAMLGroupMappings.Nodes
		for _, node := range nodes {
			for _, item := range groupMappings {
				groupMapping := item.(map[string]interface{})
				providerGroupID := groupMapping["provider_group_id"].(string)
				roleID := groupMapping["role"].(string)
				projectIDs := utils.ConvertListToString(groupMapping["projects"].([]interface{}))
				nodeProjectIDs := extractProjectIDs(node.Projects)

				// If we find a match, store the node
				if node.ProviderGroupID == providerGroupID && node.Role.ID == roleID && slices.Equal(projectIDs, nodeProjectIDs) {
					matchingNodes = append(matchingNodes, node)
				}
			}
		}
	}

	return matchingNodes, nil
}
