package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const (
	HostComplianceEndpoint         = "api/v1/policies/compliance/host"
)

// Create/Update/Delete host compliance policy
func UpsertHostCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy CompliancePolicy) error {
	return c.Request(http.MethodPut, HostComplianceEndpoint, nil, policy, nil)
}

// Get host compliance policy
func GetHostCompliancePolicy(c api.PrismaCloudComputeAPIClient) (*CompliancePolicy, error) {
    var ans CompliancePolicy
    if err := c.Request(http.MethodGet, HostComplianceEndpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting host compliance policy: %s", err)
    }
    
    return &ans, nil
}
