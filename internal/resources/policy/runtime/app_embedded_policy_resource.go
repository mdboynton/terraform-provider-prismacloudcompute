package policy

import (
	"context"
	"fmt"
    "time"
    "slices"

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

var _ resource.Resource = &AppEmbeddedRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &AppEmbeddedRuntimePolicyResource{}

func (r *AppEmbeddedRuntimePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_app_embedded_runtime_policy"
}

func (r *AppEmbeddedRuntimePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema(ctx)
}

func (r *AppEmbeddedRuntimePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AppEmbeddedRuntimePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan models.RuntimeAppEmbeddedPolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    data, diags := AppEmbeddedRuntimePolicySchemaToTerraform(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new app-embedded runtime policy 
    err := policyAPI.UpsertRuntimeAppEmbeddedPolicyFiltered(*r.client, data, []string{})
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating App-Embedded Runtime Policy resource", 
            "Failed to create app-embedded runtime policy: " + err.Error(),
        )
        return
	}

    // Retrieve newly created app-embedded runtime policy 
    response, err := policyAPI.GetRuntimeAppEmbeddedPolicyFiltered(*r.client, plan.GetRuleNames())
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created App-Embedded Runtime Policy resource", 
            "Failed to retrieve created app-embedded runtime policy: " + err.Error(),
        )
        return
    }

    createdPolicy, diags := AppEmbeddedRuntimePolicyTerraformToSchema(ctx, response, plan, r.client)
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

func (r *AppEmbeddedRuntimePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state models.RuntimeAppEmbeddedPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get policy value from Prisma Cloud
    data, err := policyAPI.GetRuntimeAppEmbeddedPolicyFiltered(*r.client, state.GetRuleNames())
    if err != nil {
		resp.Diagnostics.AddError(
            "Error reading App-Embedded Runtime Policy resource", 
            "Failed to read app-embedded runtime policy: " + err.Error(),
        )
        return
    }

    // Overwrite state values with Prisma Cloud data
    policySchema, diags := AppEmbeddedRuntimePolicyTerraformToSchema(ctx, data, state, r.client)
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

func (r *AppEmbeddedRuntimePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state models.RuntimeAppEmbeddedPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan models.RuntimeAppEmbeddedPolicyResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    planPolicy, diags := AppEmbeddedRuntimePolicySchemaToTerraform(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Find any rules being deleted
    deletedRuleNames := []string{}
    planRuleNames := plan.GetRuleNames()
    for _, stateRule := range *state.Rules {
        if !slices.Contains(planRuleNames, stateRule.Name.ValueString()) {
            deletedRuleNames = append(deletedRuleNames, stateRule.Name.ValueString())
        }
    }

    // Update existing policy
    err := policyAPI.UpsertRuntimeAppEmbeddedPolicyFiltered(*r.client, planPolicy, deletedRuleNames)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating App-Embedded Runtime Policy resource", 
            "Failed to update app-embedded runtime policy: " + err.Error(),
        )
        return
	}

    // Get updated policy value from Prisma Cloud
    updatedPolicy, err := policyAPI.GetRuntimeAppEmbeddedPolicyFiltered(*r.client, plan.GetRuleNames())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading App-Embedded Runtime Policy resource", 
            "Failed to read App-Embedded Runtime Policy: " + err.Error(),
        )
        return
    }

    // Convert updated policy into schema
    policySchema, diags := AppEmbeddedRuntimePolicyTerraformToSchema(ctx, updatedPolicy, plan, r.client)
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

func (r *AppEmbeddedRuntimePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state models.RuntimeAppEmbeddedPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    ruleNames := state.GetRuleNames()

    // Clear policy rules
    state.Rules = &[]models.RuntimeAppEmbeddedPolicyRuleResourceModel{}

    // Generate API request body from plan
    updatedPlan, diags := AppEmbeddedRuntimePolicySchemaToTerraform(ctx, &state, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing policy 
    err := policyAPI.UpsertRuntimeAppEmbeddedPolicyFiltered(*r.client, updatedPlan, ruleNames)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting App-Embedded Runtime Policy resource", 
            "Failed to delete app-embedded runtime policy: " + err.Error(),
        )
        return
	}
}

func (r *AppEmbeddedRuntimePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AppEmbeddedRuntimePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyPlan")

    //policy.ModifyPolicyResourcePlan(ctx, r.client, req.Plan, resp)

    util.DLog(ctx, "exiting ModifyPlan")
}

func AppEmbeddedRuntimePolicySchemaToTerraform(ctx context.Context, plan *models.RuntimeAppEmbeddedPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (policyAPI.RuntimeAppEmbeddedPolicy, diag.Diagnostics) {
    util.DLog(ctx, "Executing AppEmbeddedRuntimePolicySchemaToTerraform")

    var (
        diags diag.Diagnostics
        rules []policyAPI.RuntimeAppEmbeddedPolicyRule
    )

    if plan.Rules != nil {
        customRuleIdMap, err := ruleAPI.GetCustomRuleIdToNameMappings(*client) 
        if err != nil {
            diags.AddError(
                "API Error",
                "Error during retrieval of custom rules data: " + err.Error(),
            )
            return policyAPI.RuntimeAppEmbeddedPolicy{}, diags
        }

        rules, diags = AppEmbeddedRuntimePolicyRulesSchemaToTerraform(ctx, *plan.Rules, client, customRuleIdMap)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedPolicy{}, diags
        }
    } else {
        rules = []policyAPI.RuntimeAppEmbeddedPolicyRule{}
    }

    tfPolicy := policyAPI.RuntimeAppEmbeddedPolicy{
        Id:         "appEmbeddedRuntime",
        Rules:      &rules,
    }

    tfPolicy.SortRules(ctx, plan.Rules)

    util.DLog(ctx, "Finishing AppEmbeddedRuntimePolicySchemaToTerraform execution")

    return tfPolicy, diags
}

func AppEmbeddedRuntimePolicyRulesSchemaToTerraform(ctx context.Context, schemaRules []models.RuntimeAppEmbeddedPolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient, customRuleIdMap map[string]int) ([]policyAPI.RuntimeAppEmbeddedPolicyRule, diag.Diagnostics) {
    util.DLog(ctx, "Executing AppEmbeddedRuntimePolicyRulesSchemaToTerraform")

    var (
        diags diag.Diagnostics
        collectionNames []string
    )
    
    rules := []policyAPI.RuntimeAppEmbeddedPolicyRule{}

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

        network, dns, diags := appEmbeddedRuntimeNetworkingToTerraform(ctx, schemaRule.Networking)
        if diags.HasError() {
            return rules, diags
        }

        fileSystem, advancedProtection, wildFireAnalysis, diags := appEmbeddedRuntimeFileSystemToTerraform(ctx, schemaRule.FileSystem)
        if diags.HasError() {
            return rules, diags
        }

        processes, diags := appEmbeddedRuntimeProcessesToTerraform(ctx, schemaRule.Processes)
        if diags.HasError() {
            return rules, diags
        }

        customRules, diags := customRulesToTerraform(ctx, schemaRule.CustomRules, customRuleIdMap)
        if diags.HasError() {
            return rules, diags
        }

        rule := policyAPI.RuntimeAppEmbeddedPolicyRule{
            AdvancedProtection: advancedProtection,
            Collections: collections,
            CustomRules: customRules,
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
            WildFireAnalysis: wildFireAnalysis,
        }

        rules = append(rules, rule)
    }

    util.DLog(ctx, "Finishing AppEmbeddedRuntimePolicyRulesSchemaToTerraform execution")

    return rules, diags
}

func AppEmbeddedRuntimePolicyTerraformToSchema(ctx context.Context, policy policyAPI.RuntimeAppEmbeddedPolicy, plan models.RuntimeAppEmbeddedPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (models.RuntimeAppEmbeddedPolicyResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "Executing AppEmbeddedRuntimePolicyTerraformToSchema")

    var (
        diags diag.Diagnostics
        rules []models.RuntimeAppEmbeddedPolicyRuleResourceModel
    )
    
    if policy.Rules != nil {
        customRuleIdMap, err := ruleAPI.GetCustomRuleNameToIdMappings(*client) 
        if err != nil {
            diags.AddError(
                "API Error",
                "Error during retrieval of custom rules data: " + err.Error(),
            )
            return models.RuntimeAppEmbeddedPolicyResourceModel{}, diags
        }

        rules, diags = AppEmbeddedRuntimePolicyRulesTerraformToSchema(ctx, *policy.Rules, plan.Rules, customRuleIdMap)
        if diags.HasError() {
            return models.RuntimeAppEmbeddedPolicyResourceModel{}, diags
        }
    } else {
        rules = []models.RuntimeAppEmbeddedPolicyRuleResourceModel{}
    }

    schema := models.RuntimeAppEmbeddedPolicyResourceModel{
        Rules:         &rules,
    }

    schema.SortRules(ctx, plan.Rules)

    util.DLog(ctx, "Finishing RuntimeAppEmbeddedPolicyTerraformToSchema execution")

    return schema, diags
}

