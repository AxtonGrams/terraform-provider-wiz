package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func TestFlattenRoleScopes(t *testing.T) {
	expected := []interface{}{
		"read:all",
		"create:reports",
	}

	var scopes = []string{
		"read:all",
		"create:reports",
	}

	flattened := utils.ConvertSliceToGenericArray(scopes)

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}
func TestFlattenIdpUsers(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"id":                     "5480a8dc-68e1-5835-82d2-0d9a9eac7caf",
			"email":                  "john.doe@w.com",
			"name":                   "John Doe",
			"is_suspended":           false,
			"identity_provider_type": "SAML",
			"identity_provider": []interface{}{
				map[string]interface{}{
					"name": "myIdp",
				},
			},
			"effective_role": []interface{}{
				map[string]interface{}{
					"id":     "GLOBAL_READER",
					"name":   "GlobalReader",
					"scopes": []string{"read:all", "create:reports"},
				},
			},
		},
		map[string]interface{}{
			"id":                     "1480a8dc-68e1-5835-82d2-0d9a9eac7cac",
			"email":                  "jane.doe@w.com",
			"name":                   "Jane Doe",
			"is_suspended":           true,
			"identity_provider_type": "SAML",
			"identity_provider": []interface{}{
				map[string]interface{}{
					"name": "myIdp",
				},
			},
			"effective_role": []interface{}{
				map[string]interface{}{
					"id":     "GLOBAL_READER",
					"name":   "GlobalReader",
					"scopes": []string{"read:all", "create:reports"},
				},
			},
		},
	}

	readUsers1 := ReadUsers{
		Users: wiz.UserConnection{
			Nodes: []*wiz.User{
				{
					ID:                   "5480a8dc-68e1-5835-82d2-0d9a9eac7caf",
					Email:                "john.doe@w.com",
					Name:                 "John Doe",
					IsSuspended:          false,
					IdentityProviderType: "SAML",
					IdentityProvider:     wiz.SAMLIdentityProvider{Name: "myIdp"},
					EffectiveRole: wiz.UserRole{
						ID:   "GLOBAL_READER",
						Name: "GlobalReader",
						Scopes: []string{
							"read:all",
							"create:reports",
						},
					},
				},
			},
		},
	}

	readUsers2 := &ReadUsers{
		Users: wiz.UserConnection{
			Nodes: []*wiz.User{
				{
					ID:                   "1480a8dc-68e1-5835-82d2-0d9a9eac7cac",
					Email:                "jane.doe@w.com",
					Name:                 "Jane Doe",
					IsSuspended:          true,
					IdentityProviderType: "SAML",
					IdentityProvider: wiz.SAMLIdentityProvider{
						Name: "myIdp",
					},
					EffectiveRole: wiz.UserRole{
						ID:   "GLOBAL_READER",
						Name: "GlobalReader",
						Scopes: []string{
							"read:all",
							"create:reports",
						},
					},
				},
			},
		},
	}

	// Convert users to the correct type ([]interface{})
	users := make([]interface{}, 0)
	users = append(users, &readUsers1)
	users = append(users, readUsers2)

	result := flattenUsers(ctx, users)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected: %v, but got: %v", expected, result)
	}
}

func TestFlattenLocalUsers(t *testing.T) {
	ctx := context.Background()

	expected := []interface{}{
		map[string]interface{}{
			"id":                     "5480a8dc-68e1-5835-82d2-0d9a9eac7caf",
			"email":                  "john.doe@w.com",
			"name":                   "John Doe",
			"is_suspended":           false,
			"identity_provider_type": "WIZ",
			"identity_provider": []interface{}{
				map[string]interface{}{
					"name": "",
				},
			},
			"effective_role": []interface{}{
				map[string]interface{}{
					"id":   "GLOBAL_ADMIN",
					"name": "GlobalAdmin",
					"scopes": []string{
						"admin:all",
						"create:all",
						"delete:all",
						"read:all",
						"update:all",
					},
				},
			},
		},
	}

	readUsers := ReadUsers{
		Users: wiz.UserConnection{
			Nodes: []*wiz.User{
				{
					ID:                   "5480a8dc-68e1-5835-82d2-0d9a9eac7caf",
					Email:                "john.doe@w.com",
					Name:                 "John Doe",
					IsSuspended:          false,
					IdentityProviderType: "WIZ",
					IdentityProvider:     wiz.SAMLIdentityProvider{Name: ""},
					EffectiveRole: wiz.UserRole{
						ID:   "GLOBAL_ADMIN",
						Name: "GlobalAdmin",
						Scopes: []string{
							"admin:all",
							"create:all",
							"delete:all",
							"read:all",
							"update:all",
						},
					},
				},
			},
		},
	}

	// Convert users to the correct type ([]interface{})
	users := make([]interface{}, 0)
	users = append(users, &readUsers)

	result := flattenUsers(ctx, users)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected: %v, but got: %v", expected, result)
	}
}
