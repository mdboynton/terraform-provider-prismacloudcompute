package auth

import (
	//"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    //"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithImportState = &RoleResource{}
var _ resource.ResourceWithModifyPlan = &RoleResource{}

func NewRoleResource() resource.Resource {
    return &RoleResource{}
}

type RoleResource struct {
    client *api.PrismaCloudComputeAPIClient
}

type RoleResourceModel struct {
    Name            types.String                    `tfsdk:"name"`
    Description     types.String                    `tfsdk:"description"`
    System          types.Bool                      `tfsdk:"system"`
    Permissions     *[]RolePermissionResourceModel `tfsdk:"permissions"`
}

type RolePermissionResourceModel struct {
    Name        types.String    `tfsdk:"name"`
    ReadWrite   types.Bool      `tfsdk:"read_write"`
}

func (r *RoleResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
                //PlanModifiers: []planmodifier.String{
                //    stringplanmodifier.RequiresReplace(),
                //},
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
            },
            "system": schema.BoolAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
            },
            "permissions": schema.SetNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Set{
                    planmodifiers.UseEmptySetForUnknownRolePermissions(),
                    //planmodifiers.UseStateForNullRolePermissions(),
                },
                // TODO: figure out how to make it so that changes to permissions dont appear as
                // the previous value being removed/set to null and the new value being created as
                // a separate object
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "name": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            //Required: true,
                            Optional: true,
                            Computed: true,
                            //PlanModifiers: []planmodifier.String{
                            //    planmodifiers.UseStateForUnknownRolePermissionName(),
                            //},
                        },
                        "read_write": schema.BoolAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            //PlanModifiers: []planmodifier.Bool{
                            //    planmodifiers.UseStateForUnknownRolePermissionReadWrite(),
                            //},
                        },
                    },
                },
            },
        },
    }
}
