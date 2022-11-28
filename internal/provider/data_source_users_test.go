package provider

import (
	"context"
	"reflect"
	"testing"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func TestFlattenRoleScopes(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		"read:all",
		"create:reports",
	}

	var scopes = []string{
		"read:all",
		"create:reports",
	}

	flattened := flattenRoleScopes(ctx, scopes)

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
			"id": "71ef857e9018809bff5c7a5666f4f3eba2f8d141",
			"users": []interface{}{
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
							"id":   "GLOBAL_READER",
							"name": "GlobalReader",
							"scopes": []interface{}{
								"read:all",
								"create:reports",
							},
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
							"id":   "GLOBAL_READER",
							"name": "GlobalReader",
							"scopes": []interface{}{
								"read:all",
								"create:reports",
							},
						},
					},
				},
			},
		},
	}

	var users = &[]*vendor.User{
		{
			ID:                   "5480a8dc-68e1-5835-82d2-0d9a9eac7caf",
			Email:                "john.doe@w.com",
			Name:                 "John Doe",
			IsSuspended:          false,
			IdentityProviderType: "SAML",
			IdentityProvider: vendor.SAMLIdentityProvider{
				Name: "myIdp",
			},
			EffectiveRole: vendor.UserRole{
				ID:   "GLOBAL_READER",
				Name: "GlobalReader",
				Scopes: []string{
					"read:all",
					"create:reports",
				},
			},
		},
		{
			ID:                   "1480a8dc-68e1-5835-82d2-0d9a9eac7cac",
			Email:                "jane.doe@w.com",
			Name:                 "Jane Doe",
			IsSuspended:          true,
			IdentityProviderType: "SAML",
			IdentityProvider: vendor.SAMLIdentityProvider{
				Name: "myIdp",
			},
			EffectiveRole: vendor.UserRole{
				ID:   "GLOBAL_READER",
				Name: "GlobalReader",
				Scopes: []string{
					"read:all",
					"create:reports",
				},
			},
		},
	}

	flattened := []interface{}{
		map[string]interface{}{
			"id":    "71ef857e9018809bff5c7a5666f4f3eba2f8d141",
			"users": flattenUsers(ctx, users),
		},
	}

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}

func TestFlattenLocalUsers(t *testing.T) {
	ctx := context.Background()
	expected := []interface{}{
		map[string]interface{}{
			"id": "71ef857e9018809bff5c7a5666f4f3eba2f8d141",
			"users": []interface{}{
				map[string]interface{}{
					"id":                     "auth0|6328329c6914d76a4ff0d209",
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
							"scopes": []interface{}{
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
		},
	}

	var users = &[]*vendor.User{
		{
			ID:                   "auth0|6328329c6914d76a4ff0d209",
			Email:                "john.doe@w.com",
			Name:                 "John Doe",
			IsSuspended:          false,
			IdentityProviderType: "WIZ",
			IdentityProvider: vendor.SAMLIdentityProvider{
				Name: "",
			},
			EffectiveRole: vendor.UserRole{
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
	}

	flattened := []interface{}{
		map[string]interface{}{
			"id":    "71ef857e9018809bff5c7a5666f4f3eba2f8d141",
			"users": flattenUsers(ctx, users),
		},
	}

	if !reflect.DeepEqual(flattened, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			flattened,
			expected,
		)
	}
}
