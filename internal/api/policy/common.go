package policy

import (
	"context"
	"fmt"
    "errors"
	"net/http"
	"sort"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
)

type Settings struct {
    FormattedName string
    Module string 
    Endpoint string
    Context string
    IsApplicableVuln func(v system.Vulnerability) bool
}

const BaseEndpoint = "api/v1"

// TODO: review all policies to make sure there's no policy types missing from these IsApplicableVuln functions
var SettingsMap map[string]Settings = map[string]Settings {
        PolicyTypeComplianceContainer: {
            FormattedName: "Deployed Containers and Images Compliance",
            Module: "compliance",
            Context: "container",
            Endpoint: fmt.Sprintf("%s/policies/compliance/container", BaseEndpoint),
            IsApplicableVuln: func(v system.Vulnerability) bool {
                if v.Type == "container" {
                    return true
                }

                if v.Type == "image" {
                    return true
                }

                if v.Type == "istio" {
                    return true
                }

                return false
            },
        },
        PolicyTypeComplianceCiImage: {
            FormattedName: "CI Images Compliance",
            Module: "compliance",
            Context: "ciImages",
            Endpoint: fmt.Sprintf("%s/policies/compliance/ci/container", BaseEndpoint),
            IsApplicableVuln: func(v system.Vulnerability) bool {
                if v.Type == "image" && (v.Id == 406 || v.Id == 408 || v.Id == 41 || v.Id == 422 || v.Id == 424 || v.Id == 425 || v.Id == 426 || v.Id == 448 || v.Id == 5041 || v.Id == 5048) {
                    return true
                }

                return false
            },
        },
        PolicyTypeComplianceHost: {
            FormattedName: "Host Compliance",
            Module: "compliance",
            Context: "host",
            Endpoint: fmt.Sprintf("%s/policies/compliance/host", BaseEndpoint),
            IsApplicableVuln: func(v system.Vulnerability) bool {
                if v.Type == "host_config" {
                    return true
                }

                if v.Type == "daemon_config" {
                    return true
                }

                if v.Type == "daemon_config_files" {
                    return true
                }

                if v.Type == "security_operations" {
                    return true
                }

                if v.Type == "linux" {
                    return true
                }

                if v.Type == "windows" {
                    return true
                }

                if v.Type == "k8s_worker" {
                    return true
                }
                
                if v.Type == "gke_worker" {
                    return true
                }

                if v.Type == "eks_worker" {
                    return true
                }

                if v.Type == "aks_worker" {
                    return true
                }

                if v.Type == "openshift_worker" {
                    return true
                }

                if v.Type == "k8s_master" {
                    return true
                }

                if v.Type == "openshift_master" {
                    return true
                }

                if v.Type == "k8s_federation" {
                    return true
                }
                
                if v.Type == "docker_stig" {
                    return true
                }

                return false
            },
        },
        PolicyTypeComplianceVmImage: {
            FormattedName: "VM Image Compliance",
            Module: "compliance",
            Context: "vms",
            Endpoint: fmt.Sprintf("%s/policies/compliance/vms", BaseEndpoint),
            IsApplicableVuln: func(v system.Vulnerability) bool {
                if v.Type == "host" {
                    return true
                }

                if v.Type == "host_config" {
                    return true
                }

                if v.Type == "daemon_config" {
                    return true
                }

                if v.Type == "daemon_config_files" {
                    return true
                }

                if v.Type == "security_operations" {
                    return true
                }

                if v.Type == "linux" {
                    return true
                }

                return false
            },
        },
        PolicyTypeComplianceFunction: {
            FormattedName: "Serverless Compliance",
            Module: "compliance",
            Context: "serverless",
            Endpoint: fmt.Sprintf("%s/policies/compliance/serverless", BaseEndpoint),
            IsApplicableVuln: func(v system.Vulnerability) bool {
                if v.Type == "serverless" {
                    return true
                }

                return false
            },
        },
        PolicyTypeComplianceCiFunction: {
            FormattedName: "CI Serverless Compliance",
            Module: "compliance",
            Context: "ciServerless",
            Endpoint: fmt.Sprintf("%s/policies/compliance/ci/serverless", BaseEndpoint),
            IsApplicableVuln: func(v system.Vulnerability) bool {
                if v.Type == "serverless" {
                    return true
                }

                return false
            },
        },
   }

