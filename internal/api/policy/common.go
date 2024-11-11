package policy

import (
	"fmt"
	"net/http"
    "sort"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
)

const (
	PolicyTypeAdmission                         = "admission"
	PolicyTypeComplianceCiImage                 = "ciImagesCompliance"
	PolicyTypeComplianceCiImageFormatted        = "CI images compliance"
	PolicyTypeComplianceContainer               = "containerCompliance"
	PolicyTypeComplianceContainerFormatted      = "container compliance"
	PolicyTypeComplianceHost                    = "hostCompliance"
	PolicyTypeComplianceHostFormatted           = "host compliance"
	PolicyTypeComplianceVmImage                 = "vmCompliance"
	PolicyTypeComplianceVmImageFormatted        = "VM compliance"
	PolicyTypeComplianceFunction                = "serverlessCompliance"
	PolicyTypeComplianceFunctionFormatted       = "serverless compliance"
	PolicyTypeComplianceCiFunction              = "ciServerlessCompliance"
	PolicyTypeComplianceCiFunctionFormatted     = "CI serverless compliance"
    PolicyTypeComplianceTrustedImages           = "trust"
	PolicyTypeRuntimeContainer                  = "containerRuntime"
	PolicyTypeRuntimeHost                       = "hostRuntime"
	PolicyTypeVulnerabilityCiImage              = "ciImagesVulnerability"
	PolicyTypeVulnerabilityHost                 = "hostVulnerability"
	PolicyTypeVulnerabilityDeployedImage        = "containerVulnerability"
	PolicyTypeVulnerabilityDeployedImageFormatted = "deployed image vulnerability"
	PolicyContextVulnerabilityDeployedImage     = "images"
	TypeVulnerability                           = "vulnerability"
)

type CompliancePolicy struct {
    Id          string                      `json:"_id"`
	Rules       *[]CompliancePolicyRule `json:"rules"`
	PolicyType  string                      `json:"policyType"`
}

