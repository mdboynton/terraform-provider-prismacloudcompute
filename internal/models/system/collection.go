package models

import (
    "context"

	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

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

type CollectionDataSourceModel struct {
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

func (d CollectionDataSourceModel) RefreshPropertyValues(ctx context.Context, collection *collectionAPI.Collection) (CollectionDataSourceModel, diag.Diagnostics) {
    accountIds, diags := types.SetValueFrom(ctx, types.StringType, collection.AccountIDs)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    appIds, diags := types.SetValueFrom(ctx, types.StringType, collection.AppIDs)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    clusters, diags := types.SetValueFrom(ctx, types.StringType, collection.Clusters)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    containers, diags := types.SetValueFrom(ctx, types.StringType, collection.Containers)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    functions, diags := types.SetValueFrom(ctx, types.StringType, collection.Functions)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    hosts, diags := types.SetValueFrom(ctx, types.StringType, collection.Hosts)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    images, diags := types.SetValueFrom(ctx, types.StringType, collection.Images)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    labels, diags := types.SetValueFrom(ctx, types.StringType, collection.Labels)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    namespaces, diags := types.SetValueFrom(ctx, types.StringType, collection.Namespaces)
    if diags.HasError() {
        return CollectionDataSourceModel{}, diags
    }

    return CollectionDataSourceModel{
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
    }, diags
}
