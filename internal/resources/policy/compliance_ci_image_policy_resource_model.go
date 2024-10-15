package policy

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

    "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ resource.Resource = &CiImageCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &CiImageCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &CiImageCompliancePolicyResource{}

func NewCiImageCompliancePolicyResource() resource.Resource {
    return &CiImageCompliancePolicyResource{}
}

type CiImageCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *CiImageCompliancePolicyResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("ciImagesCompliance"),
            },
            "policy_type": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("ciImagesCompliance"),
            },
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Validators: []validator.List{
                    validators.PolicyRuleNameIsUnique("CI image compliance"),
                },
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "order": schema.Int32Attribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                        },
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
                                },
                            },
                        },
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
                        "owner": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.UseStateForUnknown(),
                            },
                        },
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

func (r *CiImageCompliancePolicyResource) GetCollectionsSchema() schema.ListNestedAttribute {
    return schema.ListNestedAttribute{
        MarkdownDescription: "TODO",
        Optional: true,
        Computed: true,
        // TODO: see if we can omit all but the name field and get away with it
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