package provider

import (
	"context"
	"encoding/json"
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

func resourceWizCICDScanPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Configure CI/CD Scan Policies.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Internal identifier",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the Scan Policy.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the Scan Policy.",
				Optional:    true,
			},
			"builtin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The scan policy type",
				Computed:    true,
			},
			"disk_vulnerabilities_params": {
				Type:        schema.TypeSet,
				Description: "Vulnerability scan parameters.",
				Optional:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"disk_vulnerabilities_params",
					"disk_secrets_params",
					"iac_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:     schema.TypeString,
							Required: true,
							Description: fmt.Sprintf(
								"Severity.\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.DiskScanVulnerabilitySeverity,
								),
							),
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.DiskScanVulnerabilitySeverity,
									false,
								),
							),
						},
						"package_count_threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"ignore_unfixed": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"package_allow_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"disk_secrets_params": {
				Type:        schema.TypeSet,
				Description: "Secret scan parameters.",
				Optional:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"disk_vulnerabilities_params",
					"disk_secrets_params",
					"iac_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count_threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"path_allow_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"iac_params": {
				Type:        schema.TypeSet,
				Description: "IaC scan parameters.",
				Optional:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"disk_vulnerabilities_params",
					"disk_secrets_params",
					"iac_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity_threshold": {
							Type:     schema.TypeString,
							Required: true,
							Description: fmt.Sprintf(
								"Severity threshold.\n    - Allowed values: %s",
								utils.SliceOfStringToMDUList(
									wiz.IACScanSeverity,
								),
							),
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(
									wiz.IACScanSeverity,
									false,
								),
							),
						},
						"count_threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"ignored_rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"builtin_ignore_tags_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"custom_ignore_tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"rule_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"ignore_all_rules": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"security_frameworks": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		CreateContext: resourceWizCICDScanPolicyCreate,
		ReadContext:   resourceWizCICDScanPolicyRead,
		UpdateContext: resourceWizCICDScanPolicyUpdate,
		DeleteContext: resourceWizCICDScanPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateCICDScanPolicy struct
type CreateCICDScanPolicy struct {
	CreateCICDScanPolicy wiz.CreateCICDScanPolicyPayload `json:"createCICDScanPolicy"`
}

func getDiskVulnerabilitiesParams(ctx context.Context, d *schema.ResourceData) *wiz.CreateCICDScanPolicyDiskVulnerabilitiesInput {
	tflog.Info(ctx, "getDiskVulnerabilitiesParams called...")

	// return var
	var output wiz.CreateCICDScanPolicyDiskVulnerabilitiesInput

	// fetch and walk the structure
	params := d.Get("disk_vulnerabilities_params").(*schema.Set).List()
	for _, a := range params {
		tflog.Trace(ctx, fmt.Sprintf("param: %T %s", a, utils.PrettyPrint(a)))
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "severity":
				output.Severity = c.(string)
			case "package_count_threshold":
				output.PackageCountThreshold = c.(int)
			case "ignore_unfixed":
				output.IgnoreUnfixed = c.(bool)
			case "package_allow_list":
				output.PackageAllowList = utils.ConvertListToString(c.([]interface{}))
			}
		}
	}

	return &output
}

func getDiskSecretsParams(ctx context.Context, d *schema.ResourceData) *wiz.CreateCICDScanPolicyDiskSecretsInput {
	tflog.Info(ctx, "getDiskSecretsParams called...")

	// return var
	var output wiz.CreateCICDScanPolicyDiskSecretsInput

	// fetch and walk the structure
	params := d.Get("disk_secrets_params").(*schema.Set).List()
	for _, a := range params {
		tflog.Trace(ctx, fmt.Sprintf("param: %T %s", a, utils.PrettyPrint(a)))
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "count_threshold":
				output.CountThreshold = c.(int)
			case "path_allow_list":
				output.PathAllowList = utils.ConvertListToString(c.([]interface{}))
			}
		}
	}

	return &output
}

