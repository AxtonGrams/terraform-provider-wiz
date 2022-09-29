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
	"wiz.io/hashicorp/terraform-provider-wiz/internal/vendor"
)

func resourceWizAutomationAction() *schema.Resource {
	return &schema.Resource{
		Description: "Automation actions define actions to perform for findings.",
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
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The automation action type",
				Required:    true,
				ForceNew:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						vendor.AutomationActionType,
						false,
					),
				),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.IsUUID,
				),
			},
			"is_accessible_to_all_projects": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"email_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is EMAIL, define these paramemters.",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"note": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to": {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cc": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"attach_evidence_csv": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"webhook_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is WEBHOOK, define these parameters.",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"body": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_certificate": {
							Type:        schema.TypeString,
							Description: "optional client cert",
							Optional:    true,
						},
						"auth_username": {
							Type:        schema.TypeString,
							Description: "For basic authorization specify username and password, do NOT specify authToken",
							Optional:    true,
						},
						"auth_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"auth_token": {
							Type:        schema.TypeString,
							Description: "For auth bearer specify token, do not specify username/password",
							Optional:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"slack_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is SLACK_MESSAGE, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Description: "slack url in the format: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX ; see https://api.slack.com/messaging/webhooks",
							Required:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IsURLWithHTTPorHTTPS,
							),
						},
						"note": {
							Type:        schema.TypeString,
							Description: "This is the message body that will be sent to the slack channel it is a template that will be resolved with the following parameters: TBD, the format is Mustache",
							Optional:    true,
						},
						"channel": {
							Type:        schema.TypeString,
							Description: "The slack webhook has default channel, you can specify a different channel here #your-public-channel, or a specific user with @ e.g: @myself",
							Optional:    true,
						},
					},
				},
			},
			"google_chat_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is GOOGLE_CHAT_MESSAGE, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Description: "google chat webhook url in the format: https://chat.googleapis.com/v1/spaces/AAAA0000000/messages?key=XXXXX&token=XXXXX; see https://developers.google.com/hangouts/chat/how-tos/webhooks#define_an_incoming_webhook",
							Required:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IsURLWithHTTPorHTTPS,
							),
						},
						"note": {
							Type:        schema.TypeString,
							Description: "This is an optional note which will be added to the message that will be sent to the google chat room it is a template that will be resolved with the following parameters: TBD, the format is Mustache",
							Optional:    true,
						},
					},
				},
			},
			"jira_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is JIRA_TICKET, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_url": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IsURLWithHTTPorHTTPS,
							),
						},
						"is_onprem": {
							Type:        schema.TypeBool,
							Description: "Is the Jira service is only accessible on-premise?",
							Required:    true,
						},
						"onprem_tunnel_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"onprem_tunnel_token": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"tls_config": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Description: "custom TLS config (custom server CA, client certificate etc..)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_insecure_tls": {
										Type:        schema.TypeBool,
										Description: "Setting this to true will ignore any TLS validation errors on the server side certificate Warning: should only be used to validate that the action works regardless of TLS validation, if for example your server is presenting self signed or expired TLS certificate",
										Optional:    true,
										Default:     false,
									},
									"server_ca": {
										Type:        schema.TypeString,
										Description: "a PEM of the certificate authority that your server presents (if you use self signed, or custom CA)",
										Optional:    true,
									},
									"client_certificate_and_private_key": {
										Type:        schema.TypeString,
										Description: "a PEM of the client certificate as well as the certificate private key",
										Optional:    true,
										Sensitive:   true,
									},
								},
							},
						},
						"jira_authentication_basic": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"jira_authentication_token": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},
						"ticket_fields": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Description: "Ticket fields",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"summary": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Required: true,
									},
									"issue_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"assignee": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"components": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"fix_version": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"labels": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"priority": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"project": {
										Type:     schema.TypeString,
										Required: true,
									},
									"alternative_description_field": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_fields": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value should be wrapped in jsonencode() to prevent false diffs.",
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.StringIsJSON,
										),
									},
									"attach_evidence_csv": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
			"jira_transition_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is JIRA_TICKET_TRANSITION, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_url": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IsURLWithHTTPorHTTPS,
							),
						},
						"is_onprem": {
							Type:        schema.TypeBool,
							Description: "Is the Jira service is only accessible on-premise?",
							Required:    true,
						},
						"onprem_tunnel_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"onprem_tunnel_token": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						// this is marked required even though it is optional in the graphql mutation
						// the mutation, when tls_config is ommitted, defines tls_config.allow_insecure_tls=false
						"tls_config": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Description: "custom TLS config (custom server CA, client certificate etc..)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_insecure_tls": {
										Type:        schema.TypeBool,
										Description: "Setting this to true will ignore any TLS validation errors on the server side certificate Warning: should only be used to validate that the action works regardless of TLS validation, if for example your server is presenting self signed or expired TLS certificate",
										Optional:    true,
										Default:     false,
									},
									"server_ca": {
										Type:        schema.TypeString,
										Description: "a PEM of the certificate authority that your server presents (if you use self signed, or custom CA)",
										Optional:    true,
									},
									"client_certificate_and_private_key": {
										Type:        schema.TypeString,
										Description: "a PEM of the client certificate as well as the certificate private key",
										Optional:    true,
										Sensitive:   true,
									},
								},
							},
						},
						"jira_authentication_basic": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"jira_authentication_token": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
								},
							},
						},
						"project": {
							Type:     schema.TypeString,
							Required: true,
						},
						"transition_id": {
							Type:        schema.TypeString,
							Description: "Transition Id or Name",
							Required:    true,
						},
						"fields": {
							Type:        schema.TypeString,
							Description: "JSON representation of field updates. Value should be wrapped in jsonencode() to prevent false diffs.",
							Optional:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringIsJSON,
							),
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"comment_on_transition": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Comment on transition?",
						},
					},
				},
			},
			"servicenow_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is SERVICENOW_TICKET, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_url": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IsURLWithHTTPorHTTPS,
							),
						},
						"user": {
							Type:        schema.TypeString,
							Description: "Email of a Jira user with permissions to create tickets",
							Required:    true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"ticket_fields": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table_name": {
										Type:        schema.TypeString,
										Description: "Table name to which new tickets will be added to, e.g: 'Incident'",
										Required:    true,
									},
									"custom_fields": {
										Type:        schema.TypeString,
										Description: "Custom configuration fields as specified in Service Now. Make sure you add the fields that are configured as required in Service Now Project, otherwise ticket creation will fail. Value should be wrapped in jsonencode() to prevent false diffs.",
										Optional:    true,
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.StringIsJSON,
										),
									},
									"summary": {
										Type:        schema.TypeString,
										Description: "Ticket summary",
										Required:    true,
									},
									"description": {
										Type:        schema.TypeString,
										Description: "Ticket description",
										Required:    true,
									},
									"attach_evidence_csv": {
										Type:        schema.TypeBool,
										Description: "Attache evidence.",
										Optional:    true,
										Default:     false,
									},
								},
							},
						},
						"client_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"client_secret": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"servicenow_update_ticket_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is SERVICENOW_UPDATE_TICKET, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_url": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IsURLWithHTTPorHTTPS,
							),
						},
						"user": {
							Type:        schema.TypeString,
							Description: "Email of a Jira user with permissions to create tickets",
							Required:    true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fields": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value should be wrapped in jsonencode() to prevent false diffs.",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringIsJSON,
							),
						},
						"client_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"client_secret": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"aws_message_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is AWS_SNS, define these parameters",
				ConflictsWith: []string{
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sns_topic_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_method": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connector_for_access": {
										Type:        schema.TypeString,
										Description: "Required if and only if access method is ASSUME_CONNECTOR_ROLE, this should be a valid existing AWS connector ID",
										Optional:    true,
									},
									"customer_role_arn": {
										Type:        schema.TypeString,
										Description: "Required if and only if access method is ASSUME_SPECIFIED_ROLE, this is the role that should be assumed, the ExternalID of the role must be your Wiz Tenant ID (a GUID)",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"azure_service_bus_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is AZURE_SERVICE_BUS, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"email_params",
					"google_chat_params",
					"google_pub_sub_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_method": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connector_for_access": {
										Type:        schema.TypeString,
										Description: "Required if and only if access method is CONNECTOR_CREDENTIALS, this should be a valid existing Azure connector ID",
										Optional:    true,
									},
									"connection_string_with_sas": {
										Type:        schema.TypeString,
										Description: "Required if and only if access method is CONNECTION_STRING_WITH_SAS, this should be the connection string that contains the Shared access secret SAS For example: Endpoint=sb://my-sb-namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},

			"google_pub_sub_params": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Description: "If type is GOOGLE_PUB_SUB, define these parameters",
				ConflictsWith: []string{
					"aws_message_params",
					"azure_service_bus_params",
					"email_params",
					"google_chat_params",
					"jira_params",
					"jira_transition_params",
					"servicenow_params",
					"servicenow_update_ticket_params",
					"slack_params",
					"webhook_params",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"topic_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_method": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connector_for_access": {
										Type:        schema.TypeString,
										Description: "Required if and only if access method is CONNECTOR_CREDENTIALS, this should be a valid existing GCP connector ID",
										Optional:    true,
									},
									"service_account_key": {
										Type:        schema.TypeString,
										Sensitive:   true,
										Description: "Required if and only if access method is SERVICE_ACCOUNT_KEY, this should be the Service account key JSON file you downloaded from GCP. Value should be wrapped in jsonencode() to prevent false diffs.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
		CreateContext: resourceWizAutomationActionCreate,
		ReadContext:   resourceWizAutomationActionRead,
		UpdateContext: resourceWizAutomationActionUpdate,
		DeleteContext: resourceWizAutomationActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAutomationAction struct
type CreateAutomationAction struct {
	CreateAutomationAction vendor.CreateAutomationActionPayload `json:"createAutomationAction"`
}

func getServicenowParams(ctx context.Context, set interface{}) *vendor.CreateServiceNowAutomationActionParamInput {
	tflog.Info(ctx, "getServicenowParams called...")

	// return var
	var output vendor.CreateServiceNowAutomationActionParamInput

	// fetch and walk the structure
	params := set.(*schema.Set).List()
	for _, a := range params {
		tflog.Trace(ctx, fmt.Sprintf("param: %T %s", a, utils.PrettyPrint(a)))
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "base_url":
				output.BaseURL = c.(string)
			case "user":
				output.User = c.(string)
			case "password":
				output.Password = c.(string)
			case "client_id":
				output.ClientID = c.(string)
			case "client_secret":
				output.ClientSecret = c.(string)
			case "ticket_fields":
				for _, f := range c.(*schema.Set).List() {
					tflog.Trace(ctx, fmt.Sprintf("f: %T %s", f, f))
					for g, h := range f.(map[string]interface{}) {
						tflog.Trace(ctx, fmt.Sprintf("g: %T %s", g, g))
						tflog.Trace(ctx, fmt.Sprintf("h: %T %s", h, h))
						switch g {
						case "table_name":
							output.TicketFields.TableName = h.(string)
						case "custom_fields":
							output.TicketFields.CustomFields = json.RawMessage(h.(string))
						case "summary":
							output.TicketFields.Summary = h.(string)
						case "description":
							output.TicketFields.Description = h.(string)
						case "attach_evidence_csv":
							output.TicketFields.AttachEvidenceCSV = utils.ConvertBoolToPointer(h.(bool))
						}
					}
				}
			}
		}
	}

	return &output
}

func getServicenowUpdateTicketParams(ctx context.Context, set interface{}) *vendor.CreateServiceNowUpdateTicketAutomationActionParamInput {
	tflog.Info(ctx, "getServicenowUpdateTicketParams called...")

	// return var
	var output vendor.CreateServiceNowUpdateTicketAutomationActionParamInput

	// fetch and walk the structure
	params := set.(*schema.Set).List()
	for _, a := range params {
		tflog.Trace(ctx, fmt.Sprintf("param: %T %s", a, utils.PrettyPrint(a)))
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "base_url":
				output.BaseURL = c.(string)
			case "user":
				output.User = c.(string)
			case "password":
				output.Password = c.(string)
			case "client_id":
				output.ClientID = c.(string)
			case "client_secret":
				output.ClientSecret = c.(string)
			case "fields":
				output.Fields = json.RawMessage(c.(string))
			case "table_name":
				output.TableName = c.(string)
			}
		}
	}

	return &output
}

