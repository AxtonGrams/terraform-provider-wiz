# retrieve account by aws account id
data "wiz_cloud_accounts" "accounts_by_id" {
  search = [
    "012345678912",
    "987654321098",
  ]
}

# retrieve one account by wiz internal identifier
data "wiz_cloud_accounts" "accounts_by_wiz_id" {
  ids = [
    "d33a2072-4b95-481b-8153-c0b9089992aa",
  ]
}

# retrieve all ccounts with multiple source connectors
data "wiz_cloud_accounts" "multiple_connectors" {
  has_multiple_connector_sources = true
}
