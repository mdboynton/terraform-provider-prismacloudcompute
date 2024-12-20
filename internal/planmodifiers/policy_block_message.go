package planmodifiers

import (
	"context"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetBlockMessageNullIfUnsupported(policyType string) planmodifier.String {
    return setBlockMessageNullIfUnsupported{
        PolicyType: policyType,
    } 
}

type setBlockMessageNullIfUnsupported struct {
    PolicyType string
}

func (m setBlockMessageNullIfUnsupported) Description(_ context.Context) string {
    return ""
}

func (m setBlockMessageNullIfUnsupported) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setBlockMessageNullIfUnsupported) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    isSupportedPolicy := (
        m.PolicyType == policyAPI.PolicyTypeComplianceContainer|| 
        m.PolicyType == policyAPI.PolicyTypeComplianceHost)
    
    if !isSupportedPolicy {
        resp.PlanValue = types.StringNull()
        return
    } else {
        if !req.ConfigValue.IsNull() {
            return
        }

        resp.PlanValue = types.StringNull()
    }
}
