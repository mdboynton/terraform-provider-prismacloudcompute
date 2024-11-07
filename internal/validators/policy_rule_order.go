package validators

import (
    "context"
    "fmt"
	
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	//"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

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
    util.DLog(ctx, "Executing PolicyRuleOrderIsPositiveNonZero")

    rules := []models.CompliancePolicyRuleResourceModel{}
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &rules, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    for _, rule := range rules {
        // Skip if null, as this will be set in the plan modifier
        if rule.Order.IsNull() {
            continue
        }

        name := rule.Name.ValueString()
        order := rule.Order.ValueInt32()

        if int(order) < 1 {
            resp.Diagnostics.AddError(
		    	"Invalid Resource Configuration",
		    	fmt.Sprintf("%s policy rule \"%s\" configured with a negative or non-zero order (%d)", v.PolicyType, name, order),
            )
        }
    }

    return
}

func PolicyRuleOrderIsPositiveNonZero(policyType string) policyRuleOrderIsPositiveNonZero {
    return policyRuleOrderIsPositiveNonZero{
        PolicyType: policyType,
    }
}




type policyRuleOrderIsPositiveNonZeroVuln struct {
    PolicyType string
}

func (v policyRuleOrderIsPositiveNonZeroVuln) Description(ctx context.Context) string {
    return ""
}

func (v policyRuleOrderIsPositiveNonZeroVuln) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v policyRuleOrderIsPositiveNonZeroVuln) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
    util.DLog(ctx, "Executing PolicyRuleOrderIsPositiveNonZeroVuln")

    rules := []models.VulnerabilityPolicyRuleResourceModel{}
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &rules, false)...)
    if resp.Diagnostics.HasError() {
        util.DLog(ctx, "error in PolicyRuleOrderIsPositiveNonZeroVuln")
        return
    }

    for _, rule := range rules {
        // Skip if null, as this will be set in the plan modifier
        if rule.Order.IsNull() {
            continue
        }

        name := rule.Name.ValueString()
        order := rule.Order.ValueInt32()

        if int(order) < 1 {
            resp.Diagnostics.AddError(
		    	"Invalid Resource Configuration",
		    	fmt.Sprintf("%s policy rule \"%s\" configured with a negative or non-zero order (%d)", v.PolicyType, name, order),
            )
        }
    }

    util.DLog(ctx, "Finishing PolicyRuleOrderIsPositiveNonZeroVuln execution")

    return
}

func PolicyRuleOrderIsPositiveNonZeroVuln(policyType string) policyRuleOrderIsPositiveNonZeroVuln {
    return policyRuleOrderIsPositiveNonZeroVuln{
        PolicyType: policyType,
    }
}
