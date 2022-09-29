package provider

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func TestFlattenRiskProfileAll(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"business_impact":       "HBI",
			"has_authentication":    "Yes",
			"has_exposed_api":       "Yes",
			"is_actively_developed": "Yes",
			"is_customer_facing":    "Yes",
			"is_internet_facing":    "Yes",
			"is_regulated":          "Yes",
			"regulatory_standards": []interface{}{
				"ISO_27001",
				"ISO_27017",
				"ISO_27018",
				"ISO_27701",
			},
			"sensitive_data_types": []interface{}{
				"HEALTH",
				"PII",
			},
			"stores_data": "Yes",
		},
	}

	var expanded = &vendor.ProjectRiskProfile{
		BusinessImpact:      "HBI",
		HasAuthentication:   "Yes",
		HasExposedAPI:       "Yes",
		IsActivelyDeveloped: "Yes",
		IsCustomerFacing:    "Yes",
		IsInternetFacing:    "Yes",
		IsRegulated:         "Yes",
		RegulatoryStandards: []string{
			"ISO_27001",
			"ISO_27017",
			"ISO_27018",
			"ISO_27701",
		},
		SensitiveDataTypes: []string{
			"HEALTH",
			"PII",
		},
		StoresData: "Yes",
	}

	riskProfile := flattenRiskProfile(ctx, expanded)

	if !reflect.DeepEqual(riskProfile, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			riskProfile,
			expected,
		)
	}
}

func TestFlattenRiskProfileRequired(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"business_impact":       "HBI",
			"has_authentication":    "",
			"has_exposed_api":       "",
			"is_actively_developed": "",
			"is_customer_facing":    "",
			"is_internet_facing":    "",
			"is_regulated":          "",
			"regulatory_standards":  []interface{}{},
			"sensitive_data_types":  []interface{}{},
			"stores_data":           "",
		},
	}

	var expanded = &vendor.ProjectRiskProfile{
		BusinessImpact: "HBI",
	}

	riskProfile := flattenRiskProfile(ctx, expanded)

	if !reflect.DeepEqual(riskProfile, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			riskProfile,
			expected,
		)
	}
}

func TestFlattenRiskProfileDefaults(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"business_impact":       "MBI",
			"has_authentication":    "Unknown",
			"has_exposed_api":       "Unknown",
			"is_actively_developed": "Unknown",
			"is_customer_facing":    "Unknown",
			"is_internet_facing":    "Unknown",
			"is_regulated":          "Unknown",
			"regulatory_standards":  []interface{}{},
			"sensitive_data_types":  []interface{}{},
			"stores_data":           "Unknown",
		},
	}

	var expanded = &vendor.ProjectRiskProfile{
		BusinessImpact:      "MBI",
		HasAuthentication:   "Unknown",
		HasExposedAPI:       "Unknown",
		IsActivelyDeveloped: "Unknown",
		IsCustomerFacing:    "Unknown",
		IsInternetFacing:    "Unknown",
		IsRegulated:         "Unknown",
		RegulatoryStandards: []string{},
		SensitiveDataTypes:  []string{},
		StoresData:          "Unknown",
	}

	riskProfile := flattenRiskProfile(ctx, expanded)

	if !reflect.DeepEqual(riskProfile, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			riskProfile,
			expected,
		)
	}
}

func TestFlattenCloudOrganizationLinksWithTags(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"cloud_organization": "f2b48c0b-57c6-4e1c-9bea-09c92c2fe0ed",
			"shared":             true,
			"environment":        "PRODUCTION",
			"resource_tags": []interface{}{
				map[string]interface{}{
					"key":   "k1",
					"value": "v1",
				},
				map[string]interface{}{
					"key":   "k2",
					"value": "v2",
				},
			},
		},
	}

	var projectCloudOrganizationLink1 = &vendor.ProjectCloudOrganizationLink{
		CloudOrganization: vendor.CloudOrganization{
			ID: "f2b48c0b-57c6-4e1c-9bea-09c92c2fe0ed",
		},
		Shared:      true,
		Environment: "PRODUCTION",
		ResourceTags: []*vendor.ResourceTag{
			{
				Key:   "k1",
				Value: "v1",
			},
			{
				Key:   "k2",
				Value: "v2",
			},
		},
	}

	expanded := []*vendor.ProjectCloudOrganizationLink{}
	expanded = append(expanded, projectCloudOrganizationLink1)

	cloudOrganizationLinks := flattenCloudOrganizationLinks(ctx, expanded)

	if !reflect.DeepEqual(cloudOrganizationLinks, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			cloudOrganizationLinks,
			expected,
		)
	}
}

