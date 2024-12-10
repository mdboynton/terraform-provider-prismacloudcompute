package models

import (
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

