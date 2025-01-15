package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	//collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	authAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/auth"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"

	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	//"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
)

var _ resource.Resource = &GcpCloudAccountResource{}
var _ resource.ResourceWithImportState = &GcpCloudAccountResource{}
//var _ resource.ResourceWithModifyPlan = &CollectionResource{}

type GcpCloudAccountResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func NewGcpCloudAccountResource() resource.Resource {
    return &GcpCloudAccountResource{}
}

func (r *GcpCloudAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_gcp_cloud_account"
}

func (r *GcpCloudAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "account_name": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                // TODO: validation (cant be longer than 30 characters)
            },
            "credential_level": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                // TODO: validation
                // TODO: default
            },
            "service_account_key": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
                Sensitive: true,
                // TODO: validation
            },
            "api_token": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Sensitive: true,
                // TODO: validation
            },
            "agentless_scanning": schema.SingleNestedAttribute{
                MarkdownDescription: "TODO",
                Required: true,
                Validators: []validator.Object{
                    validators.HubAccountNotTrueIfHubCredentialIdSet(),
                },
                Attributes: map[string]schema.Attribute{
                    "enabled": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Required: true,
                    },
                    "is_hub_account": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: booldefault.StaticBool(false),
                    },
                    "hub_account_id": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        // TODO: validation
                    },
                    "console_url": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "port": schema.Int32Attribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "proxy_address": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "proxy_certificate": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "auto_scale": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: booldefault.StaticBool(false),
                    },
                    "custom_labels": schema.SetNestedAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        NestedObject: schema.NestedAttributeObject{
                            Attributes: map[string]schema.Attribute{
                                "key": schema.StringAttribute{
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                },
                                "value": schema.StringAttribute{
                                    MarkdownDescription: "TODO",
                                    Optional: true,
                                },
                            },
                        },
                    },
                    "max_number_of_scanners": schema.Int32Attribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: int32default.StaticInt32(int32(1)),
                        // TODO: validation
                    },
                    "enforce_permissions_check": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: booldefault.StaticBool(true),
                    },
                    "scan_scope": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "scan_non_running_hosts": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "scope_by_labels": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "subnet": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                },
            },
            "serverless_scanning": schema.SingleNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Attributes: map[string]schema.Attribute{
                    "enabled": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "limit": schema.Int32Attribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                },
            },
            "enable_cloud_discovery": schema.BoolAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                // TODO: default
            },
        },
    }
}


func (r *GcpCloudAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*api.PrismaCloudComputeAPIClient)

    if !ok {
        resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

        return
    }

    r.client = client
}


func (r *GcpCloudAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan models.GcpCloudAccountResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    credential, cloudScanRule, diags := cloudAccountToTerraform(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new credential 
    err := authAPI.UpdateCredential(*r.client, credential)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Cloud Account resource", 
            "Failed to create credential: " + err.Error(),
        )
        return
	}

    // Retrieve newly created credential
    createdCredential, err := authAPI.GetCredential(*r.client, credential.Id)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Cloud Account resource", 
            fmt.Sprintf("Failed to retrieve created credential with ID \"%s\": ", credential.Id) + err.Error(),
        )
        return
    }

    // Create new cloud scan rule
    // TODO: update this func to accept a single CloudScanRule and do the encapsulation on the other end
    err = systemAPI.UpdateCloudScanRule(*r.client, []systemAPI.CloudScanRule{ cloudScanRule })
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Cloud Account resource", 
            "Failed to create cloud scan rule: " + err.Error(),
        )
        return
	}

    // Retrieve newly created cloud scan rule 
    createdCloudScanRule, err := systemAPI.GetCloudScanRuleByCredentialId(*r.client, cloudScanRule.CredentialId)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Cloud Account resource", 
            fmt.Sprintf("Failed to retrieve created cloud scan rule with credential ID \"%s\": ", cloudScanRule.CredentialId) + err.Error(),
        )
        return
	}
    
    state, diags := terraformToCloudAccount(ctx, plan, createdCredential, createdCloudScanRule)

    // Set state to collection data
    diags = resp.State.Set(ctx, state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *GcpCloudAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    //// Get current state
    //var state models.GcpCloudAccountResourceModel 
    //diags := req.State.Get(ctx, &state)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //// Get collection value from Prisma Cloud
    //collection, err := collectionAPI.GetCollection(*r.client, state.Name.ValueString())
    //if err != nil {
    //    resp.Diagnostics.AddError(
    //        "Error reading Collection resource", 
    //        "Failed to read collection name " + state.Name.ValueString()  + ": " + err.Error(),
    //    )
    //    return
    //}
  
    //// Overwrite state values with Prisma Cloud data
    //state, diags = collectionToSchema(ctx, *collection) 
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //// Set refreshed state
    //diags = resp.State.Set(ctx, &state)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}
}

