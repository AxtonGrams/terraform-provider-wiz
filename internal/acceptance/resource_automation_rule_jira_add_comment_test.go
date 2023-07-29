package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWizAutomationRuleJiraAddComment_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix(ResourcePrefix)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, TestCase(TcServiceNow)) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceWizAutomationRuleJiraAddCommentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wiz_integration_jira.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_add_comment.foo",
						"name",
						rName,
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_add_comment.foo",
						"description",
						"Provider Acceptance Test",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_add_comment.foo",
						"enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_add_comment.foo",
						"trigger_source",
						"CONTROL",
					),
					resource.TestCheckResourceAttr(
						"wiz_automation_rule_jira_add_comment.foo",
						"trigger_type.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttr(
						"wiz_automation_rule_jira_add_comment.foo",
						"trigger_type.*",
						"UPDATED",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"id",
						"wiz_automation_rule_jira_add_comment.foo",
						"integration_id",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_project_key",
						"wiz_automation_rule_jira_add_comment.foo",
						os.Getenv("WIZ_INTEGRATION_JIRA_PROJECT"),
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_comment",
						"wiz_automation_rule_jira_add_comment.foo",
						"Comment added via Wiz automation",
					),
					resource.TestCheckResourceAttrPair(
						"wiz_integration_jira.foo",
						"jira_add_issues_report",
						"wiz_automation_rule_jira_add_comment.foo",
						"false",
					),
				),
			},
		},
	})
}

func testResourceWizAutomationRuleJiraAddCommentBasic(rName string) string {
	return fmt.Sprintf(`
resource "wiz_integration_jira" "foo" {
  name                = "%s"
  scope               = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_jira_add_comment" "foo" {
  name           = "%s"
  description    = "Provider Acceptance Test"
  enabled        = false
  integration_id = wiz_integration_jira.foo.id
  trigger_source = "CONTROL"
  trigger_type = [
    "UPDATED",
  ]
  filters = jsonencode({
      "severity": [
        "CRITICAL"
      ]
    })
  jira_project_key = "%s"
  jira_comment = "Comment added via Wiz automation"
  jira_add_issues_report = false
}
`, rName, rName, os.Getenv("WIZ_INTEGRATION_JIRA_PROJECT"))
}
