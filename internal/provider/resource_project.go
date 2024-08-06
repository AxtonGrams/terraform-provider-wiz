package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizProject() *schema.Resource {
	return &schema.Resource{
		Description: "Projects let you group your cloud resources according to their users and/or purposes.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The project name to display in Wiz.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The project description.",
				Optional:    true,
			},
			"identifiers": {
				Type:        schema.TypeList,
				Description: "Identifiers for the project.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"archived": {
				Type:        schema.TypeBool,
				Description: "Whether the project is archived/inactive",
				Optional:    true,
				Default:     false,
			},
			"is_folder": {
				Type:        schema.TypeBool,
				Description: "Whether the project is a folder.",
				Optional:    true,
				Default:     false,
			},
			"business_unit": {
				Type:        schema.TypeString,
				Description: "The business unit to which the project belongs.",
				Optional:    true,
			},
			"parent_project_id": {
				Type:        schema.TypeString,
				Description: "The parent project ID.",
				Optional:    true,
			},
			"project_owners": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of project owner IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_champions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of security champions IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"slug": {
				Type:        schema.TypeString,
				Description: "Short identifier for the project. The value must be unique, even against archived projects, so a uuid is generated and used as the slug value.",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "Unique identifier for the project.",
				Computed:    true,
			},
			"risk_profile": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Contains risk profile related properties for the project",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"business_impact": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Business impact.\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.BusinessImpact,
								),
							),
							Optional: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.BusinessImpact,
									false,
								),
							),
						},
						"is_actively_developed": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project under active development?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.YesNoUnknown,
									false,
								),
							),
						},
						"has_authentication": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Does the project require authentication?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.YesNoUnknown,
									false,
								),
							),
						},
						"has_exposed_api": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Does the project expose an API?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.YesNoUnknown,
									false,
								),
							),
						},
						"is_internet_facing": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project Internet facing?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.YesNoUnknown,
									false,
								),
							),
						},
						"is_customer_facing": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project customer facing?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.YesNoUnknown,
									false,
								),
							),
						},
						"stores_data": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Does the project store data?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.YesNoUnknown,
									false,
								),
							),
						},
						"is_regulated": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project regulated?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.YesNoUnknown,
									false,
								),
							),
						},
						"sensitive_data_types": {
							Type: schema.TypeList,
							Description: fmt.Sprintf(
								"Sensitive Data Types.\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.ProjectDataType,
								),
							),
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(
									validation.StringInSlice(
										wiz.ProjectDataType,
										false,
									),
								),
							},
						},
						"regulatory_standards": {
							Type: schema.TypeList,
							Description: fmt.Sprintf(
								"Regulatory Standards.\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.RegulatoryStandard,
								),
							),
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(
									validation.StringInSlice(
										wiz.RegulatoryStandard,
										false,
									),
								),
							},
						},
					},
				},
			},
			"kubernetes_cluster_link": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associate the project with kubernetes clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kubernetes_cluster": {
							Type:        schema.TypeString,
							Description: "The Wiz internal identifier for the kubernetes cluster.",
							Required:    true,
						},
						"shared": {
							Type:        schema.TypeBool,
							Description: "Mark the kubernetes cluster as shared, in which case, specific namespaces can be linked. This needs to be set to `true` if `namespaces` are set.",
							Optional:    true,
							Default:     true,
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
						"namespaces": {
							Type:        schema.TypeList,
							Description: "The kubernetes namespaces to link. `shared` must be set to `true` if namespaces are set.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"cloud_account_link": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Please either use this embedded set or the resource wiz_project_cloud_account_link. " +
					"Associate the project directly with a cloud account by wiz identifier UID to organize all the subscription resources, issues, and findings within this project.",

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_account_id": {
							Type:        schema.TypeString,
							Description: "The Wiz internal identifier for the Cloud Account Subscription.",
							Required:    true,
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
				},
			},
			"cloud_organization_link": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associate the project with an organizational link to organize all the subscription resources, issues, and findings within this project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_organization": {
							Type:        schema.TypeString,
							Description: "The Wiz internal identifier for the Organizational Unit.",
							Required:    true,
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
							Default:     true,
						},
						"resource_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Please provide a list of strings for filtering by resource groups. `shared` must be true to define resource_groups.",
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
				},
			},
		},
		CustomizeDiff: customdiff.All(
			// projects cannot be changed from folder to non-folder or vice versa.
			// to accommodate for importing resources into state, we can't use `ForceNew` in the schema definition.
			customdiff.ForceNewIfChange("is_folder", func(ctx context.Context, old, new, meta any) bool {
				return old.(bool) != new.(bool)
			},
			),
			func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
				var errMsgs []string

				if isFolder, ok := diff.Get("is_folder").(bool); ok && isFolder {
					resources := []string{"cloud_account_link", "cloud_organization_link", "kubernetes_cluster_link"}
					for _, link := range resources {
						if _, ok := diff.GetOk(link); ok {
							errMsgs = append(errMsgs, fmt.Sprintf("'%s' cannot be set if 'is_folder' is true", link))
						}
					}
				}
				if len(errMsgs) > 0 {
					var sb strings.Builder
					for _, errMsg := range errMsgs {
						sb.WriteString(errMsg)
						sb.WriteString("; ")
					}
					// Remove the last delimiters "; "
					return fmt.Errorf(sb.String()[:sb.Len()-2])
				}
				return nil
			},
		),
		CreateContext: resourceWizProjectCreate,
		ReadContext:   resourceWizProjectRead,
		UpdateContext: resourceWizProjectUpdate,
		DeleteContext: resourceWizProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func getOrganizationLinksVar(ctx context.Context, d *schema.ResourceData) []*wiz.ProjectCloudOrganizationLinkInput {
	linkSet := d.Get("cloud_organization_link").(*schema.Set).List()
	var myLinks []*wiz.ProjectCloudOrganizationLinkInput
	for _, y := range linkSet {
		var localLink wiz.ProjectCloudOrganizationLinkInput
		for a, b := range y.(map[string]interface{}) {
			switch a {
			case "environment":
				localLink.Environment = b.(string)
			case "cloud_organization":
				localLink.CloudOrganization = b.(string)
			case "shared":
				localLink.Shared = b.(bool)
			case "resource_groups":
				rgs := utils.ConvertListToString(b.([]interface{}))
				if len(rgs) > 0 {
					localLink.ResourceGroups = rgs
				}
			case "resource_tags":
				var myResourceTags []*wiz.ResourceTagInput
				for _, d := range b.(*schema.Set).List() {
					var localResourceTag wiz.ResourceTagInput
					for e, f := range d.(map[string]interface{}) {
						if e == "key" {
							localResourceTag.Key = f.(string)
						}
						if e == "value" {
							localResourceTag.Value = f.(string)
						}
					}
					myResourceTags = append(myResourceTags, &localResourceTag)
				}
				localLink.ResourceTags = myResourceTags
			}
		}
		myLinks = append(myLinks, &localLink)
	}
	return myLinks
}

