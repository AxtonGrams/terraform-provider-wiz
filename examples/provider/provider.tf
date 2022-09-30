terraform {
  required_providers {
    wiz = {
      source  = "AxtonGrams/wiz"
      version = "1.0.2"
    }
  }
}

provider "wiz" {
  wiz_url                = var.wiz_url
  wiz_auth_client_id     = var.wiz_auth_client_id
  wiz_auth_client_secret = var.wiz_auth_client_secret
  wiz_auth_audience      = "wiz-api"
}
