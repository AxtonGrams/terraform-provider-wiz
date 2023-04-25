package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizAutomationRuleServiceNowCreateTicket_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckIntegrationServiceNow(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizAutomationRuleServiceNowCreateTicketBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_aws_sns.foo",
						"name",
						"test-acc-WizAutomationRuleAwsSNS_basic",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_aws_sns.foo",
						"name",
						"test-acc-WizAutomationRuleAwsSNS_basic",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_aws_sns.foo",
						"description",
						"Terraform provider acceptance test TestAccResourceWizAutomationRuleAwsSNS_basic",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_aws_sns.foo",
						"enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_aws_sns.foo",
						"trigger_source",
						"ISSUES",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_aws_sns.foo",
						"trigger_type.#",
						"2",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_automation_rule_aws_sns.foo",
						"trigger_type.*",
						"CREATED",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_automation_rule_aws_sns.foo",
						"trigger_type.*",
						"REOPENED",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_aws_sns.foo",
						"id",
						"wiz_automation_rule_aws_sns.foo",
						"integration_id",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_aws_sns.foo",
						"aws_sns_body",
						"{\n  \"trigger\": {\n    \"source\": \"{{triggerSource}}\",\n    \"type\": \"{{triggerType}}\",\n    \"ruleId\": \"{{ruleId}}\",\n    \"ruleName\": \"{{ruleName}}\"\n  },\n  \"issue\": {\n    \"id\": \"{{issue.id}}\",\n    \"status\": \"{{issue.status}}\",\n    \"severity\": \"{{issue.severity}}\",\n    \"created\": \"{{issue.createdAt}}\",\n    \"projects\": \"{{#issue.projects}}{{name}}, {{/issue.projects}}\"\n  },\n  \"resource\": {\n    \"id\": \"{{issue.entitySnapshot.providerId}}\",\n    \"name\": \"{{issue.entitySnapshot.name}}\",\n    \"type\": \"{{issue.entitySnapshot.nativeType}}\",\n    \"cloudPlatform\": \"{{issue.entitySnapshot.cloudPlatform}}\",\n    \"subscriptionId\": \"{{issue.entitySnapshot.subscriptionExternalId}}\",\n    \"subscriptionName\": \"{{issue.entitySnapshot.subscriptionName}}\",\n    \"region\": \"{{issue.entitySnapshot.region}}\",\n    \"status\": \"{{issue.entitySnapshot.status}}\",\n    \"cloudProviderURL\": \"{{issue.entitySnapshot.cloudProviderURL}}\"\n  },\n  \"control\": {\n    \"id\": \"{{issue.control.id}}\",\n    \"name\": \"{{issue.control.name}}\",\n    \"description\": \"{{issue.control.description}}\",\n    \"severity\": \"{{issue.control.severity}}\",\n    \"sourceCloudConfigurationRuleId\": \"{{issue.control.sourceCloudConfigurationRule.shortId}}\",\n    \"sourceCloudConfigurationRuleName\": \"{{issue.control.sourceCloudConfigurationRule.name}}\"\n  }\n}",
					),
					resource.TestMatchResourceAttr(
						"wiz_automation_rule_aws_sns.foo",
						"filters",
						regexp.MustCompile("b95efbdb-ac2e-4deb-b9a7-23211f3a5d0a"),
					),
				),
			},
		},
	})
}

func testResourceWizAutomationRuleServiceNowCreateTicketBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_servicenow" "foo" {
  name  = "%s"
  scope = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_servicenow_create_ticket" "foo" {
  name           = "%s"
  description    = "Provider Acceptance Test"
  enabled        = false
  integration_id = wiz_integration_servicenow.foo.id
  trigger_source = "ISSUES"
  trigger_type = [
    "CREATED",
  ]
  filters = jsonencode({
      "severity": [
        "CRITICAL"
      ]
    })
  servicenow_table_name  = "incident"
  servicenow_summary     = "Wiz Issue: {{issue.control.name}}"
  servicenow_description = <<EOT
