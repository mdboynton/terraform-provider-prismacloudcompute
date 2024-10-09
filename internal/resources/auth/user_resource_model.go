package auth

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
    //"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
    return &UserResource{}
}

type UserResource struct {
    client *api.PrismaCloudComputeAPIClient
}

type UserResourceModel struct {
    AuthenticationType types.String                    `tfsdk:"authentication_type"`
    Username           types.String                    `tfsdk:"username"`
    Password           types.String                    `tfsdk:"password"`
    Role               types.String                    `tfsdk:"role"`
    Permissions        *[]UserPermissionsResourceModel `tfsdk:"permissions"`
}

type UserPermissionsResourceModel struct {
    Project     types.String `tfsdk:"project"`
    Collections types.List   `tfsdk:"collections"`
}


func (r *UserResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "authentication_type": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
            },
            "username": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "password": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
                Sensitive: true,
            },
            "role": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplaceIf(
                        func(
                            ctx context.Context,
                            sr planmodifier.StringRequest,
                            rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse,
                        ) {
                            rrifr.RequiresReplace = (sr.PlanValue.ValueString() != sr.StateValue.ValueString()) && (sr.PlanValue.ValueString() == "admin" || sr.PlanValue.ValueString() == "operator")
                        },
                        "TODO",
                        "TODO",
                    ),
                },
            },
            "permissions": schema.SetNestedAttribute{
                NestedObject: schema.NestedAttributeObject {
                    Attributes: map[string]schema.Attribute{
                        "project": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Required: true,
                        },
                        "collections": schema.ListAttribute{
                            ElementType: types.StringType,
                            MarkdownDescription: "TODO",
                            Required: true,
                        },
                    },
                },
                Optional: true,
            },
        },
    }
}
