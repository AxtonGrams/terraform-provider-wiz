package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizProject_basic(t *testing.T) {
	subscriptionID := os.Getenv("WIZ_SUBSCRIPTION_ID")
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcProject)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizProjectBasic(rName, subscriptionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_project.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_project.foo",
						"business_unit",
						"Technology",
					),
					resource.TestMatchResourceAttr(
						"wiz_project.foo",
						"slug",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestCheckResourceAttr(
						"wiz_project.foo",
						"archived",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_project.foo",
						"is_folder",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_project.foo",
						"cloud_account_link.0.cloud_account_id",
						"477ea00a-4d4d-5bb4-9fa6-634691e68de7",
					),
					resource.TestCheckResourceAttr(
						"wiz_project.foo",
						"risk_profile.0.business_impact",
						"MBI",
					),
				),
			},
		},
	})
}

func testResourceWizProjectBasic(rName string, subscriptionID string) string {
	return fmt.Sprintf(`
	resource "wiz_project" "foo" {
		name              = "%s"
		risk_profile {
		  business_impact = "MBI"
		}
		business_unit = "Technology"
		cloud_account_link {
		  cloud_account_id = "%s"
		  environment      = "PRODUCTION"
		  shared           = true
		}
	  }
`, rName, subscriptionID)
}
