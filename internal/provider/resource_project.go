package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
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
			"archived": {
				Type:        schema.TypeBool,
				Description: "Whether the project is archived/inactive",
				Optional:    true,
				Default:     false,
			},
			"business_unit": {
				Type:        schema.TypeString,
				Description: "The business unit to which the project belongs.",
				Optional:    true,
			},
			"slug": {
				Type:        schema.TypeString,
				Description: "Short identifier for the project. The value must be unique, even against archived projects, so a uuid is generated and used as the slug value.",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "Unique identifier for the project",
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
									vendor.BusinessImpact,
								),
							),
							Optional: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.BusinessImpact,
									false,
								),
							),
						},
						"is_actively_developed": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project under active development?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.YesNoUnknown,
									false,
								),
							),
						},
						"has_authentication": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Does the project require authentication?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.YesNoUnknown,
									false,
								),
							),
						},
						"has_exposed_api": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Does the project expose an API?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.YesNoUnknown,
									false,
								),
							),
						},
						"is_internet_facing": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project Internet facing?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.YesNoUnknown,
									false,
								),
							),
						},
						"is_customer_facing": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project customer facing?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.YesNoUnknown,
									false,
								),
							),
						},
						"stores_data": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Does the project store data?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.YesNoUnknown,
									false,
								),
							),
						},
						"is_regulated": {
							Type: schema.TypeString,
							Description: fmt.Sprintf(
								"Is the project regulated?\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.YesNoUnknown,
								),
							),
							Optional: true,
							Default:  "UNKNOWN",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.YesNoUnknown,
									false,
								),
							),
						},
						"sensitive_data_types": {
							Type: schema.TypeList,
							Description: fmt.Sprintf(
								"Sensitive Data Types.\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									vendor.ProjectDataType,
								),
							),
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(
									validation.StringInSlice(
										vendor.ProjectDataType,
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
									vendor.RegulatoryStandard,
								),
							),
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(
									validation.StringInSlice(
										vendor.RegulatoryStandard,
										false,
									),
								),
							},
						},
					},
				},
			},
			"cloud_organization_link": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Associate the project with the resources and subscriptions to organize all the resources, issues, and findings within this project.",
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
									vendor.Environment,
								),
							),
							Optional: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									vendor.Environment,
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
		CreateContext: resourceWizProjectCreate,
		ReadContext:   resourceWizProjectRead,
		UpdateContext: resourceWizProjectUpdate,
		DeleteContext: resourceWizProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func getOrganizationLinksVar(ctx context.Context, d *schema.ResourceData) []*vendor.ProjectCloudOrganizationLinkInput {
	linkSet := d.Get("cloud_organization_link").(*schema.Set).List()
	var myLinks []*vendor.ProjectCloudOrganizationLinkInput
	for _, y := range linkSet {
		var localLink vendor.ProjectCloudOrganizationLinkInput
		for a, b := range y.(map[string]interface{}) {
			if a == "environment" {
				localLink.Environment = b.(string)
			}
			if a == "cloud_organization" {
				localLink.CloudOrganization = b.(string)
			}
			if a == "shared" {
				localLink.Shared = b.(bool)
			}
			if a == "resource_tags" {
				var myResourceTags []*vendor.ResourceTag
				for _, d := range b.(*schema.Set).List() {
					var localResourceTag vendor.ResourceTag
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

// CreateProject struct
type CreateProject struct {
	CreateProject vendor.CreateProjectPayload `json:"createProject"`
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
	vars := &vendor.CreateProjectInput{}
	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	vars.BusinessUnit = d.Get("business_unit").(string)
	vars.Slug = uuid.New().String()
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
	vars.CloudOrganizationLinks = getOrganizationLinksVar(ctx, d)

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

func flattenRiskProfile(ctx context.Context, riskProfile *vendor.ProjectRiskProfile) []interface{} {
	var output = make([]interface{}, 0, 0)
	riskProfileMap := make(map[string]interface{})
	riskProfileMap["business_impact"] = riskProfile.BusinessImpact
	riskProfileMap["is_actively_developed"] = riskProfile.IsActivelyDeveloped
	riskProfileMap["has_authentication"] = riskProfile.HasAuthentication
	riskProfileMap["has_exposed_api"] = riskProfile.HasExposedAPI
	riskProfileMap["is_internet_facing"] = riskProfile.IsInternetFacing
	riskProfileMap["is_customer_facing"] = riskProfile.IsCustomerFacing
	riskProfileMap["stores_data"] = riskProfile.StoresData
	riskProfileMap["is_regulated"] = riskProfile.IsRegulated

	var sensitiveDataTypes = make([]interface{}, 0, 0)
	for _, a := range riskProfile.SensitiveDataTypes {
		tflog.Trace(ctx, fmt.Sprintf("a: %T %s", a, utils.PrettyPrint(a)))
		sensitiveDataTypes = append(sensitiveDataTypes, a)
	}
	riskProfileMap["sensitive_data_types"] = sensitiveDataTypes

	var regulatoryStandards = make([]interface{}, 0, 0)
	for _, a := range riskProfile.RegulatoryStandards {
		tflog.Trace(ctx, fmt.Sprintf("a: %T %s", a, utils.PrettyPrint(a)))
		regulatoryStandards = append(regulatoryStandards, a)
	}
	riskProfileMap["regulatory_standards"] = regulatoryStandards

	output = append(output, riskProfileMap)
	return output
}

func flattenCloudOrganizationLinks(ctx context.Context, cloudOrganizationLink []*vendor.ProjectCloudOrganizationLink) []interface{} {
	var output = make([]interface{}, 0, 0)

	for _, b := range cloudOrganizationLink {
		cloudOrganizatinLinksMap := make(map[string]interface{})
		cloudOrganizatinLinksMap["cloud_organization"] = b.CloudOrganization.ID
		cloudOrganizatinLinksMap["shared"] = b.Shared
		cloudOrganizatinLinksMap["environment"] = b.Environment

		var resourceTags = make([]interface{}, 0, 0)
		for _, d := range b.ResourceTags {
			var resourceTag = make(map[string]interface{})
			resourceTag["key"] = d.Key
			resourceTag["value"] = d.Value
			resourceTags = append(resourceTags, resourceTag)
		}
		cloudOrganizatinLinksMap["resource_tags"] = resourceTags

		output = append(output, cloudOrganizatinLinksMap)
	}
	return output
}

// ReadProjectPayload struct -- updates
type ReadProjectPayload struct {
	Project vendor.Project `json:"project"`
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
	        description
	        slug
	        archived
	        businessUnit
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
	            shared
	            environment
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

	return diags
}

// UpdateProject struct
type UpdateProject struct {
	UpdateProject vendor.UpdateProjectPayload `json:"updateProject"`
}

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
	vars := &vendor.UpdateProjectInput{}
	vars.ID = d.Id()

	if d.HasChange("archived") {
		vars.Patch.Archived = utils.ConvertBoolToPointer(d.Get("archived").(bool))
	}
	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		vars.Patch.Description = d.Get("description").(string)
	}
	if d.HasChange("business_unit") {
		vars.Patch.BusinessUnit = d.Get("business_unit").(string)
	}
	// The API treats a patch to riskProfile as an override so we set all values
	riskProfile := &vendor.ProjectRiskProfileInput{}
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
	vars.Patch.RiskProfile = riskProfile

	// if cloud organization links are altered, we must send them all org links
	var updateOrgLinks = []*vendor.ProjectCloudOrganizationLinkInput{}
	if d.HasChange("cloud_organization_link") {
		links := d.Get("cloud_organization_link").(*schema.Set).List()
		for _, b := range links {
			var updateOrgLink = &vendor.ProjectCloudOrganizationLinkInput{}
			for c, d := range b.(map[string]interface{}) {
				if c == "environment" {
					updateOrgLink.Environment = d.(string)
				}
				if c == "shared" {
					updateOrgLink.Shared = d.(bool)
				}
				if c == "cloud_organization" {
					updateOrgLink.CloudOrganization = d.(string)
				}
				if c == "resource_tags" {
					var updateResourceTags = []*vendor.ResourceTag{}
					for _, f := range d.(*schema.Set).List() {
						var updateResourceTag = &vendor.ResourceTag{}
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
		vars.Patch.CloudOrganizationLinks = updateOrgLinks
	}

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
  We also change the naem to avoid conflicts since project names must be unique to the org
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
	vars := &vendor.UpdateProjectInput{}
	vars.ID = d.Id()
	vars.Patch.Name = d.Get("slug").(string)
	vars.Patch.Archived = utils.ConvertBoolToPointer(d.Get("archived").(bool))

	// process the request
	data := &UpdateProject{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "project", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
