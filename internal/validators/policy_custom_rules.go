package validators

import (
    "fmt"
    "context"
    "slices"

    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func CustomRulesAreValid() customRulesAreValid {
    return customRulesAreValid{}
}

type customRulesAreValid struct {}

func (v customRulesAreValid) Description(ctx context.Context) string {
    return ""
}

func (v customRulesAreValid) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v customRulesAreValid) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
    if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
        return
    }

    // TODO: how to determine which rules are "network-outgoing" type?
    //  these rules cannot be paired with "prevent" effect
    // (probably need to do this in ModifyPlan)

    customRules := []models.RuntimeHostPolicyCustomRuleResourceModel{}
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &customRules, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    names := []string{}
    for idx, customRule := range customRules {
        name := customRule.Name.ValueString()
        effect := customRule.Effect.ValueString()
        logAs := customRule.LogAs.ValueString()
        
        if (effect == "allow" && logAs != "") {
            resp.Diagnostics.AddAttributeError(
                req.Path, 
                "Invalid Resource Configuration", 
                fmt.Sprintf("custom_rule[%d] contains a log_as value. Custom rules with an effect of \"allow\" cannot specify a log_as value.", idx), 
            )
            continue
        }

        if slices.Contains(names, name) {
            resp.Diagnostics.AddAttributeError(
                req.Path, 
                "Invalid Resource Configuration", 
                fmt.Sprintf("custom_rule[%d] is a duplicate of an existing custom rule attached to this policy rule.", idx), 
            )
            continue
        } else {
            names = append(names, name)
        }
    }


    return
}
