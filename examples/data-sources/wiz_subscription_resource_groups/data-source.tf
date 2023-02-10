# Get the first 3 resource groups for an Azure subscription ID

data "wiz_subscription_resource_groups" "rgs" {
  subscription_id = "1689bd5b-4df3-5dc8-9046-2f0a15faa62f"
  first           = 3
}