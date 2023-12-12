# Importing Considerations:
#
# Please note this is considered experimental, exercise caution and consider the following:
#
# - Make sure that the `auth_params` field is set to the same values as set when the resource was created outside of Terraform.
#   This is due to the way we need to handle change as under normal diff conditions, `auth_params` requires a resource recreation.
#
# - For `auth_params` include `isManagedIdentity`. If using outposts, also include `outPostId` and `diskAnalyzer` structure.
#
# For more information, refer to the examples in the documentation.
#
terraform import wiz_connector_gcp.import_example "7be792ba-bfd1-46d0-9fba-5f6bc19df4a8"

# Optional - this is to set auth_params in state.
#
# If not run post-import, the next `terraform apply` will take care of it.
# Note any speculative changes to `auth_params` are for setting state for the one-time import only, any further changes would require a resource recreation as normal.
terraform apply --target=wiz_connector_gcp.import_example