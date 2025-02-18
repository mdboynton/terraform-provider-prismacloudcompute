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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
                // TODO: when writing the description, make sure to warn users that changing this field will require the resource
                // to be re-made
                MarkdownDescription: "TODO",
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString(""),
                // TODO: validation (cant be longer than 30 characters)
            },
            "credential_level": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("project"),
                // TODO: validation
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
                // TODO: this should be null (or at least "enabled" should be set to "false") when credential_level is set to "organization"
                MarkdownDescription: "TODO",
                Required: true,
                Validators: []validator.Object{
                    validators.HubAccountNotTrueIfHubCredentialIdSet(),
                    validators.PermissionsCheckNotEnforcedIfScanningWithHub(),
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
                        Computed: true,
                        Default: stringdefault.StaticString(""),
                    },
                    "proxy_certificate": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: stringdefault.StaticString(""),
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
                        Default: int32default.StaticInt32(int32(10)),
                        // TODO: validation (this attribute cannot be set when creating a "scan with hub account" type account
                    },
                    "enforce_permissions_check": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: booldefault.StaticBool(false),
                    },
                    "regions": schema.SetAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        ElementType: types.StringType,
                        // TODO: validation
                    },
                    "scan_scope": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "scan_non_running_hosts": schema.BoolAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: booldefault.StaticBool(false),
                    },
                    "scope_by_labels": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                    },
                    "subnet": schema.StringAttribute{
                        MarkdownDescription: "TODO",
                        Optional: true,
                        Computed: true,
                        Default: stringdefault.StaticString(""),
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
                Computed: true,
                Default: booldefault.StaticBool(false),
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
        
        // Remove hanging credential
        err = authAPI.DeleteCredential(*r.client, createdCredential.Id)
        if err != nil {
		    resp.Diagnostics.AddWarning(
                "Error creating Cloud Account resource", 
                "Failed to delete hanging credential created prior to cloud scan rule. Error: " + err.Error(),
            )
        }

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
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set state to collection data
    diags = resp.State.Set(ctx, state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *GcpCloudAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state models.GcpCloudAccountResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get values from Prisma Cloud
    credential, err := authAPI.GetCredential(*r.client, state.AccountName.ValueString())
    if err != nil {
		resp.Diagnostics.AddError(
            "Error reading Cloud Account resource", 
            fmt.Sprintf("Failed to retrieve credential with ID \"%s\": ", state.AccountName.ValueString()) + err.Error(),
        )
        return
    }

    cloudScanRule, err := systemAPI.GetCloudScanRuleByCredentialId(*r.client, state.AccountName.ValueString())
	if err != nil {
	    resp.Diagnostics.AddError(
            "Error reading Cloud Account resource", 
            fmt.Sprintf("Failed to retrieve cloud scan rule with account name \"%s\": ", state.AccountName.ValueString()) + err.Error(),
        )
        return
	}
  
    // Convert values to schema model
    state, diags = terraformToCloudAccount(ctx, state, credential, cloudScanRule)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *GcpCloudAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state models.GcpCloudAccountResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan models.GcpCloudAccountResourceModel
    diags = req.Plan.Get(ctx, &plan)
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

    // Update credential 
    err := authAPI.UpdateCredential(*r.client, credential)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Cloud Account resource", 
            "Failed to create credential: " + err.Error(),
        )
        return
	}

    // Retrieve updated credential
    updatedCredential, err := authAPI.GetCredential(*r.client, credential.Id)
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
       
        // TODO: Can we roll back the update to the credential in the event of an error at this step?
        //err = authAPI.DeleteCredential(*r.client, createdCredential.Id)
        //if err != nil {
		//    resp.Diagnostics.AddWarning(
        //        "Error creating Cloud Account resource", 
        //        "Failed to delete hanging credential created prior to cloud scan rule. Error: " + err.Error(),
        //    )
        //}

        return
	}

    // Retrieve newly created cloud scan rule 
    updatedCloudScanRule, err := systemAPI.GetCloudScanRuleByCredentialId(*r.client, cloudScanRule.CredentialId)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Cloud Account resource", 
            fmt.Sprintf("Failed to retrieve created cloud scan rule with credential ID \"%s\": ", cloudScanRule.CredentialId) + err.Error(),
        )
        return
	}
   
    // Convert updated values to schema
    schema, diags := terraformToCloudAccount(ctx, plan, updatedCredential, updatedCloudScanRule)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set updated state
    diags = resp.State.Set(ctx, schema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *GcpCloudAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state models.GcpCloudAccountResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Delete cloud scan rule
    err := systemAPI.DeleteCloudScanRule(*r.client, state.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Cloud Account resource", 
            "Failed to delete cloud scan rule: " + err.Error(),
        )
        return
	}
    
    // Delete credential
    // TODO: Print warning about dangling credential?
    err = authAPI.DeleteCredential(*r.client, state.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Cloud Account resource", 
            "Failed to delete credential: " + err.Error(),
        )
        return
	}
}

