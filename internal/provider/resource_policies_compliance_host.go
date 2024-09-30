package provider

import (
    "context"
    //"time"
	"fmt"
    "reflect"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
    //"github.com/hashicorp/terraform-plugin-log/tflog"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/convert"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	

    "github.com/hashicorp/terraform-plugin-framework/types"
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
    //Action types.Set `tfsdk:"action"`
    //Modified types.String `tfsdk:"modified"`
    Name types.String `tfsdk:"name"`
    Collections types.Set `tfsdk:"collections"`
    //Action types.Set `tfsdk:"action"`
    ////AlertThreshold *HostCompliancePolicyRuleAlertThresholdResourceModel `tfsdk:"alert_threshold"`
    //ReportAllCompliance types.Bool `tfsdk:"report_all_compliance"`
    ////AuditAllowed types.Bool `tfsdk:"audit_allowed"`
    //BlockMessage types.String `tfsdk:"block_message"`
    ////BlockThreshold  *HostCompliancePolicyRuleBlockThresholdResourceModel `tfsdk:"block_threshold"`
    //Condition *HostCompliancePolicyRuleConditionResourceModel `tfsdk:"condition"`
    Condition types.Object `tfsdk:"condition"`
    //CreatePR types.Bool `tfsdk:"create_pr"
    //CVERules *[]HostCompliancePolicyRuleCVERuleResourceModel `tfsdk:"cve_rules"`
    ////Disabled types.Bool `tfsdk:"disabled"`
    Effect types.String `tfsdk:"effect"`
    ////ExcludeBaseImageVulns types.Bool `tfsdk:"exclude_base_image_vulns"`
    //GraceDays types.Int32 `tfsdk:"grace_days"`
    //GraceDaysPolicy *HostCompliancePolicyRuleGraceDaysPolicyResourceModel `tfsdk:"grace_days_policy"`
    ////GraceDaysPolicy types.Object `tfsdk:"grace_days_policy"`
    //Groups types.Set `tfsdk:"groups"`
    //License *HostCompliancePolicyRuleLicenseResourceModel `tfsdk:"license"`
    //License types.Object `tfsdk:"license"`
    //Modified types.String `tfsdk:"modified"`
    //Name types.String `tfsdk:"name"`
    Notes types.String `tfsdk:"notes"`
    //OnlyFixed types.Bool `tfsdk:"only_fixed"`
    //Owner types.String `tfsdk:"owner"`
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

// condition set modifier 
type useStateForUnknownCondition struct {}

func (m useStateForUnknownCondition) Description(_ context.Context) string {
    return ""
}

func (m useStateForUnknownCondition) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useStateForUnknownCondition) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    fmt.Println("***********************")
    fmt.Println("entering PlanModifySet for UseStateForUnknownCondition")
    fmt.Println("***********************")
    
    if resp.PlanValue.IsUnknown() {
        return
        //vulnerabilityListType := types.ListType{
        //    ElemType: types.ObjectType{
        //        AttrTypes: map[string]attr.Type{
        //            "id": types.Int32Type,
        //            "block": types.BoolType,
        //        },
        //    },
        //}

        //conditionObjectValue := types.ObjectValueMust(
        //    map[string]attr.Type{
        //        "device":  types.StringType,
        //        "read_only":  types.BoolType,
        //        "vulnerabilities": vulnerabilityListType,
        //    },
        //    map[string]attr.Value{},
        //)

        //req.PlanValue = conditionObjectValue
    }
}

func UseStateForUnknownCondition() planmodifier.Object {
    return useStateForUnknownCondition{} 
}

// string plan modifier
func UseAlertBlockForBlockEffect() planmodifier.String {
    return &useAlertBlockForBlockEffect{}
}


type useAlertBlockForBlockEffect struct {}

func (m *useAlertBlockForBlockEffect) Description(_ context.Context) string {
    return ""
}

func (m *useAlertBlockForBlockEffect) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useAlertBlockForBlockEffect) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    fmt.Println("***********************")
    fmt.Println("entering PlanModifyString")
    fmt.Println(req.PlanValue.ValueString())
    fmt.Println("***********************")

    //resp.PlanValue = 

    //if req.PlanValue.ValueString() == "block" {
    //    req.PlanValue = types.StringValue("alert, block")
    //    return
    //}

    //resp.PlanValue = req.StateValue
}

