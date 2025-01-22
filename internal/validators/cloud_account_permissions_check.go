package validators

import (
    "context"
	
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

type permissionsCheckNotEnforcedIfScanningWithHub struct {}

func PermissionsCheckNotEnforcedIfScanningWithHub() permissionsCheckNotEnforcedIfScanningWithHub {
    return permissionsCheckNotEnforcedIfScanningWithHub{}
}

func (v permissionsCheckNotEnforcedIfScanningWithHub) Description(ctx context.Context) string {
    return ""
}

func (v permissionsCheckNotEnforcedIfScanningWithHub) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v permissionsCheckNotEnforcedIfScanningWithHub) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
    util.DLog(ctx, "Executing PermissionsCheckNotEnforcedIfScanningWithHub")

    var (
        hubAccountId basetypes.StringValue 
        enforcePermissionsCheck basetypes.BoolValue
    )

    resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("agentless_scanning").AtName("hub_account_id"), &hubAccountId)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("agentless_scanning").AtName("enforce_permissions_check"), &enforcePermissionsCheck)...)
    if resp.Diagnostics.HasError() {
        return
    }


    if (!hubAccountId.IsNull() && !hubAccountId.IsUnknown() && hubAccountId.ValueString() != "") && (!enforcePermissionsCheck.IsNull() && !enforcePermissionsCheck.IsUnknown() && enforcePermissionsCheck.ValueBool()) {
        resp.Diagnostics.AddError(
			"Invalid Resource Configuration",
            "Cloud Accounts cannot be configured with agentless scanning in hub mode (hub_credential_id set to non-empty string) and enforce_permissions_check set to true",
        )
    }

    return 
}
