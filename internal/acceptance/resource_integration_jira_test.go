package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizIntegrationJira_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TcJira) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizIntegrationJiraBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_jira.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_jira.foo",
						"jira_url",
						os.Getenv("WIZ_INTEGRATION_JIRA_URL"),
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_jira.foo",
						"jira_username",
						os.Getenv("WIZ_INTEGRATION_JIRA_USERNAME"),
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_jira.foo",
						"scope",
						"All Resources, Restrict this Integration to global roles only",
					),
				),
			},
		},
	})
}

func testResourceWizIntegrationJiraBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_jira" "foo" {
  name  = "%s"
  scope = "All Resources, Restrict this Integration to global roles only"
}
`, rName)
}
