package system

import (
	"context"
    "fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
    models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/system"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *CollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_collection"
}

func (r *CollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema()
}


func (r *CollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *CollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan models.CollectionResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    collection, diags := schemaToCollection(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new collection 
    err := collectionAPI.CreateCollection(*r.client, collection)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Collection resource", 
            "Failed to create collection: " + err.Error(),
        )
        return
	}

    // Retrieve newly created collection
    response, err := collectionAPI.GetCollection(*r.client, collection.Name)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created Collection resource", 
            "Failed to retrieve created collection: " + err.Error(),
        )
        return
	}
    
    updatedCollection, diags := collectionToSchema(ctx, *response)

    // Set state to collection data
    diags = resp.State.Set(ctx, updatedCollection)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *CollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state models.CollectionResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get collection value from Prisma Cloud
    collection, err := collectionAPI.GetCollection(*r.client, state.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Collection resource", 
            "Failed to read collection name " + state.Name.ValueString()  + ": " + err.Error(),
        )
        return
    }
  
    // Overwrite state values with Prisma Cloud data
    state, diags = collectionToSchema(ctx, *collection) 
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

func (r *CollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state models.CollectionResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan models.CollectionResourceModel
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    collection, diags := schemaToCollection(ctx, &plan)

    // Update existing collection 
	err := collectionAPI.UpdateCollection(*r.client, state.Name.ValueString(), collection)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating Collection resource", 
            "Failed to update collection: " + err.Error(),
        )
        return
	}

    // Fetch updated collection from Prisma Cloud
    updatedCollection, err := collectionAPI.GetCollection(*r.client, plan.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating Collection resource", 
            "Failed to read name" + plan.Name.ValueString()  + ": " + err.Error(),
        )
        return
    }

    // Convert updated collection to schema
    plan, diags = collectionToSchema(ctx, *updatedCollection)
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

func (r *CollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state models.CollectionResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing collection 
    collection := state.Name.ValueString()
    err := collectionAPI.DeleteCollection(*r.client, collection)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Collection resource", 
            "Failed to delete collection: " + err.Error(),
        )
        return
	}
}

func (r *CollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func GetAllCollectionSet() basetypes.SetValue {
    return types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("All") })
}

func schemaToCollection(ctx context.Context, plan *models.CollectionResourceModel) (collectionAPI.Collection, diag.Diagnostics) {
    var diags diag.Diagnostics

    accountIds := make([]string, 0, len(plan.AccountIDs.Elements()))
    diags.Append(plan.AccountIDs.ElementsAs(ctx, &accountIds, false)...)

    appIds := make([]string, 0, len(plan.AppIDs.Elements()))
    diags.Append(plan.AppIDs.ElementsAs(ctx, &appIds, false)...)

    clusters := make([]string, 0, len(plan.Clusters.Elements()))
    diags.Append(plan.Clusters.ElementsAs(ctx, &clusters, false)...)

    containers := make([]string, 0, len(plan.Containers.Elements()))
    diags.Append(plan.Containers.ElementsAs(ctx, &containers, false)...)

    functions := make([]string, 0, len(plan.Functions.Elements()))
    diags.Append(plan.Functions.ElementsAs(ctx, &functions, false)...)

    hosts := make([]string, 0, len(plan.Hosts.Elements()))
    diags.Append(plan.Hosts.ElementsAs(ctx, &hosts, false)...)

    images := make([]string, 0, len(plan.Images.Elements()))
    diags.Append(plan.Images.ElementsAs(ctx, &images, false)...)

    labels := make([]string, 0, len(plan.Labels.Elements()))
    diags.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)

    namespaces := make([]string, 0, len(plan.Namespaces.Elements()))
    diags.Append(plan.Namespaces.ElementsAs(ctx, &namespaces, false)...)

    if diags.HasError() {   
        return collectionAPI.Collection{}, diags
    }

    collection := collectionAPI.Collection{
        AccountIDs: accountIds,
        AppIDs: appIds,
        Clusters: clusters,
        Color: plan.Color.ValueString(),
        Containers: containers,
        Description: plan.Description.ValueString(),
        Functions: functions,
        Hosts: hosts,
        Images: images,
        Labels: labels,
        Name: plan.Name.ValueString(),
        Namespaces: namespaces,
        Prisma: plan.Prisma.ValueBool(),
        System: plan.System.ValueBool(),
    }

	return collection, diags 
}

func collectionToSchema(ctx context.Context, collection collectionAPI.Collection) (models.CollectionResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    accountIds, diags := types.SetValueFrom(ctx, types.StringType, collection.AccountIDs)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    appIds, diags := types.SetValueFrom(ctx, types.StringType, collection.AppIDs)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    clusters, diags := types.SetValueFrom(ctx, types.StringType, collection.Clusters)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    containers, diags := types.SetValueFrom(ctx, types.StringType, collection.Containers)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    functions, diags := types.SetValueFrom(ctx, types.StringType, collection.Functions)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    hosts, diags := types.SetValueFrom(ctx, types.StringType, collection.Hosts)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    images, diags := types.SetValueFrom(ctx, types.StringType, collection.Images)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    labels, diags := types.SetValueFrom(ctx, types.StringType, collection.Labels)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    namespaces, diags := types.SetValueFrom(ctx, types.StringType, collection.Namespaces)
    if diags.HasError() {
        return models.CollectionResourceModel{}, diags
    }

    schema := models.CollectionResourceModel{
        AccountIDs: accountIds,
        AppIDs: appIds,
        Clusters: clusters,
        Color: types.StringValue(collection.Color),
        Containers: containers,
        Description: types.StringValue(collection.Description),
        Functions: functions,
        Hosts: hosts,
        Images: images,
        Labels: labels,
        Modified: types.StringValue(collection.Modified),
        Name: types.StringValue(collection.Name),
        Namespaces: namespaces,
        Owner: types.StringValue(collection.Owner),
        Prisma: types.BoolValue(collection.Prisma),
        System: types.BoolValue(collection.System),
    }

    return schema, diags
}
