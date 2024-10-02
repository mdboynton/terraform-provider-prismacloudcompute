package validators

import (
    "context"
    "fmt"
	
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// TODO: figure out a way to import this from resource_policy_compliance_host
type HostCompliancePolicyRuleResourceModel struct {
    Name types.String `tfsdk:"name"`
    Order types.Int32 `tfsdk:"order"`
    Collections types.Set `tfsdk:"collections"`
    ReportAllPassedAndFailedChecks types.Bool `tfsdk:"report_passed_and_failed_checks"`
    BlockMessage types.String `tfsdk:"block_message"`
    Condition types.Object `tfsdk:"condition"`
    Disabled types.Bool `tfsdk:"disabled"`
    Effect types.String `tfsdk:"effect"`
    Modified types.String `tfsdk:"modified"`
    Notes types.String `tfsdk:"notes"`
    Owner types.String `tfsdk:"owner"`
    Verbose types.Bool `tfsdk:"verbose"`
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
    rulesModel := []HostCompliancePolicyRuleResourceModel{}
    diags := req.ConfigValue.ElementsAs(ctx, &rulesModel, false)
    if diags.HasError() {
        resp.Diagnostics.AddError(
            "Value Conversion Error",
            fmt.Sprintf("Error occured when converting %s policy resource's rules to Terraform type.", v.PolicyType),
        )
        return
    }

    unique := make(map[string]bool, len(req.ConfigValue.Elements()))
    
    for _, rule := range rulesModel {
        name := rule.Name.ValueString()
        if !unique[name] {
            unique[name] = true
        } else {
            resp.Diagnostics.AddAttributeError(
                req.Path,
                "Duplicate Rule Name",
                fmt.Sprintf("Found duplicate value for rule name \"%s\" in %s policy. All policy rules must have a unique name value.", name, v.PolicyType),
            )
            return
        }
    }

    return
}

func PolicyRuleNameIsUnique(policyType string) policyRuleNameIsUniqueValidator {
    return policyRuleNameIsUniqueValidator{
        PolicyType: policyType,
    }
}
