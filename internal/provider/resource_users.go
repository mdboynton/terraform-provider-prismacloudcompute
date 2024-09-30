package provider

import (
	"context"
    "fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/auth"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
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

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
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

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan UserResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Handle state changes that would cause API errors
    if (plan.Role.ValueString() == "admin" || plan.Role.ValueString() == "operator") && plan.Permissions != nil {
        resp.Diagnostics.AddError(
            "Invalid Resource Configuration",
            "Users with role 'admin' or 'operator' cannot have assigned permissions",
        )
        return
    }

    // TODO: dont think we need this, since sending an empty array for permissions in this case seems to work 
    //if (plan.Role.ValueString() != "admin" && plan.Role.ValueString() != "operator") && plan.Permissions == nil {
    //    resp.Diagnostics.AddError(
    //        "Invalid Resource Configuration",
    //        "Users with role other than 'admin' or 'operator' must have assigned permissions",
    //    )
    //    return
    //}

    // Generate API request body from plan
    user, diags := schemaToUser(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new user
    _, err := auth.CreateUser(*r.client, user)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating User resource", 
            "Failed to create user: " + err.Error(),
        )
        return
	}

    // TODO: retrieve newly created resource and use that data to populate
    //       state instead of using plan data (see below)

    // Set state to plan data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state UserResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get user value from Prisma Cloud
    user, err := auth.GetUser(*r.client, state.Username.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading User resource", 
            "Failed to read username " + state.Username.ValueString()  + ": " + err.Error(),
        )
        return
    }
  
    // Convert user value to schema
    state, diags = userToSchema(ctx, *user) 
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Retrieve values from plan
    var plan UserResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    user, diags := schemaToUser(ctx, &plan)

    // Update existing user
	err := auth.UpdateUser(*r.client, user)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating User resource", 
            "Failed to update user: " + err.Error(),
        )
        return
	}

    // Fetch updated user from Prisma Cloud
    updatedUser, err := auth.GetUser(*r.client, plan.Username.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating User resource", 
            "Failed to read username " + plan.Username.ValueString()  + ": " + err.Error(),
        )
        return
    }

    // Convert updated user to schema
    plan, diags = userToSchema(ctx, *updatedUser)
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

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state UserResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing user
    user := state.Username.ValueString()
    err := auth.DeleteUser(*r.client, user)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting User resource", 
            "Failed to delete user: " + err.Error(),
        )
        return
	}
}

// TODO: Define ImportState to work properly with this resource
func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func schemaToUser(ctx context.Context, plan *UserResourceModel) (auth.User, diag.Diagnostics) {
    var diags diag.Diagnostics

    user := auth.User{
        AuthType: plan.AuthenticationType.ValueString(),
        Username: plan.Username.ValueString(),
        Password: plan.Password.ValueString(),
        Role: plan.Role.ValueString(),
    }

    if plan.Permissions != nil {
        permissions := []auth.UserPermission{}
        for _, permissionsObject := range *plan.Permissions {
            collections := make([]string, 0, len(permissionsObject.Collections.Elements()))
            diags = permissionsObject.Collections.ElementsAs(ctx, &collections, false)
            if diags.HasError() {
                return user, diags
            }
            permissions = append(permissions, auth.UserPermission{
                Project: permissionsObject.Project.ValueString(),
                Collections: collections,
            })
        }
        user.Permissions = permissions
    }

	return user, diags 
}

func userToSchema(ctx context.Context, user auth.User) (UserResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    schema := UserResourceModel{
        AuthenticationType: types.StringValue(user.AuthType),
        Username: types.StringValue(user.Username),
        Password: types.StringValue(user.Password),
        Role: types.StringValue(user.Role),
    }

    if user.Permissions != nil {
        permissions := []UserPermissionsResourceModel{}
        for _, permissionsObject := range user.Permissions {
            collections, diags := types.ListValueFrom(ctx, types.StringType, permissionsObject.Collections)
            if diags.HasError() {
                return schema, diags
            }

            permissions = append(permissions, UserPermissionsResourceModel{
                Project: types.StringValue(permissionsObject.Project),
                Collections: collections,
            })
        }
        schema.Permissions = &permissions
    }

    return schema, diags
}
