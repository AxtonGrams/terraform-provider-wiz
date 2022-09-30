variable "wiz_url" {
  type        = string
  description = "Wiz api endpoint. This varies for each Wiz deployment. See https://docs.wiz.io/wiz-docs/docs/using-the-wiz-api#the-graphql-endpoint"
}

variable "wiz_auth_client_id" {
  type        = string
  description = "Your application's Client ID. You can find this value on the Settings > Service Accounts page."
}

variable "wiz_auth_client_secret" {
  type        = string
  description = "Your application's Client Secret. You can find this value on the Settings > Service Accounts page."
  sensitive   = true
}

variable "wiz_auth_audience" {
  type        = string
  description = "Use 'beyond-api' if using auth0, otherwise use 'wiz-api'"
  default     = "wiz-api"
}

terraform {
  required_providers {
    wiz = {
      source  = "wiz.io/hashicorp/wiz"
      version = "1.0.0"
    }
  }
}

provider "wiz" {
  wiz_url                = var.wiz_url
  wiz_auth_client_id     = var.wiz_auth_client_id
  wiz_auth_client_secret = var.wiz_auth_client_secret
  wiz_auth_audience      = var.wiz_auth_audience
}
