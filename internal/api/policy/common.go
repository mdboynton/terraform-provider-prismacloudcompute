package policy

import (
	"fmt"
    "context"
    "sort"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
)

const (
    // Policy Type
	PolicyTypeAdmission                         = "admission"
	PolicyTypeComplianceCiImage                 = "ciImagesCompliance"
	PolicyTypeComplianceContainer               = "containerCompliance"
	PolicyTypeComplianceHost                    = "hostCompliance"
	PolicyTypeComplianceVmImage                 = "vmCompliance"
	PolicyTypeComplianceFunction                = "serverlessCompliance"
	PolicyTypeComplianceCiFunction              = "ciServerlessCompliance"
    PolicyTypeComplianceTrustedImages           = "trust"
    // Policy Type Formatted
	PolicyTypeComplianceCiImageFormatted        = "CI images compliance"
	PolicyTypeComplianceContainerFormatted      = "container compliance"
	PolicyTypeComplianceHostFormatted           = "host compliance"
	PolicyTypeComplianceVmImageFormatted        = "VM compliance"
	PolicyTypeComplianceFunctionFormatted       = "serverless compliance"
	PolicyTypeComplianceCiFunctionFormatted     = "CI serverless compliance"
	PolicyTypeRuntimeContainer                  = "containerRuntime"
	PolicyTypeRuntimeHost                       = "hostRuntime"
	PolicyTypeVulnerabilityCiImage              = "ciImagesVulnerability"
	PolicyTypeVulnerabilityHost                 = "hostVulnerability"
	PolicyTypeVulnerabilityDeployedImage        = "containerVulnerability"
	PolicyTypeVulnerabilityDeployedImageFormatted = "deployed image vulnerability"
    // Policy Context
	PolicyContextComplianceContainer              = "container"
	PolicyContextComplianceCiFunction              = "ciServerless"
	PolicyContextComplianceCiImage              = "ciImages"
	PolicyContextComplianceFunction              = "serverless"
	PolicyContextComplianceHost              = "host"
	PolicyContextComplianceVmImage              = "vms"
	PolicyContextVulnerabilityDeployedImage     = "images"
    // Type
	TypeCompliance                              = "compliance"
	TypeVulnerability                           = "vulnerability"
)

type Policy struct {
    Id              string                          `json:"_id"`
	PolicyType      string                          `json:"policyType"`
	PolicyContext   string                          `json:"policyContext"`
	Rules           *[]PolicyRule                   `json:"rules"`
    Type            string                          `json:"type"`
}

func (p *Policy) SortRules(ctx context.Context, planRules *[]models.PolicyRuleResourceModel) {
    util.DLog(ctx, "Executing api.Policy.SortRules()")

    rulesOrderMap := generatePolicyRulesOrderMap(*planRules)
    sort.Slice((*p.Rules), func(i, j int) bool {
        return rulesOrderMap[(*p.Rules)[i].Name] < rulesOrderMap[(*p.Rules)[j].Name]
    })
    
    util.DLog(ctx, "Finishing api.Policy.SortRules() execution")
}

// TODO: remove this duplicate function when we can move the logic somewhere that can be used here
// and by internal/resources/policy/common.go
func generatePolicyRulesOrderMap(rules []models.PolicyRuleResourceModel) map[string]int {
    orderedRulesMap := make(map[int][]string)

    for _, rule := range rules {
        order := int(rule.Order.ValueInt32())
        if _, exists := orderedRulesMap[order]; exists {
            orderedRulesMap[order] = append(orderedRulesMap[order], rule.Name.ValueString())
        } else {
            orderedRulesMap[order] = []string{rule.Name.ValueString()}
        }
    }
        
    sortedKeys := make([]int, 0, len(orderedRulesMap))
    for key := range orderedRulesMap {
        sortedKeys = append(sortedKeys, key)
    }
    sort.Ints(sortedKeys)

    ruleOrders := make(map[string]int)
    lastOrderValue := -1
    for _, key := range sortedKeys {
        offset := 0
        if lastOrderValue != -1 && lastOrderValue >= key {
            offset = lastOrderValue - key + 1
        }

        for sliceIndex, ruleName := range orderedRulesMap[key] {
            orderValue := key + sliceIndex + offset
            ruleOrders[ruleName] = orderValue
            lastOrderValue = orderValue
        }
    }

    return ruleOrders
}

type PolicyRule struct {
    AlertThreshold                  AlertThreshold                  `json:"alertThreshold" tfsdk:"alert_threshold"`
	BlockMessage                    string                              `json:"blockMsg"`
    BlockThreshold                  BlockThreshold                  `json:"blockThreshold" tfsdk:"block_threshold"`
    Collections                     []collection.Collection             `json:"collections" tfsdk:"collections"`
    Condition                       *CompliancePolicyRuleCondition  `json:"condition" tfsdk:"condition"`
    CVERules                        []Exception                         `json:"cveRules" tfsdk:"cve_rules"`
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


func UpsertPolicy(c api.PrismaCloudComputeAPIClient, policy Policy) error {
    endpoint, policyName, err := getEndpointAndPolicyName(policy.PolicyType)
    if err != nil {
		return err
    }

    if err := c.Request(http.MethodPut, endpoint, nil, policy, nil); err != nil {
		return fmt.Errorf("error upserting %s policy: %s", policyName, err)
    }
    
    return nil
}


func GetPolicy(c api.PrismaCloudComputeAPIClient, policyType string) (*Policy, error) {
    var ans Policy 

    endpoint, policyName, err := getEndpointAndPolicyName(policyType)
    if err != nil {
		return &ans, err
    }
    
    if err := c.Request(http.MethodGet, endpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error retrieving %s policy: %s", policyName, err)
    }
    
    return &ans, nil
}
