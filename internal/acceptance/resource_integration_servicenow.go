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
		PreCheck:          func() { testAccPreCheckIntegrationServiceNow(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizIntegrationServiceNowBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"name",
						"test-acc-WizIntegrationServiceNow_basic",
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"servicenow_url",
						os.Getenv("WIZ_AUTH_CLIENT_ID"),
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"servicenow_url",
						os.Getenv("WIZ_INTEGRATION_SERVICENOW_URL"),
					),
				),
			},
		},
	})
}

func testResourceWizIntegrationServiceNowBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_servicenow" "test" {
  name  = "%s"
  scope = "All Resources, Restrict this Integration to global roles only"
}
`, rName)
}
