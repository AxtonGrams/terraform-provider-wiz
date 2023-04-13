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

func TestGetIACMatchers(t *testing.T) {
	ctx := context.Background()

	var expected1 = &wiz.CreateCloudConfigurationRuleMatcherInput{
		Type:     "668e943b-84d4-4bcd-b255-58b4bc6f8587",
		RegoCode: "060231e4-31a5-47a0-a1ec-2fe166c939fd",
	}

	var expected2 = &wiz.CreateCloudConfigurationRuleMatcherInput{
		Type:     "c92208fc-c9c5-45fe-9f09-6aa074bc0fd0",
		RegoCode: "7f8b8994-aa4d-42ca-bda2-159a8602d4e0",
	}

	var expected = []*wiz.CreateCloudConfigurationRuleMatcherInput{}
	expected = append(expected, expected1)
	expected = append(expected, expected2)

	d := schema.TestResourceDataRaw(
		t,
		resourceWizCloudConfigurationRule().Schema,
		map[string]interface{}{
			"name":                "d5590810-2a09-4986-b63b-ab0f993a3c34",
			"target_native_types": "801fcd6a-c603-4275-8ed3-4206cc1508d7",
			"iac_matchers": []interface{}{
				map[string]interface{}{
					"type":      "668e943b-84d4-4bcd-b255-58b4bc6f8587",
					"rego_code": "060231e4-31a5-47a0-a1ec-2fe166c939fd",
				},
				map[string]interface{}{
					"type":      "c92208fc-c9c5-45fe-9f09-6aa074bc0fd0",
					"rego_code": "7f8b8994-aa4d-42ca-bda2-159a8602d4e0",
				},
			},
		},
	)

	iacMatchers := getIACMatchers(ctx, d)

	sort.SliceStable(expected, func(i, j int) bool { return expected[i].RegoCode < expected[j].RegoCode })
	sort.SliceStable(iacMatchers, func(i, j int) bool { return iacMatchers[i].RegoCode < iacMatchers[j].RegoCode })

	if !reflect.DeepEqual(expected, iacMatchers) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(iacMatchers),
			utils.PrettyPrint(expected),
		)
	}
}

func TestFlattenSecuritySubCategoriesID(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"c14e688d-62ca-42f0-bf27-01195288d4ec",
		"09f38436-80ac-493a-b554-e000e4f9cbd1",
		"70f8c07e-aa05-4354-b589-e67605c72d44",
	}
	var expanded1 = &wiz.SecuritySubCategory{
		ID: "c14e688d-62ca-42f0-bf27-01195288d4ec",
	}
	var expanded2 = &wiz.SecuritySubCategory{
		ID: "09f38436-80ac-493a-b554-e000e4f9cbd1",
	}
	var expanded3 = &wiz.SecuritySubCategory{
		ID: "70f8c07e-aa05-4354-b589-e67605c72d44",
	}
	var expanded = []*wiz.SecuritySubCategory{}
	expanded = append(expanded, expanded1)
	expanded = append(expanded, expanded2)
	expanded = append(expanded, expanded3)
	securitySubCategories := flattenSecuritySubCategoriesID(ctx, expanded)
	sort.SliceStable(expected, func(i, j int) bool { return expected[i].(string) < expected[j].(string) })
	sort.SliceStable(securitySubCategories, func(i, j int) bool { return securitySubCategories[i].(string) < securitySubCategories[j].(string) })
	if !reflect.DeepEqual(securitySubCategories, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			securitySubCategories,
			expected,
		)
	}
}

func TestFattenScopeAccountIDs(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"c14e688d-62ca-42f0-bf27-01195288d4ec",
		"09f38436-80ac-493a-b554-e000e4f9cbd1",
		"70f8c07e-aa05-4354-b589-e67605c72d44",
	}
	var expanded1 = &wiz.CloudAccount{
		ID: "c14e688d-62ca-42f0-bf27-01195288d4ec",
	}
	var expanded2 = &wiz.CloudAccount{
		ID: "09f38436-80ac-493a-b554-e000e4f9cbd1",
	}
	var expanded3 = &wiz.CloudAccount{
		ID: "70f8c07e-aa05-4354-b589-e67605c72d44",
	}
	var expanded = []*wiz.CloudAccount{}
	expanded = append(expanded, expanded1)
	expanded = append(expanded, expanded2)
	expanded = append(expanded, expanded3)
	scopeAccountIDs := flattenScopeAccountIDs(ctx, expanded)
	sort.SliceStable(expected, func(i, j int) bool { return expected[i].(string) < expected[j].(string) })
	sort.SliceStable(scopeAccountIDs, func(i, j int) bool { return scopeAccountIDs[i].(string) < scopeAccountIDs[j].(string) })
	if !reflect.DeepEqual(scopeAccountIDs, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			scopeAccountIDs,
			expected,
		)
	}
}

func TestFlattenIACMatchers(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"type":      "484e96e0-8bac-4b86-ab9f-e43dd950ec05",
			"rego_code": "6f8344d7-6252-4672-980f-86ed9f19b883",
		},
		map[string]interface{}{
			"type":      "87a7f4a1-a789-4fae-9bbc-3fe6d8f16218",
			"rego_code": "c1ebfd5c-a483-4559-95e7-ce6b87f9a0fe",
		},
		map[string]interface{}{
			"type":      "d53bafe7-d821-4bbd-86f9-9634bb6f9b14",
			"rego_code": "fa637c75-889a-4e85-9017-5cb16b36cc7c",
		},
	}
	var expanded1 = &wiz.CloudConfigurationRuleMatcher{
		Type:     "484e96e0-8bac-4b86-ab9f-e43dd950ec05",
		RegoCode: "6f8344d7-6252-4672-980f-86ed9f19b883",
	}
	var expanded2 = &wiz.CloudConfigurationRuleMatcher{
		Type:     "87a7f4a1-a789-4fae-9bbc-3fe6d8f16218",
		RegoCode: "c1ebfd5c-a483-4559-95e7-ce6b87f9a0fe",
	}
	var expanded3 = &wiz.CloudConfigurationRuleMatcher{
		Type:     "d53bafe7-d821-4bbd-86f9-9634bb6f9b14",
		RegoCode: "fa637c75-889a-4e85-9017-5cb16b36cc7c",
	}
	var expanded = []*wiz.CloudConfigurationRuleMatcher{}
	expanded = append(expanded, expanded1)
	expanded = append(expanded, expanded2)
	expanded = append(expanded, expanded3)
	iacMatchers := flattenIACMatchers(ctx, expanded)
	if !reflect.DeepEqual(iacMatchers, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			iacMatchers,
			expected,
		)
	}
}
