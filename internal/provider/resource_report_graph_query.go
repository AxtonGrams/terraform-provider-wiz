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

const reportRunStartsAtLayout = "2006-01-02 15:04:05 +0000 UTC"

func resourceWizReportGraphQuery() *schema.Resource {
	return &schema.Resource{
		Description: "A GraphQL Query Report is an automated query that can be scheduled to run at hourly intervals.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Report.",
			},
			"project_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "*",
				Description: "The ID of the project that this report belongs to (changing this requires re-creatting the report). Defaults to all projects.",
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
				Type:     schema.TypeString,
				Optional: true,
				Description: fmt.Sprintf(
					"String representing the time and date when the scheduling should start (required when run_interval_hours is set). Must be in the following format: %s. Also, Wiz will always round this down by the hour.",
					reportRunStartsAtLayout,
				),
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

func setScheduling(diags diag.Diagnostics, d *schema.ResourceData, vars interface{}) diag.Diagnostics {
	runIntervalHours, hasOk := d.GetOk("run_interval_hours")
	if !hasOk {
		return nil
	}

	runIntervalHoursVal, _ := runIntervalHours.(int)
	runStartsAt, hasOk := d.GetOk("run_starts_at")
	if !hasOk {
		return append(diags, diag.FromErr(fmt.Errorf("both run_interval_hours ad run_starts_at must be set to enable scheduling"))...)
	}

	runStartsAtVal, _ := runStartsAt.(string)
	dt, err := time.Parse(reportRunStartsAtLayout, runStartsAtVal)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("run_starts_at %s does not match layout %s", runStartsAtVal, reportRunStartsAtLayout))...)
	}

	switch vars := vars.(type) {
	case *wiz.CreateReportInput:
		vars.RunIntervalHours = &runIntervalHoursVal
		vars.RunStartsAt = &dt
	case *wiz.UpdateReportInput:
		vars.Override.RunIntervalHours = &runIntervalHoursVal
		vars.Override.RunStartsAt = &dt
	default:
		return append(diags, diag.FromErr(fmt.Errorf("vars is an invalid ReportInput type"))...)
	}

	return nil
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
			runStartsAt
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

	if diags := setScheduling(diags, d, vars); diags != nil {
		return diags
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
		runStartsAt
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
	projectID := "*"
	if data.Report.Project != nil {
		projectID = data.Report.Project.ID
	}

	err = d.Set("project_id", projectID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if data.Report.RunIntervalHours != nil {
		err = d.Set("run_interval_hours", data.Report.RunIntervalHours)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if data.Report.RunStartsAt != nil {
		runStartsAt := data.Report.RunStartsAt.Format(reportRunStartsAtLayout)
		err = d.Set("run_starts_at", runStartsAt)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
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
			runStartsAt
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

	if diags := setScheduling(diags, d, vars); diags != nil {
		return diags
	}

	data := &UpdateReport{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizReportGraphQueryRead(ctx, d, m)
}
