variable "servicenow_password" {
  type        = string
  sensitive   = true
  description = "ServiceNow credential for Wiz integration"
}

variable "servicenow_user" {
  type        = string
  description = "ServiceNow username for Wiz integration"
}

variable "servicenow_url" {
  type        = string
  description = "ServiceNow URL"
}

variable "jira_token" {
  type        = string
  description = "Jira authentication token"
  sensitive   = true
}

# Service Now Ticket
resource "wiz_automation_action" "servicenow" {
  name                          = "provider_dev"
  type                          = "SERVICENOW_TICKET"
  is_accessible_to_all_projects = false
  servicenow_params {
    base_url = var.servicenow_url
    password = var.servicenow_password
    user     = var.servicenow_user
    ticket_fields {
      description         = <<-EOT
Description:  {{description}}
Status:       {{status}}
Created:      {{createdAt}}
Severity:     {{severity}}
Project:      {{#projects}}{{name}}, {{/projects}}

---
Resource:                   {{entitySnapshot.name}}
Type:                   {{entitySnapshot.nativeType}}
Cloud Platform:         {{entitySnapshot.cloudPlatform}}
Cloud Resource URL:     {{entitySnapshot.cloudProviderURL}}
Subscription Name (ID): {{entitySnapshot.subscriptionName}} ({{entitySnapshot.subscriptionExternalId}})
Region:                 {{entitySnapshot.region}}
Please click the following link to proceed to investigate the issue:
https://{{wizDomain}}/issues#~(issue~'{{id}})
Source Automation Rule: {{ruleName}}
EOT
      summary             = "Wiz Issue: {{control.name}}"
      table_name          = "incident"
      attach_evidence_csv = true
      custom_fields = jsonencode(
        {
          "assignment_group" : "HELPDESK",
          "category" : "Security",
          "impact" : "3 - Low",
          "subcategory" : "IDS",
          "urgency" : "1 - High"
        }
      )
    }
  }
}

# Service Now Update
resource "wiz_automation_action" "servicenow_update" {
  name = "provider_dev_update"
  type = "SERVICENOW_UPDATE_TICKET"
  servicenow_update_ticket_params {
    base_url   = var.servicenow_url
    password   = var.servicenow_password
    user       = var.servicenow_user
    table_name = "incident"
  }
}

# Jira Ticket
resource "wiz_automation_action" "jira_ticket" {
  name                          = "Jira Ticket"
  is_accessible_to_all_projects = true
  type                          = "JIRA_TICKET"
  jira_params {
    is_onprem  = false
    server_url = "https://jira.atlassian.net"
    token      = var.jira_token
    user       = "someone@example.com"
    ticket_fields {
      fix_version         = ["Q2.2022"]
      description         = "Wiz Issue"
      issue_type          = "story"
      project             = "WIZ"
      summary             = "Wiz Finding"
      attach_evidence_csv = true
    }
    tls_config {
      allow_insecure_tls = false
    }
  }
}

# Jira Ticket Transition
resource "wiz_automation_action" "jira_transition" {
  name                          = "Jira Transition"
  is_accessible_to_all_projects = true
  type                          = "JIRA_TICKET_TRANSITION"
  jira_transition_params {
    is_onprem     = false
    project       = "WIZ"
    server_url    = "https://jira.atlassian.net"
    token         = var.jira_token
    user          = "someone@example.com"
    transition_id = "Defined"
    tls_config {
      allow_insecure_tls = false
    }
  }
}

resource "wiz_automation_action" "pagerduty_create" {
  name                          = "terraform-test-pagerduty-create"
  type                          = "PAGER_DUTY_CREATE_INCIDENT"
  is_accessible_to_all_projects = true
  webhook_params {
    body = <<EOT
{
  "dedup_key": "{{id}}",
  "event_action": "trigger",
  "routing_key": "testtesttesttesttesttesttesttest",
  "payload": {
    "custom_details": [
      {
        "trigger": {
          "source": "{{triggerSource}}",
          "type": "{{triggerType}}",
          "ruleId": "{{ruleId}}",
          "ruleName": "{{ruleName}}"
        },
        "issue": {
          "id": "{{id}}",
          "status": "{{status}}",
          "severity": "{{severity}}",
          "created": "{{createdAt}}",
          "projects": "{{#projects}}{{name}}, {{/projects}}"
        },
        "resource": {
          "id": "{{entitySnapshot.providerId}}",
          "name": "{{entitySnapshot.name}}",
          "type": "{{entitySnapshot.nativeType}}",
          "cloudPlatform": "{{entitySnapshot.cloudPlatform}}",
          "subscriptionId": "{{entitySnapshot.subscriptionExternalId}}",
          "subscriptionName": "{{entitySnapshot.subscriptionName}}",
          "region": "{{entitySnapshot.region}}",
          "status": "{{entitySnapshot.status}}",
          "cloudProviderURL": "{{entitySnapshot.cloudProviderURL}}"
        },
        "control": {
          "id": "{{control.id}}",
          "name": "{{control.name}}",
          "description": "{{description}}",
          "severity": "{{severity}}"
        }
      }
    ],
    "summary": "wiz alert",
    "severity": "critical",
    "source": "wiz"
  }
}
EOT
    url  = "https://events.pagerduty.com/v2/enqueue"
  }
}

resource "wiz_automation_action" "pagerduty_resolve" {
  name                          = "terraform-test-pagerduty-resolve"
  type                          = "PAGER_DUTY_RESOLVE_INCIDENT"
  is_accessible_to_all_projects = true
  webhook_params {
    body = <<EOT
{
  "dedup_key": "{{id}}",
  "event_action": "resolve",
  "routing_key": "testtesttesttesttesttesttesttest"
}
EOT
    url  = "https://events.pagerduty.com/v2/enqueue"
  }
}
