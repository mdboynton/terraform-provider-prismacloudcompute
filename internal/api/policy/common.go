package policy

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
)

type CompliancePolicy struct {
    Id          string                      `json:"_id"`
	Rules       *[]CompliancePolicyRule `json:"rules"`
	PolicyType  string                      `json:"policyType"`
}

type CompliancePolicyRule struct {
	BlockMessage                    string                              `json:"blockMsg"`
    Collections                     []collection.Collection             `json:"collections" tfsdk:"collections"`
    Condition                       *CompliancePolicyRuleCondition  `json:"condition" tfsdk:"condition"`
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

type CompliancePolicyRuleCondition struct {
    Vulnerabilities []CompliancePolicyRuleVulnerability `json:"vulnerabilities" tfsdk:"vulnerabilities"`
}

type CompliancePolicyRuleVulnerability struct {
    Id          int         `json:"id" tfsdk:"id"`
    Block       bool        `json:"block" tfsdk:"block"`
}
