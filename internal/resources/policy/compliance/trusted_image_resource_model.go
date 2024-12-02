package policy

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"

    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
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
                        Default: stringdefault.StaticString(policyAPI.PolicyTypeComplianceTrustedImages),
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
                                    Default: setdefault.StaticValue(
                                        types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("All") }),
                                    ),
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
