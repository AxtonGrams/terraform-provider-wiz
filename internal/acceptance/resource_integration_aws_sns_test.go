package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizIntegrationAwsSNS_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcCommon)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWizIntegrationAwsSNSBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_aws_sns.foo",
						"name",
						"test-acc-WizIntegrationAwsSNS_basic",
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_aws_sns.foo",
						"aws_sns_topic_arn",
						"arn:aws:sns:us-east-1:123456789012:Wiz-Remediation-Issues-Topic",
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_aws_sns.foo",
						"aws_sns_access_method",
						"ASSUME_SPECIFIED_ROLE",
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_aws_sns.foo",
						"aws_sns_customer_role_arn",
						"arn:aws:iam::123456789012:role/WizAccess-Role",
					),
					resource.TestCheckResourceAttr(
						"wiz_integration_aws_sns.foo",
						"scope",
						"All Resources, Restrict this Integration to global roles only",
					),
				),
			},
		},
	})
}

const testAccResourceWizIntegrationAwsSNSBasic = `
resource "wiz_integration_aws_sns" "foo" {
  name                      = "test-acc-WizIntegrationAwsSNS_basic"
  aws_sns_topic_arn         = "arn:aws:sns:us-east-1:123456789012:Wiz-Remediation-Issues-Topic"
  aws_sns_access_method     = "ASSUME_SPECIFIED_ROLE"
  aws_sns_customer_role_arn = "arn:aws:iam::123456789012:role/WizAccess-Role"
  scope                     = "All Resources, Restrict this Integration to global roles only"
}
`
