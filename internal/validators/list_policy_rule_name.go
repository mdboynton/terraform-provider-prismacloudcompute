package validators

import (
    "context"
    "fmt"
	
	//models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

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

    //rulesModel := []models.CompliancePolicyRuleResourceModel{}
    rulesModel := []models.PolicyRuleResourceModel{}
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

    util.DLog(ctx, "Finishing PolicyRuleNameIsUnique execution")

    return
}

func PolicyRuleNameIsUnique(policyType string) policyRuleNameIsUniqueValidator {
    return policyRuleNameIsUniqueValidator{
        PolicyType: policyType,
    }
}



type policyRuleNameIsUniqueValidatorVuln struct {
    PolicyType string
}

func (v policyRuleNameIsUniqueValidatorVuln) Description(ctx context.Context) string {
    return ""
}

func (v policyRuleNameIsUniqueValidatorVuln) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v policyRuleNameIsUniqueValidatorVuln) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
    util.DLog(ctx, "Executing PolicyRuleNameIsUniqueVuln")

    rulesModel := []models.VulnerabilityPolicyRuleResourceModel{}
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

func PolicyRuleNameIsUniqueVuln(policyType string) policyRuleNameIsUniqueValidatorVuln {
    return policyRuleNameIsUniqueValidatorVuln{
        PolicyType: policyType,
    }
}
