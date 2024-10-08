package policy

import (
    "context"
	//"fmt"
    //"reflect"
    //"slices"
    //"cmp"
    //"sort"
    //"time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"
    //"github.com/hashicorp/terraform-plugin-log/tflog"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
    //"github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ resource.Resource = &HostCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &HostCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &HostCompliancePolicyResource{}

func NewHostCompliancePolicyResource() resource.Resource {
    return &HostCompliancePolicyResource{}
}

func (r *HostCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_host_compliance_policy"
}

type HostCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

type HostCompliancePolicyResourceModel struct {
    Id types.String `tfsdk:"id"`
    PolicyType types.String `tfsdk:"policy_type"`
    Rules *[]HostCompliancePolicyRuleResourceModel `tfsdk:"rules"`
}

type HostCompliancePolicyRuleResourceModel struct {
    Order types.Int32 `tfsdk:"order"`
    //Action types.Set `tfsdk:"action"`
    //Modified types.String `tfsdk:"modified"`
    Name types.String `tfsdk:"name"`
    //Collections types.Set `tfsdk:"collections"`
    Collections types.List `tfsdk:"collections"`
    //Action types.Set `tfsdk:"action"`
    ////AlertThreshold *HostCompliancePolicyRuleAlertThresholdResourceModel `tfsdk:"alert_threshold"`
    ReportAllPassedAndFailedChecks types.Bool `tfsdk:"report_passed_and_failed_checks"`
    ////AuditAllowed types.Bool `tfsdk:"audit_allowed"`
    BlockMessage types.String `tfsdk:"block_message"`
    ////BlockThreshold  *HostCompliancePolicyRuleBlockThresholdResourceModel `tfsdk:"block_threshold"`
    //Condition *HostCompliancePolicyRuleConditionResourceModel `tfsdk:"condition"`
    Condition types.Object `tfsdk:"condition"`
    //CreatePR types.Bool `tfsdk:"create_pr"
    //CVERules *[]HostCompliancePolicyRuleCVERuleResourceModel `tfsdk:"cve_rules"`
    Disabled types.Bool `tfsdk:"disabled"`
    Effect types.String `tfsdk:"effect"`
    ////ExcludeBaseImageVulns types.Bool `tfsdk:"exclude_base_image_vulns"`
    //GraceDays types.Int32 `tfsdk:"grace_days"`
    //GraceDaysPolicy *HostCompliancePolicyRuleGraceDaysPolicyResourceModel `tfsdk:"grace_days_policy"`
    ////GraceDaysPolicy types.Object `tfsdk:"grace_days_policy"`
    //Groups types.Set `tfsdk:"groups"`
    //License *HostCompliancePolicyRuleLicenseResourceModel `tfsdk:"license"`
    //License types.Object `tfsdk:"license"`
    Modified types.String `tfsdk:"modified"`
    //Name types.String `tfsdk:"name"`
    Notes types.String `tfsdk:"notes"`
    //OnlyFixed types.Bool `tfsdk:"only_fixed"`
    Owner types.String `tfsdk:"owner"`
    //PkgTypesThresholds *[]HostCompliancePolicyRulePkgTypesThresholdsResourceModel `tfsdk:"package_types_thresholds"`
    ////PreviousName types.String `tfsdk:"previous_name"`
    //Principal types.Set `tfsdk:"principal"`
    //RiskFactorsEffects *[]HostCompliancePolicyRuleRiskFactorsEffectsResourceModel `tfsdk:"risk_factors_effects"`
    //Tags *[]HostCompliancePolicyRuleTagsResourceModel `tfsdk:"tags"`
    Verbose types.Bool `tfsdk:"verbose"`
}

// Alert Threshold
type HostCompliancePolicyRuleAlertThresholdResourceModel struct {
    Disabled types.Bool `tfsdk:"disabled"`
    Value types.Int32 `tfsdk:"value"`
}

// Block Threshold
type HostCompliancePolicyRuleBlockThresholdResourceModel struct {
    Enabled types.Bool `tfsdk:"enabled"`
    Value types.Int32 `tfsdk:"value"`
}

// Condition
type HostCompliancePolicyRuleConditionResourceModel struct {
    Device types.String `tfsdk:"device"`
    ReadOnly types.Bool `tfsdk:"read_only"`
    Vulnerabilities []HostCompliancePolicyRuleVulnerabilityResourceModel `tfsdk:"vulnerabilities"`
}

type HostCompliancePolicyRuleVulnerabilityResourceModel struct {
    Block types.Bool `tfsdk:"block"`
    Id types.Int32 `tfsdk:"id"`
}

// CVE Rule
type HostCompliancePolicyRuleCVERuleResourceModel struct {
    Description types.String `tfsdk:"description"`
    Effect types.String `tfsdk:"effect"`
    Expiration HostCompliancePolicyRuleCVERuleExpirationResourceModel `tfsdk:"expiration"`
    Id types.String `tfsdk:"id"`
}

type HostCompliancePolicyRuleCVERuleExpirationResourceModel struct {
    Date types.String `tfsdk:"date"`
    Enabled types.Bool `tfsdk:"enabled"`
}

// Grace Days Policy
type HostCompliancePolicyRuleGraceDaysPolicyResourceModel struct {
    Critical types.Int32 `tfsdk:"critical"`
    Enabled types.Bool `tfsdk:"enabled"`
    High types.Int32 `tfsdk:"high"`
    Low types.Int32 `tfsdk:"low"`
    Medium types.Int32 `tfsdk:"medium"`
}

// License
type HostCompliancePolicyRuleLicenseResourceModel struct {
    AlertThreshold HostCompliancePolicyRuleAlertThresholdResourceModel `tfsdk:"alert_threshold"`
    BlockThreshold HostCompliancePolicyRuleBlockThresholdResourceModel `tfsdk:"block_threshold"`
    Critical types.Set `tfsdk:"critical"`
    High types.Set `tfsdk:"high"`
    Low types.Set `tfsdk:"low"`
    Medium types.Set `tfsdk:"medium"`
}

// Package Types Thresholds
type HostCompliancePolicyRulePkgTypesThresholdsResourceModel struct {
    AlertThreshold HostCompliancePolicyRuleAlertThresholdResourceModel `tfsdk:"alert_threshold"`
    BlockThreshold HostCompliancePolicyRuleBlockThresholdResourceModel `tfsdk:"block_threshold"`
    Type types.String `tfsdk:"type"` 
}

// Risk Factors Effects
type HostCompliancePolicyRuleRiskFactorsEffectsResourceModel struct {
    Effect types.String `tfsdk:"effect"`
    RiskFactor types.String `tfsdk:"risk_factor"`
}

// Tags
type HostCompliancePolicyRuleTagsResourceModel struct {
    Description types.String `tfsdk:"description"`
    Effect types.String `tfsdk:"effect"`
    Expiration HostCompliancePolicyRuleExpirationResourceModel `tfsdk:"expiration"`
    Name types.String `tfsdk:"name"`
}

type HostCompliancePolicyRuleExpirationResourceModel struct {
    Date types.String `tfsdk:"date"`
    Enabled types.Bool `tfsdk:"enabled"`
}

func (r *HostCompliancePolicyResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("hostCompliance"),
            },
            "policy_type": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("hostCompliance"),
            },
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                //Computed: true,
                Validators: []validator.List{
                    validators.PolicyRuleNameIsUnique("host compliance"),
                },
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "order": schema.Int32Attribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                        },
                        //"action": schema.SetAttribute{
                        //    MarkdownDescription: "TODO",
                        //    ElementType: types.StringType,
                        //    Optional: true,
                        //    Computed: true,
                        //    //Default: setdefault.StaticValue(
                        //    //    types.SetValueMust(
                        //    //        types.StringType,
                        //    //        []attr.Value{
                        //    //            types.StringValue("*"),
                        //    //        },
                        //    //    ),
                        //    //),
                        //},
                        ////"alert_threshold": schema.SingleNestedAttribute{
                        ////    MarkdownDescription: "TODO",
                        ////    Optional: true,
                        ////    Computed: true,
                        ////    Attributes: map[string]schema.Attribute{
                        ////        "disabled": schema.BoolAttribute{
                        ////            MarkdownDescription: "TODO",
                        ////            Optional: true,
                        ////            Computed: true,
                        ////        },
                        ////        "value": schema.Int32Attribute{
                        ////            MarkdownDescription: "TODO",
                        ////            Optional: true,
                        ////            Computed: true,
                        ////        },
                        ////    },
                        ////    //Default: objectdefault.StaticValue(
                        ////    //    types.ObjectValueMust(
                        ////    //        map[string]attr.Type{
                        ////    //            "disabled": types.BoolType,
                        ////    //            "value": types.Int32Type,
                        ////    //        },
                        ////    //        map[string]attr.Value{
                        ////    //            "disabled": types.BoolValue(false),
                        ////    //            "value": types.Int32Value(int32(0)),
                        ////    //        },
                        ////    //    ),
                        ////    //),
                        ////},
                        ////"block_threshold": schema.SingleNestedAttribute{
                        ////    MarkdownDescription: "TODO",
                        ////    Optional: true,
                        ////    Computed: true,
                        ////    Attributes: map[string]schema.Attribute{
                        ////        "enabled": schema.BoolAttribute{
                        ////            MarkdownDescription: "TODO",
                        ////            Optional: true,
                        ////            Computed: true,
                        ////        },
                        ////        "value": schema.Int32Attribute{
                        ////            MarkdownDescription: "TODO",
                        ////            Optional: true,
                        ////            Computed: true,
                        ////        },
                        ////    },
                        ////    //Default: objectdefault.StaticValue(
                        ////    //    types.ObjectValueMust(
                        ////    //        map[string]attr.Type{
                        ////    //            "enabled": types.BoolType,
                        ////    //            "value": types.Int32Type,
                        ////    //        },
                        ////    //        map[string]attr.Value{
                        ////    //            "enabled": types.BoolValue(false),
                        ////    //            "value": types.Int32Value(int32(0)),
                        ////    //        },
                        ////    //    ),
                        ////    //),
                        ////},
                        "report_passed_and_failed_checks": schema.BoolAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            Default: booldefault.StaticBool(false), 
                        },
                        "block_message": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            Default: stringdefault.StaticString(""),
                        },
                        "collections": r.GetCollectionsSchema(),
                        "condition": schema.SingleNestedAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            Attributes: map[string]schema.Attribute{
                                //"device": schema.StringAttribute{
                                //    MarkdownDescription: "TODO",
                                //    Optional: true,
                                //    Computed: true,
                                //    //Default: stringdefault.StaticString(""),
                                //},
                                //"read_only": schema.BoolAttribute{
                                //    MarkdownDescription: "TODO",
                                //    Optional: true,
                                //    Computed: true,
                                //    //Default: booldefault.StaticBool(false),
                                //},
                                "vulnerabilities": schema.ListNestedAttribute{
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
                                    NestedObject: schema.NestedAttributeObject{
                                        Attributes: map[string]schema.Attribute{
                                            "id": schema.Int32Attribute{
                                                MarkdownDescription: "TODO",
                                                Optional: true,
                                                Computed: true,
                                            },
                                            "block": schema.BoolAttribute{
                                                MarkdownDescription: "TODO",
                                                Optional: true,
                                                Computed: true,
                                            },
                                        },
                                    },
                                    // TODO: default value
                                },
                            },
                        },
                        //"create_pr": schema.BoolAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: booldefault.StaticBool(false),
                        //},
                        //"cve_rules": schema.SetNestedAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    PlanModifiers: []planmodifier.Set{
                        //        setplanmodifier.UseStateForUnknown(),
                        //    },
                        //    NestedObject: schema.NestedAttributeObject{
                        //        Attributes: map[string]schema.Attribute{
                        //            "description": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //            },
                        //            "effect": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Validators: []validator.String{
                        //                    stringvalidator.OneOf("alert", "block", "ignore"),
                        //                },
                        //            },
                        //            "expiration": schema.SingleNestedAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Attributes: map[string]schema.Attribute{
                        //                    "date": schema.StringAttribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                    "enabled": schema.BoolAttribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                },
                        //            },
                        //        },
                        //    },
                        //},
                        "disabled": schema.BoolAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            Default: booldefault.StaticBool(false), 
                        },
                        "effect": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            Validators: []validator.String{
                                stringvalidator.OneOf("ignore", "alert", "block", "alert, block"),
                            },
                        },
                        ////"exclude_base_image_vulns": schema.BoolAttribute{
                        ////    MarkdownDescription: "TODO",
                        ////    Optional: true,
                        ////    Computed: true,
                        ////    // TODO: do we need a default value? 
                        ////},
                        //"grace_days": schema.Int32Attribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: int32default.StaticInt32(0), 
                        //},
                        //"grace_days_policy": schema.SingleNestedAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    PlanModifiers: []planmodifier.Object{
                        //        objectplanmodifier.UseStateForUnknown(),
                        //    },
                        //    Attributes: map[string]schema.Attribute{
                        //        "enabled": schema.BoolAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //        },
                        //        "low": schema.Int32Attribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //        },
                        //        "medium": schema.Int32Attribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //        },
                        //        "high": schema.Int32Attribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //        },
                        //        "critical": schema.Int32Attribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //        },
                        //    },
                        //    //PlanModifiers: UseStateForUnknown(),
                        //    //PlanModifiers: []planmodifier.Object{
                        //    //    UseStateForUnknown(),
                        //    //},
                        //    //PlanModifiers: []planmodifier.Object{
                        //        //objectdefault.StaticValue(
                        //        //    types.ObjectValueFrom(ctx, map[string]attr.Type{}, nil),
                        //        //    //types.ObjectValue(
                        //        //    //    map[string]attr.Type{}, 
                        //        //    //    map[string]attr.Value{},
                        //        //    //),
                        //        //),
                        //        //objectplanmodifier.DefaultValue({}),
                        //    //},
                        //    //Default: objectdefault.StaticValue(
                        //    //    types.ObjectValueMust(
                        //    //        map[string]attr.Type{
                        //    //    //        //"disabled": types.BoolType,
                        //    //    //        //"value": types.Int32Type,
                        //    //        },
                        //    //        map[string]attr.Value{
                        //    //    //        //"disabled": types.BoolValue(false),
                        //    //    //        //"value": types.Int32Value(int32(0)),
                        //    //        },
                        //    //    ),
                        //    //),
                        //},
                        //"groups": schema.SetAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    ElementType: types.StringType,
                        //    PlanModifiers: []planmodifier.Set{
                        //        setplanmodifier.UseStateForUnknown(),
                        //    },
                        //    Default: setdefault.StaticValue(
                        //        types.SetValueMust(
                        //            types.StringType,
                        //            []attr.Value{
                        //                types.StringValue("*"),
                        //            },
                        //        ),
                        //    ),
                        //},
                        //"license": schema.SingleNestedAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Attributes: map[string]schema.Attribute{
                        //        "alert_threshold": schema.SingleNestedAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            Attributes: map[string]schema.Attribute{
                        //                "disabled": schema.BoolAttribute{
                        //                    MarkdownDescription: "TODO",
                        //                    Optional: true,
                        //                    Computed: true,
                        //                },
                        //                "value": schema.Int32Attribute{
                        //                    MarkdownDescription: "TODO",
                        //                    Optional: true,
                        //                    Computed: true,
                        //                },
                        //            },
                        //            //Default: objectdefault.StaticValue(
                        //            //    types.ObjectValueMust(
                        //            //        map[string]attr.Type{
                        //            //            "disabled": types.BoolType,
                        //            //            "value": types.Int32Type,
                        //            //        },
                        //            //        map[string]attr.Value{
                        //            //            "disabled": types.BoolValue(false),
                        //            //            "value": types.Int32Value(int32(0)),
                        //            //        },
                        //            //    ),
                        //            //),
                        //        },
                        //        "block_threshold": schema.SingleNestedAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            Attributes: map[string]schema.Attribute{
                        //                "enabled": schema.BoolAttribute{
                        //                    MarkdownDescription: "TODO",
                        //                    Optional: true,
                        //                    Computed: true,
                        //                },
                        //                "value": schema.Int32Attribute{
                        //                    MarkdownDescription: "TODO",
                        //                    Optional: true,
                        //                    Computed: true,
                        //                },
                        //            },
                        //            //Default: objectdefault.StaticValue(
                        //            //    types.ObjectValueMust(
                        //            //        map[string]attr.Type{
                        //            //            "enabled": types.BoolType,
                        //            //            "value": types.Int32Type,
                        //            //        },
                        //            //        map[string]attr.Value{
                        //            //            "enabled": types.BoolValue(false),
                        //            //            "value": types.Int32Value(int32(0)),
                        //            //        },
                        //            //    ),
                        //            //),
                        //        },
                        //        "low": schema.SetAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            ElementType: types.StringType,
                        //        },
                        //        "medium": schema.SetAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            ElementType: types.StringType,
                        //        },
                        //        "high": schema.SetAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            ElementType: types.StringType,
                        //        },
                        //        "critical": schema.SetAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            ElementType: types.StringType,
                        //        },
                        //    },
                        //    //PlanModifiers: []planmodifier.Object{
                        //    //    UseStateForUnknown(),
                        //    //},
                        //},
                        "modified": schema.StringAttribute{
                           MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            //Default: stringdefault.StaticString(time.Now().Format("2006-01-02T15:04:05.000Z")),
                            PlanModifiers: []planmodifier.String{
                                //UseStateForUnknown(),
                                //UsePlanForUnknownString(),
                                planmodifiers.UseEmptyStringForNull(),
                            },
                        },
                        "name": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Required: true,
                            //Optional: true,
                            //Computed: true,
                        },
                        "notes": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.UseStateForUnknown(),
                            },
                        },
                        //"only_fixed": schema.BoolAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: booldefault.StaticBool(false),
                        //},
                        "owner": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.UseStateForUnknown(),
                            },
                        },
                        //"package_types_thresholds": schema.SetNestedAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    PlanModifiers: []planmodifier.Set{
                        //        setplanmodifier.UseStateForUnknown(),
                        //    },
                        //    NestedObject: schema.NestedAttributeObject{
                        //        Attributes: map[string]schema.Attribute{
                        //            "alert_threshold": schema.SingleNestedAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Attributes: map[string]schema.Attribute{
                        //                    "disabled": schema.BoolAttribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                    "value": schema.Int32Attribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                },
                        //                Default: objectdefault.StaticValue(
                        //                    types.ObjectValueMust(
                        //                        map[string]attr.Type{
                        //                            "disabled": types.BoolType,
                        //                            "value": types.Int32Type,
                        //                        },
                        //                        map[string]attr.Value{
                        //                            "disabled": types.BoolValue(false),
                        //                            "value": types.Int32Value(int32(0)),
                        //                        },
                        //                    ),
                        //                ),
                        //            },
                        //            "block_threshold": schema.SingleNestedAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Attributes: map[string]schema.Attribute{
                        //                    "enabled": schema.BoolAttribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                    "value": schema.Int32Attribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                },
                        //                Default: objectdefault.StaticValue(
                        //                    types.ObjectValueMust(
                        //                        map[string]attr.Type{
                        //                            "enabled": types.BoolType,
                        //                            "value": types.Int32Type,
                        //                        },
                        //                        map[string]attr.Value{
                        //                            "enabled": types.BoolValue(false),
                        //                            "value": types.Int32Value(int32(0)),
                        //                        },
                        //                    ),
                        //                ),
                        //            },
                        //            "type": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Validators: []validator.String{
                        //                    stringvalidator.OneOf(
                        //                        "nodejs", 
                        //                        "gem", 
                        //                        "python", 
                        //                        "jar",
                        //                        "package",
                        //                        "windows",
                        //                        "binary",
                        //                        "nuget",
                        //                        "go",
                        //                        "app",
                        //                        "unknown",
                        //                    ),
                        //                },
                        //            },
                        //        },
                        //    },
                        //},
                        ////"previous_name": schema.StringAttribute{
                        ////    MarkdownDescription: "TODO",
                        ////    Optional: true,
                        ////    Computed: true,
                        ////},
                        //"principal": schema.SetAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    ElementType: types.StringType,
                        //    PlanModifiers: []planmodifier.Set{
                        //        setplanmodifier.UseStateForUnknown(),
                        //    },
                        //    Default: setdefault.StaticValue(
                        //        types.SetValueMust(
                        //            types.StringType,
                        //            []attr.Value{},
                        //        ),
                        //    ),
                        //},
                        //"risk_factors_effects": schema.SetNestedAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    PlanModifiers: []planmodifier.Set{
                        //        setplanmodifier.UseStateForUnknown(),
                        //    },
                        //    NestedObject: schema.NestedAttributeObject{
                        //        Attributes: map[string]schema.Attribute{
                        //            "effect": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Validators: []validator.String{
                        //                    stringvalidator.OneOf("alert", "block", "ignore"),
                        //                },
                        //                // possible values: [ignore,alert,block]
                        //            },
                        //            "risk_factor": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Validators: []validator.String{
                        //                    stringvalidator.OneOf(
                        //                        "Critical severity",
                        //                        "High severity",
                        //                        "Medium severity",
                        //                        "Has fix",
                        //                        "Remote Execution",
                        //                        "DoS - Low",
                        //                        "DoS - High",
                        //                        "Recent vulnerability",
                        //                        "Exploit exists - in the wild",
                        //                        "Exploit exists - POC",
                        //                        "Attack complexity: low",
                        //                        "Attack vector: network",
                        //                        "Reachable from the internet",
                        //                        "Listening ports",
                        //                        "Container is running as root",
                        //                        "No mandatory security profile applied",
                        //                        "Running as privileged container",
                        //                        "Package in use",
                        //                        "Sensitive information",
                        //                        "Root mount",
                        //                        "Runtime socket",
                        //                        "Host access",
                        //                    ),
                        //                },
                        //            },
                        //        },
                        //    },
                        //},
                        //"tags": schema.SetNestedAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    PlanModifiers: []planmodifier.Set{
                        //        setplanmodifier.UseStateForUnknown(),
                        //    },
                        //    NestedObject: schema.NestedAttributeObject{
                        //        Attributes: map[string]schema.Attribute{
                        //            "description": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //            },
                        //            "effect": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //                Validators: []validator.String{
                        //                    stringvalidator.OneOf("alert", "block", "ignore"),
                        //                },
                        //            },
                        //            "expiration": schema.SingleNestedAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Attributes: map[string]schema.Attribute{
                        //                    "date": schema.StringAttribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                    "enabled": schema.BoolAttribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                },
                        //            },
                        //            "name": schema.StringAttribute{
                        //                MarkdownDescription: "TODO",
                        //                Optional: true,
                        //                Computed: true,
                        //            },
                        //        },
                        //    },
                        //},
                        "verbose": schema.BoolAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            Default: booldefault.StaticBool(false),
                        },
                    },
                },
            },
        },
    }
}