func getWebhookAutomationActionParams(ctx context.Context, set interface{}) *vendor.CreateWebhookAutomationActionParamsInput {
	tflog.Info(ctx, "getWebhookAutomationActionParams called...")

	// return var
	var output vendor.CreateWebhookAutomationActionParamsInput

	// fetch and walk the structure
	params := set.(*schema.Set).List()
	for _, a := range params {
		tflog.Trace(ctx, fmt.Sprintf("param: %T %s", a, utils.PrettyPrint(a)))
		for b, c := range a.(map[string]interface{}) {
			tflog.Trace(ctx, fmt.Sprintf("b: %T %s", b, b))
			tflog.Trace(ctx, fmt.Sprintf("c: %T %s", c, c))
			switch b {
			case "url":
				output.URL = c.(string)
			case "body":
				output.Body = c.(string)
			case "client_certificate":
				output.ClientCertificate = c.(string)
			case "auth_username":
				output.AuthUsername = c.(string)
			case "auth_password":
				output.AuthPassword = c.(string)
			case "auth_token":
				output.AuthToken = c.(string)
			}
		}
	}
	return &output
}

func resourceWizAutomationActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationActionCreate called...")

	// define the graphql query
	query := `mutation CreateAutomationAction (
	    $input: CreateAutomationActionInput!
	) {
	    createAutomationAction(
	        input: $input
	    ) {
	        automationAction {
	            id
	            createdAt
	        }
	    }
	}`

	// populate the graphql variables
	vars := &vendor.CreateAutomationActionInput{}
	vars.Name = d.Get("name").(string)
	vars.ProjectID = d.Get("project_id").(string)
	vars.IsAccessibleToAllProjects = d.Get("is_accessible_to_all_projects").(bool)
	vars.Type = d.Get("type").(string)

	var params interface{}
	if d.Get("servicenow_params").(*schema.Set).Len() > 0 {
		params = &vendor.CreateServiceNowAutomationActionParamInput{}
		vars.ServicenowParams = getServicenowParams(ctx, d.Get("servicenow_params"))
	}
	if d.Get("servicenow_update_ticket_params").(*schema.Set).Len() > 0 {
		params = &vendor.CreateServiceNowUpdateTicketAutomationActionParamInput{}
		vars.ServicenowUpdateTicketParams = getServicenowUpdateTicketParams(ctx, d.Get("servicenow_update_ticket_params"))
	}
	if d.Get("webhook_params").(*schema.Set).Len() > 0 {
		params = &vendor.CreateWebhookAutomationActionParamsInput{}
		vars.WebhookParams = getWebhookAutomationActionParams(ctx, d.Get("webhook_params"))
	}
	tflog.Debug(ctx, fmt.Sprintf("Type is %s: %T", vars.Type, params))

	// process the request
	data := &CreateAutomationAction{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_action", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id and computed values
	d.SetId(data.CreateAutomationAction.AutomationAction.ID)
	d.Set("created_at", data.CreateAutomationAction.AutomationAction.CreatedAt)

	return resourceWizAutomationActionRead(ctx, d, m)
}

