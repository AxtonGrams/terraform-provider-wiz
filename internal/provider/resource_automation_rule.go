package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

// CreateAutomationRule struct
type CreateAutomationRule struct {
	CreateAutomationRule vendor.CreateAutomationRulePayload `json:"createAutomationRule"`
}

// ReadAutomationRulePayload struct -- updates
type ReadAutomationRulePayload struct {
	AutomationRule vendor.AutomationRule `json:"automationRule"`
}

// UpdateAutomationRule struct
type UpdateAutomationRule struct {
	UpdateAutomationRule vendor.UpdateAutomationRulePayload `json:"updateAutomationRule"`
}

// DeleteAutomationRule struct
type DeleteAutomationRule struct {
	DeleteAutomationRule vendor.DeleteAutomationRulePayload `json:"deleteAutomationRule"`
}

func resourceWizAutomationRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationRuleDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteAutomationRule (
            $input: DeleteAutomationRuleInput!
        ) {
            deleteAutomationRule (
                input: $input
            ) {
                _stub
            }
        }`

	// populate the graphql variables
	vars := &vendor.DeleteAutomationRuleInput{}
	vars.ID = d.Id()

	// process the request
	data := &DeleteAutomationRule{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_rule", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
