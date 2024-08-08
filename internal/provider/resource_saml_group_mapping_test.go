package provider

import (
	"reflect"
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
			input:           "link|samlIdpID|providerGroupID|projectID1,projectID2|role",
			expectedMapping: SAMLGroupMappingsImport{SamlIdpID: "samlIdpID", ProviderGroupID: "providerGroupID", ProjectIDs: []string{"projectID1", "projectID2"}, Role: "role"},
			expectErr:       false,
		},
		{
			name:            "Valid ID global mapping",
			input:           "link|samlIdpID|providerGroupID|global|role",
			expectedMapping: SAMLGroupMappingsImport{SamlIdpID: "samlIdpID", ProviderGroupID: "providerGroupID", ProjectIDs: nil, Role: "role"},
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
