package provider

import (
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

// CreateHostConfigurationRule struct
type CreateHostConfigurationRule struct {
	CreateHostConfigurationRule wiz.CreateHostConfigurationRulePayload `json:"createHostConfigurationRule"`
}

// ReadHostConfigurationRulePayload struct -- updates
type ReadHostConfigurationRulePayload struct {
	HostConfigurationRule wiz.HostConfigurationRule `json:"hostConfigurationRule"`
}

// UpdateHostConfigurationRule struct
type UpdateHostConfigurationRule struct {
	UpdateHostConfigurationRule wiz.UpdateHostConfigurationRulePayload `json:"updateHostConfigurationRule"`
}

// DeleteHostConfigurationRule struct
type DeleteHostConfigurationRule struct {
	DeleteHostConfigurationRule wiz.DeleteHostConfigurationRulePayload `json:"deleteHostConfigurationRule"`
}
