package provider

import (
    "context"
	"encoding/json"
	"io"
	"os"
    "strings"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/auth"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure PrismaCloudComputeProvider satisfies provider interface
var (
	_ provider.Provider = &PrismaCloudComputeProvider{}
)

func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &PrismaCloudComputeProvider {
            version: version,
        }
    }
}

type PrismaCloudComputeProvider struct {
    version string
}

type PrismaCloudComputeProviderModel struct {
    ConsoleUrl  types.String    `tfsdk:"console_url"`
    Username    types.String    `tfsdk:"username"`
    Password    types.String    `tfsdk:"password"`
    Insecure    types.Bool      `tfsdk:"insecure"`
    ConfigFile  types.String    `tfsdk:"config_file"`
    //Project types.String `tfsdk:"project"`
}

func (p *PrismaCloudComputeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "console_url": schema.StringAttribute{
				Optional:    true,
				Description: "The Prisma Cloud Compute Console URL",
            },
            "username": schema.StringAttribute{
				Optional:    true,
				Description: "Prisma Cloud Compute username",
            },
            "password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Prisma Cloud Compute password",
            },
            "insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether Prisma Cloud Compute host should be accessed without verifying the TLS certificate",
            },
            "config_file": schema.StringAttribute{
				Optional:    true,
				Description: "Configuration file in JSON format. See examples/creds.json",
            },
            //"project": schema.StringAttribute{
			//	Optional:    true,
			//	Description: "The Prisma Cloud Compute project",
            //},
        },
    }
}

func (p *PrismaCloudComputeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "prismacloudcompute"
    resp.Version = p.version
}

func (p *PrismaCloudComputeProvider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{
       auth.NewUserResource,
       auth.NewRoleResource,
       system.NewCollectionResource,
       policy.NewHostCompliancePolicyResource,
       policy.NewContainerCompliancePolicyResource,
       policy.NewCiImageCompliancePolicyResource,
       policy.NewVmImageCompliancePolicyResource,
       policy.NewFunctionCompliancePolicyResource,
       policy.NewCiFunctionCompliancePolicyResource,
       policy.NewApplicationControlPolicyResource,
       policy.NewCustomComplianceCheckResource,
    }
}

func (p *PrismaCloudComputeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{}
}

func (p *PrismaCloudComputeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    tflog.Debug(ctx, "Provider configure start")

    var config api.PrismaCloudComputeAPIClientConfig
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Attempt to read in configuration values from file, using provider block or environment variables if not provided
    if config.ConfigFile != nil && !createConfigurationFromFile(ctx, &config, resp) {
        overwriteConfigurationWithEnvVars(ctx, &config, resp)
    }

    if resp.Diagnostics.HasError() {
        return
    }

    tflog.Debug(ctx, "Provider config created, initializing API client")

    client, err := api.Client(config)
	if err != nil {
	    resp.Diagnostics.AddError("API Client Configuration Error", err.Error())
    }
    
    tflog.Debug(ctx, "Provider initialized API client")

    resp.DataSourceData = client
    resp.ResourceData = client
}

func createConfigurationFromFile(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig, resp *provider.ConfigureResponse) (bool) {
    tflog.Debug(ctx, "Attempting to create provider config from file")
    
    configFile, err := os.Open(*config.ConfigFile)

    if err != nil {
        resp.Diagnostics.AddWarning(
            "Provider Configuration File Error",
            fmt.Sprintf("Error configuring provider: Configuration file specified but could not be opened. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
        )
        return false
    } else {
        defer configFile.Close()

        configFileContent, err := io.ReadAll(configFile)
        if err != nil {
            resp.Diagnostics.AddWarning(
                "Provider Configuration File Error",
                fmt.Sprintf("Error configuring provider: Failed to read configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
            )
            return false
        } else {
            err = json.Unmarshal(configFileContent, &config)
            if err != nil {
                resp.Diagnostics.AddWarning(
                    "Provider Configuration File Error",
                    fmt.Sprintf("Error configuring provider: Failed to unmarshal configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
                )
                return false
            }
        }
    }
    
    return true 
}

func overwriteConfigurationWithEnvVars(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig, resp *provider.ConfigureResponse) {
    envVarValues := make([]string, 0, 3)

    // Assign environment variable values if provider block is not configured
    if config.ConsoleURL == nil {
        consoleUrl := os.Getenv("PRISMACLOUDCOMPUTE_CONSOLE_URL")
        config.ConsoleURL = &consoleUrl
        envVarValues = append(envVarValues, "console_url")
    }

    if config.Username == nil {
        username := os.Getenv("PRISMACLOUDCOMPUTE_USERNAME")
        config.Username = &username
        envVarValues = append(envVarValues, "username")
    }

    if config.Password == nil {
        password := os.Getenv("PRISMACLOUDCOMPUTE_PASSWORD")
        config.Password = &password
        envVarValues = append(envVarValues, "password")
    }

    if len(envVarValues) > 0 {
        tflog.Debug(ctx, fmt.Sprintf("Using environment variable values for configuration fields: %s", strings.Join(envVarValues, ", ")))
    }
    
    // Raise errors if configuration values are not found in configuration file/provider block/environment variables
    if *config.ConsoleURL == "" {
        resp.Diagnostics.AddError(
            "Provider Configuration Error",
            "Error configuring provider: No console URL value supplied. Specify the console URL value in a configuration file, the provider block or in the PRISMACLOUDCOMPUTE_CONSOLE_URL environment variable. Refer to provider documentation for configuration options and examples.",
        )
    }

    if *config.Username == "" {
        resp.Diagnostics.AddError(
            "Provider Configuration Error",
            "Error configuring provider: No username value supplied. Specify the username value in a configuration file, the provider block or in the PRISMACLOUDCOMPUTE_USERNAME environment variable. Refer to provider documentation for configuration options and examples.",
        )
    }
    
    if *config.Password == "" {
        resp.Diagnostics.AddError(
            "Provider Configuration Error",
            "Error configuring provider: No password value supplied. Specify the password value in a configuration file, the provider block or in the PRISMACLOUDCOMPUTE_PASSWORD environment variable. Refer to provider documentation for configuration options and examples.",
        )
    }
}