func (r *GcpCloudAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    //// Get current state
    //var state models.GcpCloudAccountResourceModel 
    //diags := req.State.Get(ctx, &state)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //// Retrieve values from plan
    //var plan models.GcpCloudAccountResourceModel
    //diags = req.Plan.Get(ctx, &plan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //// Generate API request body from plan
    //collection, diags := schemaToCollection(ctx, &plan)

    //// Update existing collection 
	//err := collectionAPI.UpdateCollection(*r.client, state.Name.ValueString(), collection)
	//if err != nil {
	//	resp.Diagnostics.AddError(
    //        "Error updating Collection resource", 
    //        "Failed to update collection: " + err.Error(),
    //    )
    //    return
	//}

    //// Fetch updated collection from Prisma Cloud
    //updatedCollection, err := collectionAPI.GetCollection(*r.client, plan.Name.ValueString())
    //if err != nil {
    //    resp.Diagnostics.AddError(
    //        "Error updating Collection resource", 
    //        "Failed to read name" + plan.Name.ValueString()  + ": " + err.Error(),
    //    )
    //    return
    //}

    //// Convert updated collection to schema
    //plan, diags = collectionToSchema(ctx, *updatedCollection)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}
    //
    //// Set updated state
    //diags = resp.State.Set(ctx, plan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}
}

func (r *GcpCloudAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    //// Retrieve values from state
	//var state models.GcpCloudAccountResourceModel
    //diags := req.State.Get(ctx, &state)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}
    //
    //// Delete existing collection 
    //collection := state.Name.ValueString()
    //err := collectionAPI.DeleteCollection(*r.client, collection)
	//if err != nil {
	//	resp.Diagnostics.AddError(
    //        "Error deleting Collection resource", 
    //        "Failed to delete collection: " + err.Error(),
    //    )
    //    return
	//}
}

