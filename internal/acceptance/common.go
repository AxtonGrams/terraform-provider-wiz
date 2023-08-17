package acceptance

// TestCase is a string type used to identify the type of test case being run
type TestCase string

const (
	// ResourcePrefix is the default prefix used for all acceptance test resources
	ResourcePrefix = "tf-acc-test"
	// UUIDPattern is a regex pattern that matches the hexadecimal pattern commonly used by Wiz for internal IDs
	UUIDPattern = `^[A-Fa-f0-9]{8}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{12}`

	// TcCommon test case
	TcCommon TestCase = "COMMON"
	// TcUser test case
	TcUser TestCase = "USER"
	// TcServiceNow test case
	TcServiceNow TestCase = "SERVICE_NOW"
	// TcJira test case
	TcJira TestCase = "JIRA"
	// TcSubscriptionResourceGroups test case
	TcSubscriptionResourceGroups TestCase = "SUBSCRIPTION_RESOURCE_GROUPS"
	// TcProject test case
	TcProject TestCase = "PROJECT"
	// TcCloudConfigRule test case
	TcCloudConfigRule TestCase = "CLOUD_CONFIG_RULE"
)
