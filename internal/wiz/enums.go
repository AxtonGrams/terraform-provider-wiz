package wiz

// YesNoUnknown enum
var YesNoUnknown = []string{
	"YES",
	"NO",
	"UNKNOWN",
}

// ProjectDataType enum
var ProjectDataType = []string{
	"CLASSIFIED",
	"HEALTH",
	"PII",
	"PCI",
	"FINANCIAL",
	"CUSTOMER",
}

// RegulatoryStandard enum
var RegulatoryStandard = []string{
	"ISO_20000_1_2011",
	"ISO_22301",
	"ISO_27001",
	"ISO_27017",
	"ISO_27018",
	"ISO_27701",
	"ISO_9001",
	"SOC",
	"FEDRAMP",
	"NIST_800_171",
	"NIST_CSF",
	"HIPPA_HITECH",
	"HITRUST",
	"PCI_DSS",
	"SEC_17A_4",
	"SEC_REGULATION_SCI",
	"SOX",
	"GDPR",
}

// BusinessImpact enum
var BusinessImpact = []string{
	"LBI",
	"MBI",
	"HBI",
}

// AuthenticationSource enum
var AuthenticationSource = []string{
	"LEGACY",
	"MODERN",
}

// Environment enum
var Environment = []string{
	"PRODUCTION",
	"STAGING",
	"DEVELOPMENT",
	"TESTING",
	"OTHER",
}

// CloudAccountStatus enum
var CloudAccountStatus = []string{
	"CONNECTED",
	"ERROR",
	"DISABLED",
	"INITIAL_SCANNING",
	"PARTIALLY_CONNECTED",
	"DISCONNECTED",
	"DISCOVERED",
}

// UserIdentityProviderType enum
var UserIdentityProviderType = []string{
	"WIZ",
	"SAML",
}

// AutomationActionType enum
var AutomationActionType = []string{
	"AWS_SNS",
	"AZURE_DEVOPS",
	"AZURE_LOGIC_APPS",
	"AZURE_SENTINEL",
	"AZURE_SERVICE_BUS",
	"CISCO_WEBEX",
	"CORTEX_XSOAR",
	"CYWARE",
	"EMAIL",
	"EVENT_BRIDGE",
	"FRESHSERVICE",
	"GOOGLE_CHAT_MESSAGE",
	"GOOGLE_PUB_SUB",
	"JIRA_TICKET",
	"JIRA_TICKET_TRANSITION",
	"MICROSOFT_TEAMS",
	"OPSGENIE_CLOSE_ALERT",
	"OPSGENIE_CREATE_ALERT",
	"PAGER_DUTY_CREATE_INCIDENT",
	"PAGER_DUTY_RESOLVE_INCIDENT",
	"SECURITY_HUB",
	"SERVICENOW_TICKET",
	"SERVICENOW_UPDATE_TICKET",
	"SLACK_MESSAGE",
	"SNOWFLAKE",
	"SPLUNK",
	"SUMO_LOGIC",
	"TORQ",
	"WEBHOOK",
}

// AwsMessageAutomationActionAccessMethodType enum
var AwsMessageAutomationActionAccessMethodType = []string{
	"ASSUME_CONNECTOR_ROLE",
	"ASSUME_SPECIFIED_ROLE",
}

// AzureServiceBusAutomationActionAccessMethodType enum
var AzureServiceBusAutomationActionAccessMethodType = []string{
	"CONNECTOR_CREDENTIALS",
	"CONNECTION_STRING_WITH_SAS",
}

// AutomationActionStatus enum
var AutomationActionStatus = []string{
	"SUCCESS",
	"FAILURE",
}

// ConnectorStatus enum
var ConnectorStatus = []string{
	"INITIAL_SCANNING",
	"PARTIALLY_CONNECTED",
	"ERROR",
	"CONNECTED",
	"DISABLED",
}

// ConnectorErrorCode enum
var ConnectorErrorCode = []string{
	"CONNECTION_ERROR",
	"DISK_SCAN_ERROR",
}

// GooglePubSubAutomationActionAccessMethodType enum
var GooglePubSubAutomationActionAccessMethodType = []string{
	"CONNECTOR_CREDENTIALS",
	"SERVICE_ACCOUNT_KEY",
}

