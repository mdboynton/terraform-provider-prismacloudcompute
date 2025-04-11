package validators

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type errorIfAllowListNonEmptyWithBlock struct{}

func ErrorIfAllowListNonEmptyWithBlock() errorIfAllowListNonEmptyWithBlock {
	return errorIfAllowListNonEmptyWithBlock{}
}

func (v errorIfAllowListNonEmptyWithBlock) Description(ctx context.Context) string {
	return ""
}

func (v errorIfAllowListNonEmptyWithBlock) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (v errorIfAllowListNonEmptyWithBlock) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	util.LogDebug(ctx, "Executing ErrorIfAllowListNonEmptyWithBlock")

	var (
		allowList                   basetypes.ListValue
		blockAllProcessesExceptMain basetypes.BoolValue
	)

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName("allowed_processes"), &allowList)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName("block_all_processes_except_main"), &blockAllProcessesExceptMain)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(allowList.Elements()) > 0 && blockAllProcessesExceptMain.ValueBool() {
		resp.Diagnostics.AddError(
			"Invalid Resource Configuration",
			fmt.Sprintf("Invalid configuration at %s \nallowed_processes cannot be non-empty if block_all_processes_except_main is enabled", req.Path.String()),
		)
	}

	return
}