// TODO: Ensure this works properly
func (r *GcpCloudAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("account_name"), req, resp)
}

func cloudAccountToTerraform(ctx context.Context, cloudAccount models.GcpCloudAccountResourceModel) (authAPI.Credential, systemAPI.CloudScanRule, diag.Diagnostics) {
    var diags diag.Diagnostics

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
            Encrypted: "",
            Plain: cloudAccount.ServiceAccountKey.ValueString(),
        },
        //SkipVerify:
        Type: "gcp",
        //UseAWSRole: 
        //UseSTSRegionalEndpoint: 
    }

    if cloudAccount.CredentialLevel.ValueString() == "organization" {
        credential.Global = true
    } else if cloudAccount.CredentialLevel.ValueString() == "project" {
        credential.Global = false
    } else {
        // TODO: put this in a validator
        diags.AddError(
            "Value Conversion Error",
            fmt.Sprintf("Invalid value for credential_level: %s", cloudAccount.CredentialLevel),
        )
    }

    var regions []string
    if (!cloudAccount.AgentlessScanning.Regions.IsNull() && !cloudAccount.AgentlessScanning.Regions.IsUnknown()) {
        diags = cloudAccount.AgentlessScanning.Regions.ElementsAs(ctx, &regions, false)
        if diags.HasError() {
            return authAPI.Credential{}, systemAPI.CloudScanRule{}, diags
        }
    } else {
        regions = []string{}
    }

    // If the "Organization" setting is selected for scan scope, agentless scan spec must be "{}" 
    cloudScanRule := systemAPI.CloudScanRule{
        CredentialId: cloudAccount.AccountName.ValueString(),    
        DiscoveryEnabled: cloudAccount.EnableCloudDiscovery.ValueBool(),
        AwsRegionType: "regular",
        AgentlessScanSpec: systemAPI.AgentlessScanSpec{
            Enabled: cloudAccount.AgentlessScanning.Enabled.ValueBool(),
            HubAccount: cloudAccount.AgentlessScanning.IsHubAccount.ValueBool(),
            //HubCredentialID: cloudAccount.AgentlessScanning.HubAccountId.ValueString(),
            ConsoleAddress: fmt.Sprintf("%s:%d", cloudAccount.AgentlessScanning.ConsoleURL.ValueString(), cloudAccount.AgentlessScanning.Port.ValueInt32()),
            //ProxyAddress
            //ProxyCA
            AutoScale: cloudAccount.AgentlessScanning.AutoScale.ValueBool(),
            //Scanners: int(cloudAccount.AgentlessScanning.MaxNumberOfScanners.ValueInt32()), // "Agentless target account with hub account scan cannot be configured to specify max scanners"
            Regions: regions,
            ScanNonRunning: cloudAccount.AgentlessScanning.ScanNonRunningHosts.ValueBool(),
            //SecurityGroup: cloudAccount.AgentlessScanning.SecurityGroup.ValueString(), // This might just be for AWS?
            SkipPermissionsCheck: !cloudAccount.AgentlessScanning.EnforcePermissionsCheck.ValueBool(),
            Subnet: cloudAccount.AgentlessScanning.Subnet.ValueString(),
        },
        ServerlessScanSpec: systemAPI.ServerlessScanSpec{},
    }

    if !cloudAccount.AgentlessScanning.HubAccountId.IsNull() && !cloudAccount.AgentlessScanning.HubAccountId.IsUnknown() && cloudAccount.AgentlessScanning.HubAccountId.ValueString() != "" {
        cloudScanRule.AgentlessScanSpec.HubCredentialID = cloudAccount.AgentlessScanning.HubAccountId.ValueString()
    } else {
        cloudScanRule.AgentlessScanSpec.Scanners = int(cloudAccount.AgentlessScanning.MaxNumberOfScanners.ValueInt32())
    }

    return credential, cloudScanRule, diags
}

