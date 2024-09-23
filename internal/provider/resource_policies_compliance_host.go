package provider

import (
    "context"
    "time"
	"fmt"
    "reflect"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
    //"github.com/hashicorp/terraform-plugin-log/tflog"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
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
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	//"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	

    "github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
    //Collections *[]CollectionResourceModel `tfsdk:"collections"`
    //Collections []CollectionResourceModel `tfsdk:"collections"`
    Collections types.Set `tfsdk:"collections"`
    //Collections types.Object `tfsdk:"collections"`
    

    //Action types.Set `tfsdk:"action"`
    ////AlertThreshold *HostCompliancePolicyRuleAlertThresholdResourceModel `tfsdk:"alert_threshold"`
    //ReportAllCompliance types.Bool `tfsdk:"report_all_compliance"`
    ////AuditAllowed types.Bool `tfsdk:"audit_allowed"`
    //BlockMessage types.String `tfsdk:"block_message"`
    ////BlockThreshold  *HostCompliancePolicyRuleBlockThresholdResourceModel `tfsdk:"block_threshold"`
    //Collections *[]CollectionResourceModel `tfsdk:"collections"`
    //Collections []CollectionResourceModel `tfsdk:"collections"`
    //Collections types.Set `tfsdk:"collections"`
    //Collections types.Object `tfsdk:"collections"`
    //Condition *HostCompliancePolicyRuleConditionResourceModel `tfsdk:"condition"`
    ////Condition types.Object `tfsdk:"condition"`
    //CreatePR types.Bool `tfsdk:"create_pr"
    //CVERules *[]HostCompliancePolicyRuleCVERuleResourceModel `tfsdk:"cve_rules"`
    ////Disabled types.Bool `tfsdk:"disabled"`
    //Effect types.String `tfsdk:"effect"`
    ////ExcludeBaseImageVulns types.Bool `tfsdk:"exclude_base_image_vulns"`
    //GraceDays types.Int32 `tfsdk:"grace_days"`
    //GraceDaysPolicy *HostCompliancePolicyRuleGraceDaysPolicyResourceModel `tfsdk:"grace_days_policy"`
    ////GraceDaysPolicy types.Object `tfsdk:"grace_days_policy"`
    //Groups types.Set `tfsdk:"groups"`
    //License *HostCompliancePolicyRuleLicenseResourceModel `tfsdk:"license"`
    //License types.Object `tfsdk:"license"`
    //Modified types.String `tfsdk:"modified"`
    //Name types.String `tfsdk:"name"`
    ////Notes types.String `tfsdk:"notes"`
    //OnlyFixed types.Bool `tfsdk:"only_fixed"`
    //Owner types.String `tfsdk:"owner"`
    //PkgTypesThresholds *[]HostCompliancePolicyRulePkgTypesThresholdsResourceModel `tfsdk:"package_types_thresholds"`
    ////PreviousName types.String `tfsdk:"previous_name"`
    //Principal types.Set `tfsdk:"principal"`
    //RiskFactorsEffects *[]HostCompliancePolicyRuleRiskFactorsEffectsResourceModel `tfsdk:"risk_factors_effects"`
    //Tags *[]HostCompliancePolicyRuleTagsResourceModel `tfsdk:"tags"`
    //Verbose types.Bool `tfsdk:"verbose"`
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

// CVE Rules
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

// set plan modifier
type useStateForUnknownSetModifier struct {}

func (m useStateForUnknownSetModifier) Description(_ context.Context) string {
    return ""
}

func (m useStateForUnknownSetModifier) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useStateForUnknownSetModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    fmt.Println("***********************")
    fmt.Println("entering PlanModifySet")
    fmt.Println("***********************")

    if req.ConfigValue.IsUnknown() {
        fmt.Println("***********************")
        fmt.Println("entering ConfigValue.IsUnknown block of PlanModifySet")
        fmt.Println("***********************")
        return
    }

    if req.PlanValue.IsUnknown() {
        fmt.Println("***********************")
        fmt.Println("entering isUnknown block of PlanModifySet")
        fmt.Println("***********************")
        c := types.ObjectValueMust(
            map[string]attr.Type{
                "account_ids":  types.SetType{ElemType: types.StringType},
                "app_ids":  types.SetType{ElemType: types.StringType},
                "clusters":  types.SetType{ElemType: types.StringType},
                "color":        types.StringType,
                "containers":  types.SetType{ElemType: types.StringType},
                "description":  types.StringType,
                "functions":  types.SetType{ElemType: types.StringType},
                "hosts":  types.SetType{ElemType: types.StringType},
                "images":  types.SetType{ElemType: types.StringType},
                "labels":  types.SetType{ElemType: types.StringType},
                "modified": types.StringType,
                "name":         types.StringType,
                "namespaces":  types.SetType{ElemType: types.StringType},
                "owner":        types.StringType,
                "prisma":       types.BoolType,
                "system":       types.BoolType,
            },
            map[string]attr.Value{
                "account_ids": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "app_ids": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "clusters": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "color": types.StringValue("#FFFFFF"),
                "containers": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "description": types.StringValue(""),
                "functions": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "hosts": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "images": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "labels": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "modified": types.StringValue(time.Now().Format("2006-01-02T15:04:05.000Z")),
                "name": types.StringValue("All"),
                "namespaces": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
                "owner": types.StringValue(""),
                "prisma": types.BoolValue(true),
                "system": types.BoolValue(true),
            },
        )
        c2 := types.SetValueMust(
            types.ObjectType{
                AttrTypes: map[string]attr.Type{
                    "account_ids":  types.SetType{ElemType: types.StringType},
                    "app_ids":  types.SetType{ElemType: types.StringType},
                    "clusters":  types.SetType{ElemType: types.StringType},
                    "color":        types.StringType,
                    "containers":  types.SetType{ElemType: types.StringType},
                    "description":  types.StringType,
                    "functions":  types.SetType{ElemType: types.StringType},
                    "hosts":  types.SetType{ElemType: types.StringType},
                    "images":  types.SetType{ElemType: types.StringType},
                    "labels":  types.SetType{ElemType: types.StringType},
                    "modified": types.StringType,
                    "name":         types.StringType,
                    "namespaces":  types.SetType{ElemType: types.StringType},
                    "owner":        types.StringType,
                    "prisma":       types.BoolType,
                    "system":       types.BoolType,
                },
            },
            []attr.Value{
                c,
            },
        )
        req.PlanValue = c2

        return
    } 

    resp.PlanValue = req.StateValue
}

