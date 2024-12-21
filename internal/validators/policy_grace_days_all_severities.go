package validators

import (
    "context"
    "fmt"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func GraceDaysAllSeveritiesIsSupported(policyType string) graceDaysAllSeveritiesIsSupported {
    return graceDaysAllSeveritiesIsSupported{
        PolicyType: policyType,
    }
}

type graceDaysAllSeveritiesIsSupported struct {
    PolicyType string
}

func (v graceDaysAllSeveritiesIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v graceDaysAllSeveritiesIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v graceDaysAllSeveritiesIsSupported) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityDeployedImage || 
        v.PolicyType == policyAPI.PolicyTypeVulnerabilityCiImage)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := fmt.Sprintf("Policy type does not support grace_days_all_severities argument")
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
