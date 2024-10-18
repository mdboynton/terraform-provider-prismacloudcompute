package policy

import (
    "fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const ApplicationControlEndpoint = "api/v1/application-control/host"

//type ApplicationControlPolicy struct {
//    Rules   
//}

type ApplicationControlPolicyRule struct {
    Id              int                                         `json:"_id"`
    Applications    []ApplicationControlPolicyRuleApplication   `json:"applications"` 
    Description     string                                      `json:"description"` 
    //Disabled        string                                      `json:"disabled"` 
    Modified        string                                      `json:"modified"` 
    Name            string                                      `json:"name"` 
    Notes           string                                      `json:"notes"` 
    Owner           string                                      `json:"owner"` 
    PreviousName    string                                      `json:"previousName"` 
    Severity        string                                      `json:"severity"` 
}

type ApplicationControlPolicyRuleApplication struct {
    Name            string      `json:"name" tfsdk:"name"`
    AllowedVersions [][]string  `json:"allowedVersions" tfsdk:"allowed_versions"` 
}

//type ApplicationControlPolicyRuleApplicationStruct struct {
//    AllowedVersions kj
//}

// Create/Update application control policy rule
func UpsertApplicationControlPolicyRule(c api.PrismaCloudComputeAPIClient, policyRule ApplicationControlPolicyRule) error {
    if err := c.Request(http.MethodPut, ApplicationControlEndpoint, nil, policyRule, nil); err != nil {
		return fmt.Errorf("error upserting application control policy rule: %s", err)
    }
    
    return nil
}

// Get application control policy
func GetApplicationControlPolicy(c api.PrismaCloudComputeAPIClient) (*[]ApplicationControlPolicyRule, error) {
    var ans []ApplicationControlPolicyRule 

    if err := c.Request(http.MethodGet, ApplicationControlEndpoint, nil, nil, &ans); err != nil {
		return &ans, fmt.Errorf("error getting application control policy: %s", err)
    }
    
    return &ans, nil
}

// Delete application control policy rule
func DeleteApplicationControlPolicyRule(c api.PrismaCloudComputeAPIClient, policyRule ApplicationControlPolicyRule) error {
    if err := c.Request(http.MethodDelete, fmt.Sprintf("%s/%d", ApplicationControlEndpoint, policyRule.Id), nil, nil, nil); err != nil {
		return fmt.Errorf("error deleting application control policy rule: %s", err)
    }
    
    return nil
}
