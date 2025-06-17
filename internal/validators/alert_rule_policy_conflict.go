package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type rulesNotConfiguredWithAllRulesEnabled struct{}

func RulesNotConfiguredWithAllRulesEnabled() rulesNotConfiguredWithAllRulesEnabled {
	return rulesNotConfiguredWithAllRulesEnabled{}
}

func (v rulesNotConfiguredWithAllRulesEnabled) Description(ctx context.Context) string {
	return ""
}

func (v rulesNotConfiguredWithAllRulesEnabled) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (v rulesNotConfiguredWithAllRulesEnabled) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	var (
		allRules basetypes.BoolValue
		rules    basetypes.ListValue
	)

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName("all_rules"), &allRules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName("rules"), &rules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allRulesEnabled := (!allRules.IsNull() && !allRules.IsUnknown() && allRules.ValueBool())
	rulesConfigured := (!rules.IsNull() && !rules.IsUnknown() && len(rules.Elements()) > 0)

	if allRulesEnabled && rulesConfigured {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Combination",
			"Attribute \"rules\" cannot be configured with a non-empty list "+
				"when \"all_rules\" is known and set to true",
		)
	}

	return
}
