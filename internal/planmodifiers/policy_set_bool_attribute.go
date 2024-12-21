package planmodifiers

import (
	"context"
    "fmt"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

func SetBoolAttributeByPolicyType(attributeName string, policyType string) planmodifier.Bool {
    return setBoolAttributeByPolicyType{
        AttributeName: attributeName,
        PolicyType: policyType,
    } 
}

type setBoolAttributeByPolicyType struct {
    AttributeName string
    PolicyType string
}

func (m setBoolAttributeByPolicyType) Description(_ context.Context) string {
    return ""
}

func (m setBoolAttributeByPolicyType) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setBoolAttributeByPolicyType) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    isSupported, ok := policyAPI.IsAttributeSupported[m.AttributeName][m.PolicyType]
    if !ok {
        tflog.Warn(ctx, fmt.Sprintf("Unable to determine if policy type \"%s\" supports attribute %s", m.PolicyType, m.AttributeName))
        return
    }

    if !isSupported {
        resp.PlanValue = basetypes.NewBoolNull()
    } else {
        if req.ConfigValue.IsNull() {
            resp.PlanValue = basetypes.NewBoolNull()
        } else {
            resp.PlanValue = req.ConfigValue
        }
    }

    return
}
