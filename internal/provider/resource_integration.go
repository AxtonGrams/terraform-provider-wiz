package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

// CreateIntegration struct
type CreateIntegration struct {
	CreateIntegration wiz.CreateIntegrationPayload `json:"createIntegration"`
}

// ReadIntegrationPayload struct
type ReadIntegrationPayload struct {
	Integration wiz.Integration `json:"integration"`
}

// UpdateIntegration struct
type UpdateIntegration struct {
	UpdateIntegration wiz.UpdateIntegrationPayload `json:"updateIntegration"`
}

// DeleteIntegration struct
type DeleteIntegration struct {
	DeleteIntegration wiz.DeleteIntegrationPayload `json:"deleteIntegration"`
}

// resourceWizIntegrationDelete deletes a Wiz integration resource
func resourceWizIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsSNSDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteIntegration (
	  $input: DeleteIntegrationInput!
	) {
	  deleteIntegration(
	    input: $input
	  ) {
	    _stub
	  }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteIntegrationInput{}
	vars.ID = d.Id()

	// process the request
	data := &DeleteIntegration{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "integration", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}

// convertIntegrationScopeToBool converts the literal string representation of the 'scope' to the boolean expected by Wiz
func convertIntegrationScopeToBool(integrationScope string) *bool {
	var value bool

	switch integrationScope {
	case "Select Project":
		value = false
	case "All Resources":
		value = true
	}

	return &value
}
