package provider

import (
	"reflect"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"

	"testing"
)

func TestExtractIDsFromSamlIdpGroupMappingImportID(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		expectedMapping SAMLGroupMappingsImport
		expectErr       bool
	}{
		{
			name:            "Valid ID",
			input:           "link|samlIdpID|providerGroupID:role:projectID1,projectID2",
			expectedMapping: SAMLGroupMappingsImport{SamlIdpID: "samlIdpID", GroupMappings: []wiz.SAMLGroupDetailsInput{{ProviderGroupID: "providerGroupID", Role: "role", Projects: []string{"projectID1", "projectID2"}}}},
			expectErr:       false,
		},
		{
			name:            "Valid ID global mapping",
			input:           "link|samlIdpID|providerGroupID:role",
			expectedMapping: SAMLGroupMappingsImport{SamlIdpID: "samlIdpID", GroupMappings: []wiz.SAMLGroupDetailsInput{{ProviderGroupID: "providerGroupID", Role: "role", Projects: nil}}},
			expectErr:       false,
		},
		{
			name:            "Invalid ID",
			input:           "invalidId",
			expectedMapping: SAMLGroupMappingsImport{},
			expectErr:       true,
		},
		{
			name:            "Invalid ID length",
			input:           "link|samlIdpId",
			expectedMapping: SAMLGroupMappingsImport{},
			expectErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mapping, err := extractIDsFromSamlIdpGroupMappingImportID(tc.input)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}
			if !reflect.DeepEqual(mapping, tc.expectedMapping) {
				t.Errorf("Expected mapping: %+v, got: %+v", tc.expectedMapping, mapping)
			}
		})
	}
}
