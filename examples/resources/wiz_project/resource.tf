# A simple example
resource "wiz_project" "test" {
  name        = "Test App"
  description = "My project description"
  risk_profile {
    business_impact = "MBI"
  }
  business_unit = "Technology"
}

# Folder projects example
resource "wiz_project" "root" {
  name        = "root"
  description = "root"
  is_folder   = true

  risk_profile {
    business_impact = "MBI"
  }
  business_unit = "Technology"
}

resource "wiz_project" "child" {
  name              = "project_with_accounts"
  parent_project_id = wiz_project.root.id
  risk_profile {
    business_impact = "MBI"
  }
  business_unit = "Technology"
  cloud_account_link {
    cloud_account_id = "477ea00a-4d4d-5bb4-9fa6-634691e68fff"
    environment      = "PRODUCTION"
  }
}

# This resource contains multiple organization links, one with tags and another without
resource "wiz_project" "test" {
  name        = "Test App"
  description = "My project description"
  risk_profile {
    business_impact = "MBI"
  }
  business_unit = "Technology"
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

# This resource contains a single cloud account link, with tag
resource "wiz_project" "test" {
  name        = "Test App"
  description = "My project description"
  risk_profile {
    business_impact = "MBI"
  }
  business_unit = "Technology"
  cloud_account_link {
    cloud_account_id = "3225def3-0e0e-5cb8-955a-3583f696f778"
    environment      = "PRODUCTION"
    resource_tags {
      key   = "created_by"
      value = "terraform"
    }
  }
}

# This resource contains a single kubernetes cluster link
resource "wiz_project" "test" {
  name        = "My Kubernetes Project"
  description = "My project description"
  risk_profile {
    business_impact = "MBI"
  }
  business_unit = "Technology"
  kubernetes_cluster_link {
    kubernetes_cluster = "77de7ca1-02f9-5ed2-a94b-5d19c683efaf"
    environment        = "STAGING"
    shared             = true
    namespaces         = ["kube-system"]
  }
}