const (
	// Policy Type
	PolicyTypeAdmission               = "admission"
	PolicyTypeComplianceCiImage       = "ciImagesCompliance"
	PolicyTypeComplianceContainer     = "containerCompliance"
	PolicyTypeComplianceHost          = "hostCompliance"
	PolicyTypeComplianceVmImage       = "vmCompliance"
	PolicyTypeComplianceFunction      = "serverlessCompliance"
	PolicyTypeComplianceCiFunction    = "ciServerlessCompliance"
	PolicyTypeComplianceTrustedImages = "trust"
	PolicyTypeVulnerabilityDeployedImage          = "containerVulnerability"
	PolicyTypeVulnerabilityCiImage                = "ciImagesVulnerability"
	PolicyTypeVulnerabilityHost                   = "hostVulnerability"
	PolicyTypeVulnerabilityVmImage                = "vmVulnerability"
	PolicyTypeVulnerabilityFunction               = "serverlessVulnerability"
	PolicyTypeVulnerabilityCiFunction             = "ciServerlessVulnerability"
	PolicyTypeRuntimeContainer                    = "containerRuntime"
	PolicyTypeRuntimeHost                         = "hostRuntime"
	PolicyTypeRuntimeServerless                   = "serverlessRuntime"
	PolicyTypeRuntimeAppEmbedded                  = "appEmbeddedRuntime"
	// Policy Type Formatted
	PolicyTypeComplianceCiImageFormatted          = "CI images compliance"
	PolicyTypeComplianceContainerFormatted        = "container compliance"
	PolicyTypeComplianceHostFormatted             = "host compliance"
	PolicyTypeComplianceVmImageFormatted          = "VM compliance"
	PolicyTypeComplianceFunctionFormatted         = "serverless compliance"
	PolicyTypeComplianceCiFunctionFormatted       = "CI serverless compliance"
	PolicyTypeVulnerabilityDeployedImageFormatted = "deployed image vulnerability"
	PolicyTypeVulnerabilityCiImageFormatted       = "CI image vulnerability"
	PolicyTypeVulnerabilityHostFormatted          = "host vulnerability"
	PolicyTypeVulnerabilityVmImageFormatted       = "VM vulnerability"
	PolicyTypeVulnerabilityFunctionFormatted      = "serverless vulnerability"
	PolicyTypeVulnerabilityCiFunctionFormatted    = "CI serverless vulnerability"
	PolicyTypeRuntimeContainerFormatted           = "container runtime"
	PolicyTypeRuntimeHostFormatted                = "host runtime"
	PolicyTypeRuntimeServerlessFormatted          = "serverless runtime"
	PolicyTypeRuntimeAppEmbeddedFormatted         = "app-embedded runtime"
	// Policy Context
	PolicyContextImage            = "images"
	PolicyContextContainer        = "container"
	PolicyContextCiFunction       = "ciServerless"
	PolicyContextCiImage          = "ciImages"
	PolicyContextFunction         = "serverless"
	PolicyContextHost             = "host"
	PolicyContextVmImage          = "vms"
	// Type
	TypeCompliance    = "compliance"
	TypeVulnerability = "vulnerability"

    //BaseEndpoint                = "api/v1"
)

