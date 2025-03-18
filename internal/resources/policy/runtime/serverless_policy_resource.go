package policy

import (
	"context"
	"fmt"
    "time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	ruleAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/rule"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &ServerlessRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &ServerlessRuntimePolicyResource{}

func (r *ServerlessRuntimePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_serverless_runtime_policy"
}

func (r *ServerlessRuntimePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema(ctx)
}

func (r *ServerlessRuntimePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServerlessRuntimePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan models.RuntimeServerlessPolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    data, diags := ServerlessRuntimePolicySchemaToTerraform(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new serverless runtime policy 
    err := policyAPI.UpsertRuntimeServerless(*r.client, data)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Serverless Runtime Policy resource", 
            "Failed to create serverless runtime policy: " + err.Error(),
        )
        return
	}

    // Retrieve newly created serverless runtime policy 
    response, err := policyAPI.GetRuntimeServerless(*r.client)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created Serverless Runtime Policy resource", 
            "Failed to retrieve created serverless runtime policy: " + err.Error(),
        )
        return
    }

    createdPolicy, diags := ServerlessRuntimePolicyTerraformToSchema(ctx, response, plan, r.client)
    if diags.HasError() {
        return
    }

    // Set state to collection data
    diags = resp.State.Set(ctx, createdPolicy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *ServerlessRuntimePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state models.RuntimeServerlessPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get policy value from Prisma Cloud
    data, err := policyAPI.GetRuntimeServerless(*r.client)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error reading Serverless Runtime Policy resource", 
            "Failed to read serverless runtime policy: " + err.Error(),
        )
        return
    }

    // Overwrite state values with Prisma Cloud data
    policySchema, diags := ServerlessRuntimePolicyTerraformToSchema(ctx, data, state, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set refreshed state
    diags = resp.State.Set(ctx, &policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *ServerlessRuntimePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state models.RuntimeServerlessPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan models.RuntimeServerlessPolicyResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    planPolicy, diags := ServerlessRuntimePolicySchemaToTerraform(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Update existing policy
    err := policyAPI.UpsertRuntimeServerless(*r.client, planPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating Serverless Runtime Policy resource", 
            "Failed to update serverless runtime policy: " + err.Error(),
        )
        return
	}

    // Get updated policy value from Prisma Cloud
    updatedPolicy, err := policyAPI.GetRuntimeServerless(*r.client)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Serverless Runtime Policy resource", 
            "Failed to read Serverless Runtime Policy: " + err.Error(),
        )
        return
    }

    // Convert updated policy into schema
    policySchema, diags := ServerlessRuntimePolicyTerraformToSchema(ctx, updatedPolicy, plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set updated state
    diags = resp.State.Set(ctx, policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *ServerlessRuntimePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state models.RuntimeServerlessPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Clear policy rules
    state.Rules = &[]models.RuntimeServerlessPolicyRuleResourceModel{}

    // Generate API request body from plan
    updatedPlan, diags := ServerlessRuntimePolicySchemaToTerraform(ctx, &state, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing policy 
    err := policyAPI.UpsertRuntimeServerless(*r.client, updatedPlan)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Serverless Runtime Policy resource", 
            "Failed to delete serverless runtime policy: " + err.Error(),
        )
        return
	}
}

func (r *ServerlessRuntimePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ServerlessRuntimePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyPlan")

    //policy.ModifyPolicyResourcePlan(ctx, r.client, req.Plan, resp)

    util.DLog(ctx, "exiting ModifyPlan")
}

func ServerlessRuntimePolicySchemaToTerraform(ctx context.Context, plan *models.RuntimeServerlessPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (policyAPI.RuntimeServerlessPolicy, diag.Diagnostics) {
    util.DLog(ctx, "Executing ServerlessRuntimePolicySchemaToTerraform")

    var (
        diags diag.Diagnostics
        rules []policyAPI.RuntimeServerlessPolicyRule
    )

    if plan.Rules != nil {
        rules, diags = ServerlessRuntimePolicyRulesSchemaToTerraform(ctx, *plan.Rules, client)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessPolicy{}, diags
        }
    } else {
        rules = []policyAPI.RuntimeServerlessPolicyRule{}
    }

    tfPolicy := policyAPI.RuntimeServerlessPolicy{
        Id:         "serverlessRuntime",
        Rules:      &rules,
        LearningDisabled: false, 
    }

    tfPolicy.SortRules(ctx, plan.Rules)

    util.DLog(ctx, "Finishing ServerlessRuntimePolicySchemaToTerraform execution")

    return tfPolicy, diags
}

