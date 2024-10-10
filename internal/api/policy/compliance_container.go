package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
)

const (
	ContainerComplianceEndpoint         = "api/v1/policies/compliance/container"
)

// Create/Update/Delete container compliance policy
func UpsertContainerCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy CompliancePolicy) error {
	return c.Request(http.MethodPut, ContainerComplianceEndpoint, nil, policy, nil)
}

// Get container compliance policy
func GetContainerCompliancePolicy(c api.PrismaCloudComputeAPIClient) (*CompliancePolicy, error) {
    var ans CompliancePolicy
    if err := c.Request(http.MethodGet, ContainerComplianceEndpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting container compliance policy: %s", err)
    }
    
    return &ans, nil
}