var (
    // Compliance endpoints
	ContainerComplianceEndpoint = fmt.Sprintf("%s/policies/compliance/container", BaseEndpoint)
	CiImageComplianceEndpoint   = fmt.Sprintf("%s/policies/compliance/ci/images", BaseEndpoint)
	HostComplianceEndpoint      = fmt.Sprintf("%s/policies/compliance/host", BaseEndpoint)
	VmImageComplianceEndpoint   = fmt.Sprintf("%s/policies/compliance/vms", BaseEndpoint)
    ApplicationControlEndpoint = fmt.Sprintf("%s/application-control/host", BaseEndpoint)
    FunctionComplianceEndpoint  = fmt.Sprintf("%s/policies/compliance/serverless", BaseEndpoint)
    CiFunctionComplianceEndpoint  = fmt.Sprintf("%s/policies/compliance/ci/serverless", BaseEndpoint)
    TrustedImagesEndpoint = fmt.Sprintf("%s/trust/data", BaseEndpoint)
    CustomComplianceChecksEndpoint = fmt.Sprintf("%s/custom-compliance", BaseEndpoint)
    // Vulnerabilities endpoints
	DeployedImageVulnerabilityEndpoint      = fmt.Sprintf("%s/policies/policies/vulnerability/images", BaseEndpoint)
	CiImageVulnerabilityEndpoint            = fmt.Sprintf("%s/policies/policies/vulnerability/ci/images", BaseEndpoint)
	HostVulnerabilityEndpoint               = fmt.Sprintf("%s/policies/vulnerability/host", BaseEndpoint)
	VmImageVulnerabilityEndpoint            = fmt.Sprintf("%s/policies/vulnerability/vms", BaseEndpoint)
	FunctionVulnerabilityEndpoint           = fmt.Sprintf("%s/policies/vulnerability/serverless", BaseEndpoint)
	CiFunctionVulnerabilityEndpoint         = fmt.Sprintf("%s/policies/vulnerability/ci/serverless", BaseEndpoint)
    // Runtime endpoints
    RuntimeContainerEndpoint = fmt.Sprintf("%s/policies/runtime/container", BaseEndpoint)
    RuntimeHostEndpoint = fmt.Sprintf("%s/policies/runtime/host", BaseEndpoint)
    RuntimeServerlessEndpoint = fmt.Sprintf("%s/policies/runtime/serverless", BaseEndpoint)
    RuntimeAppEmbeddedEndpoint = fmt.Sprintf("%s/policies/runtime/app-embedded", BaseEndpoint)

    IsAttributeSupported = map[string]map[string]bool{
        "block_message": {
            PolicyTypeComplianceContainer: true, 
            PolicyTypeComplianceCiImage: false, 
            PolicyTypeComplianceHost: true, 
            PolicyTypeComplianceVmImage: true,
            PolicyTypeComplianceFunction: false,
            PolicyTypeComplianceCiFunction: false,
            PolicyTypeVulnerabilityDeployedImage: true,
            PolicyTypeVulnerabilityCiImage: false,
            PolicyTypeVulnerabilityHost: false,
            PolicyTypeVulnerabilityVmImage: false,
            PolicyTypeVulnerabilityFunction: false,
            PolicyTypeVulnerabilityCiFunction: false,
        },
        "report_passed_and_failed_checks": {
            PolicyTypeComplianceContainer: true, 
            PolicyTypeComplianceCiImage: false, 
            PolicyTypeComplianceHost: true, 
            PolicyTypeComplianceVmImage: true,
            PolicyTypeComplianceFunction: false,
            PolicyTypeComplianceCiFunction: false,
            PolicyTypeVulnerabilityDeployedImage: false,
            PolicyTypeVulnerabilityCiImage: false,
            PolicyTypeVulnerabilityHost: false,
            PolicyTypeVulnerabilityVmImage: false,
            PolicyTypeVulnerabilityFunction: false,
            PolicyTypeVulnerabilityCiFunction: false,
        },
        "exclude_base_image_vulns": {
            PolicyTypeComplianceContainer: true, 
            PolicyTypeComplianceCiImage: false, 
            PolicyTypeComplianceHost: true, 
            PolicyTypeComplianceVmImage: true,
            PolicyTypeComplianceFunction: false,
            PolicyTypeComplianceCiFunction: false,
            PolicyTypeVulnerabilityDeployedImage: true,
            PolicyTypeVulnerabilityCiImage: true,
            PolicyTypeVulnerabilityHost: false,
            PolicyTypeVulnerabilityVmImage: false,
            PolicyTypeVulnerabilityFunction: false,
            PolicyTypeVulnerabilityCiFunction: false,
        },
    }
)

func PolicyTypeToFormattedString(policyType string) (string, error) {
    switch policyType {
        case PolicyTypeAdmission:
            return "Admission Control", nil
        case PolicyTypeComplianceCiImage:
            return "CI Images Compliance", nil
        case PolicyTypeComplianceContainer:
            return "Deployed Containers and Images Compliance", nil
        case PolicyTypeComplianceHost:
            return "Host Compliance", nil
        case PolicyTypeComplianceVmImage:
            return "VM Image Compliance", nil
        case PolicyTypeComplianceFunction:
            return "Function Compliance", nil
        case PolicyTypeComplianceCiFunction:
            return "CI Function Compliance", nil
        case PolicyTypeComplianceTrustedImages:
            return "Trusted Images Compliance", nil
        case PolicyTypeVulnerabilityDeployedImage:
            return "Deployed Images Vulnerability", nil
        case PolicyTypeVulnerabilityCiImage:
            return "CI Images Vulnerability", nil
        case PolicyTypeVulnerabilityHost:
            return "Host Vulnerability", nil
        case PolicyTypeVulnerabilityVmImage:
            return "VM Image Vulnerability", nil
        case PolicyTypeVulnerabilityFunction:
            return "Function Vulnerability", nil
        case PolicyTypeVulnerabilityCiFunction:
            return "CI Function Vulnerability", nil
        case PolicyTypeRuntimeContainer:
            return "Container Runtime", nil
        case PolicyTypeRuntimeHost:
            return "Host Runtime", nil
        case PolicyTypeRuntimeServerless:
            return "Serverless Runtime", nil
        case PolicyTypeRuntimeAppEmbedded:
            return "Application Embedded Runtime", nil
        default:
            return "", errors.New(fmt.Sprintf("unknown policy type \"%s\" specified", policyType))
    }
}