func ServerlessRuntimePolicyRulesSchemaToTerraform(ctx context.Context, schemaRules []models.RuntimeServerlessPolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient) ([]policyAPI.RuntimeServerlessPolicyRule, diag.Diagnostics) {
    util.DLog(ctx, "Executing ServerlessRuntimePolicyRulesSchemaToTerraform")

    var (
        diags diag.Diagnostics
        collectionNames []string
    )
    
    rules := []policyAPI.RuntimeServerlessPolicyRule{}

    for _, schemaRule := range schemaRules {
        collectionNames = []string{}
        diags = schemaRule.Collections.ElementsAs(ctx, &collectionNames, false)
        if diags.HasError() {
            return rules, diags
        }

        collections, err := collectionAPI.GetCollections(*client, collectionNames)
        if err != nil {
            diags.AddError(
                "Value Conversion Error",
                fmt.Sprintf("Error retrieving collection names while converting compliance policy rules to schema: %s", err.Error()),
            )
            return rules, diags
        }

        network, dns, diags := serverlessRuntimeNetworkingToTerraform(ctx, schemaRule.Networking)
        if diags.HasError() {
            return rules, diags
        }

        fileSystem, diags := serverlessRuntimeFileSystemToTerraform(ctx, schemaRule.FileSystem)
        if diags.HasError() {
            return rules, diags
        }

        processes, diags := serverlessRuntimeProcessesToTerraform(ctx, schemaRule.Processes)
        if diags.HasError() {
            return rules, diags
        }

        rule := policyAPI.RuntimeServerlessPolicyRule{
            AdvancedProtection: schemaRule.AdvancedThreatProtection.ValueBool(),
            Collections: collections,
            Network: network,
            Dns: dns,
            Filesystem: fileSystem,
            Processes: processes,
            Disabled: schemaRule.Disabled.ValueBool(),
            Modified: time.Now().Format("2006-01-02T15:04:05.000Z"),
            Name: schemaRule.Name.ValueString(),
            Notes: schemaRule.Notes.ValueString(),
            Owner: schemaRule.Owner.ValueString(),
            PreviousName: schemaRule.PreviousName.ValueString(),
        }

        rules = append(rules, rule)
    }

    util.DLog(ctx, "Finishing ServerlessRuntimePolicyRulesSchemaToTerraform execution")

    return rules, diags
}

func ServerlessRuntimePolicyTerraformToSchema(ctx context.Context, policy policyAPI.RuntimeServerlessPolicy, plan models.RuntimeServerlessPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (models.RuntimeServerlessPolicyResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "Executing ServerlessRuntimePolicyTerraformToSchema")

    var (
        diags diag.Diagnostics
        rules []models.RuntimeServerlessPolicyRuleResourceModel
    )
    
    if policy.Rules != nil {
        customRuleIdMap, err := ruleAPI.GetCustomRuleNameToIdMappings(*client) 
        if err != nil {
            diags.AddError(
                "API Error",
                "Error during retrieval of custom rules data: " + err.Error(),
            )
            return models.RuntimeServerlessPolicyResourceModel{}, diags
        }

        rules, diags = ServerlessRuntimePolicyRulesTerraformToSchema(ctx, *policy.Rules, plan.Rules, customRuleIdMap)
        if diags.HasError() {
            return models.RuntimeServerlessPolicyResourceModel{}, diags
        }
    } else {
        rules = []models.RuntimeServerlessPolicyRuleResourceModel{}
    }

    schema := models.RuntimeServerlessPolicyResourceModel{
        Rules:         &rules,
    }

    schema.SortRules(ctx, plan.Rules)

    util.DLog(ctx, "Finishing RuntimeServerlessPolicyTerraformToSchema execution")

    return schema, diags
}