func AppEmbeddedRuntimePolicyRulesTerraformToSchema(ctx context.Context, rules []policyAPI.RuntimeAppEmbeddedPolicyRule, planRules *[]models.RuntimeAppEmbeddedPolicyRuleResourceModel, customRuleIdMap map[int]string) ([]models.RuntimeAppEmbeddedPolicyRuleResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "Executing RuntimeAppEmbeddedPolicyRulesTerraformToSchema")

    var diags diag.Diagnostics

    schemaRules := []models.RuntimeAppEmbeddedPolicyRuleResourceModel{}

    if len(rules) == 0 {
        return schemaRules, diags
    }

    for _, rule := range rules {
        var (
            planRule models.RuntimeAppEmbeddedPolicyRuleResourceModel
            processes *models.RuntimeAppEmbeddedPolicyProcessesResourceModel
            networking *models.RuntimeAppEmbeddedPolicyNetworkingResourceModel
            fileSystem *models.RuntimeAppEmbeddedPolicyFileSystemResourceModel
            customRules *[]models.RuntimeHostPolicyCustomRuleResourceModel
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
            return []models.RuntimeAppEmbeddedPolicyRuleResourceModel{}, diags
        }

        if planRule.Processes == nil {
            processes = nil
        } else {
            processesValue, diags := appEmbeddedRuntimeProcessesToSchema(ctx, rule.Processes, planRule.Processes)
            if diags.HasError() {
                return []models.RuntimeAppEmbeddedPolicyRuleResourceModel{}, diags
            }

            processes = &processesValue
        }
    
        if planRule.Networking == nil {
            networking = nil
        } else {
            networkingValue, diags := appEmbeddedRuntimeNetworkingToSchema(ctx, rule.Network, rule.Dns, planRule.Networking)
            if diags.HasError() {
                return []models.RuntimeAppEmbeddedPolicyRuleResourceModel{}, diags
            }

            networking = &networkingValue
        }

        if planRule.FileSystem == nil {
            fileSystem = nil
        } else {
            fileSystemValue, diags := appEmbeddedRuntimeFileSystemToSchema(ctx, rule.Filesystem, rule.AdvancedProtection, rule.WildFireAnalysis, planRule.FileSystem)
            if diags.HasError() {
                return []models.RuntimeAppEmbeddedPolicyRuleResourceModel{}, diags
            }

            fileSystem = &fileSystemValue
        }

        if planRule.CustomRules == nil {
            customRules = nil
        } else {
            customRulesValue, diags := customRulesToSchema(ctx, rule.CustomRules, customRuleIdMap)
            if diags.HasError() {
                return []models.RuntimeAppEmbeddedPolicyRuleResourceModel{}, diags
            }

            customRules = &customRulesValue
        }

        if planRule.Notes.IsNull() {
            notes = types.StringNull()
        } else {
            notes = types.StringValue(rule.Notes)
        }

        schemaRule := models.RuntimeAppEmbeddedPolicyRuleResourceModel{
            Processes: processes,
            Networking: networking,
            FileSystem: fileSystem,
            Collections: collections,
            CustomRules: customRules,
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

    util.DLog(ctx, "Finishing RuntimeAppEmbeddedPolicyRulesTerraformToSchema exection")

    return schemaRules, diags
}

func appEmbeddedRuntimeProcessesToTerraform(ctx context.Context, schemaProcesses *models.RuntimeAppEmbeddedPolicyProcessesResourceModel) (policyAPI.RuntimeAppEmbeddedProcesses, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        allowedProcesses []string
        deniedProcesses []string
    )

    if schemaProcesses == nil {
        return policyAPI.RuntimeAppEmbeddedProcesses{}, diags
    }

    if schemaProcesses.AllowedProcesses.IsNull() {
        allowedProcesses = []string{}
    } else {
        diags = schemaProcesses.AllowedProcesses.ElementsAs(ctx, &allowedProcesses, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedProcesses{}, diags
        }
    }

    if schemaProcesses.DeniedProcesses.IsNull() {
        deniedProcesses = []string{}
    } else {
        diags = schemaProcesses.DeniedProcesses.ElementsAs(ctx, &deniedProcesses, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedProcesses{}, diags
        }
    }

    var effect string
    if !schemaProcesses.Enabled.ValueBool() {
        effect = "disable"
    } else {
        effect = schemaProcesses.DeniedProcessesEffect.ValueString()
    }

    return policyAPI.RuntimeAppEmbeddedProcesses{
        Effect: effect,
        Whitelist: allowedProcesses,
        Blacklist: deniedProcesses,
        CheckCryptoMiners: schemaProcesses.CryptoMiners.ValueBool(),
        CheckNewBinaries: schemaProcesses.ProcessesFromModifiedBinaries.ValueBool(),
    }, diags
}

