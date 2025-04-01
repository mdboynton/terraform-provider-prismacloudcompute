package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strings"
    "strconv"
    "math"
    "net/url"

    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/auth"
    compliance "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/compliance"
    vulnerability "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/vulnerability"
    runtime "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/runtime"
    custom "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/custom"
    systemResource "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
    systemDataSource "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/data_sources/system"
    //"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
    ConsoleUrlEnvVar = "PRISMACLOUDCOMPUTE_CONSOLE_URL"
    UsernameEnvVar = "PRISMACLOUDCOMPUTE_USERNAME"
    PasswordEnvVar = "PRISMACLOUDCOMPUTE_PASSWORD"
    InsecureEnvVar = "PRISMACLOUDCOMPUTE_INSECURE"
    RequestTimeoutEnvVar = "PRISMACLOUDCOMPUTE_REQUEST_TIMEOUT"
)

var (
    _ provider.Provider = &PrismaCloudComputeProvider{}
)

func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &PrismaCloudComputeProvider{
            version: version,
        }
    }
}

type PrismaCloudComputeProvider struct {
    version string
}

type PrismaCloudComputeProviderModel struct {
    ConsoleUrl      types.String `tfsdk:"console_url"`
    Username        types.String `tfsdk:"username"`
    Password        types.String `tfsdk:"password"`
    Insecure        types.Bool   `tfsdk:"insecure"`
    RequestTimeout  types.Int32  `tfsdk:"request_timeout"`
    ConfigFile      types.String `tfsdk:"config_file"`
}

func (p *PrismaCloudComputeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "console_url": schema.StringAttribute{
                Optional:    true,
                Description: "URL for Prisma Cloud Compute console. Do not include anything after the hostname.",
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
            "request_timeout": schema.Int32Attribute{
                Optional:    true,
                Description: "Time in seconds to wait for API requests to return before timing out",
            },
            "config_file": schema.StringAttribute{
                Optional:    true,
                Description: "Configuration file in JSON format. See examples/creds.json",
            },
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
        systemResource.NewCollectionResource,
        systemResource.NewGcpCloudAccountResource,
        // Compliance policies
        compliance.NewHostCompliancePolicyResource,
        compliance.NewContainerCompliancePolicyResource,
        compliance.NewCiImageCompliancePolicyResource,
        compliance.NewVmImageCompliancePolicyResource,
        compliance.NewFunctionCompliancePolicyResource,
        compliance.NewCiFunctionCompliancePolicyResource,
        compliance.NewApplicationControlPolicyResource,
        compliance.NewTrustedImagesPolicyResource,
        compliance.NewCustomComplianceCheckResource,
        // Vulnerability policies
        // TODO: need to rename deployed and ci images policy structs to be singular (Images -> Image)
        vulnerability.NewDeployedImagesVulnerabilityPolicyResource,
        vulnerability.NewCiImagesVulnerabilityPolicyResource,
        vulnerability.NewHostVulnerabilityPolicyResource,
        vulnerability.NewVmImageVulnerabilityPolicyResource,
        vulnerability.NewFunctionVulnerabilityPolicyResource,
        vulnerability.NewCiFunctionVulnerabilityPolicyResource,
        // Runtime policies
        runtime.NewHostRuntimePolicyResource,
        runtime.NewContainerRuntimePolicyResource,
        runtime.NewServerlessRuntimePolicyResource,
        runtime.NewAppEmbeddedRuntimePolicyResource,
        // Custom rules
        custom.NewCustomRuntimeRuleResource,
    }
}

func (p *PrismaCloudComputeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        systemDataSource.NewCollectionDataSource, 
    }
}

func (p *PrismaCloudComputeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    tflog.Debug(ctx, "Starting provider configuration")

    var config api.PrismaCloudComputeAPIClientConfig
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    if p.version == "test" {
        consoleUrl := os.Getenv(ConsoleUrlEnvVar)
        username := os.Getenv(UsernameEnvVar)
        password := os.Getenv(PasswordEnvVar)
        insecure := true
        requestTimeout := 60

        config = api.PrismaCloudComputeAPIClientConfig{
            ConsoleURL: &consoleUrl,
            Username: &username,
            Password: &password,
            Insecure: &insecure,
            RequestTimeout: &requestTimeout,
        }

        client, err := api.Client(config)
        if err != nil {
            resp.Diagnostics.AddError("API Client Configuration Error", err.Error())
        }

        resp.DataSourceData = client
        resp.ResourceData = client

        return
    }

    // Attempt to read in configuration values from file, using provider block or environment variables if not provided
    // TODO: make it more obvious whats happening here
    if config.ConfigFile != nil && !createConfigurationFromFile(ctx, &config, resp) {
        overwriteConfigurationWithEnvVars(ctx, &config, resp)
    }
   

    // TODO: fix this 
    //// Overwrite console URL value in config
    //err := validateConsoleUrl(config.ConsoleURL)
    //if err != nil {
    //    resp.Diagnostics.AddError("API Client Configuration Error", *err)
    //}

    // Set request timeout to 60 if not specified in provider configuration
    if config.RequestTimeout == nil {
        tflog.Warn(ctx, "No request timeout configured. Using default value of 60 seconds.")
        defaultTimeout := 60
        config.RequestTimeout = &defaultTimeout
    }

    // Initialize API client
    tflog.Debug(ctx, "Provider config created, initializing API client")
    client, err := api.Client(config)
    if err != nil {
        resp.Diagnostics.AddError("API Client Configuration Error", err.Error())
    }
    tflog.Debug(ctx, "API client initialized successfully")

    resp.DataSourceData = client
    resp.ResourceData = client
}

