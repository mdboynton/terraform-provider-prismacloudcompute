package validators

import (
    "context"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func UserPermissionsNotConfiguredWithAdminRoles() userPermissionsNotConfiguredWithAdminRoles {
    return userPermissionsNotConfiguredWithAdminRoles{}
}

type userPermissionsNotConfiguredWithAdminRoles struct {}

func (v userPermissionsNotConfiguredWithAdminRoles) Description(ctx context.Context) string {
    return ""
}

func (v userPermissionsNotConfiguredWithAdminRoles) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v userPermissionsNotConfiguredWithAdminRoles) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
    if (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
        return
    }

    var role basetypes.StringValue
    resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("role"), &role)...)
    if resp.Diagnostics.HasError() {
        return
    }

    if (role.ValueString() == "admin" || role.ValueString() == "operator") {
        resp.Diagnostics.AddError(
            "Invalid Resource Configuration",
            "Permissions attribute cannot be used in conjunction with \"admin\" or \"operator\" roles.",
        )
    }

    return
}