func appEmbeddedRuntimeProcessesToSchema(ctx context.Context, tfProcesses policyAPI.RuntimeAppEmbeddedProcesses, planProcesses *models.RuntimeAppEmbeddedPolicyProcessesResourceModel) (models.RuntimeAppEmbeddedPolicyProcessesResourceModel, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
    )

    allowedProcesses, diags := types.ListValueFrom(ctx, types.StringType, tfProcesses.Whitelist)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyProcessesResourceModel{}, diags
    }

    deniedProcesses, diags := types.ListValueFrom(ctx, types.StringType, tfProcesses.Blacklist)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyProcessesResourceModel{}, diags
    }
    
    var enabled bool
    if tfProcesses.Effect == "disable" {
        enabled = false 
    } else {
        enabled = true 
    }

    return models.RuntimeAppEmbeddedPolicyProcessesResourceModel{
        Enabled: types.BoolValue(enabled),
        AllowedProcesses: allowedProcesses,
        DeniedProcesses: deniedProcesses,
        DeniedProcessesEffect: planProcesses.DeniedProcessesEffect,
        CryptoMiners: types.BoolValue(tfProcesses.CheckCryptoMiners),
        ProcessesFromModifiedBinaries: types.BoolValue(tfProcesses.CheckNewBinaries),
    }, diags
}

