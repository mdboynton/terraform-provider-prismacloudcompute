package auth

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const (
    RolesEndpoint = "api/v1/rbac/roles"
    PermissionRadarsContainers = "radarsContainers"
    PermissionRadarsHosts = "radarsHosts"
    PermissionRadarsServerless = "radarsServerless"
    PermissionRadarsCloud = "radarsCloud"
    PermissionPolicyContainers = "policyContainers"
    PermissionPolicyHosts = "policyHosts"
    PermissionPolicyServerless = "policyServerless"
    PermissionPolicyCloud = "policyCloud"
    PermissionPolicyComplianceCustomRules = "policyComplianceCustomRules"
    PermissionPolicyRuntimeContainer = "policyRuntimeContainer"
)

type Role struct {
	Description string              `json:"description,omitempty"`
	Name        string              `json:"name,omitempty"`
    System      bool                `json:"system"`
	Permissions []RolePermission    `json:"perms,omitempty"`
}

type RolePermission struct {
	Name      string `json:"name,omitempty"`
	ReadWrite bool   `json:"readWrite,omitempty"`
}

func ListRoles(c api.PrismaCloudComputeAPIClient) ([]Role, error) {
	var ans []Role
	if err := c.Request(http.MethodGet, RolesEndpoint, nil, nil, &ans); err != nil {
		return nil, fmt.Errorf("error listing roles: %s", err)
	}
	return ans, nil
}

func GetRole(c api.PrismaCloudComputeAPIClient, name string) (*Role, error) {
	roles, err := ListRoles(c)
	if err != nil {
		return nil, err
	}
	for _, val := range roles {
		if val.Name == name {
			return &val, nil
		}
	}
	return nil, fmt.Errorf("role '%s' not found", name)
}

func CreateRole(c api.PrismaCloudComputeAPIClient, role Role) error {
	return c.Request(http.MethodPost, RolesEndpoint, nil, role, nil)
}

func UpdateRole(c api.PrismaCloudComputeAPIClient, role Role) error {
	return c.Request(http.MethodPut, RolesEndpoint, nil, role, nil)
}

func DeleteRole(c api.PrismaCloudComputeAPIClient, name string) error {
	return c.Request(http.MethodDelete, fmt.Sprintf("%s/%s", RolesEndpoint, name), nil, nil, nil)
}
