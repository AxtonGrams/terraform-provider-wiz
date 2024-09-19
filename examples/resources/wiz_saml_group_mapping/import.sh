# The id for importing resources has to be in this format: 'mapping|<saml_idp_id>|<provider_group_id>:<project_ids>:<role>#...'.
# Import with saml mapping to multiple projects
terraform import wiz_saml_group_mapping.example_import "mapping|wiz-azure-ad-saml|88990357-fe36-421b-aedc-fcdd602b91d7:bb62aac7-e8bd-5d5e-b205-2dbafe106e1a,ee25cc95-82b0-4543-8934-5bc655b86786:PROJECT_READER"

# Import with mapping to single project
terraform import wiz_saml_group_mapping.example_import "mapping|wiz-azure-ad-saml|88990357-fe36-421b-aedc-fcdd602b91d7:bb62aac7-e8bd-5d5e-b205-2dbafe106e1a:PROJECT_READER"

# Import with global mapping
terraform import wiz_saml_group_mapping.example_import "mapping|wiz-azure-ad-saml|88990357-fe36-421b-aedc-fcdd602b91d7::PROJECT_READER"

# Import with multiple group mappings
terraform import wiz_saml_group_mapping.example_import "mapping|wiz-azure-ad-saml|88990357-fe36-421b-aedc-fcdd602b91d7:bb62aac7-e8bd-5d5e-b205-2dbafe106e1a:PROJECT_READER#12345678-1234-1234-1234-123456789012:ee25cc95-82b0-4543-8934-5bc655b86786:PROJECT_WRITER"