# get aws cloud configuration rules for access keys
data "wiz_cloud_config_rules" "aws_access_key" {
  search = "Access key"
  cloud_provider = [
    "AWS",
  ]
}

# get high and critical aws cloud configuration rules that have remediation
data "wiz_cloud_config_rules" "aws_critical" {
  cloud_provider = [
    "AWS",
  ]
  severity = [
    "CRITICAL",
    "HIGH",
  ]
  has_remediation = true
}
