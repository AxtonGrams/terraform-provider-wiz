# Create a project reader service account
resource "wiz_service_account" "project_reader" {
  name = "project_reader"
  scopes = [
    "read:projects",
  ]
}

# Create a helm (broker) service account
resource "wiz_service_account" "helm" {
  name = "helm"
  type = "BROKER"
}
