package provider

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"testing"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

const (
	environmentProd = "PRODUCTION"
	environmentDev  = "DEVELOPMENT"
	tag1            = "tag1"
	value1          = "value1"
)

var cloudAccountId = uuid.NewString()
var resourceGroups = []string{"group1", "group2"}
var resourceGroupsInterface = []interface{}{"group1", "group2"}

func TestGetAccountLinkVar(t *testing.T) {
	expected := &wiz.ProjectCloudAccountLinkInput{
		CloudAccount:   cloudAccountId,
		Environment:    environmentProd,
		Shared:         utils.ConvertBoolToPointer(true),
		ResourceGroups: resourceGroups,
		ResourceTags: []*wiz.ResourceTagInput{
			{
				Key:   tag1,
				Value: value1,
			},
		},
	}

	d := schema.TestResourceDataRaw(
		t,
		resourceWizProjectCloudAccountLink().Schema,
		map[string]interface{}{
			"cloud_account_id": cloudAccountId,
			"environment":      environmentProd,
			"shared":           true,
			"resource_groups":  resourceGroupsInterface,
			"resource_tags": []interface{}{
				map[string]interface{}{
					"key":   tag1,
					"value": value1,
				},
			},
		},
	)

	accountLink := getAccountLinkVar(d, cloudAccountId)
	if !reflect.DeepEqual(expected, accountLink) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(accountLink),
			utils.PrettyPrint(expected),
		)
	}
}

func TestAccountLinkToAccountLinkInput(t *testing.T) {
	link := &wiz.ProjectCloudAccountLink{
		CloudAccount: wiz.CloudAccount{
			ID: cloudAccountId,
		},
		Environment: environmentProd,
		ResourceTags: []*wiz.ResourceTag{
			{
				Key:   tag1,
				Value: value1,
			},
		},
		Shared: true,
	}

	expected := &wiz.ProjectCloudAccountLinkInput{
		CloudAccount: cloudAccountId,
		Environment:  environmentProd,
		ResourceTags: []*wiz.ResourceTagInput{
			{
				Key:   tag1,
				Value: value1,
			},
		},
		Shared: utils.ConvertBoolToPointer(true),
	}

	result := accountLinkToAccountLinkInput(link)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(result),
			utils.PrettyPrint(expected),
		)
	}
}

func TestExtractCloudAccountLink(t *testing.T) {
	cloudAccountLinks := []*wiz.ProjectCloudAccountLink{
		{
			CloudAccount: wiz.CloudAccount{
				ID: cloudAccountId,
			},
			Environment: environmentProd,
			Shared:      true,
		},
		{
			CloudAccount: wiz.CloudAccount{
				ID: "other-id",
			},
			Environment: environmentDev,
			Shared:      false,
		},
	}

	expected := &wiz.ProjectCloudAccountLink{
		CloudAccount: wiz.CloudAccount{
			ID: cloudAccountId,
		},
		Environment: environmentProd,
		Shared:      true,
	}

	result, err := extractCloudAccountLink(cloudAccountLinks, cloudAccountId)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			utils.PrettyPrint(result),
			utils.PrettyPrint(expected),
		)
	}
}

func TestExtractIds(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedProj  string
		expectedCloud string
		expectErr     bool
	}{
		{
			name:          "Valid ID",
			input:         "link|projectId|cloudAccountUpstreamId",
			expectedProj:  "projectId",
			expectedCloud: "cloudAccountUpstreamId",
			expectErr:     false,
		},
		{
			name:          "Invalid ID",
			input:         "invalidId",
			expectedProj:  "",
			expectedCloud: "",
			expectErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projId, cloudId, err := extractIds(tc.input)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}
			if projId != tc.expectedProj {
				t.Errorf("Expected project ID: %s, got: %s", tc.expectedProj, projId)
			}
			if cloudId != tc.expectedCloud {
				t.Errorf("Expected cloud ID: %s, got: %s", tc.expectedCloud, cloudId)
			}
		})
	}
}
