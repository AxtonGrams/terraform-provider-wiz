package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

type CreateReport struct {
	CreateReport wiz.CreateReportPayload `json:"createReport"`
}

type UpdateReport struct {
	UpdateReport wiz.Report `json:"updateReport"`
}

type DeleteReport struct {
	DeleteReport wiz.DeleteReportPayload `json:"deleteReport"`
}

type ReadReportPayload struct {
	Report wiz.Report `json:"report"`
}

func resourceWizReportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportDelete called...")

	if d.Id() == "" {
		return nil
	}

	query := `mutation DeleteReport (
	    $input: DeleteReportInput!
	) {
	    deleteReport(
		input: $input
	    ) {
		_stub
	    }
	}`

	vars := &wiz.DeleteReportInput{}
	vars.ID = d.Id()

	data := &UpdateReport{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
