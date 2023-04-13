package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenScanPolicyParamsIACNoTags(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"builtin_ignore_tags_enabled": false,
			"count_threshold":             3,
			"custom_ignore_tags":          []interface{}{},
			"ignored_rules": []interface{}{
				"fd7dd0c6-4953-4b36-bc39-004ec3d870db",
				"063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
			"security_frameworks": []interface{}{
				"fd7dd0c6-4953-4b36-bc39-004ec3d870db",
				"063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
			"severity_threshold": "CRITICAL",
		},
	}
	var expanded = &wiz.CICDScanPolicyParamsIAC{
		BuiltinIgnoreTagsEnabled: false,
		CountThreshold:           3,
		SeverityThreshold:        "CRITICAL",
		IgnoredRules: []*wiz.CloudConfigurationRule{
			{
				ID: "fd7dd0c6-4953-4b36-bc39-004ec3d870db",
			},
			{
				ID: "063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
		},
		SecurityFrameworks: []*wiz.SecurityFramework{
			{
				ID: "fd7dd0c6-4953-4b36-bc39-004ec3d870db",
			},
			{
				ID: "063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
		},
	}
	scanPolicyParamsIAC := flattenScanPolicyParams(ctx, "CICDScanPolicyParamsIAC", expanded)
	if !reflect.DeepEqual(scanPolicyParamsIAC, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			scanPolicyParamsIAC,
			expected,
		)
	}
}

func TestFlattenScanPolicyParamsIACTags(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"builtin_ignore_tags_enabled": false,
			"count_threshold":             3,
			"custom_ignore_tags": []interface{}{
				map[string]interface{}{
					"ignore_all_rules": false,
					"key":              "testkey1",
					"rule_ids": []interface{}{
						"063fb380-9eda-4c08-a31b-9211ee37bd42",
					},
					"value": "testval1",
				},
				map[string]interface{}{
					"ignore_all_rules": false,
					"key":              "testkey2",
					"rule_ids": []interface{}{
						"1f0ee3b5-5404-4b40-bbc8-33a990330ac3",
						"a1958aa1-b810-4df6-bd82-487cb37c6039",
					},
					"value": "testval2",
				},
				map[string]interface{}{
					"ignore_all_rules": true,
					"key":              "testkey3",
					"value":            "testval3",
					"rule_ids":         []interface{}{},
				},
			},
			"ignored_rules": []interface{}{
				"fd7dd0c6-4953-4b36-bc39-004ec3d870db",
				"063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
			"security_frameworks": []interface{}{
				"fd7dd0c6-4953-4b36-bc39-004ec3d870db",
				"063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
			"severity_threshold": "CRITICAL",
		},
	}
	var expanded = &wiz.CICDScanPolicyParamsIAC{
		BuiltinIgnoreTagsEnabled: false,
		CountThreshold:           3,
		SeverityThreshold:        "CRITICAL",
		IgnoredRules: []*wiz.CloudConfigurationRule{
			{
				ID: "fd7dd0c6-4953-4b36-bc39-004ec3d870db",
			},
			{
				ID: "063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
		},
		SecurityFrameworks: []*wiz.SecurityFramework{
			{
				ID: "fd7dd0c6-4953-4b36-bc39-004ec3d870db",
			},
			{
				ID: "063fb380-9eda-4c08-a31b-9211ee37bd42",
			},
		},
		CustomIgnoreTags: []*wiz.CICDPolicyCustomIgnoreTag{
			{
				IgnoreAllRules: false,
				Key:            "testkey1",
				Value:          "testval1",
				Rules: []*wiz.CloudConfigurationRule{
					{
						ID: "063fb380-9eda-4c08-a31b-9211ee37bd42",
					},
				},
			},
			{
				IgnoreAllRules: false,
				Key:            "testkey2",
				Value:          "testval2",
				Rules: []*wiz.CloudConfigurationRule{
					{
						ID: "1f0ee3b5-5404-4b40-bbc8-33a990330ac3",
					},
					{
						ID: "a1958aa1-b810-4df6-bd82-487cb37c6039",
					},
				},
			},
			{
				IgnoreAllRules: true,
				Key:            "testkey3",
				Value:          "testval3",
			},
		},
	}
	scanPolicyParamsIAC := flattenScanPolicyParams(ctx, "CICDScanPolicyParamsIAC", expanded)
	if !reflect.DeepEqual(scanPolicyParamsIAC, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			scanPolicyParamsIAC,
			expected,
		)
	}
}

func TestFlattenScanPolicyParamsSecrets(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"count_threshold": 3,
			"path_allow_list": []interface{}{
				"/root",
				"/etc",
			},
		},
	}
	var expanded = &wiz.CICDScanPolicyParamsSecrets{
		CountThreshold: 3,
		PathAllowList: []string{
			"/root",
			"/etc",
		},
	}
	scanPolicyParamsSecrets := flattenScanPolicyParams(ctx, "CICDScanPolicyParamsSecrets", expanded)
	if !reflect.DeepEqual(scanPolicyParamsSecrets, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			scanPolicyParamsSecrets,
			expected,
		)
	}
}

