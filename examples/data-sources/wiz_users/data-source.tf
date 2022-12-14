# Get Wiz user(s) based on an email address
data "wiz_users" "by_email" {
  search = "johnny@domain.com"
}

# Get first 4 Wiz user(s) based on role
data "wiz_users" "by_role" {
  roles = ["GLOBAL_READER"]
  first = 4
}