func ServerlessRuntimePolicyRulesTerraformToSchema(ctx context.Context, rules []policyAPI.RuntimeServerlessPolicyRule, planRules *[]models.RuntimeServerlessPolicyRuleResourceModel, customRuleIdMap map[int]string) ([]models.RuntimeServerlessPolicyRuleResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "Executing RuntimeServerlessPolicyRulesTerraformToSchema")

    var diags diag.Diagnostics

    schemaRules := []models.RuntimeServerlessPolicyRuleResourceModel{}

    if len(rules) == 0 {
        return schemaRules, diags
    }

    for _, rule := range rules {
        var (
            planRule models.RuntimeServerlessPolicyRuleResourceModel
            processes *models.RuntimeServerlessPolicyProcessesResourceModel
            networking *models.RuntimeServerlessPolicyNetworkingResourceModel
            fileSystem *models.RuntimeServerlessPolicyFileSystemResourceModel
            notes basetypes.StringValue
        )

        // Find the matching plan rule
        for idx, pRule := range *planRules {
            if pRule.Name.ValueString() == rule.Name {
                planRule = (*planRules)[idx]
                break
            }
        }

        collectionNames := []string{}
        for _, collection := range rule.Collections {
            collectionNames = append(collectionNames, collection.Name)
        }

        collections, diags := types.SetValueFrom(ctx, types.StringType, collectionNames)
        if diags.HasError() {
            return []models.RuntimeServerlessPolicyRuleResourceModel{}, diags
        }

        if planRule.Processes == nil {
            processes = nil
        } else {
            processesValue, diags := serverlessRuntimeProcessesToSchema(ctx, rule.Processes, planRule.Processes)
            if diags.HasError() {
                return []models.RuntimeServerlessPolicyRuleResourceModel{}, diags
            }

            processes = &processesValue
        }
    
        if planRule.Networking == nil {
            networking = nil
        } else {
            networkingValue, diags := serverlessRuntimeNetworkingToSchema(ctx, rule.Network, rule.Dns, planRule.Networking)
            if diags.HasError() {
                return []models.RuntimeServerlessPolicyRuleResourceModel{}, diags
            }

            networking = &networkingValue
        }

        if planRule.FileSystem == nil {
            fileSystem = nil
        } else {
            fileSystemValue, diags := serverlessRuntimeFileSystemToSchema(ctx, rule.Filesystem, planRule.FileSystem)
            if diags.HasError() {
                return []models.RuntimeServerlessPolicyRuleResourceModel{}, diags
            }

            fileSystem = &fileSystemValue
        }

        if planRule.Notes.IsNull() {
            notes = types.StringNull()
        } else {
            notes = types.StringValue(rule.Notes)
        }

        schemaRule := models.RuntimeServerlessPolicyRuleResourceModel{
            Processes: processes,
            Networking: networking,
            FileSystem: fileSystem,
            Collections: collections,
            AdvancedThreatProtection: types.BoolValue(rule.AdvancedProtection),
            Disabled: types.BoolValue(rule.Disabled),
            Modified: types.StringValue(""),
            Name: types.StringValue(rule.Name),
            Notes: notes,
            Order: planRule.Order,
            Owner: types.StringValue(rule.Owner),
            PreviousName: types.StringValue(rule.PreviousName),
        }

        schemaRules = append(schemaRules, schemaRule)
    }

    util.DLog(ctx, "Finishing RuntimeServerlessPolicyRulesTerraformToSchema exection")

    return schemaRules, diags
}

func serverlessRuntimeProcessesToTerraform(ctx context.Context, schemaProcesses *models.RuntimeServerlessPolicyProcessesResourceModel) (policyAPI.RuntimeServerlessProcesses, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        allowedProcesses []string
        //deniedProcessesPaths []string
    )

    if schemaProcesses == nil {
        return policyAPI.RuntimeServerlessProcesses{}, diags
    }


    if schemaProcesses.AllowedProcesses.IsNull() {
        allowedProcesses = []string{}
    } else {
        diags = schemaProcesses.AllowedProcesses.ElementsAs(ctx, &allowedProcesses, false)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessProcesses{}, diags
        }
    }

    var effect string
    if !schemaProcesses.Enabled.ValueBool() {
        effect = "disable"
    } else {
        effect = schemaProcesses.DeniedProcessesEffect.ValueString()
    }

    return policyAPI.RuntimeServerlessProcesses{
        Effect: effect,
        Whitelist: allowedProcesses,
        CheckCryptoMiners: schemaProcesses.CryptoMiners.ValueBool(),
        BlockAllBinaries: schemaProcesses.BlockAllProcessesExceptMain.ValueBool(),
    }, diags
}

func serverlessRuntimeProcessesToSchema(ctx context.Context, tfProcesses policyAPI.RuntimeServerlessProcesses, planProcesses *models.RuntimeServerlessPolicyProcessesResourceModel) (models.RuntimeServerlessPolicyProcessesResourceModel, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
    )

    allowedProcesses, diags := types.ListValueFrom(ctx, types.StringType, tfProcesses.Whitelist)
    if diags.HasError() {
        return models.RuntimeServerlessPolicyProcessesResourceModel{}, diags
    }
    
    var enabled bool
    if tfProcesses.Effect == "disable" {
        enabled = false 
    } else {
        enabled = true 
    }

    return models.RuntimeServerlessPolicyProcessesResourceModel{
        Enabled: types.BoolValue(enabled),
        AllowedProcesses: allowedProcesses,
        DeniedProcessesEffect: planProcesses.DeniedProcessesEffect,
        CryptoMiners: types.BoolValue(tfProcesses.CheckCryptoMiners),
        BlockAllProcessesExceptMain: types.BoolValue(tfProcesses.BlockAllBinaries),
    }, diags
}

