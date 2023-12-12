package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizConnectorGcp_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizConnectorGcpBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_connector_gcp.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_gcp.foo",
						"folder_id",
						"123456",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_gcp.foo",
						"auth_params",
						"{\"folder_id\":\"123456\",\"isManagedIdentity\":true}",
					),
					resource.TestMatchResourceAttr(
						"wiz_connector_gcp.foo",
						"id",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_gcp.foo",
						"enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_gcp.foo",
						"disk_analyzer_inflight_disabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_gcp.foo",
						"extra_config",
						"{\"auditLogMonitorEnabled\":false,\"diskAnalyzerInFlightDisabled\":false,\"excludedFolders\":[],\"excludedProjects\":[],\"includedFolders\":[],\"projects\":[]}",
					),
				),
			},
		},
	})
}

func testResourceWizConnectorGcpBasic(rName string) string {
	return fmt.Sprintf(`
	resource "wiz_connector_gcp" "foo" {
		name = "%[1]s"
		auth_params = jsonencode({
		  "isManagedIdentity" : true,
		  "folder_id" : "123456",
		})
		extra_config = jsonencode(
		  {
			"projects" : [],
			"excludedProjects" : [],
			"includedFolders" : [],
			"excludedFolders" : [],
			"diskAnalyzerInFlightDisabled" : false,
			"auditLogMonitorEnabled" : false,
		  }
		)
	  }
`, rName)
}
