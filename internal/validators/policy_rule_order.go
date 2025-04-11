package validators

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func PolicyRuleOrderIsPositiveNonZero(policyType string) policyRuleOrderIsPositiveNonZero {
	return policyRuleOrderIsPositiveNonZero{
		PolicyType: policyType,
	}
}

type policyRuleOrderIsPositiveNonZero struct {
	PolicyType string
}

func (v policyRuleOrderIsPositiveNonZero) Description(ctx context.Context) string {
	return ""
}

func (v policyRuleOrderIsPositiveNonZero) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (v policyRuleOrderIsPositiveNonZero) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	util.LogDebug(ctx, "Executing PolicyRuleOrderIsPositiveNonZero")

	var (
		name  basetypes.StringValue
		order basetypes.Int32Value
	)

	for idx := range len(req.ConfigValue.Elements()) {
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("name"), &name)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("order"), &order)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !order.IsNull() && int(order.ValueInt32()) < 1 {
			resp.Diagnostics.AddError(
				"Invalid Resource Configuration",
				fmt.Sprintf("%s policy rule %s is configured with a negative or zero order value: %d", v.PolicyType, name, int(order.ValueInt32())),
			)
		}
	}

	util.LogDebug(ctx, "Finishing PolicyRuleOrderIsPositiveNonZero execution")

	return
}
