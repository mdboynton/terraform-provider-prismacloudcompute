package policy

import (
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

var _ resource.Resource = &VmImageCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &VmImageCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &VmImageCompliancePolicyResource{}

func NewVmImageCompliancePolicyResource() resource.Resource {
    return &VmImageCompliancePolicyResource{}
}

type VmImageCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *VmImageCompliancePolicyResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("vmCompliance"),
            },
            "policy_type": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("vmCompliance"),
            },
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Validators: []validator.List{
                    validators.PolicyRuleNameIsUnique("vm compliance"),
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
