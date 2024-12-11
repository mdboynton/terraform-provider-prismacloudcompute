package models

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

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
