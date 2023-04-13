package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/provider"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"wiz": func() (*schema.Provider, error) {
		return provider.New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := provider.New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("WIZ_URL"); v == "" {
		t.Fatal("WIZ_URL must be set for acceptance tests")
	}
	if v := os.Getenv("WIZ_AUTH_CLIENT_ID"); v == "" {
		t.Fatal("WIZ_AUTH_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("WIZ_AUTH_CLIENT_SECRET"); v == "" {
		t.Fatal("WIZ_AUTH_CLIENT_SECRET must be set for acceptance tests")
	}
}
