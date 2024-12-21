package planmodifiers

import (
	"context"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetExcludeBaseImageVulnsNullIfUnsupported(policyType string) planmodifier.Bool {
    return setExcludeBaseImageVulnsNullIfUnsupported{
        PolicyType: policyType,
    } 
}

type setExcludeBaseImageVulnsNullIfUnsupported struct {
    PolicyType string
}

func (m setExcludeBaseImageVulnsNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setExcludeBaseImageVulnsNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setExcludeBaseImageVulnsNullIfUnsupported) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage)
    
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
