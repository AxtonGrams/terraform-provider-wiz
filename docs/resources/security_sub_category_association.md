---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wiz_security_sub_category_association Resource - terraform-provider-wiz"
subcategory: ""
description: |-
  Manage associations between security sub-categories and policies. This resource can only be used with custom security sub-categories. Wiz managed or custom policies can be referenced. When the association is removed from state, all associations managed by this resource will be removed. Associations managed outside this resouce declaration will remain untouched through the lifecycle of this resource.
---

# wiz_security_sub_category_association (Resource)

Manage associations between security sub-categories and policies. This resource can only be used with custom security sub-categories. Wiz managed or custom policies can be referenced. When the association is removed from state, all associations managed by this resource will be removed. Associations managed outside this resouce declaration will remain untouched through the lifecycle of this resource.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `security_sub_category_id` (String) Security sub-category ID.

### Optional

- `cloud_config_rule_ids` (List of String) List of cloud config rule IDs.
    - Required exactly one of: `[cloud_config_rule_ids control_ids host_config_rule_ids]`.
- `control_ids` (List of String) List of control IDs.
    - Required exactly one of: `[cloud_config_rule_ids control_ids host_config_rule_ids]`.
- `details` (String) Details of the association. This information is not used to manage resources, but can serve as notes for the associations.
- `host_config_rule_ids` (List of String) List of host config rule IDs.
    - Required exactly one of: `[cloud_config_rule_ids control_ids host_config_rule_ids]`.

### Read-Only

- `id` (String) Internal identifier for the association.

