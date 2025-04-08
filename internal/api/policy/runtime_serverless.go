package policy

import (
    //"context"
	"fmt"
	"net/http"
    //"sort"
    "slices"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
)

type RuntimeServerlessPolicy struct {
    Id                  string                          `json:"_id"`
	Rules               *[]RuntimeServerlessPolicyRule    `json:"rules,omitempty"`
	LearningDisabled    bool                            `json:"learningDisabled,omitempty"`
}

type RuntimeServerlessPolicyRule struct {
	AdvancedProtection              bool                       `json:"advancedProtection"`
	CloudMetadataEnforcementEffect  string                       `json:"cloudMetadataEnforcementEffect"`
	Collections                    []collectionAPI.Collection      `json:"collections,omitempty"`
	CustomRules                    []CustomRule `json:"customRules,omitempty"`
	Disabled                       bool                         `json:"disabled"`
	Dns                            RuntimeServerlessDns          `json:"dns,omitempty"`
	Filesystem                     RuntimeServerlessFilesystem   `json:"filesystem,omitempty"`
	KubernetesEnforcementEffect    string                       `json:"kubernetesEnforcementEffect"`
	Name                           string                       `json:"name,omitempty"`
	Network                        RuntimeServerlessNetwork      `json:"network,omitempty"`
	Notes                          string                       `json:"notes,omitempty"`
	Owner                          string                       `json:"owner,omitempty"`
	PreviousName                   string                       `json:"previousName,omitempty"`
	Processes                      RuntimeServerlessProcesses    `json:"processes,omitempty"`
	SkipExecSessions               bool                         `json:"skipExecSessions,omitempty"`
	Modified                       string                       `json:"modified"`
	WildFireAnalysis               string                       `json:"wildFireAnalysis,omitempty"`
}

type RuntimeServerlessDns struct {
	Effect      string      `json:"effect,omitempty"`
	Blacklist   []string    `json:"blacklist,omitempty"`
	Whitelist   []string    `json:"whitelist,omitempty"`
}

type RuntimeServerlessFilesystem struct {
	BackdoorFiles           bool        `json:"backdoorFiles,omitempty"`
    Blacklist               []string    `json:"blacklist,omitempty"`
	CheckNewFiles           bool        `json:"checkNewFiles,omitempty"`
    Effect                  string      `json:"effect,omitempty"` 
	SkipEncryptedBinares    bool        `json:"skipEncryptedBinaries,omitempty"`
	SuspiciousElfHeaders    bool        `json:"suspiciousElfHeaders,omitempty"`
    Whitelist               []string    `json:"whitelist,omitempty"`
}

type RuntimeServerlessNetwork struct {
	BlacklistIps         []string                     `json:"blacklistIPs,omitempty"`
    BlacklistListeningPorts []PortRange             `json:"blacklistListeningPorts,omitempty"`
    BlacklistOutboundPorts []PortRange             `json:"blacklistOutboundPorts,omitempty"`
	Effect              string                     `json:"effect,omitempty"`
	WhitelistIps         []string                     `json:"whitelistIPs,omitempty"`
    WhitelistListeningPorts []PortRange             `json:"whitelistListeningPorts,omitempty"`
    WhitelistOutboundPorts []PortRange             `json:"whitelistOutboundPorts,omitempty"`
}

type RuntimeServerlessDnsDomainList struct {
	Allowed []string `json:"allowed,omitempty"`
	Denied  []string `json:"denied,omitempty"`
	Effect  string   `json:"effect,omitempty"`
}

type RuntimeServerlessPort struct {
	Deny  bool `json:"deny"`
	End   int  `json:"end,omitempty"`
	Start int  `json:"start,omitempty"`
}

type RuntimeServerlessProcesses struct {
	Blacklist            []string                     `json:"blacklist,omitempty"`
	BlockAllBinaries      bool                       `json:"blockAllBinaries,omitempty"`
	CheckCryptoMiners      bool                       `json:"checkCryptoMiners,omitempty"`
	CheckLateralMovement      bool                       `json:"checkLateralMovement,omitempty"`
	CheckNewBinaries      bool                       `json:"checkNewBinaries,omitempty"`
	Effect                  string                     `json:"effect,omitempty"`
	SkipModified      bool                       `json:"skipModified,omitempty"`
	Whitelist            []string                     `json:"whitelist,omitempty"`
}

type RuntimeServerlessDeniedList struct {
	Effect string   `json:"effect,omitempty"`
	Paths  []string `json:"paths,omitempty"`
}

func (p *RuntimeServerlessPolicy) GetRuleNames() []string {
    ruleNames := []string{}

    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
        return ruleNames
    }

    for _, rule := range *p.Rules {
        ruleNames = append(ruleNames, rule.Name)
    }

    return ruleNames
}