func TestFlattenScanPolicyParamsVulnerabilitiesTrue(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"ignore_unfixed": true,
			"package_allow_list": []interface{}{
				"lsof",
				"tcpdump",
			},
			"package_count_threshold": 1,
			"severity":                "HIGH",
		},
	}
	var expanded = &wiz.CICDScanPolicyParamsVulnerabilities{
		IgnoreUnfixed: true,
		PackageAllowList: []string{
			"lsof",
			"tcpdump",
		},
		PackageCountThreshold: 1,
		Severity:              "HIGH",
	}
	scanPolicyParamsVulnerabilities := flattenScanPolicyParams(ctx, "CICDScanPolicyParamsVulnerabilities", expanded)
	if !reflect.DeepEqual(scanPolicyParamsVulnerabilities, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			scanPolicyParamsVulnerabilities,
			expected,
		)
	}
}

func TestFlattenScanPolicyParamsVulnerabilitiesFalse(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"ignore_unfixed": false,
			"package_allow_list": []interface{}{
				"lsof",
				"tcpdump",
			},
			"package_count_threshold": 1,
			"severity":                "HIGH",
		},
	}
	var expanded = &wiz.CICDScanPolicyParamsVulnerabilities{
		IgnoreUnfixed: false,
		PackageAllowList: []string{
			"lsof",
			"tcpdump",
		},
		PackageCountThreshold: 1,
		Severity:              "HIGH",
	}
	scanPolicyParamsVulnerabilities := flattenScanPolicyParams(ctx, "CICDScanPolicyParamsVulnerabilities", expanded)
	if !reflect.DeepEqual(scanPolicyParamsVulnerabilities, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			scanPolicyParamsVulnerabilities,
			expected,
		)
	}
}

func TestGetDiskVulnerabilitiesParams(t *testing.T) {
	ctx := context.Background()

	var expected = &wiz.CreateCICDScanPolicyDiskVulnerabilitiesInput{
		Severity:              "1525fe10-2575-43ef-84bc-6969f81625e7",
		PackageCountThreshold: 3,
		IgnoreUnfixed:         false,
		PackageAllowList: []string{
			"f9de6434-38bc-4da7-b6ea-ff02ad55073f",
			"675a4ecc-71cb-444a-920e-582b06bbadcb",
		},
	}

	d := schema.TestResourceDataRaw(
		t,
		resourceWizCICDScanPolicy().Schema,
		map[string]interface{}{
			"name": "5fdb33cc-5b36-46ee-9d71-b7282d06271a",
			"disk_vulnerabilities_params": []interface{}{
				map[string]interface{}{
					"severity":                "1525fe10-2575-43ef-84bc-6969f81625e7",
					"package_count_threshold": 3,
					"ignore_unfixed":          false,
					"package_allow_list": []interface{}{
						"f9de6434-38bc-4da7-b6ea-ff02ad55073f",
						"675a4ecc-71cb-444a-920e-582b06bbadcb",
					},
				},
			},
		},
	)

	cicdParams := getDiskVulnerabilitiesParams(ctx, d)

	if !reflect.DeepEqual(expected, cicdParams) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			cicdParams,
			expected,
		)
	}
}

