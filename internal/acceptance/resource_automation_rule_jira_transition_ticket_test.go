package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizAutomationRuleJiraTransitionTicket_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcServiceNow)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizAutomationRuleJiraTransitionTicketBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_jira.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_transition_ticket.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_transition_ticket.foo",
						"description",
						"Provider Acceptance Test",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_transition_ticket.foo",
						"enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_transition_ticket.foo",
						"trigger_source",
						"ISSUES",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_transition_ticket.foo",
						"trigger_type.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_automation_rule_jira_transition_ticket.foo",
						"trigger_type.*",
						"RESOLVED",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"id",
						"wiz_automation_rule_jira_transition_ticket.foo",
						"integration_id",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_project",
						"wiz_automation_rule_jira_transition_ticket.foo",
						os.Getenv("WIZ_INTEGRATION_JIRA_PROJECT"),
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_transition_id",
						"wiz_automation_rule_jira_transition_ticket.foo",
						"Resolved",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_advanced_fields",
						"wiz_automation_rule_jira_transition_ticket.foo",
						"Wiz Issue: {{issue.control.name}}",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_comment",
						"wiz_automation_rule_jira_transition_ticket.foo",
						"Jira comment from Wiz",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_comment_on_transition",
						"wiz_automation_rule_jira_transition_ticket.foo",
						"false",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_attach_evidence_csv",
						"wiz_automation_rule_jira_transition_ticket.foo",
						"false",
					),
				),
			},
		},
	})
}

func testResourceWizAutomationRuleJiraTransitionTicketBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_jira" "foo" {
  name                = "%s"
  scope               = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_jira_transition_ticket" "foo" {
  name           = "%s"
  description    = "Provider Acceptance Test"
  enabled        = false
  integration_id = wiz_integration_jira.foo.id
  trigger_source = "ISSUES"
  trigger_type = [
    "RESOLVED",
  ]
  filters = jsonencode({
      "severity": [
        "CRITICAL"
      ]
    })
  jira_project = "%s"
  jira_transition_id = "Resolved"
  jira_advanced_fields = jsonencode({
    "resolution" : "Done"
  })
  jira_comment = "Resolved via Wiz Automation"
  jira_comment_on_transition = true
  jira_attach_evidence_csv = false
}
`, rName, rName, os.Getenv("WIZ_INTEGRATION_JIRA_PROJECT"))
}
