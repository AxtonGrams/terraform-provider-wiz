resource "wiz_service_account" "project_reader" {
  name = "project_reader"
  scopes = [
    "read:projects",
  ]
}
