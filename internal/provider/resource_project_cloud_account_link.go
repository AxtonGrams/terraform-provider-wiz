package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

type CloudAccountSearchResponse struct {
	GraphSearch struct {
		Nodes []struct {
			Entities []struct {
				Id string `json:"id"`
			} `json:"entities"`
		} `json:"nodes"`
	} `json:"graphSearch"`
}

type SearchForCloudAccountVars struct {
	ExternalId string `json:"externalId"`
	ProjectId  string `json:"projectId"`
	Quick      bool   `json:"quick"`
}

type PartialProjectWithCloudAccountLinks struct {
	Project PartialProject `json:"project"`
}

type PartialProject struct {
	CloudAccountLinks []*wiz.ProjectCloudAccountLink
}

type UpdateProjectCloudAccountLinks struct {
	ID    string                        `json:"id"`
	Patch PatchProjectCloudAccountLinks `json:"patch"`
}

type PatchProjectCloudAccountLinks struct {
	CloudAccountLinks []*wiz.ProjectCloudAccountLinkInput `json:"cloudAccountLinks"`
}

func resourceWizProjectCloudAccountLink() *schema.Resource {
	return &schema.Resource{
		Description: "Please either use this resource or the embedded set of Cloud Account Links in the wiz_project resource. " +
			"Link of a Project to a Cloud Account.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique tf-internal identifier for the project cloud account link",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "The Wiz internal identifier of the Wiz project to link the cloud account to",
				Required:    true,
				ForceNew:    true,
			},
			"cloud_account_id": {
				Type:        schema.TypeString,
				Description: "The Wiz internal identifier for the Cloud Account Subscription.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"external_cloud_account_id": {
				Type:        schema.TypeString,
				Description: "The external identifier for the Cloud Account, e.g. an azure subscription id or an aws account id.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"environment": {
				Type: schema.TypeString,
				Description: fmt.Sprintf(
					"The environment.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.Environment,
					),
				),
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.Environment,
						false,
					),
				),
				Default: "PRODUCTION",
			},
			"shared": {
				Type:        schema.TypeBool,
				Description: "Subscriptions that host a few projects can be marked as ‘shared subscriptions’ and resources can be filtered by tags.",
				Optional:    true,
			},
			"resource_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Please provide a list of resource group identifiers for filtering by resource groups. `shared` must be true to define resource_groups.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_tags": {
				Type:        schema.TypeSet,
				Description: "Provide a key and value pair for filtering resources. `shared` must be true to define resource_tags.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		CreateContext: resourceWizProjectCloudAccountLinkCreate,
		ReadContext:   resourceWizProjectCloudAccountLinkRead,
		UpdateContext: resourceWizProjectCloudAccountLinkUpdate,
		DeleteContext: resourceWizProjectCloudAccountLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// schema for import id: link|<project_id>|<cloud_account_id>
				projectId, cloudAccountId, err := extractIds(d.Id())
				if err != nil {
					return nil, err
				}

				err = d.Set("project_id", projectId)
				if err != nil {
					return nil, err
				}

				err = d.Set("cloud_account_id", cloudAccountId)
				if err != nil {
					return nil, err
				}

				d.SetId(uuid.NewString())

				return []*schema.ResourceData{d}, nil
			},
		},
		// allow the user to supply both 'cloud_account_id' and 'external_cloud_account_id'
		// if none is given, we return and error
		// if they do not match, we also return an error
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			cloudAccountId, cloudAccountIdOk := diff.GetOk("cloud_account_id")
			externalCloudAccountId, externalCloudAccountIdOk := diff.GetOk("external_cloud_account_id")
			if !cloudAccountIdOk && !externalCloudAccountIdOk {
				return fmt.Errorf("either cloud_account_id or external_cloud_account_id must be set")
			}

			if cloudAccountIdOk && externalCloudAccountIdOk {
				queriedAccountId, diags := searchForCloudAccount(ctx, externalCloudAccountId.(string), v)
				if len(diags) != 0 {
					return fmt.Errorf("error while searching for cloud account in wiz")
				}

				if queriedAccountId != cloudAccountId {
					return fmt.Errorf("cloud_account_id and external_cloud_account_id must correspond to the same account")
				}
			}

			return nil
		},
	}
}

