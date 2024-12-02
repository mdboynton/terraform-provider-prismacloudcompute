package planmodifiers

import (
	"context"

	//models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

func UseIndexForUnknownOrderVuln(policyType string) planmodifier.List {
    return useIndexForUnknownOrderVuln{
        PolicyType: policyType,
    } 
}

type useIndexForUnknownOrderVuln struct {
    PolicyType string
}

func (m useIndexForUnknownOrderVuln) Description(_ context.Context) string {
    return ""
}

func (m useIndexForUnknownOrderVuln) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useIndexForUnknownOrderVuln) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    util.DLog(ctx, "Executing UseIndexForUnknownOrderVuln")

    var rules []models.VulnerabilityPolicyRuleResourceModel
    diags := req.PlanValue.ElementsAs(ctx, &rules, false)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        util.DLog(ctx, "error in UseIndexForUnknownOrderVuln")
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
        util.DLog(ctx, "error in UseIndexForUnknownOrderVuln")
        return
    }

    resp.PlanValue = updatedPlan 

    return
}
