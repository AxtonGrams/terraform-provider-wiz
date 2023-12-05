package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizReport() *schema.Resource {
	return &schema.Resource{
		Description: "A Report consists of a pre-defined Security Graph query and a severity level—if a Report's query returns any results, an Issue is generated for every result. Each Report is assigned to a category in one or more Policy Frameworks.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Report.",
			},
			"column_selection": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom columns to include in the report. When disabled, the default columns will be included.",
			},
			"csv_delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Delimiter to use for the CSV output.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.ReportCsvDelimiters,
						false,
					),
				),
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project that this report belongs to.",
				Default:     false,
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The query that the report will run. Required by the GRAPH_QUERY report type.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
			"graph_query_params": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The query that represents the control's scope.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringIsJSON,
				),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Type for this report.\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.ReportTypes,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.ReportTypes,
						false,
					),
				),
			},
		},
		CreateContext: resourceWizProjectCreate,
		ReadContext:   resourceWizProjectRead,
		UpdateContext: resourceWizProjectUpdate,
		DeleteContext: resourceWizProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateReport struct
type CreateReport struct {
	CreateReport *wiz.Report `json:"createReport"`
}

func resourceWizReportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportCreate called...")

	// define the graphql query
	query := `mutation createReport(
	    $input: CreateReportInput!
	) {
	    createReport(
	        input: $input
	    ) {
	        control {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.CreateReportInput{}
	vars.Name = d.Get("name").(string)

	csvDelimiter, _ := d.Get("csv_delimiter").(wiz.CSVDelimiter)
	vars.CSVDelimiter = &csvDelimiter
	columnSelection, hasVal := d.GetOk("column_selection")
	if hasVal {
		vars.ColumnSelection = columnSelection.([]string)
	}

	projectID, _ := d.Get("project_id").(string)
	vars.ProjectID = &projectID

	// process the request
	data := &CreateReport{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateReport.ID)

	return resourceWizReportRead(ctx, d, m)
}

// ReadReportPayload struct -- updates
type ReadReportPayload struct {
	Report wiz.Report `json:"report"`
}

func resourceWizReportRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query Report (
	    $id: ID!
	){
	    control(
	        id: $id
	    ) {
	        id
	        name
	        params {
		  query
		  entityOptions {
		    entityType
		    propertyOptions {
		      key
		    }
		  }
		}
	        type
	        project {
	            id
	            name
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: oops! an internal error has occurred. for reference purposes, this is your request id
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

	// set the resource parameters
	err := d.Set("name", data.Report.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("type", data.Report.Type)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("project_id", data.Report.Project.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateReport struct
type UpdateReport struct {
	UpdateReport wiz.Report `json:"updateReport"`
}

func resourceWizReportUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizReportUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateReport(
	    $input: UpdateReportInput!
	) {
	    updateReport(
		input: $input
	    ) {
		control {
		    id
		}
	    }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateReportInput{}
	vars.ID = d.Id()

	if d.HasChange("name") {
		vars.Override.Name = d.Get("name").(string)
	}

	// process the request
	data := &UpdateReport{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "report", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizReportRead(ctx, d, m)
}

// DeleteReport struct
type DeleteReport struct {
	DeleteReport wiz.DeleteReportPayload `json:"deleteReport"`
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
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