//func GetConfigurationFromFile(filePath string) (*api.PrismaCloudComputeAPIClientConfig, error) {
//    // Open config file specified 
//    configFile, err := os.Open(filePath)
//    if err != nil {
//        return nil, fmt.Errorf(
//            "Provider Configuration File Error",
//            fmt.Sprintf("Error configuring provider: Configuration file specified but could not be opened. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
//        )
//    }
//
//    defer configFile.Close()
//
//    // Read contents of config file
//    configFileContent, err := io.ReadAll(configFile)
//    if err != nil {
//
//        return nil, fmt.Errorf(
//            "Provider Configuration File Error",
//            fmt.Sprintf("Error configuring provider: Failed to read configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
//        )
//    }
//
//    // Unmarshal config file contents
//    config := api.PrismaCloudComputeAPIClientConfig{}
//    err = json.Unmarshal(configFileContent, &config)
//    if err != nil {
//        return nil, fmt.Errorf(
//            "Provider Configuration File Error",
//            fmt.Sprintf("Error configuring provider: Failed to unmarshal configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
//        )
//    }
//
//    return &config, nil
//}

func createConfigurationFromFile(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig, resp *provider.ConfigureResponse) bool {
    tflog.Debug(ctx, "Attempting to create provider config from file")

    // Open config file specified 
    configFile, err := os.Open(*config.ConfigFile)
    if err != nil {
        resp.Diagnostics.AddWarning(
            "Provider Configuration File Error",
            fmt.Sprintf("Error configuring provider: Configuration file specified but could not be opened. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
        )
        return false
    }

    defer configFile.Close()

    // Read contents of config file
    configFileContent, err := io.ReadAll(configFile)
    if err != nil {
        resp.Diagnostics.AddWarning(
            "Provider Configuration File Error",
            fmt.Sprintf("Error configuring provider: Failed to read configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
        )
        return false
    }

    // Unmarshal config file contents
    err = json.Unmarshal(configFileContent, &config)
    if err != nil {
        resp.Diagnostics.AddWarning(
            "Provider Configuration File Error",
            fmt.Sprintf("Error configuring provider: Failed to unmarshal configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
        )
        return false
    }

    return true
}

func overwriteConfigurationWithEnvVars(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig, resp *provider.ConfigureResponse) {
    envVarValues := make([]string, 0, 3)

    // Check each environment variable for provider, overwriting the provider configuration values from Terraform if set

    // TODO: add checks to see if env vars are set
    if config.ConsoleURL == nil {
        consoleUrl := os.Getenv(ConsoleUrlEnvVar)
        config.ConsoleURL = &consoleUrl
        envVarValues = append(envVarValues, "console_url")
    }

    if config.Username == nil {
        username := os.Getenv(UsernameEnvVar)
        config.Username = &username
        envVarValues = append(envVarValues, "username")
    }

    if config.Password == nil {
        password := os.Getenv(PasswordEnvVar)
        config.Password = &password
        envVarValues = append(envVarValues, "password")
    }
    
    if config.Insecure == nil {
        insecureString := os.Getenv(InsecureEnvVar)
        if insecureString != "" {
            insecure, err := strconv.ParseBool(insecureString)
            if err != nil {
                // TODO
                //resp.Diagnostics.AddError(
                //    "Provider Configuration Error",
                //    fmt.Sprintf("Error configuring provider: Invalid value specified for \"request_timeout\" in configuration file. Value must be an integer between 1 and %d", math.MaxInt),
                //)
                return
            }
            config.Insecure = &insecure
            envVarValues = append(envVarValues, "insecure")
        }
    }

    if config.RequestTimeout == nil {
        requestTimeoutStr := os.Getenv(RequestTimeoutEnvVar)
        if requestTimeoutStr != "" {
            requestTimeout, err := strconv.Atoi(requestTimeoutStr)
            if err != nil {
                resp.Diagnostics.AddError(
                    "Provider Configuration Error",
                    fmt.Sprintf("Error configuring provider: Invalid value specified for \"request_timeout\" in configuration file. Value must be an integer between 1 and %d", math.MaxInt),
                )
            }
            config.RequestTimeout = &requestTimeout
            envVarValues = append(envVarValues, "password")
        }
    }

    if len(envVarValues) > 0 {
        tflog.Debug(ctx, fmt.Sprintf("Using environment variable values for configuration fields: %s", strings.Join(envVarValues, ", ")))
    }

    // Raise errors if required values are not configured in configuration file, provider block or environment variables
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

func validateConsoleUrl(consoleUrl *string) *string {
    if consoleUrl == nil {
        errorMessage := "Error occured while attempting to parse console URL: nil pointer reference"
        return &errorMessage 
    }

    parsedConsoleUrl, err := url.Parse(*consoleUrl) 
    if err != nil {
        errorMessage := fmt.Sprintf("Error occured while attempting to parse console URL: %s", err.Error())
        return &errorMessage
    }

    // Set URL scheme to https if not specified
    if parsedConsoleUrl.Scheme == "" {
        parsedConsoleUrl.Scheme = "https"
    }

    resp := parsedConsoleUrl.String()
    consoleUrl = &resp

    return nil
}
