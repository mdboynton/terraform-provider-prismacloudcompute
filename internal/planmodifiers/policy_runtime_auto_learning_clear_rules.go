package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ClearPolicyRulesIfAutoLearningDisabled() planmodifier.Bool {
	return clearPolicyRulesIfAutoLearningDisabled{}
}

type clearPolicyRulesIfAutoLearningDisabled struct{}

func (m clearPolicyRulesIfAutoLearningDisabled) Description(_ context.Context) string {
	return ""
}

func (m clearPolicyRulesIfAutoLearningDisabled) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m clearPolicyRulesIfAutoLearningDisabled) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	var stateValue basetypes.BoolValue
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, req.Path, &stateValue)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Exit if state value is null/unknown or is already set to false
	if stateValue.IsNull() || stateValue.IsUnknown() || stateValue.ValueBool() == false {
		return
	}

	var planValue basetypes.BoolValue
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, req.Path, &planValue)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planValue.ValueBool() == false {
		var planRules basetypes.ListValue
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, req.Path.ParentPath().AtName("rules"), &planRules)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, req.Path.ParentPath().AtName("rules"), types.ListValueMust(planRules.ElementType(ctx), []attr.Value{}))...)
		// TODO: need to re-write this to be a modifier for the whole resource object
	}

	return
}