func appEmbeddedRuntimeNetworkingToTerraform(ctx context.Context, schemaNetworking *models.RuntimeAppEmbeddedPolicyNetworkingResourceModel) (policyAPI.RuntimeAppEmbeddedNetwork, policyAPI.RuntimeAppEmbeddedDns, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        allowedListeningPorts []string
        allowedOutboundInternetPorts []string
        allowedOutboundIPs []string
        deniedListeningPorts []string
        deniedOutboundInternetPorts []string
        deniedOutboundIPs []string
        allowedDnsDomains []string 
    )

    if schemaNetworking == nil {
        return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
    }

    if schemaNetworking.AllowedListeningPorts.IsNull() {
        allowedListeningPorts = []string{}
    } else {
        diags = schemaNetworking.AllowedListeningPorts.ElementsAs(ctx, &allowedListeningPorts, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
        }
    }

    allowedListeningPortRanges, diags := policy.PortRangesToTerraform(allowedListeningPorts)
    if diags.HasError() {
        return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
    }

    if schemaNetworking.AllowedOutboundInternetPorts.IsNull() {
        allowedOutboundInternetPorts = []string{}
    } else {
        diags = schemaNetworking.AllowedOutboundInternetPorts.ElementsAs(ctx, &allowedOutboundInternetPorts, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
        }
    }

    allowedOutboundInternetPortRanges, diags := policy.PortRangesToTerraform(allowedOutboundInternetPorts)
    if diags.HasError() {
        return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
    }

    if schemaNetworking.AllowedOutboundIPs.IsNull() {
        allowedOutboundIPs = []string{}
    } else {
        diags = schemaNetworking.AllowedOutboundIPs.ElementsAs(ctx, &allowedOutboundIPs, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
        }
    }

    if schemaNetworking.DeniedListeningPorts.IsNull() {
        deniedListeningPorts = []string{}
    } else {
        diags = schemaNetworking.DeniedListeningPorts.ElementsAs(ctx, &deniedListeningPorts, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
        }
    }

    deniedListeningPortRanges, diags := policy.PortRangesToTerraform(deniedListeningPorts)
    if diags.HasError() {
        return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
    }

    if schemaNetworking.DeniedOutboundInternetPorts.IsNull() {
        deniedOutboundInternetPorts = []string{}
    } else {
        diags = schemaNetworking.DeniedOutboundInternetPorts.ElementsAs(ctx, &deniedOutboundInternetPorts, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
        }
    }

    deniedOutboundInternetPortRanges, diags := policy.PortRangesToTerraform(deniedOutboundInternetPorts)
    if diags.HasError() {
        return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
    }

    if schemaNetworking.DeniedOutboundIPs.IsNull() {
        deniedOutboundIPs = []string{}
    } else {
        diags = schemaNetworking.DeniedOutboundIPs.ElementsAs(ctx, &deniedOutboundIPs, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
        }
    }

    if schemaNetworking.AllowedDnsDomains.IsNull() {
        allowedDnsDomains = []string{}
    } else {
        diags = schemaNetworking.AllowedDnsDomains.ElementsAs(ctx, &allowedDnsDomains, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedNetwork{}, policyAPI.RuntimeAppEmbeddedDns{}, diags
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

    network := policyAPI.RuntimeAppEmbeddedNetwork{
        Effect: networkEffect,
        WhitelistListeningPorts: allowedListeningPortRanges,
        WhitelistOutboundPorts: allowedOutboundInternetPortRanges,
        WhitelistIPs: allowedOutboundIPs,
        BlacklistListeningPorts: deniedListeningPortRanges,
        BlacklistOutboundPorts: deniedOutboundInternetPortRanges,
        BlacklistIPs: deniedOutboundIPs,
    }

    dns := policyAPI.RuntimeAppEmbeddedDns{
        Effect: dnsEffect,
        Whitelist: allowedDnsDomains,
    }

    return network, dns, diags
}

func appEmbeddedRuntimeNetworkingToSchema(ctx context.Context, tfNetwork policyAPI.RuntimeAppEmbeddedNetwork, tfDns policyAPI.RuntimeAppEmbeddedDns, planNetworking *models.RuntimeAppEmbeddedPolicyNetworkingResourceModel) (models.RuntimeAppEmbeddedPolicyNetworkingResourceModel, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
    )

    whitelistListeningPortsValues := policy.PortRangeToStringSlice(tfNetwork.WhitelistListeningPorts)
    whitelistListeningPorts, diags := types.ListValueFrom(ctx, types.StringType, whitelistListeningPortsValues)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{}, diags
    }

    whitelistOutboundPortsValues := policy.PortRangeToStringSlice(tfNetwork.WhitelistOutboundPorts)
    whitelistOutboundPorts, diags := types.ListValueFrom(ctx, types.StringType, whitelistOutboundPortsValues)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{}, diags
    }

    whitelistIps, diags := types.ListValueFrom(ctx, types.StringType, tfNetwork.WhitelistIPs)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{}, diags
    }

    blacklistListeningPortsValues := policy.PortRangeToStringSlice(tfNetwork.BlacklistListeningPorts)
    blacklistListeningPorts, diags := types.ListValueFrom(ctx, types.StringType, blacklistListeningPortsValues)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{}, diags
    }

    blacklistOutboundPortsValues := policy.PortRangeToStringSlice(tfNetwork.BlacklistOutboundPorts)
    blacklistOutboundPorts, diags := types.ListValueFrom(ctx, types.StringType, blacklistOutboundPortsValues)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{}, diags
    }

    blacklistIps, diags := types.ListValueFrom(ctx, types.StringType, tfNetwork.BlacklistIPs)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{}, diags
    }

    allowedDnsDomains, diags := types.ListValueFrom(ctx, types.StringType, tfDns.Whitelist)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{}, diags
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

    return models.RuntimeAppEmbeddedPolicyNetworkingResourceModel{
        IpConnectivityEnabled: types.BoolValue(networkEnabled),
        AllowedListeningPorts: whitelistListeningPorts,
        AllowedOutboundInternetPorts: whitelistOutboundPorts,
        AllowedOutboundIPs: whitelistIps,
        DeniedIPsPortsEffect: planNetworking.DeniedIPsPortsEffect,
        DeniedListeningPorts: blacklistListeningPorts,
        DeniedOutboundInternetPorts: blacklistOutboundPorts,
        DeniedOutboundIPs: blacklistIps,
        DnsEnabled: types.BoolValue(dnsEnabled),
        AllowedDnsDomains: allowedDnsDomains,
        DeniedDnsDomainsEffect: planNetworking.DeniedDnsDomainsEffect,
    }, diags
}