func getAccountLinkVar(d *schema.ResourceData, cloudAccountId string) *wiz.ProjectCloudAccountLinkInput {
	var localAccount wiz.ProjectCloudAccountLinkInput

	localAccount.Environment = d.Get("environment").(string)
	localAccount.CloudAccount = cloudAccountId
	localAccount.Shared = utils.ConvertBoolToPointer(d.Get("shared").(bool))
	rgs := utils.ConvertListToString(d.Get("resource_groups").([]interface{}))
	if len(rgs) > 0 {
		localAccount.ResourceGroups = rgs
	}

	// var myResourceTags []*wiz.ResourceTagInput
	for _, d := range d.Get("resource_tags").(*schema.Set).List() {
		var localResourceTag wiz.ResourceTagInput
		for e, f := range d.(map[string]interface{}) {
			if e == "key" {
				localResourceTag.Key = f.(string)
			}
			if e == "value" {
				localResourceTag.Value = f.(string)
			}
		}
		// myResourceTags = append(myResourceTags, &localResourceTag)
		localAccount.ResourceTags = append(localAccount.ResourceTags, &localResourceTag)
	}

	return &localAccount
}

// this is needed, as we query for existing cloud account links, then need
// to send the list with an appended entry back as mutation - but the types are different between
// GET and PATCH.
func accountLinkToAccountLinkInput(link *wiz.ProjectCloudAccountLink) *wiz.ProjectCloudAccountLinkInput {
	resourceTagsInput := make([]*wiz.ResourceTagInput, len(link.ResourceTags))
	for i, tag := range link.ResourceTags {
		resourceTagsInput[i] = &wiz.ResourceTagInput{
			Key:   tag.Key,
			Value: tag.Value,
		}
	}

	return &wiz.ProjectCloudAccountLinkInput{
		CloudAccount:   link.CloudAccount.ID,
		Environment:    link.Environment,
		ResourceGroups: link.ResourceGroups,
		ResourceTags:   resourceTagsInput,
		Shared:         &link.Shared,
	}
}

func searchForCloudAccount(ctx context.Context, externalId string, m interface{}) (string, diag.Diagnostics) {
	tflog.Info(ctx, "searching for account in wiz inventory...")

	readCloudAccountsQuery := `query SearchForCloudAccount($externalId: String!, $projectId: String!, $quick: Boolean) {
		graphSearch(query: {
			type: [SUBSCRIPTION],
			where: {
				externalId: {
					EQUALS: [$externalId]
				}
			}
		},
			projectId: $projectId, quick: $quick) {
				nodes {
					entities {
						id
					}
				}
			}
		}`

	vars := &SearchForCloudAccountVars{
		ExternalId: externalId,
		ProjectId:  "*",
		Quick:      true,
	}

	respData := &CloudAccountSearchResponse{}
	diags := client.ProcessRequest(ctx, m, vars, respData, readCloudAccountsQuery, "SearchForCloudAccount", "read")
	if len(diags) > 0 {
		return "", diags

	}

	if len(respData.GraphSearch.Nodes) == 0 || len(respData.GraphSearch.Nodes[0].Entities) == 0 {
		return "", diag.Errorf("cloud account %s not found in wiz inventory", externalId)
	}

	return respData.GraphSearch.Nodes[0].Entities[0].Id, nil
}

func resourceWizProjectCloudAccountLinkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectCloudAccountLinkCreate called...")
	projectId := d.Get("project_id").(string)
	var cloudAccountWizId string

	if v, ok := d.GetOk("cloud_account_id"); ok {
		cloudAccountWizId = v.(string)
	} else {
		cloudAccountUpstreamId := d.Get("external_cloud_account_id").(string)
		var diagsSearch diag.Diagnostics
		cloudAccountWizId, diagsSearch = searchForCloudAccount(ctx, cloudAccountUpstreamId, m)
		if len(diagsSearch) > 0 {
			return diagsSearch
		}

	}
	// verify that the link does not already exist in wiz
	// if it does, abort and throw an error, as is standard
	// terraform behavior (no overwrite or implicit import).
	partialProject := &PartialProjectWithCloudAccountLinks{}
	linkExists, requestDiags := checkCloudAccountLinkExistence(ctx, m, projectId, cloudAccountWizId, partialProject)
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	if linkExists {
		return diag.Errorf("cloud account %s is already linked to project %s", cloudAccountWizId, projectId)
	}

	// link not present, add it to the project
	newCloudAccountLinksList := make([]*wiz.ProjectCloudAccountLinkInput, len(partialProject.Project.CloudAccountLinks)+1)
	for i, link := range partialProject.Project.CloudAccountLinks {
		newCloudAccountLinksList[i] = accountLinkToAccountLinkInput(link)
	}
	newCloudAccountLinksList[len(newCloudAccountLinksList)-1] = getAccountLinkVar(d, cloudAccountWizId)

	// define the graphql query for adding the link by taking the existing list and appending
	// the new entry to it - then patch this property on the wiz project
	query := `mutation LinkCloudAccountToProject($input: UpdateProjectInput!) {
		updateProject(input: $input) {
			project {
				id
			}
		}
	}`

	// populate the graphql variables
	vars := &UpdateProjectCloudAccountLinks{
		ID: projectId,
		Patch: PatchProjectCloudAccountLinks{
			CloudAccountLinks: newCloudAccountLinksList,
		},
	}

	// process the request
	data := &UpdateProject{}
	requestDiags = client.ProcessRequest(ctx, m, vars, data, query, "LinkCloudAccountToProject", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	d.SetId(uuid.NewString())
	err := d.Set("cloud_account_id", cloudAccountWizId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return resourceWizProjectCloudAccountLinkRead(ctx, d, m)
}

func resourceWizProjectCloudAccountLinkRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectCloudAccountLinkRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query project  (
	    $id: ID
	){
	    project(
	        id: $id
	    ) {
	        id
	        name
	        isFolder
	        ancestorProjects {
	          id
	        }	
	        description
	        identifiers
	        slug
	        archived
	        businessUnit
	        projectOwners {
	            id
	            name
	            email
	        }
	        securityChampions {
	            id
	            name
	            email
	        }
	        riskProfile {
	            businessImpact
	            isActivelyDeveloped
	            hasAuthentication
	            hasExposedAPI
	            isInternetFacing
	            isCustomerFacing
	            storesData
	            sensitiveDataTypes
	            isRegulated
	            regulatoryStandards
	        }
	        cloudOrganizationLinks {
	            cloudOrganization {
	                externalId
	                id
	                name
	                path
	            }
	            resourceTags {
	                key
	                value
	            }
	            resourceGroups
	            shared
	            environment
	        }
	        cloudAccountLinks {
	            cloudAccount {
	                externalId
	                id
	                name
	            }
	            resourceTags {
	                key
	                value
	            }
	            resourceGroups
	            shared
	            environment
	        }
	        kubernetesClustersLinks {
	            kubernetesCluster {
	                id
	            }
	            environment
	            namespaces
	            shared
              }
	    }
	}`

	projectId := d.Get("project_id").(string)
	cloudAccountWizId := d.Get("cloud_account_id").(string)

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = projectId

	// process the request
	data := &ReadProjectPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "project", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	err := d.Set("project_id", data.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// extract the single cloud account link we want
	cloudAccountLink, err := extractCloudAccountLink(data.Project.CloudAccountLinks, cloudAccountWizId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("cloud_account_id", cloudAccountLink.CloudAccount.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("external_cloud_account_id", cloudAccountLink.CloudAccount.ExternalID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("environment", cloudAccountLink.Environment)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("shared", cloudAccountLink.Shared)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("resource_groups", cloudAccountLink.ResourceGroups)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("resource_tags", cloudAccountLink.ResourceTags)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

func extractIds(id string) (string, string, error) {
	parts := strings.Split(id, "|")
	if len(parts) != 3 {
		return "", "", errors.New("invalid ID format")
	}

	return parts[1], parts[2], nil
}

func extractCloudAccountLink(cloudAccountLinks []*wiz.ProjectCloudAccountLink, wizCloudAccountId string) (*wiz.ProjectCloudAccountLink, error) {
	for _, cloudAccountLink := range cloudAccountLinks {
		if cloudAccountLink.CloudAccount.ID == wizCloudAccountId {
			return cloudAccountLink, nil
		}
	}

	return nil, fmt.Errorf("cloud account with id %s not found in cloud account links of project", wizCloudAccountId)
}

func resourceWizProjectCloudAccountLinkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectCloudAccountLinkUpdate called...")
	projectId := d.Get("project_id").(string)
	cloudAccountWizId := d.Get("cloud_account_id").(string)

	// verify that the link exists in wiz
	partialProject := &PartialProjectWithCloudAccountLinks{}
	linkExists, requestDiags := checkCloudAccountLinkExistence(ctx, m, projectId, cloudAccountWizId, partialProject)
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	if !linkExists {
		return diag.Errorf("cloud account with id %s not found in cloud account links of project %s", cloudAccountWizId, projectId)
	}

	newCloudAccountLinksList := make([]*wiz.ProjectCloudAccountLinkInput, len(partialProject.Project.CloudAccountLinks)+1)
	for i, link := range partialProject.Project.CloudAccountLinks {
		newCloudAccountLinksList[i] = accountLinkToAccountLinkInput(link)
	}
	newCloudAccountLinksList[len(newCloudAccountLinksList)-1] = getAccountLinkVar(d, cloudAccountWizId)

	query := `mutation LinkCloudAccountToProject($input: UpdateProjectInput!) {
		updateProject(input: $input) {
			project {
				id
			}
		}
	}`

	// populate the graphql variables
	vars := &UpdateProjectCloudAccountLinks{
		ID: projectId,
		Patch: PatchProjectCloudAccountLinks{
			CloudAccountLinks: newCloudAccountLinksList,
		},
	}

	// process the request
	data := &UpdateProject{}
	requestDiags = client.ProcessRequest(ctx, m, vars, data, query, "LinkCloudAccountToProject", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizProjectCloudAccountLinkRead(ctx, d, m)
}

func resourceWizProjectCloudAccountLinkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectCloudAccountLinkDelete called...")
	projectId := d.Get("project_id").(string)
	cloudAccountWizId := d.Get("cloud_account_id").(string)

	// verify that the link exists in wiz
	partialProject := &PartialProjectWithCloudAccountLinks{}
	linkExists, requestDiags := checkCloudAccountLinkExistence(ctx, m, projectId, cloudAccountWizId, partialProject)
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	if !linkExists {
		return diag.Errorf("cloud account with id %s not found in cloud account links of project %s", cloudAccountWizId, projectId)
	}

	newCloudAccountLinksList := make([]*wiz.ProjectCloudAccountLinkInput, 0, len(partialProject.Project.CloudAccountLinks))
	for _, link := range partialProject.Project.CloudAccountLinks {
		if link.CloudAccount.ID != cloudAccountWizId {
			newCloudAccountLinksList = append(newCloudAccountLinksList, accountLinkToAccountLinkInput(link))
		}
	}

	query := `mutation LinkCloudAccountToProject($input: UpdateProjectInput!) {
		updateProject(input: $input) {
			project {
				id
			}
		}
	}`

	// populate the graphql variables
	vars := &UpdateProjectCloudAccountLinks{
		ID: projectId,
		Patch: PatchProjectCloudAccountLinks{
			CloudAccountLinks: newCloudAccountLinksList,
		},
	}

	// process the request
	data := &UpdateProject{}
	requestDiags = client.ProcessRequest(ctx, m, vars, data, query, "LinkCloudAccountToProject", "update")
	diags = append(diags, requestDiags...)

	return diags
}

func checkCloudAccountLinkExistence(ctx context.Context, m interface{}, projectId string, cloudAccountWizId string, partialProject *PartialProjectWithCloudAccountLinks) (exists bool, diags diag.Diagnostics) {
	readExistingLinksQuery := `query project ($id: ID) {
	    project(
	        id: $id
	    ) {
	        cloudAccountLinks {
	            cloudAccount {
	                externalId
	                id
	                name
	            }
	            resourceTags {
	                key
	                value
	            }
	            resourceGroups
	            shared
	            environment
	        }
		}
	}`

	// read existing cloud account links
	requestDiags := client.ProcessRequest(ctx, m,
		&internal.QueryVariables{ID: projectId}, partialProject, readExistingLinksQuery,
		"project_cloud_account_link", "read")

	// handle errors from read
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return false, diags
	}

	// check if desired link exists
	linkExists := slices.ContainsFunc(
		partialProject.Project.CloudAccountLinks,
		func(link *wiz.ProjectCloudAccountLink) bool {
			return link.CloudAccount.ID == cloudAccountWizId
		},
	)

	return linkExists, diags
}
