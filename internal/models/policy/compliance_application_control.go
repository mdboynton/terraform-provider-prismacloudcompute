package models

import (
    "context"
    "time"

    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
)

type ApplicationControlPolicyResourceModel struct {
    Rules       *[]ApplicationControlPolicyRuleResourceModel    `tfsdk:"rules"`
}

type ApplicationControlPolicyRuleResourceModel struct {
    Id              types.Int32     `tfsdk:"id"`
    Applications    types.List      `tfsdk:"applications"`
    Description     types.String    `tfsdk:"description"`
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
    Severity        types.String    `tfsdk:"severity"`
}

var (
    applicationTypeMap = map[string]attr.Type{
        "name": types.StringType,
        "allowed_versions": types.ListType{
            ElemType: types.ListType{
                ElemType: types.StringType,
            },
        },
    }
    applicationObjectType = types.ObjectType{
        AttrTypes: applicationTypeMap,
    }
)

func (m ApplicationControlPolicyResourceModel) GetRuleIDs(ctx context.Context) []string {
    ruleIDs := []string{}
    
    if (m.Rules == nil || len(*m.Rules) == 0) {
        return ruleIDs 
    }

    for _, rule := range *m.Rules {
        ruleIDs = append(ruleIDs, rule.Name.ValueString())
    }

    return ruleIDs
}

func (m ApplicationControlPolicyResourceModel) ToTerraform(ctx context.Context) ([]policyAPI.ApplicationControlPolicyRule, diag.Diagnostics) {
    var diags diag.Diagnostics
    
    tfPolicy := []policyAPI.ApplicationControlPolicyRule{}

    if (m.Rules == nil || len(*m.Rules) == 0) {
        return tfPolicy, diags
    }

    for _, schemaRule := range *m.Rules {
        tfApplications := []policyAPI.ApplicationControlPolicyRuleApplication{}
        diags = schemaRule.Applications.ElementsAs(ctx, &tfApplications, false)
        if diags.HasError() {
            return tfPolicy, diags
        }

        tfRule := policyAPI.ApplicationControlPolicyRule{
            Id: int(schemaRule.Id.ValueInt32()),
            Name: schemaRule.Name.ValueString(),
            Description: schemaRule.Description.ValueString(),
            Applications: tfApplications,
            Modified: time.Now().Format("2006-01-02T15:04:05.000Z"),
            PreviousName: schemaRule.PreviousName.ValueString(),
            Severity: schemaRule.Severity.ValueString(),
        }

        tfPolicy = append(tfPolicy, tfRule)
    }

    return tfPolicy, diags
}

func (m *ApplicationControlPolicyResourceModel) FromTerraform(ctx context.Context, tfRules []policyAPI.ApplicationControlPolicyRule) diag.Diagnostics {
    var diags diag.Diagnostics

    rules := []ApplicationControlPolicyRuleResourceModel{}

    for _, tfRule := range tfRules {
        applicationValues := []attr.Value{}
        for _, tfApplication := range tfRule.Applications {
            versionConstraints := []attr.Value{}
            for _, versionContraint := range tfApplication.AllowedVersions {
                versionConstraint, diags := types.ListValueFrom(ctx, types.StringType, versionContraint)
                if diags.HasError() {
                    return diags
                }

                versionConstraints = append(versionConstraints, versionConstraint)
            }
            
            allowedVersions, diags := types.ListValueFrom(ctx, types.ListType{ ElemType: types.StringType }, versionConstraints)
            if diags.HasError() {
                return diags
            }

            application := types.ObjectValueMust(
                applicationTypeMap,
                map[string]attr.Value{
                    "name": types.StringValue(tfApplication.Name),
                    "allowed_versions": allowedVersions,
                },
            )

            applicationValues = append(applicationValues, application)
        }

        applications, diags := types.ListValueFrom(ctx, applicationObjectType, applicationValues)
        if diags.HasError() {
            return diags
        }
        
        rules = append(rules, ApplicationControlPolicyRuleResourceModel{
            Id: types.Int32Value(int32(tfRule.Id)),
            Name: types.StringValue(tfRule.Name),
            Description: types.StringValue(tfRule.Description),
            Applications: applications,
            Modified: types.StringValue(tfRule.Modified),
            Owner: types.StringValue(tfRule.Owner),
            PreviousName: types.StringValue(tfRule.PreviousName),
            Severity: types.StringValue(tfRule.Severity),
        })
    }

    m.Rules = &rules

    return diags
}

func (m ApplicationControlPolicyRuleResourceModel) ApplicationsTypeMap() map[string]attr.Type {
    return map[string]attr.Type{
        "name": types.StringType,
        "allowed_versions": types.ListType{
            ElemType: types.ListType{
                ElemType: types.StringType,
            },
        },
    }
}
