package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const (
	VmImageComplianceEndpoint         = "api/v1/policies/compliance/vms"
)

// Create/Update/Delete vm image compliance policy
func UpsertVmImageCompliancePolicy(c api.PrismaCloudComputeAPIClient, policy CompliancePolicy) error {
	return c.Request(http.MethodPut, VmImageComplianceEndpoint, nil, policy, nil)
}

// Get vm image compliance policy
func GetVmImageCompliancePolicy(c api.PrismaCloudComputeAPIClient) (*CompliancePolicy, error) {
    var ans CompliancePolicy
    if err := c.Request(http.MethodGet, VmImageComplianceEndpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting vm image compliance policy: %s", err)
    }
    
    return &ans, nil
}
