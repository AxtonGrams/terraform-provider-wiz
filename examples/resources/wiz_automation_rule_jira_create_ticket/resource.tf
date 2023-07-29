resource "wiz_integration_jira" "default" {
  name          = "default"
  jira_url      = var.jira_url
  jira_username = var.jira_username
  jira_password = var.jira_password
  scope         = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_jira_create_ticket" "example" {
  name           = "example"
  description    = "example description"
  enabled        = true
  integration_id = wiz_integration_jira.default.id
  trigger_source = "ISSUES"
  trigger_type = [
    "CREATED",
  ]
  filters = jsonencode({
    "severity" : [
      "CRITICAL"
    ]
  })
  jira_summary     = "Wiz Issue: {{issue.control.name}}"
  jira_project     = "PROJ"
  jira_description = <<EOT
Description:  {{issue.description}}
Status:       {{issue.status}}
Created:      {{issue.createdAt}}
Severity:     {{issue.severity}}
Project:      {{#issue.projects}}{{name}}, {{/issue.projects}}
EOT
}
