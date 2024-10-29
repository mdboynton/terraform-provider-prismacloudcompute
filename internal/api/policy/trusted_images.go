package policy

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
)

const TrustedImagesEndpoint = "api/v1/trust/data"

type TrustedImages struct {
    Groups  []TrustGroup        `json:"groups" tfsdk:"groups"`
    Policy  TrustedImagesPolicy `json:"policy" tfsdk:"policy"`
}

type TrustGroup struct {
    Id              string      `json:"_id" tfsdk:"id"`
    Images          []string    `json:"images" tfsdk:"images"`
    Modified        string      `json:"modified" tfsdk:"modified"`
    Name            string      `json:"name" tfsdk:"name"`
    Owner           string      `json:"owner" tfsdk:"owner"`
    PreviousName    string      `json:"previousName" tfsdk:"previous_name"`
}

type TrustedImagesPolicy struct {
    Id          string                      `json:"_id" tfsdk:"id"` 
    Enabled     bool                        `json:"enabled" tfsdk:"enabled"`
    Rules       []TrustedImagesPolicyRule   `json:"rules" tfsdk:"rules"`
}

type TrustedImagesPolicyRule struct {
    Action          []string                `json:"action" tfsdk:"action"` 
    AllowedGroups   []string                `json:"allowedGroups" tfsdk:"allowed_groups"`
    Collections     []collection.Collection `json:"collections" tfsdk:"collections"`
    DeniedGroups    []string                `json:"deniedGroups" tfsdk:"denied_groups"`
    Effect          string                  `json:"effect" tfsdk:"effect"`
    Modified        *string                  `json:"modified" tfsdk:"modified"`
    Name            string                  `json:"name" tfsdk:"name"`
    Notes           string                  `json:"notes" tfsdk:"notes"`
    Owner           string                  `json:"owner" tfsdk:"owner"`
    PreviousName    string                  `json:"previousName" tfsdk:"previous_name"`
}

func GetTrustedImagesPolicy(c api.PrismaCloudComputeAPIClient) (TrustedImages, error) {
	var ans TrustedImages
	if err := c.Request(http.MethodGet, TrustedImagesEndpoint, nil, nil, &ans); err != nil {
		return ans, fmt.Errorf("Error retrieving trusted images: %s", err)
	}
	return ans, nil
}

func UpsertTrustedImagesPolicy(c api.PrismaCloudComputeAPIClient, trustedImages TrustedImages) (error) {
    if err := c.Request(http.MethodPut, TrustedImagesEndpoint, nil, trustedImages, nil); err != nil {
        return fmt.Errorf("Error upserting trusted images: %s", err)
    }
    return nil 
}
