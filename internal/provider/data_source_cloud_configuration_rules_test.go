package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenCloudConfigurationRules(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"id":          "efc0d7e9-6a84-45f6-869c-529e337db6a8",
			"name":        "1d401010-0e50-4101-93c7-35b7f9226997",
			"short_id":    "ce84a427-7358-4fb8-853e-29fbdd4f3f55",
			"description": "27669d80-612d-4e25-9c85-369625e2e990",
			"enabled":     false,
			"severity":    "565a7bb8-9f08-4e4f-96e0-da5af64e446a",
			"target_native_types": []string{
				"a3f18510-600c-42a0-854a-117e83965f31",
				"5bcf4cca-8bc1-46bf-8dbc-404486397885",
				"d88ccf3b-3c58-4209-b2cb-57291512b093",
			},
			"supports_nrt":             false,
			"subject_entity_type":      "c53b5b06-5565-48e4-9ad7-2ecdb9596681",
			"cloud_provider":           "62c7de2d-d884-4b68-ae8b-f920ed1d60ba",
			"service_type":             "d65d3e25-36ed-45db-bf6b-50dc91c2bfc3",
			"builtin":                  false,
			"opa_policy":               "49d5692f-7681-435d-8eb5-e5f22758ffb2",
			"function_as_control":      false,
			"control_id":               "ef708a85-0c78-4a50-b5eb-f44501dc5fc2",
			"graph_id":                 "a81dc546-e3e3-4fbc-83c4-0ed5edbfa0c9",
			"has_auto_remediation":     false,
			"remediation_instructions": "eb3e8071-49dd-4679-90d6-46bdf4387ee8",
			"external_references": []interface{}{
				map[string]interface{}{
					"id":   "0bbe2b15-31f1-4786-b2b9-c57704c10fbf",
					"name": "7afbbec4-76c6-40e2-aeda-24e89f616df0",
				},
				map[string]interface{}{
					"id":   "43960b7b-09b3-4db7-90e7-7ad471071ba2",
					"name": "0bb079f6-d315-4432-9192-d64a30c2cf0f",
				},
				map[string]interface{}{
					"id":   "f60d29dd-a16f-4fc0-949b-981b95dcdc16",
					"name": "8feb5283-1e99-4ba5-b963-1ad631632be3",
				},
			},
			"scope_accounts": []interface{}{
				"42cecee6-45ec-41d0-8d66-b0571e2b6f62",
				"868eb547-5b76-4b5e-b9a1-9aecec79846d",
				"a9caf002-b075-43e9-b1cb-9cce1c809083",
			},
			"security_sub_category_ids": []interface{}{
				"26b40cb3-5730-4097-a61f-c12ed69970bc",
				"49817d03-5a33-4c52-b7bc-51d800d173de",
				"83abe403-dc4d-43a1-b5e5-8a6d4b9f5182",
			},
			"iac_matcher_ids": []interface{}{
				"4256c551-a971-4fd6-948c-d883c4964fae",
				"4d5e2d5c-1157-487b-b9c1-7aa10bc83f83",
				"56344421-3a12-42da-9c83-8377e3420733",
			},
		},
	}

	var configRules = &[]*wiz.CloudConfigurationRule{
		{
			Builtin:       utils.ConvertBoolToPointer(false),
			CloudProvider: "62c7de2d-d884-4b68-ae8b-f920ed1d60ba",
			Control: &wiz.Control{
				ID: "ef708a85-0c78-4a50-b5eb-f44501dc5fc2",
			},
			CreatedBy: &wiz.User{
				ID: "7aa9f806-18b2-48fd-83d1-2c151d462ea7",
			},
			Description: "27669d80-612d-4e25-9c85-369625e2e990",
			Enabled:     utils.ConvertBoolToPointer(false),
			ExternalReferences: []*wiz.CloudConfigurationRuleExternalReference{
				{
					ID:   "0bbe2b15-31f1-4786-b2b9-c57704c10fbf",
					Name: "7afbbec4-76c6-40e2-aeda-24e89f616df0",
				},
				{
					ID:   "f60d29dd-a16f-4fc0-949b-981b95dcdc16",
					Name: "8feb5283-1e99-4ba5-b963-1ad631632be3",
				},
				{
					ID:   "43960b7b-09b3-4db7-90e7-7ad471071ba2",
					Name: "0bb079f6-d315-4432-9192-d64a30c2cf0f",
				},
			},
			FunctionAsControl:  utils.ConvertBoolToPointer(false),
			GraphID:            "a81dc546-e3e3-4fbc-83c4-0ed5edbfa0c9",
			HasAutoRemediation: utils.ConvertBoolToPointer(false),
			IACMatchers: []*wiz.CloudConfigurationRuleMatcher{
				{
					ID: "4d5e2d5c-1157-487b-b9c1-7aa10bc83f83",
				},
				{
					ID: "4256c551-a971-4fd6-948c-d883c4964fae",
				},
				{
					ID: "56344421-3a12-42da-9c83-8377e3420733",
				},
			},
			ID:                      "efc0d7e9-6a84-45f6-869c-529e337db6a8",
			Name:                    "1d401010-0e50-4101-93c7-35b7f9226997",
			OPAPolicy:               "49d5692f-7681-435d-8eb5-e5f22758ffb2",
			RemediationInstructions: "eb3e8071-49dd-4679-90d6-46bdf4387ee8",
			ScopeAccounts: []*wiz.CloudAccount{
				{
					ID: "42cecee6-45ec-41d0-8d66-b0571e2b6f62",
				},
				{
					ID: "a9caf002-b075-43e9-b1cb-9cce1c809083",
				},
				{
					ID: "868eb547-5b76-4b5e-b9a1-9aecec79846d",
				},
			},
			SecuritySubCategories: []*wiz.SecuritySubCategory{
				{
					ID: "49817d03-5a33-4c52-b7bc-51d800d173de",
				},
				{
					ID: "83abe403-dc4d-43a1-b5e5-8a6d4b9f5182",
				},
				{
					ID: "26b40cb3-5730-4097-a61f-c12ed69970bc",
				},
			},
			ServiceType:       "d65d3e25-36ed-45db-bf6b-50dc91c2bfc3",
			Severity:          "565a7bb8-9f08-4e4f-96e0-da5af64e446a",
			ShortID:           "ce84a427-7358-4fb8-853e-29fbdd4f3f55",
			SubjectEntityType: "c53b5b06-5565-48e4-9ad7-2ecdb9596681",
			SupportsNRT:       utils.ConvertBoolToPointer(false),
			TargetNativeTypes: []string{
				"a3f18510-600c-42a0-854a-117e83965f31",
				"5bcf4cca-8bc1-46bf-8dbc-404486397885",
				"d88ccf3b-3c58-4209-b2cb-57291512b093",
			},
		},
	}

	flattened := flattenCloudConfigurationRules(ctx, configRules)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}