func UseStateForUnknownSet() planmodifier.Set {
    return useStateForUnknownSetModifier{} 
}

// object plan modifier
type useStateForUnknownObjectModifier struct {}

func (m useStateForUnknownObjectModifier) Description(_ context.Context) string {
    return ""
}

func (m useStateForUnknownObjectModifier) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useStateForUnknownObjectModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    fmt.Println("***********************")
    fmt.Println("entering PlanModifyObject")
    fmt.Println("***********************")
    if req.ConfigValue.IsUnknown() {
        fmt.Println("***********************")
        fmt.Println("entering config isUnknown block of PlanModifyObject")
        fmt.Println("***********************")
        //resp.PlanValue = types.StringValue("") 
        return
    }

    if req.PlanValue.IsUnknown() {
        fmt.Println("***********************")
        fmt.Println("entering plan isUnknown block of PlanModifyObject")
        fmt.Println("***********************")
        return
    }

    resp.PlanValue = req.StateValue
}

func UseStateForUnknownObject() planmodifier.Object {
    return useStateForUnknownObjectModifier{} 
}

// string plan modifier
func UsePlanForUnknownString() planmodifier.String {
    return &usePlanForUnknownStringModifier{}
}

type usePlanForUnknownStringModifier struct {}

func (m *usePlanForUnknownStringModifier) Description(_ context.Context) string {
    return ""
}

func (m *usePlanForUnknownStringModifier) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *usePlanForUnknownStringModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if !req.PlanValue.IsUnknown() {
        return
    }

    resp.PlanValue = req.StateValue
}

