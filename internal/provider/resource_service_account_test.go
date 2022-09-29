package provider

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceServiceAccount = `
resource "wiz_service_account" "foo" {
  name = "foo"
  scopes = [
    "read:projects",
  ]
}
`

func TestAccWizServiceAccount_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"wiz_service_account.foo",
						"name",
					),
				),
			},
		},
	})
}
*/