//func (p *RuntimeServerlessPolicy) SortRules(ctx context.Context, planRules *[]models.RuntimeServerlessPolicyRuleResourceModel) {
//	util.DLog(ctx, "Executing api.RuntimeServerlessPolicy.SortRules()")
//    if (p == nil || (*p).Rules == nil || len(*p.Rules) == 0) {
//        return
//    }
//
//	rulesOrderMap := generateRuntimeServerlessPolicyRulesOrderMap(*planRules)
//	sort.Slice((*p.Rules), func(i, j int) bool {
//		return rulesOrderMap[(*p.Rules)[i].Name] < rulesOrderMap[(*p.Rules)[j].Name]
//	})
//
//	util.DLog(ctx, "Finishing api.RuntimeServerlessPolicy.SortRules() execution")
//}
//
//// TODO: remove this duplicate function when we can move the logic somewhere that can be used here
//// and by internal/resources/policy/common.go
//func generateRuntimeServerlessPolicyRulesOrderMap(rules []models.RuntimeServerlessPolicyRuleResourceModel) map[string]int {
//	orderedRulesMap := make(map[int][]string)
//
//	for _, rule := range rules {
//		order := int(rule.Order.ValueInt32())
//		if _, exists := orderedRulesMap[order]; exists {
//			orderedRulesMap[order] = append(orderedRulesMap[order], rule.Name.ValueString())
//		} else {
//			orderedRulesMap[order] = []string{rule.Name.ValueString()}
//		}
//	}
//
//	sortedKeys := make([]int, 0, len(orderedRulesMap))
//	for key := range orderedRulesMap {
//		sortedKeys = append(sortedKeys, key)
//	}
//	sort.Ints(sortedKeys)
//
//	ruleOrders := make(map[string]int)
//	lastOrderValue := -1
//	for _, key := range sortedKeys {
//		offset := 0
//		if lastOrderValue != -1 && lastOrderValue >= key {
//			offset = lastOrderValue - key + 1
//		}
//
//		for sliceIndex, ruleName := range orderedRulesMap[key] {
//			orderValue := key + sliceIndex + offset
//			ruleOrders[ruleName] = orderValue
//			lastOrderValue = orderValue
//		}
//	}
//
//	return ruleOrders
//}

// Get serverless runtime policy
func GetRuntimeServerlessPolicy(c api.PrismaCloudComputeAPIClient) (RuntimeServerlessPolicy, error) {
	var ans RuntimeServerlessPolicy
	if err := c.Request(http.MethodGet, RuntimeServerlessEndpoint, nil, nil, &ans); err != nil {
		return ans, fmt.Errorf("error getting serverless runtime policy: %s", err)
	}
	return ans, nil
}

// Get serverless runtime policy filtered on rule names
func GetRuntimeServerlessPolicyFiltered(c api.PrismaCloudComputeAPIClient, ruleNames []string) (RuntimeServerlessPolicy, error) {
    policy, err := GetRuntimeServerlessPolicy(c)
    if err != nil {
        return RuntimeServerlessPolicy{}, err
    }

    if (policy.Rules == nil || len(*policy.Rules) == 0) {
        return policy, nil
    }

    filteredRules := []RuntimeServerlessPolicyRule{}
    for _, rule := range *policy.Rules {
        if slices.Contains(ruleNames, rule.Name) {
            filteredRules = append(filteredRules, rule)
        }
    }

    response := RuntimeServerlessPolicy{
        Id: policy.Id,
        LearningDisabled: policy.LearningDisabled,
        Rules: &filteredRules,
    }

    return response, nil
}

// Update serverless runtime policy
func UpsertRuntimeServerlessPolicy(c api.PrismaCloudComputeAPIClient, policy RuntimeServerlessPolicy) error {
	return c.Request(http.MethodPut, RuntimeServerlessEndpoint, nil, policy, nil)
}

// Update host runtime policy, scoped to specified rule names
// Rules not managed through Terraform will be preserved
func UpsertRuntimeServerlessPolicyFiltered(c api.PrismaCloudComputeAPIClient, policy RuntimeServerlessPolicy, deletedRuleNames []string) error {
    currentPolicy, err := GetRuntimeServerlessPolicy(c)
    if err != nil {
        return err
    }

    if (currentPolicy.Rules == nil || len(*currentPolicy.Rules) == 0) {
	    return c.Request(http.MethodPut, RuntimeServerlessEndpoint, nil, policy, nil)
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

	return c.Request(http.MethodPut, RuntimeServerlessEndpoint, nil, policy, nil)
}
