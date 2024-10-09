package system

import (
	"context"
    "fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &CollectionResource{}
var _ resource.ResourceWithImportState = &CollectionResource{}

func NewCollectionResource() resource.Resource {
    return &CollectionResource{}
}

type CollectionResource struct {
    client *api.PrismaCloudComputeAPIClient
}

type CollectionResourceModel struct {
    AccountIDs types.Set `tfsdk:"account_ids"`
    AppIDs types.Set `tfsdk:"app_ids"`
    Clusters types.Set `tfsdk:"clusters"`
    Color types.String `tfsdk:"color"`
    Containers types.Set `tfsdk:"containers"`
    Description types.String `tfsdk:"description"`
    Functions types.Set `tfsdk:"functions"`
    Hosts types.Set `tfsdk:"hosts"`
    Images types.Set `tfsdk:"images"`
    Labels types.Set `tfsdk:"labels"`
    Modified types.String `tfsdk:"modified"`
    Name types.String `tfsdk:"name"`
    Namespaces types.Set `tfsdk:"namespaces"`
    Owner types.String `tfsdk:"owner"`
    Prisma types.Bool `tfsdk:"prisma"`
    System types.Bool `tfsdk:"system"`
}

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
    util.DLog(ctx, "retrieving plan and serializing into CollectionResourceModel")
    // Retrieve values from plan
    var plan CollectionResourceModel
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
    util.DLog(ctx, fmt.Sprintf("creating collection resource with payload:\n\n %+v", collection))
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
    var state CollectionResourceModel 
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
    var state CollectionResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan CollectionResourceModel
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
	var state CollectionResourceModel
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

// TODO: Define ImportState to work properly with this resource
func (r *CollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func CollectionObjectType() types.ObjectType {
    return types.ObjectType{
        AttrTypes: CollectionObjectAttrTypeMap(),
    }
}

func CollectionObjectAttrTypeMap() map[string]attr.Type {
    return map[string]attr.Type{
        "account_ids":  types.SetType{ElemType: types.StringType},
        "app_ids":  types.SetType{ElemType: types.StringType},
        "clusters":  types.SetType{ElemType: types.StringType},
        "color":        types.StringType,
        "containers":  types.SetType{ElemType: types.StringType},
        "description":  types.StringType,
        "functions":  types.SetType{ElemType: types.StringType},
        "hosts":  types.SetType{ElemType: types.StringType},
        "images":  types.SetType{ElemType: types.StringType},
        "labels":  types.SetType{ElemType: types.StringType},
        "modified": types.StringType,
        "name":         types.StringType,
        "namespaces":  types.SetType{ElemType: types.StringType},
        "owner":        types.StringType,
        "prisma":       types.BoolType,
        "system":       types.BoolType,
    }
}

func CollectionObjectDefaultAttrValueMap() map[string]attr.Value {
    return map[string]attr.Value{
        "account_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "app_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "clusters": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "color": types.StringValue("#3FA2F7"),
        "containers": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "description": types.StringValue("System - all resources collection"),
        "functions": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "hosts": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "images": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "labels": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "modified": types.StringValue(""),
        "name": types.StringValue("All"),
        "namespaces": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "owner": types.StringValue("system"),
        "prisma": types.BoolValue(false),
        "system": types.BoolValue(true),
    }
}

func schemaToCollection(ctx context.Context, plan *CollectionResourceModel) (collectionAPI.Collection, diag.Diagnostics) {
    var diags diag.Diagnostics

    collection := collectionAPI.Collection{
        Color: plan.Color.ValueString(),
        Description: plan.Description.ValueString(),
        Name: plan.Name.ValueString(),
        //Modified: plan.Modified.ValueString(),
        Prisma: plan.Prisma.ValueBool(),
        System: plan.System.ValueBool(),
    }

    accountIds := make([]string, 0, len(plan.AccountIDs.Elements()))
    diags = plan.AccountIDs.ElementsAs(ctx, &accountIds, false)
    if diags.HasError() {
        return collection, diags
    }
    collection.AccountIDs = accountIds


    appIds := make([]string, 0, len(plan.AppIDs.Elements()))
    diags = plan.AppIDs.ElementsAs(ctx, &appIds, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.AppIDs = appIds 

    clusters := make([]string, 0, len(plan.Clusters.Elements()))
    diags = plan.Clusters.ElementsAs(ctx, &clusters, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.Clusters = clusters 

    containers := make([]string, 0, len(plan.Containers.Elements()))
    diags = plan.Containers.ElementsAs(ctx, &containers, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.Containers = containers 

    functions := make([]string, 0, len(plan.Functions.Elements()))
    diags = plan.Functions.ElementsAs(ctx, &functions, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.Functions = functions 

    hosts := make([]string, 0, len(plan.Hosts.Elements()))
    diags = plan.Hosts.ElementsAs(ctx, &hosts, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.Hosts = hosts 

    images := make([]string, 0, len(plan.Images.Elements()))
    diags = plan.Images.ElementsAs(ctx, &images, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.Images = images 

    labels := make([]string, 0, len(plan.Labels.Elements()))
    diags = plan.Labels.ElementsAs(ctx, &labels, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.Labels = labels 

    namespaces := make([]string, 0, len(plan.Namespaces.Elements()))
    diags = plan.Namespaces.ElementsAs(ctx, &namespaces, false)
    if diags.HasError() {   
        return collection, diags
    }
    collection.Namespaces = namespaces 

	return collection, diags 
}

func collectionToSchema(ctx context.Context, collection collectionAPI.Collection) (CollectionResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    schema := CollectionResourceModel{
        Color: types.StringValue(collection.Color),
        Description: types.StringValue(collection.Description),
        //Modified: types.StringValue(collection.Modified),
        Name: types.StringValue(collection.Name),
        Prisma: types.BoolValue(collection.Prisma),
        System: types.BoolValue(collection.System),
    }

    //if collection.Modified != nil {
    //    schema.Modified = collection.Modified
    //}

    if collection.AccountIDs != nil {
        accountIds, diags := types.SetValueFrom(ctx, types.StringType, collection.AccountIDs)
        if diags.HasError() {
            return schema, diags
        }

        schema.AccountIDs = accountIds
    }

    if collection.AppIDs != nil {
        appIds, diags := types.SetValueFrom(ctx, types.StringType, collection.AppIDs)
        if diags.HasError() {
            return schema, diags
        }

        schema.AppIDs = appIds
    }

    if collection.Clusters != nil {
        clusters, diags := types.SetValueFrom(ctx, types.StringType, collection.Clusters)
        if diags.HasError() {
            return schema, diags
        }

        schema.Clusters = clusters
    }

    if collection.Containers != nil {
        containers, diags := types.SetValueFrom(ctx, types.StringType, collection.Containers)
        if diags.HasError() {
            return schema, diags
        }

        schema.Containers = containers
    }

    if collection.Functions != nil {
        functions, diags := types.SetValueFrom(ctx, types.StringType, collection.Functions)
        if diags.HasError() {
            return schema, diags
        }

        schema.Functions = functions
    }

    if collection.Hosts != nil {
        hosts, diags := types.SetValueFrom(ctx, types.StringType, collection.Hosts)
        if diags.HasError() {
            return schema, diags
        }

        schema.Hosts = hosts
    }

    if collection.Images != nil {
        images, diags := types.SetValueFrom(ctx, types.StringType, collection.Images)
        if diags.HasError() {
            return schema, diags
        }

        schema.Images = images
    }

    if collection.Labels != nil {
        labels, diags := types.SetValueFrom(ctx, types.StringType, collection.Labels)
        if diags.HasError() {
            return schema, diags
        }

        schema.Labels = labels
    }
    
    if collection.Namespaces != nil {
        namespaces, diags := types.SetValueFrom(ctx, types.StringType, collection.Namespaces)
        if diags.HasError() {
            return schema, diags
        }

        schema.Namespaces = namespaces
    }

    return schema, diags
}
