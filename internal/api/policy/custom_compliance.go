package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const CustomComplianceChecksEndpoint = "api/v1/custom-compliance"

type CustomComplianceCheck struct {
	Id              int    `json:"_id,omitempty"`
    Owner           string `json:"owner,omitempty"`
    Modified        string `json:"modified,omitempty"`
	Name            string `json:"name,omitempty"`
	PreviousName    string `json:"previousName,omitempty"`
	Severity        string `json:"severity,omitempty"`
	Script          string `json:"script,omitempty"`
    Title           string `json:"title,omitempty"`
}

func ListCustomComplianceChecks(c api.PrismaCloudComputeAPIClient) ([]CustomComplianceCheck, error) {
	var ans []CustomComplianceCheck
	if err := c.Request(http.MethodGet, CustomComplianceChecksEndpoint, nil, nil, &ans); err != nil {
		return nil, fmt.Errorf("error listing custom Compliances: %s", err)
	}
	return ans, nil
}

func GetCustomComplianceCheckById(c api.PrismaCloudComputeAPIClient, id int) (*CustomComplianceCheck, error) {
	compliances, err := ListCustomComplianceChecks(c)
	if err != nil {
		return nil, err
	}

	for _, val := range compliances {
		if val.Id == id {
			return &val, nil
		}
	}

	return nil, fmt.Errorf("Custom compliance check with ID \"%d\" not found", id)
}

func GetCustomComplianceCheckByName(c api.PrismaCloudComputeAPIClient, name string) (*CustomComplianceCheck, error) {
	compliances, err := ListCustomComplianceChecks(c)
	if err != nil {
		return nil, err
	}

	for _, val := range compliances {
		if val.Name == name {
			return &val, nil
		}
	}

	return nil, fmt.Errorf("Custom compliance check with name \"%s\" not found", name)
}

func UpsertCustomComplianceCheck(c api.PrismaCloudComputeAPIClient, check CustomComplianceCheck) (*CustomComplianceCheck, error) {
    var ans CustomComplianceCheck
    if err := c.Request(http.MethodPut, CustomComplianceChecksEndpoint, nil, check, &ans); err != nil {
        return nil, fmt.Errorf("Error upserting custom compliance check: %s", err)
    }
    return &ans, nil 
}

func DeleteCustomComplianceCheck(c api.PrismaCloudComputeAPIClient, id int) error {
    if err := c.Request(http.MethodDelete, fmt.Sprintf("%s/%d", CustomComplianceChecksEndpoint, id), nil, nil, nil); err != nil {
        return fmt.Errorf("Error deleting custom compliance check: %s", err)
    }

    return nil
}