// AutomationRuleTriggerSource enum
var AutomationRuleTriggerSource = []string{
	"ISSUES",
	"CLOUD_EVENTS",
	"CONTROL",
	"CONFIGURATION_FINDING",
}

// AutomationRuleTriggerType enum
var AutomationRuleTriggerType = []string{
	"CREATED",
	"UPDATED",
	"RESOLVED",
	"REOPENED",
}

// ServiceAccountType enum
var ServiceAccountType = []string{
	"FIRST_PARTY",
	"THIRD_PARTY",
	"SENSOR",
	"KUBERNETES_ADMISSION_CONTROLLER",
	"BROKER",
}

// DiskScanVulnerabilitySeverity enum
var DiskScanVulnerabilitySeverity = []string{
	"INFORMATIONAL",
	"LOW",
	"MEDIUM",
	"HIGH",
	"CRITICAL",
}

// IACScanSeverity enum
var IACScanSeverity = []string{
	"INFORMATIONAL",
	"LOW",
	"MEDIUM",
	"HIGH",
	"CRITICAL",
}

// Severity enum
var Severity = []string{
	"INFORMATIONAL",
	"LOW",
	"MEDIUM",
	"HIGH",
	"CRITICAL",
}

// CloudConfigurationRuleServiceType enum
var CloudConfigurationRuleServiceType = []string{
	"AWS",
	"Azure",
	"GCP",
	"OCI",
	"Alibaba",
	"AKS",
	"EKS",
	"GKE",
	"Kubernetes",
	"OKE",
}

// CloudConfigurationRuleMatcherType enum
var CloudConfigurationRuleMatcherType = []string{
	"TERRAFORM",
	"CLOUD_FORMATION",
	"KUBERNETES",
	"AZURE_RESOURCE_MANAGER",
	"DOCKER_FILE",
	"ADMISSION_CONTROLLER",
}

// CloudProvider enum
var CloudProvider = []string{
	"GCP",
	"AWS",
	"Azure",
	"OCI",
	"Alibaba",
	"vSphere",
	"OpenShift",
	"Kubernetes",
}

// KubernetesClusterKind enum
var KubernetesClusterKind = []string{
	"EKS",
	"GKE",
	"AKS",
	"OKE",
	"OPEN_SHIFT",
	"SELF_HOSTED",
}

// ControlType enum
var ControlType = []string{
	"SECURITY_GRAPH",
	"CLOUD_CONFIGURATION",
}

// ControlCreatorType enum
var ControlCreatorType = []string{
	"USER",
	"BUILTIN",
}

// TechnologyStackLayer enum
var TechnologyStackLayer = []string{
	"APPLICATION_AND_DATA",
	"CI_CD",
	"SECURITY_AND_IDENTITY",
	"COMPUTE_PLATFORMS",
	"CODE",
	"CLOUD_ENTITLEMENTS",
}

// IssueStatus enum
var IssueStatus = []string{
	"OPEN",
	"IN_PROGRESS",
	"RESOLVED",
	"REJECTED",
}

// IssueResolutionReason enum
var IssueResolutionReason = []string{
	"OBJECT_DELETED",
	"ISSUE_FIXED",
	"CONTROL_CHANGED",
	"CONTROL_DISABLED",
	"FALSE_POSITIVE",
	"EXCEPTION",
	"WONT_FIX",
}

// CloudResourceStatus enum
var CloudResourceStatus = []string{
	"Active",
	"Inactive",
	"Error",
}

// CloudPlatform enum
var CloudPlatform = []string{
	"GCP",
	"AWS",
	"Azure",
	"OCI",
	"Alibaba",
	"vSphere",
	"AKS",
	"EKS",
	"GKE",
	"Kubernetes",
	"OpenShift",
	"OKE",
}

// ConfigurationBenchmarkTypeID enum
var ConfigurationBenchmarkTypeID = []string{
	"AWS_CIS_1_2_0",
	"AWS_CIS_1_3_0",
	"AZURE_CIS_1_1_0",
	"AZURE_CIS_1_3_0",
	"GCP_CIS_1_1_0",
}

// CloudConfigurationRuleMatcherTypeFilter enum
var CloudConfigurationRuleMatcherTypeFilter = []string{
	"CLOUD",
	"TERRAFORM",
	"CLOUD_FORMATION",
	"KUBERNETES",
	"AZURE_RESOURCE_MANAGER",
	"DOCKER_FILE",
	"ADMISSION_CONTROLLER",
}

