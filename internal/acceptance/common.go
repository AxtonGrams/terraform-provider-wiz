package acceptance

const (
	// ResourcePrefix is the default prefix used for all acceptance test resources
	ResourcePrefix = "tf-acc-test"
	// UUIDPattern is a regex pattern that matches the hexadecimal pattern commonly used by Wiz for internal IDs
	UUIDPattern = `^[A-Fa-f0-9]{8}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{4}\-[A-Fa-f0-9]{12}`
)
