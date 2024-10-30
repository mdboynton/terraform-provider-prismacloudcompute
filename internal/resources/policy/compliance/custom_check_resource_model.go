package policy

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var _ resource.Resource = &CustomComplianceCheckResource{}
var _ resource.ResourceWithImportState = &CustomComplianceCheckResource{}
var _ resource.ResourceWithModifyPlan = &CustomComplianceCheckResource{}

func NewCustomComplianceCheckResource() resource.Resource {
    return &CustomComplianceCheckResource{}
}

type CustomComplianceCheckResource struct {
    client *api.PrismaCloudComputeAPIClient
}

type CustomComplianceCheckResourceModel struct {
    Id              types.Int32     `tfsdk:"id"`
    Description     types.String    `tfsdk:"description"`
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
    Script          types.String    `tfsdk:"script"`
    Severity        types.String    `tfsdk:"severity"`
}

func (r *CustomComplianceCheckResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "id": schema.Int32Attribute{
                MarkdownDescription: "TODO",
                Computed: true,
                PlanModifiers: []planmodifier.Int32{
                    int32planmodifier.UseStateForUnknown(),
                },
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
            },
            "owner": schema.StringAttribute{
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
                Required: true,
            },
            "previous_name": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
            },
            "script": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
            },
            "severity": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
                // TODO: default to "high"
                Validators: []validator.String{
                    stringvalidator.OneOf("low", "medium", "high", "critical"),
                },
            },
        },
    }
}
