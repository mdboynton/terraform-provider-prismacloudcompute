package policy

import (
    "context"
	"fmt"
	"net/http"
    "sort"
    "slices"

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
	DomainList    DnsDomainList                 `json:"domainList,omitempty"`
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

func (p *RuntimeContainerPolicy) GetRuleNames() []string {
    ruleNames := []string{}

    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
        return ruleNames
    }

    for _, rule := range *p.Rules {
        ruleNames = append(ruleNames, rule.Name)
    }

    return ruleNames
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

// Get container runtime policy
func GetRuntimeContainerPolicy(c api.PrismaCloudComputeAPIClient) (RuntimeContainerPolicy, error) {
	var ans RuntimeContainerPolicy
	if err := c.Request(http.MethodGet, RuntimeContainerEndpoint, nil, nil, &ans); err != nil {
		return ans, fmt.Errorf("error getting container runtime policy: %s", err)
	}
	return ans, nil
}

// Get container runtime policy filtered on rule names
func GetRuntimeContainerPolicyFiltered(c api.PrismaCloudComputeAPIClient, ruleNames []string) (RuntimeContainerPolicy, error) {
    policy, err := GetRuntimeContainerPolicy(c)
    if err != nil {
        return RuntimeContainerPolicy{}, err
    }

    if (policy.Rules == nil || len(*policy.Rules) == 0) {
        return policy, nil
    }

    filteredRules := []RuntimeContainerPolicyRule{}
    for _, rule := range *policy.Rules {
        if slices.Contains(ruleNames, rule.Name) {
            filteredRules = append(filteredRules, rule)
        }
    }

    response := RuntimeContainerPolicy{
        Id: policy.Id,
        LearningDisabled: policy.LearningDisabled,
        Rules: &filteredRules,
    }

    return response, nil
}

// Update container runtime policy
func UpsertRuntimeContainerPolicy(c api.PrismaCloudComputeAPIClient, policy RuntimeContainerPolicy) error {
	return c.Request(http.MethodPut, RuntimeContainerEndpoint, nil, policy, nil)
}

// Update container runtime policy, scoped to specified rule names
// Rules not managed through Terraform will be preserved
func UpsertRuntimeContainerPolicyFiltered(c api.PrismaCloudComputeAPIClient, policy RuntimeContainerPolicy, deletedRuleNames []string) error {
    currentPolicy, err := GetRuntimeContainerPolicy(c)
    if err != nil {
        return err
    }

    if (currentPolicy.Rules == nil || len(*currentPolicy.Rules) == 0) {
	    return c.Request(http.MethodPut, RuntimeContainerEndpoint, nil, policy, nil)
    }

    updatedRuleNames := policy.GetRuleNames()
    updatedRules := *policy.Rules

    // Loop through the current policy rules and add any rules that are not present in the updated policy's rules and are not being deleted
    for _, currentPolicyRule := range *currentPolicy.Rules {
        // If the rule is being deleted, do not include it in the updated ruleset
        if slices.Contains(deletedRuleNames, currentPolicyRule.Name) {
            continue
        }

        if !slices.Contains(updatedRuleNames, currentPolicyRule.Name) {
            updatedRules = append(updatedRules, currentPolicyRule)
        }
    }

    policy.Rules = &updatedRules

	return c.Request(http.MethodPut, RuntimeContainerEndpoint, nil, policy, nil)
}