func flattenAutomationActionParams(ctx context.Context, stateParams interface{}, paramType string, params interface{}) []interface{} {
	tflog.Info(ctx, "flattenAutomationActionParams called...")

	// initialize the return var
	var output = make([]interface{}, 0, 0)

	// initialize the member
	var myParams = make(map[string]interface{})

	// log the incoming data
	tflog.Debug(ctx, fmt.Sprintf("Type %s", paramType))
	tflog.Trace(ctx, fmt.Sprintf("Params %T %s", params, utils.PrettyPrint(params)))
	tflog.Trace(ctx, fmt.Sprintf("stateParams %T %+v", stateParams, stateParams))

	// populate the structure
	switch paramType {
	case "SERVICENOW_UPDATE_TICKET":
		tflog.Debug(ctx, "Handling SERVICENOW_UPDATE_TICKET")
		// get the password from state since the api doesn't return it
		var pass string
		var clientSecret string
		for _, y := range stateParams.(*schema.Set).List() {
			for s, t := range y.(map[string]interface{}) {
				switch s {
				case "password":
					pass = t.(string)
				case "client_secret":
					clientSecret = t.(string)
				}
			}
		}
		// convert generic params to specific type
		tflog.Debug(ctx, fmt.Sprintf("params %T %s", params, utils.PrettyPrint(params)))
		jsonString, _ := json.Marshal(params)
		myServiceNowUpdateTicketAutomationActionParams := &vendor.ServiceNowUpdateTicketAutomationActionParams{}
		json.Unmarshal(jsonString, &myServiceNowUpdateTicketAutomationActionParams)
		// set the parameters
		myParams["base_url"] = myServiceNowUpdateTicketAutomationActionParams.BaseURL
		myParams["user"] = myServiceNowUpdateTicketAutomationActionParams.User
		myParams["client_id"] = myServiceNowUpdateTicketAutomationActionParams.ClientID
		myParams["client_secret"] = clientSecret
		myParams["table_name"] = myServiceNowUpdateTicketAutomationActionParams.TableName
		myParams["password"] = pass
		if string(myServiceNowUpdateTicketAutomationActionParams.Fields) != "null" {
			j, err := json.Marshal(&myServiceNowUpdateTicketAutomationActionParams.Fields)
			_ = err
			myParams["fields"] = string(j)
		}
	case "SERVICENOW_TICKET":
		tflog.Debug(ctx, "Handling SERVICENOW_TICKET")
		// get the password from state since the api doesn't return it
		var pass string
		var clientSecret string
		for _, y := range stateParams.(*schema.Set).List() {
			for s, t := range y.(map[string]interface{}) {
				switch s {
				case "password":
					pass = t.(string)
				case "client_secret":
					clientSecret = t.(string)
				}
			}
		}
		// convert generic params to specific type
		tflog.Debug(ctx, fmt.Sprintf("params %T %s", params, utils.PrettyPrint(params)))
		jsonString, _ := json.Marshal(params)
		myServiceNowAutomationActionParams := &vendor.ServiceNowAutomationActionParams{}
		json.Unmarshal(jsonString, &myServiceNowAutomationActionParams)
		// set the parameters
		myParams["base_url"] = myServiceNowAutomationActionParams.BaseURL
		myParams["user"] = myServiceNowAutomationActionParams.User
		myParams["client_id"] = myServiceNowAutomationActionParams.ClientID
		myParams["client_secret"] = clientSecret
		myParams["password"] = pass
		// set the lists and sets
		var ticketFields = make([]interface{}, 0, 0)
		var ticketField = make(map[string]interface{}, 0)
		ticketField["attach_evidence_csv"] = myServiceNowAutomationActionParams.TicketFields.AttachEvidenceCSV
		ticketField["description"] = myServiceNowAutomationActionParams.TicketFields.Description
		ticketField["summary"] = myServiceNowAutomationActionParams.TicketFields.Summary
		ticketField["table_name"] = myServiceNowAutomationActionParams.TicketFields.TableName
		if string(myServiceNowAutomationActionParams.TicketFields.CustomFields) != "null" {
			j, err := json.Marshal(&myServiceNowAutomationActionParams.TicketFields.CustomFields)
			_ = err
			ticketField["custom_fields"] = string(j)
		}
		ticketFields = append(ticketFields, ticketField)
		myParams["ticket_fields"] = ticketFields
	case "WEBHOOK", "PAGER_DUTY_CREATE_INCIDENT", "PAGER_DUTY_RESOLVE_INCIDENT":
		tflog.Debug(ctx, "Handling WEBHOOK")
		// get the password and token from state since the api doesn't return it
		var authPassword string
		var clientSecret string
		for _, y := range stateParams.(*schema.Set).List() {
			tflog.Debug(ctx, fmt.Sprintf("y: %T %s", y, utils.PrettyPrint(y)))
			for s, t := range y.(map[string]interface{}) {
				tflog.Debug(ctx, fmt.Sprintf("s: %T %s", s, utils.PrettyPrint(s)))
				tflog.Debug(ctx, fmt.Sprintf("t: %T %s", t, utils.PrettyPrint(t)))
				switch s {
				case "auth_password":
					authPassword = t.(string)
				case "auth_token":
					clientSecret = t.(string)
				}
			}
		}
		// convert generic params to specific type
		tflog.Debug(ctx, fmt.Sprintf("params %T %s", params, utils.PrettyPrint(params)))
		jsonString, _ := json.Marshal(params)
		myWebhookAutomationActionParams := &vendor.WebhookAutomationActionParams{}
		json.Unmarshal(jsonString, &myWebhookAutomationActionParams)
		tflog.Debug(ctx, fmt.Sprintf("myWebhookAutomationActionParams: %T %s", myWebhookAutomationActionParams, utils.PrettyPrint(myWebhookAutomationActionParams)))

		tflog.Debug(ctx, fmt.Sprintf("auth type: %s", myWebhookAutomationActionParams.AuthenticationType.Type))

		// set the auth parameters
		switch myWebhookAutomationActionParams.AuthenticationType.Type {
		case "WebhookAutomationActionAuthenticationBasic":
			tflog.Debug(ctx, "Found authentication type WebhookAutomationActionAuthenticationBasic")
			jsonString, _ := json.Marshal(myWebhookAutomationActionParams.Authentication)
			myWebhookAutomationActionAuthenticationBasic := &vendor.WebhookAutomationActionAuthenticationBasic{}
			json.Unmarshal(jsonString, &myWebhookAutomationActionAuthenticationBasic)
			myParams["auth_username"] = myWebhookAutomationActionAuthenticationBasic.Username
			myParams["auth_password"] = authPassword
		case "WebhookAutomationActionAuthenticationTokenBearer":
			tflog.Debug(ctx, "Found authentication type WebhookAutomationActionAuthenticationTokenBearer")
			myParams["auth_token"] = clientSecret
		}
		// set the action parameters
		myParams["url"] = myWebhookAutomationActionParams.URL
		myParams["body"] = myWebhookAutomationActionParams.Body
		myParams["client_certificate"] = myWebhookAutomationActionParams.ClientCertificate
	}
	output = append(output, myParams)
	tflog.Debug(ctx, fmt.Sprintf("flattenAutomationActionParams output: %s", utils.PrettyPrint(output)))
	return output
}