func (p *CompliancePolicy) SortRules(orderMap map[string]int) {
    sort.Slice(*p.Rules, func(i, j int) bool {
        return orderMap[(*p.Rules)[i].Name] < orderMap[(*p.Rules)[j].Name]
    })
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

// TODO: rename to PolicyRuleCondition
type CompliancePolicyRuleCondition struct {
    Vulnerabilities []CompliancePolicyRuleVulnerability `json:"vulnerabilities" tfsdk:"vulnerabilities"`
}

// TODO: rename to PolicyRuleVulnerability
type CompliancePolicyRuleVulnerability struct {
    Id          int         `json:"id" tfsdk:"id"`
    Block       bool        `json:"block" tfsdk:"block"`
}

type VulnerabilityPolicy struct {
    Id              string                          `json:"_id"`
	PolicyType      string                          `json:"policyType"`
	PolicyContext   string                          `json:"policyContext"`
	Rules           *[]VulnerabilityPolicyRule    `json:"rules"`
    Type            string                          `json:"type"`
}

type VulnerabilityPolicyRule struct {
    AlertThreshold                  AlertThreshold                  `json:"alertThreshold" tfsdk:"alert_threshold"`
	BlockMessage                    string                              `json:"blockMsg"`
    BlockThreshold                  BlockThreshold                  `json:"blockThreshold" tfsdk:"block_threshold"`
    Collections                     []collection.Collection             `json:"collections" tfsdk:"collections"`
    //Condition                       *CompliancePolicyRuleCondition  `json:"condition" tfsdk:"condition"`
    CVERules                        []Exception                         `json:"cveRules"`
	Disabled                        bool                                `json:"disabled"`
	Effect                          string                              `json:"effect"`
    ExcludeBaseImageVulns           bool                                `json:"excludeBaseImageVulns"`
    GraceDays                       int                                 `json:"graceDays"`
    GraceDaysPolicy                 GraceDaysPolicy                     `json:"graceDaysPolicy"`
    Modified                        string                              `json:"modified"`
	Name                            string                              `json:"name"`
	Notes                           string                              `json:"notes"`
    OnlyFixed                       bool                                `json:"onlyFixed"`
    Order                           int                                 `json:"order"`
    Owner                           string                              `json:"owner"`
    PkgTypesThresholds              []PkgTypesThreshold                 `json:"pkgTypesThresholds"`
	ReportAllPassedAndFailedChecks  bool                                `json:"allCompliance"`
    RiskFactorEffects               []RiskFactorsEffect                 `json:"riskFactorsEffects"` 
    Tags                            []Exception                         `json:"tags"`
	Verbose                         bool                                `json:"verbose"`
}

type RiskFactorsEffect struct {
    Effect          string      `json:"effect"`
    RiskFactor      string      `json:"riskFactor"`
}

type AlertThreshold struct {
    Disabled    bool    `json:"disabled" tfsdk:"disabled"`
    Value       int     `json:"value" tfsdk:"value"`
}

type BlockThreshold struct {
    Enabled     bool    `json:"enabled" tfsdk:"enabled"`
    Value       int     `json:"value" tfsdk:"value"`
}

type GraceDaysPolicy struct {
    Enabled     bool    `json:"enabled"`
    Low         int     `json:"low"`
    Medium      int     `json:"medium"`
    High        int     `json:"high"`
    Critical    int     `json:"critical"`
}

type PkgTypesThreshold struct {
    Type        string  `json:"type"`
    AlertThreshold  AlertThreshold  `json:"alertThreshold"`
    BlockThreshold  BlockThreshold  `json:"blockThreshold"`
}

type Exception struct {
    Name        string              `json:"name"`
    Effect      string              `json:"effect"`
    Id          string              `json:"id"`
    Description string              `json:"description"`
    Type        string              `json:"type"`
    Expiration  ExceptionExpiration `json:"expiration"`
}

type ExceptionExpiration struct {
    Enabled bool    `json:"enabled"`
    Date    string  `json:"date"`
}

func getEndpointAndPolicyName(policyType string) (string, string, error) {
    switch policyType{
        case "hostCompliance":
            return HostComplianceEndpoint, "host", nil
        case "containerCompliance":
            return ContainerComplianceEndpoint, "container", nil
        case "ciImagesCompliance":
            return CiImageComplianceEndpoint, "CI image", nil
        case "vmCompliance":
            return VmImageComplianceEndpoint, "VM image", nil
        case "serverlessCompliance":
            return FunctionComplianceEndpoint, "function", nil
        case "ciServerlessCompliance":
            return CiFunctionComplianceEndpoint, "CI function", nil
        case "containerVulnerability":
            return DeployedImagesVulnerabilityEndpoint, "deployed images", nil
        default:
            return "", "", fmt.Errorf("invalid policy type specified")
    }
}

// Create/Update/Delete compliance policy
func UpsertCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy CompliancePolicy) error {
    endpoint, policyName, err := getEndpointAndPolicyName(policy.PolicyType)
    if err != nil {
		return err
    }

    if err := c.Request(http.MethodPut, endpoint, nil, policy, nil); err != nil {
		return fmt.Errorf("error upserting %s compliance policy: %s", policyName, err)
    }
    
    return nil
}

// Create/Update/Delete vulnerability policy
func UpsertVulnerabilityPolicy(c api.PrismaCloudComputeAPIClient, policy VulnerabilityPolicy) error {
    endpoint, policyName, err := getEndpointAndPolicyName(policy.PolicyType)
    if err != nil {
		return err
    }

    if err := c.Request(http.MethodPut, endpoint, nil, policy, nil); err != nil {
		return fmt.Errorf("error upserting %s vulnerability policy: %s", policyName, err)
    }
    
    return nil
}

// Get compliance policy
func GetCompliancePolicy(c api.PrismaCloudComputeAPIClient, policyType string) (*CompliancePolicy, error) {
    var ans CompliancePolicy

    endpoint, policyName, err := getEndpointAndPolicyName(policyType)
    if err != nil {
		return &ans, err
    }
    
    if err := c.Request(http.MethodGet, endpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error retrieving %s compliance policy: %s", policyName, err)
    }
    
    return &ans, nil
}

// Get vulnerability policy
func GetVulnerabilityPolicy(c api.PrismaCloudComputeAPIClient, policyType string) (*VulnerabilityPolicy, error) {
    var ans VulnerabilityPolicy 

    endpoint, policyName, err := getEndpointAndPolicyName(policyType)
    if err != nil {
		return &ans, err
    }
    
    if err := c.Request(http.MethodGet, endpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error retrieving %s vulnerability policy: %s", policyName, err)
    }
    
    return &ans, nil
}