type Policy struct {
	Id            string        `json:"_id"`
	PolicyType    string        `json:"policyType"`
	PolicyContext string        `json:"policyContext"`
	Rules         *[]PolicyRule `json:"rules"`
	Type          string        `json:"type"`
}

func (p *Policy) SortRules(ctx context.Context, planRules *[]models.PolicyRuleResourceModel) {
	util.DLog(ctx, "Executing api.Policy.SortRules()")

	rulesOrderMap := generatePolicyRulesOrderMap(*planRules)
	sort.Slice((*p.Rules), func(i, j int) bool {
		return rulesOrderMap[(*p.Rules)[i].Name] < rulesOrderMap[(*p.Rules)[j].Name]
	})

	util.DLog(ctx, "Finishing api.Policy.SortRules() execution")
}

func (p *Policy) EndpointUrl() string {
	switch p.PolicyType {
	    case PolicyTypeComplianceHost:
		    return HostComplianceEndpoint
	    case PolicyTypeComplianceContainer:
	    	return ContainerComplianceEndpoint
        case PolicyTypeComplianceCiImage:
	    	return CiImageComplianceEndpoint
        case PolicyTypeComplianceVmImage:
	    	return VmImageComplianceEndpoint
        case PolicyTypeComplianceFunction:
	    	return FunctionComplianceEndpoint
        case PolicyTypeComplianceCiFunction:
	    	return CiFunctionComplianceEndpoint
        case PolicyTypeVulnerabilityDeployedImage: 
	    	return DeployedImageVulnerabilityEndpoint
        case PolicyTypeVulnerabilityCiImage:
	    	return CiImageVulnerabilityEndpoint
        case PolicyTypeVulnerabilityHost:
	    	return HostVulnerabilityEndpoint
        case PolicyTypeVulnerabilityVmImage:
	    	return VmImageVulnerabilityEndpoint
        case PolicyTypeVulnerabilityFunction:
	    	return FunctionVulnerabilityEndpoint
        case PolicyTypeVulnerabilityCiFunction:
	    	return CiFunctionVulnerabilityEndpoint
	    default:
	    	return ""
	}
}

func (p *Policy) FormattedType() string {
	switch p.PolicyType {
	    case PolicyTypeComplianceHost:
		    return PolicyTypeComplianceHostFormatted
	    case PolicyTypeComplianceContainer:
	    	return PolicyTypeComplianceContainerFormatted 
        case PolicyTypeComplianceCiImage:
	    	return PolicyTypeComplianceCiImageFormatted
        case PolicyTypeComplianceVmImage:
	    	return PolicyTypeComplianceVmImageFormatted
        case PolicyTypeComplianceFunction:
	    	return PolicyTypeComplianceFunctionFormatted
        case PolicyTypeComplianceCiFunction:
	    	return PolicyTypeComplianceCiFunctionFormatted
        case PolicyTypeVulnerabilityDeployedImage: 
	    	return PolicyTypeVulnerabilityDeployedImageFormatted
        case PolicyTypeVulnerabilityCiImage:
	    	return PolicyTypeVulnerabilityCiImageFormatted
        case PolicyTypeVulnerabilityHost:
	    	return PolicyTypeVulnerabilityHostFormatted
        case PolicyTypeVulnerabilityVmImage:
	    	return PolicyTypeVulnerabilityVmImageFormatted
        case PolicyTypeVulnerabilityFunction:
	    	return PolicyTypeVulnerabilityFunctionFormatted
        case PolicyTypeVulnerabilityCiFunction:
	    	return PolicyTypeVulnerabilityCiFunctionFormatted
	    default:
	    	return ""
	}
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
	AlertThreshold                 AlertThreshold          `json:"alertThreshold" tfsdk:"alert_threshold"`
	BlockMessage                   string                  `json:"blockMsg"`
	BlockThreshold                 BlockThreshold          `json:"blockThreshold" tfsdk:"block_threshold"`
	Collections                    []collection.Collection `json:"collections" tfsdk:"collections"`
	Condition                      *Condition              `json:"condition" tfsdk:"condition"`
	CVERules                       []Exception             `json:"cveRules" tfsdk:"cve_rules"`
	Disabled                       bool                    `json:"disabled"`
	Effect                         string                  `json:"effect"`
	ExcludeBaseImageVulns          bool                    `json:"excludeBaseImageVulns"`
	GraceDays                      int                     `json:"graceDays"`
	GraceDaysPolicy                GraceDaysPolicy         `json:"graceDaysPolicy"`
	Modified                       string                  `json:"modified"`
	Name                           string                  `json:"name"`
	Notes                          string                  `json:"notes"`
	OnlyFixed                      bool                    `json:"onlyFixed"`
	Order                          int                     `json:"order"`
	Owner                          string                  `json:"owner"`
	PkgTypesThresholds             []PkgTypesThreshold     `json:"pkgTypesThresholds"`
	ReportAllPassedAndFailedChecks bool                    `json:"allCompliance"`
	RiskFactorEffects              []RiskFactorsEffect     `json:"riskFactorsEffects"`
	Tags                           []Exception             `json:"tags"`
	Verbose                        bool                    `json:"verbose"`
}