func serverlessRuntimeNetworkingToTerraform(ctx context.Context, schemaNetworking *models.RuntimeServerlessPolicyNetworkingResourceModel) (policyAPI.RuntimeServerlessNetwork, policyAPI.RuntimeServerlessDns, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        allowedListeningPorts []string
        allowedOutboundInternetPorts []string
        allowedOutboundIPs []string
        allowedDnsDomains []string 
    )

    if schemaNetworking == nil {
        return policyAPI.RuntimeServerlessNetwork{}, policyAPI.RuntimeServerlessDns{}, diags
    }

    if schemaNetworking.AllowedListeningPorts.IsNull() {
        allowedListeningPorts = []string{}
    } else {
        diags = schemaNetworking.AllowedListeningPorts.ElementsAs(ctx, &allowedListeningPorts, false)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessNetwork{}, policyAPI.RuntimeServerlessDns{}, diags
        }
    }

    allowedListeningPortRanges, diags := policy.PortRangesToTerraform(allowedListeningPorts)
    if diags.HasError() {
        return policyAPI.RuntimeServerlessNetwork{}, policyAPI.RuntimeServerlessDns{}, diags
    }

    if schemaNetworking.AllowedOutboundInternetPorts.IsNull() {
        allowedOutboundInternetPorts = []string{}
    } else {
        diags = schemaNetworking.AllowedOutboundInternetPorts.ElementsAs(ctx, &allowedOutboundInternetPorts, false)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessNetwork{}, policyAPI.RuntimeServerlessDns{}, diags
        }
    }

    allowedOutboundInternetPortRanges, diags := policy.PortRangesToTerraform(allowedOutboundInternetPorts)
    if diags.HasError() {
        return policyAPI.RuntimeServerlessNetwork{}, policyAPI.RuntimeServerlessDns{}, diags
    }

    if schemaNetworking.AllowedOutboundIPs.IsNull() {
        allowedOutboundIPs = []string{}
    } else {
        diags = schemaNetworking.AllowedOutboundIPs.ElementsAs(ctx, &allowedOutboundIPs, false)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessNetwork{}, policyAPI.RuntimeServerlessDns{}, diags
        }
    }

    if schemaNetworking.AllowedDnsDomains.IsNull() {
        allowedDnsDomains = []string{}
    } else {
        diags = schemaNetworking.AllowedDnsDomains.ElementsAs(ctx, &allowedDnsDomains, false)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessNetwork{}, policyAPI.RuntimeServerlessDns{}, diags
        }
    }

    var (
        networkEffect string
        dnsEffect string
    )
    
    if !schemaNetworking.IpConnectivityEnabled.ValueBool() {
        networkEffect = "disable"
    } else {
        networkEffect = schemaNetworking.DeniedIPsPortsEffect.ValueString()
    }

    if !schemaNetworking.DnsEnabled.ValueBool() {
        dnsEffect = "disable"
    } else {
        dnsEffect = schemaNetworking.DeniedDnsDomainsEffect.ValueString()
    }

    network := policyAPI.RuntimeServerlessNetwork{
        Effect: networkEffect,
        WhitelistListeningPorts: allowedListeningPortRanges,
        WhitelistOutboundPorts: allowedOutboundInternetPortRanges,
        WhitelistIps: allowedOutboundIPs,
    }

    dns := policyAPI.RuntimeServerlessDns{
        Effect: dnsEffect,
        Whitelist: allowedDnsDomains,
    }

    return network, dns, diags
}

