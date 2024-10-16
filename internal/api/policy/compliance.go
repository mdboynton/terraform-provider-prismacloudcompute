package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
)

const (
	HostComplianceEndpoint      = "api/v1/policies/compliance/host"
	ContainerComplianceEndpoint = "api/v1/policies/compliance/container"
	CiImageComplianceEndpoint   = "api/v1/policies/compliance/ci/images"
	VmImageComplianceEndpoint   = "api/v1/policies/compliance/vms"
    FunctionComplianceEndpoint  = "api/v1/policies/compliance/serverless"
)

func getEndpointAndPolicyName(policyType string) (string, string, error) {
    switch policyType{
        case "hostCompliance":
            return HostComplianceEndpoint, "host", nil
        case "containerCompliance":
            return ContainerComplianceEndpoint, "container", nil
        case "ciImagesCompliance":
            return CiImageComplianceEndpoint, "CI image", nil
        case "vmCompliance":
            return VmImageComplianceEndpoint, "VM image", nil
        case "serverlessCompliance":
            return FunctionComplianceEndpoint, "function", nil
        default:
            return "", "", fmt.Errorf("invalid policy type specified")
    }
}

// Create/Update/Delete compliance policy
func UpsertCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy CompliancePolicy) error {
    endpoint, policyName, err := getEndpointAndPolicyName(policy.PolicyType)
    if err != nil {
		return err
    }

    if err := c.Request(http.MethodPut, endpoint, nil, policy, nil); err != nil {
		return fmt.Errorf("error getting %s compliance policy: %s", policyName, err)
    }
    
    return nil
}

// Get compliance policy
func GetCompliancePolicy(c api.PrismaCloudComputeAPIClient, policyType string) (*CompliancePolicy, error) {
    var ans CompliancePolicy

    endpoint, policyName, err := getEndpointAndPolicyName(policyType)
    if err != nil {
		return &ans, err
    }
    
    if err := c.Request(http.MethodGet, endpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting %s compliance policy: %s", policyName, err)
    }
    
    return &ans, nil
}