func getKubernetesClusterLinksVar(ctx context.Context, d *schema.ResourceData) []*wiz.ProjectKubernetesClusterLinkInput {
	clusterSet := d.Get("kubernetes_cluster_link").(*schema.Set).List()
	var myClusters []*wiz.ProjectKubernetesClusterLinkInput
	for _, y := range clusterSet {
		var localCluster wiz.ProjectKubernetesClusterLinkInput
		for a, b := range y.(map[string]interface{}) {
			switch a {
			case "environment":
				localCluster.Environment = b.(string)
			case "shared":
				localCluster.Shared = b.(bool)
			case "kubernetes_cluster":
				localCluster.KubernetesCluster = b.(string)
			case "namespaces":
				localCluster.Namespaces = utils.ConvertListToString(b.([]interface{}))
			}
		}
		myClusters = append(myClusters, &localCluster)
	}

	return myClusters
}

func getAccountLinksVar(ctx context.Context, d *schema.ResourceData) []*wiz.ProjectCloudAccountLinkInput {
	accountSet := d.Get("cloud_account_link").(*schema.Set).List()
	var myAccounts []*wiz.ProjectCloudAccountLinkInput
	for _, y := range accountSet {
		var localAccount wiz.ProjectCloudAccountLinkInput
		for a, b := range y.(map[string]interface{}) {
			switch a {
			case "environment":
				localAccount.Environment = b.(string)
			case "cloud_account_id":
				localAccount.CloudAccount = b.(string)
			case "shared":
				localAccount.Shared = utils.ConvertBoolToPointer(b.(bool))
			case "resource_groups":
				rgs := utils.ConvertListToString(b.([]interface{}))
				if len(rgs) > 0 {
					localAccount.ResourceGroups = rgs
				}
			case "resource_tags":
				var myResourceTags []*wiz.ResourceTagInput
				for _, d := range b.(*schema.Set).List() {
					var localResourceTag wiz.ResourceTagInput
					for e, f := range d.(map[string]interface{}) {
						if e == "key" {
							localResourceTag.Key = f.(string)
						}
						if e == "value" {
							localResourceTag.Value = f.(string)
						}
					}
					myResourceTags = append(myResourceTags, &localResourceTag)
				}
				localAccount.ResourceTags = myResourceTags
			}
		}

		myAccounts = append(myAccounts, &localAccount)
	}
	return myAccounts
}

