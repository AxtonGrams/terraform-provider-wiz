package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceWizUsers_basic tests the basic functionality of the datasource wiz_users
// the assumption is that at least two users exist in the Wiz tenant in order
// to validate pagination functionality
func TestAccDatasourceWizUsers_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWizUsersBasic(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check that the first user has an id that matches the UUID regex
						"data.wiz_users.foo",
						"users.0.id",
						regexp.MustCompile(uuidPattern),
					),
				),
			},
			{
				Config: testAccDatasourceWizUsersBasic(2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						// check that the second user has an id that matches the UUID regex
						"data.wiz_users.foo",
						"users.1.id",
						regexp.MustCompile(uuidPattern),
					),
				),
			},
		},
	})
}

const uuidPattern = `^[A-Fa-f0-9]{8}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{12}`

func testAccDatasourceWizUsersBasic(maxPages int) string {
	return fmt.Sprintf(`
	data "wiz_users" "foo" {
		search   = ""
		first = 1
		max_pages = %d
	  }
`, maxPages)
}
