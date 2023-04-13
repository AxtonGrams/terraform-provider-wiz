package provider

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
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

	var expanded = &wiz.ProjectRiskProfile{
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

	var expanded = &wiz.ProjectRiskProfile{
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

	var expanded = &wiz.ProjectRiskProfile{
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

func TestFlattenKubernetesClusterLinks(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"environment":        "STAGING",
			"kubernetes_cluster": "77de7ca1-02f9-5ed2-a94b-5d19c683efa2",
			"shared":             true,
			"namespaces": []interface{}{
				"kube-system",
			},
		},
	}
	var clusterlink = &wiz.ProjectKubernetesClusterLink{
		Environment: "STAGING",
		KubernetesCluster: wiz.KubernetesCluster{
			ID: "77de7ca1-02f9-5ed2-a94b-5d19c683efa2",
		},
		Shared: true,
		Namespaces: []string{
			"kube-system",
		},
	}

	expanded := []*wiz.ProjectKubernetesClusterLink{}
	expanded = append(expanded, clusterlink)

	var flattened = flattenKubernetesClusterLinks(ctx, expanded)
	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			expanded,
			expected,
		)
	}
}

func TestFlattenCloudAccountLinksWithTags(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"cloud_account_id": "3225def3-0e0e-5cb8-955a-3583f696f77f",
			"shared":           true,
			"environment":      "PRODUCTION",
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
			"resource_groups": []interface{}{
				"v1",
				"v2",
			},
		},
	}

	var projectCloudAccountLink1 = &wiz.ProjectCloudAccountLink{
		CloudAccount: wiz.CloudAccount{
			ID: "3225def3-0e0e-5cb8-955a-3583f696f77f",
		},
		Shared:      true,
		Environment: "PRODUCTION",
		ResourceTags: []*wiz.ResourceTag{
			{
				Key:   "k1",
				Value: "v1",
			},
			{
				Key:   "k2",
				Value: "v2",
			},
		},
		ResourceGroups: []string{
			"v1",
			"v2",
		},
	}

	expanded := []*wiz.ProjectCloudAccountLink{}
	expanded = append(expanded, projectCloudAccountLink1)

	cloudAccountLinks := flattenCloudAccountLinks(ctx, expanded)

	if !reflect.DeepEqual(cloudAccountLinks, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			cloudAccountLinks,
			expected,
		)
	}
}

func TestFlattenCloudAccountLinksNoTags(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"cloud_account_id": "3225def3-0e0e-5cb8-955a-3583f696f77f",
			"shared":           true,
			"environment":      "PRODUCTION",
			"resource_tags":    []interface{}{},
			"resource_groups":  []interface{}{},
		},
	}

	var projectCloudAccountLink1 = &wiz.ProjectCloudAccountLink{
		CloudAccount: wiz.CloudAccount{
			ID: "3225def3-0e0e-5cb8-955a-3583f696f77f",
		},
		Shared:         true,
		Environment:    "PRODUCTION",
		ResourceTags:   []*wiz.ResourceTag{},
		ResourceGroups: []string{},
	}

	expanded := []*wiz.ProjectCloudAccountLink{}
	expanded = append(expanded, projectCloudAccountLink1)

	cloudAccountLinks := flattenCloudAccountLinks(ctx, expanded)

	if !reflect.DeepEqual(cloudAccountLinks, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			cloudAccountLinks,
			expected,
		)
	}
}

