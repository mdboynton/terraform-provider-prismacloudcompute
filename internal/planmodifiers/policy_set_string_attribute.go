package planmodifiers

import (
	"context"
    "fmt"

	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

func SetStringAttributeByPolicyType(attributeName string, policyType string) planmodifier.String {
    return setStringAttributeByPolicyType{
        AttributeName: attributeName,
        PolicyType: policyType,
    } 
}

type setStringAttributeByPolicyType struct {
    AttributeName string
    PolicyType string
}

func (m setStringAttributeByPolicyType) Description(_ context.Context) string {
    return ""
}

func (m setStringAttributeByPolicyType) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setStringAttributeByPolicyType) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    isSupported, ok := policyAPI.IsAttributeSupported[m.AttributeName][m.PolicyType]
    if !ok {
        tflog.Warn(ctx, fmt.Sprintf("Unable to determine if policy type \"%s\" supports attribute %s", m.PolicyType, m.AttributeName))
        return
    }

    if !isSupported {
        resp.PlanValue = basetypes.NewStringNull()
    } else {
        if req.ConfigValue.IsNull() {
            resp.PlanValue = basetypes.NewStringNull()
        } else {
            resp.PlanValue = req.ConfigValue
        }
    }

    return
}
