package validators

import (
    "context"
    "fmt"
	
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
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
    rules := []models.CompliancePolicyRuleResourceModel{}
    //resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &rules, false)...)
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &rules, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    for _, rule := range rules {
        // Skip if unknown, as this will be set in the plan modifier
        //if rule.Order.IsUnknown() {
        if rule.Order.IsNull() {
            continue
        }

        name := rule.Name.ValueString()
        order := rule.Order.ValueInt32()

        if int(order) < 1 {
            resp.Diagnostics.AddError(
		    	"Invalid Resource Configuration",
		    	fmt.Sprintf("policy rule \"%s\" configured with a negative or non-zero order (%d)", name, order),
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
