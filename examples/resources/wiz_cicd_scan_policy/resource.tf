resource "wiz_cicd_scan_policy" "iac" {
  name        = "terraform-test-iac"
  description = "terraform-test-iac description"
  iac_params {
    count_threshold             = 3
    severity_threshold          = "CRITICAL"
    builtin_ignore_tags_enabled = false
    ignored_rules = [
      "fd7dd0c6-4953-4b36-bc39-004ec3d870db",
      "063fb380-9eda-4c08-a31b-9211ee37bd42",
    ]
    custom_ignore_tags {
      key              = "testkey1"
      value            = "testval1"
      ignore_all_rules = false
      rule_ids = [
        "063fb380-9eda-4c08-a31b-9211ee37bd42",
      ]
    }
    custom_ignore_tags {
      key              = "testkey2"
      value            = "testval2"
      ignore_all_rules = false
      rule_ids = [
        "1f0ee3b5-5404-4b40-bbc8-33a990330ac3",
        "a1958aa1-b810-4df6-bd82-487cb37c6039",
      ]
    }
  }
}
resource "wiz_cicd_scan_policy" "secrets" {
  name        = "terraform-test-secrets2"
  description = "terraform-test-secrets description"
  disk_secrets_params {
    count_threshold = 3
    path_allow_list = [
      "/etc",
      "/opt",
      "/root",
    ]
  }
}

resource "wiz_cicd_scan_policy" "vulnerabilities" {
  name        = "terraform-test-vulnerabilities"
  description = "terraform-test-vulnerabilities description"
  disk_vulnerabilities_params {
    ignore_unfixed = true
    package_allow_list = [
      "lsof",
      "sudo",
      "apt",
    ]
    package_count_threshold = 3
    severity                = "LOW"
  }
}
