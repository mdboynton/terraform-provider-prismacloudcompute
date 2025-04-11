package models

import (
	"cmp"
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	SinglePortValueRegexp *regexp.Regexp = regexp.MustCompile(`^(\d{1,5})$`)
	RangePortValueRegxp   *regexp.Regexp = regexp.MustCompile(`^(\d{0,5})-(\d{1,5})$`)
)

//
// Structs
//

type RuntimePolicyDeniedProcessesResourceModel struct {
	Effect types.String `tfsdk:"effect"`
	Paths  types.List   `tfsdk:"paths"`
}

type RuntimePolicyCustomRuleResourceModel struct {
	ID     types.Int64  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	LogAs  types.String `tfsdk:"log_as"`
	Effect types.String `tfsdk:"effect"`
}

func (m RuntimePolicyCustomRuleResourceModel) ToTerraform(id *int) policyAPI.CustomRule {
	tfCustomRule := policyAPI.CustomRule{
		Action: m.LogAs.ValueString(),
		Effect: m.Effect.ValueString(),
	}

	if id == nil {
		tfCustomRule.ID = int(m.ID.ValueInt64())
	} else {
		tfCustomRule.ID = *id
	}

	return tfCustomRule
}

// Generic Policy Resource Models

type PolicyResourceModel struct {
	Id            types.String               `tfsdk:"id"`
	PolicyType    types.String               `tfsdk:"policy_type"`
	PolicyContext types.String               `tfsdk:"policy_context"`
	Rules         *[]PolicyRuleResourceModel `tfsdk:"rules"`
	Type          types.String               `tfsdk:"type"`
}

type PolicyRuleResourceModel struct {
	AlertThreshold                 *PolicyRuleThresholdResourceModel           `tfsdk:"alert_threshold" json:"alertThreshold"`
	BlockMessage                   types.String                                `tfsdk:"block_message"`
	BlockThreshold                 *PolicyRuleThresholdResourceModel           `tfsdk:"block_threshold" json:"blockThreshold"`
	Collections                    types.Set                                   `tfsdk:"collections"`
	CVERules                       *[]PolicyRuleExceptionResourceModel         `tfsdk:"cve_rules"`
	ComplianceActions              *PolicyRuleComplianceActionsResourceModel   `tfsdk:"compliance_actions"`
	Disabled                       types.Bool                                  `tfsdk:"disabled"`
	Effect                         types.String                                `tfsdk:"effect"`
	ExcludeBaseImageVulns          types.Bool                                  `tfsdk:"exclude_base_image_vulns"`
	GraceDays                      types.Int32                                 `tfsdk:"grace_days_all_severities"`
	GraceDaysPolicy                *PolicyRuleGraceDaysPolicyResourceModel     `tfsdk:"grace_days_by_severity"`
	Modified                       types.String                                `tfsdk:"modified"`
	Name                           types.String                                `tfsdk:"name"`
	Notes                          types.String                                `tfsdk:"notes"`
	OnlyFixed                      types.Bool                                  `tfsdk:"apply_only_when_fix_available"`
	Order                          types.Int32                                 `tfsdk:"order"`
	Owner                          types.String                                `tfsdk:"owner"`
	PkgTypesThresholds             *[]PolicyRulePkgTypesThresholdResourceModel `tfsdk:"package_types_thresholds"`
	ReportAllPassedAndFailedChecks types.Bool                                  `tfsdk:"report_passed_and_failed_checks"`
	Tags                           *[]PolicyRuleExceptionResourceModel         `tfsdk:"tags"`
	Verbose                        types.Bool                                  `tfsdk:"verbose"`
}

type PolicyRuleThresholdResourceModel struct {
	Threshold   types.String `tfsdk:"threshold"`
	RiskFactors types.Set    `tfsdk:"risk_factors"`
}

type PolicyRuleGraceDaysPolicyResourceModel struct {
	Low      types.Int32 `tfsdk:"low"`
	Medium   types.Int32 `tfsdk:"medium"`
	High     types.Int32 `tfsdk:"high"`
	Critical types.Int32 `tfsdk:"critical"`
}

type PolicyRulePkgTypesThresholdResourceModel struct {
	Type           types.String `tfsdk:"type"`
	AlertThreshold types.String `tfsdk:"alert_threshold"`
	BlockThreshold types.String `tfsdk:"block_threshold"`
}

type PolicyRuleExceptionResourceModel struct {
	//Name            types.String                                            `tfsdk:"name"`
	Effect      types.String `tfsdk:"effect"`
	Id          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	//Type            types.String                                            `tfsdk:"type"`
	Expiration PolicyRuleExceptionExpirationResourceModel `tfsdk:"expiration"`
}

type PolicyRuleExceptionExpirationResourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Date    types.String `tfsdk:"date"`
}

type PolicyRuleComplianceActionsResourceModel struct {
	Template   types.String                                   `tfsdk:"template"`
	Types      types.Set                                      `tfsdk:"types"`
	Severities types.Set                                      `tfsdk:"severities"`
	Checks     []PolicyRuleComplianceActionCheckResourceModel `tfsdk:"checks"`
}