func (r *HostCompliancePolicyResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                //Default: stringdefault.StaticString("hostCompliance"),
            },
            "policy_type": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("hostCompliance"),
            },
            //"rules": schema.SetNestedAttribute{
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                //Computed: true,
                //PlanModifiers: []planmodifier.Set{
                //    setplanmodifier.UseStateForUnknown(),
                //},
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
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
                        //"report_all_compliance": schema.BoolAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: booldefault.StaticBool(false), 
                        //},
                        //"block_message": schema.StringAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: stringdefault.StaticString(""),
                        //},
                        "collections": r.GetCollectionsSchema(),
                        "condition": schema.SingleNestedAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            //PlanModifiers: []planmodifier.Object{
                            //    UseStateForUnknownCondition(),
                            //},
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
                        ////"disabled": schema.BoolAttribute{
                        ////    MarkdownDescription: "TODO",
                        ////    Optional: true,
                        ////    Computed: true,
                        ////    //Default: booldefault.StaticBool(false), 
                        ////},
                        "effect": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                UseAlertBlockForBlockEffect(),
                            },
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
                        //"modified": schema.StringAttribute{
                        //   MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    //Default: stringdefault.StaticString(time.Now().Format("2006-01-02T15:04:05.000Z")),
                        //    //PlanModifiers: []planmodifier.String{
                        //    //    //UseStateForUnknown(),
                        //    //    UsePlanForUnknownString(),
                        //    //},
                        //},
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
                        },
                        //"only_fixed": schema.BoolAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: booldefault.StaticBool(false),
                        //},
                        //"owner": schema.StringAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    //Default: stringdefault.StaticString(username),
                        //},
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

func (r *HostCompliancePolicyResource) GetCollectionsSchema() schema.SetNestedAttribute {
    return schema.SetNestedAttribute{
        MarkdownDescription: "TODO",
        Optional: true,
        Computed: true,
        Default: setdefault.StaticValue(
            types.SetValueMust(
                collectionObjectType(),
                []attr.Value{
                    types.ObjectValueMust(
                        collectionObjectAttrTypeMap(),
                        map[string]attr.Value{
                            "account_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "app_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "clusters": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "color": types.StringValue("#3FA2F7"),
                            "containers": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "description": types.StringValue("System - all resources collection"),
                            "functions": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "hosts": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "images": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "labels": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "modified": basetypes.NewStringUnknown(),
                            //"modified": types.StringValue(""),
                            "name": types.StringValue("All"),
                            "namespaces": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "owner": types.StringValue("admin"),
                            "prisma": types.BoolValue(false),
                            "system": types.BoolValue(true),
                        },
                    ),
                },
            ),
        ),
        NestedObject: schema.NestedAttributeObject{
            Attributes: map[string]schema.Attribute{
                "account_ids": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "app_ids": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "clusters": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "color": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.String{
                        planmodifiers.UseDefaultColorForDefaultCollectionColor(), 
                    },
                },
                "containers": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "description": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.String{
                        planmodifiers.UseDefaultForDefaultCollectionDescription(),
                    },
                },
                "functions": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "hosts": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "images": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "labels": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "modified": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Computed: true,
                    PlanModifiers: []planmodifier.String{
                        planmodifiers.UseCurrentTimeForDefaultCollectionModified(),
                    },
                },
                "name": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.String{
                        planmodifiers.UseAllForDefaultCollectionName(),
                    },
                },
                "namespaces": schema.SetAttribute{
                    MarkdownDescription: "TODO",
                    ElementType: types.StringType,
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Set{
                        planmodifiers.UseDefaultForUnknownCollectionSets(),
                    },
                },
                "owner": schema.StringAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.String{
                        planmodifiers.UseSystemForDefaultCollectionOwner(), 
                    },
                },
                "prisma": schema.BoolAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Bool{
                        planmodifiers.UseFalseForDefaultCollectionBools(), 
                    },
                },
                "system": schema.BoolAttribute{
                    MarkdownDescription: "TODO",
                    Optional: true,
                    Computed: true,
                    PlanModifiers: []planmodifier.Bool{
                        planmodifiers.UseTrueForDefaultCollectionBools(), 
                    },
                },
            },
        },
    }
}

