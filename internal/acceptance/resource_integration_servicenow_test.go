package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizIntegrationServiceNow_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TcServiceNow) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizIntegrationServiceNowBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"servicenow_url",
						os.Getenv("WIZ_INTEGRATION_SERVICENOW_URL"),
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"servicenow_username",
						os.Getenv("WIZ_INTEGRATION_SERVICENOW_USERNAME"),
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"scope",
						"All Resources, Restrict this Integration to global roles only",
					),
				),
			},
		},
	})
}

func testResourceWizIntegrationServiceNowBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_servicenow" "foo" {
  name  = "%s"
  scope = "All Resources, Restrict this Integration to global roles only"
}
`, rName)
}
