package validators

import (
    "context"
    "fmt"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ConditionIsSupported(policyType string) conditionIsSupported {
    return conditionIsSupported{
        PolicyType: policyType,
    }
}

type conditionIsSupported struct {
    PolicyType string
}

func (v conditionIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v conditionIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v conditionIsSupported) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeComplianceContainer || 
        v.PolicyType == policyAPI.PolicyTypeComplianceCiImage || 
        v.PolicyType == policyAPI.PolicyTypeComplianceHost || 
        v.PolicyType == policyAPI.PolicyTypeComplianceVmImage || 
        v.PolicyType == policyAPI.PolicyTypeComplianceFunction || 
        v.PolicyType == policyAPI.PolicyTypeComplianceCiFunction)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := fmt.Sprintf("Policy type does not support condition argument")
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
