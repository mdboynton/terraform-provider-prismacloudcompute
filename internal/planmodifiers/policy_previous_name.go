package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func UseStateIfRuleNameUnchanged() planmodifier.String {
	return useStateIfRuleNameUnchanged{}
}

type useStateIfRuleNameUnchanged struct{}

func (m useStateIfRuleNameUnchanged) Description(_ context.Context) string {
	return ""
}

func (m useStateIfRuleNameUnchanged) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m useStateIfRuleNameUnchanged) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var stateName basetypes.StringValue
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, req.Path.ParentPath().AtName("name"), &stateName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Exit if name is null or unknown in state
	if stateName.IsNull() || stateName.IsUnknown() {
		return
	}

	var planName basetypes.StringValue
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, req.Path.ParentPath().AtName("name"), &planName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planName.IsNull() && !planName.IsUnknown() && planName.ValueString() == stateName.ValueString() {
		var statePreviousName basetypes.StringValue
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, req.Path, &statePreviousName)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.PlanValue = statePreviousName
	}

	return
}
