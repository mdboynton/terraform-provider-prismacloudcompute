package validators

import (
    "context"
    "fmt"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ApplyOnlyWhenFixAvailableIsSupported(policyType string) applyOnlyWhenFixAvailableIsSupported {
    return applyOnlyWhenFixAvailableIsSupported{
        PolicyType: policyType,
    }
}

type applyOnlyWhenFixAvailableIsSupported struct {
    PolicyType string
}

func (v applyOnlyWhenFixAvailableIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v applyOnlyWhenFixAvailableIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v applyOnlyWhenFixAvailableIsSupported) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityHost ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityVmImage ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityFunction ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiFunction)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := fmt.Sprintf("Policy type does not support apply_only_when_fix_available argument")
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
