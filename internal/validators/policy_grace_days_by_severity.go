package validators

import (
    "context"
    "fmt"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func GraceDaysBySeverityIsSupported(policyType string) graceDaysBySeverityIsSupported {
    return graceDaysBySeverityIsSupported{
        PolicyType: policyType,
    }
}

type graceDaysBySeverityIsSupported struct {
    PolicyType string
}

func (v graceDaysBySeverityIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v graceDaysBySeverityIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v graceDaysBySeverityIsSupported) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := fmt.Sprintf("Policy type does not support grace_days_by_severity argument")
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
