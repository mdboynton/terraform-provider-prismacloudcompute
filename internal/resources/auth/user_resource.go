package auth

import (
	"context"
    "fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/auth"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/auth"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema()
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
    var plan models.UserResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

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
            fmt.Sprintf("Failed to create user \"%s\": %s", user.Username, err.Error()),
        )
        return
	}

    // Retrieve newly created user
    createdUser, err := auth.GetUser(*r.client, user.Username)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created User resource", 
            fmt.Sprintf("Failed to retrieve created user \"%s\": %s", user.Username, err.Error()),
        )
        return
    }

    // Convert to schema
    createdUserSchema, diags := userToSchema(ctx, *createdUser, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set state
    diags = resp.State.Set(ctx, createdUserSchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state models.UserResourceModel
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
            fmt.Sprintf("Failed to read user \"%s\": %s", state.Username.ValueString(), err.Error()),
        )
        return
    }
  
    // Convert user value to schema
    state, diags = userToSchema(ctx, *user, state) 
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
    var plan models.UserResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    user, diags := schemaToUser(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

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
    plan, diags = userToSchema(ctx, *updatedUser, plan)
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
	var state models.UserResourceModel
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

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("username"), req, resp)
}

func schemaToUser(ctx context.Context, plan *models.UserResourceModel) (auth.User, diag.Diagnostics) {
    var diags diag.Diagnostics

    user := auth.User{
        AuthType: plan.AuthenticationType.ValueString(),
        Username: plan.Username.ValueString(),
        Password: plan.Password.ValueString(),
        Role: plan.Role.ValueString(),
    }

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

        if len(permissions) > 0 {
            user.Permissions = permissions
        }

	return user, diags 
}

func userToSchema(ctx context.Context, user auth.User, plan models.UserResourceModel) (models.UserResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    schema := models.UserResourceModel{
        AuthenticationType: types.StringValue(user.AuthType),
        Username: types.StringValue(user.Username),
        Password: plan.Password,
        Role: types.StringValue(user.Role),
    }
   
    if user.Permissions != nil {
        permissions := []models.UserPermissionsResourceModel{}
        for _, permissionsObject := range user.Permissions {
            collections, diags := types.ListValueFrom(ctx, types.StringType, permissionsObject.Collections)
            if diags.HasError() {
                return schema, diags
            }

            permissions = append(permissions, models.UserPermissionsResourceModel{
                Project: types.StringValue(permissionsObject.Project),
                Collections: collections,
            })
        }
        schema.Permissions = &permissions
    }

    return schema, diags
}