//func (r *HostCompliancePolicyResource) GetCollectionsSchema() schema.SetNestedAttribute {
func (r *HostCompliancePolicyResource) GetCollectionsSchema() schema.ListNestedAttribute {
    //return schema.SetNestedAttribute{
    return schema.ListNestedAttribute{
        MarkdownDescription: "TODO",
        Optional: true,
        Computed: true,
        // TODO: see if we can omit all but the name field and get away with it
        //Default: setdefault.StaticValue(
        //    types.SetValueMust(
        //        system.CollectionObjectType(),
        //        []attr.Value{
        //            types.ObjectValueMust(
        //                system.CollectionObjectAttrTypeMap(),
        //                map[string]attr.Value{
        //                    "account_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "app_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "clusters": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "color": types.StringValue("#3FA2F7"),
        //                    "containers": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "description": types.StringValue("System - all resources collection"),
        //                    "functions": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "hosts": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "images": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "labels": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "modified": basetypes.NewStringUnknown(),
        //                    //"modified": types.StringValue(""),
        //                    "name": types.StringValue("All"),
        //                    "namespaces": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        //                    "owner": types.StringValue("admin"),
        //                    "prisma": types.BoolValue(false),
        //                    "system": types.BoolValue(true),
        //                },
        //            ),
        //        },
        //    ),
        //),
        //PlanModifiers: []planmodifier.Set{
        PlanModifiers: []planmodifier.List{
            planmodifiers.RemoveNullObjects(),
        },
        NestedObject: schema.NestedAttributeObject{
            Attributes: map[string]schema.Attribute{
                "account_ids": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "app_ids": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "clusters": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "color": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.String{
                    //    planmodifiers.UseDefaultColorForDefaultCollectionColor(), 
                    //},
                },
                "containers": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "description": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.String{
                    //    planmodifiers.UseDefaultForDefaultCollectionDescription(),
                    //},
                },
                "functions": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "hosts": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "images": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "labels": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "modified": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.String{
                        //UseStateForUnknown(),
                        //UsePlanForUnknownString(),
                        planmodifiers.UseEmptyStringForNull(),
                    },
                },
                "name": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.String{
                    //    planmodifiers.UseAllForDefaultCollectionName(),
                    //},
                },
                "namespaces": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Set{
                    //    planmodifiers.UseDefaultForUnknownCollectionSets(),
                    //},
                },
                "owner": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.String{
                    //    planmodifiers.UseSystemForDefaultCollectionOwner(), 
                    //},
                },
                "prisma": schema.BoolAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Bool{
                    //    planmodifiers.UseFalseForDefaultCollectionBools(), 
                    //},
                },
                "system": schema.BoolAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    //PlanModifiers: []planmodifier.Bool{
                    //    planmodifiers.UseTrueForDefaultCollectionBools(), 
                    //},
                },
            },
        },
    }
}
