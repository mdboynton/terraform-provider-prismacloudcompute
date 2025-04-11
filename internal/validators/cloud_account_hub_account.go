package validators

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type hubAccountNotTrueIfHubCredentialIdSet struct{}

func HubAccountNotTrueIfHubCredentialIdSet() hubAccountNotTrueIfHubCredentialIdSet {
	return hubAccountNotTrueIfHubCredentialIdSet{}
}

func (v hubAccountNotTrueIfHubCredentialIdSet) Description(ctx context.Context) string {
	return ""
}

func (v hubAccountNotTrueIfHubCredentialIdSet) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (v hubAccountNotTrueIfHubCredentialIdSet) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	util.LogDebug(ctx, "Executing HubAccountNotTrueIfHubCredentialIdSet")

	var (
		isHubAccount basetypes.BoolValue
		hubAccountId basetypes.StringValue
	)

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("agentless_scanning").AtName("is_hub_account"), &isHubAccount)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("agentless_scanning").AtName("hub_account_id"), &hubAccountId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if isHubAccount.ValueBool() && !hubAccountId.IsNull() && !hubAccountId.IsUnknown() && hubAccountId.ValueString() != "" {
		resp.Diagnostics.AddError(
			"Invalid Resource Configuration",
			"Cloud Accounts cannot be configured with is_hub_account set to true and hub_credential_id set to non-empty string",
		)
	}

	return
}
