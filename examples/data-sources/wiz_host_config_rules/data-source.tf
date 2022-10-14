# get the first five host configuration rules for access keys
data "wiz_host_config_rules" "access" {
  first  = 5
  search = "access"
}
