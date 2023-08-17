package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizCloudConfigRule_basic(t *testing.T) {
	subscriptionID := os.Getenv("WIZ_SUBSCRIPTION_ID")
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCloudConfigRule)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizCloudConfigRuleBasic(rName, subscriptionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"description",
						"test description",
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"remediation_instructions",
						"fix it",
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"target_native_types.0",
						"account",
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"scope_account_ids.0",
						subscriptionID,
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"function_as_control",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"severity",
						"HIGH",
					),
					resource.TestCheckResourceAttr(
						"wiz_cloud_config_rule.foo",
						"iac_matchers.0.type",
						"ADMISSION_CONTROLLER",
					),
					resource.TestMatchResourceAttr(
						"wiz_cloud_config_rule.foo",
						"iac_matchers.0.rego_code",
						regexp.MustCompile(`\w`),
					),
				),
			},
		},
	})
}

func testResourceWizCloudConfigRuleBasic(rName string, subscriptionID string) string {
	return fmt.Sprintf(`
	resource "wiz_cloud_config_rule" "foo" {
		name        = "%s"
		description = "test description"
		target_native_types = [
		  "account",
		]
		scope_account_ids = [
		  "%s",
		]
		function_as_control      = false
		remediation_instructions = "fix it"
		enabled                  = false
		severity                 = "HIGH"
		opa_policy               = <<EOT
	  package wiz

	  default result = "pass"
	  EOT
		iac_matchers {
		  type      = "ADMISSION_CONTROLLER"
		  rego_code = <<EOT
	  package wiz

	  import data.generic.cloudformation as cloudFormationLib

	  import data.generic.common as common_lib

	  WizPolicy[result] {
			  resource := input.document[i].Resources[name]
			  resource.Type == "AWS::Config::ConfigRule"
			  not hasAccessKeyRotationRule(resource)

			  result := {
					  "documentId": input.document[i].id,
					  "searchKey": sprintf("Resources.%%s", [name]),
					  "issueType": "MissingAttribute",
					  "keyExpectedValue": sprintf("Resources.%%s has a ConfigRule defining rotation period on AccessKeys.", [name]),
					  "keyActualValue": sprintf("Resources.%%s doesn't have a ConfigRule defining rotation period on AccessKeys.", [name]),
					  "resourceTags": cloudFormationLib.getCFTags(resource),
			  }
	  }

	  hasAccessKeyRotationRule(configRule) {
			  configRule.Properties.Source.SourceIdentifier == "ACCESS_KEYS_ROTATED"
	  } else = false {
			  true
	  }
	  EOT
		}
	  }
`, rName, subscriptionID)
}
