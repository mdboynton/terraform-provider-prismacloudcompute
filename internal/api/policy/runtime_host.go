package policy

import (
    "context"
	"fmt"
	"net/http"
    "sort"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
)

//type RuntimeHostPolicy struct {
//	Rules []RuntimeHostRule `json:"rules,omitempty"`
//}
//
//type RuntimeHostRule struct {
//	AntiMalware        RuntimeHostAntiMalware         `json:"antiMalware,omitempty"`
//	Collections        []collection.Collection        `json:"collections,omitempty"`
//	CustomRules        []RuntimeHostCustomRule        `json:"customRules,omitempty"`
//	Disabled           bool                           `json:"disabled"`
//	Dns                RuntimeHostDns                 `json:"dns,omitempty"`
//	FileIntegrityRules []RuntimeHostFileIntegrityRule `json:"fileIntegrityRules,omitempty"`
//	Forensic           RuntimeHostForensic            `json:"forensic,omitempty"`
//	LogInspectionRules []RuntimeHostLogInspectionRule `json:"logInspectionRules,omitempty"`
//	Name               string                         `json:"name,omitempty"`
//	Network            RuntimeHostNetwork             `json:"network,omitempty"`
//	Notes              string                         `json:"notes,omitempty"`
//}
//
//type RuntimeHostAntiMalware struct {
//	AllowedProcesses              []string                   `json:"allowedProcesses,omitempty"`
//	CryptoMiner                   string                     `json:"cryptoMiner,omitempty"`
//	CustomFeed                    string                     `json:"customFeed,omitempty"`
//	DeniedProcesses               RuntimeHostDeniedProcesses `json:"deniedProcesses,omitempty"`
//	DetectCompilerGeneratedBinary bool                       `json:"detectCompilerGeneratedBinary"`
//	EncryptedBinaries             string                     `json:"encryptedBinaries,omitempty"`
//	ExecutionFlowHijack           string                     `json:"executionFlowHijack,omitempty"`
//	IntelligenceFeed              string                     `json:"intelligenceFeed,omitempty"`
//	ReverseShell                  string                     `json:"reverseShell,omitempty"`
//	ServiceUnknownOriginBinary    string                     `json:"serviceUnknownOriginBinary,omitempty"`
//	SkipSshTracking               bool                       `json:"skipSSHTracking,omitempty"`
//	SuspiciousElfHeaders          string                     `json:"suspiciousELFHeaders,omitempty"`
//	TempFsProcesses               string                     `json:"tempFSProc,omitempty"`
//	UserUnknownOriginBinary       string                     `json:"userUnknownOriginBinary,omitempty"`
//	WebShell                      string                     `json:"webShell,omitempty"`
//	WildFireAnalysis              string                     `json:"wildFireAnalysis,omitempty"`
//}
//
//type RuntimeHostCustomRule struct {
//	Action string `json:"action,omitempty"`
//	Effect string `json:"effect,omitempty"`
//	Id     int    `json:"_id,omitempty"`
//}
//
//type RuntimeHostDeniedProcesses struct {
//	Effect string   `json:"effect,omitempty"`
//	Paths  []string `json:"paths,omitempty"`
//}
//
//type RuntimeHostDns struct {
//	Allowed          []string `json:"allow,omitempty"`
//	Denied           []string `json:"deny,omitempty"`
//	DenyEffect       string   `json:"denyListEffect,omitempty"`
//	IntelligenceFeed string   `json:"intelligenceFeed,omitempty"`
//}
//
//type RuntimeHostFileIntegrityRule struct {
//	AllowedProcesses []string `json:"procWhitelist,omitempty"`
//	ExcludedFiles    []string `json:"exclusions,omitempty"`
//	Metadata         bool     `json:"metadata"`
//	Path             string   `json:"path,omitempty"`
//	Read             bool     `json:"read"`
//	Recursive        bool     `json:"recursive"`
//	Write            bool     `json:"write"`
//}
//
//type RuntimeHostForensic struct {
//	ActivitiesDisabled       bool `json:"activitiesDisabled"`
//	DockerEnabled            bool `json:"dockerEnabled"`
//	ReadonlyDockerEnabled    bool `json:"readonlyDockerEnabled"`
//	ServiceActivitiesEnabled bool `json:"serviceActivitiesEnabled"`
//	SshdEnabled              bool `json:"sshdEnabled"`
//	SudoEnabled              bool `json:"sudoEnabled"`
//}
//
//type RuntimeHostLogInspectionRule struct {
//	Path  string   `json:"path,omitempty"`
//	Regex []string `json:"regex,omitempty"`
//}
//
//type RuntimeHostNetwork struct {
//	AllowedOutboundIps   []string          `json:"allowedOutboundIPs,omitempty"`
//	CustomFeed           string            `json:"customFeed,omitempty"`
//	DeniedListeningPorts []RuntimeHostPort `json:"deniedListeningPorts,omitempty"`
//	DeniedOutboundIps    []string          `json:"deniedOutboundIPs,omitempty"`
//	DeniedOutboundPorts  []RuntimeHostPort `json:"deniedOutboundPorts,omitempty"`
//	DenyEffect           string            `json:"denyListEffect,omitempty"`
//	IntelligenceFeed     string            `json:"intelligenceFeed,omitempty"`
//}
//
//type RuntimeHostPort struct {
//	Deny  bool `json:"deny"`
//	End   int  `json:"end,omitempty"`
//	Start int  `json:"start,omitempty"`
//}

type RuntimeHostPolicy struct {
    Id      string      `json:"_id"`
    Owner   string      `json:"owner"`
    Rules *[]RuntimeHostPolicyRule
}

