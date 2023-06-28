package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizConnectorAws_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizConnectorAwsBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_connector_aws.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_aws.foo",
						"customer_role_arn",
						fmt.Sprintf("arn:aws:iam::000000000000:role/%s", rName),
					),
					resource.TestMatchResourceAttr(
						"wiz_connector_aws.foo",
						"external_id_nonce",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_connector_aws.foo",
						"opted_in_regions.*",
						"us-east-1",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_connector_aws.foo",
						"excluded_ous.*",
						"DEV",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_connector_aws.foo",
						"excluded_accounts.*",
						"100000000009",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_aws.foo",
						"skip_organization_scan",
						"true",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_aws.foo",
						"enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_aws.foo",
						"disk_analyzer_inflight_disabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_aws.foo",
						"extra_config",
						"{\"auditLogMonitorEnabled\":false,\"diskAnalyzerInFlightDisabled\":false,\"excludedAccounts\":[\"100000000009\"],\"excludedOUs\":[\"DEV\"],\"optedInRegions\":[\"us-east-1\"],\"skipOrganizationScan\":true}",
					),
					resource.TestCheckResourceAttr(
						"wiz_connector_aws.foo",
						"auth_params",
						fmt.Sprintf("{\"customerRoleARN\":\"arn:aws:iam::000000000000:role/%s\"}", rName),
					),
				),
			},
		},
	})
}

func testResourceWizConnectorAwsBasic(rName string) string {
	return fmt.Sprintf(`
	resource "wiz_connector_aws" "foo" {
		name = "%[1]s"
		auth_params = jsonencode({
			"customerRoleARN" : "arn:aws:iam::000000000000:role/%[1]s",
		})
		extra_config = jsonencode(
			{
				"skipOrganizationScan" : true,
				"diskAnalyzerInFlightDisabled" : false,
				"optedInRegions" : ["us-east-1"],
				"excludedAccounts" : ["100000000009"],
				"excludedOUs" : ["DEV"],
				"auditLogMonitorEnabled" : false
			}
		)
	}
`, rName)
}
