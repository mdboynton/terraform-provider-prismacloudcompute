package system

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/auth"
)

const CloudScanRulesEndpoint = "api/v1/cloud-scan-rules"

type CloudScanRule struct {
	CredentialId                string             `json:"credentialId"`
	Credential                  auth.Credential    `json:"credential,omitempty"`
	DiscoveryEnabled            bool               `json:"discoveryEnabled,omitempty"`
    //ServerlessRadarEnabled      bool               `json:"serverlessRadarEnabled,omitempty"`
	//VmTagsEnabled               bool               `json:"vmTagsEnabled,omitempty"`
	//DiscoverAllFunctionVersions bool               `json:"discoverAllFunctionVersions,omitempty"`
	//ServerlessRadarCap          int                `json:"serverlessRadarCap,omitempty"`
	AgentlessScanSpec           AgentlessScanSpec  `json:"agentlessScanSpec,omitempty"`
	ServerlessScanSpec          ServerlessScanSpec `json:"serverlessScanSpec,omitempty"`
	AwsRegionType               string             `json:"awsRegionType,omitempty"`
}

type ServerlessScanSpec struct {
	Enabled         bool `json:"enabled,omitempty"`
	Cap             int  `json:"cap,omitempty"`
	ScanAllVersions bool `json:"scanAllVersions,omitempty"`
	ScanLayers      bool `json:"scanLayers,omitempty"`
}

type AgentlessScanSpec struct {
	Enabled              bool     `json:"enabled,omitempty"`
	HubAccount           bool     `json:"hubAccount,omitempty"`
    HubCredentialID      string   `json:"hubCredentialID"`
	ConsoleAddress       string   `json:"consoleAddr,omitempty"`
	ScanNonRunning       bool     `json:"scanNonRunning,omitempty"`
	ProxyAddress         string   `json:"proxyAddress,omitempty"`
	ProxyCA              string   `json:"proxyCA,omitempty"`
	SkipPermissionsCheck bool     `json:"skipPermissionsCheck,omitempty"`
	AutoScale            bool     `json:"autoScale,omitempty"`
	Scanners             int      `json:"scanners,omitempty"`
	SecurityGroup        string   `json:"securityGroup,omitempty"`
	Subnet               string   `json:"subnet,omitempty"`
	Regions              []string `json:"regions,omitempty"`
	CustomTags           []Tag    `json:"customTags,omitempty"`
	IncludedTags         []Tag    `json:"includedTags,omitempty"`
}

type Tag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func ListCloudScanRules(c api.PrismaCloudComputeAPIClient) ([]CloudScanRule, error) {
	var ans []CloudScanRule
	if err := c.Request(http.MethodGet, CloudScanRulesEndpoint, nil, nil, &ans); err != nil {
		return nil, fmt.Errorf("error listing Cloud Scan Rules: %s", err)
	}
	return ans, nil
}

func GetCloudScanRuleByCredentialId(c api.PrismaCloudComputeAPIClient, credentialId string) (*CloudScanRule, error) {
	var ans []CloudScanRule

	if err := c.Request(http.MethodGet, CloudScanRulesEndpoint, map[string]string{"search": credentialId}, nil, &ans); err != nil {
		return nil, fmt.Errorf("Error retrieving cloud scan rule: %s", err)
	}

	for _, val := range ans {
		if val.CredentialId == credentialId {
			return &val, nil
		}
	}
	return nil, fmt.Errorf("cloud scan rule with credential ID \"%s\" not found", credentialId)
}

func UpdateCloudScanRule(c api.PrismaCloudComputeAPIClient, rule []CloudScanRule) error {
	return c.Request(http.MethodPut, CloudScanRulesEndpoint, nil, rule, nil)
}

func DeleteCloudScanRule(c api.PrismaCloudComputeAPIClient, name string) error {
	return c.Request(http.MethodDelete, fmt.Sprintf("%s/%s", CloudScanRulesEndpoint, name), nil, nil, nil)
}