// CreateProject struct
type CreateProject struct {
	CreateProject wiz.CreateProjectPayload `json:"createProject"`
}

func resourceWizProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectCreate called...")

	// define the graphql query
	query := `mutation CreateProject($input: CreateProjectInput!) {
	  createProject(input: $input) {
	    project {
	      id
	    }
	  }
	}`

	// populate the graphql variables
	vars := &wiz.CreateProjectInput{}
	isFolder := d.Get("is_folder").(bool)
	vars.IsFolder = &isFolder

	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	vars.BusinessUnit = d.Get("business_unit").(string)
	vars.CloudOrganizationLinks = getOrganizationLinksVar(ctx, d)
	vars.CloudAccountLinks = getAccountLinksVar(ctx, d)
	vars.Identifiers = utils.ConvertListToString(d.Get("identifiers").([]interface{}))
	vars.KubernetesClusterLinks = getKubernetesClusterLinksVar(ctx, d)
	vars.RiskProfile.BusinessImpact = d.Get("risk_profile.0.business_impact").(string)
	vars.RiskProfile.IsActivelyDeveloped = d.Get("risk_profile.0.is_actively_developed").(string)
	vars.RiskProfile.HasAuthentication = d.Get("risk_profile.0.has_authentication").(string)
	vars.RiskProfile.HasExposedAPI = d.Get("risk_profile.0.has_exposed_api").(string)
	vars.RiskProfile.IsInternetFacing = d.Get("risk_profile.0.is_internet_facing").(string)
	vars.RiskProfile.IsCustomerFacing = d.Get("risk_profile.0.is_customer_facing").(string)
	vars.RiskProfile.StoresData = d.Get("risk_profile.0.stores_data").(string)
	vars.RiskProfile.IsRegulated = d.Get("risk_profile.0.is_regulated").(string)
	vars.RiskProfile.SensitiveDataTypes = utils.ConvertListToString(d.Get("risk_profile.0.sensitive_data_types").([]interface{}))
	vars.RiskProfile.RegulatoryStandards = utils.ConvertListToString(d.Get("risk_profile.0.regulatory_standards").([]interface{}))
	vars.ParentProjectID = d.Get("parent_project_id").(string)
	vars.ProjectOwners = utils.ConvertListToString(d.Get("project_owners").([]interface{}))
	vars.SecurityChampion = utils.ConvertListToString(d.Get("security_champions").([]interface{}))
	vars.Slug = uuid.New().String()

	// process the request
	data := &CreateProject{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "project", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateProject.Project.ID)

	return resourceWizProjectRead(ctx, d, m)
}

func flattenRiskProfile(ctx context.Context, riskProfile *wiz.ProjectRiskProfile) []interface{} {
	var output = make([]interface{}, 0)
	riskProfileMap := make(map[string]interface{})
	riskProfileMap["business_impact"] = riskProfile.BusinessImpact
	riskProfileMap["is_actively_developed"] = riskProfile.IsActivelyDeveloped
	riskProfileMap["has_authentication"] = riskProfile.HasAuthentication
	riskProfileMap["has_exposed_api"] = riskProfile.HasExposedAPI
	riskProfileMap["is_internet_facing"] = riskProfile.IsInternetFacing
	riskProfileMap["is_customer_facing"] = riskProfile.IsCustomerFacing
	riskProfileMap["stores_data"] = riskProfile.StoresData
	riskProfileMap["is_regulated"] = riskProfile.IsRegulated

	var sensitiveDataTypes = make([]interface{}, 0)
	for _, a := range riskProfile.SensitiveDataTypes {
		tflog.Trace(ctx, fmt.Sprintf("a: %T %s", a, utils.PrettyPrint(a)))
		sensitiveDataTypes = append(sensitiveDataTypes, a)
	}
	riskProfileMap["sensitive_data_types"] = sensitiveDataTypes

	var regulatoryStandards = make([]interface{}, 0)
	for _, a := range riskProfile.RegulatoryStandards {
		tflog.Trace(ctx, fmt.Sprintf("a: %T %s", a, utils.PrettyPrint(a)))
		regulatoryStandards = append(regulatoryStandards, a)
	}
	riskProfileMap["regulatory_standards"] = regulatoryStandards

	output = append(output, riskProfileMap)
	return output
}

