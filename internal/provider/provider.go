package provider

import (
    "context"
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	//"os"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
	//"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure KubernetesProvider satisfies various provider interfaces.
var (
	_ provider.Provider              = &PrismaCloudComputeProvider{}
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
    ConsoleUrl types.String `tfsdk:"console_url"`
    //Project types.String `tfsdk:"project"`
    Username types.String `tfsdk:"username"`
    Password types.String `tfsdk:"password"`
    Insecure types.Bool `tfsdk:"insecure"`
    //ConfigFile types.String `tfsdk:"config_file"`
}

func (p *PrismaCloudComputeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "console_url": schema.StringAttribute{
				Optional:    true,
				Description: "The Prisma Cloud Compute Console URL",
            },
            //"project": schema.StringAttribute{
			//	Optional:    true,
			//	Description: "The Prisma Cloud Compute project",
            //},
            "username": schema.StringAttribute{
				Optional:    true,
				Description: "Prisma Cloud Compute username",
            },
            "password": schema.StringAttribute{
				Optional:    true,
				Description: "Prisma Cloud Compute password",
				Sensitive:   true,
            },
            "insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether Prisma Cloud Compute host should be accessed without verifying the TLS certificate",
            },
            //"config_file": schema.StringAttribute{
			//	Optional:    true,
			//	Description: "Configuration file in JSON format. See examples/creds.json",
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
       NewUserResource,
       NewCollectionResource,
    }
}

func (p *PrismaCloudComputeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{}
}

func (p *PrismaCloudComputeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var config api.PrismaCloudComputeAPIClientConfig
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    tflog.Debug(ctx, "provider configure: config created")

    //// TODO: Error handling for attributes

    //// Default values to environment variables, but override
    //// with Terraform configuration if set
    //consoleUrl := os.Getenv("PRISMACLOUDCOMPUTE_CONSOLE_URL")
    //username := os.Getenv("PRISMACLOUDCOMPUTE_USERNAME")
    //password := os.Getenv("PRISMACLOUDCOMPUTE_PASSWORD")
    //// TODO: add remaining provider value defaults

    //// TODO: add error handling for defaults/env vars

    //config.ConsoleUrl = consoleUrl
    //config.Username = username
    //config.Password = password

    //if !config.ConsoleUrl.IsNull() {
    //    consoleUrl = config.ConsoleUrl.ValueString()
    //}

    //if !config.Username.IsNull() {
    //    username = config.Username.ValueString()
    //}

    //if !config.Password.IsNull() {
    //    password = config.Password.ValueString()
    //}

    client, err := api.Client(config)
	if err != nil {
	    resp.Diagnostics.AddError("error creating API client", err.Error())
    }

	//return client
    resp.DataSourceData = client
    resp.ResourceData = client
}
