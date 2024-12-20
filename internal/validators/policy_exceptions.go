package validators

import (
    "context"
    "fmt"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExceptionsIsSupported(policyType string, attributeName string) exceptionsIsSupported {
    return exceptionsIsSupported{
        PolicyType: policyType,
        AttributeName: attributeName,
    }
}

type exceptionsIsSupported struct {
    PolicyType string
    AttributeName string
}

func (v exceptionsIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v exceptionsIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v exceptionsIsSupported) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityHost ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityVmImage ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityFunction ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiFunction)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := fmt.Sprintf("Policy type does not support %s argument", v.AttributeName)
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
