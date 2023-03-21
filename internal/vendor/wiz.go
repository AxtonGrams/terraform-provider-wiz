package vendor

import (
	"encoding/json"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
)

/*
  Notes:
  - DateTime types are included if immutable (created vs last update)
  - Type references for struct members are a pointer if the member is optional; this allows for proper JSON
  - json.RawMessage is used for fields that house raw json strings
  - bool elements must be bool* if omitempty is specified. see https://stackoverflow.com/questions/37756236/json-golang-boolean-omitempty
  - GraphQL union types require that the schema be extended so the type is known
  - GraphQL enums are represented by slices; validation is performed by the resource
  - Sometimes deviations are required from the Wiz schema definition; please note the deviation in the struct comments
*/

// PageInfo struct
type PageInfo struct {
	EndCursor   string `json:"endCursor,omitempty"`
	HasNextPage bool   `json:"hasNextPage"`
}

// CloudOrganizationFilters struct
type CloudOrganizationFilters struct {
	CloudProvider []string `json:"cloudProvider,omitempty"` // enum CloudProvider
	ProjectID     string   `json:"projectId,omitempty"`
	Search        []string `json:"search,omitempty"`
}

// KubernetesClusterFilters struct
type KubernetesClusterFilters struct {
	Search       string   `json:"search,omitempty"`
	Kind         []string `json:"kind,omitempty"`
	CloudAccount []string `json:"cloudAccount,omitempty"`
}

// UserFilters struct
type UserFilters struct {
	Search               string   `json:"search,omitempty"`
	Roles                []string `json:"role,omitempty"`
	AuthenticationSource string   `json:"source,omitempty"`
}

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

// ProjectRiskProfileInput struct
type ProjectRiskProfileInput struct {
	IsActivelyDeveloped string   `json:"isActivelyDeveloped,omitempty"` // enum YesNoUnknown
	HasAuthentication   string   `json:"hasAuthentication,omitempty"`   // enum YesNoUnknown
	HasExposedAPI       string   `json:"hasExposedAPI,omitempty"`       // enum YesNoUnknown
	IsInternetFacing    string   `json:"isInternetFacing,omitempty"`    // enum YesNoUnknown
	IsCustomerFacing    string   `json:"isCustomerFacing,omitempty"`    // enum YesNoUnknown
	StoresData          string   `json:"storesData,omitempty"`          // enum YesNoUnknown
	SensitiveDataTypes  []string `json:"sensitiveDataTypes,omitempty"`  // enum ProjectDataType
	BusinessImpact      string   `json:"businessImpact,omitempty"`      // enum BusinessImpact
	IsRegulated         string   `json:"isRegulated,omitempty"`         // enum YesNoUnknown
	RegulatoryStandards []string `json:"regulatoryStandards,omitempty"` // enum RegulatoryStandard
}

// CreateProjectPayload struct -- updates
type CreateProjectPayload struct {
	Project Project `json:"project"`
}

// CreateProjectInput struct -- updates
type CreateProjectInput struct {
	Name                   string                               `json:"name"`
	Slug                   string                               `json:"slug,omitempty"`
	Description            string                               `json:"description,omitempty"`
	Archived               *bool                                `json:"archived,omitempty"`
	Identifiers            []string                             `json:"identifiers,omitempty"`
	BusinessUnit           string                               `json:"businessUnit,omitempty"`
	ProjectOwners          []string                             `json:"projectOwners,omitempty"`
	SecurityChampion       []string                             `json:"securityChampions,omitempty"`
	RiskProfile            ProjectRiskProfileInput              `json:"riskProfile"`
	CloudOrganizationLinks []*ProjectCloudOrganizationLinkInput `json:"cloudOrganizationLinks,omitempty"`
	CloudAccountLinks      []*ProjectCloudAccountLinkInput      `json:"cloudAccountLinks,omitempty"`
	KubernetesClusterLinks []*ProjectKubernetesClusterLinkInput `json:"kubernetesClusterLinks,omitempty"`
}

// Environment enum
var Environment = []string{
	"PRODUCTION",
	"STAGING",
	"DEVELOPMENT",
	"TESTING",
	"OTHER",
}

// ProjectCloudOrganizationLinkInput struct
type ProjectCloudOrganizationLinkInput struct {
	CloudOrganization string              `json:"cloudOrganization"`
	Environment       string              `json:"environment"` // enum Environment
	ResourceGroups    []string            `json:"resourceGroups,omitempty"`
	ResourceTags      []*ResourceTagInput `json:"resourceTags,omitempty"`
	Shared            bool                `json:"shared"`
}

// ProjectCloudAccountLinkInput struct
type ProjectCloudAccountLinkInput struct {
	CloudAccount   string              `json:"cloudAccount"`
	Environment    string              `json:"environment"` // enum Environment
	ResourceGroups []string            `json:"resourceGroups,omitempty"`
	ResourceTags   []*ResourceTagInput `json:"resourceTags,omitempty"`
	Shared         *bool               `json:"shared,omitempty"`
}

// ProjectKubernetesClusterLinkInput struct
type ProjectKubernetesClusterLinkInput struct {
	KubernetesCluster string   `json:"kubernetesCluster"`
	Environment       string   `json:"environment"` // enum Environment
	Namespaces        []string `json:"namespaces,omitempty"`
	Shared            bool     `json:"shared"`
}

// ResourceTagInput struct -- updates
type ResourceTagInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ResourceTag struct
type ResourceTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ProjectCloudOrganizationLink struct
type ProjectCloudOrganizationLink struct {
	CloudOrganization CloudOrganization `json:"cloudOrganization"`
	Environment       string            `json:"environment"` // enum Environment
	ResourceTags      []*ResourceTag    `json:"resourceTags,omitempty"`
	ResourceGroups    []string          `json:"resourceGroups,omitempty"`
	Shared            bool              `json:"shared"`
}

