# Provision a simple AWS connector, opting for a single region
resource "wiz_connector_aws" "example" {
  name = "example"
  auth_params = jsonencode({
    "customerRoleARN" : "arn:aws:iam::100000000009:role/wiz-customer",
  })

  extra_config = jsonencode(
    {
      "skipOrganizationScan" : true,
      "diskAnalyzerInFlightDisabled" : false,
      "optedInRegions" : ["us-east-1"],
      "excludedAccounts" : [],
      "excludedOUs" : [],
      "auditLogMonitorEnabled" : false
    }
  )
}

# Provision an AWS connector with Outpost that uses a custom config
resource "wiz_connector_aws" "example" {
  name = "example"
  auth_params = jsonencode({
    "customerRoleARN" : "arn:aws:iam::100000000009:role/wiz-customer",
    "outpostId" : "078862d0-a62f-406c-b966-13445af34c0d",
    "diskAnalyzer" : {
      "scanner" : {
        "roleARN" : "arn:aws:iam::100000000009:role/outpost-scanner"
      }
    },
  })

  extra_config = jsonencode(
    {
      "auditLogMonitorEnabled" : false,
      "excludedAccounts" : ["100000000009", "100000000010", "100000000013"],
      "excludedOUs" : ["EXCLUDE-ME"],
      "auditLogMonitorEnabled" : false,
      "diskAnalyzerInFlightDisabled" : false,
      "skipOrganizationScan" : true,
      "optedInRegions" : [],
      "cloudTrailConfig" : {
        "bucketName" : "buckethere",
        "bucketSubAccount" : "000000000012",
        "trailOrg" : "o-myorg"
      }
    }
  )

}