func (r *HostCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    //c := NewCollectionResource()
    //cSchemaRequest := resource.SchemaRequest{}
    //cSchemaResponse := &resource.SchemaResponse{}
    //c.Schema(ctx, cSchemaRequest, cSchemaResponse)

    //fmt.Println(cSchemaResponse)
    //fmt.Println(reflect.TypeOf(cSchemaResponse))
    //fmt.Println(*cSchemaResponse)
    resp.Schema = schema.Schema{
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
            "rules": schema.SetNestedAttribute{
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
                        "collections": schema.SetNestedAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            //Required: true,
                            PlanModifiers: []planmodifier.Set{
                                //setplanmodifier.UseStateForUnknown(),
                                UseStateForUnknownSet(),
                            },
                            //PlanModifiers: []planmodifier.Set{
                            //   planmodifiers.SetDefaultValue(schema.SetNestedAttribute, nil), 
                            //},
                            //Default: objectdefault.StaticValue(
                            //    types.ObjectValueMust(
                            //        map[string]attr.Type{
                            //            "account_ids":  types.SetType{ElemType: types.StringType},
                            //            "app_ids":      types.SetType{ElemType: types.StringType},
                            //            "clusters":      types.SetType{ElemType: types.StringType},
                            //            "color":        types.StringType,
                            //            "containers":      types.SetType{ElemType: types.StringType},
                            //            "description":  types.StringType,
                            //            "functions":      types.SetType{ElemType: types.StringType},
                            //            "hosts":      types.SetType{ElemType: types.StringType},
                            //            "images":      types.SetType{ElemType: types.StringType},
                            //            "labels":      types.SetType{ElemType: types.StringType},
                            //            "name":         types.StringType,
                            //            "namespaces":      types.SetType{ElemType: types.StringType},
                            //            "owner":        types.StringType,
                            //            "prisma":       types.BoolType,
                            //            "system":       types.BoolType,
                            //        },
                            //        map[string]attr.Value{
                            //            "account_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "app_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "clusters": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "color": types.StringValue("#FFFFFF"),
                            //            "containers": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "description": types.StringValue(""),
                            //            "functions": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "hosts": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "images": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "labels": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "name": types.StringValue("All"),
                            //            "namespaces": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            //            "owner": types.StringValue(""),
                            //            "prisma": types.BoolValue(true),
                            //            "system": types.BoolValue(true),
                            //        },
                            //    ),
                            //),
                            NestedObject: schema.NestedAttributeObject{
                                //Attributes: cSchemaResponse.Schema.Attributes,
                                Attributes: map[string]schema.Attribute{
                                    "account_ids": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //PlanModifiers: []planmodifier.Set{
                                        //    setplanmodifier.UseStateForUnknown(),
                                        //},
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "app_ids": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "clusters": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "color": schema.StringAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                        //Default: stringdefault.StaticString("#3FA2F7"),
                                    },
                                    "containers": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "description": schema.StringAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                        //Default: stringdefault.StaticString(""),
                                    },
                                    "functions": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "hosts": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "images": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "labels": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "modified": schema.StringAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                    },
                                    "name": schema.StringAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                    },
                                    "namespaces": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        ElementType: types.StringType,
                                        Optional: true,
                                        Computed: true,
                                        //Default: setdefault.StaticValue(
                                        //    types.SetValueMust(
                                        //        types.StringType,
                                        //        []attr.Value{
                                        //            types.StringValue("*"),
                                        //        },
                                        //    ),
                                        //),
                                    },
                                    "owner": schema.StringAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                    },
                                    "prisma": schema.BoolAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                        //Default: booldefault.StaticBool(false),
                                    },
                                    "system": schema.BoolAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                        //Default: booldefault.StaticBool(false),
                                    },
                                },
                            },
                        },
                        //"condition": schema.SingleNestedAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    PlanModifiers: []planmodifier.Object{
                        //        objectplanmodifier.UseStateForUnknown(),
                        //    },
                        //    Attributes: map[string]schema.Attribute{
                        //        "device": schema.StringAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            //Default: stringdefault.StaticString(""),
                        //        },
                        //        "read_only": schema.BoolAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            //Default: booldefault.StaticBool(false),
                        //        },
                        //        "vulnerabilities": schema.ListNestedAttribute{
                        //            MarkdownDescription: "TODO",
                        //            Optional: true,
                        //            Computed: true,
                        //            NestedObject: schema.NestedAttributeObject{
                        //                Attributes: map[string]schema.Attribute{
                        //                    "id": schema.Int32Attribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                    "block": schema.BoolAttribute{
                        //                        MarkdownDescription: "TODO",
                        //                        Optional: true,
                        //                        Computed: true,
                        //                    },
                        //                },
                        //            },
                        //            // TODO: default value
                        //        },
                        //    },
                        //},
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
                        //"effect": schema.StringAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Validators: []validator.String{
                        //        stringvalidator.OneOf("allow", "deny", "block", "alert"),
                        //    },
                        //    Default: stringdefault.StaticString("alert"),
                        //},
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
                        //"notes": schema.StringAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //},
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
                        //"verbose": schema.BoolAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: booldefault.StaticBool(false),
                        //},
                    },
                },
            },
        },
    }
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
    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    //fmt.Println("entering ModifyPlan")
    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    ////fmt.Println(reflect.TypeOf(req.Plan.Raw))

    //var plan HostCompliancePolicyResourceModel
    //diags := req.Plan.Get(ctx, &plan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}


    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    //fmt.Println("starting loop over rules")
    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    //for _, rule := range *plan.Rules {
    //    //if rule.Collections == nil {
    //    if rule.Collections.IsUnknown() {
    //        c := types.SetType{
    //            ElemType: types.ObjectType{
    //                AttrTypes: map[string]attr.Type{
    //                    "account_ids":  types.SetType{ElemType: types.StringType},
    //                    "app_ids":  types.SetType{ElemType: types.StringType},
    //                    "clusters":  types.SetType{ElemType: types.StringType},
    //                    "color":        types.StringType,
    //                    "containers":  types.SetType{ElemType: types.StringType},
    //                    "description":  types.StringType,
    //                    "functions":  types.SetType{ElemType: types.StringType},
    //                    "hosts":  types.SetType{ElemType: types.StringType},
    //                    "images":  types.SetType{ElemType: types.StringType},
    //                    "labels":  types.SetType{ElemType: types.StringType},
    //                    "name":         types.StringType,
    //                    "namespaces":  types.SetType{ElemType: types.StringType},
    //                    "owner":        types.StringType,
    //                    "prisma":       types.BoolType,
    //                    "system":       types.BoolType,
    //                },
    //            }
    //        }
    //       
    //        test, _ := types.SetValueMust(ctx, types.ObjectType{
    //                AttrTypes: map[string]attr.Type{
    //                    "account_ids":  types.SetType{ElemType: types.StringType},
    //                    "app_ids":  types.SetType{ElemType: types.StringType},
    //                    "clusters":  types.SetType{ElemType: types.StringType},
    //                    "color":        types.StringType,
    //                    "containers":  types.SetType{ElemType: types.StringType},
    //                    "description":  types.StringType,
    //                    "functions":  types.SetType{ElemType: types.StringType},
    //                    "hosts":  types.SetType{ElemType: types.StringType},
    //                    "images":  types.SetType{ElemType: types.StringType},
    //                    "labels":  types.SetType{ElemType: types.StringType},
    //                    "name":         types.StringType,
    //                    "namespaces":  types.SetType{ElemType: types.StringType},
    //                    "owner":        types.StringType,
    //                    "prisma":       types.BoolType,
    //                    "system":       types.BoolType,
    //                },
    //                []attr.Value{types.ObjectValueMust
    //            },



    //        //wildcardSetValue, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("*")})
    //        //collectionSetObject, _ := types.SetValue(types.DynamicType, []attr.Value{collectionObject})
    //        //collectionObject, _ := types.ObjectValue(
    //        //    map[string]attr.Type{
    //        //        "account_ids":  types.DynamicType,
    //        //        "app_ids":      types.DynamicType,
    //        //        "clusters":     types.DynamicType,
    //        //        "color":        types.StringType,
    //        //        "containers":   types.DynamicType,
    //        //        "description":  types.StringType,
    //        //        "functions":    types.DynamicType,
    //        //        "hosts":        types.DynamicType,
    //        //        "images":       types.DynamicType,
    //        //        "labels":       types.DynamicType,
    //        //        "name":         types.StringType,
    //        //        "namespaces":   types.DynamicType,
    //        //        "owner":        types.StringType,
    //        //        "prisma":       types.BoolType,
    //        //        "system":       types.BoolType,
    //        //    },
    //        //    map[string]attr.Value{
    //        //        "account_ids": wildcardSetValue,
    //        //        "app_ids": wildcardSetValue,
    //        //        "clusters": wildcardSetValue,
    //        //        "color": types.StringValue("#FFFFFF"),
    //        //        "containers": wildcardSetValue,
    //        //        "description": types.StringValue(""),
    //        //        "functions": wildcardSetValue,
    //        //        "hosts": wildcardSetValue,
    //        //        "images": wildcardSetValue,
    //        //        "labels": wildcardSetValue,
    //        //        "name": types.StringValue("All"),
    //        //        "namespaces": wildcardSetValue,
    //        //        "owner": types.StringValue(""),
    //        //        "prisma": types.BoolValue(true),
    //        //        "system": types.BoolValue(true),
    //        //    },
    //        //)
    //        //fmt.Printf("%+v\n", collectionObject)
    //        //rule.Collections = collectionSetObject
    //        //fmt.Println(reflect.TypeOf(collectionSetObject))
    //        //fmt.Printf("%+v\n", collectionSetObject)
    //    }
    //}

    //resp.Plan.Set(ctx, &plan)

    //var respPlan HostCompliancePolicyResourceModel
    //diags = resp.Plan.Get(ctx, &respPlan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //fmt.Printf("%+v\n", respPlan)
    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    //fmt.Println("exiting ModifyPlan")
    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
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
    fmt.Println("starting value assignment")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    policy := policyAPI.HostCompliancePolicy{
        Id: plan.Id.ValueString(),
        PolicyType: plan.PolicyType.ValueString(),
    }

    rules := []policyAPI.HostCompliancePolicyRule{}
    for _, planRule := range *plan.Rules {
        collections := []collectionAPI.Collection{}
        collections = append(collections, collectionAPI.Collection{
            AccountIDs: []string{"*"},
            AppIDs: []string{"*"},
            Clusters: []string{"*"},
            Color: "#3FA2F7",
            Containers: []string{"*"},
            Description: "System - all resources collection",
            Functions: []string{"*"},
            Hosts: []string{"*"},
            Images: []string{"*"},
            Labels: []string{"*"},
            Name: "All",
            Namespaces: []string{"*"},
            Prisma: true,
            System: true,
        })

        rules = append(rules, policyAPI.HostCompliancePolicyRule{
            Name: planRule.Name.ValueString(), 
            Collections: collections,
        })
    }
    policy.Rules = &rules
    //    if planRule.Action.IsUnknown() {
    //        planRule.Action = types.SetValueMust(
    //            types.StringType,
    //            []attr.Value{
    //                types.StringValue("*"),
    //            },
    //        )
    //    }

    //    if planRule.Collections.IsUnknown() {
    //        //planRule.Collections = types.ObjectValueMust(
    //        planRule.Collections = types.SetValueMust(
    //            types.ObjectType{ AttrTypes: map[string]attr.Type{
    //              "account_ids":  types.SetType{ElemType: types.StringType},
    //                "app_ids":      types.SetType{ElemType: types.StringType},
    //                "clusters":      types.SetType{ElemType: types.StringType},
    //                "color":        types.StringType,
    //                "containers":      types.SetType{ElemType: types.StringType},
    //                "description":  types.StringType,
    //                "functions":      types.SetType{ElemType: types.StringType},
    //                "hosts":      types.SetType{ElemType: types.StringType},
    //                "images":      types.SetType{ElemType: types.StringType},
    //                "labels":      types.SetType{ElemType: types.StringType},
    //                "name":         types.StringType,
    //                "namespaces":      types.SetType{ElemType: types.StringType},
    //                "owner":        types.StringType,
    //                "prisma":       types.BoolType,
    //                "system":       types.BoolType,
    //            }},
    //            []attr.Value{
    //                types.ObjectValueMust(
    //                    map[string]attr.Type{
    //                        "account_ids":  types.SetType{ElemType: types.StringType},
    //                        "app_ids":      types.SetType{ElemType: types.StringType},
    //                        "clusters":      types.SetType{ElemType: types.StringType},
    //                        "color":        types.StringType,
    //                        "containers":      types.SetType{ElemType: types.StringType},
    //                        "description":  types.StringType,
    //                        "functions":      types.SetType{ElemType: types.StringType},
    //                        "hosts":      types.SetType{ElemType: types.StringType},
    //                        "images":      types.SetType{ElemType: types.StringType},
    //                        "labels":      types.SetType{ElemType: types.StringType},
    //                        "name":         types.StringType,
    //                        "namespaces":      types.SetType{ElemType: types.StringType},
    //                        "owner":        types.StringType,
    //                        "prisma":       types.BoolType,
    //                        "system":       types.BoolType,
    //                    },
    //                    map[string]attr.Value{
    //                        "account_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "app_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "clusters": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "color": types.StringValue("#FFFFFF"),
    //                        "containers": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "description": types.StringValue(""),
    //                        "functions": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "hosts": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "images": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "labels": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "name": types.StringValue("All"),
    //                        "namespaces": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
    //                        "owner": types.StringValue(""),
    //                        "prisma": types.BoolValue(true),
    //                        "system": types.BoolValue(true),
    //                    },
    //                ),
    //            },
    //        )
    //    }
    //}
    
    //rules := make([]policyAPI.HostCompliancePolicyRule, 0, len(*plan.Rules))
    ////for _, planRule := range *plan.Rules {
    ////    rule := policyAPI.HostCompliancePolicyRule{}

    ////    if !planRule.Action.IsUnknown() {
    ////        action := make([]string, 0, len(planRule.Action.Elements()))
    ////        diags = planRule.Action.ElementsAs(ctx, &action, false)
    ////        if diags.HasError() {
    ////            return 
    ////        }
    ////        rule.Action = action
    ////    } else {
    ////        rule.Action = []string{"*"}
    ////    }
    ////}
    //policy.Rules = rules
    

    //// Generate API request body from plan
    ////policy, diags := schemaToPolicy(ctx, &plan, username)
    //policy, diags := schemaToPolicy(ctx, &plan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    // Create new host compliance policy 
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("creating policy resource with payload:")
    fmt.Printf("%+v\n", policy)
    r1 := *policy.Rules
    fmt.Printf("%+v\n", r1)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    err := policyAPI.CreateHostCompliancePolicy(*r.client, policy)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Host Compliance Policy resource", 
            "Failed to create host compliance policy: " + err.Error(),
        )
        return
	}

    // Retrieve newly created host compliance policy 
    response, err := policyAPI.GetHostCompliancePolicy(*r.client)
    createdPolicy, diags := policyToSchema(ctx, *response)
    if diags.HasError() {
        fmt.Println(diags)
        return
    }

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("createdPolicy in Create():")
    fmt.Println(reflect.TypeOf(createdPolicy))
    fmt.Printf("%+v\n", createdPolicy)
    fmt.Printf("%+v\n", *createdPolicy.Rules)
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
    fmt.Println("oh shit we're in Read")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")

    // Get current state
    var state HostCompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get collection value from Prisma Cloud
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
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
  
    // Overwrite state values with Prisma Cloud data
    //state, diags = collectionToSchema(ctx, *collection) 
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}
    stateSchema := HostCompliancePolicyResourceModel{
        Id: types.StringValue(policy.Id),
        PolicyType: types.StringValue(policy.PolicyType),
    }

    // Set refreshed state
    diags = resp.State.Set(ctx, &stateSchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *HostCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state CollectionResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan CollectionResourceModel
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    collection, diags := schemaToCollection(ctx, &plan)

    // Update exsting collection 
	err := collectionAPI.UpdateCollection(*r.client, state.Name.ValueString(), collection)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating Collection resource", 
            "Failed to update collection: " + err.Error(),
        )
        return
	}

    // Fetch updated collection from Prisma Cloud
    updatedCollection, err := collectionAPI.GetCollection(*r.client, plan.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating Collection resource", 
            "Failed to read name" + plan.Name.ValueString()  + ": " + err.Error(),
        )
        return
    }

    // Update plan values from collection data
    plan, diags = collectionToSchema(ctx, *updatedCollection)
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
	var state CollectionResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing collection 
    collection := state.Name.ValueString()
    err := collectionAPI.DeleteCollection(*r.client, collection)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Collection resource", 
            "Failed to delete collection: " + err.Error(),
        )
        return
	}
}

