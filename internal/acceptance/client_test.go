package acceptance

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

func TestAPIClient(t *testing.T) {
	// Open credentials file
	credsFile, err := os.Open(configFilePath)
	if err != nil {
		t.Errorf("Error opening credentials file: %v", err)
		return
	}
	defer credsFile.Close()

	// Read credentials file contents
	fileContent, err := io.ReadAll(credsFile)
	if err != nil {
		t.Errorf("Error reading credentials file: %v", err)
		return
	}

	// Unmarshal contents to config object
	var config api.PrismaCloudComputeAPIClientConfig
	if err := json.Unmarshal(fileContent, &config); err != nil {
		t.Errorf("Error unmarshalling contents of credentials file: %v", err)
		return
	}

	// Create client, initializing authentication workflow
	client, err = api.NewPrismaCloudComputeAPIClient(context.Background(), config)
	if err != nil {
		t.Errorf("Error initializing API client: %v", err)
		return
	}

	// Check for empty Java Web Token field
	if client.JWT == "" {
		t.Errorf("API client failed to authenticate (JWT empty)")
	}
}
