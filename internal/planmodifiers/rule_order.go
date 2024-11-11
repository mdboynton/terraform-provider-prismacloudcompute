package planmodifiers

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

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
    util.DLog(ctx, "Executing UseIndexForUnknownOrder")

    var rules []policy.CompliancePolicyRuleResourceModel
    diags := req.PlanValue.ElementsAs(ctx, &rules, false)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // TODO: test this with empty rule set
    if rules == nil {
        return
    }

    for index, rule := range rules {
        if rule.Order.IsUnknown() {
            rule.Order = types.Int32Value(int32(index + 1))
            rules[index] = rule
        }
    }

    updatedPlan, diags := types.ListValueFrom(ctx, req.PlanValue.ElementType(ctx), rules)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.PlanValue = updatedPlan 

    return
}