func TestFlattenExternalReferences(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"id":   "39475a03-cba0-4aa3-867d-41494b9c4ad7",
			"name": "7f5fa936-402f-4885-8b79-874ce3db8fb2",
		},
		map[string]interface{}{
			"id":   "6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
			"name": "284aab9b-88a5-4e13-a182-212de84989af",
		},
		map[string]interface{}{
			"id":   "9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
			"name": "b9f72ab5-ac09-4f5b-a57f-6691c55bcbdb",
		},
	}

	var refs = &[]*wiz.CloudConfigurationRuleExternalReference{
		{
			ID:   "6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
			Name: "284aab9b-88a5-4e13-a182-212de84989af",
		},
		{
			ID:   "9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
			Name: "b9f72ab5-ac09-4f5b-a57f-6691c55bcbdb",
		},
		{
			ID:   "39475a03-cba0-4aa3-867d-41494b9c4ad7",
			Name: "7f5fa936-402f-4885-8b79-874ce3db8fb2",
		},
	}

	flattened := flattenExternalReferences(ctx, refs)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}

func TestFlattenScopeAccounts(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"39475a03-cba0-4aa3-867d-41494b9c4ad7",
		"6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
		"9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
	}

	var accounts = &[]*wiz.CloudAccount{
		{
			ID: "6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
		},
		{
			ID: "9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
		},
		{
			ID: "39475a03-cba0-4aa3-867d-41494b9c4ad7",
		},
	}

	flattened := flattenScopeAccounts(ctx, accounts)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}

func TestFlattenSecuritySubCategoryIDs(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"39475a03-cba0-4aa3-867d-41494b9c4ad7",
		"6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
		"9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
	}

	var cats = &[]*wiz.SecuritySubCategory{
		{
			ID: "6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
		},
		{
			ID: "9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
		},
		{
			ID: "39475a03-cba0-4aa3-867d-41494b9c4ad7",
		},
	}

	flattened := flattenSecuritySubCategoryIDs(ctx, cats)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}

func TestFlattenIACMatcherIDs(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"39475a03-cba0-4aa3-867d-41494b9c4ad7",
		"6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
		"9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
	}

	var matchers = &[]*wiz.CloudConfigurationRuleMatcher{
		{
			ID: "6362a2a5-a7fc-48fd-9f38-f5a1e9c560d8",
		},
		{
			ID: "9c5d545a-bcb2-41ac-89f9-7dcca90483b7",
		},
		{
			ID: "39475a03-cba0-4aa3-867d-41494b9c4ad7",
		},
	}

	flattened := flattenIACMatcherIDs(ctx, matchers)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}
