# Configure SAML Group Role Mapping on a global scope
resource "wiz_saml_group_mapping" "test_global_scope" {
  saml_idp_id       = "test-saml-identity-provider"
  provider_group_id = "global-reader-group-id"
  role              = "PROJECT_READER"
}

# Configure SAML Group Role Mapping for a single project
resource "wiz_saml_group_mapping" "test_single_project" {
  saml_idp_id       = "test-saml-identity-provider"
  provider_group_id = "admin-group-id"
  role              = "PROJECT_ADMIN"
  projects = [
    "ee25cc95-82b0-4543-8934-5bc655b86786"
  ]
}

# Configure SAML Group Role Mapping for multiple projects
resource "wiz_saml_group_mapping" "test_multi_project" {
  saml_idp_id       = "test-saml-identity-provider"
  provider_group_id = "member-group-id"
  role              = "PROJECT_MEMBER"
  projects = [
    "ee25cc95-82b0-4543-8934-5bc655b86786",
    "e7f6542c-81f6-43cf-af48-bdd77f09650d"
  ]
}