func TestGetDiskSecretsParams(t *testing.T) {
	ctx := context.Background()

	var expected = &wiz.CreateCICDScanPolicyDiskSecretsInput{
		CountThreshold: 3,
		PathAllowList: []string{
			"f9de6434-38bc-4da7-b6ea-ff02ad55073f",
			"675a4ecc-71cb-444a-920e-582b06bbadcb",
		},
	}

	d := schema.TestResourceDataRaw(
		t,
		resourceWizCICDScanPolicy().Schema,
		map[string]interface{}{
			"name": "5fdb33cc-5b36-46ee-9d71-b7282d06271a",
			"disk_secrets_params": []interface{}{
				map[string]interface{}{
					"count_threshold": 3,
					"path_allow_list": []interface{}{
						"f9de6434-38bc-4da7-b6ea-ff02ad55073f",
						"675a4ecc-71cb-444a-920e-582b06bbadcb",
					},
				},
			},
		},
	)

	cicdParams := getDiskSecretsParams(ctx, d)

	if !reflect.DeepEqual(expected, cicdParams) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			cicdParams,
			expected,
		)
	}
}

func TestGetIACParams(t *testing.T) {
	ctx := context.Background()

	var expected = &wiz.CreateCICDScanPolicyIACInput{
		SeverityThreshold: "5f45a8d4-24b2-463d-b604-ca532e4ec4d3",
		CountThreshold:    3,
		IgnoredRules: []string{
			"1c1e4a07-8062-4c40-849f-b41417887768",
			"3f25530e-3295-462e-a300-4ef456291263",
		},
		BuiltinIgnoreTagsEnabled: utils.ConvertBoolToPointer(false),
		CustomIgnoreTags: []*wiz.CICDPolicyCustomIgnoreTagCreateInput{
			{
				Key:   "eb9b5425-1635-4cf6-a7b1-44f015795efc",
				Value: "cdebef02-fc13-472e-a4cc-2fe4d355c924",
				RuleIDs: []string{
					"f53784f1-a676-489b-aae6-6672e7005a5f",
					"16eae9f8-b2b7-4cfe-9bff-b828f65d459a",
				},
				IgnoreAllRules: utils.ConvertBoolToPointer(false),
			},
		},
		SecurityFrameworks: []string{
			"5add2652-f417-4050-85de-c1c00c4a6a3c",
			"57fb812b-1220-41c8-b71b-200abbf32c98",
		},
	}

	d := schema.TestResourceDataRaw(
		t,
		resourceWizCICDScanPolicy().Schema,
		map[string]interface{}{
			"name": "5fdb33cc-5b36-46ee-9d71-b7282d06271a",
			"iac_params": []interface{}{
				map[string]interface{}{
					"severity_threshold": "5f45a8d4-24b2-463d-b604-ca532e4ec4d3",
					"count_threshold":    3,
					"ignored_rules": []interface{}{
						"1c1e4a07-8062-4c40-849f-b41417887768",
						"3f25530e-3295-462e-a300-4ef456291263",
					},
					"builtin_ignore_tags_enabled": false,
					"custom_ignore_tags": []interface{}{
						map[string]interface{}{
							"key":   "eb9b5425-1635-4cf6-a7b1-44f015795efc",
							"value": "cdebef02-fc13-472e-a4cc-2fe4d355c924",
							"rule_ids": []interface{}{
								"f53784f1-a676-489b-aae6-6672e7005a5f",
								"16eae9f8-b2b7-4cfe-9bff-b828f65d459a",
							},
							"ignore_all_rules": false,
						},
					},
					"security_frameworks": []interface{}{
						"5add2652-f417-4050-85de-c1c00c4a6a3c",
						"57fb812b-1220-41c8-b71b-200abbf32c98",
					},
				},
			},
		},
	)

	cicdParams := getIACParams(ctx, d)

	if !reflect.DeepEqual(expected, cicdParams) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			cicdParams,
			expected,
		)
	}
}
