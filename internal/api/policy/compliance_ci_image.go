package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const (
	CiImageComplianceEndpoint         = "api/v1/policies/compliance/ci/images"
)

// Create/Update/Delete CI image compliance policy
func UpsertCiImageCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy CompliancePolicy) error {
	return c.Request(http.MethodPut, CiImageComplianceEndpoint, nil, policy, nil)
}

// Get CI image compliance policy
func GetCiImageCompliancePolicy(c api.PrismaCloudComputeAPIClient) (*CompliancePolicy, error) {
    var ans CompliancePolicy
    if err := c.Request(http.MethodGet, CiImageComplianceEndpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting CI image compliance policy: %s", err)
    }
    
    return &ans, nil
}
