package acceptance

import (
	"context"
	"testing"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	p "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)


const (
    configFileEnvVar = "PRISMACLOUDCOMPUTE_CONFIG_FILE"
    credentialsFilePath = "../../examples/creds.json"
)

var (
    providerFactory map[string]func() (tfprotov6.ProviderServer, error)
    client *api.PrismaCloudComputeAPIClient
)

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
    providerFactory = map[string]func() (tfprotov6.ProviderServer, error) {
        "prismacloudcompute": providerserver.NewProtocol6WithError(p.New("test")()),
    }

    return providerFactory
}

func TestProvider(t *testing.T) {
    ctx := context.Background()
    provider := p.New("test")()

    // Create the provider server 
    providerServer, err := createProviderServer(provider)
	if err != nil {
		t.Fatalf("Failed to create provider server: %s", err)
	}

    // Perform config validation
	validateResponse, err := providerServer.ValidateProviderConfig(ctx, &tfprotov6.ValidateProviderConfigRequest{})
	if err != nil {
		t.Fatalf("Provider config validation failed, error: %v", err)
	}

	if hasError(validateResponse.Diagnostics) {
		t.Fatalf("Provider config validation failed, diagnostics: %v", validateResponse.Diagnostics)
	}
}

func createProviderServer(provider provider.Provider) (tfprotov6.ProviderServer, error) {
	server, err := providerserver.NewProtocol6WithError(provider)()
	if err != nil {
        // TODO:
	}

	return server, nil
}

func hasError(diagnostics []*tfprotov6.Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			return true
		}
	}
	return false
}
