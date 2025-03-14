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

type RuntimeContainerPolicy struct {
    Id                  string                          `json:"_id"`
	Rules               *[]RuntimeContainerPolicyRule    `json:"rules,omitempty"`
	LearningDisabled    bool                            `json:"learningDisabled,omitempty"`
}

type RuntimeContainerPolicyRule struct {
	AdvancedProtectionEffect       string                       `json:"advancedProtectionEffect"`
	CloudMetadataEnforcementEffect string                       `json:"cloudMetadataEnforcementEffect"`
	Collections                    []collectionAPI.Collection      `json:"collections,omitempty"`
	CustomRules                    []CustomRule `json:"customRules,omitempty"`
	Disabled                       bool                         `json:"disabled"`
	Dns                            RuntimeContainerDns          `json:"dns,omitempty"`
	Filesystem                     RuntimeContainerFilesystem   `json:"filesystem,omitempty"`
	KubernetesEnforcementEffect    string                       `json:"kubernetesEnforcementEffect"`
	Name                           string                       `json:"name,omitempty"`
	PreviousName                   string                       `json:"previousName,omitempty"`
	SkipExecSessions               bool                         `json:"skipExecSessions,omitempty"`
	Modified                       string                       `json:"modified"`
	Network                        RuntimeContainerNetwork      `json:"network,omitempty"`
	Notes                          string                       `json:"notes,omitempty"`
	Owner                          string                       `json:"owner,omitempty"`
	Processes                      RuntimeContainerProcesses    `json:"processes,omitempty"`
	WildFireAnalysis               string                       `json:"wildFireAnalysis,omitempty"`
}

type RuntimeContainerDns struct {
	DefaultEffect string                        `json:"defaultEffect,omitempty"`
	Disabled      bool                          `json:"disabled,omitempty"`
	DomainList    RuntimeContainerDnsDomainList `json:"domainList,omitempty"`
}

type RuntimeContainerFilesystem struct {
	AllowedList                []string                   `json:"allowedList,omitempty"`
	BackdoorFilesEffect        string                     `json:"backdoorFilesEffect,omitempty"`
	DefaultEffect              string                     `json:"defaultEffect,omitempty"`
	DeniedList                 RuntimeContainerDeniedList `json:"deniedList,omitempty"`
	Disabled                   bool                       `json:"disabled,omitempty"`
	EncryptedBinariesEffect    string                     `json:"encryptedBinariesEffect,omitempty"`
	NewFilesEffect             string                     `json:"newFilesEffect,omitempty"`
	SuspiciousElfHeadersEffect string                     `json:"suspiciousElfHeadersEffect,omitempty"`
}

type RuntimeContainerNetwork struct {
	AllowedIps         []string                     `json:"allowedIPs,omitempty"`
	DefaultEffect      string                       `json:"defaultEffect,omitempty"`
	DeniedIps          []string                     `json:"deniedIPs,omitempty"`
	DeniedIpsEffect    string                       `json:"deniedIPsEffect,omitempty"`
	Disabled           bool                         `json:"disabled,omitempty"`
	ListeningPorts     NetworkPorts                 `json:"listeningPorts,omitempty"`
	ModifiedProcEffect string                       `json:"modifiedProcEffect,omitempty"`
	OutboundPorts      NetworkPorts                 `json:"outboundPorts,omitempty"`
	PortScanEffect     string                       `json:"portScanEffect,omitempty"`
	RawSocketsEffect   string                       `json:"rawSocketsEffect,omitempty"`
}

type RuntimeContainerDnsDomainList struct {
	Allowed []string `json:"allowed,omitempty"`
	Denied  []string `json:"denied,omitempty"`
	Effect  string   `json:"effect,omitempty"`
}

type RuntimeContainerPort struct {
	Deny  bool `json:"deny"`
	End   int  `json:"end,omitempty"`
	Start int  `json:"start,omitempty"`
}

type RuntimeContainerProcesses struct {
	ModifiedProcessEffect string                     `json:"modifiedProcessEffect,omitempty"`
	CryptoMinersEffect    string                     `json:"cryptoMinersEffect,omitempty"`
	LateralMovementEffect string                     `json:"lateralMovementEffect,omitempty"`
	ReverseShellEffect    string                     `json:"reverseShellEffect,omitempty"`
	SuidBinariesEffect    string                     `json:"suidBinariesEffect,omitempty"`
	DefaultEffect         string                     `json:"defaultEffect,omitempty"`
	CheckParentChild      bool                       `json:"checkParentChild"`
	AllowedList           []string                   `json:"allowedList,omitempty"`
	Disabled              bool                       `json:"disabled"`
	DeniedList            RuntimeContainerDeniedList `json:"deniedList"`
}

type RuntimeContainerDeniedList struct {
	Effect string   `json:"effect,omitempty"`
	Paths  []string `json:"paths,omitempty"`
}

func (p *RuntimeContainerPolicy) SortRules(ctx context.Context, planRules *[]models.RuntimeContainerPolicyRuleResourceModel) {
	util.DLog(ctx, "Executing api.RuntimeContainerPolicy.SortRules()")
    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
        return
    }

	rulesOrderMap := generateRuntimeContainerPolicyRulesOrderMap(*planRules)
	sort.Slice((*p.Rules), func(i, j int) bool {
		return rulesOrderMap[(*p.Rules)[i].Name] < rulesOrderMap[(*p.Rules)[j].Name]
	})

	util.DLog(ctx, "Finishing api.RuntimeContainerPolicy.SortRules() execution")
}

// TODO: remove this duplicate function when we can move the logic somewhere that can be used here
// and by internal/resources/policy/common.go
func generateRuntimeContainerPolicyRulesOrderMap(rules []models.RuntimeContainerPolicyRuleResourceModel) map[string]int {
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

// Get the current container runtime policy.
func GetRuntimeContainer(c api.PrismaCloudComputeAPIClient) (RuntimeContainerPolicy, error) {
	var ans RuntimeContainerPolicy
	if err := c.Request(http.MethodGet, RuntimeContainerEndpoint, nil, nil, &ans); err != nil {
		return ans, fmt.Errorf("error getting container runtime policy: %s", err)
	}
	return ans, nil
}

// Update the current container runtime policy.
func UpsertRuntimeContainer(c api.PrismaCloudComputeAPIClient, policy RuntimeContainerPolicy) error {
	return c.Request(http.MethodPut, RuntimeContainerEndpoint, nil, policy, nil)
}