func (r *HostCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema()
}

func (r *HostCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*api.PrismaCloudComputeAPIClient)

    if !ok {
        resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

        return
    }

    r.client = client
}

func (r *HostCompliancePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("entering ModifyPlan")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

    var plan HostCompliancePolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("starting loop over rules")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    for index, rule := range *plan.Rules {
        if rule.Effect.IsUnknown() {
            resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("effect"), types.StringValue("alert"))...)

            vulnerabilitiesData, err := systemAPI.GetComplianceHostVulnerabilities(*r.client)
	        if err != nil {
	        	diags.AddError(
                    "Error modifying planned policy rules", 
                    "Failed to retrieve compliance host vulnerabilities while modifying plan rules: " + err.Error(),
                )
                return
	        }

            complianceVulnerabilities := vulnerabilitiesData.ComplianceVulnerabilities

            vulnerabilityObjectValues := []attr.Value{}
            for _, vuln := range complianceVulnerabilities {
                if vuln.Severity == "high" || vuln.Severity == "critical" {
                    vulnerabilityObjectValue := types.ObjectValueMust(
                        map[string]attr.Type{
                            "id":        types.Int32Type,
                            "block":       types.BoolType,
                        },
                        map[string]attr.Value{
                            "id": types.Int32Value(int32(vuln.Id)),
                            "block": types.BoolValue(false),
                        },
                    )
                   
                    vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
                }
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
                return
            }

            conditionObject := types.ObjectValueMust(
                map[string]attr.Type{
                    //"device": types.StringType,
                    //"read_only": types.BoolType,
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
                    //"device": types.StringValue(rule.Condition.Device),
                    //"read_only": types.BoolValue(rule.Condition.ReadOnly),
                    "vulnerabilities": vulnerabilityObject,
                },
            )
            
            resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)
        } else if rule.Effect.ValueString() == "alert, block" {
            var ruleConditionVulns []policyAPI.HostCompliancePolicyRuleVulnerability
            if rule.Condition.IsUnknown() {
                return
            } else {
                ruleCondition := policyAPI.HostCompliancePolicyRuleCondition{} 
                diags = rule.Condition.As(ctx, &ruleCondition, basetypes.ObjectAsOptions{})
                if diags.HasError() {
                    return
                }
                ruleConditionVulns = ruleCondition.Vulnerabilities
            }

            vulnerabilityObjectValues := []attr.Value{}
            for _, vuln := range ruleConditionVulns {
                vulnerabilityObjectValue := types.ObjectValueMust(
                    map[string]attr.Type{
                        "id":        types.Int32Type,
                        "block":       types.BoolType,
                    },
                    map[string]attr.Value{
                        "id": types.Int32Value(int32(vuln.Id)),
                        //"block": types.BoolValue(false),
                        "block": types.BoolValue(vuln.Block),
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
                return
            }

            conditionObject := types.ObjectValueMust(
                map[string]attr.Type{
                    //"device": types.StringType,
                    //"read_only": types.BoolType,
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
                    //"device": types.StringValue(rule.Condition.Device),
                    //"read_only": types.BoolValue(rule.Condition.ReadOnly),
                    "vulnerabilities": vulnerabilityObject,
                },
            )
            
            resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)

        } else if rule.Effect.ValueString() == "alert" {
            fmt.Println("@@@@@@@@@@@@@@")
            fmt.Println("effect is alert")
            fmt.Println("@@@@@@@@@@@@@@")
            vulnerabilitiesData, err := systemAPI.GetComplianceHostVulnerabilities(*r.client)
	        if err != nil {
	        	//diags.AddError(
	        	resp.Diagnostics.AddError(
                    "Error modifying planned policy rules", 
                    "Failed to retrieve compliance host vulnerabilities while modifying plan rules: " + err.Error(),
                )
                return
	        }

            complianceVulnerabilities := vulnerabilitiesData.ComplianceVulnerabilities

            vulnerabilityObjectValues := []attr.Value{}
            for _, vuln := range complianceVulnerabilities {
                vulnerabilityObjectValue := types.ObjectValueMust(
                    map[string]attr.Type{
                        "id":        types.Int32Type,
                        "block":       types.BoolType,
                    },
                    map[string]attr.Value{
                        "id": types.Int32Value(int32(vuln.Id)),
                        "block": types.BoolValue(false),
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
                return
            }

            conditionObject := types.ObjectValueMust(
                map[string]attr.Type{
                    //"device": types.StringType,
                    //"read_only": types.BoolType,
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
                    //"device": types.StringValue(rule.Condition.Device),
                    //"read_only": types.BoolValue(rule.Condition.ReadOnly),
                    "vulnerabilities": vulnerabilityObject,
                },
            )
            
            resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)
        } else if rule.Effect.ValueString() == "ignore" {
            fmt.Println("@@@@@@@@@@@@@@")
            fmt.Println("effect is ignore")
            fmt.Println("@@@@@@@@@@@@@@")
            vulnerabilityObjectValues := []attr.Value{}

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
                return
            }

            conditionObject := types.ObjectValueMust(
                map[string]attr.Type{
                    //"device": types.StringType,
                    //"read_only": types.BoolType,
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
                    //"device": types.StringValue(rule.Condition.Device),
                    //"read_only": types.BoolValue(rule.Condition.ReadOnly),
                    "vulnerabilities": vulnerabilityObject,
                },
            )
            
            //rule.Condition = conditionObject
            resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)
        } else if rule.Effect.ValueString() == "block" {
            fmt.Println("@@@@@@@@@@@@@@")
            fmt.Println("effect is block")
            fmt.Println("@@@@@@@@@@@@@@")

            vulnerabilitiesData, err := systemAPI.GetComplianceHostVulnerabilities(*r.client)
	        if err != nil {
	        	diags.AddError(
                    "Error modifying planned policy rules", 
                    "Failed to retrieve compliance host vulnerabilities while modifying plan rules: " + err.Error(),
                )
                return
	        }
            
            complianceVulnerabilities := vulnerabilitiesData.ComplianceVulnerabilities

            vulnerabilityObjectValues := []attr.Value{}
            for _, vuln := range complianceVulnerabilities {
                var block bool
                if vuln.Type == "windows" {
                    block = false
                } else {
                    block = true
                }

                vulnerabilityObjectValue := types.ObjectValueMust(
                    map[string]attr.Type{
                        "id":        types.Int32Type,
                        "block":       types.BoolType,
                    },
                    map[string]attr.Value{
                        "id": types.Int32Value(int32(vuln.Id)),
                        "block": types.BoolValue(block),
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
                return
            }

            conditionObject := types.ObjectValueMust(
                map[string]attr.Type{
                    //"device": types.StringType,
                    //"read_only": types.BoolType,
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
                    //"device": types.StringValue(rule.Condition.Device),
                    //"read_only": types.BoolValue(rule.Condition.ReadOnly),
                    "vulnerabilities": vulnerabilityObject,
                },
            )
            
            //rule.Condition = conditionObject
            resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)
        }
    }
    //resp.Plan.Set(ctx, &plan)

    var respPlan HostCompliancePolicyResourceModel
    diags = resp.Plan.Get(ctx, &respPlan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    fmt.Printf("%+v\n", respPlan)
    fmt.Printf("%+v\n", *respPlan.Rules)
    //fmt.Printf("%+v\n", &respPlan.Rules.Elements()[0].Condition)
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("exiting ModifyPlan")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    
    var ru types.List 
    dgs := resp.Plan.GetAttribute(ctx, path.Root("rules"), &ru)
    if dgs.HasError() {
        fmt.Print(dgs)
        return
    }
    for _, li := range ru.Elements() {
        fmt.Println(li)
    }
    return
}

