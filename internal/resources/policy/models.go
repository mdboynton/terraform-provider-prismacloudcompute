package policy

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Generic compliance policy models
type CompliancePolicyResourceModel struct {
    Id          types.String                                `tfsdk:"id"`
    PolicyType  types.String                                `tfsdk:"policy_type"`
    Rules       *[]CompliancePolicyRuleResourceModel    `tfsdk:"rules"`
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

// Generic vulnerability policy models
type VulnerabilityPolicyResourceModel struct {
    Id              types.String                            `tfsdk:"id"`
    PolicyType      types.String                            `tfsdk:"policy_type"`
    PolicyContext   types.String                            `tfsdk:"policy_context"`
    Rules           *[]VulnerabilityPolicyRuleResourceModel `tfsdk:"rules"`
    Type            types.String                            `tfsdk:"type"`
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
