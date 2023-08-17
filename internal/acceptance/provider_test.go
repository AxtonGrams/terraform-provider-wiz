package acceptance

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/provider"
)

// required/common environment variables for acceptance tests
var commonEnvVars = []string{"WIZ_URL", "WIZ_AUTH_CLIENT_ID", "WIZ_AUTH_CLIENT_SECRET"}

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

func testAccPreCheck(t *testing.T, tc TestCase) {
	var envVars []string
	switch tc {
	case TcCommon:
		envVars = commonEnvVars
	case TcUser:
		envVars = append(commonEnvVars, "WIZ_SMTP_DOMAIN")
	case TcServiceNow:
		envVars = append(commonEnvVars, "WIZ_INTEGRATION_SERVICENOW_URL", "WIZ_INTEGRATION_SERVICENOW_USERNAME", "WIZ_INTEGRATION_SERVICENOW_PASSWORD")
	case TcJira:
		envVars = append(commonEnvVars, "WIZ_INTEGRATION_JIRA_URL", "WIZ_INTEGRATION_JIRA_USERNAME", "WIZ_INTEGRATION_JIRA_PASSWORD", "WIZ_INTEGRATION_JIRA_PROJECT")
	case TcSubscriptionResourceGroups:
		envVars = append(commonEnvVars, "WIZ_SUBSCRIPTION_ID")
	case TcProject:
		envVars = append(commonEnvVars, "WIZ_SUBSCRIPTION_ID")
	case TcCloudConfigRule:
		envVars = append(commonEnvVars, "WIZ_SUBSCRIPTION_ID")
	default:
		t.Fatalf("unknown testCase: %s", tc)
	}

	if err := checkEnvVars(t, envVars); err != nil {
		t.Fatal(err)
	}
}

// checkEnvVars checks that the given environment variables are set and returns the error as appropriate
func checkEnvVars(t *testing.T, names []string) error {
	var unsetVars []string
	term := os.Getenv("TERM")
	supportsColor := term == "xterm" || term == "xterm-256color" || term == "screen" || term == "screen-256color"
	for _, name := range names {
		if v := os.Getenv(name); v == "" {
			unsetVars = append(unsetVars, name)
		}
	}
	if len(unsetVars) > 0 {
		var errMsg string
		if supportsColor {
			errMsg = fmt.Sprintf("\033[31m%s\033[0m must be set for acceptance tests", strings.Join(unsetVars, ", "))
		} else {
			errMsg = fmt.Sprintf("%s must be set for acceptance tests", strings.Join(unsetVars, ", "))
		}
		return fmt.Errorf(errMsg)
	}
	return nil
}
