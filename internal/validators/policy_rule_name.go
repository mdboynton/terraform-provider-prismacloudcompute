package validators

import (
    "context"
    "fmt"
    "slices"
	
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func PolicyRuleNameIsUnique(policyType string) policyRuleNameIsUniqueValidator {
    return policyRuleNameIsUniqueValidator{
        PolicyType: policyType,
    }
}

type policyRuleNameIsUniqueValidator struct {
    PolicyType string
}

func (v policyRuleNameIsUniqueValidator) Description(ctx context.Context) string {
    return ""
}

func (v policyRuleNameIsUniqueValidator) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v policyRuleNameIsUniqueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
    util.DLog(ctx, "Executing PolicyRuleNameIsUnique")

    var name basetypes.StringValue

    names := []string{}
    for idx := range len(req.ConfigValue.Elements()) {
        resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("name"), &name)...)
        if resp.Diagnostics.HasError() {
            return
        }

        if slices.Contains(names, name.ValueString()) {
            resp.Diagnostics.AddAttributeError(
                req.Path,
                "Duplicate Rule Name",
                fmt.Sprintf("Found duplicate value for rule name %s in %s policy. All rules in a policy must have a unique name.", name, v.PolicyType),
            )
            return
        } else {
            names = append(names, name.ValueString())
        }
    }

    util.DLog(ctx, "Finishing PolicyRuleNameIsUnique execution")

    return
}
