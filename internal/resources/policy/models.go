package policy

import (
    "context"
    //"sort"
    "slices"
    "cmp"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

//type OrderedPolicy interface {
//    //SortRules(rules *[]interface{}, planRules *[]interface{})
//    //CompliancePolicyRuleResourceModel | VulnerabilityPolicyRuleResourceModel
//    CompliancePolicyResourceModel | VulnerabilityPolicyResourceModel
//    GetRules() []interface{} 
//}

type BasePolicy interface {
    GetRules() *[]BasePolicyRule 
}

//type PolicyWithOrderedRules interface {
type OrderedPolicy interface {
    GetRules() *[]OrderedPolicyRule
}

type OrderedPolicyRule interface {
    BasePolicyRule
    GetOrder() int
    SetOrder(value int)
}

type BasePolicyRule interface {
    GetModified() string
    SetModified(value string)
    GetName() string
    SetName(value string)
    GetNotes() string
    SetNotes(value string)
    GetOwner() string
    SetOwner(value string)
}

type BaseComplianceVulnerabilityRule interface {
    BasePolicyRule
    GetBlockMessage() string
    SetBlockMessage(value string)
    GetCollections(ctx context.Context) ([]string, diag.Diagnostics)
    SetCollections(ctx context.Context, value []string) diag.Diagnostics
    GetDisabled() bool 
    SetDisabled(value bool)
    GetEffect() string
    SetEffect(value string)
    GetVerbose() bool
    SetVerbose(value bool)
    GetReportAllPassedAndFailedChecks() bool
    SetReportAllPassedAndFailedChecks(value bool)
}

type ComplianceRule interface {
    BaseComplianceVulnerabilityRule
    GetCondition(ctx context.Context) (*policyAPI.CompliancePolicyRuleCondition, diag.Diagnostics)
    SetCondition(ctx context.Context, value policyAPI.CompliancePolicyRuleCondition) diag.Diagnostics
}

type OrderedComplianceRule interface {
    OrderedPolicyRule
    BaseComplianceVulnerabilityRule
    GetCondition(ctx context.Context) (*policyAPI.CompliancePolicyRuleCondition, diag.Diagnostics)
    SetCondition(ctx context.Context, value policyAPI.CompliancePolicyRuleCondition) diag.Diagnostics
}

type VulnerabilityRule interface {
    BaseComplianceVulnerabilityRule
    GetAlertThreshold(ctx context.Context) (policyAPI.AlertThreshold, diag.Diagnostics)
    SetAlertThreshold(value policyAPI.AlertThreshold, riskFactors []string)
    GetBlockThreshold(ctx context.Context) (policyAPI.BlockThreshold, diag.Diagnostics)
    SetBlockThreshold(value policyAPI.BlockThreshold)
    GetCveRules() []policyAPI.Exception
    SetCveRules(value []policyAPI.Exception)
    GetExcludeBaseImageVulns() bool
    SetExcludeBaseImageVulns(value bool)
    GetGraceDays() int
    SetGraceDays(value int)
    GetGraceDaysPolicy() policyAPI.GraceDaysPolicy
    SetGraceDaysPolicy(value policyAPI.GraceDaysPolicy)
    GetOnlyFixed() bool
    SetOnlyFixed(value bool)
    GetPkgTypesThresholds() []policyAPI.PkgTypesThreshold
    SetPkgTypesThreshold(value []policyAPI.PkgTypesThreshold)
    GetTags() []policyAPI.Exception
    SetTags(value []policyAPI.Exception)
}

type OrderedVulnerabilityRule interface {
    OrderedPolicyRule
    BaseComplianceVulnerabilityRule
    GetAlertThreshold(ctx context.Context) (policyAPI.AlertThreshold, diag.Diagnostics)
    SetAlertThreshold(value policyAPI.AlertThreshold, riskFactors []string)
    GetBlockThreshold(ctx context.Context) (policyAPI.BlockThreshold, diag.Diagnostics)
    SetBlockThreshold(value policyAPI.BlockThreshold)
    GetCveRules() []policyAPI.Exception
    SetCveRules(value []policyAPI.Exception)
    GetExcludeBaseImageVulns() bool
    SetExcludeBaseImageVulns(value bool)
    GetGraceDays() int
    SetGraceDays(value int)
    GetGraceDaysPolicy() policyAPI.GraceDaysPolicy
    SetGraceDaysPolicy(value policyAPI.GraceDaysPolicy)
    GetOnlyFixed() bool
    SetOnlyFixed(value bool)
    GetPkgTypesThresholds() []policyAPI.PkgTypesThreshold
    SetPkgTypesThreshold(value []policyAPI.PkgTypesThreshold)
    GetTags() []policyAPI.Exception
    SetTags(value []policyAPI.Exception)
}

// Generic compliance policy models
type CompliancePolicyResourceModel struct {
    Id          types.String                                `tfsdk:"id"`
    PolicyType  types.String                                `tfsdk:"policy_type"`
    Rules       *[]CompliancePolicyRuleResourceModel    `tfsdk:"rules"`
    //Rules       *[]BasePolicyRule `tfsdk:"rules"`
    //Rules       *[]OrderedPolicyRule `tfsdk:"rules"`
    //Rules       *[]ComplianceRule `tfsdk:"rules"`
    //Rules       *[]OrderedComplianceRule `tfsdk:"rules"`
}

func (m CompliancePolicyResourceModel) GetRules() *[]CompliancePolicyRuleResourceModel {
//func (m *CompliancePolicyResourceModel) GetRules() *[]OrderedPolicyRule{
//func (m *CompliancePolicyResourceModel) GetRules() *[]BasePolicyRule{
//func (m *CompliancePolicyResourceModel) GetRules() *[]OrderedComplianceRule {
    return m.Rules
}

func (m CompliancePolicyResourceModel) SortSchemaRules(planRules *[]CompliancePolicyRuleResourceModel) {
//func (m *CompliancePolicyResourceModel) SortSchemaRules(planRules *[]OrderedPolicyRule) {
//func (m *CompliancePolicyResourceModel) SortSchemaRules(planRules *[]OrderedComplianceRule) {
    rules := m.GetRules()

    if planRules == nil {
        //for i := 0; i < len(*m.Rules); i++ {
        for i := 0; i < len(*rules); i++ {
            //(*m.Rules)[i].Order = types.Int32Value(int32(i + 1))
            (*rules)[i].SetOrder(i + 1)
        }
        return
    }

    ruleOrderMap := make(map[string]int32)
    for index, planRule := range *planRules {
        //ruleOrderMap[planRule.Name.ValueString()] = int32(index)
        ruleOrderMap[planRule.GetName()] = int32(index)
    }

    //slices.SortFunc(*m.Rules, func(a, b CompliancePolicyRuleResourceModel) int {
    //slices.SortFunc(*rules, func(a, b OrderedComplianceRule) int {
    slices.SortFunc(*rules, func(a, b CompliancePolicyRuleResourceModel) int {
        //orderA, okA := ruleOrderMap[a.Name.ValueString()]
        orderA, okA := ruleOrderMap[a.GetName()]
        if !okA {
            orderA = int32(len(ruleOrderMap) + 1)
        }
        //orderB, okB := ruleOrderMap[b.Name.ValueString()]
        orderB, okB := ruleOrderMap[b.GetName()]
        if !okB {
            orderB = int32(len(ruleOrderMap) + 1)
        }
        return cmp.Compare(orderA, orderB)
    })


    for i := 0; i < len(*m.Rules); i++ {
        //(*m.Rules)[i].Order = (*planRules)[i].Order
        (*rules)[i].SetOrder((*planRules)[i].GetOrder())
    }
}

type CompliancePolicyRuleResourceModel struct {
    BlockMessage                    types.String    `tfsdk:"block_message"`
    Collections                     types.Set `tfsdk:"collections"`
    Condition                       types.Object    `tfsdk:"condition"`
    Disabled                        types.Bool      `tfsdk:"disabled"`
    Effect                          types.String    `tfsdk:"effect"`
    Modified                        types.String    `tfsdk:"modified"`
    Name                            types.String    `tfsdk:"name"`
    Notes                           types.String    `tfsdk:"notes"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner                           types.String    `tfsdk:"owner"`
    ReportAllPassedAndFailedChecks  types.Bool      `tfsdk:"report_passed_and_failed_checks"`
    Verbose                         types.Bool      `tfsdk:"verbose"`
}


func (m *CompliancePolicyRuleResourceModel) GetBlockMessage() string {
    return m.BlockMessage.ValueString()
}

func (m *CompliancePolicyRuleResourceModel) SetBlockMessage(value string) {
    m.BlockMessage = types.StringValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetCollections(ctx context.Context) ([]string, diag.Diagnostics) {
    collections := []string{}
    diags := m.Collections.ElementsAs(ctx, &collections, false)
    return collections, diags
}

func (m *CompliancePolicyRuleResourceModel) SetCollections(ctx context.Context, value []string) diag.Diagnostics {
    collections := []attr.Value{}
    for _, val := range value {
        collections = append(collections, types.StringValue(val))
    }
    //m.Collections = types.SetValueMust(types.StringType, collections) 
    collectionsSet, diags := types.SetValue(types.StringType, collections) 
    if !diags.HasError() {
        m.Collections = collectionsSet 
    }
    return diags
}


func (m *CompliancePolicyRuleResourceModel) GetCondition(ctx context.Context) (*policyAPI.CompliancePolicyRuleCondition, diag.Diagnostics) {
    var diags diag.Diagnostics
    if m.Condition.IsUnknown() || m.Condition.IsNull() {
        return nil, diags
    }
    condition := policyAPI.CompliancePolicyRuleCondition{} 
    diags = m.Condition.As(ctx, &condition, basetypes.ObjectAsOptions{})
    return &condition, diags
}

func (m *CompliancePolicyRuleResourceModel) SetCondition(ctx context.Context, value policyAPI.CompliancePolicyRuleCondition) diag.Diagnostics {
    vulnerabilityObjectValues := []attr.Value{}
    for _, vulnerability := range(value.Vulnerabilities) {
        vulnerabilityObjectValue := types.ObjectValueMust(
            map[string]attr.Type{
                "id":        types.Int32Type,
                "block":       types.BoolType,
            },
            map[string]attr.Value{
                "id": types.Int32Value(int32(vulnerability.Id)),
                "block": types.BoolValue(vulnerability.Block),
            },
        )
       
        vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
    }

    vulnerabilityObject, diags := types.ListValueFrom(
        ctx,
        types.ObjectType{
            AttrTypes: map[string]attr.Type{
                "id": types.Int32Type,
                "block": types.BoolType,
            },
        },
        vulnerabilityObjectValues,
    )

    if diags.HasError() {
        return diags
    }

    conditionObject := types.ObjectValueMust(
        map[string]attr.Type{
            "vulnerabilities": types.ListType{
                ElemType: types.ObjectType{
                    AttrTypes: map[string]attr.Type{
                        "id": types.Int32Type,
                        "block": types.BoolType,
                    },
                },
            },
        },
        map[string]attr.Value{
            "vulnerabilities": vulnerabilityObject,
        },
    )

    m.Condition = conditionObject
   
    return diags
}

func (m *CompliancePolicyRuleResourceModel) GetDisabled() bool {
    return m.Disabled.ValueBool()
}

func (m *CompliancePolicyRuleResourceModel) SetDisabled(value bool) {
    m.Disabled = types.BoolValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetEffect() string {
    return m.Effect.ValueString()
}

func (m *CompliancePolicyRuleResourceModel) SetEffect(value string) {
    m.Effect = types.StringValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetModified() string {
    return m.Modified.ValueString()
}

func (m *CompliancePolicyRuleResourceModel) SetModified(value string) {
    m.Modified = types.StringValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetName() string {
    return m.Name.ValueString()
}

func (m *CompliancePolicyRuleResourceModel) SetName(value string) {
    m.Name = types.StringValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetNotes() string {
    return m.Notes.ValueString()
}

func (m *CompliancePolicyRuleResourceModel) SetNotes(value string) {
    m.Notes = types.StringValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetOrder() int {
    return int(m.Order.ValueInt32())
}

func (m *CompliancePolicyRuleResourceModel) SetOrder(value int) {
    m.Order = types.Int32Value(int32(value))
}

func (m *CompliancePolicyRuleResourceModel) GetOwner() string {
    return m.Owner.ValueString()
}

func (m *CompliancePolicyRuleResourceModel) SetOwner(value string) {
    m.Owner = types.StringValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetReportAllPassedAndFailedChecks() bool {
    return m.ReportAllPassedAndFailedChecks.ValueBool()
}

func (m *CompliancePolicyRuleResourceModel) SetReportAllPassedAndFailedChecks(value bool) {
    m.ReportAllPassedAndFailedChecks = types.BoolValue(value)
}

func (m *CompliancePolicyRuleResourceModel) GetVerbose() bool {
    return m.Verbose.ValueBool()
}

func (m *CompliancePolicyRuleResourceModel) SetVerbose(value bool) {
    m.Verbose = types.BoolValue(value)
}

// Generic vulnerability policy models
type VulnerabilityPolicyResourceModel struct {
    Id              types.String                            `tfsdk:"id"`
    PolicyType      types.String                            `tfsdk:"policy_type"`
    PolicyContext   types.String                            `tfsdk:"policy_context"`
    Rules           *[]VulnerabilityPolicyRuleResourceModel `tfsdk:"rules"`
    //Rules           *[]OrderedPolicyRule                    `tfsdk:"rules"`
    //Rules           *[]VulnerabilityRule                    `tfsdk:"rules"`
    //Rules           *[]OrderedVulnerabilityRule                    `tfsdk:"rules"`
    Type            types.String                            `tfsdk:"type"`
}

//func (m *VulnerabilityPolicyResourceModel) GetRules() *[]OrderedPolicyRule {
//func (m *VulnerabilityPolicyResourceModel) GetRules() *[]VulnerabilityPolicyRuleResourceModel {
//func (m *VulnerabilityPolicyResourceModel) GetRules() *[]BasePolicyRule {
//func (m *VulnerabilityPolicyResourceModel) GetRules() *[]VulnerabilityRule {
//func (m *VulnerabilityPolicyResourceModel) GetRules() *[]OrderedVulnerabilityRule {
func (m *VulnerabilityPolicyResourceModel) GetRules() *[]VulnerabilityPolicyRuleResourceModel {
    return m.Rules 
}

type VulnerabilityPolicyRuleResourceModel struct {
    AlertThreshold                  *VulnerabilityPolicyRuleThresholdResourceModel `tfsdk:"alert_threshold" json:"alertThreshold"`
    //AlertThreshold                  *VulnerabilityPolicyRuleAlertThresholdResourceModel `tfsdk:"alert_threshold" json:"alertThreshold"`
    //AlertThreshold                  types.Object `tfsdk:"alert_threshold" json:"alertThreshold"`
    BlockMessage                    types.String    `tfsdk:"block_message"`
    BlockThreshold                  *VulnerabilityPolicyRuleThresholdResourceModel `tfsdk:"block_threshold" json:"blockThreshold"`
    //BlockThreshold                  *VulnerabilityPolicyRuleBlockThresholdResourceModel `tfsdk:"block_threshold" json:"blockThreshold"`
    //BlockThreshold                  types.Object `tfsdk:"block_threshold" json:"blockThreshold"`
    Collections                     types.Set `tfsdk:"collections"`
    CVERules                        *[]VulnerabilityPolicyRuleExceptionResourceModel  `tfsdk:"cve_rules"`
    //Condition                       types.Object    `tfsdk:"condition"`
    Disabled                        types.Bool      `tfsdk:"disabled"`
    Effect                          types.String    `tfsdk:"effect"`
    ExcludeBaseImageVulns           types.Bool      `tfsdk:"exclude_base_image_vulns"`
    GraceDays                       types.Int32     `tfsdk:"grace_days_all_severities"`
    GraceDaysPolicy                 *VulnerabilityPolicyRuleGraceDaysPolicyResourceModel `tfsdk:"grace_days_by_severity"` 
    Modified                        types.String    `tfsdk:"modified"`
    Name                            types.String    `tfsdk:"name"`
    Notes                           types.String    `tfsdk:"notes"`
    OnlyFixed                       types.Bool      `tfsdk:"apply_only_when_fix_available"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner                           types.String    `tfsdk:"owner"`
    PkgTypesThresholds              *[]VulnerabilityPolicyRulePkgTypesThresholdResourceModel    `tfsdk:"package_types_thresholds"`
    ReportAllPassedAndFailedChecks  types.Bool      `tfsdk:"report_passed_and_failed_checks"`
    Tags                            *[]VulnerabilityPolicyRuleExceptionResourceModel  `tfsdk:"tags"`
    Verbose                         types.Bool      `tfsdk:"verbose"`
}

func (m *VulnerabilityPolicyRuleResourceModel) GetAlertThreshold(ctx context.Context) (policyAPI.AlertThreshold, diag.Diagnostics) {
    return alertThresholdToTerraform(ctx, m.AlertThreshold)
}

func (m *VulnerabilityPolicyRuleResourceModel) SetAlertThreshold(ctx context.Context, value policyAPI.AlertThreshold, riskFactors []string) diag.Diagnostics {
    alertThreshold, diags := alertThresholdToSchema(ctx, &value, riskFactors)
    if !diags.HasError() {
        m.AlertThreshold = &alertThreshold
    }
    return diags
}

func (m *VulnerabilityPolicyRuleResourceModel) GetBlockMessage() string {
    return m.BlockMessage.ValueString()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetBlockMessage(value string) {
    m.BlockMessage = types.StringValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetBlockThreshold(ctx context.Context) (policyAPI.BlockThreshold, diag.Diagnostics) {
    return blockThresholdToTerraform(ctx, m.BlockThreshold)
}

func (m *VulnerabilityPolicyRuleResourceModel) SetBlockThreshold(ctx context.Context, value policyAPI.BlockThreshold, riskFactors []string) diag.Diagnostics {
    blockThreshold, diags := blockThresholdToSchema(ctx, &value, riskFactors)
    if !diags.HasError() {
        m.BlockThreshold = &blockThreshold
    }
    return diags
}

func (m *VulnerabilityPolicyRuleResourceModel) GetCollections(ctx context.Context) ([]string, diag.Diagnostics) {
    collections := []string{}
    diags := m.Collections.ElementsAs(ctx, &collections, false)
    return collections, diags
}

func (m *VulnerabilityPolicyRuleResourceModel) SetCollections(ctx context.Context, value []string) diag.Diagnostics {
    collections := []attr.Value{}
    for _, val := range value {
        collections = append(collections, types.StringValue(val))
    }
    //m.Collections = types.SetValueMust(types.StringType, collections) 
    collectionsSet, diags := types.SetValue(types.StringType, collections) 
    if !diags.HasError() {
        m.Collections = collectionsSet 
    }
    return diags
}

func (m *VulnerabilityPolicyRuleResourceModel) GetCveRules(ctx context.Context) ([]policyAPI.Exception, diag.Diagnostics) {
    return exceptionsToTerraform(ctx, *m.CVERules, false)
}

func (m *VulnerabilityPolicyRuleResourceModel) SetCveRules(ctx context.Context, value []policyAPI.Exception) diag.Diagnostics {
    cveRules, diags := exceptionsToSchema(ctx, value, false)
    if !diags.HasError() {
        m.CVERules = &cveRules
    }
    return diags
}

func (m *VulnerabilityPolicyRuleResourceModel) GetDisabled() bool {
    return m.Disabled.ValueBool()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetDisabled(value bool) {
    m.Disabled = types.BoolValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetEffect() string {
    return m.Effect.ValueString()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetEffect(value string) {
    m.Effect = types.StringValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetExcludeBaseImageVulns() bool {
    return m.ExcludeBaseImageVulns.ValueBool()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetExcludeBaseImageVulns(value bool) {
    m.ExcludeBaseImageVulns = types.BoolValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetGraceDays() int {
    return int(m.GraceDays.ValueInt32())
}

func (m *VulnerabilityPolicyRuleResourceModel) SetGraceDays(value int) {
    m.GraceDays = types.Int32Value(int32(value))
}

func (m *VulnerabilityPolicyRuleResourceModel) GetGraceDaysPolicy(ctx context.Context) (policyAPI.GraceDaysPolicy, diag.Diagnostics) {
    return graceDaysPolicyToTerraform(ctx, *m.GraceDaysPolicy)
}

func (m *VulnerabilityPolicyRuleResourceModel) SetGraceDaysPolicy(ctx context.Context, value policyAPI.GraceDaysPolicy) diag.Diagnostics {
    graceDaysPolicy, diags := graceDaysPolicyToSchema(ctx, value)
    if !diags.HasError() {
        m.GraceDaysPolicy = &graceDaysPolicy
    }
    return diags
}

func (m *VulnerabilityPolicyRuleResourceModel) GetModified() string {
    return m.Modified.ValueString()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetModified(value string) {
    m.Modified = types.StringValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetName() string {
    return m.Name.ValueString()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetName(value string) {
    m.Name = types.StringValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetNotes() string {
    return m.Notes.ValueString()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetNotes(value string) {
    m.Name = types.StringValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetOnlyFixed() bool {
    return m.OnlyFixed.ValueBool()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetOnlyFixed(value bool) {
    m.OnlyFixed = types.BoolValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetOrder() int {
    return int(m.Order.ValueInt32())
}

func (m *VulnerabilityPolicyRuleResourceModel) SetOrder(value int) {
    m.Order = types.Int32Value(int32(value))
}

func (m *VulnerabilityPolicyRuleResourceModel) GetOwner() string {
    return m.Owner.ValueString()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetOwner(value string) {
    m.Owner = types.StringValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetPkgTypesThresholds(ctx context.Context) ([]policyAPI.PkgTypesThreshold, diag.Diagnostics) {
    return pkgTypesThresholdsToTerraform(ctx, *m.PkgTypesThresholds)
}

func (m *VulnerabilityPolicyRuleResourceModel) SetPkgTypesThresholds(ctx context.Context, value []policyAPI.PkgTypesThreshold) diag.Diagnostics {
    pkgTypesThresholds, diags := pkgTypesThresholdsToSchema(ctx, value)
    if !diags.HasError() {
        m.PkgTypesThresholds = &pkgTypesThresholds
    }
    return diags
}

func (m *VulnerabilityPolicyRuleResourceModel) GetReportAllPassedAndFailedChecks() bool {
    return m.ReportAllPassedAndFailedChecks.ValueBool()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetReportAllPassedAndFailedChecks(value bool) {
    m.ReportAllPassedAndFailedChecks = types.BoolValue(value)
}

func (m *VulnerabilityPolicyRuleResourceModel) GetTags(ctx context.Context) ([]policyAPI.Exception, diag.Diagnostics) {
    return exceptionsToTerraform(ctx, *m.Tags, true)
}

func (m *VulnerabilityPolicyRuleResourceModel) SetTags(ctx context.Context, value []policyAPI.Exception) diag.Diagnostics {
    tags, diags := exceptionsToSchema(ctx, value, true)
    if !diags.HasError() {
        m.Tags = &tags
    }
    return diags
}

func (m *VulnerabilityPolicyRuleResourceModel) GetVerbose() bool {
    return m.Verbose.ValueBool()
}

func (m *VulnerabilityPolicyRuleResourceModel) SetVerbose(value bool) {
    m.Verbose = types.BoolValue(value)
}


type VulnerabilityPolicyRuleThresholdResourceModel struct {
    Threshold   types.String    `tfsdk:"threshold"`
    RiskFactors types.Set       `tfsdk:"risk_factors"`
}

type VulnerabilityPolicyRuleGraceDaysPolicyResourceModel struct {
    Low         types.Int32     `tfsdk:"low"`
    Medium      types.Int32     `tfsdk:"medium"`
    High        types.Int32     `tfsdk:"high"`
    Critical    types.Int32     `tfsdk:"critical"`
}

type VulnerabilityPolicyRulePkgTypesThresholdResourceModel struct {
    Type            types.String    `tfsdk:"type"`
    AlertThreshold  types.String    `tfsdk:"alert_threshold"`
    BlockThreshold  types.String    `tfsdk:"block_threshold"`
}

type VulnerabilityPolicyRuleExceptionResourceModel struct {
    Name            types.String                                            `tfsdk:"name"`
    Effect          types.String                                            `tfsdk:"effect"`
    Id              types.String                                            `tfsdk:"id"`
    Description     types.String                                            `tfsdk:"description"`
    Type            types.String                                            `tfsdk:"type"`
    Expiration      VulnerabilityPolicyRuleExceptionExpirationResourceModel `tfsdk:"expiration"`
}

type VulnerabilityPolicyRuleExceptionExpirationResourceModel struct {
    Enabled types.Bool      `tfsdk:"enabled"`
    Date    types.String    `tfsdk:"date"`
}

// Application control policy models
type ApplicationControlPolicyResourceModel struct {
    Rules       *[]ApplicationControlPolicyRuleResourceModel    `tfsdk:"rules"`
}

type ApplicationControlPolicyRuleResourceModel struct {
    Id              types.Int32     `tfsdk:"id"`
    Applications    types.Set       `tfsdk:"applications"`
    Description     types.String    `tfsdk:"description"`
    //Disabled        types.Bool      `tfsdk:"disabled"` // TODO: is this even supported? doesnt appear in UI
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Notes           types.String    `tfsdk:"notes"` // TODO: is this even supported? doesnt appear in UI
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
    Severity        types.String    `tfsdk:"severity"`
}

type ApplicationControlPolicyRuleAppliactionResourceModel struct {
    Name            types.String    `tfsdk:"name"`
    AllowedVersions types.List      `tfsdk:"allowed_versions"`
}

// Trusted images policy models
type TrustedImagesPolicyResourceModel struct {
    Groups      *[]TrustGroupResourceModel           `tfsdk:"groups"`
    Policy      TrustedImagesPolicyRulesResourceModel    `tfsdk:"policy"`
}

type TrustGroupResourceModel struct {
    Id              types.String    `tfsdk:"id"`
    Images          types.Set       `tfsdk:"images"`
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
}

type TrustedImagesPolicyRulesResourceModel struct {
    Id      types.String                            `tfsdk:"id"`
    Enabled types.Bool                              `tfsdk:"enabled"`
    Rules   *[]TrustedImagesPolicyRuleResourceModel `tfsdk:"rules"`
}

type TrustedImagesPolicyRuleResourceModel struct {
    //Action          types.Set       `tfsdk:"action"`
    AllowedGroups   types.Set       `tfsdk:"allowed_groups"`
    Collections     types.Set       `tfsdk:"collections"`
    DeniedGroups    types.Set       `tfsdk:"denied_groups"`
    Effect          types.String    `tfsdk:"effect"`
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Notes           types.String    `tfsdk:"notes"`
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
}
