resource "wiz_user" "psm" {
  for_each = local.wiz_local_users
  email    = var.wiz_local_users[each.key].email
  name     = each.key
  role     = var.wiz_local_users[each.key].role
}