type Condition struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities" tfsdk:"vulnerabilities"`
}

type Vulnerability struct {
	Id    int  `json:"id" tfsdk:"id"`
	Block bool `json:"block" tfsdk:"block"`
}

type RiskFactorsEffect struct {
	Effect     string `json:"effect"`
	RiskFactor string `json:"riskFactor"`
}

type AlertThreshold struct {
	Disabled bool `json:"disabled" tfsdk:"disabled"`
	Value    int  `json:"value" tfsdk:"value"`
}

type BlockThreshold struct {
	Enabled bool `json:"enabled" tfsdk:"enabled"`
	Value   int  `json:"value" tfsdk:"value"`
}

type GraceDaysPolicy struct {
	Enabled  bool `json:"enabled"`
	Low      int  `json:"low"`
	Medium   int  `json:"medium"`
	High     int  `json:"high"`
	Critical int  `json:"critical"`
}

type PkgTypesThreshold struct {
	Type           string         `json:"type"`
	AlertThreshold AlertThreshold `json:"alertThreshold"`
	BlockThreshold BlockThreshold `json:"blockThreshold"`
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
	Enabled bool   `json:"enabled"`
	Date    string `json:"date"`
}

type PortRange struct {
    Deny    bool `json:"deny,omitempty"`
    Start   int `json:"start,omitempty"`
    End     int `json:"end,omitempty"`
}

type NetworkPorts struct {
	Allowed []PortRange `json:"allowed,omitempty"`
	Denied  []PortRange `json:"denied,omitempty"`
	Effect  string                 `json:"effect,omitempty"`
}

type DeniedProcesses struct {
	Effect string `json:"effect"`
	Paths  []string `json:"paths"`
}

type CustomRule struct {
	ID     int      `json:"_id"`
	Action string `json:"action"`
	Effect string `json:"effect"`
}

type DnsDomainList struct {
	Allowed []string `json:"allowed,omitempty"`
	Denied  []string `json:"denied,omitempty"`
	Effect  string   `json:"effect,omitempty"`
}


func getEndpointAndPolicyName(policyType string) (string, string, error) {
	switch policyType {
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
		return DeployedImageVulnerabilityEndpoint, "deployed image", nil
	case "ciImagesVulnerability":
		return CiImageVulnerabilityEndpoint, "CI image", nil
	case "hostVulnerability":
		return HostVulnerabilityEndpoint, "host", nil
	case "vmVulnerability":
		return VmImageVulnerabilityEndpoint, "VM image", nil
	case "serverlessVulnerability":
		return FunctionVulnerabilityEndpoint, "function", nil
	case "ciServerlessVulnerability":
		return CiFunctionVulnerabilityEndpoint, "CI function", nil
	default:
		return "", "", fmt.Errorf("invalid policy type specified")
	}
}

func UpsertPolicy(c api.PrismaCloudComputeAPIClient, policy Policy) error {
	if err := c.Request(http.MethodPut, policy.EndpointUrl(), nil, policy, nil); err != nil {
		return fmt.Errorf("Error occured while attempting to upsert %s policy: %s", policy.FormattedType(), err)
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
