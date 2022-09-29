resource "wiz_cloud_config_rule" "test" {
  name        = "terraform-test-iac"
  description = "test description"
  target_native_types = [
    "account",
  ]
  security_sub_categories = [
    "wsct-id-7",
    "wsct-id-3759",
  ]
  scope_account_ids = [
    "87d3d654-41ba-53d2-aeba-8399be9f62d7",
    #"e295ffc1-636a-5f8f-959a-4bf9fdba9727",
  ]
  function_as_control      = false
  remediation_instructions = "fix it"
  enabled                  = false
  severity                 = "HIGH"
  opa_policy               = <<EOT
package wiz

default result = "pass"
EOT
  iac_matchers {
    type      = "CLOUD_FORMATION"
    rego_code = <<EOT
package wiz

import data.generic.cloudformation as cloudFormationLib

import data.generic.common as common_lib

WizPolicy[result] {
        resource := input.document[i].Resources[name]
        resource.Type == "AWS::Config::ConfigRule"
        not hasAccessKeyRotationRule(resource)

        result := {
                "documentId": input.document[i].id,
                "searchKey": sprintf("Resources.%s", [name]),
                "issueType": "MissingAttribute",
                "keyExpectedValue": sprintf("Resources.%s has a ConfigRule defining rotation period on AccessKeys.", [name]),
                "keyActualValue": sprintf("Resources.%s doesn't have a ConfigRule defining rotation period on AccessKeys.", [name]),
                "resourceTags": cloudFormationLib.getCFTags(resource),
        }
}

hasAccessKeyRotationRule(configRule) {
        configRule.Properties.Source.SourceIdentifier == "ACCESS_KEYS_ROTATED"
} else = false {
        true
}
EOT
  }
}
