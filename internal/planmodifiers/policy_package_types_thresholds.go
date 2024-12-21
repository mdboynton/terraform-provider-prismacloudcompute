package planmodifiers

import (
    "context"

    policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/attr"
)

func SetPackageTypesThresholdsNullIfUnsupported(policyType string) planmodifier.Set {
    return setPackageTypesThresholdsNullIfUnsupported{
        PolicyType: policyType,
    } 
}

type setPackageTypesThresholdsNullIfUnsupported struct {
    PolicyType string
}

func (m setPackageTypesThresholdsNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setPackageTypesThresholdsNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setPackageTypesThresholdsNullIfUnsupported) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityHost ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityVmImage ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityFunction ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiFunction)
   
    packageTypesThresholdsTypeMap := map[string]attr.Type{
        "type": types.StringType,
        "alert_threshold": types.StringType,
        "block_threshold": types.StringType,
    }

    packageTypesThresholdsType := types.ObjectType{
        AttrTypes: packageTypesThresholdsTypeMap,
    }

    if !isSupportedPolicy {
        resp.PlanValue = types.SetNull(packageTypesThresholdsType)
        return
    } else {
        if !req.ConfigValue.IsNull() {
            return
        }

        resp.PlanValue = types.SetNull(packageTypesThresholdsType)
    }
}
