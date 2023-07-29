resource "wiz_integration_jira" "default" {
  name          = "default"
  jira_url      = var.jira_url
  jira_username = var.jira_username
  jira_password = var.jira_password
  scope         = "All Resources, Restrict this Integration to global roles only"
}