func appEmbeddedRuntimeFileSystemToTerraform(ctx context.Context, schemaFileSystem *models.RuntimeAppEmbeddedPolicyFileSystemResourceModel) (policyAPI.RuntimeAppEmbeddedFilesystem, bool, string, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        allowedPaths []string
        deniedPaths []string
    )

    if schemaFileSystem == nil {
        return policyAPI.RuntimeAppEmbeddedFilesystem{}, false, "", diags
    }

    if schemaFileSystem.AllowedPaths.IsNull() {
        allowedPaths = []string{}
    } else {
        diags = schemaFileSystem.AllowedPaths.ElementsAs(ctx, &allowedPaths, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedFilesystem{}, false, "", diags
        }
    }

    if schemaFileSystem.DeniedPaths.IsNull() {
        deniedPaths = []string{}
    } else {
        diags = schemaFileSystem.DeniedPaths.ElementsAs(ctx, &deniedPaths, false)
        if diags.HasError() {
            return policyAPI.RuntimeAppEmbeddedFilesystem{}, false, "", diags
        }
    }

    var effect string
    if !schemaFileSystem.Enabled.ValueBool() {
        effect = "disable"
    } else {
        effect = schemaFileSystem.DeniedPathsEffect.ValueString()
    }

    advancedProtection := schemaFileSystem.MalwareFromCustomFeed.ValueBool()
    wildFireAnalysis := schemaFileSystem.WildFireAnalysis.ValueString()

    return policyAPI.RuntimeAppEmbeddedFilesystem{
        Effect: effect,
        Whitelist: allowedPaths,
        Blacklist: deniedPaths,
        BackdoorFiles: schemaFileSystem.ChangesToSshAdminAccountConfigFiles.ValueBool(),
        SkipEncryptedBinaries: !schemaFileSystem.DetectionOfEncryptedBinaries.ValueBool(),
        CheckNewFiles: schemaFileSystem.ChangesToBinariesAndCerts.ValueBool(),
        SuspiciousElfHeaders: schemaFileSystem.SuspiciousELFHeaders.ValueBool(),
    }, advancedProtection, wildFireAnalysis, diags
}

func appEmbeddedRuntimeFileSystemToSchema(ctx context.Context, tfFileSystem policyAPI.RuntimeAppEmbeddedFilesystem, advancedProtection bool, wildFireAnalysis string, planFileSystem *models.RuntimeAppEmbeddedPolicyFileSystemResourceModel) (models.RuntimeAppEmbeddedPolicyFileSystemResourceModel, diag.Diagnostics) {
    var (
        diags diag.Diagnostics
        enabled bool
    )

    allowedPaths, diags := types.ListValueFrom(ctx, types.StringType, tfFileSystem.Whitelist)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyFileSystemResourceModel{}, diags
    }

    deniedPaths, diags := types.ListValueFrom(ctx, types.StringType, tfFileSystem.Blacklist)
    if diags.HasError() {
        return models.RuntimeAppEmbeddedPolicyFileSystemResourceModel{}, diags
    }

    if tfFileSystem.Effect == "disable" {
        enabled = false
    } else {
        enabled = true 
    }
    
    return models.RuntimeAppEmbeddedPolicyFileSystemResourceModel{
        Enabled: types.BoolValue(enabled),
        AllowedPaths: allowedPaths, 
        DeniedPaths: deniedPaths,
        DeniedPathsEffect: planFileSystem.DeniedPathsEffect,
        ChangesToSshAdminAccountConfigFiles: types.BoolValue(tfFileSystem.BackdoorFiles),
        DetectionOfEncryptedBinaries: types.BoolValue(!tfFileSystem.SkipEncryptedBinaries),
        ChangesToBinariesAndCerts: types.BoolValue(tfFileSystem.CheckNewFiles),
        SuspiciousELFHeaders: types.BoolValue(tfFileSystem.SuspiciousElfHeaders),
        MalwareFromCustomFeed: types.BoolValue(advancedProtection),
        WildFireAnalysis: types.StringValue(wildFireAnalysis),
    }, diags
}
