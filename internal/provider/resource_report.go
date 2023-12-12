package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

// CreateReport struct
type CreateReport struct {
	CreateReport wiz.CreateReportPayload `json:"createReport"`
}

// UpdateReport struct
type UpdateReport struct {
	UpdateReport wiz.Report `json:"updateReport"`
}

// DeleteReport struct
type DeleteReport struct {
	DeleteReport wiz.DeleteReportPayload `json:"deleteReport"`
}

// ReadReportPayload struct -- updates
type ReadReportPayload struct {
	Report wiz.Report `json:"report"`
}

func resourceWizReportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteReport (
	    $input: DeleteReportInput!
	) {
	    deleteReport(
		input: $input
	    ) {
		_stub
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteReportInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateReport{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