type PolicyRuleComplianceActionCheckResourceModel struct {
	Id     types.Int32  `tfsdk:"id"`
	Action types.String `tfsdk:"action"`
}

//
// Sorting functions
//

func (m *PolicyResourceModel) GetRuleOrderMap() map[string]util.PolicyRuleOrderMapTuple {
	ruleOrderMap := map[string]util.PolicyRuleOrderMapTuple{}

	if m.Rules == nil || len(*m.Rules) == 0 {
		return ruleOrderMap
	}

	for idx, rule := range *m.Rules {
		ruleOrderMapTuple := util.PolicyRuleOrderMapTuple{
			Position: idx,
		}

		if rule.Order.IsNull() || rule.Order.IsUnknown() {
			ruleOrderMapTuple.Order = idx
		} else {
			ruleOrderMapTuple.Order = int(rule.Order.ValueInt32())
		}

		ruleOrderMap[rule.Name.ValueString()] = ruleOrderMapTuple
	}

	return ruleOrderMap
}

func SortPolicyRules(ctx context.Context, policyType string, policy interface{}, ruleOrderMap map[string]util.PolicyRuleOrderMapTuple, toTerraform bool) {
	util.LogDebug(ctx, "Executing SortPolicyRules()")

	if policy == nil {
		return
	}

	switch {
	case slices.Contains(policyAPI.OrderedGenericPolicyTypes, policyType):
		if toTerraform {
			typedPolicy, ok := (policy).(*policyAPI.Policy)
			if !ok {
				util.LogDebug(ctx, "conversion failed")
				return
			}

			SortGenericPolicyRulesTerraform(ctx, typedPolicy, ruleOrderMap)
		} else {
			typedPolicy, ok := (policy).(*PolicyResourceModel)
			if !ok {
				util.LogDebug(ctx, "conversion failed")
				return
			}

			SortGenericPolicyRulesSchema(ctx, typedPolicy, ruleOrderMap)
		}
	//case policyType == policyAPI.PolicyTypeRuntimeAppEmbedded:
	//    if toTerraform {
	//        typedPolicy, ok := (policy).(*policyAPI.RuntimeAppEmbeddedPolicy)
	//        if !ok {
	//            util.DLog(ctx, "conversion failed")
	//            return
	//        }

	//        SortGenericPolicyRulesTerraform(ctx, typedPolicy, ruleOrderMap)
	//    } else {
	//        typedPolicy, ok := (policy).(*RuntimeAppEmbeddedPolicyResourceModel)
	//        if !ok {
	//            util.DLog(ctx, "conversion failed")
	//            return
	//        }

	//        SortGenericPolicyRulesSchema(ctx, typedPolicy, ruleOrderMap)
	//    }
	default:
		// TODO: error
		return
	}

	util.LogDebug(ctx, "Finishing SortPolicyRules() execution")
}

// Change sorting of policy rules from the order returned by the API to the order specified in the plan
func SortGenericPolicyRulesSchema(ctx context.Context, policy *PolicyResourceModel, ruleOrderMap map[string]util.PolicyRuleOrderMapTuple) {
	util.LogDebug(ctx, "Executing SortGenericPolicyRulesSchema()")

	rules := *policy.Rules

	// Populate order values from map
	for i := 0; i < len(rules); i++ {
		rules[i].Order = types.Int32Value(int32(ruleOrderMap[rules[i].Name.ValueString()].Order))
	}

	// Sort rules back into their original position from the plan
	slices.SortFunc(rules, func(a, b PolicyRuleResourceModel) int {
		var orderA, orderB int
		tupleA, okA := ruleOrderMap[a.Name.ValueString()]
		if !okA {
			orderA = len(ruleOrderMap) + 1
		} else {
			orderA = tupleA.Position
		}

		tupleB, okB := ruleOrderMap[b.Name.ValueString()]
		if !okB {
			orderB = len(ruleOrderMap) + 1
		} else {
			orderB = tupleB.Position
		}

		return cmp.Compare(orderA, orderB)
	})

	// Assign sorted collection to policy object
	(*policy).Rules = &rules

	util.LogDebug(ctx, "Finishing SortGenericPolicyRulesSchema() execution")
}

// Sort policy rules according to their order values
func SortGenericPolicyRulesTerraform(ctx context.Context, policy *policyAPI.Policy, ruleOrderMap map[string]util.PolicyRuleOrderMapTuple) {
	util.LogDebug(ctx, "Executing SortGenericPolicyRulesTerraform()")

	rules := *policy.Rules

	// Sort rules into their correct position according to order values
	slices.SortFunc(rules, func(a, b policyAPI.PolicyRule) int {
		var orderA, orderB int
		tupleA, okA := ruleOrderMap[a.Name]
		if !okA {
			orderA = len(ruleOrderMap) + 1
		} else {
			orderA = tupleA.Order
		}

		tupleB, okB := ruleOrderMap[b.Name]
		if !okB {
			orderB = len(ruleOrderMap) + 1
		} else {
			orderB = tupleB.Order
		}
		return cmp.Compare(orderA, orderB)
	})

	// Assign sorted collection to policy object
	(*policy).Rules = &rules

	util.LogDebug(ctx, "Finishing SortGenericPolicyRulesTerraform() execution")
}

