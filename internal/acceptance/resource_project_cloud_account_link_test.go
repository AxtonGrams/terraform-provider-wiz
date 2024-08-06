package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizProjectCloudAccountLink_basic(t *testing.T) {
	projectID := os.Getenv("WIZ_PROJECT_ID")
	cloudAccountID := os.Getenv("WIZ_SUBSCRIPTION_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TcProjectCloudAccountLink) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizProjectCloudAccountLinkBasic(projectID, cloudAccountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_project_cloud_account_link.foo",
						"project_id",
						projectID,
					),
					resource.TestCheckResourceAttr(
						"wiz_project_cloud_account_link.foo",
						"cloud_account_id",
						cloudAccountID,
					),
				),
			},
		},
	})
}

func testResourceWizProjectCloudAccountLinkBasic(projectID string, cloudAccountID string) string {
	return fmt.Sprintf(`
resource "wiz_project_cloud_account_link" "foo" {
  project_id = "%s"
  cloud_account_id = "%s"
}
`, projectID, cloudAccountID)
}
