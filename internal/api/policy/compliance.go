package policy

import (
	"fmt"
	//"net/http"

	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
)

const (
    BaseEndpoint                = "api/v1/policies"
)

var (
    // Compliance endpoints
	HostComplianceEndpoint      = fmt.Sprintf("%s/compliance/host", BaseEndpoint)
	ContainerComplianceEndpoint = fmt.Sprintf("%s/compliance/container", BaseEndpoint)
	CiImageComplianceEndpoint   = fmt.Sprintf("%s/compliance/ci/images", BaseEndpoint)
	VmImageComplianceEndpoint   = fmt.Sprintf("%s/compliance/vms", BaseEndpoint)
    FunctionComplianceEndpoint  = fmt.Sprintf("%s/compliance/serverless", BaseEndpoint)
    CiFunctionComplianceEndpoint  = fmt.Sprintf("%s/compliance/ci/serverless", BaseEndpoint)
    // Vulnerabilities endpoints
	DeployedImageVulnerabilityEndpoint      = fmt.Sprintf("%s/vulnerability/images", BaseEndpoint)
	CiImageVulnerabilityEndpoint            = fmt.Sprintf("%s/vulnerability/ci/images", BaseEndpoint)
	HostVulnerabilityEndpoint               = fmt.Sprintf("%s/vulnerability/host", BaseEndpoint)
	VmImageVulnerabilityEndpoint            = fmt.Sprintf("%s/vulnerability/vms", BaseEndpoint)
	FunctionVulnerabilityEndpoint           = fmt.Sprintf("%s/vulnerability/serverless", BaseEndpoint)
	CiFunctionVulnerabilityEndpoint         = fmt.Sprintf("%s/vulnerability/ci/serverless", BaseEndpoint)
)

//func getEndpointAndPolicyName(policyType string) (string, string, error) {
//    switch policyType{
//        case "hostCompliance":
//            return HostComplianceEndpoint, "host", nil
//        case "containerCompliance":
//            return ContainerComplianceEndpoint, "container", nil
//        case "ciImagesCompliance":
//            return CiImageComplianceEndpoint, "CI image", nil
//        case "vmCompliance":
//            return VmImageComplianceEndpoint, "VM image", nil
//        case "serverlessCompliance":
//            return FunctionComplianceEndpoint, "function", nil
//        case "ciServerlessCompliance":
//            return CiFunctionComplianceEndpoint, "CI function", nil
//        default:
//            return "", "", fmt.Errorf("invalid policy type specified")
//    }
//}

//// Create/Update/Delete compliance policy
//func UpsertCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy CompliancePolicy) error {
//    endpoint, policyName, err := getEndpointAndPolicyName(policy.PolicyType)
//    if err != nil {
//		return err
//    }
//
//    if err := c.Request(http.MethodPut, endpoint, nil, policy, nil); err != nil {
//		return fmt.Errorf("error upserting %s compliance policy: %s", policyName, err)
//    }
//    
//    return nil
//}
//
//// Get compliance policy
//func GetCompliancePolicy(c api.PrismaCloudComputeAPIClient, policyType string) (*CompliancePolicy, error) {
//    var ans CompliancePolicy
//
//    endpoint, policyName, err := getEndpointAndPolicyName(policyType)
//    if err != nil {
//		return &ans, err
//    }
//    
//    if err := c.Request(http.MethodGet, endpoint, nil, nil, &ans); err != nil {
//		return &ans, fmt.Errorf("error getting %s compliance policy: %s", policyName, err)
//    }
//    
//    return &ans, nil
//}
