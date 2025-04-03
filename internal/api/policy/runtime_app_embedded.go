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

type RuntimeAppEmbeddedPolicy struct {
    Id      string      `json:"_id"`
    Rules *[]RuntimeAppEmbeddedPolicyRule `json:"rules"`
	LearningDisabled    bool                            `json:"learningDisabled,omitempty"`
}

type RuntimeAppEmbeddedPolicyRule struct {
	AdvancedProtection              bool                       `json:"advancedProtection"`
	Collections       []collectionAPI.Collection    `json:"collections"`
	CustomRules       []CustomRule    `json:"customRules"`
	Disabled          bool            `json:"disabled"`
	Dns               RuntimeAppEmbeddedDns             `json:"dns"`
	Filesystem        RuntimeAppEmbeddedFilesystem `json:"filesystem"`
	Modified          string          `json:"modified"`
	Name              string          `json:"name"`
	Network           RuntimeAppEmbeddedNetwork         `json:"network"`
	Notes             string          `json:"notes"`
	Owner             string          `json:"owner"`
	Processes         RuntimeAppEmbeddedProcesses    `json:"processes,omitempty"`
	PreviousName      string          `json:"previousName"`
    WildFireAnalysis  string          `json:"wildFireAnalysis"`
}

type RuntimeAppEmbeddedDns struct {
	Whitelist            []string `json:"whitelist"`
	Effect   string `json:"effect"`
}

type RuntimeAppEmbeddedFilesystem struct {
	BackdoorFiles           bool        `json:"backdoorFiles,omitempty"`
    Blacklist               []string    `json:"blacklist,omitempty"`
	CheckNewFiles           bool        `json:"checkNewFiles,omitempty"`
    Effect                  string      `json:"effect,omitempty"` 
	SkipEncryptedBinaries   bool        `json:"skipEncryptedBinaries,omitempty"`
	SuspiciousElfHeaders    bool        `json:"suspiciousElfHeaders,omitempty"`
    Whitelist               []string    `json:"whitelist,omitempty"`
}

type RuntimeAppEmbeddedNetwork struct {
	WhitelistIPs []string      `json:"whitelistIPs"`
	WhitelistListeningPorts []PortRange   `json:"whitelistListeningPorts"`
	WhitelistOutboundPorts  []PortRange   `json:"whitelistOutboundPorts"`
	BlacklistIPs  []string      `json:"blacklistIPs"`
	BlacklistListeningPorts []PortRange   `json:"blacklistListeningPorts"`
	BlacklistOutboundPorts  []PortRange   `json:"blacklistOutboundPorts"`
	Effect     string      `json:"effect"`
}

type RuntimeAppEmbeddedProcesses struct {
	Blacklist            []string                     `json:"blacklist,omitempty"`
	CheckCryptoMiners      bool                       `json:"checkCryptoMiners,omitempty"`
	CheckNewBinaries      bool                       `json:"checkNewBinaries,omitempty"`
	Effect                  string                     `json:"effect,omitempty"`
	Whitelist            []string                     `json:"whitelist,omitempty"`
}

func (p *RuntimeAppEmbeddedPolicy) GetRuleNames() []string {
    ruleNames := []string{}

    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
        return ruleNames
    }

    for _, rule := range *p.Rules {
        ruleNames = append(ruleNames, rule.Name)
    }

    return ruleNames
}

func (p *RuntimeAppEmbeddedPolicy) SortRules(ctx context.Context, planRules *[]models.RuntimeAppEmbeddedPolicyRuleResourceModel) {
	util.DLog(ctx, "Executing api.RuntimeAppEmbeddedPolicy.SortRules()")
    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
        return
    }

	rulesOrderMap := generateRuntimeAppEmbeddedPolicyRulesOrderMap(*planRules)
	sort.Slice((*p.Rules), func(i, j int) bool {
		return rulesOrderMap[(*p.Rules)[i].Name] < rulesOrderMap[(*p.Rules)[j].Name]
	})

	util.DLog(ctx, "Finishing api.RuntimeAppEmbeddedPolicy.SortRules() execution")
}

// TODO: remove this duplicate function when we can move the logic somewhere that can be used here
// and by internal/resources/policy/common.go
func generateRuntimeAppEmbeddedPolicyRulesOrderMap(rules []models.RuntimeAppEmbeddedPolicyRuleResourceModel) map[string]int {
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

// Get app-embedded runtime policy
func GetRuntimeAppEmbeddedPolicy(c api.PrismaCloudComputeAPIClient) (RuntimeAppEmbeddedPolicy, error) {
	var ans RuntimeAppEmbeddedPolicy
	if err := c.Request(http.MethodGet, RuntimeAppEmbeddedEndpoint, nil, nil, &ans); err != nil {
		return ans, fmt.Errorf("error getting app-embedded runtime policy: %s", err)
	}
	return ans, nil
}

// Get app-embedded runtime policy filtered on rule names
func GetRuntimeAppEmbeddedPolicyFiltered(c api.PrismaCloudComputeAPIClient, ruleNames []string) (RuntimeAppEmbeddedPolicy, error) {
    policy, err := GetRuntimeAppEmbeddedPolicy(c)
    if err != nil {
        return RuntimeAppEmbeddedPolicy{}, err
    }

    if (policy.Rules == nil || len(*policy.Rules) == 0) {
        return policy, nil
    }

    filteredRules := []RuntimeAppEmbeddedPolicyRule{}
    for _, rule := range *policy.Rules {
        if slices.Contains(ruleNames, rule.Name) {
            filteredRules = append(filteredRules, rule)
        }
    }

    response := RuntimeAppEmbeddedPolicy{
        Id: policy.Id,
        LearningDisabled: policy.LearningDisabled,
        Rules: &filteredRules,
    }

    return response, nil
}

// Update app-embedded runtime policy
func UpsertRuntimeAppEmbeddedPolicy(c api.PrismaCloudComputeAPIClient, policy RuntimeAppEmbeddedPolicy) error {
	return c.Request(http.MethodPut, RuntimeAppEmbeddedEndpoint, nil, policy, nil)
}

// Update app-embedded runtime policy, scoped to specified rule names
// Rules not managed through Terraform will be preserved
func UpsertRuntimeAppEmbeddedPolicyFiltered(c api.PrismaCloudComputeAPIClient, policy RuntimeAppEmbeddedPolicy, deletedRuleNames []string) error {
    currentPolicy, err := GetRuntimeAppEmbeddedPolicy(c)
    if err != nil {
        return err
    }

    if (currentPolicy.Rules == nil || len(*currentPolicy.Rules) == 0) {
	    return c.Request(http.MethodPut, RuntimeAppEmbeddedEndpoint, nil, policy, nil)
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

	return c.Request(http.MethodPut, RuntimeAppEmbeddedEndpoint, nil, policy, nil)
}