// TODO: Define ImportState to work properly with this resource
func (r *GcpCloudAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func cloudAccountToTerraform(ctx context.Context, cloudAccount models.GcpCloudAccountResourceModel) (authAPI.Credential, systemAPI.CloudScanRule, diag.Diagnostics) {
    var diags diag.Diagnostics

    //if cloudAccount == nil {
    //    diags.AddError(
    //        "Value Conversion Error",
    //        "Error occured while converting cloud account resource to Terraform types: null reference error",
    //    )
    //    return authAPI.Credential{}, systemAPI.CloudScanRule{}, diags
    //}

    /*
        Configurable fields for agentless scanning when set to "scan as hub":
        - auto scale
        - max number of scanners
        - enforce permissions check
        - custom labels
        - subnets
    */

    credential := authAPI.Credential{
        Id: cloudAccount.AccountName.ValueString(),
        AccountName: cloudAccount.AccountName.ValueString(),
        Description: cloudAccount.Description.ValueString(),
        ApiToken: authAPI.Secret{
            Encrypted: "",
            Plain: cloudAccount.ApiToken.ValueString(),
        },
        Secret: authAPI.Secret{
            Encrypted: "", // TODO: do we need this?
            Plain: cloudAccount.ServiceAccountKey.ValueString(),
            //Plain: strings.Replace(cloudAccount.ServiceAccountKey.ValueString(), "\n", "", -1),
        },
        //SkipVerify:
        Type: "gcp",
        //UseAWSRole: 
        //UseSTSRegionalEndpoint: 
    }

    // If the "Organization" setting is selected for scan scope, agentless scan spec must be "{}" 
    cloudScanRule := systemAPI.CloudScanRule{
        CredentialId: cloudAccount.AccountName.ValueString(),    
        DiscoveryEnabled: cloudAccount.EnableCloudDiscovery.ValueBool(),
        AwsRegionType: "regular",
        AgentlessScanSpec: systemAPI.AgentlessScanSpec{
            Enabled: cloudAccount.AgentlessScanning.Enabled.ValueBool(),
            HubAccount: cloudAccount.AgentlessScanning.IsHubAccount.ValueBool(),
            HubCredentialID: cloudAccount.AgentlessScanning.HubAccountId.ValueString(),
            ConsoleAddress: fmt.Sprintf("%s:%d", cloudAccount.AgentlessScanning.ConsoleURL.ValueString(), cloudAccount.AgentlessScanning.Port.ValueInt32()),
            //ProxyAddress
            //ProxyCA
            AutoScale: cloudAccount.AgentlessScanning.AutoScale.ValueBool(),
            Scanners: int(cloudAccount.AgentlessScanning.MaxNumberOfScanners.ValueInt32()), // "Agentless target account with hub account scan cannot be configured to specify max scanners"
            ScanNonRunning: cloudAccount.AgentlessScanning.ScanNonRunningHosts.ValueBool(),
            //SecurityGroup: cloudAccount.AgentlessScanning.SecurityGroup.ValueString(), // This might just be for AWS?
            SkipPermissionsCheck: !cloudAccount.AgentlessScanning.EnforcePermissionsCheck.ValueBool(),
            Subnet: cloudAccount.AgentlessScanning.Subnet.ValueString(),
        },
        ServerlessScanSpec: systemAPI.ServerlessScanSpec{},
    }

    return credential, cloudScanRule, diags
}

func terraformToCloudAccount(ctx context.Context, planCloudAccount models.GcpCloudAccountResourceModel, credential *authAPI.Credential, cloudScanRule *systemAPI.CloudScanRule) (models.GcpCloudAccountResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    //if planCloudAccount == nil {
    //    diags.AddError(
    //        "Value Conversion Error",
    //        "Error while converting Cloud Account resource to schema type: nil value provided for planned value",
    //    )
    //    return models.GcpCloudAccountResourceModel{}, diags
    //}

    var (
        consoleUrl string
        consolePort int
    )

    // TODO: add validation to make sure the value of ConsoleAddress follows the format "https://example.com:1234"
    splitConsoleAddress := strings.Split(cloudScanRule.AgentlessScanSpec.ConsoleAddress, ":")
    consoleUrl = fmt.Sprintf("%s:%s", splitConsoleAddress[0], splitConsoleAddress[1])

    consolePort, err := strconv.Atoi(splitConsoleAddress[2])
    if err != nil {
        // TODO: add error to diags
        return models.GcpCloudAccountResourceModel{}, diags
    }

    cloudAccount := models.GcpCloudAccountResourceModel{
        AccountName: types.StringValue(cloudScanRule.Credential.AccountName),
        Description: types.StringValue(cloudScanRule.Credential.Description),
        CredentialLevel: types.StringValue("project"), // TODO: find out how to actually populate this
        ServiceAccountKey: planCloudAccount.ServiceAccountKey, // Populate with planned value since API only returns encoded value
        ApiToken: planCloudAccount.ApiToken, // Populate with planned value since API only returns encoded value
        AgentlessScanning: models.AgentlessScanningResourceModel{
            Enabled: types.BoolValue(true), // TODO: find out how to actually populate this
            HubAccountId: types.StringValue(cloudScanRule.AgentlessScanSpec.HubCredentialID),
            ConsoleURL: types.StringValue(consoleUrl),
            Port: types.Int32Value(int32(consolePort)),
            ProxyAddress: types.StringValue(cloudScanRule.AgentlessScanSpec.ProxyAddress),
            ProxyCertificate: types.StringValue(cloudScanRule.AgentlessScanSpec.ProxyCA),
            //Scanners
            ScanNonRunningHosts: types.BoolValue(cloudScanRule.AgentlessScanSpec.ScanNonRunning),
            //SkipPermissionsCheck: types.BoolValue(cloudScanRule.AgentlessScanSpec.SkipPermissionsCheck),
            Subnet: types.StringValue(cloudScanRule.AgentlessScanSpec.Subnet),
        },
        ServerlessScanning: models.ServerlessScanningResourceModel{},
        EnableCloudDiscovery: types.BoolValue(cloudScanRule.DiscoveryEnabled),
    }

    return cloudAccount, diags
}
