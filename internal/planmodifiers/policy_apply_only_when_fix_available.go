package planmodifiers

import (
	"context"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetApplyOnlyWhenFixAvailableNullIfUnsupported(policyType string) planmodifier.Bool {
    return setApplyOnlyWhenFixAvailableNullIfUnsupported{
        PolicyType: policyType,
    } 
}

type setApplyOnlyWhenFixAvailableNullIfUnsupported struct {
    PolicyType string
}

func (m setApplyOnlyWhenFixAvailableNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setApplyOnlyWhenFixAvailableNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setApplyOnlyWhenFixAvailableNullIfUnsupported) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityHost ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityVmImage ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityFunction ||
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiFunction)
    
    if !isSupportedPolicy {
        resp.PlanValue = types.BoolNull()
        return
    } else {
        if !req.ConfigValue.IsNull() {
            return
        }

        resp.PlanValue = types.BoolValue(false)
    }
}
