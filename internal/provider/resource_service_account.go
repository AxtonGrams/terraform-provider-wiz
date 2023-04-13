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

func resourceWizServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Services accounts are used to integrate with Wiz.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Wiz internal identifier.",
				Computed:    true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "THIRD_PARTY",
				ForceNew: true,
				Description: fmt.Sprintf(
					"Service account type, for Helm use `BROKER` type.`\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						wiz.ServiceAccountType,
					),
				),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						wiz.ServiceAccountType,
						false,
					),
				),
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Description: fmt.Sprintf(
					"Scopes, required with THIRD_PARTY (GraphQL API type).\n    - Allowed values: %s",
					utils.SliceOfStringToMDUList(
						internal.ServiceAccountScopes,
					),
				),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice(
							internal.ServiceAccountScopes,
							false,
						),
					),
				},
			},
			"assigned_projects": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Project ID assignments, optional with THIRD_PARTY (GraphQL API type)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.IsUUID,
					),
				},
			},
			"last_rotated_at": {
				Type:        schema.TypeString,
				Description: "If a change is detected with this value, the service account will be recreated to ensure a valid secret is stored in Terraform state.",
				Computed:    true,
				ForceNew:    true,
			},
			"recreate_if_rotated": {
				Type:        schema.TypeBool,
				Description: "Recreate the resource if rotated outside Terraform? This can be used to ensure the state contains valid authentication information. This option should be disabled if external tools are used to manage the credentials for this service account.",
				Optional:    true,
				Default:     false,
			},
		},
		CreateContext: resourceWizServiceAccountCreate,
		ReadContext:   resourceWizServiceAccountRead,
		UpdateContext: resourceWizServiceAccountUpdate,
		DeleteContext: resourceWizServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateServiceAccount struct
type CreateServiceAccount struct {
	CreateServiceAccount wiz.CreateServiceAccountPayload `json:"createServiceAccount"`
}

func resourceWizServiceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizServiceAccountCreate called...")

	// define the graphql query
	query := `mutation CreateServiceAccount($input: CreateServiceAccountInput!) {
	    createServiceAccount(input: $input) {
	        serviceAccount {
	            id
	            name
	            clientId
	            clientSecret
	            scopes
	            type
	            createdAt
	            assignedProjects {
	                id
	            }
	            lastRotatedAt
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.CreateServiceAccountInput{}
	vars.Name = d.Get("name").(string)
	t := d.Get("type").(string)
	vars.Type = &t
	if t == "THIRD_PARTY" {
		vars.Scopes = utils.ConvertListToString(d.Get("scopes").([]interface{}))
		vars.AssignedProjectIDs = utils.ConvertListToString(d.Get("assigned_projects").([]interface{}))
	}

	// process the request
	data := &CreateServiceAccount{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "service_account", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id and computed values
	d.SetId(data.CreateServiceAccount.ServiceAccount.ID)
	d.Set("client_secret", data.CreateServiceAccount.ServiceAccount.ClientSecret)
	d.Set("last_rotated_at", data.CreateServiceAccount.ServiceAccount.LastRotatedAt)
	d.Set("client_id", data.CreateServiceAccount.ServiceAccount.ClientID)
	d.Set("created_at", data.CreateServiceAccount.ServiceAccount.CreatedAt)

	return resourceWizServiceAccountRead(ctx, d, m)
}

// ReadServiceAccountPayload struct -- updates
type ReadServiceAccountPayload struct {
	ServiceAccount wiz.ServiceAccount `json:"serviceAccount,omitempty"`
}

func resourceWizServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizServiceAccountRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query ServiceAccount  (
	    $id: ID!
	) {
	    serviceAccount(
	        id: $id
	    ) {
	        id
	        name
	        clientId
	        clientSecret
	        scopes
	        type
	        createdAt
	        assignedProjects {
	            id
	        }
	        lastRotatedAt
	        }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	data := &ReadServiceAccountPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "service_account", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	err := d.Set("name", data.ServiceAccount.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("scopes", data.ServiceAccount.Scopes)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("type", data.ServiceAccount.Type)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("last_rotated_at", data.ServiceAccount.LastRotatedAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("created_at", data.ServiceAccount.CreatedAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("recreate_if_rotated", d.Get("recreate_if_rotated").(bool))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// if key was rotated outside terraform, trigger a recreation
	// since terraform won't force a new resource for computed values, we change the name
	lraOld, _ := d.GetChange("last_rotated_at")
	tflog.Debug(ctx, fmt.Sprintf("old/new: %s/%s", lraOld.(string), data.ServiceAccount.LastRotatedAt))
	if lraOld.(string) != data.ServiceAccount.LastRotatedAt && lraOld != "" && d.Get("recreate_if_rotated").(bool) {
		tflog.Debug(ctx, "found change with last_rotated_at and recreate if rotated is enabled")
		d.Set("name", "key rotated outside terraform")
		return nil
	}

	return diags
}

func resourceWizServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizServiceAccountUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	return resourceWizServiceAccountRead(ctx, d, m)
}

// DeleteServiceAccount struct
type DeleteServiceAccount struct {
	DeleteServiceAccount wiz.DeleteServiceAccountPayload `json:"deleteServiceAccount"`
}

func resourceWizServiceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizServiceAccountDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteServiceAccount (
	    $input: DeleteServiceAccountInput!
	) {
	    deleteServiceAccount(
	        input: $input
	    ) {
	        _stub
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteServiceAccountInput{}
	vars.ID = d.Id()

	// process the request
	data := &wiz.DeleteServiceAccountPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "service_account", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
