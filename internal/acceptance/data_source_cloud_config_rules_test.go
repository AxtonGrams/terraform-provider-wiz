package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceWizCloudConfigRules_basic tests the basic functionality of the datasource
// wiz_cloud_config_rules.
func TestAccDatasourceWizCloudConfigRules_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWizCloudConfigRulesBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check rule has an id that matches the UUID regex
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check rule has a graph id that matches the UUID regex
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.graph_id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check rule has a builtin value set to bool
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.builtin",
						regexp.MustCompile(`true|false`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a function_as_control value set to bool
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.function_as_control",
						regexp.MustCompile(`true|false`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a name that is set to a non-empty string
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.name",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a description that is set to a non-empty string
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.description",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check rule has enabled set to bool
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.enabled",
						regexp.MustCompile(`true|false`),
					),
					resource.TestMatchResourceAttr(
						// check rule has auto remediation set to bool
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.has_auto_remediation",
						regexp.MustCompile(`true|false`),
					),
					resource.TestMatchResourceAttr(
						// check rule has supports_nrt set to bool
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.supports_nrt",
						regexp.MustCompile(`true|false`),
					),
					resource.TestMatchResourceAttr(
						// check rule has service_type set to bool
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.service_type",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check that severity is set to a non-empty string
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.severity",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check that short_id is set to a non-empty string
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.short_id",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check that subject_entity_type is set to a non-empty string
						"data.wiz_cloud_config_rules.foo",
						"cloud_configuration_rules.0.subject_entity_type",
						regexp.MustCompile(`\w`),
					),
				),
			},
		},
	})
}

const testAccDatasourceWizCloudConfigRulesBasic = `
data "wiz_cloud_config_rules" "foo" {
  cloud_provider = [  "AWS", ]
  has_remediation = true
  first = 2
}
`
