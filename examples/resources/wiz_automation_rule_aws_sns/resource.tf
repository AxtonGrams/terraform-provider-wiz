# Provision an AWS SNS integration
resource "wiz_integration_aws_sns" "example" {
  name                      = "example"
  aws_sns_topic_arn         = "arn:aws:sns:us-east-1:123456789012:Example"
  aws_sns_access_method     = "ASSUME_SPECIFIED_ROLE"
  aws_sns_customer_role_arn = "arn:aws:iam::123456789012:role/Example-Role"
  scope                     = "All Resources, Restrict this Integration to global roles only"
}

# Provision an AWS SNS automation rule
resource "wiz_automation_rule_aws_sns" "example" {
  name           = "example"
  description    = "example description"
  enabled        = true
  integration_id = wiz_integration_aws_sns.example.id
  trigger_source = "ISSUES"
  trigger_type = [
    "CREATED",
    "REOPENED",
  ]
  aws_sns_body = jsonencode({
    "trigger" : {
      "source" : "{{triggerSource}}",
      "type" : "{{triggerType}}",
      "ruleId" : "{{ruleId}}",
      "ruleName" : "{{ruleName}}"
    },
    "issue" : {
      "id" : "{{issue.id}}",
      "status" : "{{issue.status}}",
      "severity" : "{{issue.severity}}",
      "created" : "{{issue.createdAt}}",
      "projects" : "{{#issue.projects}}{{name}}, {{/issue.projects}}"
    },
    "resource" : {
      "id" : "{{issue.entitySnapshot.providerId}}",
      "name" : "{{issue.entitySnapshot.name}}",
      "type" : "{{issue.entitySnapshot.nativeType}}",
      "cloudPlatform" : "{{issue.entitySnapshot.cloudPlatform}}",
      "subscriptionId" : "{{issue.entitySnapshot.subscriptionExternalId}}",
      "subscriptionName" : "{{issue.entitySnapshot.subscriptionName}}",
      "region" : "{{issue.entitySnapshot.region}}",
      "status" : "{{issue.entitySnapshot.status}}",
      "cloudProviderURL" : "{{issue.entitySnapshot.cloudProviderURL}}"
    },
    "control" : {
      "id" : "{{issue.control.id}}",
      "name" : "{{issue.control.name}}",
      "description" : "{{issue.control.description}}",
      "severity" : "{{issue.control.severity}}",
      "sourceCloudConfigurationRuleId" : "{{issue.control.sourceCloudConfigurationRule.shortId}}",
      "sourceCloudConfigurationRuleName" : "{{issue.control.sourceCloudConfigurationRule.name}}"
    }
  })
  filters = jsonencode({
    "project" : [],
    "relatedEntity" : {
      "cloudPlatform" : [
        "AWS"
      ],
      "subscriptionId" : [
        "fccc3f07-3304-4f9d-ac2d-a43dd6128eb0",
        "a005e165-49c5-41b7-befb-a0e4d866fc6c",
      ]
    },
    "sourceControl" : [
      "b46c34d2-3624-4e1e-bb04-dda5177582c7",
      "6c27d70a-7329-42e9-b19e-0b974f556365",
    ]
  })
}
