package policy

import (
    "context"
	"fmt"
	"net/http"
    "sort"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
)

type RuntimeHostPolicy struct {
    Id      string      `json:"_id"`
    Owner   string      `json:"owner"`
    Rules *[]RuntimeHostPolicyRule `json:"rules"`
}

type RuntimeHostPolicyRule struct {
	AntiMalware       RuntimeHostAntiMalware     `json:"antiMalware"`
	Collections       []collectionAPI.Collection    `json:"collections"`
	CustomRules       []CustomRule    `json:"customRules"`
	Disabled          bool            `json:"disabled"`
	DNS               RuntimeHostDns             `json:"dns"`
	FileIntegrityRules []FileIntegrityRule `json:"fileIntegrityRules"`
	Forensic          Forensic        `json:"forensic"`
	LogInspectionRules []LogInspectionRule `json:"logInspectionRules"`
	Modified          string          `json:"modified"`
	Name              string          `json:"name"`
	Network           RuntimeHostNetwork         `json:"network"`
	Notes             string          `json:"notes"`
	//Order                          int                     `json:"order"`
	Owner             string          `json:"owner"`
	PreviousName      string          `json:"previousName"`
}

type RuntimeHostAntiMalware struct {
	AllowedProcesses          []string `json:"allowedProcesses"`
	CryptoMiner               string `json:"cryptoMiner"`
	CustomFeed                string `json:"customFeed"`
	DeniedProcesses           DeniedProcesses `json:"deniedProcesses"`
	DetectCompilerGeneratedBinary bool     `json:"detectCompilerGeneratedBinary"`
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

type RuntimeHostDns struct {
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

type RuntimeHostNetwork struct {
	AllowedOutboundIPs []string      `json:"allowedOutboundIPs"`
	CustomFeed         string      `json:"customFeed"`
	DeniedListeningPorts []PortRange   `json:"deniedListeningPorts"`
	DeniedOutboundIPs  []string      `json:"deniedOutboundIPs"`
	DeniedOutboundPorts  []PortRange   `json:"deniedOutboundPorts"`
	DenyListEffect     string      `json:"denyListEffect"`
	IntelligenceFeed   string      `json:"intelligenceFeed"`
}

//type PortRange struct {
//	Deny  bool `json:"deny"`
//	End   int  `json:"end"`
//	Start int  `json:"start"`
//}

func (p *RuntimeHostPolicy) SortRules(ctx context.Context, planRules *[]models.RuntimeHostPolicyRuleResourceModel) {
	util.DLog(ctx, "Executing api.RuntimeHostPolicy.SortRules()")
    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
        return
    }

	rulesOrderMap := generateRuntimePolicyRulesOrderMap(*planRules)
	sort.Slice((*p.Rules), func(i, j int) bool {
		return rulesOrderMap[(*p.Rules)[i].Name] < rulesOrderMap[(*p.Rules)[j].Name]
	})

	util.DLog(ctx, "Finishing api.RuntimeHostPolicy.SortRules() execution")
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
func UpsertRuntimeHost(c api.PrismaCloudComputeAPIClient, policy RuntimeHostPolicy) error {
    // TODO: error handling
	return c.Request(http.MethodPut, RuntimeHostEndpoint, nil, policy, nil)
}