func (r *HostCompliancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // TODO: refine this logic to populate Owner with the value in config, if it exists
    //var username types.String
    //diags := req.Config.GetAttribute(ctx, path.Root("username"), &username)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    // Retrieve values from plan
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("retrieving plan and serializing into HostCompliancePolicyResourceModel")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    var plan HostCompliancePolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println(*plan.Rules)
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

    // Generate API request body from plan
    policy, diags := schemaToPolicy(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new host compliance policy 
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("creating policy resource with payload:")
    fmt.Printf("%+v\n", policy)
    r1 := *policy.Rules
    fmt.Printf("%+v\n", r1)
    //fmt.Printf("%+v\n", *r1[0].Condition)
    fmt.Println("number of vulns:")
    fmt.Println(len(r1[0].Condition.Vulnerabilities))
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    err := policyAPI.UpsertHostCompliancePolicy(*r.client, policy)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Host Compliance Policy resource", 
            "Failed to create host compliance policy: " + err.Error(),
        )
        return
	}

    // Retrieve newly created host compliance policy 
    response, err := policyAPI.GetHostCompliancePolicy(*r.client)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created Host Compliance Policy resource", 
            "Failed to retrieve created host compliance policy: " + err.Error(),
        )
        return
    }

    createdPolicy, diags := policyToSchema(ctx, *response)
    if diags.HasError() {
        return
    }

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("createdPolicy in Create():")
    fmt.Println(reflect.TypeOf(createdPolicy))
    fmt.Printf("%+v\n", createdPolicy)
    fmt.Printf("%+v\n", *createdPolicy.Rules)
    fmt.Println("number of vulns:")
    fmt.Println(len(r1[0].Condition.Vulnerabilities))
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    
    // Set state to collection data
    diags = resp.State.Set(ctx, createdPolicy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%")
        fmt.Println("error in resp.State.Set")
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%")
        return
    }

}

