resource "wiz_integration_servicenow" "default" {
  name                = "default"
  servicenow_url      = var.servicename_url
  servicenow_username = var.servicenow_username
  servicenow_password = var.servicenow_password
  scope               = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_servicenow_update_ticket" "example" {
  name           = "example"
  description    = "example description"
  enabled        = true
  integration_id = wiz_integration_servicenow.default.id
  trigger_source = "ISSUES"
  trigger_type = [
    "RESOLVED",
  ]
  filters = jsonencode({
    "severity" : [
      "CRITICAL"
    ]
  })
  servicenow_table_name           = "incident"
  servicenow_attach_issues_report = true
  servicenow_fields = jsonencode({
    "state" : "Closed"
  })
}
