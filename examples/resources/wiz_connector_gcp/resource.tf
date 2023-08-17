# Provision a simple GCP connector, organization-wide
resource "wiz_connector_gcp" "example" {
  name = "example"
  auth_params = jsonencode({
    "isManagedIdentity" : true,
    "organization_id" : "o-example"
  })

  extra_config = jsonencode(
    {
      "projects" : [],
      "excludedProjects" : [],
      "includedFolders" : [],
      "excludedFolders" : [],
      "diskAnalyzerInFlightDisabled" : false,
      "auditLogMonitorEnabled" : false
    }
  )
}

# Provision a GCP connector targeting an individual Google project
resource "wiz_connector_gcp" "example" {
  name = "example"
  auth_params = jsonencode({
    "isManagedIdentity" : true,
    "project_id" : "exmaple-project-id"
  })

  extra_config = jsonencode(
    {
      "projects" : [],
      "excludedProjects" : [],
      "includedFolders" : [],
      "excludedFolders" : [],
      "diskAnalyzerInFlightDisabled" : false,
      "auditLogMonitorEnabled" : false
    }
  )
}