func (r *HostCompliancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("we're in Read")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")

    // Get current state
    var state HostCompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get policy value from Prisma Cloud
    policy, err := policyAPI.GetHostCompliancePolicy(*r.client)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Host Compliance Policy resource", 
            "Failed to read Host Compliance Policy: " + err.Error(),
        )
        return
    }

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("retrieved host compliance policy: ") 
    fmt.Printf("%+v\n", policy)
    fmt.Printf("%+v\n", *policy.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
  
    // Overwrite state values with Prisma Cloud data
    policySchema, diags := policyToSchema(ctx, *policy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("policySchema: ") 
    fmt.Printf("%+v\n", policySchema)
    fmt.Printf("%+v\n", policySchema.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")

    // Set refreshed state
    diags = resp.State.Set(ctx, &policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *HostCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state HostCompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan HostCompliancePolicyResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    planPolicy, diags := schemaToPolicy(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Update existing policy
    err := policyAPI.UpsertHostCompliancePolicy(*r.client, planPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating Host Compliance Policy resource", 
            "Failed to update host compliance policy: " + err.Error(),
        )
        return
	}

    // Get updated policy value from Prisma Cloud
    policy, err := policyAPI.GetHostCompliancePolicy(*r.client)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Host Compliance Policy resource", 
            "Failed to read Host Compliance Policy: " + err.Error(),
        )
        return
    }

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("retrieved host compliance policy during Update() execution: ") 
    fmt.Printf("%+v\n", policy)
    fmt.Printf("%+v\n", *policy.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
  
    // Convert updated policy into schema
    plan, diags = policyToSchema(ctx, *policy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set updated state
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *HostCompliancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state HostCompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Clear policy rules
    emptyRules := []HostCompliancePolicyRuleResourceModel{}
    state.Rules = &emptyRules

    // Generate API request body from plan
    updatedPlan, diags := schemaToPolicy(ctx, &state, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing policy 
    err := policyAPI.UpsertHostCompliancePolicy(*r.client, updatedPlan)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Host Compliance Policy resource", 
            "Failed to delete host compliance policy: " + err.Error(),
        )
        return
	}
}

// TODO: Define ImportState to work properly with this resource
func (r *HostCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("executing ImportState")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func schemaToPolicy(ctx context.Context, plan *HostCompliancePolicyResourceModel, client *api.PrismaCloudComputeAPIClient,/*, username types.String*/) (policyAPI.HostCompliancePolicy, diag.Diagnostics) {
    var diags diag.Diagnostics

    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("entering schemaToPolicy")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    policy := policyAPI.HostCompliancePolicy{
        Id: plan.Id.ValueString(),
        PolicyType: plan.PolicyType.ValueString(),
    }

    if plan.Rules == nil {
        rules := []policyAPI.HostCompliancePolicyRule{}
        policy.Rules = &rules
        return policy, diags
    } else {
        rules, diags := ruleSchemaToPolicy(ctx, *plan.Rules, client)
        if diags.HasError() {
            return policy, diags
        }

        policy.Rules = &rules
    }

    return policy, diags
}

func ruleSchemaToPolicy(ctx context.Context, planRules []HostCompliancePolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient, /*, username types.String*/) ([]policyAPI.HostCompliancePolicyRule, diag.Diagnostics) {
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("entering ruleSchemaToPolicy")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    _ = reflect.TypeOf(ctx)

    var diags diag.Diagnostics

    rules := []policyAPI.HostCompliancePolicyRule{}

    for _, planRule := range planRules {
        collections := []collectionAPI.Collection{}
        planCollections := make([]types.Object, 0, len(planRule.Collections.Elements()))
        diags = planRule.Collections.ElementsAs(ctx, &planCollections, false)
        if diags.HasError() {
            return rules, diags
        }

        for _, planCollection := range planCollections {
            collectionModel := CollectionResourceModel{}
            diags = planCollection.As(ctx, &collectionModel, basetypes.ObjectAsOptions{})
            if diags.HasError() {
                return rules, diags
            }

            accountIds := make([]string, 0, len(collectionModel.AccountIDs.Elements()))
            diags = collectionModel.AccountIDs.ElementsAs(ctx, &accountIds, false)
            if diags.HasError() {
                return rules, diags
            }

            appIds := make([]string, 0, len(collectionModel.AppIDs.Elements()))
            diags = collectionModel.AccountIDs.ElementsAs(ctx, &appIds, false)
            if diags.HasError() {
                return rules, diags
            }

            clusters := make([]string, 0, len(collectionModel.Clusters.Elements()))
            diags = collectionModel.Clusters.ElementsAs(ctx, &clusters, false)
            if diags.HasError() {
                return rules, diags
            }

            containers := make([]string, 0, len(collectionModel.Containers.Elements()))
            diags = collectionModel.Containers.ElementsAs(ctx, &containers, false)
            if diags.HasError() {
                return rules, diags
            }

            functions := make([]string, 0, len(collectionModel.Functions.Elements()))
            diags = collectionModel.Functions.ElementsAs(ctx, &functions, false)
            if diags.HasError() {
                return rules, diags
            }

            hosts := make([]string, 0, len(collectionModel.Hosts.Elements()))
            diags = collectionModel.Hosts.ElementsAs(ctx, &hosts, false)
            if diags.HasError() {
                return rules, diags
            }

            images := make([]string, 0, len(collectionModel.Images.Elements()))
            diags = collectionModel.Images.ElementsAs(ctx, &images, false)
            if diags.HasError() {
                return rules, diags
            }

            labels := make([]string, 0, len(collectionModel.Labels.Elements()))
            diags = collectionModel.Labels.ElementsAs(ctx, &labels, false)
            if diags.HasError() {
                return rules, diags
            }

            namespaces := make([]string, 0, len(collectionModel.Namespaces.Elements()))
            diags = collectionModel.Namespaces.ElementsAs(ctx, &namespaces, false)
            if diags.HasError() {
                return rules, diags
            }

            collections = append(collections, collectionAPI.Collection{
                AccountIDs: accountIds,
                AppIDs: appIds,
                Clusters: clusters,
                Color: collectionModel.Color.ValueString(),
                Containers: containers,
                Description: collectionModel.Description.ValueString(),
                Functions: functions,
                Hosts: hosts,
                Images: images,
                Labels: labels,
                Name: collectionModel.Name.ValueString(),
                Namespaces: namespaces,
                Owner: collectionModel.Owner.ValueString(),
                Prisma: collectionModel.Prisma.ValueBool(),
                System: collectionModel.System.ValueBool(),
            })
        }

        // This fails due to the value of "modified" being unknown
        //collections := []collectionAPI.Collection{}
        //diags = planRule.Collections.ElementsAs(ctx, &col, false)
        //if diags.HasError() {
        //    fmt.Println(diags)
        //    return rules, diags
        //}

        if planRule.Effect.ValueString() == "alert, block" && planRule.Condition.IsUnknown() {
            diags.AddError(
                "Missing condition from \"alert, block\" effect rule",
                "Condition attribute must be defined for rules with effect \"alert, block\".",
            )
            return rules, diags
        }

        condition := policyAPI.HostCompliancePolicyRuleCondition{} 
        diags = planRule.Condition.As(ctx, &condition, basetypes.ObjectAsOptions{})
        if diags.HasError() {
            return rules, diags
        }

        rule := policyAPI.HostCompliancePolicyRule{
            Name: planRule.Name.ValueString(), 
            Collections: collections,
            //Collections: col,
            Condition: &condition,
            Effect: planRule.Effect.ValueString(),
            Verbose: planRule.Verbose.ValueBool(),
        }
        
        if !planRule.Notes.IsUnknown() && !planRule.Notes.IsNull() {
            rule.Notes = planRule.Notes.ValueString()
        }

        rules = append(rules, rule)
    }

    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("exiting ruleSchemaToPolicy")
    //fmt.Printf("%+v\n", rules)
    fmt.Printf("%+v\n", *rules[0].Condition)
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    
    return rules, diags
}

func policyToSchema(ctx context.Context, policy policyAPI.HostCompliancePolicy) (HostCompliancePolicyResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    schema := HostCompliancePolicyResourceModel{
        Id: types.StringValue(policy.Id),
        PolicyType: types.StringValue(policy.PolicyType),
    }

    rules, diags := policyRulesToSchema(ctx, *policy.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("!!! returned from policyRulesToSchema: ") 
    //fmt.Printf("%+v\n", rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    if diags.HasError() {
        return schema, diags
    }

    schema.Rules = &rules

    //if policy.Action != nil {
    //    action, diags := types.SetValueFrom(ctx, types.StringType, policy.Action)
    //    if diags.HasError() {
    //        return schema, diags
    //    }

    //    schema.Action = action
    //}

    //if policy.Groups != nil {
    //    groups, diags := types.SetValueFrom(ctx, types.StringType, policy.Groups)
    //    if diags.HasError() {
    //        return schema, diags
    //    }

    //    schema.Groups = groups 
    //}

    //if policy.Principal != nil {
    //    principal, diags := types.SetValueFrom(ctx, types.StringType, policy.Principal)
    //    if diags.HasError() {
    //        return schema, diags
    //    }

    //    schema.Principal = principal 
    //}

    return schema, diags
}

func policyRulesToSchema(ctx context.Context, rules []policyAPI.HostCompliancePolicyRule) ([]HostCompliancePolicyRuleResourceModel, diag.Diagnostics) {
    fmt.Println("***********************")
    fmt.Println("entering policyRulesToSchema")
    fmt.Println("***********************")

    var diags diag.Diagnostics

    schemaRules := []HostCompliancePolicyRuleResourceModel{}

    if len(rules) == 0 {
        return nil, diags
    }

    for _, rule := range rules {
        //ruleObject := HostCompliancePolicyRuleResourceModel {
        //    //ReportAllCompliance: types.BoolValue(rule.ReportAllCompliance),
        //    //BlockMessage: types.StringValue(rule.BlockMessage),
        //    //CreatePR: types.BoolValue(rule.CreatePR),
        //    ////Disabled: types.BoolValue(rule.Disabled),
        //    //Effect: types.StringValue(rule.Effect),
        //    ////ExcludeBaseImageVulns: types.BoolValue(rule.ExcludeBaseImageVulns),
        //    //GraceDays: types.Int32Value(rule.GraceDays),
        //    //Modified: types.StringValue(rule.Modified),
        //    Name: types.StringValue(rule.Name),
        //    //Notes: types.StringValue(rule.Notes),
        //    //OnlyFixed: types.BoolValue(rule.OnlyFixed),
        //    //Owner: types.StringValue(rule.Owner),
        //    ////PreviousName: types.StringValue(rule.PreviousName),
        //    //Verbose: types.BoolValue(rule.Verbose),
        //}

        schemaRule := HostCompliancePolicyRuleResourceModel{
            Name: types.StringValue(rule.Name),
            Effect: types.StringValue(rule.Effect),
            Verbose: types.BoolValue(rule.Verbose),
        }

        if rule.Collections != nil {
            collectionObjectValues := []attr.Value{}
            for _, collection := range(rule.Collections) {
                accountIDs, diags := types.SetValueFrom(ctx, types.StringType, collection.AccountIDs)
                if diags.HasError() {
                    return schemaRules, diags
                }

                appIDs, diags := types.SetValueFrom(ctx, types.StringType, collection.AppIDs)
                if diags.HasError() {
                    return schemaRules, diags
                }

                clusters, diags := types.SetValueFrom(ctx, types.StringType, collection.Clusters)
                if diags.HasError() {
                    return schemaRules, diags
                }

                containers, diags := types.SetValueFrom(ctx, types.StringType, collection.Containers)
                if diags.HasError() {
                    return schemaRules, diags
                }

                functions, diags := types.SetValueFrom(ctx, types.StringType, collection.Functions)
                if diags.HasError() {
                    return schemaRules, diags
                }

                hosts, diags := types.SetValueFrom(ctx, types.StringType, collection.Hosts)
                if diags.HasError() {
                    return schemaRules, diags
                }

                images, diags := types.SetValueFrom(ctx, types.StringType, collection.Images)
                if diags.HasError() {
                    return schemaRules, diags
                }

                labels, diags := types.SetValueFrom(ctx, types.StringType, collection.Labels)
                if diags.HasError() {
                    return schemaRules, diags
                }
                
                namespaces, diags := types.SetValueFrom(ctx, types.StringType, collection.Namespaces)
                if diags.HasError() {
                    return schemaRules, diags
                }

                collectionObjectValue := types.ObjectValueMust(
                    collectionObjectAttrTypeMap(),
                    map[string]attr.Value{
                        "account_ids": accountIDs,
                        "app_ids": appIDs,
                        "clusters": clusters,
                        "color": types.StringValue(collection.Color),
                        "containers": containers,
                        "description": types.StringValue(collection.Description),
                        "functions": functions,
                        "hosts": hosts,
                        "images": images,
                        "labels": labels,
                        "modified": types.StringValue(collection.Modified),
                        //"modified": types.StringValue(""),
                        "name": types.StringValue(collection.Name),
                        "namespaces": namespaces,
                        "owner": types.StringValue(collection.Owner),
                        "prisma": types.BoolValue(collection.Prisma),
                        "system": types.BoolValue(collection.System),
                    },
                )

                collectionObjectValues = append(collectionObjectValues, collectionObjectValue)
            }

            collectionSet, diags := types.SetValueFrom(
                ctx,
                collectionObjectType(),
                collectionObjectValues,
            )

            if diags.HasError() {
                return schemaRules, diags
            }

            schemaRule.Collections = collectionSet
        }

        if rule.Effect == "alert, block" {
            rule.Effect = "block" 
        }

        if rule.Condition != nil {
            vulnerabilityObjectValues := []attr.Value{}
            for _, vulnerability := range(rule.Condition.Vulnerabilities) {
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
                return schemaRules, diags
            }

            conditionObject := types.ObjectValueMust(
                map[string]attr.Type{
                    //"device": types.StringType,
                    //"read_only": types.BoolType,
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
                    //"device": types.StringValue(rule.Condition.Device),
                    //"read_only": types.BoolValue(rule.Condition.ReadOnly),
                    "vulnerabilities": vulnerabilityObject,
                },
            )
            
            schemaRule.Condition = conditionObject
        }

        if rule.Notes != "" {
            schemaRule.Notes = types.StringValue(rule.Notes)
        }
            
        schemaRules = append(schemaRules, schemaRule)
    }

    fmt.Println("***********************")
    fmt.Println("exiting policyRulesToSchema")
    fmt.Printf("%+v\n", schemaRules)
    fmt.Println("***********************")

    return schemaRules, diags
}
