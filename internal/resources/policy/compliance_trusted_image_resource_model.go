package policy

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"
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
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	//"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ resource.Resource = &TrustedImagesPolicyResource{}
var _ resource.ResourceWithImportState = &TrustedImagesPolicyResource{}
var _ resource.ResourceWithModifyPlan = &TrustedImagesPolicyResource{}

func NewTrustedImagesPolicyResource() resource.Resource {
    return &TrustedImagesPolicyResource{}
}

type TrustedImagesPolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

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

func (r *TrustedImagesPolicyResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "groups": schema.SetNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                        },
                        "images": schema.SetAttribute{
                            ElementType: types.StringType,
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
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
                        "owner": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                        },
                        "previous_name": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                        },
                    },
                },
            },
            "policy": schema.SingleNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Attributes: map[string]schema.Attribute{
                    "id": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: stringdefault.StaticString("trust"),
                    },
                    "enabled": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: booldefault.StaticBool(true),
                    },
                    "rules": schema.SetNestedAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        PlanModifiers: []planmodifier.Set{
                            planmodifiers.UseDefaultForUnknownTrustedImageRules(),
                        },
                        NestedObject: schema.NestedAttributeObject{
                            Attributes: map[string]schema.Attribute{
                                //"action": schema.SetAttribute{
                                //    ElementType: types.StringType,
                                //    MarkdownDescription: "TODO",
                                //    Optional: true,
                                //    Computed: true,
                                //},
                                "allowed_groups": schema.SetAttribute{
                                    ElementType: types.StringType,
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
                                    // TODO: does this need to be a modifier or can we use Default?
                                    PlanModifiers: []planmodifier.Set{
                                        planmodifiers.UseEmptySetForUnknown(),
                                    },
                                },
                                "collections": schema.SetAttribute{
                                    ElementType: types.StringType,
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
                                    // TODO: add validator to ensure this is not empty
                                    PlanModifiers: []planmodifier.Set{
                                        planmodifiers.UseEmptySetForUnknown(),
                                    },
                                },
                                "denied_groups": schema.SetAttribute{
                                    ElementType: types.StringType,
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
                                    // TODO: does this need to be a modifier or can we use Default?
                                    PlanModifiers: []planmodifier.Set{
                                        planmodifiers.UseEmptySetForUnknown(),
                                    },
                                },
                                "effect": schema.StringAttribute{
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
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
                                "notes": schema.StringAttribute{
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
                                    PlanModifiers: []planmodifier.String{
                                        planmodifiers.UseEmptyStringForUnknown(),
                                    },
                                },
                                "owner": schema.StringAttribute{
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
                                },
                                "previous_name": schema.StringAttribute{
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                    Computed: true,
                                    PlanModifiers: []planmodifier.String{
                                        planmodifiers.UseEmptyStringForUnknown(),
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }
}
