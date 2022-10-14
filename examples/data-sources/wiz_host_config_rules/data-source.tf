# get host configuration rules for access keys
data "wiz_host_config_rules" "aws_access_key" {
  search = "Access key"
}
