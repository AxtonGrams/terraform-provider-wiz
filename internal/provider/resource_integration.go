package provider

import (
	"context"
	//"encoding/json"
	//"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	//"wiz.io/hashicorp/terraform-provider-wiz/internal"
	//"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	//"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

// CreateIntegration struct
type CreateIntegration struct {
	CreateIntegration vendor.CreateIntegrationPayload `json:"createIntegration"`
}

func resourceWizIntegrationAwsSNSDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizIntegrationAwsSNSDelete called...")

	return diags
}
