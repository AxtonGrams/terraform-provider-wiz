resource "wiz_integration_servicenow" "default" {
  name                = "default"
  servicenow_url      = var.servicename_url
  servicenow_username = var.servicenow_username
  servicenow_password = var.servicenow_password
  scope               = "All Resources, Restrict this Integration to global roles only"
}

resource "wiz_automation_rule_servicenow_create_ticket" "example" {
  name           = "example"
  description    = "example description"
  enabled        = true
  integration_id = wiz_integration_servicenow.default.id
  trigger_source = "ISSUES"
  trigger_type = [
    "CREATED",
  ]
  filters = jsonencode({
    "severity" : [
      "CRITICAL"
    ]
  })
  servicenow_table_name  = "incident"
  servicenow_summary     = "Wiz Issue: {{issue.control.name}}"
  servicenow_description = "Description:  {{issue.description}}\nStatus:       {{issue.status}}\nCreated:      {{issue.createdAt}}\nSeverity:     {{issue.severity}}\nProject:      {{#issue.projects}}{{name}}, {{/issue.projects}}\n\n---\nResource:\t            {{issue.entitySnapshot.name}}\nType:\t                {{issue.entitySnapshot.nativeType}}\nCloud Platform:\t        {{issue.entitySnapshot.cloudPlatform}}\nCloud Resource URL:     {{issue.entitySnapshot.cloudProviderURL}}\nSubscription Name (ID): {{issue.entitySnapshot.subscriptionName}} ({{issue.entitySnapshot.subscriptionExternalId}})\nRegion:\t                {{issue.entitySnapshot.region}}\nPlease click the following link to proceed to investigate the issue:\nhttps://{{wizDomain}}/issues#~(issue~'{{issue.id}})\nSource Automation Rule: {{ruleName}}"
}
