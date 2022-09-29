# This resource contains multiple organization links, one with tags and another without

resource "wiz_project" "test" {
  name        = "Test App"
  description = "My project description"
  risk_profile {
    business_impact = "MBI"
  }
  business_unit = data.insight_organization.aws.description
  cloud_organization_link {
    cloud_organization = "7edbb879-9960-513f-b56d-876e9db2a962"
    environment        = "PRODUCTION"
    shared             = false
  }
  cloud_organization_link {
    cloud_organization = "07401938-0347-5a02-80eb-db19eecfbf98"
    environment        = "PRODUCTION"
    shared             = true
    resource_tags {
      key   = "application"
      value = "Wiz"
    }
    resource_tags {
      key   = "environment"
      value = "production"
    }
  }
}

# A simple example

resource "wiz_project" "test" {
  name        = "Test App"
  description = "My project description"
  risk_profile {
    business_impact = "MBI"
  }
  business_unit = "Technology"
}
