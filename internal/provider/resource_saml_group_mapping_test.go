package provider

import (
	"reflect"
	"testing"
)

func TestExtractIdsFromSamlIdpGroupMappingImportId(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		expectedMapping SAMLGroupMappingsImport
		expectErr       bool
	}{
		{
			name:            "Valid ID",
			input:           "link|samlIdpId|providerGroupId|projectId1,projectId2|role",
			expectedMapping: SAMLGroupMappingsImport{SamlIdpID: "samlIdpId", ProviderGroupID: "providerGroupId", ProjectIDs: []string{"projectId1", "projectId2"}, Role: "role"},
			expectErr:       false,
		},
		{
			name:            "Valid ID global mapping",
			input:           "link|samlIdpId|providerGroupId|global|role",
			expectedMapping: SAMLGroupMappingsImport{SamlIdpID: "samlIdpId", ProviderGroupID: "providerGroupId", ProjectIDs: nil, Role: "role"},
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
			input:           "link|samlIdpId|providerGroupId",
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
