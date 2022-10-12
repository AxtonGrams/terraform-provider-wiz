package provider

import (
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

// CreateHostConfigurationRule struct
type CreateHostConfigurationRule struct {
	CreateHostConfigurationRule vendor.CreateHostConfigurationRulePayload `json:"createHostConfigurationRule"`
}

// ReadHostConfigurationRulePayload struct -- updates
type ReadHostConfigurationRulePayload struct {
	HostConfigurationRule vendor.HostConfigurationRule `json:"cloudConfigurationRule"`
}

// UpdateHostConfigurationRule struct
type UpdateHostConfigurationRule struct {
	UpdateHostConfigurationRule vendor.UpdateHostConfigurationRulePayload `json:"updateHostConfigurationRule"`
}

// DeleteHostConfigurationRule struct
type DeleteHostConfigurationRule struct {
	DeleteHostConfigurationRule vendor.DeleteHostConfigurationRulePayload `json:"deleteHostConfigurationRule"`
}
