package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizAutomationRuleServiceNowUpdateTicket_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcServiceNow)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizAutomationRuleServiceNowUpdateTicketBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_servicenow.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"description",
						"Provider Acceptance Test",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"trigger_source",
						"ISSUES",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"trigger_type.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"trigger_type.*",
						"RESOLVED",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_servicenow.foo",
						"id",
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"integration_id",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"servicenow_table_name",
						"incident",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"servicenow_attach_issues_report",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_servicenow_update_ticket.foo",
						"servicenow_fields",
						"{\"state\":\"Closed\"}",
					),
				),
			},
		},
	})
}

func testResourceWizAutomationRuleServiceNowUpdateTicketBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_servicenow" "foo" {
  name                = "%s"
  scope               = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_servicenow_update_ticket" "foo" {
  name           = "%s"
  description    = "Provider Acceptance Test"
  enabled        = false
  integration_id = wiz_integration_servicenow.foo.id
  trigger_source = "ISSUES"
  trigger_type = [
    "RESOLVED",
  ]
  filters = jsonencode({
      "severity": [
        "CRITICAL"
      ]
    })
  servicenow_table_name  = "incident"
  servicenow_fields = jsonencode({
    "state" : "Closed"
  })
  servicenow_attach_issues_report = false
}
`, rName, rName)
}