// DeploymentModel enum
var DeploymentModel = []string{
	"CLOUD_SERVICE",
	"CLOUD_PLATFORM_SERVICE",
	"SERVER_APPLICATION",
	"CLIENT_APPLICATION",
	"CODE_LIBRARY",
	"CODE",
	"VIRTUAL_APPLIANCE",
}

// TechnologyRisk enum
var TechnologyRisk = []string{
	"NONE",
	"LOW",
	"MEDIUM",
	"HIGH",
}

// TechnologyStatus enum
var TechnologyStatus = []string{
	"UNREVIEWED",
	"SANCTIONED",
	"UNSANCTIONED",
	"REQUIRED",
}

// TechnologyUsage enum
var TechnologyUsage = []string{
	"RARE",
	"UNCOMMON",
	"COMMON",
	"VERY_COMMON",
}

// OrderDirection enum
var OrderDirection = []string{
	"ASC",
	"DESC",
}

// CloudConfigurationRuleOrderField enum
var CloudConfigurationRuleOrderField = []string{
	"FAILED_CHECK_COUNT",
	"SEVERITY",
	"NAME",
}

// HostConfigurationRuleOrderField enum
var HostConfigurationRuleOrderField = []string{
	"FAILED_CHECK_COUNT",
	"NAME",
}

// GraphRelationshipType enum
var GraphRelationshipType = []string{
	"ANY",
	"ANY_OUTGOING",
	"ACTING_AS",
	"ADMINISTRATE",
	"ALERTED_ON",
	"ALLOWS",
	"ALLOWS_ACCESS_TO",
	"APPLIES_TO",
	"ASSIGNED_TO",
	"ATTACHED_TO",
	"BEHIND",
	"BOOTS",
	"BUILT_FROM",
	"CAUSES",
	"COLLABORATES",
	"CONNECTED_TO",
	"CONTAINS",
	"CONTAINS_DST_IP_RANGE",
	"CONTAINS_DST_PORT_RANGE",
	"CONTAINS_SRC_IP_RANGE",
	"CONTAINS_SRC_PORT_RANGE",
	"DENIES",
	"DEPENDS_ON",
	"DEPLOYED_TO",
	"ENCRYPTS",
	"ENCRYPTS_PARTITION",
	"ENTITLES",
	"EXCLUDES",
	"EXPOSES",
	"GOVERNS",
	"HAS",
	"HAS_BOUNDARY_POLICY",
	"HAS_DATA_FINDING",
	"HAS_DATA_INVENTORY",
	"HAS_DATA_SCHEMA",
	"HAS_DATA_STORE",
	"HAS_ORGANIZATION_POLICY",
	"HAS_PRINCIPAL_POLICY",
	"HAS_RESOURCE_POLICY",
	"HAS_SNAPSHOT",
	"HAS_SOURCE",
	"HAS_STANDARD_WEB_ACCESS_FROM",
	"HAS_TECH",
	"HOSTS",
	"IGNORES",
	"IMPLEMENTS",
	"INCLUDES",
	"INFECTS",
	"INSIDE",
	"INSTANCE_OF",
	"INVOKES",
	"LOGS_DATA_FOR",
	"MANAGES",
	"MOUNTS",
	"OWNS",
	"PART_OF",
	"PEERED_TO",
	"PERFORMED",
	"PERFORMED_IMPERSONATED",
	"PERMITS",
	"POINTS_TO",
	"PROTECTS",
	"READS_DATA_FROM",
	"REFERENCED_BY",
	"REPLICA_OF",
	"ROUTES_TRAFFIC_FROM",
	"ROUTES_TRAFFIC_TO",
	"RUNS",
	"SCANNED",
	"SEND_MESSAGES_TO",
	"SERVES",
	"STORES_DATA_IN",
	"TRANSIT_PEERED_TO",
	"USES",
	"VALIDATES",
}