func terraformToCloudAccount(ctx context.Context, planCloudAccount models.GcpCloudAccountResourceModel, credential *authAPI.Credential, cloudScanRule *systemAPI.CloudScanRule) (models.GcpCloudAccountResourceModel, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        credentialLevel string
        consoleUrl string
        consolePort int
        hubAccountIdValue basetypes.StringValue
        agentlessScanSpec systemAPI.AgentlessScanSpec = cloudScanRule.AgentlessScanSpec
    )

    if credential.Global {
        credentialLevel = "organization"
    } else {
        credentialLevel = "project"
    }

    // TODO: add validation to make sure the value of ConsoleAddress follows the format "https://example.com:1234"
    splitConsoleAddress := strings.Split(agentlessScanSpec.ConsoleAddress, ":")
    consoleUrl = fmt.Sprintf("%s:%s", splitConsoleAddress[0], splitConsoleAddress[1])
    consolePort, err := strconv.Atoi(splitConsoleAddress[2])
    if err != nil {
        // TODO: add error to diags
        return models.GcpCloudAccountResourceModel{}, diags
    }

    regions, diags := types.SetValueFrom(ctx, types.StringType, agentlessScanSpec.Regions)
    if diags.HasError() {
        return models.GcpCloudAccountResourceModel{}, diags
    }

    if (!planCloudAccount.AgentlessScanning.HubAccountId.IsNull() && !planCloudAccount.AgentlessScanning.HubAccountId.IsUnknown()) {
        hubAccountIdValue = types.StringValue(agentlessScanSpec.HubCredentialID)    
    } else {
        hubAccountIdValue = types.StringNull()
    }

    cloudAccount := models.GcpCloudAccountResourceModel{
        AccountName: types.StringValue(cloudScanRule.Credential.AccountName),
        Description: types.StringValue(cloudScanRule.Credential.Description),
        CredentialLevel: types.StringValue(credentialLevel),
        ServiceAccountKey: planCloudAccount.ServiceAccountKey, // Populate with planned value since API only returns encoded value
        ApiToken: planCloudAccount.ApiToken, // Populate with planned value since API only returns encoded value
        AgentlessScanning: models.AgentlessScanningResourceModel{
            Enabled: types.BoolValue(agentlessScanSpec.Enabled),
            IsHubAccount: types.BoolValue(agentlessScanSpec.HubAccount),
            HubAccountId: hubAccountIdValue,
            ConsoleURL: types.StringValue(consoleUrl),
            Port: types.Int32Value(int32(consolePort)),
            Subnet: types.StringValue(agentlessScanSpec.Subnet),
            ProxyAddress: types.StringValue(agentlessScanSpec.ProxyAddress),
            ProxyCertificate: types.StringValue(agentlessScanSpec.ProxyCA),
            Regions: regions,
            AutoScale: types.BoolValue(agentlessScanSpec.AutoScale),
            MaxNumberOfScanners: types.Int32Value(int32(agentlessScanSpec.Scanners)),
            ScanNonRunningHosts: types.BoolValue(agentlessScanSpec.ScanNonRunning),
            EnforcePermissionsCheck: types.BoolValue(!agentlessScanSpec.SkipPermissionsCheck),
        },
        ServerlessScanning: models.ServerlessScanningResourceModel{},
        EnableCloudDiscovery: types.BoolValue(cloudScanRule.DiscoveryEnabled),
    }

    return cloudAccount, diags
}
