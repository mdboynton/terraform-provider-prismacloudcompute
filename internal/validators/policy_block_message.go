package validators

import (
    "context"
    "fmt"
	
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func BlockMessageIsSupported(policyType string) blockMessageIsSupported {
    return blockMessageIsSupported{
        PolicyType: policyType,
    }
}

type blockMessageIsSupported struct {
    PolicyType string
}

func (v blockMessageIsSupported) Description(ctx context.Context) string {
    return ""
}

func (v blockMessageIsSupported) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v blockMessageIsSupported) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
    isSupportedPolicy := (
        v.PolicyType == policyAPI.PolicyTypeComplianceContainer || 
        v.PolicyType == policyAPI.PolicyTypeComplianceHost)

    if !req.ConfigValue.IsNull() && !isSupportedPolicy {
        errorMessage := fmt.Sprintf("Policy type does not support block_message argument")
        resp.Diagnostics.AddAttributeError(req.Path, "Policy Rule Validation Error", errorMessage)
    }

    return
}
