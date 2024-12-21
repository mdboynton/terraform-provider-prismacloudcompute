package planmodifiers

import (
    "context"

    policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/attr"
)

func SetExceptionsNullIfUnsupported(policyType string) planmodifier.Set {
    return setExceptionsNullIfUnsupported{
        PolicyType: policyType,
    } 
}

type setExceptionsNullIfUnsupported struct {
    PolicyType string
}

func (m setExceptionsNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setExceptionsNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setExceptionsNullIfUnsupported) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityHost ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityVmImage ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityFunction ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiFunction)
   
    exceptionsTypeMap := map[string]attr.Type{
        "effect": types.StringType,
        "id": types.StringType,
        "description": types.StringType,
        "expiration": types.ObjectType{
            AttrTypes: map[string]attr.Type{
                "enabled": types.BoolType,
                "date": types.StringType,
            },
        },
    }

    exceptionsType := types.ObjectType{
        AttrTypes: exceptionsTypeMap,
    }

    if !isSupportedPolicy {
        resp.PlanValue = types.SetNull(exceptionsType)
        return
    } else {
        if !req.ConfigValue.IsNull() {
            return
        }

        resp.PlanValue = types.SetNull(exceptionsType)
    }
}
