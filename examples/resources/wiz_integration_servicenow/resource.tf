resource "wiz_integration_servicenow" "default" {
  name                = "default"
  servicenow_url      = var.servicename_url
  servicenow_username = var.servicenow_username
  servicenow_password = var.servicenow_password
  scope               = "All Resources, Restrict this Integration to global roles only"
}
