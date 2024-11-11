package validators

import (
    "context"
    "fmt"
	
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	//"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type policyRuleGraceDaysConfigIsValid struct {
    PolicyType string
}

func (v policyRuleGraceDaysConfigIsValid) Description(ctx context.Context) string {
    return ""
}

func (v policyRuleGraceDaysConfigIsValid) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v policyRuleGraceDaysConfigIsValid) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
    util.DLog(ctx, "Executing PolicyRuleGraceDaysConfigIsValid")

    rules := []policy.VulnerabilityPolicyRuleResourceModel{}
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &rules, false)...)
    if resp.Diagnostics.HasError() {
        util.DLog(ctx, "error in PolicyRuleGraceDaysConfigIsValid")
        return
    }

    for _, rule := range rules {
        if !rule.GraceDays.IsNull() && rule.GraceDaysPolicy != nil {
            resp.Diagnostics.AddError(
                "Invalid Resource Configuration",
                fmt.Sprintf("%s policy rule \"%s\" contains configurations for both grace_days_all_severities and grace_days_by_severity. Only one grace days policy may be configured for each policy rule.", v.PolicyType, rule.Name.ValueString()),
            )
        }
    }

    util.DLog(ctx, "Finishing PolicyRuleGraceDaysConfigIsValid execution")

    return
}

func PolicyRuleGraceDaysConfigIsValid(policyType string) policyRuleGraceDaysConfigIsValid {
    return policyRuleGraceDaysConfigIsValid{
        PolicyType: policyType,
    }
}
