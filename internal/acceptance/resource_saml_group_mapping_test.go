package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizSAMLGroupMapping_basic(t *testing.T) {
	samlIdpID := os.Getenv("WIZ_SAML_IDP_ID")
	providerGroupID := os.Getenv("WIZ_PROVIDER_GROUP_ID")
	projectID := os.Getenv("WIZ_PROJECT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TcSAMLGroupMapping) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizSAMLGroupMappingBasic(samlIdpID, providerGroupID, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_saml_group_mapping.foo",
						"saml_idp_id",
						samlIdpID,
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_group_mapping.foo",
						"group_mapping.0.provider_group_id",
						providerGroupID,
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_group_mapping.foo",
						"group_mapping.0.projects.0",
						projectID,
					),
					resource.TestCheckResourceAttr(
						"wiz_saml_group_mapping.foo",
						"group_mapping.0.description",
						"test mapping.",
					),
				),
			},
		},
	})
}

func testResourceWizSAMLGroupMappingBasic(samlIdpID string, providerGroupID string, projectID string) string {
	return fmt.Sprintf(`
		resource "wiz_saml_group_mapping" "foo" {
		  saml_idp_id = "%s"
		  group_mapping {
			provider_group_id = "%s"
			role = "PROJECT_READER"
			projects = [
			  "%s"
			]
			description = "test mapping."
		  }
		}`, samlIdpID, providerGroupID, projectID)
}
