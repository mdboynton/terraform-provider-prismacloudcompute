package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
)

const (
	HostComplianceEndpoint         = "api/v1/policies/compliance/host"
)

type HostCompliancePolicy struct {
    Id          string                      `json:"_id"`
	Rules       *[]HostCompliancePolicyRule `json:"rules"`
	PolicyType  string                      `json:"policyType"`
}

type HostCompliancePolicyRule struct {
	BlockMessage                    string                              `json:"blockMsg"`
    Collections                     []collection.Collection             `json:"collections" tfsdk:"collections"`
    Condition                       *HostCompliancePolicyRuleCondition  `json:"condition" tfsdk:"condition"`
	Disabled                        bool                                `json:"disabled"`
	Effect                          string                              `json:"effect"`
    Modified                        string                              `json:"modified"`
	Name                            string                              `json:"name"`
	Notes                           string                              `json:"notes"`
    Order                           int                                 `json:"order"`
    Owner                           string                              `json:"owner"`
	ReportAllPassedAndFailedChecks  bool                                `json:"allCompliance"`
	Verbose                         bool                                `json:"verbose"`
}

type HostCompliancePolicyRuleCondition struct {
    Vulnerabilities []HostCompliancePolicyRuleVulnerability `json:"vulnerabilities" tfsdk:"vulnerabilities"`
}

type HostCompliancePolicyRuleVulnerability struct {
    Id          int         `json:"id" tfsdk:"id"`
    Block       bool        `json:"block" tfsdk:"block"`
}

// Create/Update/Delete host compliance policy
func UpsertHostCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy HostCompliancePolicy) error {
	return c.Request(http.MethodPut, HostComplianceEndpoint, nil, policy, nil)
}

// Get host compliance policy
func GetHostCompliancePolicy(c api.PrismaCloudComputeAPIClient) (*HostCompliancePolicy, error) {
    var ans HostCompliancePolicy
    if err := c.Request(http.MethodGet, HostComplianceEndpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting host compliance policy: %s", err)
    }
    
    return &ans, nil
}
