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
						// check connector has name that matches the random string
						"wiz_connector_aws.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						// check connector has customer_role_arn that matches the random string
						"wiz_connector_aws.foo",
						"customer_role_arn",
						fmt.Sprintf("arn:aws:iam::000000000000:role/%s", rName),
					),
					resource.TestMatchResourceAttr(
						// check connector has sts external_id / nonce that matches the UUID pattern
						"wiz_connector_aws.foo",
						"external_id_nonce",
						regexp.MustCompile(UUIDPattern),
					),
					resource.TestCheckTypeSetElemAttr(
						// check connector has opted_in_regions that matches requested region
						"wiz_connector_aws.foo",
						"opted_in_regions.*",
						"us-east-1",
					),
					resource.TestCheckTypeSetElemAttr(
						// check connector has excluded_ous that matches requested OU
						"wiz_connector_aws.foo",
						"excluded_ous.*",
						"DEV",
					),
					resource.TestCheckTypeSetElemAttr(
						// check connector has excluded_accounts that matches requested account number
						"wiz_connector_aws.foo",
						"excluded_accounts.*",
						"100000000009",
					),
					resource.TestCheckResourceAttr(
						// check connector has skip_organization_scan that matches request
						"wiz_connector_aws.foo",
						"skip_organization_scan",
						"true",
					),
					resource.TestCheckResourceAttr(
						// check connector has audit_log_monitor_enabled that matches request
						"wiz_connector_aws.foo",
						"enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						// check connector has disk_analyzer_inflight_disabled that matches request
						"wiz_connector_aws.foo",
						"disk_analyzer_inflight_disabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						// check connector has extra_config that matches request
						"wiz_connector_aws.foo",
						"extra_config",
						"{\"auditLogMonitorEnabled\":false,\"diskAnalyzerInFlightDisabled\":false,\"excludedAccounts\":[\"100000000009\"],\"excludedOUs\":[\"DEV\"],\"optedInRegions\":[\"us-east-1\"],\"skipOrganizationScan\":true}",
					),
					resource.TestCheckResourceAttr(
						// check connector has auth_params that matches request
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
