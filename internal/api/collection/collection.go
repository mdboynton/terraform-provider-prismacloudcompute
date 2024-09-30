package collection

import (
	"fmt"
	"net/http"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const CollectionsEndpoint = "api/v1/collections"

type Collection struct {
    AccountIDs  []string `json:"accountIDs,omitempty" tfsdk:"account_ids"`
    AppIDs      []string `json:"appIDs,omitempty" tfsdk:"app_ids"`
    Clusters    []string `json:"clusters,omitempty" tfsdk:"clusters"`
    Color       string   `json:"color,omitempty" tfsdk:"color"`
    Containers  []string `json:"containers,omitempty" tfsdk:"containers"`
    Description string   `json:"description,omitempty" tfsdk:"description"`
    Functions   []string `json:"functions,omitempty" tfsdk:"functions"`
    Hosts       []string `json:"hosts,omitempty" tfsdk:"hosts"`
    Images      []string `json:"images,omitempty" tfsdk:"images"`
    Labels      []string `json:"labels,omitempty" tfsdk:"labels"`
    Modified    string   `json:"modified,omitempty" tfsdk:"modified"`
    Name        string   `json:"name,omitempty" tfsdk:"name"`
    Namespaces  []string `json:"namespaces,omitempty" tfsdk:"namespaces"`
    Owner       string   `json:"owner,omitempty" tfsdk:"owner"`
    Prisma      bool     `json:"prisma,omitempty" tfsdk:"prisma"`
    System      bool     `json:"system,omitempty" tfsdk:"system"`
}

// Get all collections.
func ListCollections(c api.PrismaCloudComputeAPIClient) ([]Collection, error) {
	var ans []Collection
	if err := c.Request(http.MethodGet, CollectionsEndpoint, nil, nil, &ans); err != nil {
		return nil, fmt.Errorf("error listing collections: %s", err)
	}
	return ans, nil
}

// Get a specific collection.
func GetCollection(c api.PrismaCloudComputeAPIClient, name string) (*Collection, error) {
	collections, err := ListCollections(c)
	if err != nil {
		return nil, err
	}
	for _, val := range collections {
		if val.Name == name {
			return &val, nil
		}
	}
	return nil, fmt.Errorf("collection '%s' not found", name)
}

// Create a new collection.
func CreateCollection(c api.PrismaCloudComputeAPIClient, collection Collection) error {
	return c.Request(http.MethodPost, CollectionsEndpoint, nil, collection, nil)
}

// Update an existing collection.
func UpdateCollection(c api.PrismaCloudComputeAPIClient, name string, collection Collection) error {
	return c.Request(http.MethodPut, fmt.Sprintf("%s/%s", CollectionsEndpoint, name), nil, collection, nil)
}

// Delete an existing collection.
func DeleteCollection(c api.PrismaCloudComputeAPIClient, name string) error {
	return c.Request(http.MethodDelete, fmt.Sprintf("%s/%s", CollectionsEndpoint, name), nil, nil, nil)
}
