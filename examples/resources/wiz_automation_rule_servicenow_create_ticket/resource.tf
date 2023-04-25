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
  servicenow_description = <<EOT
Description:  {{issue.description}}
Status:       {{issue.status}}
Created:      {{issue.createdAt}}
Severity:     {{issue.severity}}
Project:      {{#issue.projects}}{{name}}, {{/issue.projects}}

---
Resource:	            {{issue.entitySnapshot.name}}
Type:	                {{issue.entitySnapshot.nativeType}}
Cloud Platform:	        {{issue.entitySnapshot.cloudPlatform}}
Cloud Resource URL:     {{issue.entitySnapshot.cloudProviderURL}}
Subscription Name (ID): {{issue.entitySnapshot.subscriptionName}} ({{issue.entitySnapshot.subscriptionExternalId}})
Region:	                {{issue.entitySnapshot.region}}
Please click the following link to proceed to investigate the issue:
https://{{wizDomain}}/issues#~(issue~'{{issue.id}})
Source Automation Rule: {{ruleName}}
EOT
}