// GraphEntityType enum
var GraphEntityType = []string{
	"ANY",
	"ACCESS_KEY",
	"ACCESS_ROLE",
	"ACCESS_ROLE_BINDING",
	"ACCESS_ROLE_PERMISSION",
	"API_GATEWAY",
	"APPLICATION",
	"AUTHENTICATION_CONFIGURATION",
	"AUTHENTICATION_POLICY",
	"BACKEND_BUCKET",
	"BACKUP_SERVICE",
	"BRANCH_PACKAGE",
	"BUCKET",
	"CALL_CENTER_SERVICE",
	"CDN",
	"CERTIFICATE",
	"CICD_SERVICE",
	"CLOUD_LOG_CONFIGURATION",
	"CLOUD_ORGANIZATION",
	"CLOUD_RESOURCE",
	"COMPUTE_INSTANCE_GROUP",
	"CONFIGURATION_FINDING",
	"CONFIGURATION_RULE",
	"CONFIGURATION_SCAN",
	"CONFIG_MAP",
	"CONTAINER",
	"CONTAINER_GROUP",
	"CONTAINER_IMAGE",
	"CONTAINER_INSTANCE_GROUP",
	"CONTAINER_REGISTRY",
	"CONTAINER_REPOSITORY",
	"CONTAINER_SERVICE",
	"CONTROLLER_REVISION",
	"DAEMON_SET",
	"DATABASE",
	"DATA_FINDING",
	"DATA_INVENTORY",
	"DATA_SCHEMA",
	"DATA_STORE",
	"DATA_WORKFLOW",
	"DATA_WORKLOAD",
	"DB_SERVER",
	"DEPLOYMENT",
	"DNS_RECORD",
	"DNS_ZONE",
	"DOMAIN",
	"EMAIL_SERVICE",
	"ENCRYPTION_KEY",
	"ENDPOINT",
	"EXCESSIVE_ACCESS_FINDING",
	"FILE_DESCRIPTOR",
	"FILE_DESCRIPTOR_FINDING",
	"FILE_SYSTEM_SERVICE",
	"FIREWALL",
	"GATEWAY",
	"GOVERNANCE_POLICY",
	"GOVERNANCE_POLICY_GROUP",
	"GROUP",
	"HOSTED_APPLICATION",
	"HOSTED_TECHNOLOGY",
	"HOST_CONFIGURATION_FINDING",
	"HOST_CONFIGURATION_RULE",
	"IAC_DECLARATION_INSTANCE",
	"IAC_RESOURCE_DECLARATION",
	"IAC_STATE_INSTANCE",
	"IAM_BINDING",
	"IDENTITY_PROVIDER",
	"IP_RANGE",
	"KUBERNETES_CLUSTER",
	"KUBERNETES_CRON_JOB",
	"KUBERNETES_INGRESS",
	"KUBERNETES_INGRESS_CONTROLLER",
	"KUBERNETES_JOB",
	"KUBERNETES_NETWORK_POLICY",
	"KUBERNETES_NODE",
	"KUBERNETES_PERSISTENT_VOLUME",
	"KUBERNETES_PERSISTENT_VOLUME_CLAIM",
	"KUBERNETES_POD_SECURITY_POLICY",
	"KUBERNETES_SERVICE",
	"KUBERNETES_STORAGE_CLASS",
	"KUBERNETES_VOLUME",
	"LAST_LOGIN",
	"LATERAL_MOVEMENT_FINDING",
	"LOAD_BALANCER",
	"LOCAL_USER",
	"MALWARE",
	"MALWARE_INSTANCE",
	"MANAGED_CERTIFICATE",
	"MANAGEMENT_SERVICE",
	"MAP_REDUCE_CLUSTER",
	"MESSAGING_SERVICE",
	"NAMESPACE",
	"NAT",
	"NETWORK_ADDRESS",
	"NETWORK_APPLIANCE",
	"NETWORK_INTERFACE",
	"NETWORK_ROUTING_RULE",
	"NETWORK_SECURITY_RULE",
	"PACKAGE",
	"PEERING",
	"POD",
	"PORT_RANGE",
	"PREDEFINED_GROUP",
	"PRIVATE_ENDPOINT",
	"PRIVATE_LINK",
	"PROJECT",
	"PROXY",
	"PROXY_RULE",
	"RAW_ACCESS_POLICY",
	"REGION",
	"REGISTERED_DOMAIN",
	"REPLICA_SET",
	"REPOSITORY",
	"REPOSITORY_BRANCH",
	"REPOSITORY_TAG",
	"RESOURCE_GROUP",
	"ROUTE_TABLE",
	"SEARCH_INDEX",
	"SECRET",
	"SECRET_CONTAINER",
	"SECRET_DATA",
	"SECRET_INSTANCE",
	"SECURITY_EVENT_FINDING",
	"SECURITY_TOOL_FINDING",
	"SECURITY_TOOL_FINDING_TYPE",
	"SECURITY_TOOL_SCAN",
	"SERVERLESS",
	"SERVERLESS_PACKAGE",
	"SERVICE_ACCOUNT",
	"SERVICE_CONFIGURATION",
	"SERVICE_USAGE_TECHNOLOGY",
	"SNAPSHOT",
	"STATEFUL_SET",
	"STORAGE_ACCOUNT",
	"SUBNET",
	"SUBSCRIPTION",
	"SWITCH",
	"TECHNOLOGY",
	"USER_ACCOUNT",
	"VIRTUAL_DESKTOP",
	"VIRTUAL_MACHINE",
	"VIRTUAL_MACHINE_IMAGE",
	"VIRTUAL_NETWORK",
	"VOLUME",
	"VULNERABILITY",
	"WEAKNESS",
	"WEB_SERVICE",
}