type RuntimeHostPolicyRule struct {
	AntiMalware       AntiMalware     `json:"antiMalware"`
	Collections       []collectionAPI.Collection    `json:"collections"`
	CustomRules       []CustomRule    `json:"customRules"`
	Disabled          bool            `json:"disabled"`
	DNS               DNS             `json:"dns"`
	FileIntegrityRules []FileIntegrityRule `json:"fileIntegrityRules"`
	Forensic          Forensic        `json:"forensic"`
	LogInspectionRules []LogInspectionRule `json:"logInspectionRules"`
	Modified          string          `json:"modified"`
	Name              string          `json:"name"`
	Network           Network         `json:"network"`
	Notes             string          `json:"notes"`
	//Order                          int                     `json:"order"`
	Owner             string          `json:"owner"`
	PreviousName      string          `json:"previousName"`
}

type AntiMalware struct {
	AllowedProcesses          []string `json:"allowedProcesses"`
	CryptoMiner               string `json:"cryptoMiner"`
	CustomFeed                string `json:"customFeed"`
	DeniedProcesses           DeniedProcesses `json:"deniedProcesses"`
	DetectCompilerGeneratedBinary bool     `json:"detectCompilerGeneratedBinary"` // suppress_compiler_generated_binaries
	EncryptedBinaries         string `json:"encryptedBinaries"`
	ExecutionFlowHijack       string `json:"executionFlowHijack"`
	IntelligenceFeed          string `json:"intelligenceFeed"`
	ReverseShell              string `json:"reverseShell"`
	ServiceUnknownOriginBinary string `json:"serviceUnknownOriginBinary"`
	SkipSSHTracking           bool     `json:"skipSSHTracking"`
	SuspiciousELFHeaders      string `json:"suspiciousELFHeaders"`
	TempFSProc                string `json:"tempFSProc"`
	UserUnknownOriginBinary   string `json:"userUnknownOriginBinary"`
	WebShell                  string `json:"webShell"`
	WildFireAnalysis          string `json:"wildFireAnalysis"`
}

type DeniedProcesses struct {
	Effect string `json:"effect"`
	Paths  []string `json:"paths"`
}

type CustomRule struct {
	ID     int      `json:"_id"`
	Action []string `json:"action"`
	Effect []string `json:"effect"`
}

type DNS struct {
	Allow            []string `json:"allow"`
	Deny             []string `json:"deny"`
	DenyListEffect   string `json:"denyListEffect"`
	IntelligenceFeed string `json:"intelligenceFeed"`
}

type FileIntegrityRule struct {
	Dir           bool     `json:"dir"`
	Exclusions    []string `json:"exclusions"`
	Metadata      bool     `json:"metadata"`
	Path          string   `json:"path"`
	ProcWhitelist []string `json:"procWhitelist"`
	Read          bool     `json:"read"`
	Recursive     bool     `json:"recursive"`
	Write         bool     `json:"write"`
}

type Forensic struct {
	ActivitiesDisabled    bool `json:"activitiesDisabled"`
	DockerEnabled         bool `json:"dockerEnabled"`
	ReadonlyDockerEnabled bool `json:"readonlyDockerEnabled"`
	ServiceActivitiesEnabled bool `json:"serviceActivitiesEnabled"`
	SshdEnabled           bool `json:"sshdEnabled"`
	SudoEnabled           bool `json:"sudoEnabled"`
}

type LogInspectionRule struct {
	Path  string   `json:"path"`
	Regex []string `json:"regex"`
}

type Network struct {
	AllowedOutboundIPs []string      `json:"allowedOutboundIPs"`
	CustomFeed         string      `json:"customFeed"`
	DeniedListeningPorts []PortRange   `json:"deniedListeningPorts"`
	DeniedOutboundIPs  []string      `json:"deniedOutboundIPs"`
	DeniedOutboundPorts  []PortRange   `json:"deniedOutboundPorts"`
	DenyListEffect     string      `json:"denyListEffect"`
	IntelligenceFeed   string      `json:"intelligenceFeed"`
}

type PortRange struct {
	Deny  bool `json:"deny"`
	End   int  `json:"end"`
	Start int  `json:"start"`
}

func (p *RuntimeHostPolicy) SortRules(ctx context.Context, planRules *[]models.RuntimeHostPolicyRuleResourceModel) {
	//util.DLog(ctx, "Executing api.Policy.SortRules()")
    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
        return
    }

	rulesOrderMap := generateRuntimePolicyRulesOrderMap(*planRules)
	sort.Slice((*p.Rules), func(i, j int) bool {
		return rulesOrderMap[(*p.Rules)[i].Name] < rulesOrderMap[(*p.Rules)[j].Name]
	})

	//util.DLog(ctx, "Finishing api.Policy.SortRules() execution")
}

// TODO: remove this duplicate function when we can move the logic somewhere that can be used here
// and by internal/resources/policy/common.go
func generateRuntimePolicyRulesOrderMap(rules []models.RuntimeHostPolicyRuleResourceModel) map[string]int {
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

// Get the current host runtime policy.
func GetRuntimeHost(c api.PrismaCloudComputeAPIClient) (RuntimeHostPolicy, error) {
	var ans RuntimeHostPolicy
	if err := c.Request(http.MethodGet, RuntimeHostEndpoint, nil, nil, &ans); err != nil {
		return ans, fmt.Errorf("error getting host runtime policy: %s", err)
	}
	return ans, nil
}

// Update the current host runtime policy.
//func UpdateRuntimeHost(c api.PrismaCloudComputeAPIClient, policy RuntimeHostPolicy) error {
func UpsertRuntimeHost(c api.PrismaCloudComputeAPIClient, policy RuntimeHostPolicy) error {
    // TODO: error handling
	return c.Request(http.MethodPut, RuntimeHostEndpoint, nil, policy, nil)
}