Description:  {{issue.description}}
Status:       {{issue.status}}
Created:      {{issue.createdAt}}
Severity:     {{issue.severity}}
Project:      {{#issue.projects}}{{name}}, {{/issue.projects}}

---
Resource:                   {{issue.entitySnapshot.name}}
Type:                   {{issue.entitySnapshot.nativeType}}
Cloud Platform:         {{issue.entitySnapshot.cloudPlatform}}
Cloud Resource URL:     {{issue.entitySnapshot.cloudProviderURL}}
Subscription Name (ID): {{issue.entitySnapshot.subscriptionName}} ({{issue.entitySnapshot.subscriptionExternalId}})
Region:                 {{issue.entitySnapshot.region}}
Please click the following link to proceed to investigate the issue:
https://{{wizDomain}}/issues#~(issue~'{{issue.id}})
Source Automation Rule: {{ruleName}}
EOT
`, rName, rName)
}




const testAccResourceWizAutomationRuleAwsSNSBasic = `
resource "wiz_integration_servicenow" "foo" {
  name                      = "test-acc-WizAutomationRuleAwsSNS_basic"
  aws_sns_topic_arn         = "arn:aws:sns:us-east-1:123456789012:Wiz"
  aws_sns_access_method     = "ASSUME_SPECIFIED_ROLE"
  aws_sns_customer_role_arn = "arn:aws:iam::123456789012:role/Wiz"
  scope                     = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_aws_sns" "foo" {
  name           = "test-acc-WizAutomationRuleAwsSNS_basic"
  description    = "Terraform provider acceptance test TestAccResourceWizAutomationRuleAwsSNS_basic"
  enabled        = false
  integration_id = wiz_integration_aws_sns.foo.id
  trigger_source = "ISSUES"
  trigger_type = [
    "CREATED",
  ]
  aws_sns_body = "{\n  \"trigger\": {\n    \"source\": \"{{triggerSource}}\",\n    \"type\": \"{{triggerType}}\",\n    \"ruleId\": \"{{ruleId}}\",\n    \"ruleName\": \"{{ruleName}}\"\n  },\n  \"issue\": {\n    \"id\": \"{{issue.id}}\",\n    \"status\": \"{{issue.status}}\",\n    \"severity\": \"{{issue.severity}}\",\n    \"created\": \"{{issue.createdAt}}\",\n    \"projects\": \"{{#issue.projects}}{{name}}, {{/issue.projects}}\"\n  },\n  \"resource\": {\n    \"id\": \"{{issue.entitySnapshot.providerId}}\",\n    \"name\": \"{{issue.entitySnapshot.name}}\",\n    \"type\": \"{{issue.entitySnapshot.nativeType}}\",\n    \"cloudPlatform\": \"{{issue.entitySnapshot.cloudPlatform}}\",\n    \"subscriptionId\": \"{{issue.entitySnapshot.subscriptionExternalId}}\",\n    \"subscriptionName\": \"{{issue.entitySnapshot.subscriptionName}}\",\n    \"region\": \"{{issue.entitySnapshot.region}}\",\n    \"status\": \"{{issue.entitySnapshot.status}}\",\n    \"cloudProviderURL\": \"{{issue.entitySnapshot.cloudProviderURL}}\"\n  },\n  \"control\": {\n    \"id\": \"{{issue.control.id}}\",\n    \"name\": \"{{issue.control.name}}\",\n    \"description\": \"{{issue.control.description}}\",\n    \"severity\": \"{{issue.control.severity}}\",\n    \"sourceCloudConfigurationRuleId\": \"{{issue.control.sourceCloudConfigurationRule.shortId}}\",\n    \"sourceCloudConfigurationRuleName\": \"{{issue.control.sourceCloudConfigurationRule.name}}\"\n  }\n}"
  filters = jsonencode({
    "project" : [],
    "relatedEntity" : {
      "cloudPlatform" : [
        "AWS"
      ],
      "subscriptionId" : [
        "b95efbdb-ac2e-4deb-b9a7-23211f3a5d0a",
        "2d036cf5-7062-4b3d-83ce-fad305a2fef1"
      ]
    },
    "sourceControl" : [
      "253702e2-4ef6-4f6f-af4b-f3eae38142c7",
      "b2a1243d-c701-4f83-9544-58f7ebb31c49",
    ]
  })
}
`
