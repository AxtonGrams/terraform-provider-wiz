# Provision an AWS SNS integration with a specified role
resource "wiz_integration_aws_sns" "specified_role_all_projects" {
  name                      = "test-terraform-001"
  aws_sns_topic_arn         = "arn:aws:sns:us-east-1:123456789012:RemediationTopic"
  aws_sns_access_method     = "ASSUME_SPECIFIED_ROLE"
  aws_sns_customer_role_arn = "arn:aws:iam::123456789012:role/RemediationRole"
}

# Provision and AWS SNS integration with the connector role
resource "wiz_integration_aws_sns" "connector_role_all_projects" {
  name                  = "test-terraform-003"
  aws_sns_topic_arn     = "arn:aws:sns:us-east-1:123456789012:RemediationTopic"
  aws_sns_access_method = "ASSUME_CONNECTOR_ROLE"
  aws_sns_connector_id  = "ab48ad5e-44fb-48f8-9899-24ee4ed974c1"
}

# Provision and AWS SNS integration that uses the connector role role for a specified project
resource "wiz_integration_aws_sns" "specified_role_single_project" {
  name                  = "test-terraform-004"
  aws_sns_topic_arn     = "arn:aws:sns:us-east-1:981012938874:Wiz-Remediation-Issues-Topic"
  aws_sns_access_method = "ASSUME_CONNECTOR_ROLE"
  aws_sns_connector_id  = "ef0bd8a5-165b-4498-b5d7-19871f762c21"
  scope                 = "Selected Project"
  project_id            = "1091ae77-116a-56cf-990e-db2f4f691f66"
}
