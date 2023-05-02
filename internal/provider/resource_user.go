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

func resourceWizUser() *schema.Resource {
	return &schema.Resource{
		Description: "Users let you authenticate to Wiz.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique identifier for the user",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The user name.",
				Required:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The user email address.",
				Required:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "Whether the project is archived/inactive",
				Required:    true,
			},
			"assigned_project_ids": {
				Type:        schema.TypeList,
				Description: "Assigned Project Identifiers.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.IsUUID,
					),
				},
			},
			"send_email_invite": {
				Type:        schema.TypeBool,
				Description: "Send email invite?",
				Optional:    true,
				Default:     true,
			},
		},
		CreateContext: resourceWizUserCreate,
		ReadContext:   resourceWizUserRead,
		UpdateContext: resourceWizUserUpdate,
		DeleteContext: resourceWizUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateUser struct
type CreateUser struct {
	CreateUser wiz.CreateUserPayload `json:"createUser"`
}

func resourceWizUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizUserCreate called...")

	// define the graphql query
	query := `mutation CreateUser($input: CreateUserInput!) {
	    createUser(input: $input) {
	        user {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.CreateUserInput{}
	vars.Name = d.Get("name").(string)
	vars.Email = d.Get("email").(string)
	vars.Role = d.Get("role").(string)
	vars.SendEmailInvite = d.Get("send_email_invite").(bool)
	vars.AssignedProjectIDs = utils.ConvertListToString(d.Get("assigned_project_ids").([]interface{}))

	// process the request
	data := &CreateUser{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "user", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateUser.User.ID)

	return resourceWizUserRead(ctx, d, m)
}

func flattenAssignedProjectIDs(ctx context.Context, project []wiz.Project) []interface{} {
	tflog.Info(ctx, "flattenAssignedProjectIDs called...")
	tflog.Debug(ctx, fmt.Sprintf("flattenAssignedProjectIDs input: %+v", project))
	var output = make([]interface{}, 0, 0)
	for a, b := range project {
		tflog.Trace(ctx, fmt.Sprintf("a: %T %d", a, a))
		tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
		output = append(output, b.ID)
	}
	tflog.Debug(ctx, fmt.Sprintf("output: %s", utils.PrettyPrint(output)))
	return output
}

// ReadUserPayload struct -- updates
type ReadUserPayload struct {
	User wiz.User `json:"user"`
}

func resourceWizUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizUserRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query Users(
	    $id: ID!
	) {
	    user(
	        id: $id
	    ) {
	        id
	        name
	        email
	        effectiveAssignedProjects {
	            id
	        }
	        effectiveRole {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	data := &ReadUserPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "user", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	err := d.Set("name", data.User.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("email", data.User.Email)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("role", data.User.EffectiveRole.ID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	assignedProjectIDs := flattenAssignedProjectIDs(ctx, data.User.EffectiveAssignedProjects)
	// only set the assigned projects if they are defined; avoids empty list change detection
	if len(assignedProjectIDs) > 0 {
		err = d.Set("assigned_project_ids", assignedProjectIDs)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// UpdateUser struct
type UpdateUser struct {
	UpdateUser wiz.UpdateUserPayload `json:"updateUser"`
}

func resourceWizUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizUserUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateUser($input: UpdateUserInput!) {
	    updateUser(input: $input) {
	        user {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateUserInput{}
	vars.ID = d.Id()

	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}
	if d.HasChange("email") {
		vars.Patch.Email = d.Get("email").(string)
	}
	if d.HasChange("role") {
		vars.Patch.Role = d.Get("role").(string)
	}
	if d.HasChange("assigned_project_ids") {
		vars.Patch.AssignedProjectIDs = utils.ConvertListToString(d.Get("assigned_project_ids").([]interface{}))
	}

	// process the request
	data := &UpdateUser{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "user", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizUserRead(ctx, d, m)
}

// DeleteUser struct
type DeleteUser struct {
	DeleteUser wiz.DeleteUserPayload `json:"deleteUser"`
}

func resourceWizUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizUserDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteUser (
	    $input: DeleteUserInput!
	) {
	    deleteUser(
	        input: $input
	    ) {
	        _stub
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteUserInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateUser{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "user", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
