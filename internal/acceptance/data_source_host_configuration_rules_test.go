package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceWizHostConfigurationRules_basic tests the basic functionality of the datasource
// wiz_cloud_config_rules.
func TestAccDatasourceWizHostConfigurationRules_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWizHostConfigurationRulesBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check rule has an id that matches the UUID regex
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check rule has an external_id that matches the UUID regex
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.external_id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check rule has a builtin value set to bool
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.builtin",
						regexp.MustCompile(`true|false`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a enabled value set to bool
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.enabled",
						regexp.MustCompile(`true|false`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a name that is set to a non-empty string
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.name",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a short_name that is set to a non-empty string
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.short_name",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a description that is set to a non-empty string
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.description",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check rule has a direct_oval that is set to a non-empty string
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.direct_oval",
						regexp.MustCompile(`\w`),
					),
					resource.TestCheckResourceAttrSet(
						// check that security_sub_category_ids is a set or list, can be empty
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.security_sub_category_ids.#",
					),
					resource.TestCheckResourceAttrSet(
						// check that target_platform_ids is a set or list, can be empty
						"data.wiz_host_config_rules.foo",
						"host_configuration_rules.0.target_platform_ids.#",
					),
				),
			},
		},
	})
}

const testAccDatasourceWizHostConfigurationRulesBasic = `
data "wiz_host_config_rules" "foo" {
	first  = 5
	search = "access"
  }
`
