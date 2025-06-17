package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	//"net/url"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	systemDataSource "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/data_sources/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/auth"
	compliance "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/compliance"
	custom "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/custom"
	runtime "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/runtime"
	vulnerability "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/vulnerability"
	systemResource "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	ConsoleUrlEnvVar     = "PRISMACLOUDCOMPUTE_CONSOLE_URL"
	UsernameEnvVar       = "PRISMACLOUDCOMPUTE_USERNAME"
	PasswordEnvVar       = "PRISMACLOUDCOMPUTE_PASSWORD"
	InsecureEnvVar       = "PRISMACLOUDCOMPUTE_INSECURE"
	RequestTimeoutEnvVar = "PRISMACLOUDCOMPUTE_REQUEST_TIMEOUT"

	InsecureDefault       = true
	RequestTimeoutDefault = 60
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
	ConsoleUrl     types.String `tfsdk:"console_url"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Insecure       types.Bool   `tfsdk:"insecure"`
	RequestTimeout types.Int32  `tfsdk:"request_timeout"`
	ConfigFile     types.String `tfsdk:"config_file"`
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

//func (p *PrismaCloudComputeProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
//    util.HCLogDebug(ctx, "ValidateConfig")
//
//    var config PrismaCloudComputeProviderModel
//    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
//    if resp.Diagnostics.HasError() {
//        return
//    }
//
//    // Validate required provider configuration arguments
//    if (config.ConsoleUrl.IsNull() || config.ConsoleUrl.IsUnknown()) {
//        resp.Diagnostics.AddError(
//            "Missing Provider Configuration Attribute",
//            fmt.Sprintf("Missing value for provider configuration attribute 'console_url'. Please provide a value in the provider configuration block, in the file specified by the 'config_file' attribute, or set the %s environment variable.", ConsoleUrlEnvVar),
//        )
//    }
//
//    if (config.Username.IsNull() || config.Username.IsUnknown()) {
//        resp.Diagnostics.AddError(
//            "Missing Provider Configuration Attribute",
//            fmt.Sprintf("Missing value for provider configuration attribute 'username'. Please provide a value in the provider configuration block, in the file specified by the 'config_file' attribute, or set the %s environment variable.", UsernameEnvVar),
//        )
//    }
//
//    if (config.Password.IsNull() || config.Password.IsUnknown()) {
//        resp.Diagnostics.AddError(
//            "Missing Provider Configuration Attribute",
//            fmt.Sprintf("Missing value for provider configuration attribute 'username'. Please provide a value in the provider configuration block, in the file specified by the 'config_file' attribute, or set the %s environment variable.", UsernameEnvVar),
//        )
//    }
//
//    if resp.Diagnostics.HasError() {
//        return
//    }
//
//    // Check optional parameters and assign default value if null or unknown
//    if (config.Insecure.IsNull() || config.Insecure.IsUnknown()) {
//        resp.Diagnostics.AddWarning(
//            "Missing Optional Provider Configuration Attribute",
//            fmt.Sprintf("Missing value for provider configuration attribute 'insecure'. Using default value of %t.", InsecureDefault),
//        )
//        config.Insecure = types.BoolValue(InsecureDefault)
//    }
//
//    if (config.RequestTimeout.IsNull() || config.RequestTimeout.IsUnknown()) {
//        resp.Diagnostics.AddWarning(
//            "Missing Optional Provider Configuration Attribute",
//            fmt.Sprintf("Missing value for provider configuration attribute 'request_timeout'. Using default value of %d.", RequestTimeoutDefault),
//        )
//        config.RequestTimeout = types.Int32Value(RequestTimeoutDefault)
//    }
//}

func (p *PrismaCloudComputeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		auth.NewUserResource,
		auth.NewRoleResource,
		systemResource.NewCollectionResource,
		systemResource.NewGcpCloudAccountResource,
		systemResource.NewAlertProfileResource,
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
		custom.NewCustomRuleResource,
	}
}

func (p *PrismaCloudComputeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		systemDataSource.NewCollectionDataSource,
	}
}

func (p *PrismaCloudComputeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Starting provider configuration")

	var (
		config api.PrismaCloudComputeAPIClientConfig
		diags  diag.Diagnostics
	)

	if p.version == "test" {
		consoleUrl := os.Getenv(ConsoleUrlEnvVar)
		username := os.Getenv(UsernameEnvVar)
		password := os.Getenv(PasswordEnvVar)
		insecure := true
		requestTimeout := 60

		config = api.PrismaCloudComputeAPIClientConfig{
			ConsoleURL:     &consoleUrl,
			Username:       &username,
			Password:       &password,
			Insecure:       &insecure,
			RequestTimeout: &requestTimeout,
		}

		client, err := api.NewPrismaCloudComputeAPIClient(ctx, config)
		if err != nil {
			resp.Diagnostics.AddError("API Client Initliaization Error", err.Error())
		}

		resp.DataSourceData = client
		resp.ResourceData = client

		return
	}

	config, diags = GetProviderConfiguration(ctx, req)
	if diags.HasError() {
		return
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
	client, err := api.NewPrismaCloudComputeAPIClient(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError("API Client Initalization Error", err.Error())
	}
	tflog.Debug(ctx, "API client initialized successfully")

	resp.DataSourceData = client
	resp.ResourceData = client
}

func GetProviderConfiguration(ctx context.Context, req provider.ConfigureRequest) (api.PrismaCloudComputeAPIClientConfig, diag.Diagnostics) {
	// Retrieve configuration values from provider block
	var config api.PrismaCloudComputeAPIClientConfig
	diags := req.Config.Get(ctx, &config)
	if diags.HasError() {
		return config, diags
	}

	// Attempt to retrieve configuration values from the provided config file
	// path and overwrite the provider block values with the values from the
	// config file if they exist and can be deserialized
	diags.Append(overwriteApiClientConfigWithFile(ctx, &config)...)

	// Attempt to retrieve configuration values from environment variables and
	// overwrite the provider block values if they are succesfully retrieved
	// and validated
	diags.Append(overwriteApiClientConfigWithEnvVars(ctx, &config)...)
	if diags.HasError() {
		return config, diags
	}

	// Raise diagnostic errors for any unconfigured required values
	unconfiguredValues := getUnconfiguredValues(config)
	if len(unconfiguredValues) > 0 {
		diags.AddError(
			"Provider Configuration Error",
			fmt.Sprintf("Required provider configuration values not found in provider block, config file or environment variables: %s", strings.Join(unconfiguredValues, ", ")),
		)
	}

	return config, diags
}

func overwriteApiClientConfigWithFile(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if config == nil {
		diags.AddError(
			"Provider Configuration Error",
			"Error configuring provider: Expected *api.PrismaCloudComputeAPIClientConfig, got nil pointer. Please report this issue to the provider developers.",
		)
		return diags
	}

	if config.ConfigFile != nil {
		tflog.Debug(ctx, fmt.Sprintf("Configuring provider from file %s", *config.ConfigFile))

		// Open config file
		configFile, err := os.Open(*config.ConfigFile)
		if err != nil {
			diags.AddWarning(
				"Provider Configuration File Error",
				fmt.Sprintf("Error configuring provider: Configuration file specified but could not be opened. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
			)
			return diags
		}

		defer configFile.Close()

		// Read contents of config file
		configFileContent, err := io.ReadAll(configFile)
		if err != nil {
			diags.AddWarning(
				"Provider Configuration File Error",
				fmt.Sprintf("Error configuring provider: Failed to read configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
			)
			return diags
		}

		// Unmarshal config file contents
		err = json.Unmarshal(configFileContent, &config)
		if err != nil {
			diags.AddWarning(
				"Provider Configuration File Error",
				fmt.Sprintf("Error configuring provider: Failed to unmarshal configuration file. Provider will default to using configuration values in provider block or environment variables.\nError: %s", err),
			)
		}
	}

	return diags
}

// Overwrite API client configuration values with values from environment variables if they're set, non-empty and valid
func overwriteApiClientConfigWithEnvVars(ctx context.Context, config *api.PrismaCloudComputeAPIClientConfig) diag.Diagnostics {
	var (
		diags          diag.Diagnostics
		consoleUrl     string
		username       string
		password       string
		insecure       bool
		requestTimeout int
		err            error
	)

	// For each nil/empty provider configuration parameter, check its
	// respective environment variable for a non-empty/valid value. If a valid
	// value is found, overwrite the relevent configuration value. Otherwise,
	// raise an error, as this is the final place where this value can be
	// retrieved from.
	if util.IsNilOrEmpty(config.ConsoleURL) {
		err = util.GetEnvironmentVariable(ConsoleUrlEnvVar, &consoleUrl)
		if err != nil {
			diags.AddError(
				"Provider Configuration Error",
				fmt.Sprintf("Error occured while attempting to parse provider configuration value \"console_url\" from environment variable \"%s\": %s", ConsoleUrlEnvVar, err.Error()),
			)
		} else {
			util.HCLogInfo(ctx, fmt.Sprintf("Using environment variable value \"%s\" for provider configuration value \"console_url\"", consoleUrl))
		}
	}

	if util.IsNilOrEmpty(config.Username) {
		err = util.GetEnvironmentVariable(UsernameEnvVar, &username)
		if err != nil {
			diags.AddError(
				"Provider Configuration Error",
				fmt.Sprintf("Error occured while attempting to parse provider configuration value \"username\" from environment variable \"%s\": %s", UsernameEnvVar, err.Error()),
			)
		} else {
			util.HCLogInfo(ctx, "Using environment variable for provider configuration value \"username\"")
		}
	}

	if util.IsNilOrEmpty(config.Password) {
		err = util.GetEnvironmentVariable(PasswordEnvVar, &password)
		if err != nil {
			diags.AddError(
				"Provider Configuration Error",
				fmt.Sprintf("Error occured while attempting to parse provider configuration value \"password\" from environment variable \"%s\": %s", PasswordEnvVar, err.Error()),
			)
		} else {
			util.HCLogInfo(ctx, "Using environment variable for provider configuration value \"password\"")
		}
	}

	if config.Insecure == nil {
		err = util.GetEnvironmentVariable(InsecureEnvVar, &insecure)
		if err != nil {
			diags.AddWarning(
				"Provider Configuration Error",
				fmt.Sprintf("Error occured while attempting to parse provider configuration value \"insecure\" from environment variable \"%s\": %s", InsecureEnvVar, err.Error()),
			)
		} else {
			util.HCLogInfo(ctx, "Using environment variable for provider configuration value \"insecure\"")
		}
	}

	if config.RequestTimeout == nil {
		err = util.GetEnvironmentVariable(RequestTimeoutEnvVar, &requestTimeout)
		if err != nil {
			diags.AddWarning(
				"Provider Configuration Error",
				fmt.Sprintf("Error occured while attempting to parse provider configuration value \"request_timeout\" from environment variable \"%s\": %s", RequestTimeoutEnvVar, err.Error()),
			)
		} else {
			util.HCLogInfo(ctx, "Using environment variable for provider configuration value \"request_timeout\"")
		}
	}

	return diags
}

func getUnconfiguredValues(config api.PrismaCloudComputeAPIClientConfig) []string {
	unconfiguredValues := []string{}

	if util.IsNilOrEmpty(config.ConsoleURL) {
		unconfiguredValues = append(unconfiguredValues, "console_url")
	}

	if util.IsNilOrEmpty(config.Username) {
		unconfiguredValues = append(unconfiguredValues, "username")
	}

	if util.IsNilOrEmpty(config.Password) {
		unconfiguredValues = append(unconfiguredValues, "password")
	}

	return unconfiguredValues
}

//func validateConsoleUrl(consoleUrl *string) *string {
//    if consoleUrl == nil {
//        errorMessage := "Error occured while attempting to parse console URL: nil pointer reference"
//        return &errorMessage
//    }
//
//    parsedConsoleUrl, err := url.Parse(*consoleUrl)
//    if err != nil {
//        errorMessage := fmt.Sprintf("Error occured while attempting to parse console URL: %s", err.Error())
//        return &errorMessage
//    }
//
//    // Set URL scheme to https if not specified
//    if parsedConsoleUrl.Scheme == "" {
//        parsedConsoleUrl.Scheme = "https"
//    }
//
//    resp := parsedConsoleUrl.String()
//    consoleUrl = &resp
//
//    return nil
//}
