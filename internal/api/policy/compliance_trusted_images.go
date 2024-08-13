package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
)

const TrustedImagesEndpoint = "api/v1/policies/runtime/host"
const TrustGroupsEndpoint = "api/v1/trust/data"

type TrustedImagesPolicy struct {
	Rules TrustedImagesPolicyObject `json:"policy,omitempty"`
	Groups []TrustGroup `json:"groups,omitempty"`
}

type TrustedImagesPolicyObject struct {
    Id  int `json:"_id,omitempty"`
    Enabled bool `json:"enabled,omitempty"`
    Rules []TrustedImageRule `json:"rules,omitempty"`
}

type TrustedImageRule struct {
    Modified      string   `json:"modified,omitempty"` 
    Owner         string   `json:"owner,omitempty"`
    Name          string   `json:"name,omitempty"`
    PreviousName  string   `json:"previousName,omitempty"`
    Disabled      bool     `json:"disabled,omitempty"`
    AllowedGroups []string `json:"allowedGroups,omitempty"`
    DeniedGroups  []string `json:"deniedGroups,omitempty"`
    Collections   []collection.Collection `json:"collections,omitempty"`
    Effect        string   `json:"effect,omitempty"`
}

type TrustGroup struct {
    Modified      string   `json:"modified,omitempty"` 
    Owner         string   `json:"owner,omitempty"`
    Name          string   `json:"name,omitempty"`
    PreviousName  string   `json:"previousName,omitempty"`
    Id            string   `json:"_id,omitempty"`
    Images        []string `json:"images,omitempty"`
}

func GetTrustedImages(c api.Client) (TrustedImagesPolicy, error) {
	var res TrustedImagesPolicy 
	if err := c.Request(http.MethodGet, TrustedImagesEndpoint, nil, nil, &res); err != nil {
		return res, fmt.Errorf("Error retrieving Trusted Images policy: %s", err)
	}
	return res, nil
}

func UpdateRuntimeHost(c api.Client, policy TrustedImagesPolicy) error {
	return c.Request(http.MethodPut, TrustedImagesEndpoint, nil, policy, nil)
}

