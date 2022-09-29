# Get the Wiz internal information for the Organization root based on the AWS Root ID

data "wiz_organizations" "root" {
  search = "r-1234"
}