func TestGetCloudAccountLinksVar(t *testing.T) {
	ctx := context.Background()

	var expectedAcclink1 = &wiz.ProjectCloudAccountLinkInput{
		CloudAccount: "3225def3-0e0e-5cb8-955a-3583f696f77f",
		Environment:  "PRODUCTION",
		Shared:       utils.ConvertBoolToPointer(true),
		ResourceTags: []*wiz.ResourceTagInput{
			{
				Key:   "7982c5c6-1c66-435c-a509-68fae7718bd8",
				Value: "fbf63c90-67ed-4198-af07-05ee17a58c1d",
			},
		},
		ResourceGroups: []string{
			"v1",
			"v2",
		},
	}

	var expectedAcclink2 = &wiz.ProjectCloudAccountLinkInput{
		CloudAccount: "d8181cf9-38bb-486c-8278-f95f416afb3c",
		Environment:  "PRODUCTION",
		Shared:       utils.ConvertBoolToPointer(false),
	}

	var expected = []*wiz.ProjectCloudAccountLinkInput{}
	expected = append(expected, expectedAcclink1)
	expected = append(expected, expectedAcclink2)

	d := schema.TestResourceDataRaw(
		t,
		resourceWizProject().Schema,
		map[string]interface{}{
			"name": "my project",
			"cloud_account_link": []interface{}{
				map[string]interface{}{
					"cloud_account_id": "3225def3-0e0e-5cb8-955a-3583f696f77f",
					"environment":      "PRODUCTION",
					"shared":           true,
					"resource_tags": []interface{}{
						map[string]interface{}{
							"key":   "7982c5c6-1c66-435c-a509-68fae7718bd8",
							"value": "fbf63c90-67ed-4198-af07-05ee17a58c1d",
						},
					},
					"resource_groups": []interface{}{
						"v1",
						"v2",
					},
				},
				map[string]interface{}{
					"cloud_account_id": "d8181cf9-38bb-486c-8278-f95f416afb3c",
					"environment":      "PRODUCTION",
					"shared":           false,
				},
			},
		},
	)

	accLink := getAccountLinksVar(ctx, d)

	sort.SliceStable(expected, func(i, j int) bool { return expected[i].CloudAccount < expected[j].CloudAccount })
	sort.SliceStable(accLink, func(i, j int) bool { return accLink[i].CloudAccount < accLink[j].CloudAccount })

	if !reflect.DeepEqual(expected, accLink) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(accLink),
			utils.PrettyPrint(expected),
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
			"resource_groups": []interface{}{
				"v1",
				"v2",
			},
		},
	}

	var projectCloudOrganizationLink1 = &wiz.ProjectCloudOrganizationLink{
		CloudOrganization: wiz.CloudOrganization{
			ID: "f2b48c0b-57c6-4e1c-9bea-09c92c2fe0ed",
		},
		Shared:      true,
		Environment: "PRODUCTION",
		ResourceTags: []*wiz.ResourceTag{
			{
				Key:   "k1",
				Value: "v1",
			},
			{
				Key:   "k2",
				Value: "v2",
			},
		},
		ResourceGroups: []string{
			"v1",
			"v2",
		},
	}

	expanded := []*wiz.ProjectCloudOrganizationLink{}
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

func TestFlattenUserIds(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"01882691-fb1b-5e72-b4cf-dc207a875907",
	}
	var userID = &wiz.User{
		ID: "01882691-fb1b-5e72-b4cf-dc207a875907",
	}

	expanded := []*wiz.User{}
	expanded = append(expanded, userID)

	userIds := flattenUserIds(ctx, expanded)
	if !reflect.DeepEqual(userIds, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			userIds,
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
			"resource_groups":    []interface{}{},
		},
	}

	var projectCloudOrganizationLink1 = &wiz.ProjectCloudOrganizationLink{
		CloudOrganization: wiz.CloudOrganization{
			ID: "f2b48c0b-57c6-4e1c-9bea-09c92c2fe0ed",
		},
		Shared:       true,
		Environment:  "PRODUCTION",
		ResourceTags: []*wiz.ResourceTag{},
	}

	expanded := []*wiz.ProjectCloudOrganizationLink{}
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

	var expectedOrgLink1 = &wiz.ProjectCloudOrganizationLinkInput{
		CloudOrganization: "98a51e0c-7ab9-4f40-b8d9-a4fc398ad98a",
		Environment:       "PRODUCTION",
		ResourceTags: []*wiz.ResourceTagInput{
			{
				Key:   "7982c5c6-1c66-435c-a509-68fae7718bd8",
				Value: "fbf63c90-67ed-4198-af07-05ee17a58c1d",
			},
		},
		Shared: true,
	}

	var expectedOrgLink2 = &wiz.ProjectCloudOrganizationLinkInput{
		CloudOrganization: "d8181cf9-38bb-486c-8278-f95f416afb3c",
		Environment:       "PRODUCTION",
		Shared:            false,
	}

	var expected = []*wiz.ProjectCloudOrganizationLinkInput{}
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
