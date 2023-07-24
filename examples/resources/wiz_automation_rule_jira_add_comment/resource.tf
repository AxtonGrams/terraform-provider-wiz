resource "wiz_integration_jira" "default" {
  name          = "default"
  jira_url      = var.jira_url
  jira_username = var.jira_username
  jira_password = var.jira_password
  scope         = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_jira_add_comment" "example" {
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
  jira_project_key = "PROJ"
  jira_comment     = "Comment from Wiz"
}
