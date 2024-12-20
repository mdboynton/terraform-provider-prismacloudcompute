package planmodifiers

import (
	"context"
	"fmt"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func GenerateConditionFromEffect(policyType string, complianceVulnerabilities *[]systemAPI.Vulnerability) planmodifier.List{
    return generateConditionFromEffect{
        PolicyType: policyType,
        ComplianceVulnerabilities: complianceVulnerabilities,
    } 
}

type generateConditionFromEffect struct {
    PolicyType string
    ComplianceVulnerabilities *[]systemAPI.Vulnerability
}

func (m generateConditionFromEffect) Description(_ context.Context) string {
    return ""
}

func (m generateConditionFromEffect) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m generateConditionFromEffect) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    //util.DLog(ctx, "Executing planmodifier.GenerateConditionFromEffect")
    util.DLog(ctx, "****** Executing planmodifier.GenerateConditionFromEffect")
    
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeComplianceContainer || 
        m.PolicyType == policyAPI.PolicyTypeComplianceCiImage || 
        m.PolicyType == policyAPI.PolicyTypeComplianceHost || 
        m.PolicyType == policyAPI.PolicyTypeComplianceVmImage || 
        m.PolicyType == policyAPI.PolicyTypeComplianceFunction || 
        m.PolicyType == policyAPI.PolicyTypeComplianceCiFunction)

    if !isSupportedPolicy {
        util.DLog(ctx, "Unsupported policy: " + m.PolicyType)
        return
    }
   
    var (
        rules basetypes.ListValue
        nameValue basetypes.StringValue
        name string
        effectValue basetypes.StringValue
        effect string
        conditionValue basetypes.ObjectValue
        //condition models.PolicyRuleConditionResourceModel
        block bool
        isDefault bool
        vulnerabilities []models.PolicyRuleConditionVulnerabilityResourceModel
        complianceVulnerabilities []systemAPI.Vulnerability = *m.ComplianceVulnerabilities
    )


    resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rules"), &rules)...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    //vulnerabilities = []models.PolicyRuleConditionVulnerabilityResourceModel{}

    for idx := range rules.Elements() {
        // TODO: put rule order generation logic here

        // Get name value
        resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("name"), &nameValue)...)
        if resp.Diagnostics.HasError() {
            return
        }
        name = nameValue.ValueString()

        // Get effect value
        resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("effect"), &effectValue)...)
        if resp.Diagnostics.HasError() {
            return
        }

        // If effect is not configured, use the default value of "alert"
        // TODO: before we reach this point need to validate that there's no cases of condition being null and effect being "alert, block"
        if effectValue.IsNull() || effectValue.IsUnknown() {
            util.DLog(ctx, "assigning default value to effect")
            effect = "alert"
            resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("effect"), types.StringValue(effect))...)
            isDefault = true
        } else {
            effect = effectValue.ValueString()
            isDefault = false
        }

        // Get condition value
        resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("condition"), &conditionValue)...)
        if resp.Diagnostics.HasError() {
            return
        }

        if conditionValue.IsNull() {
            conditionTypeMap := map[string]attr.Type{
                "vulnerabilities": types.ListType{
                    ElemType: types.ObjectType{
                        AttrTypes: map[string]attr.Type{
                            "id":    types.Int32Type,
                            "block": types.BoolType,
                        },
                    },
                },
            }

            conditionValueMap := map[string]attr.Value{
                "vulnerabilities": types.ListValueMust(
                    types.ObjectType{
                        AttrTypes: map[string]attr.Type{
                            "id":    types.Int32Type,
                            "block": types.BoolType,
                        },
                    },
                    []attr.Value{},
                ),
            }

            conditionValue = types.ObjectValueMust(conditionTypeMap, conditionValueMap)

            //resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("condition"), conditionValue)...)
        }

        vulnerabilities = []models.PolicyRuleConditionVulnerabilityResourceModel{}

        if m.PolicyType == "serverlessCompliance" || m.PolicyType == "ciServerlessCompliance" {
            var conditionValue basetypes.ObjectValue
            resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("condition"), &conditionValue)...)
            if effect != "ignore" {
                //if effect == "default" {
                if effect == "alert" && isDefault {
                    complianceVulnerabilities = systemAPI.GetHighOrCriticalVulnerabilities(complianceVulnerabilities)
                }
                // If this block executes and effect is default, effect must be set to "alert"
                for _, vuln := range complianceVulnerabilities {
                    conditionVulnerability := models.PolicyRuleConditionVulnerabilityResourceModel{
                        Id:    types.Int32Value(int32(vuln.Id)),
                        Block: types.BoolValue(false),
                    }
                    vulnerabilities = append(vulnerabilities, conditionVulnerability)
                }
            }
        } else {
            // If the effect is "alert, block", create condition vulnerabilities object from plan
            if effect == "alert, block" {
                // TODO: add validator that ensures a condition is configured if the effect is set to "alert, block"
                var configVulnerabilities []models.PolicyRuleConditionVulnerabilityResourceModel
                resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("condition").AtName("vulnerabilities"), &configVulnerabilities)...)
                if resp.Diagnostics.HasError() {
                    return 
                }
                vulnerabilities = configVulnerabilities
                // Otherwise, if the rule effect is not "ignore", create condition vulnerabilities using Prisma Cloud vulnerability data
            } else if effect != "ignore" {
                if effect == "alert" && isDefault {
                    complianceVulnerabilities = systemAPI.GetHighOrCriticalVulnerabilities(complianceVulnerabilities)
                }

                isBlockEffect := (effect == "block")

                util.DLog(ctx, "looping over vulns")
                util.DLogf(ctx, len(complianceVulnerabilities))
                for _, vuln := range complianceVulnerabilities {
                    switch m.PolicyType {
                        case policyAPI.PolicyTypeComplianceHost:
                            block = (isBlockEffect && !(vuln.Type == "windows"))
                        case policyAPI.PolicyTypeComplianceContainer:
                            util.DLog(ctx, "in container case")
                            block = (isBlockEffect && !(vuln.Type == "istio" || vuln.Id == 58 || vuln.Id == 596 || vuln.Id == 598))
                        case policyAPI.PolicyTypeComplianceCiImage:
                            block = isBlockEffect
                        case policyAPI.PolicyTypeComplianceVmImage:
                            block = (isBlockEffect && !(vuln.Type == "istio" || vuln.Id == 58 || vuln.Id == 596 || vuln.Id == 598))
                        default:
                            resp.Diagnostics.AddError(
                                "Policy Rule Condition Conversion Error",
                                fmt.Sprintf("policy rule \"%s\" configured with unknown policy type \"%s\"", name, m.PolicyType),
                            )
                            return

                    }

                    conditionVulnerability := models.PolicyRuleConditionVulnerabilityResourceModel{
                        Id:    types.Int32Value(int32(vuln.Id)),
                        Block: types.BoolValue(block),
                    }
                    vulnerabilities = append(vulnerabilities, conditionVulnerability)
                }
            }
        }

        //generatedCondition.Vulnerabilities = conditionVulnerabilities
        //resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("condition").AtName("vulnerabilities"), vulnerabilities)...)
        resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(idx).AtName("condition"), conditionValue)...)
        if resp.Diagnostics.HasError() {
            return
        }

        resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("rules"), &resp.PlanValue)...)
    }

    //util.DLog(ctx, "Finishing planmodifier.GenerateConditionFromEffect execution")
    util.DLog(ctx, "****** Finishing planmodifier.GenerateConditionFromEffect execution")
}