// ReadAutomationActionPayload struct
type ReadAutomationActionPayload struct {
	AutomationAction *vendor.AutomationAction `json:"automationAction"`
}

func resourceWizAutomationActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationActionRead called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `query automationAction (
	    $id: ID!
	){
	    automationAction(
	        id: $id
	    ){
	        id
                createdAt
	        name
	        type
	        isAccessibleToAllProjects
	        project {
	            id
	        }
	        paramsType: params {
	            type: __typename
	        }
	        params {
	            ... on ServiceNowAutomationActionParams {
	                baseUrl
	                user
	                password
	                clientId
	                clientSecret
	                ticketFields {
	                    tableName
	                    customFields
	                    summary
	                    description
	                    attachEvidenceCSV
	                }
	            }
	            ... on ServiceNowUpdateTicketAutomationActionParams {
	                baseUrl
	                user
	                password
	                clientId
	                clientSecret
	                tableName
	                fields
	            }
	            ... on WebhookAutomationActionParams {
	                body
	                url
	                authenticationType: authentication {
	                    type: __typename
	                }
	                authentication {
	                    ... on WebhookAutomationActionAuthenticationBasic {
	                        username
	                        password
	                    }
	                    ... on WebhookAutomationActionAuthenticationTokenBearer {
	                        token
	                    }
	                }
	            }
	        }
	    }
	}`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	// process the request
	// this query returns http 200 with a payload that contains no errors and a null data body
	data := &ReadAutomationActionPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_action", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// ensure the resource was found
	tflog.Debug(ctx, fmt.Sprintf("Found data (%T) %s", data.AutomationAction, utils.PrettyPrint(data.AutomationAction)))
	if data.AutomationAction == nil {
		tflog.Info(ctx, "Resource not found, recreating. Assuming it was deleted outside terraform.")
		d.SetId("")
		d.MarkNewResource()
		return nil
	}

	// set the resource parameters
	err := d.Set("name", data.AutomationAction.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("type", data.AutomationAction.Type)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("is_accessible_to_all_projects", data.AutomationAction.IsAccessibleToAllProjects)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	err = d.Set("created_at", data.AutomationAction.CreatedAt)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	// since project is a pointer to an optional set of parameters, we need to check if it was initialized before we read/set
	if data.AutomationAction.Project == (&vendor.Project{}) {
		err = d.Set("project", data.AutomationAction.Project.ID)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	switch data.AutomationAction.Type {
	case "SERVICENOW_UPDATE_TICKET":
		params := flattenAutomationActionParams(
			ctx,
			d.Get("servicenow_update_ticket_params"),
			data.AutomationAction.Type,
			data.AutomationAction.Params,
		)
		err = d.Set("servicenow_update_ticket_params", params)
	case "SERVICENOW_TICKET":
		params := flattenAutomationActionParams(
			ctx,
			d.Get("servicenow_params"),
			data.AutomationAction.Type,
			data.AutomationAction.Params,
		)
		err = d.Set("servicenow_params", params)
	case "WEBHOOK", "PAGER_DUTY_CREATE_INCIDENT", "PAGER_DUTY_RESOLVE_INCIDENT":
		params := flattenAutomationActionParams(
			ctx,
			d.Get("webhook_params"),
			data.AutomationAction.Type,
			data.AutomationAction.Params,
		)
		err = d.Set("webhook_params", params)
	}
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateAutomationAction struct
type UpdateAutomationAction struct {
	UpdateAutomationAction vendor.UpdateAutomationActionPayload `json:"updateAutomationAction"`
}

func resourceWizAutomationActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationActionUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation updateAutomationAction($input: UpdateAutomationActionInput!) {
	    updateAutomationAction(
	        input: $input
	    ) {
	        automationAction {
	            id
	        }
	    }
	}`

	// populate the graphql variables
	// we opt to populate all params every time to ensure sensitive values are up to date
	vars := &vendor.UpdateAutomationActionInput{}
	vars.ID = d.Id()
	vars.Override = &vendor.UpdateAutomationActionChange{}
	vars.Override.Name = d.Get("name").(string)
	switch d.Get("type").(string) {
	case "SERVICENOW_TICKET":
		varsType := &vendor.UpdateServiceNowAutomationActionParamInput{}
		varsTicketFields := &vendor.UpdateServiceNowFieldsInput{}
		for _, b := range d.Get("servicenow_params").(*schema.Set).List() {
			for c, d := range b.(map[string]interface{}) {
				tflog.Debug(ctx, fmt.Sprintf("c: %T %s", c, c))
				tflog.Debug(ctx, fmt.Sprintf("d: %T %s", d, d))
				switch c {
				case "ticket_fields":
					for _, f := range d.(*schema.Set).List() {
						for g, h := range f.(map[string]interface{}) {
							tflog.Debug(ctx, fmt.Sprintf("g: %T %s", g, g))
							tflog.Debug(ctx, fmt.Sprintf("h: %T %s", h, h))
							switch g {
							case "description":
								varsTicketFields.Description = h.(string)
							case "summary":
								varsTicketFields.Summary = h.(string)
							case "table_name":
								varsTicketFields.TableName = h.(string)
							case "custom_fields":
								varsTicketFields.CustomFields = json.RawMessage(h.(string))
							case "attach_evidence_csv":
								varsTicketFields.AttachEvidenceCSV = utils.ConvertBoolToPointer(h.(bool))
							}
						}
					}
				case "base_url":
					varsType.BaseURL = d.(string)
				case "user":
					varsType.User = d.(string)
				case "password":
					varsType.Password = d.(string)
				case "client_secret":
					varsType.ClientSecret = d.(string)
				case "client_id":
					varsType.ClientID = d.(string)
				}
			}
		}
		varsType.TicketFields = varsTicketFields
		vars.Override.ServicenowParams = varsType
	case "SERVICENOW_UPDATE_TICKET":
		varsType := &vendor.UpdateServiceNowUpdateTicketAutomationActionParamInput{}
		for _, b := range d.Get("servicenow_update_ticket_params").(*schema.Set).List() {
			for c, d := range b.(map[string]interface{}) {
				tflog.Debug(ctx, fmt.Sprintf("c: %T %s", c, c))
				tflog.Debug(ctx, fmt.Sprintf("d: %T %s", d, d))
				switch c {
				case "base_url":
					varsType.BaseURL = d.(string)
				case "password":
					varsType.Password = d.(string)
				case "table_name":
					varsType.TableName = d.(string)
				case "user":
					varsType.User = d.(string)
				case "client_secret":
					varsType.ClientSecret = d.(string)
				case "client_id":
					varsType.ClientID = d.(string)
				case "fields":
					varsType.Fields = json.RawMessage(d.(string))
				}
			}
		}
		vars.Override.ServicenowUpdateTicketParams = varsType
	case "WEBHOOK", "PAGER_DUTY_CREATE_INCIDENT", "PAGER_DUTY_RESOLVE_INCIDENT":
		varsType := &vendor.UpdateWebhookAutomationActionParamsInput{}
		for _, b := range d.Get("webhook_params").(*schema.Set).List() {
			for c, d := range b.(map[string]interface{}) {
				tflog.Debug(ctx, fmt.Sprintf("c: %T %s", c, c))
				tflog.Debug(ctx, fmt.Sprintf("d: %T %s", d, d))
				switch c {
				case "url":
					varsType.URL = d.(string)
				case "body":
					varsType.Body = d.(string)
				case "client_certificate":
					varsType.ClientCertificate = d.(string)
				case "auth_username":
					varsType.AuthUsername = d.(string)
				case "auth_password":
					varsType.AuthPassword = d.(string)
				case "auth_token":
					varsType.AuthToken = d.(string)
				}
			}
		}
		vars.Override.WebhookParams = varsType
	}

	// process the request
	data := &UpdateAutomationAction{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_action", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizAutomationActionRead(ctx, d, m)
}

// DeleteAutomationAction struct
type DeleteAutomationAction struct {
	DeleteAutomationAction vendor.DeleteAutomationActionPayload `json:"deleteAutomationAction"`
}

func resourceWizAutomationActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizAutomationActionDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteAutomationAction (
	    $input: DeleteAutomationActionInput!
	) {
	    deleteAutomationAction (
	        input: $input
	    ) {
	        _stub
	    }
	}`

	// populate the graphql variables
	vars := &vendor.DeleteAutomationActionInput{}
	vars.ID = d.Id()

	// process the request
	data := &DeleteAutomationAction{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "automation_action", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
