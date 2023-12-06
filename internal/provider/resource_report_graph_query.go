package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizReportGraphQuery() *schema.Resource {
	return &schema.Resource{
		Description: "TBD.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Report.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project that this report belongs to.",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The query that the report will run. Required by the GRAPH_QUERY report type.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
			"run_interval_hours": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Run interval for scheduled reports (in hours).",
			},
			"run_starts_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ISO8601 string representing the time and date when the scheduling should start (required when run_interval_hours is set).",
			},
		},
		CreateContext: resourceWizReportGraphQueryCreate,
		ReadContext:   resourceWizReportGraphQueryRead,
		UpdateContext: resourceWizReportGraphQueryUpdate,
		DeleteContext: resourceWizReportDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWizReportGraphQueryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportGraphQueryCreate called...")

	query := `mutation createReport(
	    $input: CreateReportInput!
	) {
	    createReport(
	        input: $input
	    ) {
	        report {
	            id
	        }
	    }
	}`

	vars := &wiz.CreateReportInput{}
	vars.Name = d.Get("name").(string)
	projectID, _ := d.Get("project_id").(string)
	vars.ProjectID = &projectID
	vars.Type = wiz.ReportTypeNameGraphQuery
	reportQuery := json.RawMessage(d.Get("query").(string))
	vars.GraphQueryParams = &wiz.CreateReportGraphQueryParamsInput{
		Query: reportQuery,
	}
	runIntervalHours, hasOk := d.GetOk("run_interval_hours")
	if hasOk {
		runIntervalHoursVal, _ := runIntervalHours.(int)
		vars.RunIntervalHours = &runIntervalHoursVal

		runStartsAt, hasOk := d.GetOk("run_starts_at")
		if !hasOk {
			return append(diags, diag.FromErr(fmt.Errorf("both run_interval_hours ad run_starts_at must be set for scheduling"))...)
		}

		runStartsAtVal, _ := runStartsAt.(string)
		dt, err := time.Parse(time.RFC3339, runStartsAtVal)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("run_starts_at %s is not a valid IS08601 timestamp", runStartsAtVal))...)
		}

		vars.RunStartsAt = &dt
	}

	data := &CreateReport{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	d.SetId(data.CreateReport.Report.ID)

	return resourceWizReportGraphQueryRead(ctx, d, m)
}

func resourceWizReportGraphQueryRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportGraphQueryRead called...")

	if d.Id() == "" {
		return nil
	}

	query := `query Report (
	    $id: ID!
	){
	    report(
	        id: $id
	    ) {
	        id
	        name
	        params {
		  ... on ReportParamsGraphQuery {
		    query
		    entityOptions {
		      entityType
		      propertyOptions {
		        key
		      }
		    }
		  }
		}
	        type {
		  id
		  name
		  description
		}
	        project {
	            id
	            name
	        }
		runIntervalHours
	    }
	}`

	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	tflog.Info(ctx, fmt.Sprintf("report ID during read: %s", vars.ID))

	data := &ReadReportPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.Report.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	err := d.Set("name", data.Report.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", data.Report.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("run_interval_hours", data.Report.RunIntervalHours)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("run_starts_at", data.Report.RunStartsAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	switch params := data.Report.Params.(type) {
	case wiz.ReportParamsGraphQuery:
		err = d.Set("query", params.Query)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

func resourceWizReportGraphQueryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportGraphQueryUpdate called...")

	if d.Id() == "" {
		return nil
	}

	query := `mutation UpdateReport(
	    $input: UpdateReportInput!
	) {
	    updateReport(
		input: $input
	    ) {
		report {
		    id
		    name
		}
	    }
	}`

	vars := &wiz.UpdateReportInput{}
	vars.ID = d.Id()
	vars.Override = &wiz.UpdateReportChange{}
	vars.Override.GraphQueryParams = &wiz.UpdateReportGraphQueryParamsInput{}
	reportQuery, _ := d.Get("query").(string)
	vars.Override.GraphQueryParams.Query = json.RawMessage(reportQuery)
	vars.Override.Name = d.Get("name").(string)
	runIntervalHours, hasOk := d.GetOk("run_interval_hours")
	if hasOk {
		runIntervalHoursVal, _ := runIntervalHours.(int)
		vars.Override.RunIntervalHours = &runIntervalHoursVal

		runStartsAt, hasOk := d.GetOk("run_starts_at")
		if !hasOk {
			return append(diags, diag.FromErr(fmt.Errorf("both run_interval_hours ad run_starts_at must be set for scheduling"))...)
		}

		runStartsAtVal, _ := runStartsAt.(string)
		dt, err := time.Parse(time.RFC3339, runStartsAtVal)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("run_starts_at %s is not a valid IS08601 timestamp", runStartsAtVal))...)
		}

		vars.Override.RunStartsAt = &dt
	}

	data := &UpdateReport{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizReportGraphQueryRead(ctx, d, m)
}