func PortRangesToTerraform(ctx context.Context, portRanges basetypes.ListValue) ([]policyAPI.PortRange, diag.Diagnostics) {
	var (
		diags        diag.Diagnostics
		values       []string              = []string{}
		tfPortRanges []policyAPI.PortRange = []policyAPI.PortRange{}
	)

	if portRanges.IsNull() {
		return tfPortRanges, diags
	}

	diags.Append(util.ListToStringSlice(ctx, &portRanges, &values)...)
	if diags.HasError() {
		return tfPortRanges, diags
	}

	for _, portRange := range values {
		start, end, diags := GetPortRangeValues(portRange)
		if diags.HasError() {
			continue
		}

		tfPortRanges = append(tfPortRanges, policyAPI.PortRange{
			Start: start,
			End:   end,
			Deny:  false,
		})
	}

	return tfPortRanges, diags
}

func GetPortRangeValues(portRange string) (int, int, diag.Diagnostics) {
	var diags diag.Diagnostics

	matchesSingleValue := SinglePortValueRegexp.MatchString(portRange)
	matchesRangeValue := RangePortValueRegxp.MatchString(portRange)

	if matchesRangeValue {
		splitRange := strings.SplitN(portRange, "-", 2)

		start, err := util.StringToInt(splitRange[0])
		if err != nil {
			diags.AddError(
				"Value Conversion Error",
				fmt.Sprintf("Error occurred while attempting to processes port range value \"%s\". Port ranges must follow the syntax \"N-N\", where ", portRange),
			)
			return 0, 0, diags
		}

		end, err := util.StringToInt(splitRange[1])
		if err != nil {
			diags.AddError("Value Conversion Error", err.Error())
			return 0, 0, diags
		}

		return start, end, diags
	} else if matchesSingleValue {
		value, err := util.StringToInt(portRange)
		if err != nil {
			diags.AddError("Value Conversion Error", err.Error())
			return 0, 0, diags
		}

		return value, value, diags
	} else {
		diags.AddError(
			"Invalid Resource Configuration",
			fmt.Sprintf("Invalid port value specified: %s\nValues must be strings containing either a non-negative integer or two non-negative integers separated by a hyphen. Integers must be between greater than or equal to 0 and less than or equal to 65535.", portRange),
		)
		return 0, 0, diags
	}
}

// Should be PortRanges (plural)
// TODO: rename?
func PortRangeToStringSlice(portRanges []policyAPI.PortRange) []string {
	response := []string{}

	for _, portRange := range portRanges {
		start := strconv.Itoa(portRange.Start)
		end := strconv.Itoa(portRange.End)

		if start == end {
			response = append(response, start)
		} else {
			response = append(response, fmt.Sprintf("%s-%s", start, end))
		}
	}

	return response
}

func CustomRuntimeRulesToTerraform(ctx context.Context, schemaCustomRules *[]RuntimePolicyCustomRuleResourceModel, customRuleIdMap map[string]int) ([]policyAPI.CustomRule, diag.Diagnostics) {
	// TODO: find a way to have the error message denote which rule has the error

	var diags diag.Diagnostics

	if schemaCustomRules == nil {
		return []policyAPI.CustomRule{}, diags
	}

	tfCustomRules := []policyAPI.CustomRule{}

	for _, schemaCustomRule := range *schemaCustomRules {
		customRuleId, ok := customRuleIdMap[schemaCustomRule.Name.ValueString()]
		if !ok {
			diags.AddError(
				"Value Conversion Error",
				fmt.Sprintf("No matching custom rule found for specified rule name \"%s\"", schemaCustomRule.Name.ValueString()),
			)

			return []policyAPI.CustomRule{}, diags
		}

		tfCustomRules = append(tfCustomRules, schemaCustomRule.ToTerraform(&customRuleId))
	}

	return tfCustomRules, diags
}

func CustomRuntimeRulesToSchema(ctx context.Context, tfCustomRuntimeRules []policyAPI.CustomRule, customRuleIdMap map[int]string) ([]RuntimePolicyCustomRuleResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	schemaCustomRuntimeRules := []RuntimePolicyCustomRuleResourceModel{}

	for _, tfCustomRuntimeRule := range tfCustomRuntimeRules {
		customRuleName, ok := customRuleIdMap[tfCustomRuntimeRule.ID]
		if !ok {
			diags.AddError(
				"Value Conversion Error",
				fmt.Sprintf("No matching custom runtime rule found for specified rule ID: %d", tfCustomRuntimeRule.ID),
			)

			return schemaCustomRuntimeRules, diags
		}

		schemaCustomRuntimeRules = append(schemaCustomRuntimeRules, RuntimePolicyCustomRuleResourceModel{
			ID:     types.Int64Value(int64(tfCustomRuntimeRule.ID)),
			Name:   types.StringValue(customRuleName),
			LogAs:  types.StringValue(tfCustomRuntimeRule.Action),
			Effect: types.StringValue(tfCustomRuntimeRule.Effect),
		})
	}

	return schemaCustomRuntimeRules, diags
}
