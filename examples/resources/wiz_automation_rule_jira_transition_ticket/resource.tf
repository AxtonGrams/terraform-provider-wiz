resource "wiz_integration_jira" "default" {
  name          = "default"
  jira_url      = var.jira_url
  jira_username = var.jira_username
  jira_password = var.jira_password
  scope         = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_jira_transition_ticket" "example" {
  name           = "example"
  description    = "example description"
  enabled        = true
  integration_id = wiz_integration_jira.default.id
  trigger_source = "ISSUES"
  trigger_type = [
    "RESOLVED",
  ]
  filters = jsonencode({
    "severity" : [
      "CRITICAL"
    ]
  })
  jira_project_key   = "PROJ"
  jira_transition_id = "Resolved"
  jira_advanced_fields = jsonencode({
    "resolution" : "Done"
  })
  jira_comment               = "Resolved via Wiz Automation"
  jira_add_issues_report     = false
  jira_comment_on_transition = true
  jira_attach_evidence_csv   = false
}
