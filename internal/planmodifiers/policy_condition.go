package planmodifiers

import (
    "context"

    policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/attr"
)

func SetConditionNullIfUnsupported(policyType string) planmodifier.Object {
    return setConditionNullIfUnsupported{
        PolicyType: policyType,
    } 
}

type setConditionNullIfUnsupported struct {
    PolicyType string
}

func (m setConditionNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setConditionNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setConditionNullIfUnsupported) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeComplianceContainer || 
        m.PolicyType == policyAPI.PolicyTypeComplianceCiImage || 
        m.PolicyType == policyAPI.PolicyTypeComplianceHost || 
        m.PolicyType == policyAPI.PolicyTypeComplianceVmImage || 
        m.PolicyType == policyAPI.PolicyTypeComplianceFunction || 
        m.PolicyType == policyAPI.PolicyTypeComplianceCiFunction)
   
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

    if !isSupportedPolicy {
        resp.PlanValue = types.ObjectNull(conditionTypeMap)
        return
    } else {
        if !req.ConfigValue.IsNull() {
            return
        }

        resp.PlanValue = types.ObjectNull(conditionTypeMap)
    }
}
