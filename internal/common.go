package internal

// QueryVariables struct
type QueryVariables struct {
	Query    interface{} `json:"query,omitempty"`
	ID       string      `json:"id,omitempty"`
	FilterBy interface{} `json:"filterBy,omitempty"`
	After    string      `json:"after,omitempty"`
	First    int         `json:"first,omitempty"`
}

// FilterBy struct
type FilterBy struct {
	Search []string `json:"search,omitempty"`
}

// EnumType struct
type EnumType struct {
	Type string `json:"type,omitempty"`
}

// ServiceAccountScopes defines the allowed service account rights
var ServiceAccountScopes = []string{
	"admin:all",
	"admin:audit",
	"admin:digital_trust_settings",
	"admin:identity_providers",
	"admin:projects",
	"admin:reports",
	"admin:security_settings",
	"admin:users",
	"create:action_templates",
	"create:admission_controllers",
	"create:all",
	"create:automation_actions",
	"create:automation_rules",
	"create:cloud_configuration",
	"create:cloud_event_rules",
	"create:connectors",
	"create:controls",
	"create:host_configuration",
	"create:integrations",
	"create:outposts",
	"create:reports",
	"create:run_action",
	"create:run_control",
	"create:saved_cloud_event_filters",
	"create:saved_graph_queries",
	"create:scan_policies",
	"create:security_frameworks",
	"create:security_scans",
	"create:service_accounts",
	"create:service_tickets",
	"delete:action_templates",
	"delete:all",
	"delete:automation_actions",
	"delete:automation_rules",
	"delete:cloud_configuration",
	"delete:cloud_event_rules",
	"delete:connectors",
	"delete:controls",
	"delete:host_configuration",
	"delete:integrations",
	"delete:outposts",
	"delete:reports",
	"delete:saved_cloud_event_filters",
	"delete:saved_graph_queries",
	"delete:scan_policies",
	"delete:security_frameworks",
	"delete:security_scans",
	"delete:service_accounts",
	"read:action_templates",
	"read:admission_controllers",
	"read:all",
	"read:automation_actions",
	"read:automation_rules",
	"read:benchmarks",
	"read:cloud_accounts",
	"read:cloud_configuration",
	"read:cloud_event_rules",
	"read:cloud_events",
	"read:connectors",
	"read:controls",
	"read:digital_trust_settings",
	"read:host_configuration",
	"read:integrations",
	"read:inventory",
	"read:issue_settings",
	"read:issues",
	"read:kubernetes_clusters",
	"read:licenses",
	"read:outposts",
	"read:projects",
	"read:reports",
	"read:resources",
	"read:saved_cloud_event_filters",
	"read:saved_graph_queries",
	"read:scan_policies",
	"read:scanner_settings",
	"read:security_frameworks",
	"read:security_scans",
	"read:security_settings",
	"read:service_accounts",
	"read:system_activities",
	"read:users",
	"read:vulnerabilities",
	"update:admission_controllers",
	"update:all",
	"update:automation_actions",
	"update:automation_rules",
	"update:cloud_configuration",
	"update:cloud_event_rules",
	"update:connectors",
	"update:controls",
	"update:host_configuration",
	"update:integrations",
	"update:inventory",
	"update:issue_settings",
	"update:issues",
	"update:outposts",
	"update:reports",
	"update:resources",
	"update:saved_cloud_event_filters",
	"update:saved_graph_queries",
	"update:scan_policies",
	"update:scanner_settings",
	"update:security_frameworks",
	"update:security_scans",
	"update:service_accounts",
	"update:vulnerabilities",
	"write:all",
	"write:automation_actions",
	"write:automation_rules",
	"write:cloud_configuration",
	"write:cloud_event_rules",
	"write:connectors",
	"write:controls",
	"write:host_configuration",
	"write:issue_settings",
	"write:issues",
	"write:outposts",
	"write:reports",
	"write:saved_cloud_event_filters",
	"write:saved_graph_queries",
	"write:scan_policies",
	"write:scanner_settings",
	"write:security_frameworks",
	"write:security_scans",
	"write:service_accounts",
}

// IntegrationScope is used in a Wiz Integration to define the scope of the Integration
var IntegrationScope = []string{
	"Selected Project",
	"All Resources",
	"All Resources, Restrict this Integration to global roles only",
}

// ProviderServiceNowAuthorizationType is used to infer the type of Authorization struct used in wiz.ServiceNowIntegrationParams
type ProviderServiceNowAuthorizationType struct {
	Type string `json:"type"`
}