func flattenCloudOrganizationLinks(ctx context.Context, cloudOrganizationLink []*wiz.ProjectCloudOrganizationLink) []interface{} {
	var output = make([]interface{}, 0)

	for _, b := range cloudOrganizationLink {
		cloudOrganizatinLinksMap := make(map[string]interface{})
		cloudOrganizatinLinksMap["cloud_organization"] = b.CloudOrganization.ID
		cloudOrganizatinLinksMap["shared"] = b.Shared
		cloudOrganizatinLinksMap["environment"] = b.Environment

		var resourceTags = make([]interface{}, 0)
		for _, d := range b.ResourceTags {
			var resourceTag = make(map[string]interface{})
			resourceTag["key"] = d.Key
			resourceTag["value"] = d.Value
			resourceTags = append(resourceTags, resourceTag)
		}

		var resourceGroups = make([]interface{}, 0)
		for _, d := range b.ResourceGroups {
			resourceGroups = append(resourceGroups, d)
		}

		cloudOrganizatinLinksMap["resource_tags"] = resourceTags
		cloudOrganizatinLinksMap["resource_groups"] = resourceGroups

		output = append(output, cloudOrganizatinLinksMap)
	}
	return output
}

func flattenCloudAccountLinks(ctx context.Context, cloudAccountLink []*wiz.ProjectCloudAccountLink) []interface{} {
	var output = make([]interface{}, 0)

	for _, b := range cloudAccountLink {
		cloudAccountLinksMap := make(map[string]interface{})
		cloudAccountLinksMap["cloud_account_id"] = b.CloudAccount.ID
		cloudAccountLinksMap["shared"] = b.Shared
		cloudAccountLinksMap["environment"] = b.Environment

		var resourceTags = make([]interface{}, 0)
		for _, d := range b.ResourceTags {
			var resourceTag = make(map[string]interface{})
			resourceTag["key"] = d.Key
			resourceTag["value"] = d.Value
			resourceTags = append(resourceTags, resourceTag)
		}

		var resourceGroups = make([]interface{}, 0)
		for _, d := range b.ResourceGroups {
			resourceGroups = append(resourceGroups, d)
		}

		cloudAccountLinksMap["resource_tags"] = resourceTags
		cloudAccountLinksMap["resource_groups"] = resourceGroups

		output = append(output, cloudAccountLinksMap)
	}
	return output
}

func flattenKubernetesClusterLinks(ctx context.Context, kubernetesClusterLink []*wiz.ProjectKubernetesClusterLink) []interface{} {
	var output = make([]interface{}, 0)

	for _, b := range kubernetesClusterLink {
		clusterLinksMap := make(map[string]interface{})
		clusterLinksMap["kubernetes_cluster"] = b.KubernetesCluster.ID
		clusterLinksMap["shared"] = b.Shared
		clusterLinksMap["environment"] = b.Environment

		var namespaces = make([]interface{}, 0)
		for _, d := range b.Namespaces {
			namespaces = append(namespaces, d)
		}
		clusterLinksMap["namespaces"] = namespaces

		output = append(output, clusterLinksMap)
	}
	return output
}

func flattenUserIds(ctx context.Context, userIds []*wiz.User) []interface{} {
	var output = make([]interface{}, 0)
	for _, b := range userIds {
		output = append(output, b.ID)
	}
	return output
}

// ReadProjectPayload struct -- updates
type ReadProjectPayload struct {
	Project wiz.Project `json:"project"`
}

func resourceWizProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectRead called...")

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

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	data := &ReadProjectPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "project", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	err := d.Set("name", data.Project.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("is_folder", data.Project.IsFolder)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// the parent project will be the first element of the list of ancestor projects
	if len(data.Project.AncestorProjects) > 0 {
		err = d.Set("parent_project_id", data.Project.AncestorProjects[0].ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	err = d.Set("description", data.Project.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("archived", data.Project.Archived)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("slug", data.Project.Slug)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	identifiers := utils.ConvertSliceToGenericArray(data.Project.Identifiers)
	if err = d.Set("identifiers", identifiers); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	projectOwners := flattenUserIds(ctx, data.Project.ProjectOwners)
	if err := d.Set("project_owners", projectOwners); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	securityChampions := flattenUserIds(ctx, data.Project.SecurityChampions)
	if err := d.Set("security_champions", securityChampions); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("business_unit", data.Project.BusinessUnit)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	RiskProfile := flattenRiskProfile(ctx, &data.Project.RiskProfile)
	if err := d.Set("risk_profile", RiskProfile); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	cloudOrganizationLinks := flattenCloudOrganizationLinks(ctx, data.Project.CloudOrganizationLinks)
	if err := d.Set("cloud_organization_link", cloudOrganizationLinks); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	cloudAccountLinks := flattenCloudAccountLinks(ctx, data.Project.CloudAccountLinks)
	if err := d.Set("cloud_account_link", cloudAccountLinks); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	kubernetesClusterLinks := flattenKubernetesClusterLinks(ctx, data.Project.KubernetesClustersLinks)
	if err := d.Set("kubernetes_cluster_link", kubernetesClusterLinks); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateProject struct
type UpdateProject struct {
	UpdateProject wiz.UpdateProjectPayload `json:"updateProject"`
}

/*
  In order to effectively manage the lifecycle of project settings, we need to override instead of patch
  Resource state change detection is not observed for UpdateContext, all attributes are captured in the request
*/

func resourceWizProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateProject($input: UpdateProjectInput!) {
	  updateProject(input: $input) {
	    project {
	      id
	    }
	  }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateProjectInput{}
	vars.ID = d.Id()

	vars.Override.Archived = utils.ConvertBoolToPointer(d.Get("archived").(bool))
	vars.Override.Name = d.Get("name").(string)
	vars.Override.Description = d.Get("description").(string)
	vars.Override.BusinessUnit = d.Get("business_unit").(string)
	vars.Override.Identifiers = utils.ConvertListToString((d.Get("identifiers")).([]interface{}))
	vars.Override.Slug = d.Get("slug").(string)
	vars.Override.ParentProjectID = d.Get("parent_project_id").(string)

	// The API treats a patch to riskProfile as an override so we set all values
	riskProfile := &wiz.ProjectRiskProfileInput{}
	riskProfile.BusinessImpact = d.Get("risk_profile.0.business_impact").(string)
	riskProfile.IsActivelyDeveloped = d.Get("risk_profile.0.is_actively_developed").(string)
	riskProfile.HasAuthentication = d.Get("risk_profile.0.has_authentication").(string)
	riskProfile.HasExposedAPI = d.Get("risk_profile.0.has_exposed_api").(string)
	riskProfile.IsInternetFacing = d.Get("risk_profile.0.is_internet_facing").(string)
	riskProfile.IsCustomerFacing = d.Get("risk_profile.0.is_customer_facing").(string)
	riskProfile.StoresData = d.Get("risk_profile.0.stores_data").(string)
	riskProfile.IsRegulated = d.Get("risk_profile.0.is_regulated").(string)
	riskProfile.SensitiveDataTypes = utils.ConvertListToString((d.Get("risk_profile.0.sensitive_data_types")).([]interface{}))
	riskProfile.RegulatoryStandards = utils.ConvertListToString((d.Get("risk_profile.0.regulatory_standards")).([]interface{}))
	vars.Override.RiskProfile = riskProfile

	vars.Override.ProjectOwners = utils.ConvertListToString((d.Get("project_owners")).([]interface{}))
	vars.Override.SecurityChampions = utils.ConvertListToString((d.Get("security_champions")).([]interface{}))

	var updateOrgLinks = []*wiz.ProjectCloudOrganizationLinkInput{}
	cloudOrganizationlinks := d.Get("cloud_organization_link").(*schema.Set).List()
	for _, b := range cloudOrganizationlinks {
		var updateOrgLink = &wiz.ProjectCloudOrganizationLinkInput{}
		for c, d := range b.(map[string]interface{}) {
			switch c {
			case "environment":
				updateOrgLink.Environment = d.(string)
			case "shared":
				updateOrgLink.Shared = d.(bool)
			case "cloud_organization":
				updateOrgLink.CloudOrganization = d.(string)
			case "resource_groups":
				updateOrgLink.ResourceGroups = utils.ConvertListToString(d.([]interface{}))
			case "resource_tags":
				var updateResourceTags = []*wiz.ResourceTagInput{}
				for _, f := range d.(*schema.Set).List() {
					var updateResourceTag = &wiz.ResourceTagInput{}
					for g, h := range f.(map[string]interface{}) {
						if g == "key" {
							updateResourceTag.Key = h.(string)
						}
						if g == "value" {
							updateResourceTag.Value = h.(string)
						}
					}
					updateResourceTags = append(updateResourceTags, updateResourceTag)
				}
				updateOrgLink.ResourceTags = updateResourceTags
			}
		}
		updateOrgLinks = append(updateOrgLinks, updateOrgLink)
	}
	vars.Override.CloudOrganizationLinks = updateOrgLinks

	var updateAccountLinks = []*wiz.ProjectCloudAccountLinkInput{}
	cloudAccountLinks := d.Get("cloud_account_link").(*schema.Set).List()
	for _, b := range cloudAccountLinks {
		var updateAccountLink = &wiz.ProjectCloudAccountLinkInput{}
		for c, d := range b.(map[string]interface{}) {
			switch c {
			case "environment":
				updateAccountLink.Environment = d.(string)
			case "shared":
				updateAccountLink.Shared = utils.ConvertBoolToPointer(d.(bool))
			case "cloud_account_id":
				updateAccountLink.CloudAccount = d.(string)
			case "resource_groups":
				updateAccountLink.ResourceGroups = utils.ConvertListToString(d.([]interface{}))
			case "resource_tags":
				var updateResourceTags = []*wiz.ResourceTagInput{}
				for _, f := range d.(*schema.Set).List() {
					var updateResourceTag = &wiz.ResourceTagInput{}
					for g, h := range f.(map[string]interface{}) {
						if g == "key" {
							updateResourceTag.Key = h.(string)
						}
						if g == "value" {
							updateResourceTag.Value = h.(string)
						}
					}
					updateResourceTags = append(updateResourceTags, updateResourceTag)
				}
				updateAccountLink.ResourceTags = updateResourceTags
			}
		}
		updateAccountLinks = append(updateAccountLinks, updateAccountLink)
	}
	vars.Override.CloudAccountLinks = updateAccountLinks

	var updateKubernetesClusterLinks = []*wiz.ProjectKubernetesClusterLinkInput{}
	kubernetesClusterLinks := d.Get("kubernetes_cluster_link").(*schema.Set).List()
	for _, b := range kubernetesClusterLinks {
		var updateClusterLink = &wiz.ProjectKubernetesClusterLinkInput{}
		for c, d := range b.(map[string]interface{}) {
			switch c {
			case "environment":
				updateClusterLink.Environment = d.(string)
			case "shared":
				updateClusterLink.Shared = d.(bool)
			case "namespaces":
				updateClusterLink.Namespaces = utils.ConvertListToString(d.([]interface{}))
			case "kubernetes_cluster":
				updateClusterLink.KubernetesCluster = d.(string)
			}
		}
		updateKubernetesClusterLinks = append(updateKubernetesClusterLinks, updateClusterLink)
	}
	vars.Override.KubernetesClusterLinks = updateKubernetesClusterLinks

	// process the request
	data := &UpdateProject{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "project", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizProjectRead(ctx, d, m)
}

/*
  Wiz does not support deleting projects, so we fake it by setting archived=true
  Project names must be unique to the tenant, we set the project name to the slug uid to avoid conflicts (i.e. destroy/create)
  The slug itself also needs to be included in update/delete contexts
*/

func resourceWizProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizProjectDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateProject($input: UpdateProjectInput!) {
          updateProject(input: $input) {
            project {
              id
            }
          }
        }`

	// populate the graphql variables
	vars := &wiz.UpdateProjectInput{}
	vars.ID = d.Id()
	vars.Override.Name = d.Get("slug").(string)
	vars.Override.Slug = d.Get("slug").(string)
	archive := true
	vars.Override.Archived = &archive

	// process the request
	data := &UpdateProject{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "project", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