// ActionTemplateType enum
var ActionTemplateType = []string{
	"AWS_EVENT_BRIDGE",
	"AWS_SECURITY_HUB",
	"AWS_SNS",
	"AZURE_DEVOPS",
	"AZURE_LOGIC_APPS",
	"AZURE_SENTINEL",
	"AZURE_SERVICE_BUS",
	"CISCO_WEBEX",
	"CLICK_UP_CREATE_TASK",
	"CORTEX_XSOAR",
	"CYWARE",
	"EMAIL",
	"FRESHSERVICE",
	"GCP_PUB_SUB",
	"GOOGLE_CHAT",
	"HUNTERS",
	"JIRA_ADD_COMMENT",
	"JIRA_CREATE_TICKET",
	"JIRA_TRANSITION_TICKET",
	"MICROSOFT_TEAMS",
	"OPSGENIE_CLOSE_ALERT",
	"OPSGENIE_CREATE_ALERT",
	"PAGER_DUTY_CREATE_INCIDENT",
	"PAGER_DUTY_RESOLVE_INCIDENT",
	"SERVICE_NOW_CREATE_TICKET",
	"SERVICE_NOW_UPDATE_TICKET",
	"SLACK",
	"SLACK_BOT",
	"SPLUNK",
	"SUMO_LOGIC",
	"TINES",
	"TORQ",
	"WEBHOOK",
}

// IntegrationType enum
var IntegrationType = []string{
	"AWS_SECURITY_HUB",
	"AWS_SNS",
	"AZURE_DEVOPS",
	"AZURE_LOGIC_APPS",
	"AZURE_SENTINEL",
	"AZURE_SERVICE_BUS",
	"CISCO_WEBEX",
	"CORTEX_XSOAR",
	"CYWARE",
	"EMAIL",
	"AWS_EVENT_BRIDGE",
	"GOOGLE_CHAT",
	"GCP_PUB_SUB",
	"JIRA",
	"MICROSOFT_TEAMS",
	"PAGER_DUTY",
	"SERVICE_NOW",
	"SLACK",
	"SLACK_BOT",
	"SPLUNK",
	"SUMO_LOGIC",
	"TORQ",
	"WEBHOOK",
	"FRESHSERVICE",
	"OPSGENIE",
	"TINES",
	"HUNTERS",
	"CLICK_UP",
}

// JiraServerType enum
var JiraServerType = []string{
	"CLOUD",
	"SELF_HOSTED",
}

// AwsSNSIntegrationAccessMethodType enum
var AwsSNSIntegrationAccessMethodType = []string{
	"ASSUME_CONNECTOR_ROLE",
	"ASSUME_SPECIFIED_ROLE",
}

// AzureServiceBusIntegrationAccessMethodType enum
var AzureServiceBusIntegrationAccessMethodType = []string{
	"CONNECTOR_CREDENTIALS",
	"CONNECTION_STRING_WITH_SAS",
}

// GcpPubSubIntegrationAccessMethodType enum
var GcpPubSubIntegrationAccessMethodType = []string{
	"CONNECTOR_CREDENTIALS",
	"SERVICE_ACCOUNT_KEY",
}