func serverlessRuntimeNetworkingToSchema(ctx context.Context, tfNetwork policyAPI.RuntimeServerlessNetwork, tfDns policyAPI.RuntimeServerlessDns, planNetworking *models.RuntimeServerlessPolicyNetworkingResourceModel) (models.RuntimeServerlessPolicyNetworkingResourceModel, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
    )

    whitelistListeningPortsValues := policy.PortRangeToStringSlice(tfNetwork.WhitelistListeningPorts)
    whitelistListeningPorts, diags := types.ListValueFrom(ctx, types.StringType, whitelistListeningPortsValues)
    if diags.HasError() {
        return models.RuntimeServerlessPolicyNetworkingResourceModel{}, diags
    }

    whitelistOutboundPortsValues := policy.PortRangeToStringSlice(tfNetwork.WhitelistOutboundPorts)
    whitelistOutboundPorts, diags := types.ListValueFrom(ctx, types.StringType, whitelistOutboundPortsValues)
    if diags.HasError() {
        return models.RuntimeServerlessPolicyNetworkingResourceModel{}, diags
    }

    whitelistIps, diags := types.ListValueFrom(ctx, types.StringType, tfNetwork.WhitelistIps)
    if diags.HasError() {
        return models.RuntimeServerlessPolicyNetworkingResourceModel{}, diags
    }

    allowedDnsDomains, diags := types.ListValueFrom(ctx, types.StringType, tfDns.Whitelist)
    if diags.HasError() {
        return models.RuntimeServerlessPolicyNetworkingResourceModel{}, diags
    }

    var (
        networkEnabled bool 
        dnsEnabled bool 
    )
   
    if tfNetwork.Effect == "disable" {
        networkEnabled = false
    } else {
        networkEnabled = true
    }

    if tfDns.Effect == "disable" {
        dnsEnabled = false 
    } else {
        dnsEnabled = true
    }

    return models.RuntimeServerlessPolicyNetworkingResourceModel{
        IpConnectivityEnabled: types.BoolValue(networkEnabled),
        AllowedListeningPorts: whitelistListeningPorts,
        AllowedOutboundInternetPorts: whitelistOutboundPorts,
        AllowedOutboundIPs: whitelistIps,
        DeniedIPsPortsEffect: planNetworking.DeniedIPsPortsEffect,
        DnsEnabled: types.BoolValue(dnsEnabled),
        AllowedDnsDomains: allowedDnsDomains,
        DeniedDnsDomainsEffect: planNetworking.DeniedDnsDomainsEffect,
    }, diags
}

func serverlessRuntimeFileSystemToTerraform(ctx context.Context, schemaFileSystem *models.RuntimeServerlessPolicyFileSystemResourceModel) (policyAPI.RuntimeServerlessFilesystem, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        allowedPaths []string
        deniedPaths []string
    )

    if schemaFileSystem == nil {
        return policyAPI.RuntimeServerlessFilesystem{}, diags
    }

    if schemaFileSystem.AllowedPaths.IsNull() {
        allowedPaths = []string{}
    } else {
        diags = schemaFileSystem.AllowedPaths.ElementsAs(ctx, &allowedPaths, false)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessFilesystem{}, diags
        }
    }

    if schemaFileSystem.DeniedPaths.IsNull() {
        deniedPaths = []string{}
    } else {
        diags = schemaFileSystem.DeniedPaths.ElementsAs(ctx, &deniedPaths, false)
        if diags.HasError() {
            return policyAPI.RuntimeServerlessFilesystem{}, diags
        }
    }

    var effect string
    if !schemaFileSystem.Enabled.ValueBool() {
        effect = "disable"
    } else {
        effect = schemaFileSystem.DeniedPathsEffect.ValueString()
    }

    return policyAPI.RuntimeServerlessFilesystem{
        Effect: effect,
        Whitelist: allowedPaths,
        Blacklist: deniedPaths,
    }, diags
}

func serverlessRuntimeFileSystemToSchema(ctx context.Context, tfFileSystem policyAPI.RuntimeServerlessFilesystem, planFileSystem *models.RuntimeServerlessPolicyFileSystemResourceModel) (models.RuntimeServerlessPolicyFileSystemResourceModel, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        enabled bool
    )

    allowedPaths, diags := types.ListValueFrom(ctx, types.StringType, tfFileSystem.Whitelist)
    if diags.HasError() {
        return models.RuntimeServerlessPolicyFileSystemResourceModel{}, diags
    }

    deniedPaths, diags := types.ListValueFrom(ctx, types.StringType, tfFileSystem.Blacklist)
    if diags.HasError() {
        return models.RuntimeServerlessPolicyFileSystemResourceModel{}, diags
    }

    if tfFileSystem.Effect == "disable" {
        enabled = false
    } else {
        enabled = true 
    }
    
    return models.RuntimeServerlessPolicyFileSystemResourceModel{
        Enabled: types.BoolValue(enabled),
        AllowedPaths: allowedPaths, 
        DeniedPaths: deniedPaths,
        DeniedPathsEffect: planFileSystem.DeniedPathsEffect,
    }, diags
}
