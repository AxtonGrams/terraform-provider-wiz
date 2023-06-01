package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceWizCloudAccounts_basic tests the basic functionality of the datasource
// wiz_cloud_accounts. The assumption is that at least two accounts exist in the Wiz tenant in order
// to validate pagination functionality
func TestAccDatasourceWizCloudAccounts_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWizCloudAccountsBasic(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check that the first cloud account has an id that matches the UUID regex
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.0.id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestMatchResourceAttr(
						// check that cloud_provider is set to a non-empty string
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.0.cloud_provider",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check that external_id is set to a non-empty string, different cloud providers have different formats
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.0.external_id",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check that name is set to a non-empty string
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.0.name",
						regexp.MustCompile(`\w`),
					),
					resource.TestMatchResourceAttr(
						// check that status is set to a non-empty string
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.0.status",
						regexp.MustCompile(`\w`),
					),
					resource.TestCheckResourceAttrSet(
						// check that linked_project_ids is a set or list, can be empty
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.0.linked_project_ids.#",
					),
					resource.TestCheckResourceAttrSet(
						// check that source_connector_ids is a set or list, can be empty
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.0.source_connector_ids.#",
					),
				),
			},
			{
				Config: testAccDatasourceWizCloudAccountsBasic(2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check that the second cloud account has an id that matches the UUID regex
						"data.wiz_cloud_accounts.foo",
						"cloud_accounts.1.id",
						regexp.MustCompile(UUIDPattern),
					),
				),
			},
		},
	})
}

func testAccDatasourceWizCloudAccountsBasic(maxPages int) string {
	return fmt.Sprintf(`
	data "wiz_cloud_accounts" "foo" {
		search   = [""]
		max_pages    = %d
		first = 1
	}
`, maxPages)
}