// TODO: Define ImportState to work properly with this resource
func (r *HostCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func schemaToPolicy(ctx context.Context, plan *HostCompliancePolicyResourceModel/*, username types.String*/) (policyAPI.HostCompliancePolicy, diag.Diagnostics) {
    var diags diag.Diagnostics

    policy := policyAPI.HostCompliancePolicy{
        Id: plan.Id.ValueString(),
        PolicyType: plan.PolicyType.ValueString(),
    }

    //if plan.Rules != nil && len(*plan.Rules) > 0 {
    //    policy.Rules, diags = ruleSchemaToPolicy(ctx, *plan.Rules/*, username*/) 
    //    if diags.HasError() {
    //        return policy, diags 
    //    }
    //} else {
    //    policy.Rules = nil // TODO: dunno if this works 
    //}

    return policy, diags
}

func ruleSchemaToPolicy(ctx context.Context, planRules []HostCompliancePolicyRuleResourceModel/*, username types.String*/) ([]policyAPI.HostCompliancePolicyRule, diag.Diagnostics) {
    _ = reflect.TypeOf(ctx)

    var diags diag.Diagnostics
    rules := []policyAPI.HostCompliancePolicyRule{}
    for _, ruleObject := range planRules {
        rule := policyAPI.HostCompliancePolicyRule{
            //ReportAllCompliance: ruleObject.ReportAllCompliance.ValueBool(),
            //BlockMessage: ruleObject.BlockMessage.ValueString(),
            //CreatePR: ruleObject.CreatePR.ValueBool(),
            ////Disabled: ruleObject.Disabled.ValueBool(),
            //Effect: ruleObject.Effect.ValueString(),
            ////ExcludeBaseImageVulns: ruleObject.ExcludeBaseImageVulns.ValueBool(),
            //GraceDays: ruleObject.GraceDays.ValueInt32(),
            //Modified: ruleObject.Modified.ValueString(),
            Name: ruleObject.Name.ValueString(),
            //Notes: ruleObject.Notes.ValueString(),
            //OnlyFixed: ruleObject.OnlyFixed.ValueBool(),
            ////Owner: ruleObject.Owner.ValueString(),
            ////PreviousName: ruleObject.PreviousName.ValueString(),
            //Verbose: ruleObject.Verbose.ValueBool(),
        }

        //if !ruleObject.Action.IsNull() {
        //    action := make([]string, 0, len(ruleObject.Action.Elements()))
        //    diags = ruleObject.Action.ElementsAs(ctx, &action, false)
        //    if diags.HasError() {
        //        return rules, diags
        //    }
        //    rule.Action = action
        //} else {
        //    rule.Action = []string{"*"}
        //}

        //if len(ruleObject.Modified.ValueString()) == 0 {
        //    rule.Modified = time.Now().Format("2006-01-02T15:04:05.000Z")
        //}

        //if ruleObject.License.IsNull() {
        //   rule.License = policyAPI.HostCompliancePolicyRuleLicense{}
        //} /*else {*/
        //   //rule.License = policyAPI.HostCompliancePolicyRuleLicense{
        ////}

        // TODO: refine this logic to populate Owner with the value in config, if it exists
        //if ruleObject.Owner == types.StringValue("") {
        //    rule.Owner = username.ValueString()
        //}

        //if ruleObject.AlertThreshold != nil {
        //    alertThreshold := policyAPI.HostCompliancePolicyRuleAlertThreshold {
        //        Disabled: ruleObject.AlertThreshold.Disabled.ValueBool(),
        //        Value: ruleObject.AlertThreshold.Value.ValueInt32(),
        //    }
        //    rule.AlertThreshold = alertThreshold
        //}

        //if ruleObject.BlockThreshold != nil {
        //    blockThreshold := policyAPI.HostCompliancePolicyRuleBlockThreshold {
        //        Enabled: ruleObject.BlockThreshold.Enabled.ValueBool(),
        //        Value: ruleObject.BlockThreshold.Value.ValueInt32(),
        //    }
        //    rule.BlockThreshold = blockThreshold
        //}

        //if ruleObject.Collections.IsUnknown() {
        //    collections := make([]collectionAPI.Collection, 0, 1)
        //    wildcardSetValue, _ := types.SetValue(types.String, []attr.Value{types.StringValue("*")})
        //    collections = append(collections, collectionAPI.Collection {
        //        //AccountIDs: []string{"*"},
        //        //AppIDs: []string{"*"},
        //        //Clusters: []string{"*"},
        //        AccountIDs: wildcardSetValue,
        //        AppIDs: wildcardSetValue,
        //        Clusters: wildcardSetValue,
        //        Color: "#3FA2F7",
        //        //Containers: []string{"*"},
        //        Containers: wildcardSetValue,
        //        Description: "System - all resources collection",
        //        //Functions: []string{"*"},
        //        //Hosts: []string{"*"},
        //        //Images: []string{"*"},
        //        //Labels: []string{"*"},
        //        Functions: wildcardSetValue,
        //        Hosts: wildcardSetValue,
        //        Images: wildcardSetValue,
        //        Labels: wildcardSetValue,
        //        Name: "All",
        //        //Namespaces: []string{"*"},
        //        Namespaces: wildcardSetValue,
        //        Prisma: true,
        //        System: true,
        //    })
        //                    //        []attr.Value{
        //                    //            types.StringValue("*"),
        //                    //        },
        //    //diags = ruleObject.Collections.ElementsAs(ctx, &collections, false)
        //    //fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%")
        //    //fmt.Println(collections)
        //    //fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%")
        //    //for _, collectionObject := range ruleObject.Collections.ElementsAs() {
        //    //    collection, diags := schemaToCollection(ctx, &collectionObject)
        //    //    if diags.HasError() {
        //    //        return rules, diags
        //    //    }
        //    //    collections = append(collections, collection)
        //    //}
        //    rule.Collections = collections
        //}
        //for _, collectionObject := range *ruleObject.Collections {
        //    fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
        //    collection, diags := schemaToCollection(ctx, &collectionObject) 
        //    if diags.HasError() {
        //        return rules, diags
        //    }
        //    collections = append(collections, collection)
        //}
        //rule.Collections = collections
        //if ruleObject.Collections != nil {
        ////if !ruleObject.Collections.IsNull() {
        ////    for _, collectionObject := range ruleObject.Collections {
        ////        //fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%")
        ////        //fmt.Println(reflect.TypeOf(collectionObject))
        ////        //fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%")
        ////        collection, diags := schemaToCollection(ctx, &collectionObject) 
        ////        if diags.HasError() {
        ////            return rules, diags
        ////        }
        ////        collections = append(collections, collection)
        ////    }
        ////    rule.Collections = collections
        ////} else {
        ////    rule.Collections = []collectionAPI.Collection{
        ////        collectionAPI.Collection{
        ////        },
        ////    }
        ////}
        //collections := []collectionAPI.Collection{}
        //if ruleObject.Collections.IsNull() {
        //    collections = append(collections, collectionAPI.Collection {
        //        AccountIDs: []string{"*"},
        //        AppIDs: []string{"*"},
        //        Clusters: []string{"*"},
        //        Color: "#3FA2F7",
        //        Containers: []string{"*"},
        //        Description: "System - all resources collection",
        //        Functions: []string{"*"},
        //        Hosts: []string{"*"},
        //        Images: []string{"*"},
        //        Labels: []string{"*"},
        //        Name: "All",
        //        Namespaces: []string{"*"},
        //        Prisma: true,
        //        System: true,
        //    })
        //    rule.Collections = collections 
        //}

        //if ruleObject.CVERules != nil {
        //    cveRules := []policyAPI.HostCompliancePolicyRuleCVERule{}
        //    for _, cveRuleObject := range *ruleObject.CVERules {
        //        cveRule := policyAPI.HostCompliancePolicyRuleCVERule{
        //            Description: cveRuleObject.Description.ValueString(),
        //            Effect: cveRuleObject.Effect.ValueString(),
        //            Expiration: policyAPI.HostCompliancePolicyRuleExpiration{
        //                Date: cveRuleObject.Expiration.Date.ValueString(),
        //                Enabled: cveRuleObject.Expiration.Enabled.ValueBool(),
        //            },
        //            Id: cveRuleObject.Id.ValueString(),
        //        }
        //        cveRules = append(cveRules, cveRule)
        //    }
        //    rule.CVERules = cveRules
        //}

        //if ruleObject.GraceDaysPolicy == nil {
        //    fmt.Println("aaah! aaaaaaaahhhh!!!!")
        //    rule.GraceDaysPolicy = policyAPI.HostCompliancePolicyRuleGraceDaysPolicy{}
        //} else {
        //    fmt.Println("... oh...")
        //}

        //if ruleObject.PkgTypesThresholds != nil {
        //    pkgTypesThresholds := []policyAPI.HostCompliancePolicyRulePkgTypesThresholds{}
        //    for _, pkgTypesThresholdsObject := range *ruleObject.PkgTypesThresholds {
        //        pkgTypesThreshold := policyAPI.HostCompliancePolicyRulePkgTypesThresholds{
        //            AlertThreshold: policyAPI.HostCompliancePolicyRuleAlertThreshold{
        //                Disabled: pkgTypesThresholdsObject.AlertThreshold.Disabled.ValueBool(),
        //                Value: pkgTypesThresholdsObject.AlertThreshold.Value.ValueInt32(),
        //            },
        //            BlockThreshold: policyAPI.HostCompliancePolicyRuleBlockThreshold{
        //                Enabled: pkgTypesThresholdsObject.BlockThreshold.Enabled.ValueBool(),
        //                Value: pkgTypesThresholdsObject.BlockThreshold.Value.ValueInt32(),
        //            },
        //            Type: pkgTypesThresholdsObject.Type.ValueString(),
        //        }
        //        pkgTypesThresholds = append(pkgTypesThresholds, pkgTypesThreshold)
        //    }
        //    rule.PkgTypesThresholds = pkgTypesThresholds
        //}

        //if ruleObject.RiskFactorsEffects != nil {
        //    riskFactorsEffects := []policyAPI.HostCompliancePolicyRuleRiskFactorsEffect{}
        //    for _, riskFactorsEffectObject := range *ruleObject.RiskFactorsEffects {
        //        riskFactorsEffect := policyAPI.HostCompliancePolicyRuleRiskFactorsEffect{
        //            Effect: riskFactorsEffectObject.Effect.ValueString(),
        //            RiskFactor: riskFactorsEffectObject.RiskFactor.ValueString(),
        //        }
        //        riskFactorsEffects = append(riskFactorsEffects, riskFactorsEffect)
        //    }
        //    rule.RiskFactorsEffects = riskFactorsEffects
        //}

        //if ruleObject.Tags != nil {
        //    tags := []policyAPI.HostCompliancePolicyRuleTag{}
        //    for _, tagObject := range *ruleObject.Tags {
        //        tag := policyAPI.HostCompliancePolicyRuleTag{
        //            Name: tagObject.Name.ValueString(),
        //            Description: tagObject.Description.ValueString(),
        //            Effect: tagObject.Effect.ValueString(),
        //            Expiration: policyAPI.HostCompliancePolicyRuleExpiration{
        //                Date: tagObject.Expiration.Date.ValueString(),
        //                Enabled: tagObject.Expiration.Enabled.ValueBool(),
        //            },
        //        }
        //        tags = append(tags, tag)
        //    }
        //    rule.Tags = tags
        //}

        rules = append(rules, rule)
    }

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
    fmt.Printf("%+v\n", rules)
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

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("returning from policy collection values: ") 
    fmt.Printf("%+v\n", schema)
    fmt.Printf("%+v\n", *schema.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")

    return schema, diags
}

func policyRulesToSchema(ctx context.Context, rules []policyAPI.HostCompliancePolicyRule) ([]HostCompliancePolicyRuleResourceModel, diag.Diagnostics) {
//func policyRulesToSchema(ctx context.Context, rules []policyAPI.HostCompliancePolicyRule) (types.SetType, diag.Diagnostics) {
    fmt.Println("***********************")
    fmt.Println("entering policyRulesToSchema")
    fmt.Println("***********************")

    var diags diag.Diagnostics

    schemaRules := []HostCompliancePolicyRuleResourceModel{}

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
        }

        if rule.Collections != nil {
            values := []attr.Value{}
            for _, collection := range(rule.Collections) {
                fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
                fmt.Println("iterating over collection with values: ") 
                fmt.Printf("%+v\n", collection)
                //fmt.Println(reflect.TypeOf(collection.AccountIDs))
                //fmt.Println(collection.AccountIDs)
                fmt.Println("#$%#$%#$%#$%#$%#$%#$%")

                //accountIDs := make([]attr.Value, 0, len(collection.AccountIDs))
                //for _, val := range collection.AccountIDs {
                //    accountIDs = append(accountIDs, types.StringValue(val))
                //}
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

                c := types.ObjectValueMust(
                    map[string]attr.Type{
                        "account_ids":  types.SetType{ElemType: types.StringType},
                        "app_ids":  types.SetType{ElemType: types.StringType},
                        "clusters":  types.SetType{ElemType: types.StringType},
                        "color":        types.StringType,
                        "containers":  types.SetType{ElemType: types.StringType},
                        "description":  types.StringType,
                        "functions":  types.SetType{ElemType: types.StringType},
                        "hosts":  types.SetType{ElemType: types.StringType},
                        "images":  types.SetType{ElemType: types.StringType},
                        "labels":  types.SetType{ElemType: types.StringType},
                        "modified": types.StringType,
                        "name":         types.StringType,
                        "namespaces":  types.SetType{ElemType: types.StringType},
                        "owner":        types.StringType,
                        "prisma":       types.BoolType,
                        "system":       types.BoolType,
                    },
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
                        "name": types.StringValue(collection.Name),
                        "namespaces": namespaces,
                        "owner": types.StringValue(collection.Owner),
                        "prisma": types.BoolValue(collection.Prisma),
                        "system": types.BoolValue(collection.System),
                    },
                )
                values = append(values, c)

                //collectionSetValue, diags := types.SetValueFrom(
                //    ctx,
                //    types.ObjectType{
                //        AttrTypes: map[string]attr.Type{
                //            "account_ids":  types.SetType{ElemType: types.StringType},
                //            "app_ids":  types.SetType{ElemType: types.StringType},
                //            "clusters":  types.SetType{ElemType: types.StringType},
                //            "color":        types.StringType,
                //            "containers":  types.SetType{ElemType: types.StringType},
                //            "description":  types.StringType,
                //            "functions":  types.SetType{ElemType: types.StringType},
                //            "hosts":  types.SetType{ElemType: types.StringType},
                //            "images":  types.SetType{ElemType: types.StringType},
                //            "labels":  types.SetType{ElemType: types.StringType},
                //            "name":         types.StringType,
                //            "namespaces":  types.SetType{ElemType: types.StringType},
                //            "owner":        types.StringType,
                //            "prisma":       types.BoolType,
                //            "system":       types.BoolType,
                //        },
                //    },
                //    collection,
                //)
                
                //collectionSchema := collectionToSchema(ctx, collectionSetValue)
            }

            //c2 := types.SetValueMust(
            c2, diags := types.SetValueFrom(
                ctx,
                types.ObjectType{
                    AttrTypes: map[string]attr.Type{
                        "account_ids":  types.SetType{ElemType: types.StringType},
                        "app_ids":  types.SetType{ElemType: types.StringType},
                        "clusters":  types.SetType{ElemType: types.StringType},
                        "color":        types.StringType,
                        "containers":  types.SetType{ElemType: types.StringType},
                        "description":  types.StringType,
                        "functions":  types.SetType{ElemType: types.StringType},
                        "hosts":  types.SetType{ElemType: types.StringType},
                        "images":  types.SetType{ElemType: types.StringType},
                        "labels":  types.SetType{ElemType: types.StringType},
                        "modified": types.StringType,
                        "name":         types.StringType,
                        "namespaces":  types.SetType{ElemType: types.StringType},
                        "owner":        types.StringType,
                        "prisma":       types.BoolType,
                        "system":       types.BoolType,
                    },
                },
                values,
            )

            if diags.HasError() {
                fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
                fmt.Println("error calling SetValueFrom on collection data") 
                fmt.Println(diags)
                fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
                return schemaRules, diags
            }

            schemaRule.Collections = c2
        }

        schemaRules = append(schemaRules, schemaRule)
    }

    fmt.Println("***********************")
    fmt.Println("exiting policyRulesToSchema")
    fmt.Printf("%+v\n", schemaRules)
    fmt.Println("***********************")

    return schemaRules, diags
}