// CloudOrganization struct -- updates
type CloudOrganization struct {
	ID            string `json:"id"`
	ExternalID    string `json:"externalId"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	CloudProvider string `json:"cloudProvider,omitempty"` // enum CloudProvider
}

// UpdateProjectInput struct
type UpdateProjectInput struct {
	ID       string              `json:"id"`
	Override UpdateProjectPatch  `json:"override,omitempty"`
	Patch    *UpdateProjectPatch `json:"patch,omitempty"`
}

// UpdateProjectPayload struct
type UpdateProjectPayload struct {
	Project Project `json:"project"`
}

// UpdateProjectPatch struct
// We deviate from the GraphQL schema to include resource links because the update requires an empty value to nullify removed attributes
// The slug is required in the request in order to override update and deletion contexts
type UpdateProjectPatch struct {
	Name                   string                               `json:"name,omitempty"`
	Archived               *bool                                `json:"archived,omitempty"`
	Description            string                               `json:"description,omitempty"`
	BusinessUnit           string                               `json:"businessUnit,omitempty"`
	SecurityChampions      []string                             `json:"securityChampions,omitempty"`
	Identifiers            []string                             `json:"identifiers,omitempty"`
	ProjectOwners          []string                             `json:"projectOwners,omitempty"`
	RiskProfile            *ProjectRiskProfileInput             `json:"riskProfile,omitempty"`
	Slug                   string                               `json:"slug"`
	CloudOrganizationLinks []*ProjectCloudOrganizationLinkInput `json:"cloudOrganizationLinks"`
	CloudAccountLinks      []*ProjectCloudAccountLinkInput      `json:"cloudAccountLinks"`
	KubernetesClusterLinks []*ProjectKubernetesClusterLinkInput `json:"kubernetesClusterLinks"`
}

// UpdateSAMLIdentityProviderInput struct
type UpdateSAMLIdentityProviderInput struct {
	ID    string                          `json:"id"`
	Patch UpdateSAMLIdentityProviderPatch `json:"patch"`
}

// UpdateSAMLIdentityProviderPatch struct
// We deviate from the GraphQL schema to include all attributes because the update requires an empty value to nullify removed attributes
type UpdateSAMLIdentityProviderPatch struct {
	Name                     string                        `json:"name"`
	IssuerURL                string                        `json:"issuerURL"`
	LoginURL                 string                        `json:"loginURL"`
	LogoutURL                string                        `json:"logoutURL"`
	UseProviderManagedRoles  *bool                         `json:"useProviderManagedRoles"`
	AllowManualRoleOverride  *bool                         `json:"allowManualRoleOverride"`
	Certificate              string                        `json:"certificate"`
	Domains                  []string                      `json:"domains"`
	GroupMapping             []SAMLGroupMappingUpdateInput `json:"groupMapping"`
	MergeGroupsMappingByRole *bool                         `json:"mergeGroupsMappingByRole"`
}

// UpdateSAMLIdentityProviderPayload struct -- updates
type UpdateSAMLIdentityProviderPayload struct {
	SAMLIdentityProvider SAMLIdentityProvider `json:"samlIdentityProvider"`
}

// SAMLGroupMappingUpdateInput struct
type SAMLGroupMappingUpdateInput struct {
	ProviderGroupID string   `json:"providerGroupId"`
	Role            string   `json:"role"`
	Projects        []string `json:"projects"`
}

// CreateSAMLIdentityProviderInput struct -- updates
type CreateSAMLIdentityProviderInput struct {
	Name                     string                         `json:"name"`
	IssuerURL                string                         `json:"issuerURL,omitempty"`
	LoginURL                 string                         `json:"loginURL"`
	LogoutURL                string                         `json:"logoutURL,omitempty"`
	UseProviderManagedRoles  bool                           `json:"useProviderManagedRoles"`
	AllowManualRoleOverride  *bool                          `json:"allowManualRoleOverride,omitempty"`
	Certificate              string                         `json:"certificate"`
	Domains                  []string                       `json:"domains"`
	GroupMapping             []*SAMLGroupMappingCreateInput `json:"groupMapping,omitempty"`
	MergeGroupsMappingByRole *bool                          `json:"mergeGroupsMappingByRole,omitempty"`
}

// CreateSAMLIdentityProviderPayload struct -- updates
type CreateSAMLIdentityProviderPayload struct {
	SAMLIdentityProvider SAMLIdentityProvider `json:"samlIdentityProvider,omitempty"`
}

// SAMLGroupMappingCreateInput struct -- updates
type SAMLGroupMappingCreateInput struct {
	ProviderGroupID string   `json:"providerGroupId"`
	Role            string   `json:"role"`
	Projects        []string `json:"projects"`
}

// SAMLIdentityProvider struct -- updates
type SAMLIdentityProvider struct {
	AllowManualRoleOverride  *bool               `json:"allowManualRoleOverride"`
	Certificate              string              `json:"certificate"`
	Domains                  []string            `json:"domains"`
	GroupMapping             []*SAMLGroupMapping `json:"groupMapping,omitempty"`
	ID                       string              `json:"id"`
	IssuerURL                string              `json:"issuerURL,omitempty"`
	LoginURL                 string              `json:"loginURL"`
	LogoutURL                string              `json:"logoutURL"`
	MergeGroupsMappingByRole bool                `json:"mergeGroupsMappingByRole"`
	Name                     string              `json:"name"`
	UseProviderManagedRoles  bool                `json:"useProviderManagedRoles"`
}

// SAMLGroupMapping struct -- updates
type SAMLGroupMapping struct {
	Projects        []Project `json:"projects"`
	ProviderGroupID string    `json:"providerGroupId"`
	Role            UserRole  `json:"role"`
}

// DeleteSAMLIdentityProviderInput struct
type DeleteSAMLIdentityProviderInput struct {
	ID string `json:"id"`
}

// DeleteSAMLIdentityProviderPayload struct -- updated
type DeleteSAMLIdentityProviderPayload struct {
	Stub string `json:"_stub"`
}

// DeleteAutomationActionInput struct -- updates
type DeleteAutomationActionInput struct {
	ID string `json:"id"`
}

// DeleteAutomationActionPayload struct -- updates
type DeleteAutomationActionPayload struct {
	Stub string `json:"_stub,omitempty"`
}

// UpdateAutomationActionInput struct -- updates
type UpdateAutomationActionInput struct {
	ID       string                        `json:"id"`
	Merge    *UpdateAutomationActionChange `json:"merge,omitempty"`
	Override *UpdateAutomationActionChange `json:"override,omitempty"`
}

// UpdateAutomationActionChange struct
type UpdateAutomationActionChange struct {
	Name                         string                                                  `json:"name,omitempty"`
	EmailParams                  *UpdateEmailAutomationActionParamsInput                 `json:"emailParams,omitempty"`
	WebhookParams                *UpdateWebhookAutomationActionParamsInput               `json:"webhookParams,omitempty"`
	SlackParams                  *UpdateSlackMessageAutomationActionParamsInput          `json:"slackParams,omitempty"`
	GoogleChatParams             *UpdateGoogleChatMessageAutomationActionParamsInput     `json:"googleChatParams,omitempty"`
	JiraParams                   *UpdateJiraAutomationActionParamInput                   `json:"jiraParams,omitempty"`
	JiraTransitionParams         *UpdateJiraTransitionAutomationActionParamInput         `json:"jiraTransitionParams,omitempty"`
	ServicenowParams             *UpdateServiceNowAutomationActionParamInput             `json:"servicenowParams,omitempty"`
	ServicenowUpdateTicketParams *UpdateServiceNowUpdateTicketAutomationActionParamInput `json:"servicenowUpdateTicketParams,omitempty"`
	AWSMessageParams             *UpdateAwsMessageAutomationActionParamsInput            `json:"awsMessageParams,omitempty"`
	AzureServiceBusParams        *UpdateAzureServiceBusAutomationActionParamsInput       `json:"azureServiceBusParams,omitempty"`
	GooglePubSubParams           *UpdateGooglePubSubAutomationActionParamsInput          `json:"googlePubSubParams,omitempty"`
}

// Project struct
type Project struct {
	Archived                bool                            `json:"archived"`
	BusinessUnit            string                          `json:"businessUnit"`
	CloudAccountCount       int                             `json:"cloudAccountCount"`
	CloudAccountLinks       []*ProjectCloudAccountLink      `json:"cloudAccountLinks"`
	CloudOrganizationCount  int                             `json:"cloudOrganizationCount"`
	CloudOrganizationLinks  []*ProjectCloudOrganizationLink `json:"cloudOrganizationLinks"`
	Description             string                          `json:"description"`
	EntityCount             int                             `json:"entityCount"`
	Entrypoints             []*ProjectEntrypoint            `json:"entrypoints"`
	ID                      string                          `json:"id"`
	Identifiers             []string                        `json:"identifiers"`
	KubernetesClustersLinks []*ProjectKubernetesClusterLink `json:"kubernetesClustersLinks"`
	Name                    string                          `json:"name"`
	ProfileCompletion       int                             `json:"profileCompletion"`
	ProjectOwners           []*User                         `json:"projectOwners"`
	RepositoryCount         int                             `json:"repositoryCount"`
	RiskProfile             ProjectRiskProfile              `json:"riskProfile"`
	SecurityChampions       []*User                         `json:"securityChampions"`
	Slug                    string                          `json:"slug"`
	TeamMemberCount         int                             `json:"teamMemberCount"`
	TechnologyCount         int                             `json:"technologyCount"`
}

// ProjectCloudAccountLink struct
type ProjectCloudAccountLink struct {
	CloudAccount   CloudAccount   `json:"cloudAccount"`
	Environment    string         `json:"environment"` // enum Environment
	ResourceGroups []string       `json:"resourceGroups,omitempty"`
	ResourceTags   []*ResourceTag `json:"ResourceTags,omitempty"`
	Shared         bool           `json:"shared"`
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

// CloudAccount struct
type CloudAccount struct {
	CloudProvider       string      `json:"cloudProvider"` // enum CloudProvider
	ContainerCount      int         `json:"containerCount"`
	ExternalID          string      `json:"externalId"`
	FirstScannedAt      string      `json:"firstScannedAt"`
	ID                  string      `json:"id"`
	LastScannedAt       string      `json:"lastScannedAt"`
	LinkedProjects      []*Project  `json:"linkedProjects,omitempty"`
	Name                string      `json:"name"`
	ResourceCount       int         `json:"resourceCount"`
	SourceConnectors    []Connector `json:"sourceConnectors"`
	Status              string      `json:"status"` // enum CloudAccountStatus
	VirtualMachineCount int         `json:"virtualMachineCount"`
}

// ProjectEntrypoint struct
type ProjectEntrypoint struct {
	Environment string `json:"environment"` // enum Environment
	URL         string `json:"url"`
}

// UserIdentityProviderType enum
var UserIdentityProviderType = []string{
	"WIZ",
	"SAML",
}

// User struct
type User struct {
	AssignedSAMLGroups               []SAMLGroupMapping   `json:"assignedSAMLGroups"`
	EffectiveAssignedProjects        []Project            `json:"effectiveAssignedProjects"`
	CreatedAt                        string               `json:"createdAt"`
	EffectiveRole                    UserRole             `json:"effectiveRole"`
	Email                            string               `json:"email"`
	ID                               string               `json:"id"`
	IdentityProvider                 SAMLIdentityProvider `json:"identityProvider"` // UserIdentityProviderType
	IdentityProviderAssignedProjects []Project            `json:"identityProviderAssignedProjects"`
	IdentityProviderRole             UserRole             `json:"identityProviderRole"`
	IdentityProviderType             string               `json:"identityProviderType"`
	IntercomUserHash                 string               `json:"intercomUserHash"`
	IPAddress                        string               `json:"ipAddress"`
	IsProjectScoped                  bool                 `json:"isProjectScoped"`
	IsSuspended                      bool                 `json:"isSuspended"`
	ManualOverrideAssignedProjects   []Project            `json:"manualOverrideAssignedProjects"`
	ManualOverrideRole               UserRole             `json:"manualOverrideRole"`
	Name                             string               `json:"name"`
	Preferences                      string               `json:"preferences"`
	ReadmeAuthToken                  string               `json:"readmeAuthToken"`
	ZendeskAuthToken                 string               `json:"zendeskAuthToken"`
}

// KubernetesCluster struct
type KubernetesCluster struct {
	ID                     string       `json:"id"`
	ExternalID             string       `json:"externalId"`
	Name                   string       `json:"name"`
	Kind                   string       `json:"kind"`
	Status                 string       `json:"status"`
	CloudAccount           CloudAccount `json:"cloudAccount"`
	IsPrivate              bool         `json:"isPrivate"`
	Project                []Project    `json:"projects"`
	Connectors             []Connector  `json:"connectors"`
	IsConnectedUsingBroker bool         `json:"isConnectedUsingBroker"`
}

// ProjectKubernetesClusterLink struct
type ProjectKubernetesClusterLink struct {
	KubernetesCluster KubernetesCluster `json:"kubernetesCluster"`
	Environment       string            `json:"environment"`
	Namespaces        []string          `json:"namespaces,omitempty"`
	Shared            bool              `json:"shared"`
}

// UserRole struct
type UserRole struct {
	Description     string   `json:"description"`
	ID              string   `json:"id"`
	IsProjectScoped bool     `json:"isProjectScoped"`
	Name            string   `json:"name"`
	Scopes          []string `json:"scopes"`
}

// UserPreferences struct
type UserPreferences struct {
	SelectedSAMLGroup SAMLGroupMapping `json:"selectedSAMLGroup"`
}

// ProjectRiskProfile struct -- updates
type ProjectRiskProfile struct {
	BusinessImpact      string   `json:"businessImpact,omitempty"` // enum BusinessImpact
	HasAuthentication   string   `json:"hasAuthentication"`        // enum YesNoUnknown
	HasExposedAPI       string   `json:"hasExposedAPI"`            // enum YesNoUnknown
	IsActivelyDeveloped string   `json:"isActivelyDeveloped"`      // enum YesNoUnknown
	IsCustomerFacing    string   `json:"isCustomerFacing"`         // enum YesNoUnknown
	IsInternetFacing    string   `json:"isInternetFacing"`         // enum YesNoUnknown
	IsRegulated         string   `json:"isRegulated"`              // enum YesNoUnknown
	RegulatoryStandards []string `json:"regulatoryStandards"`      // enum RegulatoryStandard
	SensitiveDataTypes  []string `json:"sensitiveDataTypes"`       // enum ProjectDataType
	StoresData          string   `json:"storesData"`               // enum YesNoUnknown
}

// CloudOrganizationConnection struct -- updates
type CloudOrganizationConnection struct {
	Edges      []*CloudOrganizationEdge `json:"edges,omitempty"`
	Nodes      []CloudOrganization      `json:"nodes,omitempty"`
	PageInfo   PageInfo                 `json:"pageInfo"`
	TotalCount int                      `json:"totalCount"`
}

// CloudOrganizationEdge struct -- updates
type CloudOrganizationEdge struct {
	Cursor string            `json:"cursor"`
	Node   CloudOrganization `json:"node"`
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

// CreateAutomationActionInput struct -- updates
type CreateAutomationActionInput struct {
	Name                         string                                                  `json:"name"`
	Type                         string                                                  `json:"type,omitempty"` // enum AutomationActionType
	ProjectID                    string                                                  `json:"projectId,omitempty"`
	IsAccessibleToAllProjects    bool                                                    `json:"isAccessibleToAllProjects"`
	EmailParams                  *CreateEmailAutomationActionParamsInput                 `json:"emailParams,omitempty"`
	WebhookParams                *CreateWebhookAutomationActionParamsInput               `json:"webhookParams,omitempty"`
	SlackParams                  *CreateSlackMessageAutomationActionParamsInput          `json:"slackParams,omitempty"`
	GoogleChatParams             *CreateGoogleChatMessageAutomationActionParamsInput     `json:"googleChatParams,omitempty"`
	JiraParams                   *CreateJiraAutomationActionParamInput                   `json:"jiraParams,omitempty"`
	JiraTransitionParams         *CreateJiraTransitionAutomationActionParamInput         `json:"jiraTransitionParams,omitempty"`
	ServicenowParams             *CreateServiceNowAutomationActionParamInput             `json:"servicenowParams,omitempty"`
	ServicenowUpdateTicketParams *CreateServiceNowUpdateTicketAutomationActionParamInput `json:"servicenowUpdateTicketParams,omitempty"`
	AWSMessageParams             *CreateAwsMessageAutomationActionParamsInput            `json:"awsMessageParams,omitempty"`
	AzureServiceBusParams        *CreateAzureServiceBusAutomationActionParamsInput       `json:"azureServiceBusParams,omitempty"`
	GooglePubSubParams           *CreateGooglePubSubAutomationActionParamsInput          `json:"googlePubSubParams,omitempty"`
}

// CreateEmailAutomationActionParamsInput struct -- updates
type CreateEmailAutomationActionParamsInput struct {
	Note              string   `json:"note,omitempty"`
	To                []string `json:"to"`
	CC                []string `json:"cc,omitempty"`
	AttachEvidenceCSV *bool    `json:"attachEvidenceCSV,omitempty"`
}

// CreateWebhookAutomationActionParamsInput struct -- updates
type CreateWebhookAutomationActionParamsInput struct {
	URL               string `json:"url"`
	Body              string `json:"body"`
	ClientCertificate string `json:"clientCertificate,omitempty"`
	AuthUsername      string `json:"authUsername,omitempty"`
	AuthPassword      string `json:"authPassword,omitempty"`
	AuthToken         string `json:"authToken,omitempty"`
}

// CreateSlackMessageAutomationActionParamsInput struct -- updates
type CreateSlackMessageAutomationActionParamsInput struct {
	URL     string `json:"url"`
	Note    string `json:"body,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// CreateGoogleChatMessageAutomationActionParamsInput struct -- upates
type CreateGoogleChatMessageAutomationActionParamsInput struct {
	URL  string `json:"url"`
	Note string `json:"body,omitempty"`
}

// CreateJiraAutomationActionParamInput struct -- updates
type CreateJiraAutomationActionParamInput struct {
	ServerURL    string                          `json:"serverUrl"`
	IsOnPrem     bool                            `json:"isOnPrem"`
	TLSConfig    *AutomationActionTLSConfigInput `json:"tlsConfig,omitempty"`
	User         string                          `json:"user"`
	Token        string                          `json:"token"`
	TicketFields CreateJiraTicketFieldsInput     `json:"ticketFields"`
}

// AutomationActionTLSConfigInput struct -- updates
type AutomationActionTLSConfigInput struct {
	AllowInsecureTLS               *bool  `json:"allowInsecureTLS,omitempty"`
	ServerCA                       string `json:"serverCA,omitempty"`
	ClientCertificateAndPrivateKey string `json:"clientCertificateAndPrivateKey,omitempty"`
}

// CreateJiraTicketFieldsInput struct -- updates
type CreateJiraTicketFieldsInput struct {
	Summary                     string          `json:"summary"`
	Description                 string          `json:"description"`
	IssueType                   string          `json:"issueType"`
	Assignee                    string          `json:"assignee,omitempty"`
	Components                  []string        `json:"components,omitempty"`
	FixVersion                  []string        `json:"fixVersion,omitempty"`
	Labels                      []string        `json:"labels,omitempty"`
	Priority                    string          `json:"priority,omitempty"`
	Project                     string          `json:"project"`
	AlternativeDescriptionField string          `json:"alternativeDescriptionField,omitempty"`
	CustomFields                json.RawMessage `json:"customFields,omitempty"`
	AttachEvidenceCSV           *bool           `json:"attachEvidenceCSV,omitempty"`
}

// CreateJiraTransitionAutomationActionParamInput struct -- updates
type CreateJiraTransitionAutomationActionParamInput struct {
	ServerURL           string                          `json:"serverUrl"`
	IsOnPrem            bool                            `json:"isOnPrem"`
	TLSConfig           *AutomationActionTLSConfigInput `json:"tlsConfig,omitempty"`
	User                string                          `json:"user"`
	Token               string                          `json:"token"`
	Project             string                          `json:"project"`
	TransitionID        string                          `json:"transitionId"`
	Fields              json.RawMessage                 `json:"fields,omitempty"`
	Comment             string                          `json:"comment,omitempty"`
	CommentOnTransition *bool                           `json:"commentOnTransition,omitempty"`
}

// CreateServiceNowAutomationActionParamInput struct -- updates
type CreateServiceNowAutomationActionParamInput struct {
	BaseURL      string                      `json:"baseUrl"`
	User         string                      `json:"user"`
	Password     string                      `json:"password"`
	TicketFields CreateServiceNowFieldsInput `json:"ticketFields,omitempty"`
	ClientID     string                      `json:"clientId,omitempty"`
	ClientSecret string                      `json:"clientSecret,omitempty"`
}

// CreateServiceNowFieldsInput struct -- updates
type CreateServiceNowFieldsInput struct {
	TableName         string          `json:"tableName"`
	CustomFields      json.RawMessage `json:"customFields,omitempty"`
	Summary           string          `json:"summary"`
	Description       string          `json:"description"`
	AttachEvidenceCSV *bool           `json:"attachEvidenceCSV,omitempty"`
}

// CreateServiceNowUpdateTicketAutomationActionParamInput struct -- updates
type CreateServiceNowUpdateTicketAutomationActionParamInput struct {
	BaseURL      string          `json:"baseUrl"`
	User         string          `json:"user"`
	Password     string          `json:"password"`
	TableName    string          `json:"tableName"`
	Fields       json.RawMessage `json:"fields,omitempty"`
	ClientID     string          `json:"clientId,omitempty"`
	ClientSecret string          `json:"clientSecret,omitempty"`
}

// CreateAwsMessageAutomationActionParamsInput struct -- updates
type CreateAwsMessageAutomationActionParamsInput struct {
	SNSTopicARN  string                                      `json:"snsTopicARN"`
	Body         string                                      `json:"body"`
	AccessMethod AwsMessageAutomationActionAccessMethodInput `json:"AwsMessageAutomationActionAccessMethodInput"`
}

// AwsMessageAutomationActionAccessMethodType enum
var AwsMessageAutomationActionAccessMethodType = []string{
	"ASSUME_CONNECTOR_ROLE",
	"ASSUME_SPECIFIED_ROLE",
}

// AwsMessageAutomationActionAccessMethodInput struct -- updates
type AwsMessageAutomationActionAccessMethodInput struct {
	Type               string `json:"type"` // enum AwsMessageAutomationActionAccessMethodType
	ConnectorForAccess string `json:"connectorForAccess,omitempty"`
	CustomerRoleARN    string `json:"customerRoleARN,omitempty"`
}

// CreateAzureServiceBusAutomationActionParamsInput struct -- updates
type CreateAzureServiceBusAutomationActionParamsInput struct {
	QueueURL     string                                           `json:"queueUrl"`
	Body         string                                           `json:"body"`
	AccessMethod AzureServiceBusAutomationActionAccessMethodInput `json:"accessMethod"`
}

// AzureServiceBusAutomationActionAccessMethodType enum
var AzureServiceBusAutomationActionAccessMethodType = []string{
	"CONNECTOR_CREDENTIALS",
	"CONNECTION_STRING_WITH_SAS",
}

// AzureServiceBusAutomationActionAccessMethodInput struct -- updates
type AzureServiceBusAutomationActionAccessMethodInput struct {
	Type                    string `json:"type"` // enum AzureServiceBusAutomationActionAccessMethodType
	ConnectorForAccess      string `json:"connectorForAccess,omitempty"`
	ConnectionStringWithSAS string `json:"connectionStringWithSAS,omitempty"`
}

// CreateAutomationActionPayload struct -- updates
type CreateAutomationActionPayload struct {
	AutomationAction *AutomationAction `json:"automationAction,omitempty"`
}

// AutomationActionStatus enum
var AutomationActionStatus = []string{
	"SUCCESS",
	"FAILURE",
}

// AutomationAction struct -- updates; this is incomplete.  missing usedByRules. added paramsType
type AutomationAction struct {
	CreatedAt                 string            `json:"createdAt"`
	ID                        string            `json:"id"`
	IsAccessibleToAllProjects bool              `json:"isAccessibleToAllProjects"`
	Name                      string            `json:"name"`
	ParamsType                internal.EnumType `json:"paramsType"`
	Params                    interface{}       `json:"params"`
	Project                   *Project          `json:"project,omitempty"`
	Status                    string            `json:"AutomationActionStatus,omitempty"` // enum AutomationActionStatus
	Type                      string            `json:"type"`                             // enum AutomationActionType
}

// UpdateWebhookAutomationActionParamsInput struct -- updates
type UpdateWebhookAutomationActionParamsInput struct {
	URL               string `json:"url,omitempty"`
	Body              string `json:"body,omitempty"`
	ClientCertificate string `json:"clientCertificate,omitempty"`
	AuthUsername      string `json:"authUsername,omitempty"`
	AuthPassword      string `json:"authPassword,omitempty"`
	AuthToken         string `json:"authToken,omitempty"`
}

// SlackMessageAutomationActionParams struct -- updates
type SlackMessageAutomationActionParams struct {
	Channel string `json:"channel,omitempty"`
	Note    string `json:"note"`
	URL     string `json:"url"`
}

// GoogleChatMessageAutomationActionParams struct -- updates
type GoogleChatMessageAutomationActionParams struct {
	Note string `json:"note,omitempty"`
	URL  string `json:"url"`
}

// JiraAutomationActionParams struct -- updates
type JiraAutomationActionParams struct {
	IsOnPrem               bool                       `json:"isOnPrem"`
	JiraAuthenticationType internal.EnumType          `json:"jiraAuthenticationType"`
	JiraAuthentication     interface{}                `json:"jiraAuthentication"`
	OnPremTunnelDomain     string                     `json:"onPremTunnelDomain,omitempty"`
	OnPremTunnelToken      string                     `json:"onPremTunnelToken,omitempty"`
	ServerURL              string                     `json:"serverUrl"`
	TicketFields           JiraTicketFields           `json:"ticketFields"`
	TLSConfig              *AutomationActionTLSConfig `json:"tlsConfig,omitempty"`
}

// JiraTicketFields struct -- updates
type JiraTicketFields struct {
	AlternativeDescriptionField string          `json:"alternativeDescriptionField,omitempty"`
	Assignee                    string          `json:"assignee,omitempty"`
	AttachEvidenceCSV           *bool           `json:"attachEvidenceCSV,omitempty"`
	Components                  []string        `json:"components,omitempty"`
	CustomFields                json.RawMessage `json:"customFields,omitempty"`
	Description                 string          `json:"description"`
	FixVersion                  []string        `json:"fixVersion,omitempty"`
	IssueType                   string          `json:"issueType"`
	Labels                      []string        `json:"labels,omitempty"`
	Priority                    string          `json:"priority,omitempty"`
	Project                     string          `json:"project"`
	Summary                     string          `json:"summary"`
}

// AutomationActionTLSConfig struct -- updates
type AutomationActionTLSConfig struct {
	AllowInsecureTLS               *bool  `json:"allowInsecureTLS,omitempty"`
	ClientCertificateAndPrivateKey string `json:"clientCertificateAndPrivateKey,omitempty"`
	ServerCA                       string `json:"serverCA,omitempty"`
}

// JiraTransitionAutomationActionParams struct -- updates
type JiraTransitionAutomationActionParams struct {
	Comment                string                     `json:"comment,omitempty"`
	CommentOnTransition    *bool                      `json:"commentOnTransition,omitempty"`
	Fields                 json.RawMessage            `json:"fields,omitempty"`
	IsOnPrem               bool                       `json:"isOnPrem"`
	JiraAuthenticationType internal.EnumType          `json:"jiraAuthenticationType"`
	JiraAuthentication     interface{}                `json:"jiraAuthentication"`
	OnPremTunnelDomain     string                     `json:"onPremTunnelDomain,omitempty"`
	OnPremTunnelToken      string                     `json:"onPremTunnelToken,omitempty"`
	Project                string                     `json:"project"`
	ServerURL              string                     `json:"serverUrl"`
	TLSConfig              *AutomationActionTLSConfig `json:"tlsConfig,omitempty"`
	TransitionID           string                     `json:"transitionId"`
}

// JiraAutomationActionAuthenticationBasic struct
type JiraAutomationActionAuthenticationBasic struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// JiraAutomationActionAuthenticationTokenBearer struct
type JiraAutomationActionAuthenticationTokenBearer struct {
	Token string `json:"token"`
}

// ServiceNowAutomationActionParams struct -- updates
type ServiceNowAutomationActionParams struct {
	BaseURL      string                 `json:"baseUrl"`
	ClientID     string                 `json:"clientId,omitempty"`
	ClientSecret string                 `json:"clientSecret,omitempty"`
	Password     string                 `json:"password"`
	TicketFields ServiceNowTicketFields `json:"ticketFields"`
	User         string                 `json:"user"`
}

// ServiceNowTicketFields struct -- updates
type ServiceNowTicketFields struct {
	AttachEvidenceCSV *bool           `json:"attachEvidenceCSV,omitempty"`
	CustomFields      json.RawMessage `json:"customFields,omitempty"`
	Description       string          `json:"description"`
	Summary           string          `json:"summary"`
	TableName         string          `json:"tableName"`
}

// ServiceNowUpdateTicketAutomationActionParams struct -- updates
type ServiceNowUpdateTicketAutomationActionParams struct {
	BaseURL      string          `json:"baseUrl"`
	ClientID     string          `json:"clientId,omitempty"`
	ClientSecret string          `json:"clientSecret,omitempty"`
	Fields       json.RawMessage `json:"fields,omitempty"`
	Password     string          `json:"password"`
	TableName    string          `json:"tableName"`
	User         string          `json:"user"`
}

// AwsMessageAutomationActionParams struct -- updates
type AwsMessageAutomationActionParams struct {
	AccessMethod       string    `json:"accessMethod"` // enum AwsMessageAutomationActionAccessMethodType
	Body               string    `json:"body"`
	ConnectorForAccess Connector `json:"connectorForAccess,omitempty"`
	CustomerRoleARN    string    `json:"customerRoleARN,omitempty"`
	SNSTopicARN        string    `json:"snsTopicARN"`
}

// AzureServiceBusAutomationActionParams struct -- updates
type AzureServiceBusAutomationActionParams struct {
	AccessMethod            string    `json:"accessMethod"` // enum AzureServiceBusAutomationActionAccessMethodType
	Body                    string    `json:"body"`
	ConnectionStringWithSAS string    `json:"connectionStringWithSAS,omitempty"`
	ConnectorForAccess      Connector `json:"connectorForAccess,omitempty"`
	QueueURL                string    `json:"queueUrl"`
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

// Connector struct -- updates.  this resource is incomplete.  missing ConnectorConfigs ConnectorModules TenantIdentityClient
type Connector struct {
	AddedBy           User            `json:"addedBy"`
	AuthParams        json.RawMessage `json:"authParams"`
	CloudAccountCount int             `json:"cloudAccountCount"`
	CreatedAt         string          `json:"createdAt"`
	Enabled           bool            `json:"enabled"`
	ExtraConfig       json.RawMessage `json:"extraConfig,omitempty"`
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Status            string          `json:"status"` // enum ConnectorStatus
}

// CreateGooglePubSubAutomationActionParamsInput struct -- updates
type CreateGooglePubSubAutomationActionParamsInput struct {
	ProjectID    string                                        `json:"projectId"`
	TopicID      string                                        `json:"topicId"`
	Body         string                                        `json:"body"`
	AccessMethod GooglePubSubAutomationActionAccessMethodInput `json:"accessMethod"`
}

// GooglePubSubAutomationActionAccessMethodType enum
var GooglePubSubAutomationActionAccessMethodType = []string{
	"CONNECTOR_CREDENTIALS",
	"SERVICE_ACCOUNT_KEY",
}

// GooglePubSubAutomationActionAccessMethodInput struct -- updates
type GooglePubSubAutomationActionAccessMethodInput struct {
	Type               string          `json:"type"` // enum GooglePubSubAutomationActionAccessMethodType
	ConnectorForAccess string          `json:"connectorForAccess,omitempty"`
	ServiceAccountKey  json.RawMessage `json:"serviceAccountKey,omitempty"`
}

// UpdateEmailAutomationActionParamsInput struct -- updates
type UpdateEmailAutomationActionParamsInput struct {
	Note              string   `json:"note,omitempty"`
	To                []string `json:"to,omitempty"`
	CC                []string `json:"cc,omitempty"`
	AttachEvidenceCSV *bool    `json:"attachEvidenceCSV,omitempty"`
}

// UpdateSlackMessageAutomationActionParamsInput struct -- updates
type UpdateSlackMessageAutomationActionParamsInput struct {
	URL     string `json:"url,omitempty"`
	Note    string `json:"note,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// UpdateGoogleChatMessageAutomationActionParamsInput struct -- updates
type UpdateGoogleChatMessageAutomationActionParamsInput struct {
	URL  string `json:"url,omitempty"`
	Note string `json:"note,omitempty"`
}

// UpdateJiraAutomationActionParamInput struct -- updates
type UpdateJiraAutomationActionParamInput struct {
	ServerURL    string                          `json:"serverUrl,omitempty"`
	IsOnPrem     *bool                           `json:"isOnPrem,omitempty"`
	TLSConfig    *AutomationActionTLSConfigInput `json:"tlsConfig,omitempty"`
	User         string                          `json:"user,omitempty"`
	Token        string                          `json:"token,omitempty"`
	TicketFields *UpdateJiraTicketFieldsInput    `json:"ticketFields,omitempty"`
}

// UpdateJiraTransitionAutomationActionParamInput struct -- updates
type UpdateJiraTransitionAutomationActionParamInput struct {
	ServerURL           string                          `json:"serverUrl,omitempty"`
	IsOnPrem            *bool                           `json:"isOnPrem,omitempty"`
	TLSConfig           *AutomationActionTLSConfigInput `json:"tlsConfig,omitempty"`
	User                string                          `json:"user,omitempty"`
	Token               string                          `json:"token,omitempty"`
	Project             string                          `json:"project,omitempty"`
	TransitionID        string                          `json:"transitionId,omitemptyd"`
	Fields              json.RawMessage                 `json:"fields,omitempty"`
	Comment             string                          `json:"comment,omitempty"`
	CommentOnTransition *bool                           `json:"commentOnTransition,omitempty"`
}

// UpdateServiceNowAutomationActionParamInput struct -- updates
type UpdateServiceNowAutomationActionParamInput struct {
	BaseURL      string                       `json:"baseUrl,omitempty"`
	User         string                       `json:"user,omitempty"`
	Password     string                       `json:"password,omitempty"`
	TicketFields *UpdateServiceNowFieldsInput `json:"ticketFields,omitempty"`
	ClientID     string                       `json:"clientId,omitempty"`
	ClientSecret string                       `json:"clientSecret,omitempty"`
}

// UpdateServiceNowUpdateTicketAutomationActionParamInput struct -- updates
type UpdateServiceNowUpdateTicketAutomationActionParamInput struct {
	BaseURL      string          `json:"baseUrl,omitempty"`
	User         string          `json:"user,omitempty"`
	Password     string          `json:"password,omitempty"`
	TableName    string          `json:"tableName,omitempty"`
	Fields       json.RawMessage `json:"fields,omitempty"`
	ClientID     string          `json:"clientId,omitempty"`
	ClientSecret string          `json:"clientSecret,omitempty"`
}

// UpdateAwsMessageAutomationActionParamsInput struct -- updates
type UpdateAwsMessageAutomationActionParamsInput struct {
	SNSTopicARN  string                                       `json:"snsTopicARN,omitempty"`
	Body         string                                       `json:"body,omitempty"`
	AccessMethod *AwsMessageAutomationActionAccessMethodInput `json:"accessMethod,omitempty"`
}

// UpdateAzureServiceBusAutomationActionParamsInput struct -- updates
type UpdateAzureServiceBusAutomationActionParamsInput struct {
	QueueURL     string                                            `json:"queueUrl,omitempty"`
	Body         string                                            `json:"body,omitempty"`
	AccessMethod *AzureServiceBusAutomationActionAccessMethodInput `json:"accessMethod,omitempty"`
}

// UpdateGooglePubSubAutomationActionParamsInput struct -- updates
type UpdateGooglePubSubAutomationActionParamsInput struct {
	ProjectID    string                                         `json:"projectId,omitempty"`
	TopicID      string                                         `json:"topicId,omitempty"`
	Body         string                                         `json:"body,omitempty"`
	AccessMethod *GooglePubSubAutomationActionAccessMethodInput `json:"accessMethod,omitempty"`
}

// UpdateServiceNowFieldsInput struct -- updates
type UpdateServiceNowFieldsInput struct {
	TableName         string          `json:"tableName,omitempty"`
	CustomFields      json.RawMessage `json:"customFields,omitempty"`
	Summary           string          `json:"summary,omitempty"`
	Description       string          `json:"description,omitempty"`
	AttachEvidenceCSV *bool           `json:"attachEvidenceCSV,omitempty"`
}

// UpdateJiraTicketFieldsInput struct -- updates
type UpdateJiraTicketFieldsInput struct {
	Summary                     string          `json:"summary,omitempty"`
	Description                 string          `json:"description,omitempty"`
	IssueType                   string          `json:"issueType,omitempty"`
	Assignee                    string          `json:"assignee,omitempty"`
	Components                  []string        `json:"components,omitempty"`
	FixVersion                  []string        `json:"fixVersion,omitempty"`
	Labels                      []string        `json:"labels,omitempty"`
	Priority                    string          `json:"priority,omitempty"`
	Project                     string          `json:"project,omitempty"`
	AlternativeDescriptionField string          `json:"alternativeDescriptionField,omitempty"`
	CustomFields                json.RawMessage `json:"customFields,omitempty"`
	AttachEvidenceCSV           *bool           `json:"attachEvidenceCSV,omitempty"`
}

// UpdateAutomationActionPayload struct -- updates
type UpdateAutomationActionPayload struct {
	AutomationAction AutomationAction `json:"automationAction,omitempty"`
}

// AutomationRuleTriggerSource enum
var AutomationRuleTriggerSource = []string{
	"ISSUES",
	"CLOUD_EVENTS",
	"CONTORL",
}

// AutomationRuleTriggerType enum
var AutomationRuleTriggerType = []string{
	"CREATED",
	"UPDATED",
	"RESOLVED",
	"REOPENED",
}

// AutomationRule struct -- updates
type AutomationRule struct {
	Action               AutomationAction        `json:"action"`
	Actions              []*AutomationRuleAction `json:"actions"`
	CreatedAt            string                  `json:"createdAt"`
	Description          string                  `json:"description,omitempty"`
	Enabled              bool                    `json:"enabled"`
	Filters              json.RawMessage         `json:"filters,omitempty"`
	ID                   string                  `json:"id"`
	Name                 string                  `json:"name"`
	OverrideActionParams json.RawMessage         `json:"overrideActionParams,omitempty"`
	Project              Project                 `json:"project,omitempty"`
	TriggerSource        string                  `json:"triggerSource"` // enum AutomationRuleTriggerSource
	TriggerType          []string                `json:"triggerType"`   // enum AutomationRuleTriggerType
}

// CreateAutomationRuleInput struct -- updates
type CreateAutomationRuleInput struct {
	Name          string                      `json:"name"`
	Description   string                      `json:"description,omitempty"`
	TriggerSource string                      `json:"triggerSource"` // enum AutomationRuleTriggerSource
	TriggerType   []string                    `json:"triggerType"`   // enum AutomationRuleTriggerType
	Filters       json.RawMessage             `json:"filters,omitempty"`
	Enabled       *bool                       `json:"enabled,omitempty"`
	ProjectID     string                      `json:"projectId,omitempty"`
	Actions       []AutomationRuleActionInput `json:"actions"`
}

// CreateAutomationRulePayload struct -- updates
type CreateAutomationRulePayload struct {
	AutomationRule AutomationRule `json:"automationRule"`
}

// UpdateAutomationRuleInput struct -- updates
type UpdateAutomationRuleInput struct {
	ID    string                    `json:"id"`
	Patch UpdateAutomationRulePatch `json:"patch"`
}

// UpdateAutomationRulePatch struct -- updates
type UpdateAutomationRulePatch struct {
	Name                 string          `json:"name,omitempty"`
	Description          string          `json:"description,omitempty"`
	TriggerSource        string          `json:"triggerSource,omitempty"` // enum AutomationRuleTriggerSource
	TriggerType          []string        `json:"triggerType,omitempty"`   // enum AutomationRuleTriggerType
	Filters              json.RawMessage `json:"filters,omitempty"`
	ActionID             string          `json:"actionId,omitempty"`
	OverrideActionParams json.RawMessage `json:"overrideActionParams,omitempty"`
	Enabled              *bool           `json:"enabled,omitempty"`
}

// UpdateAutomationRulePayload struct -- updates
type UpdateAutomationRulePayload struct {
	AutomationRule AutomationRule `json:"automationRule"`
}

// DeleteAutomationRuleInput struct -- updates
type DeleteAutomationRuleInput struct {
	ID string `json:"id"`
}

// DeleteAutomationRulePayload struct -- updates
type DeleteAutomationRulePayload struct {
	Stub string `json:"_stub,omitempty"`
}

// ServiceAccount struct -- updates
type ServiceAccount struct {
	AssignedProjects []*Project `json:"assignedProjects,omitempty"`
	ClientID         string     `json:"clientId"`
	ClientSecret     string     `json:"clientSecret"`
	CreatedAt        string     `json:"createdAt"`
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Scopes           []string   `json:"scopes"`
	Type             string     `json:"type"` //enum ServiceAccountType
	LastRotatedAt    string     `json:"lastRotatedAt"`
}

// ServiceAccountType enum
var ServiceAccountType = []string{
	"THIRD_PARTY",
	"SENSOR",
	"KUBERNETES_ADMISSION_CONTROLLER",
	"BROKER",
}

// CreateServiceAccountInput struct -- updates
type CreateServiceAccountInput struct {
	Name               string   `json:"name"`
	Scopes             []string `json:"scopes"`
	AssignedProjectIDs []string `json:"assignedProjectIds,omitempty"`
	Type               *string  `json:"type,omitempty"`
}

// CreateServiceAccountPayload struct -- updates
type CreateServiceAccountPayload struct {
	ServiceAccount ServiceAccount `json:"serviceAccount,omitempty"`
}

// DeleteServiceAccountInput struct -- updates
type DeleteServiceAccountInput struct {
	ID string `json:"id"`
}

// DeleteServiceAccountPayload struct -- updates
type DeleteServiceAccountPayload struct {
	Stub string `json:"_stub"`
}

// CICDScanPolicy struct -- updates -- added paramsType
type CICDScanPolicy struct {
	Builtin     bool              `json:"builtin"`
	Description string            `json:"description,omitempty"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	ParamsType  internal.EnumType `json:"paramsType"`
	Params      interface{}       `json:"params"`
}

// DiskScanVulnerabilitySeverity enum
var DiskScanVulnerabilitySeverity = []string{
	"INFORMATIONAL",
	"LOW",
	"MEDIUM",
	"HIGH",
	"CRITICAL",
}

// CICDScanPolicyParamsVulnerabilities struct - updates
type CICDScanPolicyParamsVulnerabilities struct {
	IgnoreUnfixed         bool     `json:"ignoreUnfixed"`
	PackageAllowList      []string `json:"packageAllowList"`
	PackageCountThreshold int      `json:"packageCountThreshold"`
	Severity              string   `json:"severity"` // enum DiskScanVulnerabilitySeverity
}

// CICDScanPolicyParamsSecrets struct -- updates
type CICDScanPolicyParamsSecrets struct {
	CountThreshold int      `json:"countThreshold"`
	PathAllowList  []string `json:"pathAllowList"`
}

// IACScanSeverity enum
var IACScanSeverity = []string{
	"INFORMATIONAL",
	"LOW",
	"MEDIUM",
	"HIGH",
	"CRITICAL",
}

// CICDScanPolicyParamsIAC struct -- updates
type CICDScanPolicyParamsIAC struct {
	BuiltinIgnoreTagsEnabled bool                         `json:"builtinIgnoreTagsEnabled"`
	CountThreshold           int                          `json:"countThreshold"`
	CustomIgnoreTags         []*CICDPolicyCustomIgnoreTag `json:"customIgnoreTags"`
	IgnoredRules             []*CloudConfigurationRule    `json:"ignoredRules"`
	SecurityFrameworks       []*SecurityFramework         `json:"securityFrameworks,omitempty"`
	SeverityThreshold        string                       `json:"severityThreshold"` // enum IACScanSeverity
}

// CICDPolicyCustomIgnoreTag struct -- updates
type CICDPolicyCustomIgnoreTag struct {
	IgnoreAllRules bool                      `json:"ignoreAllRules"`
	Key            string                    `json:"key"`
	Rules          []*CloudConfigurationRule `json:"rules"`
	Value          string                    `json:"value"`
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

// CloudConfigurationRule struct -- updates
type CloudConfigurationRule struct {
	Builtin                 *bool                                      `json:"builtin"`
	CloudProvider           string                                     `json:"cloudProvider,omitempty"` // enum CloudProvider
	Control                 *Control                                   `json:"control,omitempty"`
	CreatedBy               *User                                      `json:"createdBy,omitempty"`
	Description             string                                     `json:"description,omitempty"`
	Enabled                 *bool                                      `json:"enabled"`
	ExternalReferences      []*CloudConfigurationRuleExternalReference `json:"externalReferences,omitempty"`
	FunctionAsControl       *bool                                      `json:"functionAsControl"`
	GraphID                 string                                     `json:"graphId"`
	HasAutoRemediation      *bool                                      `json:"hasAutoRemediation"`
	IACMatchers             []*CloudConfigurationRuleMatcher           `json:"iacMatchers,omitempty"`
	ID                      string                                     `json:"id"`
	Name                    string                                     `json:"name"`
	OPAPolicy               string                                     `json:"opaPolicy,omitempty"`
	RemediationInstructions string                                     `json:"remediationInstructions,omitempty"`
	ScopeAccounts           []*CloudAccount                            `json:"scopeAccounts"` // removed omitempty
	SecuritySubCategories   []*SecuritySubCategory                     `json:"securitySubCategories"`
	ServiceType             string                                     `json:"serviceType,omitempty"` // enum CloudConfigurationRuleServiceType
	Severity                string                                     `json:"severity"`              // enum Severity
	ShortID                 string                                     `json:"shortId"`
	SubjectEntityType       string                                     `json:"subjectEntityType"`
	SupportsNRT             *bool                                      `json:"supportsNRT"`
	TargetNativeType        string                                     `json:"targetNativeType,omitempty"`
	TargetNativeTypes       []string                                   `json:"targetNativeTypes,omitempty"`
}

// SecurityFramework struct -- updates
type SecurityFramework struct {
	Builtin     bool               `json:"builtin"`
	Categories  []SecurityCategory `json:"categories"`
	Description string             `json:"description,omitempty"`
	Enabled     bool               `json:"enabled,omitempty"`
	ID          string             `json:"id"`
	Name        string             `json:"name"`
}

// CreateCICDScanPolicyPayload struct -- updates
type CreateCICDScanPolicyPayload struct {
	ScanPolicy *CICDScanPolicy `json:"scanPolicy,omitempty"`
}

// CreateCICDScanPolicyInput struct -- updates
type CreateCICDScanPolicyInput struct {
	Name                      string                                        `json:"name"`
	Description               string                                        `json:"description,omitempty"`
	DiskVulnerabilitiesParams *CreateCICDScanPolicyDiskVulnerabilitiesInput `json:"diskVulnerabilitiesParams,omitempty"`
	DiskSecretsParams         *CreateCICDScanPolicyDiskSecretsInput         `json:"diskSecretsParams,omitempty"`
	IACParams                 *CreateCICDScanPolicyIACInput                 `json:"iacParams,omitempty"`
}

// CreateCICDScanPolicyDiskVulnerabilitiesInput struct -- updates
type CreateCICDScanPolicyDiskVulnerabilitiesInput struct {
	Severity              string   `json:"severity"` // enum DiskScanVulnerabilitySeverity
	PackageCountThreshold int      `json:"packageCountThreshold"`
	IgnoreUnfixed         bool     `json:"ignoreUnfixed"`
	PackageAllowList      []string `json:"packageAllowList,omitempty"`
}

// CreateCICDScanPolicyDiskSecretsInput struct -- updates
type CreateCICDScanPolicyDiskSecretsInput struct {
	CountThreshold int      `json:"countThreshold"`
	PathAllowList  []string `json:"pathAllowList,omitempty"`
}

// CreateCICDScanPolicyIACInput struct -- updates
type CreateCICDScanPolicyIACInput struct {
	SeverityThreshold        string                                  `json:"severityThreshold"` // enum IACScanSeverity
	CountThreshold           int                                     `json:"countThreshold"`
	IgnoredRules             []string                                `json:"ignoredRules,omitempty"`
	BuiltinIgnoreTagsEnabled *bool                                   `json:"builtinIgnoreTagsEnabled,omitempty"`
	CustomIgnoreTags         []*CICDPolicyCustomIgnoreTagCreateInput `json:"customIgnoreTags,omitempty"`
	SecurityFrameworks       []string                                `json:"securityFrameworks,omitempty"`
}

// CICDPolicyCustomIgnoreTagCreateInput struct -- updates
type CICDPolicyCustomIgnoreTagCreateInput struct {
	Key            string   `json:"key"`
	Value          string   `json:"value"`
	RuleIDs        []string `json:"ruleIDs,omitempty"`
	IgnoreAllRules *bool    `json:"ignoreAllRules,omitempty"`
}

// UpdateCICDScanPolicyPayload struct -- updates
type UpdateCICDScanPolicyPayload struct {
	ScanPolicy *CICDScanPolicy `json:"scanPolicy,omitempty"`
}

// UpdateCICDScanPolicyInput struct -- updates
type UpdateCICDScanPolicyInput struct {
	ID    string                    `json:"id"`
	Patch UpdateCICDScanPolicyPatch `json:"patch"`
}

// UpdateCICDScanPolicyPatch struct -- updates
type UpdateCICDScanPolicyPatch struct {
	Name                      string                                        `json:"name,omitempty"`
	Description               string                                        `json:"description,omitempty"`
	DiskVulnerabilitiesParams *UpdateCICDScanPolicyDiskVulnerabilitiesPatch `json:"diskVulnerabilitiesParams,omitempty"`
	DiskSecretsParams         *UpdateCICDScanPolicyDiskSecretsPatch         `json:"diskSecretsParams,omitempty"`
	IACParams                 *UpdateCICDScanPolicyIACPatch                 `json:"iacParams,omitempty"`
}

// UpdateCICDScanPolicyDiskVulnerabilitiesPatch struct -- updates
// the omitempty tag is ignored because the patch requires empty values to remove parameter settings
type UpdateCICDScanPolicyDiskVulnerabilitiesPatch struct {
	Severity              string   `json:"severity"` // enum DiskScanVulnerabilitySeverity
	PackageCountThreshold int      `json:"packageCountThreshold"`
	IgnoreUnfixed         *bool    `json:"ignoreUnfixed"`
	PackageAllowList      []string `json:"packageAllowList"`
}

// UpdateCICDScanPolicyDiskSecretsPatch struct -- updates
// the omitempty tag is ignored because the patch requires empty values to remove parameter settings
type UpdateCICDScanPolicyDiskSecretsPatch struct {
	CountThreshold int      `json:"countThreshold"`
	PathAllowList  []string `json:"pathAllowList"`
}

// UpdateCICDScanPolicyIACPatch struct -- updates
// the omitempty tag is ignored because the patch requires empty values to remove parameter settings
type UpdateCICDScanPolicyIACPatch struct {
	SeverityThreshold        string                                  `json:"severityThreshold"` // enum IACScanSeverity
	CountThreshold           int                                     `json:"countThreshold"`
	IgnoredRules             []string                                `json:"ignoredRules"`
	BuiltinIgnoreTagsEnabled *bool                                   `json:"builtinIgnoreTagsEnabled"`
	CustomIgnoreTags         []*CICDPolicyCustomIgnoreTagUpdateInput `json:"customIgnoreTags"`
	SecurityFrameworks       []string                                `json:"securityFrameworks"`
}

// CICDPolicyCustomIgnoreTagUpdateInput struct -- updates
type CICDPolicyCustomIgnoreTagUpdateInput struct {
	Key            string   `json:"key"`
	Value          string   `json:"value"`
	RuleIDs        []string `json:"ruleIDs,omitempty"`
	IgnoreAllRules *bool    `json:"ignoreAllRules,omitempty"`
}

// DeleteCICDScanPolicyPayload struct -- updates
type DeleteCICDScanPolicyPayload struct {
	ID string `json:"id"`
}

// DeleteCICDScanPolicyInput struct -- updates
type DeleteCICDScanPolicyInput struct {
	ID string `json:"id"`
}

// WebhookAutomationActionParams struct -- updates -- added authenticationType for authentication
type WebhookAutomationActionParams struct {
	AuthenticationType internal.EnumType `json:"authenticationType"`
	Authentication     interface{}       `json:"authentication,omitempty"`
	Body               string            `json:"body"`
	ClientCertificate  string            `json:"clientCertificate,omitempty"`
	URL                string            `json:"url"`
}

// WebhookAutomationActionAuthenticationBasic struct -- updates
type WebhookAutomationActionAuthenticationBasic struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// WebhookAutomationActionAuthenticationTokenBearer struct -- updates
type WebhookAutomationActionAuthenticationTokenBearer struct {
	Token string `json:"token"`
}

/*
// CreateCloudConfigurationRuleInput struct -- updates
type CreateCloudConfigurationRuleInput struct {
	Name                    string                                   `json:"name"`
	Description             string                                   `json:"description,omitempty"`
	TargetNativeType        string                                   `json:"targetNativeType,omitempty"`
	TargetNativeTypes       []string                                 `json:"targetNativeTypes,omitempty"`
	OPAPolicy               string                                   `json:"opaPolicy,omitempty"`
	Severity                string                                   `json:"severity,omitempty"` // enum Severity
	Enabled                 *bool                                    `json:"enabled,omitempty"`
	RemediationInstructions string                                   `json:"remediationInstructions,omitempty"`
	ScopeAccountIDs         []string                                 `json:"scopeAccountIds,omitempty"`
	FunctionAsControl       *bool                                    `json:"functionAsControl,omitempty"`
	SecuritySubCategories   []string                                 `json:"securitySubCategories,omitempty"`
	IACMatchers             []*CreateCloudConfigurationRuleMatcherInput `json:"iacMatchers,omitempty"`
}
*/

// CreateCloudConfigurationRuleInput struct -- updates
type CreateCloudConfigurationRuleInput struct {
	Name                    string                                      `json:"name"`
	Description             string                                      `json:"description"`
	TargetNativeType        string                                      `json:"targetNativeType"`
	TargetNativeTypes       []string                                    `json:"targetNativeTypes"`
	OPAPolicy               string                                      `json:"opaPolicy"`
	Severity                string                                      `json:"severity"` // enum Severity
	Enabled                 *bool                                       `json:"enabled"`
	RemediationInstructions string                                      `json:"remediationInstructions"`
	ScopeAccountIDs         []string                                    `json:"scopeAccountIds"`
	FunctionAsControl       *bool                                       `json:"functionAsControl"`
	SecuritySubCategories   []string                                    `json:"securitySubCategories"`
	IACMatchers             []*CreateCloudConfigurationRuleMatcherInput `json:"iacMatchers"`
}

// CreateCloudConfigurationRuleMatcherInput struct -- updates
type CreateCloudConfigurationRuleMatcherInput struct {
	Type     string `json:"type"` // enum CloudConfigurationRuleMatcherType
	RegoCode string `json:"regoCode"`
}

// CreateCloudConfigurationRulePayload struct -- updates
type CreateCloudConfigurationRulePayload struct {
	Rule CloudConfigurationRule `json:"rule,omitempty"`
}

// CloudConfigurationRuleMatcher struct -- updates
type CloudConfigurationRuleMatcher struct {
	Enabled   *bool  `json:"enabled"`
	ID        string `json:"id"`
	RegoCode  string `json:"regoCode"`
	ShortName string `json:"shortName"`
	Type      string `json:"type"` // enum CloudConfigurationRuleMatcherType
}

// SecuritySubCategory struct -- updates
type SecuritySubCategory struct {
	Category                 SecurityCategory `json:"category"`
	Description              string           `json:"description"`
	ID                       string           `json:"id"`
	ResolutionRecommendation string           `json:"resolutionRecommendation,omitempty"`
	Title                    string           `json:"title"`
	ExternalID               string           `json:"external_id"`
}

// SecurityCategory struct -- updates
type SecurityCategory struct {
	Description   string                `json:"description"`
	Framework     SecurityFramework     `json:"framework"`
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	SubCategories []SecuritySubCategory `json:"subCategories"`
}

// DeleteCloudConfigurationRuleInput struct -- updates
type DeleteCloudConfigurationRuleInput struct {
	ID string `json:"id"`
}

// DeleteCloudConfigurationRulePayload struct -- updates
type DeleteCloudConfigurationRulePayload struct {
	Stub string `json:"_stub"`
}

// UpdateCloudConfigurationRuleInput struct -- updates
type UpdateCloudConfigurationRuleInput struct {
	ID    string                            `json:"id"`
	Patch UpdateCloudConfigurationRulePatch `json:"patch"`
}

// UpdateCloudConfigurationRulePatch struct -- updates
// the omitempty tag is ignored because the patch requires empty values to remove parameter settings
type UpdateCloudConfigurationRulePatch struct {
	Name                    string                                      `json:"name,omitempty"`
	Description             string                                      `json:"description,omitempty"`
	TargetNativeTypes       []string                                    `json:"targetNativeTypes,omitempty"`
	OPAPolicy               string                                      `json:"opaPolicy"`
	RemediationInstructions string                                      `json:"remediationInstructions,omitempty"`
	Severity                string                                      `json:"severity"`
	Enabled                 *bool                                       `json:"enabled"`
	ScopeAccountIds         []string                                    `json:"scopeAccountIds"`
	FunctionAsControl       *bool                                       `json:"functionAsControl"`
	SecuritySubCategories   []string                                    `json:"securitySubCategories,omitempty"`
	IACMatchers             []*UpdateCloudConfigurationRuleMatcherInput `json:"iacMatchers"`
}

// UpdateCloudConfigurationRuleMatcherInput struct -- updates
type UpdateCloudConfigurationRuleMatcherInput struct {
	Type     string `json:"type"` // enum CloudConfigurationRuleMatcherType
	RegoCode string `json:"regoCode"`
}

// UpdateCloudConfigurationRulePayload struct -- updates
type UpdateCloudConfigurationRulePayload struct {
	Rule CloudConfigurationRule `json:"rule,omitempty"`
}

// CloudConfigurationRuleMatcherType enum
var CloudConfigurationRuleMatcherType = []string{
	"TERRAFORM",
	"CLOUD_FORMATION",
	"KUBERNETES",
	"AZURE_RESOURCE_MANAGER",
	"DOCKER_FILE",
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

// Control struct -- updates
type Control struct {
	CreatedAt                    string                      `json:"createdAt,omitempty"`
	Description                  string                      `json:"description"`
	Enabled                      bool                        `json:"enabled"`
	EnabledForHBI                bool                        `json:"enabledForHBI"`
	EnabledForLBI                bool                        `json:"enabledForLBI"`
	EnabledForMBI                bool                        `json:"enabledForMBI"`
	EnabledForUnattributed       bool                        `json:"enabledForUnattributed"`
	ExternalReferences           []*ControlExternalReference `json:"externalReferences,omitempty"`
	ID                           string                      `json:"id"`
	Name                         string                      `json:"name"`
	Query                        interface{}                 `json:"query,omitempty"` // scalar
	ResolutionRecommendation     string                      `json:"resolutionRecommendation,omitempty"`
	ScopeProject                 Project                     `json:"scopeProject,omitempty"`
	ScopeQuery                   interface{}                 `json:"scopeQuery,omitempty"` // scalar
	SecuritySubCategories        []*SecuritySubCategory      `json:"securitySubCategories,omitempty"`
	Severity                     string                      `json:"severity,omitempty"` // enum Severity
	SourceCloudConfigurationRule CloudConfigurationRule      `json:"sourceCloudConfigurationRule,omitempty"`
	SupportsNRT                  bool                        `json:"supportsNRT"`
	Tags                         []string                    `json:"tags"`
	Type                         string                      `json:"type"` // enum type
}

// ControlExternalReference struct -- updates
type ControlExternalReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ControlType enum
var ControlType = []string{
	"SECURITY_GRAPH",
	"CLOUD_CONFIGURATION",
}

// CreateUserInput struct
type CreateUserInput struct {
	Email              string   `json:"email"`
	Name               string   `json:"name"`
	Role               string   `json:"role"`
	AssignedProjectIDs []string `json:"assignedProjectIds,omitempty"`
	SendEmailInvite    bool     `json:"sendEmailInvite"`
}

// CreateUserPayload struct
type CreateUserPayload struct {
	User User `json:"user"`
}

// UpdateUserInput struct
type UpdateUserInput struct {
	ID    string          `json:"id"`
	Patch UpdateUserPatch `json:"patch"`
}

// UpdateUserPatch struct
type UpdateUserPatch struct {
	Name               string   `json:"name,omitempty"`
	Email              string   `json:"email,omitempty"`
	Role               string   `json:"role,omitempty"`
	AssignedProjectIDs []string `json:"assignedProjectIds,omitempty"`
	IsSuspended        bool     `json:"isSuspended,omitempty"`
}

// UpdateUserPayload struct
type UpdateUserPayload struct {
	User User `json:"user,omitempty"`
}

// DeleteUserInput struct
type DeleteUserInput struct {
	ID string `json:"id"`
}

// DeleteUserPayload struct
type DeleteUserPayload struct {
	Stub string `json:"_stub"`
}

// CreateSecurityFrameworkInput struct
type CreateSecurityFrameworkInput struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Enabled     *bool                   `json:"enabled,omitempty"`
	Categories  []SecurityCategoryInput `json:"categories"`
}

// SecurityCategoryInput struct
type SecurityCategoryInput struct {
	ID            string                     `json:"id,omitempty"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description,omitempty"`
	ExternalID    string                     `json:"externalId,omitempty"`
	SubCategories []SecuritySubCategoryInput `json:"subCategories"`
}

// SecuritySubCategoryInput struct
// Description was made required because the API nullifies the value if not provided
type SecuritySubCategoryInput struct {
	ID                       string `json:"id,omitempty"`
	Title                    string `json:"title"`
	Description              string `json:"description"`
	ExternalID               string `json:"externalId,omitempty"`
	ResolutionRecommendation string `json:"resolutionRecommendation,omitempty"`
}

// CreateSecurityFrameworkPayload struct
type CreateSecurityFrameworkPayload struct {
	Framework SecurityFramework `json:"framework,omitempty"`
}

// DeleteSecurityFrameworkInput struct
type DeleteSecurityFrameworkInput struct {
	ID string `json:"id"`
}

// DeleteSecurityFrameworkPayload struct
type DeleteSecurityFrameworkPayload struct {
	Stub string `json:"_stub,omitempty"`
}

// UpdateSecurityFrameworkInput struct
type UpdateSecurityFrameworkInput struct {
	ID    string                 `json:"id"`
	Patch SecurityFrameworkPatch `json:"patch"`
}

// UpdateSecurityFrameworkPayload struct
type UpdateSecurityFrameworkPayload struct {
	Framework SecurityFramework `json:"framework,omitempty"`
}

// SecurityFrameworkPatch struct
type SecurityFrameworkPatch struct {
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description"`
	Enabled     *bool                   `json:"enabled,omitempty"`
	Categories  []SecurityCategoryInput `json:"categories"`
}

// CreateControlInput struct
type CreateControlInput struct {
	Query                    json.RawMessage `json:"query"`
	Name                     string          `json:"name"`
	Description              string          `json:"description,omitempty"`
	ResolutionRecommendation string          `json:"resolutionRecommendation,omitempty"`
	Severity                 string          `json:"severity"` // enum Severity
	ScopeQuery               json.RawMessage `json:"scopeQuery"`
	SecuritySubCategories    []string        `json:"securitySubCategories,omitempty"`
	ProjectID                string          `json:"projectId"`
}

// CreateControlPayload struct
type CreateControlPayload struct {
	Control Control `json:"control,omitempty"`
}

// DeleteControlInput struct
type DeleteControlInput struct {
	ID string `json:"id"`
}

// DeleteControlPayload struct
type DeleteControlPayload struct {
	Stub string `json:"_stub,omitempty"`
}

// UpdateControlInput struct
type UpdateControlInput struct {
	ID    string             `json:"id"`
	Patch UpdateControlPatch `json:"patch"`
}

// UpdateControlPatch struct
type UpdateControlPatch struct {
	Query                    json.RawMessage `json:"query,omitempty"`
	ScopeQuery               json.RawMessage `json:"scopeQuery,omitempty"`
	Name                     string          `json:"name,omitempty"`
	Description              string          `json:"description,omitempty"`
	ResolutionRecommendation string          `json:"resolutionRecommendation,omitempty"`
	Severity                 string          `json:"severity,omitempty"` // enum Severity
	Enabled                  *bool           `json:"enabled,omitempty"`
	EnabledForLBI            *bool           `json:"enabledForLBI,omitempty"`
	EnabledForMBI            *bool           `json:"enabledForMBI,omitempty"`
	EnabledForHBI            *bool           `json:"enabledForHBI,omitempty"`
	EnabledForUnattributed   *bool           `json:"enabledForUnattributed,omitempty"`
	SecuritySubCategories    []string        `json:"securitySubCategories,omitempty"`
}

// UpdateControlPayload struct
type UpdateControlPayload struct {
	Control Control `json:"control,omitempty"`
}

// UpdateControlsInput struct
type UpdateControlsInput struct {
	IDs                           []string             `json:"ids,omitempty"`
	Filters                       *ControlFilters      `json:"filters,omitempty"`
	Patch                         *UpdateControlsPatch `json:"patch,omitempty"`
	SecuritySubCategoriesToAdd    []string             `json:"securitySubCategoriesToAdd,omitempty"`
	SecuritySubCategoriesToRemove []string             `json:"securitySubCategoriesToRemove,omitempty"`
}

// ControlFilters struct
type ControlFilters struct {
	ID                  []string      `json:"id,omitempty"`
	Search              string        `json:"search,omitempty"`
	Type                []string      `json:"type,omitempty"` // enum ControlType
	Project             []string      `json:"project,omitempty"`
	CreatedBy           string        `json:"createdBy,omitempty"` // enum ControlCreatorType
	SecurityFramework   []string      `json:"securityFramework,omitempty"`
	SecuritySubCategory []string      `json:"securitySubCategory,omitempty"`
	SecurityCategory    []string      `json:"securityCategory,omitempty"`
	FrameworkCategory   []string      `json:"frameworkCategory,omitempty"`
	Tag                 string        `json:"tag,omitempty"`
	EntityType          string        `json:"entityType,omitempty"` // scalar
	Severity            string        `json:"severity,omitempty"`   // enum Severity
	WithIssues          *IssueFilters `json:"withIssues,omitempty"`
	Enabled             *bool         `json:"enabled,omitempty"`
	RiskEqualsAny       []string      `json:"riskEqualsAny,omitempty"`
	RiskEqualsAll       []string      `json:"riskEqualsAll,omitempty"`
}

// IssueFilters struct
type IssueFilters struct {
	ID                  []string           `json:"id,omitempty"`
	Search              string             `json:"search,omitempty"`
	SecurityFramework   string             `json:"securityFramework,omitempty"`
	SecuritySubCategory []string           `json:"securitySubCategory,omitempty"`
	SecurityCategory    []string           `json:"securityCategory,omitempty"`
	FrameworkCategory   []string           `json:"frameworkCategory,omitempty"`
	StackLayer          []string           `json:"stackLayer,omitempty"` // enum TechnologyStackLayer
	Project             []string           `json:"project,omitempty"`
	Severity            string             `json:"severity,omitempty"` // enum Severity
	Status              []string           `json:"status,omitempty"`   // enum IssueStatus
	RelatedEntity       IssueEntityFilters `json:"relatedEntity,omitempty"`
	SourceSecurityScan  string             `json:"sourceSecurityScan,omitempty"`
	SourceControl       []string           `json:"sourceControl,omitempty"`
	CreatedAt           IssueDateFilter    `json:"createdAt,omitempty"`
	ResolvedAt          IssueDateFilter    `json:"resolvedAt,omitempty"`
	ResolutionReason    []string           `json:"resolutionReason,omitempty"` // enum IssueResolutionReason
	DueAt               IssueDateFilter    `json:"dueAt,omitempty"`
	HasServiceTicket    *bool              `json:"hasServiceTicket,omitempty"`
	HasNote             *bool              `json:"hasNote,omitempty"`
	HasRemediation      *bool              `json:"hasRemediation,,omitempty"`
	SourceControlType   []string           `json:"sourceControlType"` // enum ControlType
	RiskEqualsAny       []string           `json:"riskEqualsAny,omitempty"`
	RiskEqualsAll       []string           `json:"riskEqualsAll,omitempty"`
}

// IssueDateFilter struct
type IssueDateFilter struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

// IssueEntityFilters struct
type IssueEntityFilters struct {
	ID              string               `json:"id,omitempty"`
	IDs             []string             `json:"ids,omitempty"`
	Type            string               `json:"type,omitempty"`   // scalar GraphEntityTypeValue
	Status          []string             `json:"status,omitempty"` // enum CloudResourceStatus
	Region          []string             `json:"region,omitempty"`
	SubscriptionID  []string             `json:"subscriptionId,omitempty"`
	ResourceGroupID []string             `json:"resourceGroupId,omitempty"`
	NativeType      []string             `json:"nativeType,omitempty"`
	CloudPlatform   []string             `json:"cloudPlatform,omitempty"` // enum CloudPlatform
	Tag             IssueEntityTagFilter `json:"tag,omitempty"`
}

// IssueEntityTagFilter struct
type IssueEntityTagFilter struct {
	ContainsAll       []IssueEntityTag `json:"containsAll,omitempty"`
	ContainsAny       []IssueEntityTag `json:"IssueEntityTag,omitempty"`
	DoesNotContainAll []IssueEntityTag `json:"doesNotContainAll,omitempty"`
	DoesNotContainAny []IssueEntityTag `json:"doesNotContainAny,omitempty"`
}

// IssueEntityTag struct
type IssueEntityTag struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// UpdateControlsPatch struct
type UpdateControlsPatch struct {
	Severity              string   `json:"severity,omitempty"`
	Enabled               *bool    `json:"enabled,omitempty"`
	SecuritySubCategories []string `json:"securitySubCategories,omitempty"`
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

// UpdateControlsPayload struct
type UpdateControlsPayload struct {
	Errors       []UpdateControlsError `json:"errors,omitempty"`
	FailCount    int                   `json:"failCount"`
	SuccessCount int                   `json:"successCount"`
}

// UpdateControlsError struct
type UpdateControlsError struct {
	Control Control `json:"control"`
	Reason  string  `json:"reason,omitempty"`
}

// UpdateCloudConfigurationRulesInput struct
type UpdateCloudConfigurationRulesInput struct {
	IDs                           []string                            `json:"ids,omitempty"`
	Filters                       *CloudConfigurationRuleFilters      `json:"filters,omitempty"`
	Patch                         *UpdateCloudConfigurationRulesPatch `json:"patch,omitempty"`
	SecuritySubCategoriesToAdd    []string                            `json:"securitySubCategoriesToAdd,omitempty"`
	SecuritySubCategoriesToRemove []string                            `json:"securitySubCategoriesToRemove,omitempty"`
}

// CloudConfigurationRuleFilters struct
type CloudConfigurationRuleFilters struct {
	Search              string   `json:"search,omitempty"`
	ScopeAccountIDs     []string `json:"scopeAccountIds,omitempty"`
	CloudProvider       []string `json:"cloudProvider,omitempty"`     // enum CloudProvider
	ServiceType         []string `json:"serviceType,omitempty"`       // enum CloudConfigurationRuleServiceType
	SubjectEntityType   []string `json:"subjectEntityType,omitempty"` // enum GraphEntityTypeValue
	Severity            []string `json:"severity,omitempty"`          // enum Severity
	Enabled             *bool    `json:"enabled,omitempty"`
	HasAutoRemediation  *bool    `json:"hasAutoRemediation,omitempty"`
	HasRemediation      *bool    `json:"hasRemediation,omitempty"`
	Benchmark           []string `json:"benchmark,omitempty"` // enum ConfigurationBenchmarkTypeId
	SecurityFramework   []string `json:"securityFramework,omitempty"`
	SecuritySubCategory []string `json:"securitySubCategory,omitempty"`
	SecurityCategory    []string `json:"securityCategory,omitempty"`
	FrameworkCategory   []string `json:"frameworkCategory,omitempty"`
	TargetNativeType    []string `json:"targetNativeType,omitempty"`
	CreatedBy           []string `json:"createdBy,omitempty"`
	IsOPAPolicy         *bool    `json:"isOPAPolicy,omitempty"`
	Project             []string `json:"project,omitempty"`
	MatcherType         []string `json:"matcherType,omitempty"` // enum CloudConfigurationRuleMatcherTypeFilter
	ID                  []string `json:"id,omitempty"`
	FunctionAsControl   *bool    `json:"functionAsControl,omitempty"`
	RiskEqualsAny       []string `json:"riskEqualsAny,omitempty"`
	RiskEqualsAll       []string `json:"riskEqualsAll,omitempty"`
}

// UpdateCloudConfigurationRulesPatch struct
type UpdateCloudConfigurationRulesPatch struct {
	Severity              string   `json:"severity,omitempty"` // enum Severity
	Enabled               *bool    `json:"enabled,omitempty"`
	FunctionAsControl     *bool    `json:"functionAsControl,omitempty"`
	SecuritySubCategories []string `json:"securitySubCategories,omitempty"`
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
}

// UpdateCloudConfigurationRulesPayload struct
type UpdateCloudConfigurationRulesPayload struct {
	Errors       []UpdateCloudConfigurationRulesError `json:"errors,omitempty"`
	FailCount    int                                  `json:"failCount"`
	SuccessCount int                                  `json:"successCount"`
}

// UpdateCloudConfigurationRulesError struct
type UpdateCloudConfigurationRulesError struct {
	Reason string                 `json:"reason,omitempty"`
	Rule   CloudConfigurationRule `json:"rule"`
}

// UpdateHostConfigurationRulesInput struct
type UpdateHostConfigurationRulesInput struct {
	IDs                           []string                          `json:"ids,omitempty"`
	Filters                       HostConfigurationRuleFilters      `json:"filters,omitempty"`
	Patch                         UpdateHostConfigurationRulesPatch `json:"patch,omitempty"`
	SecuritySubCategoriesToAdd    []string                          `json:"securitySubCategoriesToAdd,omitempty"`
	SecuritySubCategoriesToRemove []string                          `json:"securitySubCategoriesToRemove,omitempty"`
}

// HostConfigurationRuleFilters struct
type HostConfigurationRuleFilters struct {
	Search            string   `json:"search,omitempty"`
	Enabled           *bool    `json:"enabled,omitempty"`
	FrameworkCategory []string `json:"frameworkCategory,omitempty"`
	TargetPlatforms   []string `json:"targetPlatforms,omitempty"`
}

// UpdateHostConfigurationRulesPatch struct
type UpdateHostConfigurationRulesPatch struct {
	Enabled               *bool    `json:"enabled,omitempty"`
	SecuritySubCategories []string `json:"securitySubCategories,omitempty"`
}

// UpdateHostConfigurationRulesPayload struct
type UpdateHostConfigurationRulesPayload struct {
	Errors       []*UpdateHostConfigurationRulesError `json:"errors,omitempty"`
	FailCount    int                                  `json:"failCount"`
	SuccessCount int                                  `json:"successCount"`
}

// UpdateHostConfigurationRulesError struct
type UpdateHostConfigurationRulesError struct {
	Reason string                `json:"reason,omitempty"`
	Rule   HostConfigurationRule `json:"rule"`
}

// HostConfigurationRule struct
type HostConfigurationRule struct {
	Analytics             HostConfigurationRuleAnalytics `json:"analytics"`
	Builtin               bool                           `json:"builtin"`
	Description           string                         `json:"description,omitempty"`
	DirectOVAL            string                         `json:"directOVAL,omitempty"`
	Enabled               bool                           `json:"enabled"`
	ExternalID            string                         `json:"externalId,omitempty"`
	ID                    string                         `json:"id"`
	Name                  string                         `json:"name"`
	SecuritySubCategories []*SecuritySubCategory         `json:"securitySubCategories,omitempty"`
	ShortName             string                         `json:"shortName"`
	TargetPlatforms       []Technology                   `json:"targetPlatforms"`
}

// HostConfigurationRuleAnalytics struct
type HostConfigurationRuleAnalytics struct {
	ErrorCount       int `json:"errorCount"`
	FailCount        int `json:"failCount"`
	NotAssessedCount int `json:"notAssessedCount"`
	PassCount        int `json:"passCount"`
	TotalCount       int `json:"totalCount"`
}

// Technology struct
type Technology struct {
	Categories                []TechnologyCategory              `json:"categories"`
	CloudAccountCount         int                               `json:"cloudAccountCount"`
	CodeRepoCount             int                               `json:"codeRepoCount"`
	Color                     string                            `json:"color,omitempty"`
	DeploymentModel           string                            `json:"deploymentModel,omitempty"` // enum DeploymentModel
	Description               string                            `json:"description"`
	Icon                      string                            `json:"icon,omitempty"`
	ID                        string                            `json:"id"`
	InstanceEntityTypes       []string                          `json:"instanceEntityTypes"`
	Name                      string                            `json:"name"`
	Note                      string                            `json:"note,omitempty"`
	OnlyServiceUsageSupported bool                              `json:"onlyServiceUsageSupported"`
	ProjectCount              int                               `json:"projectCount"`
	PropertySections          []TechnologyPropertySection       `json:"propertySections"`
	ResourceCount             int                               `json:"resourceCount"`
	Risk                      string                            `json:"risk"`       // enum TechnologyRisk
	StackLayer                string                            `json:"stackLayer"` // enum TechnologyStackLayer
	Status                    string                            `json:"status"`     // enum TechnologyStatus
	Usage                     string                            `json:"usage"`      // enum TechnologyUsage
	VulnerabilityAnalytics    *TechnologyVulnerabilityAnalytics `json:"vulnerabilityAnalytics,omitempty"`
}

// TechnologyCategory struct
type TechnologyCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TechnologyPropertySection struct
type TechnologyPropertySection struct {
	Name       string               `json:"name"`
	Properties []TechnologyProperty `json:"properties"`
}

// TechnologyVulnerabilityAnalytics struct
type TechnologyVulnerabilityAnalytics struct {
	CategoryBreakdown []*VulnerabilityCountByCategory `json:"categoryBreakdown,omitempty"`
	TotalCount        int                             `json:"totalCount"`
	YearBreakdown     []*VulnerabilityCountByYear     `json:"yearBreakdown,omitempty"`
}

// TechnologyProperty struct
type TechnologyProperty struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// VulnerabilityCountByCategory struct
type VulnerabilityCountByCategory struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// VulnerabilityCountByYear struct
type VulnerabilityCountByYear struct {
	Count int `json:"count"`
	Year  int `json:"year"`
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

// CreateHostConfigurationRuleInput struct
type CreateHostConfigurationRuleInput struct {
	Name                  string   `json:"name"`
	Description           string   `json:"description,omitempty"`
	DirectOVAL            string   `json:"directOVAL"`
	TargetPlatformIds     []string `json:"targetPlatformIds,omitempty"`
	Enabled               *bool    `json:"enabled,omitempty"`
	SecuritySubCategories []string `json:"securitySubCategories,omitempty"`
}

// CreateHostConfigurationRulePayload struct
type CreateHostConfigurationRulePayload struct {
	Rule HostConfigurationRule `json:"rule,omitempty"`
}

// DeleteHostConfigurationRuleInput struct
type DeleteHostConfigurationRuleInput struct {
	ID string `json:"id"`
}

// UpdateHostConfigurationRuleInput struct
type UpdateHostConfigurationRuleInput struct {
	ID    string                           `json:"id"`
	Patch UpdateHostConfigurationRulePatch `json:"patch"`
}

// UpdateHostConfigurationRulePatch struct
type UpdateHostConfigurationRulePatch struct {
	Enabled               *bool    `json:"enabled,omitempty"`
	SecuritySubCategories []string `json:"securitySubCategories,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Description           string   `json:"description,omitempty"`
	DirectOVAL            string   `json:"directOVAL,omitempty"`
	TargetPlatformIds     []string `json:"targetPlatformIds,omitempty"`
}

// DeleteHostConfigurationRulePayload struct
type DeleteHostConfigurationRulePayload struct {
	Stub string `json:"_stub,omitempty"`
}

// UpdateHostConfigurationRulePayload struct
type UpdateHostConfigurationRulePayload struct {
	Rule HostConfigurationRule `json:"rule,omitempty"`
}

// CloudAccountFilters struct
type CloudAccountFilters struct {
	ID                          []string `json:"id,omitempty"`
	Search                      []string `json:"search,omitempty"`
	ProjectID                   string   `json:"projectId,omitempty"`
	CloudProvider               []string `json:"cloudProvider,omitempty"` // enum CloudProvider
	Status                      []string `json:"status,omitempty"`        // enum CloudAccountStatus
	ConnectorID                 []string `json:"connectorId,omitempty"`
	ConnectorIssueID            []string `json:"connectorIssueId,omitempty"`
	AssignedToProject           *bool    `json:"assignedToProject,omitempty"`
	HasMultipleConnectorSources *bool    `json:"hasMultipleConnectorSources"`
}

// CloudAccountConnection struct
type CloudAccountConnection struct {
	Edges      []*CloudAccountEdge `json:"edges,omitempty"`
	Nodes      []*CloudAccount     `json:"nodes,omitempty"`
	PageInfo   PageInfo            `json:"pageInfo"`
	TotalCount int                 `json:"totalCount"`
}

// KubernetesClusterConnection struct
type KubernetesClusterConnection struct {
	Nodes      []*KubernetesCluster `json:"nodes,omitempty"`
	PageInfo   PageInfo             `json:"pageInfo"`
	TotalCount int                  `json:"totalCount"`
}

// UserConnection struct
type UserConnection struct {
	Nodes      []*User  `json:"nodes,omitempty"`
	PageInfo   PageInfo `json:"pageInfo"`
	TotalCount int      `json:"totalCount"`
}

// CloudAccountEdge struct
type CloudAccountEdge struct {
	Cursor string       `json:"cursor"`
	Node   CloudAccount `json:"node"`
}

// CloudConfigurationRuleOrder struct
type CloudConfigurationRuleOrder struct {
	Direction string `json:"direction"` // enum OrderDirection
	Field     string `json:"field"`     // enum CloudConfigurationRuleOrderField
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

// CloudConfigurationRuleConnection struct
type CloudConfigurationRuleConnection struct {
	AnalyticsUpdatedAt    string                        `json:"analyticsUpdatedAt"`
	Edges                 []*CloudConfigurationRuleEdge `json:"edges,omitempty"`
	EnabledAsControlCount int                           `json:"enabledAsControlCount"`
	Nodes                 []*CloudConfigurationRule     `json:"nodes,omitempty"`
	PageInfo              PageInfo                      `json:"pageInfo"`
	TotalCount            int                           `json:"totalCount"`
}

// CloudConfigurationRuleEdge struct
type CloudConfigurationRuleEdge struct {
	Cursor string                 `json:"cursor"`
	Node   CloudConfigurationRule `json:"node"`
}

// CloudConfigurationRuleExternalReference struct
type CloudConfigurationRuleExternalReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HostConfigurationRuleOrder struct
type HostConfigurationRuleOrder struct {
	Direction string `json:"direction"` // enum OrderDirection
	Field     string `json:"field"`     // enum HostConfigurationRuleOrderField
}

// HostConfigurationRuleOrderField enum
var HostConfigurationRuleOrderField = []string{
	"FAILED_CHECK_COUNT",
	"NAME",
}

// HostConfigurationRuleConnection struct
type HostConfigurationRuleConnection struct {
	Edges      []*HostConfigurationRuleEdge `json:"edges,omitempty"`
	Nodes      []*HostConfigurationRule     `json:"nodes,omitempty"`
	PageInfo   PageInfo                     `json:"pageInfo"`
	TotalCount int                          `json:"totalCount"`
}

// HostConfigurationRuleEdge struct
type HostConfigurationRuleEdge struct {
	Cursor string                `json:"cursor"`
	Node   HostConfigurationRule `json:"node"`
}

// GraphEntityQueryInput struct
type GraphEntityQueryInput struct {
	Type                []string                       `json:"type"` //enum GraphEntityType
	As                  *string                        `json:"as,omitempty"`
	Relationships       []*GraphRelationshipQueryInput `json:"relationships,omitempty"`
	Where               interface{}                    `json:"where,omitempty"` //scalar GraphPropertyPredicate
	Select              *bool                          `json:"select,omitempty"`
	Aggregate           *bool                          `json:"aggregate,omitempty"`
	AggregateConstraint []interface{}                  `json:"aggregateConstraint,omitempty"` //scalar GraphPropertyPredicate
	BlockName           *string                        `json:"blockName,omitempty"`
	BlockExpanded       *bool                          `json:"blockExpanded,omitempty"`
}

// GraphRelationshipQueryInput struct
type GraphRelationshipQueryInput struct {
	Type     []GraphDirectedRelationshipTypeInput `json:"type"`
	With     GraphEntityQueryInput                `json:"with"`
	Negate   *bool                                `json:"negate,omitempty"`
	Optional *bool                                `json:"optional,omitempty"`
}

// GraphDirectedRelationshipTypeInput struct
type GraphDirectedRelationshipTypeInput struct {
	Type    string `json:"type"` //enum GraphRelationshipType
	Reverse *bool  `json:"reverse,omitempty"`
}

// GraphEntity struct
// We deviate from the GraphQL schema and omit unused fields due to the high number of vertices
type GraphEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GraphSearchResult struct
type GraphSearchResult struct {
	AggregateCount *int          `json:"aggregateCount"`
	Entities       []GraphEntity `json:"entities"`
}

// GraphSearchResultConnection struct
type GraphSearchResultConnection struct {
	Nodes           []*GraphSearchResult     `json:"nodes"`
	Edges           []*GraphSearchResultEdge `json:"edges"`
	PageInfo        PageInfo                 `json:"pageInfo"`
	TotalCount      int                      `json:"totalCount"`
	MaxCountReached bool                     `json:"maxCountReached"`
}

// GraphSearchResultEdge struct
type GraphSearchResultEdge struct {
	AggregateCount *int          `json:"aggregateCount"`
	Entities       []GraphEntity `json:"entities"`
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

// AutomationRuleAction struct
type AutomationRuleAction struct {
	ID                   string      `json:"id"`
	ActionTemplateParams interface{} `json:"actionTemplateParams,omitempty"` // union ActionTemplateParams
	ActionTemplateType   string      `json:"actionTemplateType"`             // enum ActionTemplateType
	Integration          Integration `json:"integration"`
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

// Integration struct
// usedByRules: [AutomationRule!]! was omitted from the struct
type Integration struct {
	CreatedAt                 string      `json:"createdAt"`
	ID                        string      `json:"id"`
	IsAccessibleToAllProjects *bool       `json:"isAccessibleToAllProjects,omitempty"`
	Name                      string      `json:"name"`
	Params                    interface{} `json:"params"` // union IntegrationParams
	Project                   *Project    `json:"project,omitempty"`
	Type                      string      `json:"type"` // enum IntegrationType
	UpdatedAt                 string      `json:"updatedAt"`
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

// AutomationRuleActionInput struct
type AutomationRuleActionInput struct {
	ID                   string                    `json:"id"`
	IntegrationID        string                    `json:"integrationId"`
	ActionTemplateParams ActionTemplateParamsInput `json:"actionTemplateParams,omitempty"`
	ActionTemplateType   ActionTemplateType        `json:"actionTemplateType"`
}

// ActionTemplateParamsInput struct
type ActionTemplateParamsInput struct {
	AwsSNS                  AwsSNSActionTemplateParamsInput                  `json:"awsSNS,omitempty"`
	Email                   EmailActionTemplateParamsInput                   `json:"email,omitempty"`
	Webhook                 WebhookActionTemplateParamsInput                 `json:"webhook,omitempty"`
	Slack                   SlackActionTemplateParamsInput                   `json:"slack,omitempty"`
	SlackBot                SlackBotActionTemplateParamsInput                `json:"slackBot,omitempty"`
	AzureServiceBus         AzureServiceBusActionTemplateParamsInput         `json:"azureServiceBus,omitempty"`
	GoogleChat              GoogleChatActionTemplateParamsInput              `json:"googleChat,omitempty"`
	GcpPubSub               GcpPubSubActionTemplateParamsInput               `json:"gcpPubSub,omitempty"`
	PagerDutyCreateIncident PagerDutyActionCreateIncidentTemplateParamsInput `json:"pagerDutyCreateIncident,omitempty"`
	JiraCreateTicket        JiraActionCreateTicketTemplateParamsInput        `json:"jiraCreateTicket,omitempty"`
	JiraAddComment          JiraActionAddCommentTemplateParamsInput          `json:"jiraAddComment,omitempty"`
	JiraTransitionTicket    JiraActionTransitionTicketTemplateParamsInput    `json:"jiraTransitionTicket,omitempty"`
	ServiceNowCreateTicket  ServiceNowActionCreateTicketTemplateParamsInput  `json:"serviceNowCreateTicket,omitempty"`
	ServiceNowUpdateTicket  ServiceNowActionUpdateTicketTemplateParamsInput  `json:"serviceNowUpdateTicket,omitempty"`
	OpsgenieCreateAlert     OpsgenieCreateAlertTemplateParamsInput           `json:"opsgenieCreateAlert,omitempty"`
	OpsgenieCloseAlert      OpsgenieCloseAlertTemplateParamsInput            `json:"opsgenieCloseAlert,omitempty"`
	ClickUpCreateTask       ClickUpCreateTaskActionTemplateParamsInput       `json:"clickUpCreateTask,omitempty"`
}

// AwsSNSActionTemplateParamsInput struct
type AwsSNSActionTemplateParamsInput struct {
	Body string `json:"body"`
}

// EmailActionTemplateParamsInput struct
type EmailActionTemplateParamsInput struct {
	Note              string   `json:"note,omitempty"`
	To                []string `json:"to"`
	CC                []string `json:"cc,omitempty"`
	AttachEvidenceCSV *bool    `json:"attachEvidenceCSV,omitempty"`
}

// WebhookActionTemplateParamsInput struct
type WebhookActionTemplateParamsInput struct {
	Body    string               `json:"body"`
	Headers []WebhookHeaderInput `json:"headers,omitempty"`
}

// SlackActionTemplateParamsInput struct
type SlackActionTemplateParamsInput struct {
	Note string `json:"note,omitempty"`
}

// SlackBotActionTemplateParamsInput struct
type SlackBotActionTemplateParamsInput struct {
	Note    string `json:"note,omitempty"`
	Channel string `json:"channel"`
}

// AzureServiceBusActionTemplateParamsInput struct
type AzureServiceBusActionTemplateParamsInput struct {
	Body string `json:"body"`
}

// GoogleChatActionTemplateParamsInput struct
type GoogleChatActionTemplateParamsInput struct {
	Note string `json:"note,omitempty"`
}

// GcpPubSubActionTemplateParamsInput struct
type GcpPubSubActionTemplateParamsInput struct {
	Body string `json:"body"`
}

// PagerDutyActionCreateIncidentTemplateParamsInput struct
type PagerDutyActionCreateIncidentTemplateParamsInput struct {
	Payload string `json:"payload"`
}

// JiraActionCreateTicketTemplateParamsInput struct
type JiraActionCreateTicketTemplateParamsInput struct {
	Fields CreateJiraTicketFieldsInput `json:"fields,omitempty"`
}

// JiraActionAddCommentTemplateParamsInput struct
type JiraActionAddCommentTemplateParamsInput struct {
	ProjectKey      string `json:"projectKey,omitempty"`
	Comment         string `json:"comment"`
	AddIssuesReport bool   `json:"addIssuesReport"`
}

// JiraActionTransitionTicketTemplateParamsInput struct
type JiraActionTransitionTicketTemplateParamsInput struct {
	Project             string          `json:"project"`
	TransitionID        string          `json:"transitionId"`
	AdvancedFields      json.RawMessage `json:"advancedFields,omitempty"`
	Comment             string          `json:"comment,omitempty"`
	CommentOnTransition *bool           `json:"commentOnTransition,omitempty"`
	AttachEvidenceCSV   *bool           `json:"attachEvidenceCSV,omitempty"`
}

// ServiceNowActionCreateTicketTemplateParamsInput struct
type ServiceNowActionCreateTicketTemplateParamsInput struct {
	Fields CreateServiceNowFieldsInput `json:"fields"`
}

// ServiceNowActionUpdateTicketTemplateParamsInput struct
type ServiceNowActionUpdateTicketTemplateParamsInput struct {
	TableName          string          `json:"tableName"`
	Fields             json.RawMessage `json:"fields,omitempty"`
	AttachIssuesReport *bool           `json:"attachIssuesReport,omitempty"`
}

// OpsgenieCreateAlertTemplateParamsInput struct
type OpsgenieCreateAlertTemplateParamsInput struct {
	Body string `json:"body"`
}

// OpsgenieCloseAlertTemplateParamsInput struct
type OpsgenieCloseAlertTemplateParamsInput struct {
	Body string `json:"body"`
}

// ClickUpCreateTaskActionTemplateParamsInput struct
type ClickUpCreateTaskActionTemplateParamsInput struct {
	ListID string `json:"listId"`
	Body   string `json:"body"`
}

// WebhookHeaderInput struct
type WebhookHeaderInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AWSSNSIntegrationParams struct
type AWSSNSIntegrationParams struct {
	AccessConnector Connector `json:"accessConnector,omitempty"`
	AccessMethod    string    `json:"accessMethod"` // enum AwsSNSIntegrationAccessMethodType
	CustomerRoleARN string    `json:"customerRoleARN"`
	TopicARN        string    `json:"topicARN"`
}

// AzureServiceBusIntegrationParams struct
type AzureServiceBusIntegrationParams struct {
	AccessConnector         Connector                                  `json:"accessConnector,omitempty"`
	AccessMethod            AzureServiceBusIntegrationAccessMethodType `json:"accessMethod"`
	ConnectionStringWithSAS string                                     `json:"connectionStringWithSAS,omitempty"`
	QueueURL                string                                     `json:"queueUrl"`
}

// GcpPubSubIntegrationParams struct
type GcpPubSubIntegrationParams struct {
	AccessConnector   Connector                            `json:"accessConnector,omitempty"`
	AccessMethod      GcpPubSubIntegrationAccessMethodType `json:"accessMethod"`
	ProjectID         string                               `json:"projectId"`
	ServiceAccountKey json.RawMessage                      `json:"serviceAccountKey.omitempty"`
	TopicID           string                               `json:"topicId"`
}

// JiraIntegrationParams struct
type JiraIntegrationParams struct {
	Authorization interface{}               `json:"authorization"` // union JiraIntegrationAuthorization
	OnPremConfig  OnPremIntegrationConfig   `json:"onPremConfig,omitempty"`
	ServerType    string                    `json:"serverType"` // enum JiraServerType
	TLSConfig     AutomationActionTLSConfig `json:"tlsConfig"`
	URL           string                    `json:"url"`
}

// PagerDutyIntegrationParams struct
type PagerDutyIntegrationParams struct {
	IntegrationKey string `json:"integrationKey"`
}

// ServiceNowIntegrationParams struct
type ServiceNowIntegrationParams struct {
	Authorization interface{} `json:"authorization"`
	URL           string      `json:"url"`
}

// WebhookIntegrationParams struct
type WebhookIntegrationParams struct {
	Authorization interface{}               `json:"authorization,omitempty"` // union WebhookIntegrationAuthorization
	Headers       []WebhookHeader           `json:"headers,omitempty"`
	OnPremConfig  OnPremIntegrationConfig   `json:"onPremConfig,omitempty"`
	TLSConfig     AutomationActionTLSConfig `json:"tlsConfig,omitempty"`
	URL           string                    `json:"url"`
}

// SlackIntegrationParams struct
type SlackIntegrationParams struct {
	Channel string `json:"channel,omitempty"`
	URL     string `json:"url"`
}

// SlackBotIntegrationParams struct
type SlackBotIntegrationParams struct {
	Token string `json:"token"`
}

// OpsgenieIntegrationParams struct
type OpsgenieIntegrationParams struct {
	Key string `json:"key"`
}

// ClickUpIntegrationParams struct
type ClickUpIntegrationParams struct {
	Key string `json:"key"`
}

// OnPremIntegrationConfig struct
type OnPremIntegrationConfig struct {
	IsOnPrem     bool   `json:"isOnPrem"`
	TunnelDomain string `json:"tunnelDomain,omitempty"`
	TunnelToken  string `json:"tunnelToken,omitempty"`
}

// AutomationActionTLSConfig struct {
type AutomationActionTLSConfig struct {
	AllowInsecureTLS               *bool  `json:"allowInsecureTLS,omitempty"`
	ClientCertificateAndPrivateKey string `json:"clientCertificateAndPrivateKey,omitempty"`
	ServerCA                       string `json:"serverCA,omitempty"`
}

// WebhookHeader struct
type WebhookHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// JiraIntegrationBasicAuthorization struct
type JiraIntegrationBasicAuthorization struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// JiraIntegrationTokenBearerAuthorization struct
type JiraIntegrationTokenBearerAuthorization struct {
	Token string `json:"token"`
}

// ServiceNowIntegrationBasicAuthorization struct
type ServiceNowIntegrationBasicAuthorization struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// ServiceNowIntegrationOAuthAuthorization struct
type ServiceNowIntegrationOAuthAuthorization struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Password     string `json:"password"`
	Username     string `json:"username"`
}

// WebhookIntegrationBasicAuthorization struct
type WebhookIntegrationBasicAuthorization struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// WebhookIntegrationBearerAuthorization struct
type WebhookIntegrationBearerAuthorization struct {
	Token string `json:"token"`
}

// CreateIntegrationInput struct
type CreateIntegrationInput struct {
	Name                      string                       `json:"name"`
	Type                      string                       `json:"type"` // enum IntegrationType
	ProjectID                 string                       `json:"projectId,omitempty"`
	Params                    CreateIntegrationParamsInput `json:"params"`
	IsAccessibleToAllProjects *bool                        `json:"isAccessibleToAllProjects,omitempty"`
}

// CreateIntegrationParamsInput struct
type CreateIntegrationParamsInput struct {
	AwsSNS          CreateAwsSNSIntegrationParamsInput          `json:"awsSNS,omitempty"`
	Webhook         CreateWebhookIntegrationParamsInput         `json:"webhook,omitempty"`
	Slack           CreateSlackIntegrationParamsInput           `json:"slack,omitempty"`
	SlackBot        CreateSlackBotIntegrationParamsInput        `json:"slackBot,omitempty"`
	AzureServiceBus CreateAzureServiceBusIntegrationParamsInput `json:"azureServiceBus,omitempty"`
	GcpPubSub       CreateGcpPubSubIntegrationParamsInput       `json:"gcpPubSub,omitempty"`
	PagerDuty       CreatePagerDutyIntegrationParamsInput       `json:"pagerDuty,omitempty"`
	Jira            CreateJiraIntegrationParamsInput            `json:"jira,omitempty"`
	ServiceNow      CreateServiceNowIntegrationParamsInput      `json:"serviceNow,omitempty"`
	Opsgenie        CreateOpsgenieIntegrationParamsInput        `json:"opsgenie,omitempty"`
	ClickUp         CreateClickUpIntegrationParamsInput         `json:"clickUp,omitempty"`
}

// CreateAwsSNSIntegrationParamsInput struct
type CreateAwsSNSIntegrationParamsInput struct {
	TopicARN     string                             `json:"topicARN"`
	AccessMethod AwsSNSIntegrationAccessMethodInput `json:"accessMethod"`
}

// CreateWebhookIntegrationParamsInput struct
type CreateWebhookIntegrationParamsInput struct {
	URL           string                               `json:"url"`
	IsOnPrem      *bool                                `json:"isOnPrem,omitempty"`
	Authorization WebhookIntegrationAuthorizationInput `json:"authorization,omitempty"`
	Headers       []WebhookHeaderInput                 `json:"headers,omitempty"`
	TLSConfig     IntegrationTLSConfigInput            `json:"tlsConfig,omitempty"`
}

// CreateSlackIntegrationParamsInput struct
type CreateSlackIntegrationParamsInput struct {
	URL string `json:"url"`
}

// CreateSlackBotIntegrationParamsInput struct
type CreateSlackBotIntegrationParamsInput struct {
	Token string `json:"token"`
}

// CreateAzureServiceBusIntegrationParamsInput struct
type CreateAzureServiceBusIntegrationParamsInput struct {
	QueueUrl     string                                      `json:"queueUrl"`
	AccessMethod AzureServiceBusIntegrationAccessMethodInput `json:"accessMethod"`
}

// CreateGcpPubSubIntegrationParamsInput struct
type CreateGcpPubSubIntegrationParamsInput struct {
	ProjectID    string                                   `json:"projectId"`
	TopicID      string                                   `json:"topicId"`
	AccessMethod GooglePubSubIntegrationAccessMethodInput `json:"accessMethod"`
}

// CreatePagerDutyIntegrationParamsInput struct
type CreatePagerDutyIntegrationParamsInput struct {
	IntegrationKey string `json:"integrationKey"`
}

// CreateJiraIntegrationParamsInput struct
type CreateJiraIntegrationParamsInput struct {
	ServerURL     string                            `json:"serverUrl"`
	ServerType    string                            `json:"serverType"`
	IsOnPrem      bool                              `json:"isOnPrem"`
	TLSConfig     IntegrationTLSConfigInput         `json:"tlsConfig,omitempty"`
	Authorization JiraIntegrationAuthorizationInput `json:"authorization"`
}

// CreateServiceNowIntegrationParamsInput struct
type CreateServiceNowIntegrationParamsInput struct {
	URL           string                                  `json:"url"`
	Authorization ServiceNowIntegrationAuthorizationInput `json:"authorization"`
}

// CreateOpsgenieIntegrationParamsInput struct
type CreateOpsgenieIntegrationParamsInput struct {
	Key string `json:"key"`
}

// CreateClickUpIntegrationParamsInput struct
type CreateClickUpIntegrationParamsInput struct {
	Key string `json:"key"`
}

// AwsSNSIntegrationAccessMethodInput struct
type AwsSNSIntegrationAccessMethodInput struct {
	Type              string `json:"type"` // enum AwsSNSIntegrationAccessMethodType
	AccessConnectorId string `json:"accessConnectorId,omitempty"`
	CustomerRoleARN   string `json:"customerRoleARN,omitempty"`
}

// WebhookIntegrationAuthorizationInput struct
type WebhookIntegrationAuthorizationInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// IntegrationTLSConfigInput struct
type IntegrationTLSConfigInput struct {
	AllowInsecureTLS               *bool  `json:"allowInsecureTLS,omitempty"`
	ServerCA                       string `json:"serverCA,omitempty"`
	ClientCertificateAndPrivateKey string `json:"clientCertificateAndPrivateKey,omitempty"`
}

// AzureServiceBusIntegrationAccessMethodInput struct
type AzureServiceBusIntegrationAccessMethodInput struct {
	Type                    string `json:"type"` // enum AzureServiceBusIntegrationAccessMethodType
	AccessConnectorId       string `json:"accessConnectorId,omitempty"`
	ConnectionStringWithSas string `json:"connectionStringWithSas,omitempty"`
}

// GooglePubSubIntegrationAccessMethodInput struct
type GooglePubSubIntegrationAccessMethodInput struct {
	Type              string          `json:"type"` // enum GcpPubSubIntegrationAccessMethodType
	AccessConnectorId string          `json:"accessConnectorId,omitempty"`
	ServiceAccountKey json.RawMessage `json:"serviceAccountKey,omitempty"`
}

// JiraServerType enum
var JiraServerType = []string{
	"CLOUD",
	"SELF_HOSTED",
}

// JiraIntegrationAuthorizationInput struct
type JiraIntegrationAuthorizationInput struct {
	Username            string `json:"username,omitempty"`
	Password            string `json:"password,omitempty"`
	PersonalAccessToken string `json:"personalAccessToken,omitempty"`
}

// ServiceNowIntegrationAuthorizationInput struct
type ServiceNowIntegrationAuthorizationInput struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
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

// CreateIntegrationPayload struct
type CreateIntegrationPayload struct {
	Integration Integration `json:"integration"`
}

// DeleteIntegrationInput struct
type DeleteIntegrationInput struct {
	ID string `json:"id"`
}

// DeleteIntegrationPayload struct
type DeleteIntegrationPayload struct {
	Stub string `json:"_stub"`
}

// UpdateIntegrationInput struct
type UpdateIntegrationInput struct {
	ID    string                 `json:"id"`
	Patch UpdateIntegrationPatch `json:"patch"`
}

// UpdateIntegrationPatch struct
type UpdateIntegrationPatch struct {
	Name   string                       `json:"name,omitempty"`
	Params UpdateIntegrationPatchParams `json:"params,omitempty"`
}

// UpdateIntegrationPatchParams struct
type UpdateIntegrationPatchParams struct {
	AwsSNS          UpdateAwsSNSIntegrationParamsInput          `json:"awsSNS,omitempty"`
	Webhook         UpdateWebhookIntegrationParamsInput         `json:"webhook,omitempty"`
	Slack           UpdateSlackIntegrationParamsInput           `json:"slack,omitempty"`
	SlackBot        UpdateSlackBotIntegrationParamsInput        `json:"slackBot,omitempty"`
	AzureServiceBus UpdateAzureServiceBusIntegrationParamsInput `json:"azureServiceBus,omitempty"`
	GcpPubSub       UpdateGcpPubSubIntegrationParamsInput       `json:"gcpPubSub,omitempty"`
	PagerDuty       UpdatePagerDutyIntegrationParamsInput       `json:"pagerDuty,omitempty"`
	Jira            UpdateJiraIntegrationParamsInput            `json:"jira,omitempty"`
	ServiceNow      UpdateServiceNowIntegrationParamsInput      `json:"serviceNow,omitempty"`
	Opsgenie        UpdateOpsgenieIntegrationParamsInput        `json:"opsgenie,omitempty"`
	ClickUp         UpdateClickUpIntegrationParamsInput         `json:"clickUp,omitempty"`
}

// UpdateAwsSNSIntegrationParamsInput struct
type UpdateAwsSNSIntegrationParamsInput struct {
	TopicARN     string                             `json:"topicARN,omitempty"`
	AccessMethod AwsSNSIntegrationAccessMethodInput `json:"accessMethod,omitempty"`
}

// UpdateWebhookIntegrationParamsInput struct
type UpdateWebhookIntegrationParamsInput struct {
	URL           string                               `json:"url,omitempty"`
	IsOnPrem      *bool                                `json:"isOnPrem,omitempty"`
	Authorization WebhookIntegrationAuthorizationInput `json:"authorization,omitempty"`
	Headers       []WebhookHeaderInput                 `json:"headers,omitempty"`
	TLSConfig     IntegrationTLSConfigInput            `json:"tlsConfig,omitempty"`
}

// UpdateSlackIntegrationParamsInput struct
type UpdateSlackIntegrationParamsInput struct {
	URL string `json:"url,omitempty"`
}

// UpdateSlackBotIntegrationParamsInput struct
type UpdateSlackBotIntegrationParamsInput struct {
	Token string `json:"token"`
}

// UpdateAzureServiceBusIntegrationParamsInput struct
type UpdateAzureServiceBusIntegrationParamsInput struct {
	QueueUrl     string                                      `json:"queueUrl,omitempty"`
	accessMethod AzureServiceBusIntegrationAccessMethodInput `json:"accessMethod,omitempty"`
}

// UpdateGcpPubSubIntegrationParamsInput struct
type UpdateGcpPubSubIntegrationParamsInput struct {
	ProjectID    string                                   `json:"projectId,omitempty"`
	TopicID      string                                   `json:"topicId,omitempty"`
	AccessMethod GooglePubSubIntegrationAccessMethodInput `json:"accessMethod,omitempty"`
}

// UpdatePagerDutyIntegrationParamsInput struct
type UpdatePagerDutyIntegrationParamsInput struct {
	IntegrationKey string `json:"integrationKey"`
}

// UpdateJiraIntegrationParamsInput struct
type UpdateJiraIntegrationParamsInput struct {
	ServerURL     string                            `json:"serverUrl,omitempty"`
	ServerType    sring                             `json:"serverType,omitempty"` // enum JiraServerType
	IsOnPrem      *bool                             `json:"isOnPrem,omitempty"`
	TLSConfig     IntegrationTLSConfigInput         `json:"tlsConfig,omitempty"`
	Authorization JiraIntegrationAuthorizationInput `json:"authorization,omitempty"`
}

// UpdateServiceNowIntegrationParamsInput struct
type UpdateServiceNowIntegrationParamsInput struct {
	URL           string                                  `json:"url,omitempty"`
	Authorization ServiceNowIntegrationAuthorizationInput `json:"authorization,omitempty"`
}

// UpdateOpsgenieIntegrationParamsInput struct
type UpdateOpsgenieIntegrationParamsInput struct {
	Key string `json:"key"`
}

// UpdateClickUpIntegrationParamsInput struct
type UpdateClickUpIntegrationParamsInput struct {
	Key string `json:"key"`
}
