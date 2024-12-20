package planmodifiers

import (
	"context"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

type setGraceDaysBySeverityNullIfUnsupported struct {
    PolicyType string
}

func SetGraceDaysBySeverityNullIfUnsupported(policyType string) planmodifier.Object {
    return setGraceDaysBySeverityNullIfUnsupported{
        PolicyType: policyType,
    } 
}

func (m setGraceDaysBySeverityNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setGraceDaysBySeverityNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setGraceDaysBySeverityNullIfUnsupported) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        m.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage)
    
    graceDaysBySeverityTypeMap := map[string]attr.Type{
        "low": types.Int32Type,
        "medium": types.Int32Type,
        "high": types.Int32Type,
        "critical": types.Int32Type,
    }

    if !isSupportedPolicy {
        resp.PlanValue = types.ObjectNull(graceDaysBySeverityTypeMap)
        return
    } else {
        if !req.ConfigValue.IsNull() {
            return
        }

        resp.PlanValue = types.ObjectNull(graceDaysBySeverityTypeMap)
    }
}