func TestFlattenCloudOrganizationLinksNoTags(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"cloud_organization": "f2b48c0b-57c6-4e1c-9bea-09c92c2fe0ed",
			"shared":             true,
			"environment":        "PRODUCTION",
			"resource_tags":      []interface{}{},
		},
	}

	var projectCloudOrganizationLink1 = &vendor.ProjectCloudOrganizationLink{
		CloudOrganization: vendor.CloudOrganization{
			ID: "f2b48c0b-57c6-4e1c-9bea-09c92c2fe0ed",
		},
		Shared:       true,
		Environment:  "PRODUCTION",
		ResourceTags: []*vendor.ResourceTag{},
	}

	expanded := []*vendor.ProjectCloudOrganizationLink{}
	expanded = append(expanded, projectCloudOrganizationLink1)

	cloudOrganizationLinks := flattenCloudOrganizationLinks(ctx, expanded)

	if !reflect.DeepEqual(cloudOrganizationLinks, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			cloudOrganizationLinks,
			expected,
		)
	}
}

func TestGetOrganizationLinksVar(t *testing.T) {
	ctx := context.Background()

	var expectedOrgLink1 = &vendor.ProjectCloudOrganizationLinkInput{
		CloudOrganization: "98a51e0c-7ab9-4f40-b8d9-a4fc398ad98a",
		Environment:       "PRODUCTION",
		ResourceTags: []*vendor.ResourceTag{
			{
				Key:   "7982c5c6-1c66-435c-a509-68fae7718bd8",
				Value: "fbf63c90-67ed-4198-af07-05ee17a58c1d",
			},
		},
		Shared: true,
	}

	var expectedOrgLink2 = &vendor.ProjectCloudOrganizationLinkInput{
		CloudOrganization: "d8181cf9-38bb-486c-8278-f95f416afb3c",
		Environment:       "PRODUCTION",
		Shared:            false,
	}

	var expected = []*vendor.ProjectCloudOrganizationLinkInput{}
	expected = append(expected, expectedOrgLink1)
	expected = append(expected, expectedOrgLink2)

	d := schema.TestResourceDataRaw(
		t,
		resourceWizProject().Schema,
		map[string]interface{}{
			"name": "70bbbb01-6438-4e91-82d9-e1d46e7795f8",
			"cloud_organization_link": []interface{}{
				map[string]interface{}{
					"cloud_organization": "98a51e0c-7ab9-4f40-b8d9-a4fc398ad98a",
					"environment":        "PRODUCTION",
					"shared":             true,
					"resource_tags": []interface{}{
						map[string]interface{}{
							"key":   "7982c5c6-1c66-435c-a509-68fae7718bd8",
							"value": "fbf63c90-67ed-4198-af07-05ee17a58c1d",
						},
					},
				},
				map[string]interface{}{
					"cloud_organization": "d8181cf9-38bb-486c-8278-f95f416afb3c",
					"environment":        "PRODUCTION",
					"shared":             false,
				},
			},
		},
	)

	orgLink := getOrganizationLinksVar(ctx, d)

	sort.SliceStable(expected, func(i, j int) bool { return expected[i].CloudOrganization < expected[j].CloudOrganization })
	sort.SliceStable(orgLink, func(i, j int) bool { return orgLink[i].CloudOrganization < orgLink[j].CloudOrganization })

	if !reflect.DeepEqual(expected, orgLink) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(orgLink),
			utils.PrettyPrint(expected),
		)
	}
}
