package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceWizSubscriptionResourceGroups_basic tests the basic functionality of the datasource
// wiz_subscription_resource_groups. `WIZ_SUBSCRIPTION_ID` environment variable must be set to a valid internal
// wiz identifier for a subscription that has resource groups for this test to pass
func TestAccDatasourceWizSubscriptionResourceGroups_basic(t *testing.T) {
	subscriptionID := os.Getenv("WIZ_SUBSCRIPTION_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcSubscriptionResourceGroups)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWizSubscriptionResourceGroupsBasic(subscriptionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check first resource group has an id that matches the UUID regex
						"data.wiz_subscription_resource_groups.foo",
						"resource_groups.0.id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check first resource group has a name that is set to a non-empty string
						"data.wiz_subscription_resource_groups.foo",
						"resource_groups.0.name",
						regexp.MustCompile(`\w`),
					),
				),
			},
		},
	})
}

func testAccDatasourceWizSubscriptionResourceGroupsBasic(subscriptionID string) string {
	return fmt.Sprintf(`
	data "wiz_subscription_resource_groups" "foo" {
		subscription_id = "%s"
		first           = 2
	  }
`, subscriptionID)
}
