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

func TestFlattenGroupMapping(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"projects": []interface{}{
				"cb95ced6-3ed6-5fd5-a68a-1059556fc909",
				"69229a09-f831-484e-9d3c-21f2a984a014",
			},
			"provider_group_id": "Wiz-Project-Reader",
			"role":              "PROJECT_READER",
		},
		map[string]interface{}{
			"projects": []interface{}{
				"69229a09-f831-484e-9d3c-21f2a984a014",
			},
			"provider_group_id": "Wiz-Project-Admin",
			"role":              "PROJECT_ADMIN",
		},
		map[string]interface{}{
			"projects":          []interface{}{},
			"provider_group_id": "Wiz-Global-Admin",
			"role":              "GLOBAL_ADMIN",
		},
	}

	// verify multiple projects
	var groupMapping1 = &wiz.SAMLGroupMapping{
		ProviderGroupID: "Wiz-Project-Reader",
		Role: wiz.UserRole{
			ID: "PROJECT_READER",
		},
		Projects: []wiz.Project{
			{
				ID: "cb95ced6-3ed6-5fd5-a68a-1059556fc909",
			},
			{
				ID: "69229a09-f831-484e-9d3c-21f2a984a014",
			},
		},
	}

	// verify single project
	var groupMapping2 = &wiz.SAMLGroupMapping{
		ProviderGroupID: "Wiz-Project-Admin",
		Role: wiz.UserRole{
			ID: "PROJECT_ADMIN",
		},
		Projects: []wiz.Project{
			{
				ID: "69229a09-f831-484e-9d3c-21f2a984a014",
			},
		},
	}

	// verify no projects
	var groupMapping3 = &wiz.SAMLGroupMapping{
		ProviderGroupID: "Wiz-Global-Admin",
		Role: wiz.UserRole{
			ID: "GLOBAL_ADMIN",
		},
	}

	expanded := []*wiz.SAMLGroupMapping{}
	expanded = append(expanded, groupMapping1)
	expanded = append(expanded, groupMapping2)
	expanded = append(expanded, groupMapping3)

	groupMapping := flattenGroupMapping(ctx, expanded)

	if !reflect.DeepEqual(groupMapping, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			groupMapping,
			expected,
		)
	}
}

func TestGetGroupMappingVar(t *testing.T) {
	ctx := context.Background()

	var expectedMember1 = &wiz.SAMLGroupMappingCreateInput{
		ProviderGroupID: "f11fd4a4-ba73-448d-9894-8dbd4c94f48b",
		Role:            "bd3545c2-9d3e-4c43-ab55-c9fc982d8ae4",
		Projects: []string{
			"e228be4f-1697-4c02-ab8c-7d3c526cb22c",
			"f48d4e70-7028-4f5f-8f30-a77e139f8d38",
		},
	}

	var expectedMember2 = &wiz.SAMLGroupMappingCreateInput{
		ProviderGroupID: "5591d307-49ec-41f4-acfc-c2295dc90c94",
		Role:            "b32da743-4725-4ff6-b1e5-8e330b9f0080",
	}

	var expected = []*wiz.SAMLGroupMappingCreateInput{}

	expected = append(expected, expectedMember1)
	expected = append(expected, expectedMember2)

	d := schema.TestResourceDataRaw(
		t,
		resourceWizSAMLIdP().Schema,
		map[string]interface{}{
			"name":                       "70bbbb01-6438-4e91-82d9-e1d46e7795f8",
			"login_url":                  "https://example.com",
			"use_provider_managed_roles": true,
			"allow_manual_role_override": true,
			"certificate":                "7949a0d0-bb64-43e1-9af7-1c0ee0574f7a",
			"domains": []interface{}{
				"example.com",
			},
			"group_mapping": []interface{}{
				map[string]interface{}{
					"provider_group_id": "f11fd4a4-ba73-448d-9894-8dbd4c94f48b",
					"role":              "bd3545c2-9d3e-4c43-ab55-c9fc982d8ae4",
					"projects": []interface{}{
						"e228be4f-1697-4c02-ab8c-7d3c526cb22c",
						"f48d4e70-7028-4f5f-8f30-a77e139f8d38",
					},
				},
				map[string]interface{}{
					"provider_group_id": "5591d307-49ec-41f4-acfc-c2295dc90c94",
					"role":              "b32da743-4725-4ff6-b1e5-8e330b9f0080",
				},
			},
		},
	)

	groupMapping := getGroupMappingVar(ctx, d)

	sort.SliceStable(expected, func(i, j int) bool { return expected[i].Role < expected[j].Role })
	sort.SliceStable(groupMapping, func(i, j int) bool { return groupMapping[i].Role < groupMapping[j].Role })

	if !reflect.DeepEqual(expected, groupMapping) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(groupMapping),
			utils.PrettyPrint(expected),
		)
	}
}
