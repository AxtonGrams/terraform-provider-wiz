# A link from a project to a cloud account can be created using the accounts id in wiz
resource "wiz_project_cloud_account_link" "example" {
  project_id       = "ee25cc95-82b0-4543-8934-5bc655b86786"
  cloud_account_id = "5cc3a684-44cb-4cd5-b78f-f029c25dc617"
  environment      = "PRODUCTION"
}

# Or using the external id of the cloud account
resource "wiz_project_cloud_account_link" "example" {
  project_id                = "ee25cc95-82b0-4543-8934-5bc655b86786"
  external_cloud_account_id = "04e56587-4408-402a-9c8c-f454ed45da65"
  environment               = "PRODUCTION"
}

# Both can be supplied but they have to belong to the same account
resource "wiz_project_cloud_account_link" "example" {
  project_id                = "ee25cc95-82b0-4543-8934-5bc655b86786"
  cloud_account_id          = "5cc3a684-44cb-4cd5-b78f-f029c25dc617"
  external_cloud_account_id = "04e56587-4408-402a-9c8c-f454ed45da65"
  environment               = "PRODUCTION"
}
