package auth

import (
	//"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
    return &UserResource{}
}

type UserResource struct {
    client *api.PrismaCloudComputeAPIClient
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
                    planmodifiers.ReplaceIfRoleNameChanged(),
                },
            },
            "permissions": schema.SetNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Validators: []validator.Set{
                    validators.UserPermissionsNotConfiguredWithAdminRoles(),
                },
                PlanModifiers: []planmodifier.Set{
                    planmodifiers.UseDefaultForUnknownUserPermissions(),
                },
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
            },
        },
    }
}
