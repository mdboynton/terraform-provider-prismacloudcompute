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
    Id    string            `json:"id"`
	Rules *[]HostCompliancePolicyRule `json:"rules"`
	PolicyType  string           `json:"policyType"`
}

type HostCompliancePolicyRule struct {
    Order            int                    `json:"order"`
	Name             string                  `json:"name"`
    //Action          []string        `json:"action"`
    Collections      []collection.Collection `json:"collections" tfsdk:"collections"`
    //Modified        string              `json:"modified"`
    

    //Action          []string        `json:"action"`
    //AlertThreshold  HostCompliancePolicyRuleAlertThreshold `json:"alertThreshold"`
	ReportAllPassedAndFailedChecks bool                    `json:"allCompliance"`
	BlockMessage     string                  `json:"blockMsg"`
    //BlockThreshold  HostCompliancePolicyRuleBlockThreshold `json:"blockThreshold"`
	//Collections      []collection.Collection `json:"collections"`
	//Condition        HostCompliancePolicyRuleCondition    `json:"condition"`
    Condition        *HostCompliancePolicyRuleCondition    `json:"condition" tfsdk:"condition"`
    //CreatePR        bool            `json:"createPR"`
    //CVERules        []HostCompliancePolicyRuleCVERule      `json:"cveRules"`
	Disabled         bool                    `json:"disabled"`
	Effect           string                  `json:"effect"`
    //ExcludeBaseImageVulns       bool        `json:"excludeBaseImageVulns"`
    //GraceDays       int32             `json:"graceDays"`
    //GraceDaysPolicy HostCompliancePolicyRuleGraceDaysPolicy  `json:"graceDaysPolicy"`
    //Groups          []string            `json:"groups"`
    //License         HostCompliancePolicyRuleLicense     `json:"license"`
    //Modified        string              `json:"modified"`
	//Name             string                  `json:"name"`
	Notes            string                  `json:"notes"`
    //OnlyFixed       bool            `json:"onlyFixed"`
    Owner           string          `json:"owner"`
    //PkgTypesThresholds      []HostCompliancePolicyRulePkgTypesThresholds         `json:"pkgTypesThresholds"`
    //PreviousName        string      `json:"previousName"`
    //Principal           []string    `json:"principal"`
    //RiskFactorsEffects   []HostCompliancePolicyRuleRiskFactorsEffect     `json:"riskFactorsEffects"`
    //Tags                []HostCompliancePolicyRuleTag      `json:"tags"`
	Verbose          bool                    `json:"verbose"`
}

type HostCompliancePolicyRuleAlertThreshold struct {
    Disabled    bool    `json:"disabled"`
    Value       int32 `json:"value"`
}

type HostCompliancePolicyRuleBlockThreshold struct {
    Enabled     bool    `json:"enabled"`
    Value       int32 `json:"value"`
}

type HostCompliancePolicyRuleCondition struct {
    //Device      string      `json:"device"`
    //ReadOnly    bool        `json:"readonly"`
    Vulnerabilities []HostCompliancePolicyRuleVulnerability `json:"vulnerabilities" tfsdk:"vulnerabilities"`
}

type HostCompliancePolicyRuleVulnerability struct {
    Id          int         `json:"id" tfsdk:"id"`
    Block       bool        `json:"block" tfsdk:"block"`
}

type HostCompliancePolicyRuleCVERule struct {
    Description     string      `json:"description"`
    Effect          string      `json:"effect"`
    Expiration      HostCompliancePolicyRuleExpiration `json:"expiration"`
    Id              string      `json:"id"`
}

type HostCompliancePolicyRuleExpiration struct {
    Date            string      `json:"date"`
    Enabled         bool        `json:"enabled"`
}

type HostCompliancePolicyRuleGraceDaysPolicy struct {
    Enabled         bool        `json:"enabled"`
    Low             int         `json:"low"`
    Medium          int         `json:"medium"`
    High            int         `json:"high"`
    Critical        int         `json:"critical"`
}

type HostCompliancePolicyRuleLicense struct {
    AlertThreshold  HostCompliancePolicyRuleAlertThreshold  `json:"alertThreshold"`
    BlockThreshold  HostCompliancePolicyRuleBlockThreshold  `json:"blockThreshold"`
    Low             int         `json:"low"`
    Medium          int         `json:"medium"`
    High            int         `json:"high"`
    Critical        int         `json:"critical"`
}

type HostCompliancePolicyRulePkgTypesThresholds struct {
    AlertThreshold  HostCompliancePolicyRuleAlertThreshold  `json:"alertThreshold"`
    BlockThreshold  HostCompliancePolicyRuleBlockThreshold  `json:"blockThreshold"`
    Type            string      `json:"type"`
}

type HostCompliancePolicyRuleRiskFactorsEffect struct {
    Effect          string      `json:"effect"`
    RiskFactor      string      `json:"riskFactor"`
}

type HostCompliancePolicyRuleTag struct {
    Description     string      `json:"description"`
    Effect          string      `json:"effect"`
    Expiration      HostCompliancePolicyRuleExpiration  `json:"expiration"`
    Name            string      `json:"name"`
}

// Create new host compliance policy
//func CreateHostCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy HostCompliancePolicy) error {
func UpsertHostCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy HostCompliancePolicy) error {
    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    fmt.Println("UpsertHostCompliancePolicy request body:")
    fmt.Printf("%+v\n", policy)
    //fmt.Printf("%+v\n", *policy.Rules)
    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
	return c.Request(http.MethodPut, HostComplianceEndpoint, nil, policy, nil)
}

// Get host compliance policy
func GetHostCompliancePolicy(c api.PrismaCloudComputeAPIClient) (*HostCompliancePolicy, error) {
    var ans HostCompliancePolicy
    if err := c.Request(http.MethodGet, HostComplianceEndpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting host compliance policy: %s", err)
    }
    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    fmt.Println("GetHostCompliancePolicy response body:")
    fmt.Printf("%+v\n", ans)
    //fmt.Printf("%+v\n", *ans.Rules)
    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    return &ans, nil
}

//// Update an existing collection.
//func UpdateHostCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy HostCompliancePolicy) error {
//	//return c.Request(http.MethodPut, fmt.Sprintf("%s/%s", HostComplianceEndpoint, name), nil, policy, nil)
//	return c.Request(http.MethodPut, HostComplianceEndpoint, nil, policy, nil)
//}
//
//// Delete an existing collection.
//func DeleteHostCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy HostCompliancePolicy) error {
//	return c.Request(http.MethodPut, HostComplianceEndpoint, nil, policy, nil)
//}
