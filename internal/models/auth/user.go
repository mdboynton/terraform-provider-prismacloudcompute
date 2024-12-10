package models

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

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

