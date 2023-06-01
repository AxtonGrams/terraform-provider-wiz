package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizAutomationRuleServiceNowCreateTicket_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcServiceNow)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizAutomationRuleServiceNowCreateTicketBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"description",
						"Provider Acceptance Test",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"trigger_source",
						"ISSUES",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"trigger_type.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"trigger_type.*",
						"CREATED",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_servicenow.foo",
						"id",
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"integration_id",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"servicenow_table_name",
						"incident",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"servicenow_summary",
						"Wiz Issue: {{issue.control.name}}",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_create_ticket.foo",
						"servicenow_description",
						"Description:  {{issue.description}}\nStatus:       {{issue.status}}\nCreated:      {{issue.createdAt}}\nSeverity:     {{issue.severity}}\nProject:      {{#issue.projects}}{{name}}, {{/issue.projects}}\n\n---\nResource:                   {{issue.entitySnapshot.name}}\nType:                   {{issue.entitySnapshot.nativeType}}\nCloud Platform:         {{issue.entitySnapshot.cloudPlatform}}\nCloud Resource URL:     {{issue.entitySnapshot.cloudProviderURL}}\nSubscription Name (ID): {{issue.entitySnapshot.subscriptionName}} ({{issue.entitySnapshot.subscriptionExternalId}})\nRegion:                 {{issue.entitySnapshot.region}}\nPlease click the following link to proceed to investigate the issue:\nhttps://{{wizDomain}}/issues#~(issue~'{{issue.id}})\nSource Automation Rule: {{ruleName}}\n",
					),
				),
			},
		},
	})
}

func testResourceWizAutomationRuleServiceNowCreateTicketBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_servicenow" "foo" {
  name                = "%s"
  scope               = "All Resources, Restrict this Integration to global roles only"
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
}
`, rName, rName)
}
