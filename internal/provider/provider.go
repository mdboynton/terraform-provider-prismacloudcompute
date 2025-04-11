package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strings"
    "strconv"
    "net/url"

    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
    //"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/auth"
    compliance "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/compliance"
    vulnerability "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/vulnerability"
    runtime "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/runtime"
    custom "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/custom"
    systemResource "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
    systemDataSource "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/data_sources/system"

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

    // If a config file was specified in the provider configuration, attempt to parse and populate the values into the config object
    if config.ConfigFile != nil {
        overwriteApiClientConfigWithFile(ctx, &config, resp)
        if resp.Diagnostics.HasError() {
            return 
        }
    // Otherwise, attempt to parse and populate the values from the relevant environment variables, if they are defined
    } else {
        overwriteApiClientConfigurationWithEnvVars(ctx, &config, resp)
        if resp.Diagnostics.HasError() {
            return 
        }
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

func overwriteApiClientConfigWithFile(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig, resp *provider.ConfigureResponse) {
    tflog.Debug(ctx, "Attempting to create provider config from file")

    // Open config file specified 
    configFile, err := os.Open(*config.ConfigFile)
    if err != nil {
        resp.Diagnostics.AddWarning(
            "Provider Configuration File Error",
            fmt.Sprintf("Error configuring provider: Configuration file specified but could not be opened. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
        )
        return
    }

    defer configFile.Close()

    // Read contents of config file
    configFileContent, err := io.ReadAll(configFile)
    if err != nil {
        resp.Diagnostics.AddWarning(
            "Provider Configuration File Error",
            fmt.Sprintf("Error configuring provider: Failed to read configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
        )
        return
    }

    // Unmarshal config file contents
    err = json.Unmarshal(configFileContent, &config)
    if err != nil {
        resp.Diagnostics.AddWarning(
            "Provider Configuration File Error",
            fmt.Sprintf("Error configuring provider: Failed to unmarshal configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
        )
    }

    return
}

// Overwrite API client configuration values with values from environment variables if they're set, non-empty and valid
func overwriteApiClientConfigurationWithEnvVars(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig, resp *provider.ConfigureResponse) {
    var (
        consoleUrl string
        username string
        password string
        insecureString string
        insecure bool
        //requestTimeoutString string 
        //requestTimeout int
        usedValues []string = []string{}
        emptyValues []string = []string{}
        failedParsedValues []string = []string{}
        err error
    )

    // Check each environment variable for a non-empty/valid value.
    // If a valid value is found, overwrite the relevent configuration value

    if (config.ConsoleURL == nil || *config.ConsoleURL == "") {
        consoleUrl = os.Getenv(ConsoleUrlEnvVar)
        if consoleUrl != "" {
            config.ConsoleURL = &consoleUrl
            usedValues = append(usedValues, "console_url")
        } else {
            emptyValues = append(usedValues, "console_url")
        }
    }

    if config.Username == nil {
        username = os.Getenv(UsernameEnvVar)
        if username != "" {
            config.Username = &username
            usedValues = append(usedValues, "username")
        } else {
            emptyValues = append(emptyValues, "username")
        }
    }

    if config.Password == nil {
        password = os.Getenv(PasswordEnvVar)
        if password != "" {
            config.Password = &password
            usedValues = append(usedValues, "password")
        } else {
            emptyValues = append(usedValues, "password")
        }
    }
    
    if config.Insecure == nil {
        insecureString = os.Getenv(InsecureEnvVar)
        if insecureString != "" {
            insecure, err = strconv.ParseBool(insecureString)
            if err != nil {
                failedParsedValues = append(failedParsedValues, "insecure")
            } else {
                config.Insecure = &insecure
                usedValues = append(usedValues, "insecure")
            }
        } else {
            emptyValues = append(emptyValues, "insecure")
        }
    }

    // request_timeout is optional, setting it to 60 seconds by default

    //if config.RequestTimeout == nil {
    //    requestTimeoutString = os.Getenv(RequestTimeoutEnvVar)
    //    if requestTimeoutString != "" {
    //        requestTimeout, err = strconv.Atoi(requestTimeoutString)
    //        if err != nil {
    //            failedParsedValues = append(failedParsedValues, "request_timeout")
    //        } else {
    //            config.RequestTimeout = &requestTimeout
    //            usedValues = append(usedValues, "request_timeout")
    //        }
    //    } else {
    //        emptyValues = append(emptyValues, "request_timeout")
    //    }
    //}

    if (len(failedParsedValues) > 0 || len(emptyValues) > 0) {
        errorMessage := "Error occured while attempting to populate provider configuration with environment variables for values not found in provider block/config file."
        
        if len(emptyValues) > 0 {
            errorMessage += fmt.Sprintf("\nThe following environment variables are not set or are set to empty values: %s", strings.Join(emptyValues, ", "))
        }

        if len(failedParsedValues) > 0 {
            errorMessage += fmt.Sprintf("\nThe following environment variables contained invalid values: %s", strings.Join(failedParsedValues, ", "))
        }

        resp.Diagnostics.AddError(
            "Provider Configuration Error",
            errorMessage,
        )

        return
    }

    if len(usedValues) > 0 {
        tflog.Debug(ctx, fmt.Sprintf("Using environment variable values for configuration fields: %s", strings.Join(usedValues, ", ")))
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
