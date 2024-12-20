package planmodifiers

import (
	"context"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type setGraceDaysAllSeveritiesNullIfUnsupported struct {
    PolicyType string
}

func SetGraceDaysAllSeveritiesNullIfUnsupported(policyType string) planmodifier.Int32 {
    return setGraceDaysAllSeveritiesNullIfUnsupported{
        PolicyType: policyType,
    } 
}

func (m setGraceDaysAllSeveritiesNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setGraceDaysAllSeveritiesNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setGraceDaysAllSeveritiesNullIfUnsupported) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage)
    
    if !isSupportedPolicy {
        resp.PlanValue = types.Int32Null()
        return
    } else {
        if !req.ConfigValue.IsNull() {
            return
        }

        resp.PlanValue = types.Int32Null()
    }
}
