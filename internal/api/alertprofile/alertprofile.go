package alertprofile

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const AlertprofilesEndpoint = "api/v1/alert-profiles"

// type Email struct {
// 	Enabled      bool   `json:"enabled"`
// 	SmtpAddress  string `json:"smtpAddress"`
// 	Port         int    `json:"port"`
// 	CredentialId string `json:"credentialId"`
// 	From         string `json:"from"`
// 	Ssl          string `json:"ssl"`
// }

// type Slack struct {
// 	Enabled    bool   `json:"enabled"`
// 	WebhookUrl string `json:"webhookUrl"`
// }

// type Jira struct {
// 	Enabled      bool   `json:"enabled"`
// 	BaseUrl      string `json:"baseUrl"`
// 	CredentialId int    `json:"credentialId"`
// 	CaCert       string `json:"caCert"`
// 	ProjectKey   string `json:"projectKey"`
// 	IssueType    string `json:"issueType"`
// 	Priority     string `json:"priority"`
// 	Labels       string `json:"labels"`
// 	Assignee     string `json:"assignee"`
// }

// type SecurityCenter struct {
// 	Enabled      bool   `json:"enabled"`
// 	CredentialId int    `json:"credentialId"`
// 	SourceID     string `json:"sourceID"`
// }

// type GcpPubsub struct {
// 	Enabled      bool   `json:"enabled"`
// 	CredentialId int    `json:"credentialId"`
// 	Topic        string `json:"topic"`
// }

// type SecurityHub struct {
// 	Enabled      bool   `json:"enabled"`
// 	CredentialId string `json:"credentialId"`
// 	Region       string `json:"region"`
// 	AccountID    string `json:"accountID"`
// }

type PolicyModule struct {
	Enabled  bool     `json:"enabled" tfsdk:"enabled"`
	AllRules bool     `json:"allRules" tfsdk:"all_rules"`
	Rules    []string `json:"rules,omitempty" tfsdk:"rules"`
}

// Policy struct
type Policy struct {
	Admission               PolicyModule `json:"admission" tfsdk:"admission"`
	AgentlessAppFirewall    PolicyModule `json:"agentlessAppFirewall" tfsdk:"agentless_app_firewall"`
	AppEmbeddedAppFirewall  PolicyModule `json:"appEmbeddedAppFirewall" tfsdk:"app_embedded_app_firewall"`
	AppEmbeddedRuntime      PolicyModule `json:"appEmbeddedRuntime" tfsdk:"app_embedded_runtime"`
	CloudDiscovery          PolicyModule `json:"cloudDiscovery" tfsdk:"cloud_discovery"`
	ContainerAppFirewall    PolicyModule `json:"containerAppFirewall" tfsdk:"container_app_firewall"`
	ContainerCompliance     PolicyModule `json:"containerCompliance" tfsdk:"container_compliance"`
	ContainerComplianceScan PolicyModule `json:"containerComplianceScan" tfsdk:"container_compliance_scan"`
	ContainerRuntime        PolicyModule `json:"containerRuntime" tfsdk:"container_runtime"`
	ContainerVulnerability  PolicyModule `json:"containerVulnerability" tfsdk:"container_vulnerability"`
	Defender                PolicyModule `json:"defender" tfsdk:"defender"`
	HostAppFirewall         PolicyModule `json:"hostAppFirewall" tfsdk:"host_app_firewall"`
	HostCompliance          PolicyModule `json:"hostCompliance" tfsdk:"host_compliance"`
	HostComplianceScan      PolicyModule `json:"hostComplianceScan" tfsdk:"host_compliance_scan"`
	HostRuntime             PolicyModule `json:"hostRuntime" tfsdk:"host_runtime"`
	HostVulnerability       PolicyModule `json:"hostVulnerability" tfsdk:"host_vulnerability"`
	Incident                PolicyModule `json:"incident" tfsdk:"incident"`
	KubernetesAudit         PolicyModule `json:"kubernetesAudit" tfsdk:"kubernetes_audit"`
	NetworkFirewall         PolicyModule `json:"networkFirewall" tfsdk:"network_firewall"`
	RegistryVulnerability   PolicyModule `json:"registryVulnerability" tfsdk:"registry_vulnerability"`
	ServerlessAppFirewall   PolicyModule `json:"serverlessAppFirewall" tfsdk:"serverless_app_firewall"`
	ServerlessRuntime       PolicyModule `json:"serverlessRuntime" tfsdk:"serverless_runtime"`
	VmCompliance            PolicyModule `json:"vmCompliance" tfsdk:"vm_compliance"`
	VmVulnerability         PolicyModule `json:"vmVulnerability" tfsdk:"vm_vulnerability"`
	WaasHealth              PolicyModule `json:"waasHealth" tfsdk:"waas_health"`
}

type Webhook struct {
	Enabled bool `json:"enabled" tfsdk:"enabled"`
	//CredentialId string `json:"credentialId,omitempty"`
	Url    string `json:"url,omitempty" tfsdk:"url"`
	CaCert string `json:"caCert,omitempty" tfsdk:"ca_cert"`
	Json   string `json:"json,omitempty" tfsdk:"json"`
}

type AlertProfile struct {
	Id                                  string  `json:"_id"`
	Name                                string  `json:"name"`
	VulnerabilityImmediateAlertsEnabled bool    `json:"vulnerabilityImmediateAlertsEnabled,omitempty"`
	Owner                               string  `json:"owner,omitempty"`
	Webhook                             Webhook `json:"webhook,omitempty"`
	Policy                              Policy  `json:"policy,omitempty"`
}

// Retrieve all alert profiles
func ListAlertprofiles(c api.PrismaCloudComputeAPIClient) ([]AlertProfile, error) {
	var ans []AlertProfile
	if err := c.Request(http.MethodGet, AlertprofilesEndpoint, nil, nil, &ans); err != nil {
		return nil, fmt.Errorf("error listing Alert Profiles: %s", err)
	}
	return ans, nil
}

// Retrieve an alert profile with the specified name
func GetAlertprofile(c api.PrismaCloudComputeAPIClient, name string) (*AlertProfile, error) {
	Alertprofiles, err := ListAlertprofiles(c)
	if err != nil {
		return nil, err
	}
	for _, val := range Alertprofiles {
		if val.Name == name {
			return &val, nil
		}
	}
	return nil, fmt.Errorf("error Alert Profile '%s' not found", name)
}

// Create or update an alert profile
func CreateOrUpdateAlertprofile(c api.PrismaCloudComputeAPIClient, Alertprofile AlertProfile) error {
	return c.Request(http.MethodPost, AlertprofilesEndpoint, nil, Alertprofile, nil)
}

// Delete an existing alert profile
func DeleteAlertprofile(c api.PrismaCloudComputeAPIClient, name string) error {
	return c.Request(http.MethodDelete, fmt.Sprintf("%s/%s", AlertprofilesEndpoint, name), nil, nil, nil)
}
