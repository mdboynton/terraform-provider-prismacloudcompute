package policy

import (
    //"context"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	//"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ resource.Resource = &ApplicationControlPolicyResource{}
var _ resource.ResourceWithImportState = &ApplicationControlPolicyResource{}
//var _ resource.ResourceWithModifyPlan = &ApplicationControlPolicyResource{}

func NewApplicationControlPolicyResource() resource.Resource {
    return &ApplicationControlPolicyResource{}
}

type ApplicationControlPolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

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

func (r *ApplicationControlPolicyResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.Int32Attribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            PlanModifiers: []planmodifier.Int32{
                                planmodifiers.UseStateRuleIdForUnknownIfExists(),
                            },
                        },
                        "applications": schema.SetNestedAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "name": schema.StringAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                    },
                                    "allowed_versions": schema.SetAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        //Computed: true,
                                        ElementType: types.SetType{
                                            ElemType: types.StringType,
                                        },
                                    },
                                },
                            },
                        },
                        "description": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            // TODO: validate that this is not empty
                            //Validators: []validator.String{
                            //    stringvalidator.OneOf("ignore", "alert"),
                            //},
                            //PlanModifiers: []planmodifier.String{
                            //    planmodifiers.AllowUnknownEffect(),
                            //},
                        },
                        //"disabled": schema.BoolAttribute{
                        //    MarkdownDescription: "TODO",
                        //    Optional: true,
                        //    Computed: true,
                        //    Default: booldefault.StaticBool(false), 
                        //},
                        "modified": schema.StringAttribute{
                           MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            //Default: stringdefault.StaticString(time.Now().Format("2006-01-02T15:04:05.000Z")),
                            PlanModifiers: []planmodifier.String{
                                planmodifiers.UseEmptyStringForNull(),
                            },
                        },
                        "name": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Required: true,
                            //Optional: true,
                            //Computed: true,
                            //PlanModifiers: []planmodifier.String{
                            //    stringplanmodifier.RequiresReplaceIf(
                            //        func(ctx context.Context, sr planmodifier.StringRequest, rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse) {
                            //            rrifr.RequiresReplace = (sr.PlanValue.ValueString() != sr.StateValue.ValueString())
                            //        },
                            //        "TODO",
                            //        "TODO",
                            //    ),
                            //},
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
                        "previous_name": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.UseStateForUnknown(),
                            },
                        },
                        "severity": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.UseStateForUnknown(),
                            },
                        },
                    },
                },
            },
        },
    }
}