func getIACParams(ctx context.Context, d *schema.ResourceData) *wiz.CreateCICDScanPolicyIACInput {
	tflog.Info(ctx, "getIACParams called...")

	// return var
	var output wiz.CreateCICDScanPolicyIACInput
	var customTags []*wiz.CICDPolicyCustomIgnoreTagCreateInput

	// fetch and walk the structure
	params := d.Get("iac_params").(*schema.Set).List()
	for _, a := range params {
		tflog.Trace(ctx, fmt.Sprintf("param: %T %s", a, utils.PrettyPrint(a)))
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "severity_threshold":
				output.SeverityThreshold = c.(string)
			case "count_threshold":
				output.CountThreshold = c.(int)
			case "ignored_rules":
				output.IgnoredRules = utils.ConvertListToString(c.([]interface{}))
			case "builtin_ignore_tags_enabled":
				output.BuiltinIgnoreTagsEnabled = utils.ConvertBoolToPointer(c.(bool))
			case "security_frameworks":
				output.SecurityFrameworks = utils.ConvertListToString(c.([]interface{}))
			case "custom_ignore_tags":
				for _, f := range c.(*schema.Set).List() {
					tflog.Trace(ctx, fmt.Sprintf("f: %T %s", f, f))
					customTag := &wiz.CICDPolicyCustomIgnoreTagCreateInput{}
					for g, h := range f.(map[string]interface{}) {
						tflog.Trace(ctx, fmt.Sprintf("g: %T %s", g, g))
						tflog.Trace(ctx, fmt.Sprintf("h: %T %s", h, h))
						switch g {
						case "key":
							customTag.Key = h.(string)
						case "value":
							customTag.Value = h.(string)
						case "rule_ids":
							customTag.RuleIDs = utils.ConvertListToString(h.([]interface{}))
						case "ignore_all_rules":
							customTag.IgnoreAllRules = utils.ConvertBoolToPointer(h.(bool))
						}
					}
					tflog.Debug(ctx, fmt.Sprintf("customTag: %s", utils.PrettyPrint(customTag)))
					customTags = append(customTags, customTag)
				}
			}
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("customTags: %s", utils.PrettyPrint(customTags)))
	output.CustomIgnoreTags = customTags

	return &output
}

func resourceWizCICDScanPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCICDScanPolicyCreate called...")

	// define the graphql query
	query := `mutation CreateCICDScanPolicy(
	    $input: CreateCICDScanPolicyInput!
	) {
	    createCICDScanPolicy (
	        input: $input
	    ) {
	        scanPolicy
	        {
	            id
	            builtin
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.CreateCICDScanPolicyInput{}
	var policyType string
	vars.Name = d.Get("name").(string)
	vars.Description = d.Get("description").(string)
	if d.Get("disk_vulnerabilities_params").(*schema.Set).Len() > 0 {
		policyType = "CICDScanPolicyParamsVulnerabilities"
		vars.DiskVulnerabilitiesParams = getDiskVulnerabilitiesParams(ctx, d)
	}
	if d.Get("disk_secrets_params").(*schema.Set).Len() > 0 {
		policyType = "CICDScanPolicyParamsSecrets"
		vars.DiskSecretsParams = getDiskSecretsParams(ctx, d)
	}
	if d.Get("iac_params").(*schema.Set).Len() > 0 {
		policyType = "CICDScanPolicyParamsIAC"
		vars.IACParams = getIACParams(ctx, d)
	}

	// process the request
	data := &CreateCICDScanPolicy{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cicd_scan_policy", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id and computed values
	d.SetId(data.CreateCICDScanPolicy.ScanPolicy.ID)
	d.Set("builtin", data.CreateCICDScanPolicy.ScanPolicy.Builtin)
	d.Set("type", policyType)

	return resourceWizCICDScanPolicyRead(ctx, d, m)
}

func flattenScanPolicyParams(ctx context.Context, paramType string, params interface{}) []interface{} {
	tflog.Info(ctx, "flattenParams called...")

	// initialize the return var
	var output = make([]interface{}, 0, 0)

	// initialize the member
	var myParams = make(map[string]interface{})

	// log the incoming data
	tflog.Debug(ctx, fmt.Sprintf("Type %s", paramType))
	tflog.Trace(ctx, fmt.Sprintf("Params %T %s", params, utils.PrettyPrint(params)))

	// populate the structure
	switch paramType {
	case "CICDScanPolicyParamsIAC":
		tflog.Debug(ctx, "Handling CICDScanPolicyParamsIAC")

		// convert generic params to specific type
		tflog.Debug(ctx, fmt.Sprintf("params %T %s", params, utils.PrettyPrint(params)))
		jsonString, _ := json.Marshal(params)
		myCICDScanPolicyParamsIAC := &wiz.CICDScanPolicyParamsIAC{}
		json.Unmarshal(jsonString, &myCICDScanPolicyParamsIAC)
		tflog.Debug(
			ctx,
			fmt.Sprintf(
				"myCICDScanPolicyParamsIAC %T %s",
				myCICDScanPolicyParamsIAC,
				utils.PrettyPrint(
					myCICDScanPolicyParamsIAC,
				),
			),
		)

		myParams["count_threshold"] = myCICDScanPolicyParamsIAC.CountThreshold
		myParams["builtin_ignore_tags_enabled"] = myCICDScanPolicyParamsIAC.BuiltinIgnoreTagsEnabled
		myParams["severity_threshold"] = myCICDScanPolicyParamsIAC.SeverityThreshold

		var ignoredRules = make([]interface{}, 0, 0)
		for a, b := range myCICDScanPolicyParamsIAC.IgnoredRules {
			tflog.Debug(ctx, fmt.Sprintf("a: %T %d", a, a))
			tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
			ignoredRules = append(ignoredRules, b.ID)
		}
		myParams["ignored_rules"] = ignoredRules

		var securityFrameWorks = make([]interface{}, 0, 0)
		for a, b := range myCICDScanPolicyParamsIAC.SecurityFrameworks {
			tflog.Debug(ctx, fmt.Sprintf("a: %T %d", a, a))
			tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
			securityFrameWorks = append(securityFrameWorks, b.ID)
		}
		myParams["security_frameworks"] = securityFrameWorks

		var customIgnoreTags = make([]interface{}, 0, 0)
		for a, b := range myCICDScanPolicyParamsIAC.CustomIgnoreTags {
			tflog.Debug(ctx, fmt.Sprintf("a: %T %d", a, a))
			tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
			var customIgnoreTag = make(map[string]interface{}, 0)
			customIgnoreTag["ignore_all_rules"] = b.IgnoreAllRules
			customIgnoreTag["key"] = b.Key
			customIgnoreTag["value"] = b.Value
			var rules = make([]interface{}, 0, 0)
			for c, d := range b.Rules {
				tflog.Debug(ctx, fmt.Sprintf("c: %T %d", c, c))
				tflog.Debug(ctx, fmt.Sprintf("d: %T %s", d, utils.PrettyPrint(d)))
				rules = append(rules, d.ID)
			}
			customIgnoreTag["rule_ids"] = rules
			customIgnoreTags = append(customIgnoreTags, customIgnoreTag)
		}
		myParams["custom_ignore_tags"] = customIgnoreTags
	case "CICDScanPolicyParamsSecrets":
		tflog.Debug(ctx, "Handling CICDScanPolicyParamsSecrets")

		// convert generic params to specific type
		tflog.Debug(ctx, fmt.Sprintf("params %T %s", params, utils.PrettyPrint(params)))
		jsonString, _ := json.Marshal(params)
		myCICDScanPolicyParamsSecrets := &wiz.CICDScanPolicyParamsSecrets{}
		json.Unmarshal(jsonString, &myCICDScanPolicyParamsSecrets)
		tflog.Debug(
			ctx,
			fmt.Sprintf(
				"myCICDScanPolicyParamsSecrets %T %s",
				myCICDScanPolicyParamsSecrets,
				utils.PrettyPrint(
					myCICDScanPolicyParamsSecrets,
				),
			),
		)

		myParams["count_threshold"] = myCICDScanPolicyParamsSecrets.CountThreshold

		var pathAllowList = make([]interface{}, 0, 0)
		for a, b := range myCICDScanPolicyParamsSecrets.PathAllowList {
			tflog.Debug(ctx, fmt.Sprintf("a: %T %d", a, a))
			tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
			pathAllowList = append(pathAllowList, b)
		}
		myParams["path_allow_list"] = pathAllowList
	case "CICDScanPolicyParamsVulnerabilities":
		tflog.Debug(ctx, "Handling CICDScanPolicyParamsVulnerabilities")

		// convert generic params to specific type
		tflog.Debug(ctx, fmt.Sprintf("params %T %s", params, utils.PrettyPrint(params)))
		jsonString, _ := json.Marshal(params)
		myCICDScanPolicyParamsVulnerabilities := &wiz.CICDScanPolicyParamsVulnerabilities{}
		json.Unmarshal(jsonString, &myCICDScanPolicyParamsVulnerabilities)
		tflog.Debug(
			ctx,
			fmt.Sprintf(
				"myCICDScanPolicyParamsVulnerabilities %T %s",
				myCICDScanPolicyParamsVulnerabilities,
				utils.PrettyPrint(
					myCICDScanPolicyParamsVulnerabilities,
				),
			),
		)

		myParams["ignore_unfixed"] = myCICDScanPolicyParamsVulnerabilities.IgnoreUnfixed
		myParams["package_count_threshold"] = myCICDScanPolicyParamsVulnerabilities.PackageCountThreshold
		myParams["severity"] = myCICDScanPolicyParamsVulnerabilities.Severity

		var packageAllowList = make([]interface{}, 0, 0)
		for a, b := range myCICDScanPolicyParamsVulnerabilities.PackageAllowList {
			tflog.Debug(ctx, fmt.Sprintf("a: %T %d", a, a))
			tflog.Debug(ctx, fmt.Sprintf("b: %T %s", b, utils.PrettyPrint(b)))
			packageAllowList = append(packageAllowList, b)
		}
		myParams["package_allow_list"] = packageAllowList
	}

	output = append(output, myParams)
	tflog.Info(ctx, fmt.Sprintf("flattenScanPolicyParams output: %s", utils.PrettyPrint(output)))
	return output
}

// ReadCICDScanPolicyPayload struct
type ReadCICDScanPolicyPayload struct {
	CICDScanPolicy wiz.CICDScanPolicy `json:"cicdScanPolicy"`
}

func resourceWizCICDScanPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCICDScanPolicyRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query CICDScanPolicy  (
	    $id: ID!
	) {
	    cicdScanPolicy(
                id: $id
	    ) {
	        id
	        name
	        description
	        builtin
	        paramsType: params {
	            type: __typename
	        }
	        params {
	            ... on CICDScanPolicyParamsVulnerabilities {
	                severity
	                packageCountThreshold
	                ignoreUnfixed
	                packageAllowList
	            }
	            ... on CICDScanPolicyParamsSecrets {
	                countThreshold
	                pathAllowList
	            }
	            ... on CICDScanPolicyParamsIAC {
	                builtinIgnoreTagsEnabled
	                countThreshold
	                severityThreshold
	                ignoredRules {
	                    id
	                }
	                customIgnoreTags {
	                    key
	                    value
	                    ignoreAllRules
	                    rules {
	                        id
	                    }
	                }
	                securityFrameworks {
	                    id
	                }
	            }
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: oops! an internal error has occurred. for reference purposes, this is your request id
	data := &ReadCICDScanPolicyPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cicd_scan_policy", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.CICDScanPolicy.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.CICDScanPolicy.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("description", data.CICDScanPolicy.Description)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("builtin", data.CICDScanPolicy.Builtin)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	params := flattenScanPolicyParams(ctx, data.CICDScanPolicy.ParamsType.Type, data.CICDScanPolicy.Params)
	switch data.CICDScanPolicy.ParamsType.Type {
	case "CICDScanPolicyParamsIAC":
		err = d.Set("type", "CICDScanPolicyParamsIAC")
		err = d.Set("iac_params", params)
	case "CICDScanPolicyParamsSecrets":
		err = d.Set("type", "CICDScanPolicyParamsSecrets")
		err = d.Set("disk_secrets_params", params)
	case "CICDScanPolicyParamsVulnerabilities":
		err = d.Set("type", "CICDScanPolicyParamsVulnerabilities")
		err = d.Set("disk_vulnerabilities_params", params)
	}
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateCICDScanPolicy struct
type UpdateCICDScanPolicy struct {
	UpdateCICDScanPolicy wiz.UpdateCICDScanPolicyPayload `json:"updateCICDScanPolicy"`
}

func resourceWizCICDScanPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCICDScanPolicyUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation updateCICDScanPolicy(
	    $input: UpdateCICDScanPolicyInput
	) {
	    updateCICDScanPolicy(
	        input: $input
	    ) {
	        scanPolicy {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	vars := &wiz.UpdateCICDScanPolicyInput{}
	vars.ID = d.Id()
	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		vars.Patch.Name = d.Get("description").(string)
	}
	// we need to evaluate whether the policy type changed before setting the params
	if d.Get("disk_vulnerabilities_params").(*schema.Set).Len() > 0 {
		err := d.Set("type", "CICDScanPolicyParamsVulnerabilities")
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.Get("disk_secrets_params").(*schema.Set).Len() > 0 {
		err := d.Set("type", "CICDScanPolicyParamsSecrets")
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	if d.Get("iac_params").(*schema.Set).Len() > 0 {
		err := d.Set("type", "CICDScanPolicyParamsIAC")
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	switch d.Get("type") {
	case "CICDScanPolicyParamsIAC":
		tflog.Debug(ctx, "Handling updates for CICDScanPolicyParamsIAC")
		varsType := &wiz.UpdateCICDScanPolicyIACPatch{}
		varsTypeIgnoreTags := make([]*wiz.CICDPolicyCustomIgnoreTagUpdateInput, 0)
		for a, b := range d.Get("iac_params").(*schema.Set).List() {
			tflog.Trace(ctx, fmt.Sprintf("a: (%T) %d", a, a))
			tflog.Trace(ctx, fmt.Sprintf("b: (%T) %s", b, utils.PrettyPrint(b)))
			for c, d := range b.(map[string]interface{}) {
				tflog.Trace(ctx, fmt.Sprintf("c: (%T) %s", c, c))
				tflog.Trace(ctx, fmt.Sprintf("d: (%T) %s", d, utils.PrettyPrint(d)))
				switch c {
				case "count_threshold":
					varsType.CountThreshold = d.(int)
				case "severity_threshold":
					varsType.SeverityThreshold = d.(string)
				case "builtin_ignore_tags_enabled":
					varsType.BuiltinIgnoreTagsEnabled = utils.ConvertBoolToPointer(d.(bool))
				case "ignored_rules":
					varsType.IgnoredRules = utils.ConvertListToString(d.([]interface{}))
				case "security_frameworks":
					varsType.SecurityFrameworks = utils.ConvertListToString(d.([]interface{}))
				case "custom_ignore_tags":
					for e, f := range d.(*schema.Set).List() {
						tflog.Trace(ctx, fmt.Sprintf("e: (%T) %d", e, e))
						tflog.Trace(ctx, fmt.Sprintf("f: (%T) %s", f, utils.PrettyPrint(f)))
						varsTypeIgnoreTag := &wiz.CICDPolicyCustomIgnoreTagUpdateInput{}
						for g, h := range f.(map[string]interface{}) {
							tflog.Trace(ctx, fmt.Sprintf("g: (%T) %s", g, g))
							tflog.Trace(ctx, fmt.Sprintf("h: (%T) %s", h, utils.PrettyPrint(h)))
							switch g {
							case "key":
								varsTypeIgnoreTag.Key = h.(string)
							case "value":
								varsTypeIgnoreTag.Value = h.(string)
							case "ignore_all_rules":
								varsTypeIgnoreTag.IgnoreAllRules = utils.ConvertBoolToPointer(h.(bool))
							case "rule_ids":
								varsTypeIgnoreTag.RuleIDs = utils.ConvertListToString(h.([]interface{}))
							}
						}
						tflog.Debug(ctx, fmt.Sprintf("varsTypeIgnoreTag: %s", utils.PrettyPrint(varsTypeIgnoreTag)))
						varsTypeIgnoreTags = append(varsTypeIgnoreTags, varsTypeIgnoreTag)
					}
				}
			}
		}
		varsType.CustomIgnoreTags = varsTypeIgnoreTags
		vars.Patch.IACParams = varsType
	case "CICDScanPolicyParamsSecrets":
		tflog.Debug(ctx, "Handling updates for CICDScanPolicyParamsSecrets")
		varsType := &wiz.UpdateCICDScanPolicyDiskSecretsPatch{}
		for a, b := range d.Get("disk_secrets_params").(*schema.Set).List() {
			tflog.Trace(ctx, fmt.Sprintf("a: (%T) %d", a, a))
			tflog.Trace(ctx, fmt.Sprintf("b: (%T) %s", b, utils.PrettyPrint(b)))
			for c, d := range b.(map[string]interface{}) {
				tflog.Trace(ctx, fmt.Sprintf("c: (%T) %s", c, c))
				tflog.Trace(ctx, fmt.Sprintf("d: (%T) %s", d, utils.PrettyPrint(d)))
				switch c {
				case "count_threshold":
					varsType.CountThreshold = d.(int)
				case "path_allow_list":
					varsType.PathAllowList = utils.ConvertListToString(d.([]interface{}))
				}
			}
		}
		vars.Patch.DiskSecretsParams = varsType
	case "CICDScanPolicyParamsVulnerabilities":
		tflog.Debug(ctx, "Handling updates for CICDScanPolicyParamsVulnerabilities")
		varsType := &wiz.UpdateCICDScanPolicyDiskVulnerabilitiesPatch{}
		for a, b := range d.Get("disk_vulnerabilities_params").(*schema.Set).List() {
			tflog.Trace(ctx, fmt.Sprintf("a: (%T) %d", a, a))
			tflog.Trace(ctx, fmt.Sprintf("b: (%T) %s", b, utils.PrettyPrint(b)))
			for c, d := range b.(map[string]interface{}) {
				tflog.Trace(ctx, fmt.Sprintf("c: (%T) %s", c, c))
				tflog.Trace(ctx, fmt.Sprintf("d: (%T) %s", d, utils.PrettyPrint(d)))
				switch c {
				case "ignore_unfixed":
					varsType.IgnoreUnfixed = utils.ConvertBoolToPointer(d.(bool))
				case "package_allow_list":
					varsType.PackageAllowList = utils.ConvertListToString(d.([]interface{}))
				case "package_count_threshold":
					varsType.PackageCountThreshold = d.(int)
				case "severity":
					varsType.Severity = d.(string)
				}
			}
		}
		vars.Patch.DiskVulnerabilitiesParams = varsType
	}

	// process the request
	data := &UpdateCICDScanPolicy{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cicd_scan_policy", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizCICDScanPolicyRead(ctx, d, m)
}

// DeleteCICDScanPolicy struct
type DeleteCICDScanPolicy struct {
	DeleteCICDScanPolicy wiz.DeleteCICDScanPolicyPayload `json:"deleteCICDScanPolicy"`
}

func resourceWizCICDScanPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizCICDScanPolicyDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteCICDScanPolicy (
	    $input: DeleteCICDScanPolicyInput!
	) {
	    deleteCICDScanPolicy(
	        input: $input
	    ) {
	        id
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteCICDScanPolicyInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateCICDScanPolicy{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "cicd_scan_policy", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
