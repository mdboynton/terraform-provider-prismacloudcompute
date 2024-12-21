package validators

import (
    "context"
    "fmt"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExcludeBaseImageVulnsIsSupported(policyType string) excludeBaseImageVulnsIsSupported {
    return excludeBaseImageVulnsIsSupported{
        PolicyType: policyType,
    }
}

type excludeBaseImageVulnsIsSupported struct {
    PolicyType string
}

func (v excludeBaseImageVulnsIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v excludeBaseImageVulnsIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v excludeBaseImageVulnsIsSupported) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := fmt.Sprintf("Policy type does not support exclude_base_image_vulns argument")
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
