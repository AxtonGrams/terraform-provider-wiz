package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"wiz.io/hashicorp/terraform-provider-wiz/internal/config"
)

// New creates a new provider
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"wiz_url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Wiz api endpoint.  This varies for each Wiz deployment.  See https://docs.wiz.io/wiz-docs/docs/using-the-wiz-api#the-graphql-endpoint. (default: none, environment variable: WIZ_URL)",
					DefaultFunc: schema.EnvDefaultFunc(
						"WIZ_URL",
						nil,
					),
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.IsURLWithHTTPorHTTPS,
					),
				},
				"wiz_auth_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The authentication endpoint. (default: https://auth.app.wiz.io/oauth/token, environment variable: WIZ_AUTH_URL)",
					DefaultFunc: schema.EnvDefaultFunc(
						"WIZ_AUTH_URL",
						"https://auth.app.wiz.io/oauth/token",
					),
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.IsURLWithHTTPorHTTPS,
					),
				},
				"wiz_auth_grant_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Set this to 'client_credentials'. (default: client_credentials, environment variable: WIZ_AUTH_GRANT_TYPE)",
					DefaultFunc: schema.EnvDefaultFunc(
						"WIZ_AUTH_GRANT_TYPE",
						"client_credentials",
					),
				},
				"wiz_auth_client_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Your application's Client ID. You can find this value on the Settings > Service Accounts page. (default: none, environment variable: WIZ_AUTH_CLIENT_ID)",
					DefaultFunc: schema.EnvDefaultFunc(
						"WIZ_AUTH_CLIENT_ID",
						nil,
					),
				},
				"wiz_auth_client_secret": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Your application's Client Secret. You can find this value on the Settings > Service Accounts page. (default: none, environment variable: WIZ_AUTH_CLIENT_SECRET)",
					DefaultFunc: schema.EnvDefaultFunc(
						"WIZ_AUTH_CLIENT_SECRET",
						nil,
					),
					Sensitive: true,
				},
				"wiz_auth_audience": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Set this to 'beyond-api' if using auth0 and 'wiz-api' if using Cognito. (default: wiz-api, environment variable: WIZ_AUTH_AUDIENCE)",
					DefaultFunc: schema.EnvDefaultFunc(
						"WIZ_AUTH_AUDIENCE",
						"wiz-api",
					),
				},
				"proxy": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Use an http proxy server? (default: false, environment variable: PROXY)",
					DefaultFunc: schema.EnvDefaultFunc(
						"PROXY",
						false,
					),
				},
				"proxy_server": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Proxy server address.  Syntax: http[s]://[host]:[port]. (default: none, environment variable: PROXY_SERVER)",
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.IsURLWithHTTPorHTTPS,
					),
					DefaultFunc: schema.EnvDefaultFunc(
						"PROXY_SERVER",
						nil,
					),
				},
				"ca_chain": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Base64 encoded PEM of the CA chain used when communicating with Wiz. If a proxy performs TLS interception/inspection, this will be the CA chain for the certificate used by the proxy. The default includes the CAs known to be used by Wiz: `C=IE, O=Baltimore, OU=CyberTrust, CN=Baltimore CyberTrust Root`, `C=US, O=Cloudflare, Inc., CN=Cloudflare Inc ECC CA-3`, `C=US, ST=Arizona, L=Scottsdale, O=Starfield Technologies, Inc., CN=Starfield Services Root Certificate Authority - G2`, `C=US, O=Amazon, CN=Amazon Root CA 1`, `C=US, O=Amazon, OU=Server CA 1B, CN=Amazon`. (environment variable: CA_CHAIN) ",
					DefaultFunc: schema.EnvDefaultFunc(
						"CA_CHAIN",
						`-----BEGIN CERTIFICATE-----
MIIDdzCCAl+gAwIBAgIEAgAAuTANBgkqhkiG9w0BAQUFADBaMQswCQYDVQQGEwJJ
RTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJlclRydXN0MSIwIAYD
VQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTAwMDUxMjE4NDYwMFoX
DTI1MDUxMjIzNTkwMFowWjELMAkGA1UEBhMCSUUxEjAQBgNVBAoTCUJhbHRpbW9y
ZTETMBEGA1UECxMKQ3liZXJUcnVzdDEiMCAGA1UEAxMZQmFsdGltb3JlIEN5YmVy
VHJ1c3QgUm9vdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKMEuyKr
mD1X6CZymrV51Cni4eiVgLGw41uOKymaZN+hXe2wCQVt2yguzmKiYv60iNoS6zjr
IZ3AQSsBUnuId9Mcj8e6uYi1agnnc+gRQKfRzMpijS3ljwumUNKoUMMo6vWrJYeK
mpYcqWe4PwzV9/lSEy/CG9VwcPCPwBLKBsua4dnKM3p31vjsufFoREJIE9LAwqSu
XmD+tqYF/LTdB1kC1FkYmGP1pWPgkAx9XbIGevOF6uvUA65ehD5f/xXtabz5OTZy
dc93Uk3zyZAsuT3lySNTPx8kmCFcB5kpvcY67Oduhjprl3RjM71oGDHweI12v/ye
jl0qhqdNkNwnGjkCAwEAAaNFMEMwHQYDVR0OBBYEFOWdWTCCR1jMrPoIVDaGezq1
BE3wMBIGA1UdEwEB/wQIMAYBAf8CAQMwDgYDVR0PAQH/BAQDAgEGMA0GCSqGSIb3
DQEBBQUAA4IBAQCFDF2O5G9RaEIFoN27TyclhAO992T9Ldcw46QQF+vaKSm2eT92
9hkTI7gQCvlYpNRhcL0EYWoSihfVCr3FvDB81ukMJY2GQE/szKN+OMY3EU/t3Wgx
jkzSswF07r51XgdIGn9w/xZchMB5hbgF/X++ZRGjD8ACtPhSNzkE1akxehi/oCr0
Epn3o0WC4zxe9Z2etciefC7IpJ5OCBRLbf1wbWsaY71k5h+3zvDyny67G7fyUIhz
ksLi4xaNmjICq44Y3ekQEe5+NauQrz4wlHrQMz2nZQ/1/I6eYs9HRCwBXbsdtTLS
R9I4LtD+gdwyah617jzV/OeBHRnDJELqYzmp
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDzTCCArWgAwIBAgIQCjeHZF5ftIwiTv0b7RQMPDANBgkqhkiG9w0BAQsFADBa
MQswCQYDVQQGEwJJRTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJl
clRydXN0MSIwIAYDVQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTIw
MDEyNzEyNDgwOFoXDTI0MTIzMTIzNTk1OVowSjELMAkGA1UEBhMCVVMxGTAXBgNV
BAoTEENsb3VkZmxhcmUsIEluYy4xIDAeBgNVBAMTF0Nsb3VkZmxhcmUgSW5jIEVD
QyBDQS0zMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEua1NZpkUC0bsH4HRKlAe
nQMVLzQSfS2WuIg4m4Vfj7+7Te9hRsTJc9QkT+DuHM5ss1FxL2ruTAUJd9NyYqSb
16OCAWgwggFkMB0GA1UdDgQWBBSlzjfq67B1DpRniLRF+tkkEIeWHzAfBgNVHSME
GDAWgBTlnVkwgkdYzKz6CFQ2hns6tQRN8DAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0l
BBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBIGA1UdEwEB/wQIMAYBAf8CAQAwNAYI
KwYBBQUHAQEEKDAmMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2VydC5j
b20wOgYDVR0fBDMwMTAvoC2gK4YpaHR0cDovL2NybDMuZGlnaWNlcnQuY29tL09t
bmlyb290MjAyNS5jcmwwbQYDVR0gBGYwZDA3BglghkgBhv1sAQEwKjAoBggrBgEF
BQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzALBglghkgBhv1sAQIw
CAYGZ4EMAQIBMAgGBmeBDAECAjAIBgZngQwBAgMwDQYJKoZIhvcNAQELBQADggEB
AAUkHd0bsCrrmNaF4zlNXmtXnYJX/OvoMaJXkGUFvhZEOFp3ArnPEELG4ZKk40Un
+ABHLGioVplTVI+tnkDB0A+21w0LOEhsUCxJkAZbZB2LzEgwLt4I4ptJIsCSDBFe
lpKU1fwg3FZs5ZKTv3ocwDfjhUkV+ivhdDkYD7fa86JXWGBPzI6UAPxGezQxPk1H
goE6y/SJXQ7vTQ1unBuCJN0yJV0ReFEQPaA1IwQvZW+cwdFD19Ae8zFnWSfda9J1
CZMRJCQUzym+5iPDuI9yP+kHyCREU3qzuWFloUwOxkgAyXVjBYdwRVKD05WdRerw
6DEdfgkfCv4+3ao8XnTSrLE=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEdTCCA12gAwIBAgIJAKcOSkw0grd/MA0GCSqGSIb3DQEBCwUAMGgxCzAJBgNV
BAYTAlVTMSUwIwYDVQQKExxTdGFyZmllbGQgVGVjaG5vbG9naWVzLCBJbmMuMTIw
MAYDVQQLEylTdGFyZmllbGQgQ2xhc3MgMiBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0
eTAeFw0wOTA5MDIwMDAwMDBaFw0zNDA2MjgxNzM5MTZaMIGYMQswCQYDVQQGEwJV
UzEQMA4GA1UECBMHQXJpem9uYTETMBEGA1UEBxMKU2NvdHRzZGFsZTElMCMGA1UE
ChMcU3RhcmZpZWxkIFRlY2hub2xvZ2llcywgSW5jLjE7MDkGA1UEAxMyU3RhcmZp
ZWxkIFNlcnZpY2VzIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5IC0gRzIwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDVDDrEKvlO4vW+GZdfjohTsR8/
y8+fIBNtKTrID30892t2OGPZNmCom15cAICyL1l/9of5JUOG52kbUpqQ4XHj2C0N
Tm/2yEnZtvMaVq4rtnQU68/7JuMauh2WLmo7WJSJR1b/JaCTcFOD2oR0FMNnngRo
Ot+OQFodSk7PQ5E751bWAHDLUu57fa4657wx+UX2wmDPE1kCK4DMNEffud6QZW0C
zyyRpqbn3oUYSXxmTqM6bam17jQuug0DuDPfR+uxa40l2ZvOgdFFRjKWcIfeAg5J
Q4W2bHO7ZOphQazJ1FTfhy/HIrImzJ9ZVGif/L4qL8RVHHVAYBeFAlU5i38FAgMB
AAGjgfAwge0wDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0O
BBYEFJxfAN+qAdcwKziIorhtSpzyEZGDMB8GA1UdIwQYMBaAFL9ft9HO3R+G9FtV
rNzXEMIOqYjnME8GCCsGAQUFBwEBBEMwQTAcBggrBgEFBQcwAYYQaHR0cDovL28u
c3MyLnVzLzAhBggrBgEFBQcwAoYVaHR0cDovL3guc3MyLnVzL3guY2VyMCYGA1Ud
HwQfMB0wG6AZoBeGFWh0dHA6Ly9zLnNzMi51cy9yLmNybDARBgNVHSAECjAIMAYG
BFUdIAAwDQYJKoZIhvcNAQELBQADggEBACMd44pXyn3pF3lM8R5V/cxTbj5HD9/G
VfKyBDbtgB9TxF00KGu+x1X8Z+rLP3+QsjPNG1gQggL4+C/1E2DUBc7xgQjB3ad1
l08YuW3e95ORCLp+QCztweq7dp4zBncdDQh/U90bZKuCJ/Fp1U1ervShw3WnWEQt
8jxwmKy6abaVd38PMV4s/KCHOkdp8Hlf9BRUpJVeEXgSYCfOn8J3/yNTd126/+pZ
59vPr5KW7ySaNRB6nJHGDn2Z9j8Z3/VyVOEVqQdZe4O/Ui5GjLIAZHYcSNPYeehu
VsyuLAOQ1xk4meTKCRlb/weWsKh/NEnfVqn3sF/tM+2MR7cwA130A4w=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEkjCCA3qgAwIBAgITBn+USionzfP6wq4rAfkI7rnExjANBgkqhkiG9w0BAQsF
ADCBmDELMAkGA1UEBhMCVVMxEDAOBgNVBAgTB0FyaXpvbmExEzARBgNVBAcTClNj
b3R0c2RhbGUxJTAjBgNVBAoTHFN0YXJmaWVsZCBUZWNobm9sb2dpZXMsIEluYy4x
OzA5BgNVBAMTMlN0YXJmaWVsZCBTZXJ2aWNlcyBSb290IENlcnRpZmljYXRlIEF1
dGhvcml0eSAtIEcyMB4XDTE1MDUyNTEyMDAwMFoXDTM3MTIzMTAxMDAwMFowOTEL
MAkGA1UEBhMCVVMxDzANBgNVBAoTBkFtYXpvbjEZMBcGA1UEAxMQQW1hem9uIFJv
b3QgQ0EgMTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALJ4gHHKeNXj
ca9HgFB0fW7Y14h29Jlo91ghYPl0hAEvrAIthtOgQ3pOsqTQNroBvo3bSMgHFzZM
9O6II8c+6zf1tRn4SWiw3te5djgdYZ6k/oI2peVKVuRF4fn9tBb6dNqcmzU5L/qw
IFAGbHrQgLKm+a/sRxmPUDgH3KKHOVj4utWp+UhnMJbulHheb4mjUcAwhmahRWa6
VOujw5H5SNz/0egwLX0tdHA114gk957EWW67c4cX8jJGKLhD+rcdqsq08p8kDi1L
93FcXmn/6pUCyziKrlA4b9v7LWIbxcceVOF34GfID5yHI9Y/QCB/IIDEgEw+OyQm
jgSubJrIqg0CAwEAAaOCATEwggEtMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/
BAQDAgGGMB0GA1UdDgQWBBSEGMyFNOy8DJSULghZnMeyEE4KCDAfBgNVHSMEGDAW
gBScXwDfqgHXMCs4iKK4bUqc8hGRgzB4BggrBgEFBQcBAQRsMGowLgYIKwYBBQUH
MAGGImh0dHA6Ly9vY3NwLnJvb3RnMi5hbWF6b250cnVzdC5jb20wOAYIKwYBBQUH
MAKGLGh0dHA6Ly9jcnQucm9vdGcyLmFtYXpvbnRydXN0LmNvbS9yb290ZzIuY2Vy
MD0GA1UdHwQ2MDQwMqAwoC6GLGh0dHA6Ly9jcmwucm9vdGcyLmFtYXpvbnRydXN0
LmNvbS9yb290ZzIuY3JsMBEGA1UdIAQKMAgwBgYEVR0gADANBgkqhkiG9w0BAQsF
AAOCAQEAYjdCXLwQtT6LLOkMm2xF4gcAevnFWAu5CIw+7bMlPLVvUOTNNWqnkzSW
MiGpSESrnO09tKpzbeR/FoCJbM8oAxiDR3mjEH4wW6w7sGDgd9QIpuEdfF7Au/ma
eyKdpwAJfqxGF4PcnCZXmTA5YpaP7dreqsXMGz7KQ2hsVxa81Q4gLv7/wmpdLqBK
bRRYh5TmOTFffHPLkIhqhBGWJ6bt2YFGpn6jcgAKUj6DiAdjd4lpFw85hdKrCEVN
0FE6/V1dN2RMfjCyVSRCnTawXZwXgWHxyvkQAiSr6w10kY17RSlQOYiypok1JR4U
akcjMS9cmvqtmg5iUaQqqcT5NJ0hGA==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIESTCCAzGgAwIBAgITBn+UV4WH6Kx33rJTMlu8mYtWDTANBgkqhkiG9w0BAQsF
ADA5MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRkwFwYDVQQDExBBbWF6
b24gUm9vdCBDQSAxMB4XDTE1MTAyMjAwMDAwMFoXDTI1MTAxOTAwMDAwMFowRjEL
MAkGA1UEBhMCVVMxDzANBgNVBAoTBkFtYXpvbjEVMBMGA1UECxMMU2VydmVyIENB
IDFCMQ8wDQYDVQQDEwZBbWF6b24wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDCThZn3c68asg3Wuw6MLAd5tES6BIoSMzoKcG5blPVo+sDORrMd4f2AbnZ
cMzPa43j4wNxhplty6aUKk4T1qe9BOwKFjwK6zmxxLVYo7bHViXsPlJ6qOMpFge5
blDP+18x+B26A0piiQOuPkfyDyeR4xQghfj66Yo19V+emU3nazfvpFA+ROz6WoVm
B5x+F2pV8xeKNR7u6azDdU5YVX1TawprmxRC1+WsAYmz6qP+z8ArDITC2FMVy2fw
0IjKOtEXc/VfmtTFch5+AfGYMGMqqvJ6LcXiAhqG5TI+Dr0RtM88k+8XUBCeQ8IG
KuANaL7TiItKZYxK1MMuTJtV9IblAgMBAAGjggE7MIIBNzASBgNVHRMBAf8ECDAG
AQH/AgEAMA4GA1UdDwEB/wQEAwIBhjAdBgNVHQ4EFgQUWaRmBlKge5WSPKOUByeW
dFv5PdAwHwYDVR0jBBgwFoAUhBjMhTTsvAyUlC4IWZzHshBOCggwewYIKwYBBQUH
AQEEbzBtMC8GCCsGAQUFBzABhiNodHRwOi8vb2NzcC5yb290Y2ExLmFtYXpvbnRy
dXN0LmNvbTA6BggrBgEFBQcwAoYuaHR0cDovL2NydC5yb290Y2ExLmFtYXpvbnRy
dXN0LmNvbS9yb290Y2ExLmNlcjA/BgNVHR8EODA2MDSgMqAwhi5odHRwOi8vY3Js
LnJvb3RjYTEuYW1hem9udHJ1c3QuY29tL3Jvb3RjYTEuY3JsMBMGA1UdIAQMMAow
CAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQCFkr41u3nPo4FCHOTjY3NTOVI1
59Gt/a6ZiqyJEi+752+a1U5y6iAwYfmXss2lJwJFqMp2PphKg5625kXg8kP2CN5t
6G7bMQcT8C8xDZNtYTd7WPD8UZiRKAJPBXa30/AbwuZe0GaFEQ8ugcYQgSn+IGBI
8/LwhBNTZTUVEWuCUUBVV18YtbAiPq3yXqMB48Oz+ctBWuZSkbvkNodPLamkB2g1
upRyzQ7qDn1X8nn8N8V7YJ6y68AtkHcNSRAnpTitxBKjtKPISLMVCx7i4hncxHZS
yLyKQXhw2W2Xs0qLeC1etA+jTGDK4UfLeC0SF7FSi8o5LL21L8IzApar2pR/
-----END CERTIFICATE-----`),
				},
				"http_client_retry_max": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     10,
					Description: "Maximum retry attempts.",
				},
				"http_client_retry_wait_min": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     1,
					Description: "Minimum time to wait before retrying, in seconds.",
				},
				"http_client_retry_wait_max": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     10,
					Description: "Maximum time to wait before retrying, in seconds.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"wiz_cloud_accounts":               dataSourceWizCloudAccounts(),
				"wiz_cloud_config_rules":           dataSourceWizCloudConfigurationRules(),
				"wiz_host_config_rules":            dataSourceWizHostConfigurationRules(),
				"wiz_kubernetes_clusters":          dataSourceWizKubernetesClusters(),
				"wiz_organizations":                dataSourceWizOrganizations(),
				"wiz_subscription_resource_groups": dataSourceWizSubscriptionResourceGroups(),
				"wiz_users":                        dataSourceWizUsers(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"wiz_automation_rule_aws_sns":                  resourceWizAutomationRuleAwsSns(),
				"wiz_automation_rule_servicenow_create_ticket": resourceWizAutomationRuleServiceNowCreateTicket(),
				"wiz_automation_rule_servicenow_update_ticket": resourceWizAutomationRuleServiceNowUpdateTicket(),
				"wiz_automation_rule_jira_transition_ticket":   resourceWizAutomationRuleJiraTransitionTicket(),
				"wiz_automation_rule_jira_add_comment":         resourceWizAutomationRuleJiraAddComment(),
				"wiz_automation_rule_jira_create_ticket":       resourceWizAutomationRuleJiraCreateTicket(),
				"wiz_cicd_scan_policy":                         resourceWizCICDScanPolicy(),
				"wiz_cloud_config_rule":                        resourceWizCloudConfigurationRule(),
				"wiz_cloud_config_rule_associations":           resourceWizCloudConfigRuleAssociations(),
				"wiz_control":                                  resourceWizControl(),
				"wiz_control_associations":                     resourceWizControlAssociations(),
				"wiz_connector_aws":                            resourceWizConnectorAws(),
				"wiz_host_config_rule_associations":            resourceWizHostConfigRuleAssociations(),
				"wiz_integration_aws_sns":                      resourceWizIntegrationAwsSNS(),
				"wiz_integration_servicenow":                   resourceWizIntegrationServiceNow(),
				"wiz_integration_jira":                         resourceWizIntegrationJira(),
				"wiz_project":                                  resourceWizProject(),
				"wiz_saml_idp":                                 resourceWizSAMLIdP(),
				"wiz_security_framework":                       resourceWizSecurityFramework(),
				"wiz_service_account":                          resourceWizServiceAccount(),
				"wiz_user":                                     resourceWizUser(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description

		// We have to do some magic here since some default values include golang text/template markers, which confuses tfplugindocs
		switch v := s.Default.(type) {
		case string:
			_ = v
			if s.Default != nil {
				desc += fmt.Sprintf("\n    - Defaults to `{{`%s`}}`.", s.Default.(string))
			}
		default:
			if s.Default != nil {
				desc += fmt.Sprintf("\n    - Defaults to `%v`.", s.Default)
			}
		}
		if s.ConflictsWith != nil {
			desc += fmt.Sprintf("\n    - Conflicts with `%v`.", s.ConflictsWith)
		}
		if s.AtLeastOneOf != nil {
			desc += fmt.Sprintf("\n    - Requires least one of: `%v`.", s.AtLeastOneOf)
		}
		if s.ExactlyOneOf != nil {
			desc += fmt.Sprintf("\n    - Required exactly one of: `%v`.", s.ExactlyOneOf)
		}

		return strings.TrimSpace(desc)
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		tflog.Info(ctx, "configure called...")

		// Setup a User-Agent for the API client
		userAgent := p.UserAgent("terraform-provider-wiz", version)
		tflog.Debug(ctx, fmt.Sprintf("Provider User Agent: %s", userAgent))

		tflog.Info(ctx, "configure called...")
		var diags diag.Diagnostics
		cfg, err := config.NewConfig(d)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to configure provider",
				Detail:   fmt.Sprintf("Error: %s", err),
			})
			return nil, diags
		}
		pcfg, diags := config.NewProviderConf(ctx, cfg, userAgent)
		return pcfg, diags
	}
}
