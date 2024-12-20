package validators

import (
    "context"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func PackageTypesThresholdsIsSupported(policyType string) packageTypesThresholdsIsSupported {
    return packageTypesThresholdsIsSupported{
        PolicyType: policyType,
    }
}

type packageTypesThresholdsIsSupported struct {
    PolicyType string
}

func (v packageTypesThresholdsIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v packageTypesThresholdsIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v packageTypesThresholdsIsSupported) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityHost ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityVmImage ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityFunction ||
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiFunction)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := "Policy type does not support package_types_thresholds argument"
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
