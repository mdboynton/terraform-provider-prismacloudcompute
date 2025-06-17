package planmodifiers

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseIndexForUnknownOrder(policyType string) planmodifier.List {
	return useIndexForUnknownOrder{
		PolicyType: policyType,
	}
}

type useIndexForUnknownOrder struct {
	PolicyType string
}

func (m useIndexForUnknownOrder) Description(_ context.Context) string {
	return ""
}

func (m useIndexForUnknownOrder) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m useIndexForUnknownOrder) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	util.HCLogDebug(ctx, "Executing UseIndexForUnknownOrder")

	var rules basetypes.ListValue
	diags := req.Plan.GetAttribute(ctx, path.Root("rules"), &rules)
	if diags.HasError() {
		return
	}

	var order basetypes.Int32Value
	for index := range rules.Elements() {
		diags = req.Plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("order"), &order)
		if diags.HasError() {
			return
		}

		if order.IsUnknown() {
			diags = req.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("order"), types.Int32Value(int32(index+1)))
		}
	}

	diags = req.Plan.GetAttribute(ctx, path.Root("rules"), &resp.PlanValue)

	util.HCLogDebug(ctx, "Finishing UseIndexForUnknownOrder execution")

	return
}
