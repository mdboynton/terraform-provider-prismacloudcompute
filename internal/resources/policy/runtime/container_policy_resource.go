package policy

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	ruleAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/rule"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"

	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &ContainerRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &ContainerRuntimePolicyResource{}

func (r *ContainerRuntimePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_runtime_policy"
}

func (r *ContainerRuntimePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.GetSchema(ctx)
}

func (r *ContainerRuntimePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ContainerRuntimePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.RuntimeContainerPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	data, diags := ContainerRuntimePolicySchemaToTerraform(ctx, &plan, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new container runtime policy
	err := policyAPI.UpsertRuntimeContainerPolicyFiltered(*r.client, data, []string{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Container Runtime Policy resource",
			"Failed to create container runtime policy: "+err.Error(),
		)
		return
	}

	// Retrieve newly created container runtime policy
	response, err := policyAPI.GetRuntimeContainerPolicyFiltered(*r.client, plan.GetRuleNames())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving created Container Runtime Policy resource",
			"Failed to retrieve created container runtime policy: "+err.Error(),
		)
		return
	}

	createdPolicy, diags := ContainerRuntimePolicyTerraformToSchema(ctx, response, plan, r.client)
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

func (r *ContainerRuntimePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.RuntimeContainerPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get policy value from Prisma Cloud
	data, err := policyAPI.GetRuntimeContainerPolicyFiltered(*r.client, state.GetRuleNames())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Container Runtime Policy resource",
			"Failed to read container runtime policy: "+err.Error(),
		)
		return
	}

	// Overwrite state values with Prisma Cloud data
	policySchema, diags := ContainerRuntimePolicyTerraformToSchema(ctx, data, state, r.client)
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

func (r *ContainerRuntimePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state models.RuntimeContainerPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.RuntimeContainerPolicyResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	planPolicy, diags := ContainerRuntimePolicySchemaToTerraform(ctx, &plan, r.client)
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
	err := policyAPI.UpsertRuntimeContainerPolicyFiltered(*r.client, planPolicy, deletedRuleNames)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Container Runtime Policy resource",
			"Failed to update container runtime policy: "+err.Error(),
		)
		return
	}

	// Get updated policy value from Prisma Cloud
	updatedPolicy, err := policyAPI.GetRuntimeContainerPolicyFiltered(*r.client, plan.GetRuleNames())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Container Runtime Policy resource",
			"Failed to read Container Runtime Policy: "+err.Error(),
		)
		return
	}

	// Convert updated policy into schema
	policySchema, diags := ContainerRuntimePolicyTerraformToSchema(ctx, updatedPolicy, plan, r.client)
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

func (r *ContainerRuntimePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.RuntimeContainerPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleNames := state.GetRuleNames()

	// Clear policy rules
	state.Rules = &[]models.RuntimeContainerPolicyRuleResourceModel{}

	// Generate API request body from plan
	updatedPlan, diags := ContainerRuntimePolicySchemaToTerraform(ctx, &state, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing policy
	err := policyAPI.UpsertRuntimeContainerPolicyFiltered(*r.client, updatedPlan, ruleNames)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Container Runtime Policy resource",
			"Failed to delete container runtime policy: "+err.Error(),
		)
		return
	}
}

func (r *ContainerRuntimePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ContainerRuntimePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	util.HCLogDebug(ctx, "entering ModifyPlan")

	//policy.ModifyPolicyResourcePlan(ctx, r.client, req.Plan, resp)

	util.HCLogDebug(ctx, "exiting ModifyPlan")
}

func ContainerRuntimePolicySchemaToTerraform(ctx context.Context, plan *models.RuntimeContainerPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (policyAPI.RuntimeContainerPolicy, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing ContainerRuntimePolicySchemaToTerraform")

	var (
		diags diag.Diagnostics
		rules []policyAPI.RuntimeContainerPolicyRule
	)

	if plan.Rules != nil {
		customRuleIdMap, err := ruleAPI.GetCustomRuleIdToNameMappings(*client)
		if err != nil {
			diags.AddError(
				"API Error",
				"Error during retrieval of custom rules data: "+err.Error(),
			)
			return policyAPI.RuntimeContainerPolicy{}, diags
		}

		rules, diags = ContainerRuntimePolicyRulesSchemaToTerraform(ctx, *plan.Rules, client, customRuleIdMap)
		if diags.HasError() {
			return policyAPI.RuntimeContainerPolicy{}, diags
		}
	} else {
		rules = []policyAPI.RuntimeContainerPolicyRule{}
	}

	tfPolicy := policyAPI.RuntimeContainerPolicy{
		Id:               "containerRuntime",
		Rules:            &rules,
		LearningDisabled: !plan.AutomaticRuntimeLearning.ValueBool(),
	}

	//tfPolicy.SortRules(ctx, plan.Rules)

	util.HCLogDebug(ctx, "Finishing ContainerRuntimePolicySchemaToTerraform execution")

	return tfPolicy, diags
}

func ContainerRuntimePolicyRulesSchemaToTerraform(ctx context.Context, schemaRules []models.RuntimeContainerPolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient, customRuleIdMap map[string]int) ([]policyAPI.RuntimeContainerPolicyRule, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing ContainerRuntimePolicyRulesSchemaToTerraform")

	var (
		diags           diag.Diagnostics
		collectionNames []string
	)

	rules := []policyAPI.RuntimeContainerPolicyRule{}

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

		network, dns, diags := containerRuntimeNetworkingToTerraform(ctx, schemaRule.Networking)
		if diags.HasError() {
			return rules, diags
		}

		fileSystem, diags := containerRuntimeFileSystemToTerraform(ctx, schemaRule.FileSystem)
		if diags.HasError() {
			return rules, diags
		}

		processes, allowAllActivityInAttachedSessions, diags := containerRuntimeProcessesToTerraform(ctx, schemaRule.Processes)
		if diags.HasError() {
			return rules, diags
		}

		customRules, diags := models.CustomRuntimeRulesToTerraform(ctx, schemaRule.CustomRules, customRuleIdMap)
		if diags.HasError() {
			return rules, diags
		}

		rule := policyAPI.RuntimeContainerPolicyRule{
			AdvancedProtectionEffect:       schemaRule.AntiMalware.MalwareFromAdvancedThreatProtection.ValueString(),
			KubernetesEnforcementEffect:    schemaRule.AntiMalware.KubernetesAttacks.ValueString(),
			CloudMetadataEnforcementEffect: schemaRule.AntiMalware.SuspiciousCloudProviderApiQueries.ValueString(),
			SkipExecSessions:               !allowAllActivityInAttachedSessions,
			WildFireAnalysis:               schemaRule.AntiMalware.WildFireAnalysis.ValueString(),
			Collections:                    collections,
			Network:                        network,
			Dns:                            dns,
			Filesystem:                     fileSystem,
			Processes:                      processes,
			CustomRules:                    customRules,
			Disabled:                       schemaRule.Disabled.ValueBool(),
			Modified:                       time.Now().Format("2006-01-02T15:04:05.000Z"),
			Name:                           schemaRule.Name.ValueString(),
			Notes:                          schemaRule.Notes.ValueString(),
			Owner:                          schemaRule.Owner.ValueString(),
			PreviousName:                   schemaRule.PreviousName.ValueString(),
		}

		rules = append(rules, rule)
	}

	util.HCLogDebug(ctx, "Finishing ContainerRuntimePolicyRulesSchemaToTerraform execution")

	return rules, diags
}

func ContainerRuntimePolicyTerraformToSchema(ctx context.Context, policy policyAPI.RuntimeContainerPolicy, plan models.RuntimeContainerPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (models.RuntimeContainerPolicyResourceModel, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing ContainerRuntimePolicyTerraformToSchema")

	var (
		diags diag.Diagnostics
		rules []models.RuntimeContainerPolicyRuleResourceModel
	)

	if policy.Rules != nil {
		customRuleIdMap, err := ruleAPI.GetCustomRuleNameToIdMappings(*client)
		if err != nil {
			diags.AddError(
				"API Error",
				"Error during retrieval of custom rules data: "+err.Error(),
			)
			return models.RuntimeContainerPolicyResourceModel{}, diags
		}

		rules, diags = ContainerRuntimePolicyRulesTerraformToSchema(ctx, *policy.Rules, plan.Rules, customRuleIdMap)
		if diags.HasError() {
			return models.RuntimeContainerPolicyResourceModel{}, diags
		}
	} else {
		rules = []models.RuntimeContainerPolicyRuleResourceModel{}
	}

	schema := models.RuntimeContainerPolicyResourceModel{
		Rules:                    &rules,
		AutomaticRuntimeLearning: types.BoolValue(!policy.LearningDisabled),
	}

	schema.SortRules(ctx, plan.Rules)

	util.HCLogDebug(ctx, "Finishing RuntimeContainerPolicyTerraformToSchema execution")

	return schema, diags
}

func ContainerRuntimePolicyRulesTerraformToSchema(ctx context.Context, rules []policyAPI.RuntimeContainerPolicyRule, planRules *[]models.RuntimeContainerPolicyRuleResourceModel, customRuleIdMap map[int]string) ([]models.RuntimeContainerPolicyRuleResourceModel, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing RuntimeContainerPolicyRulesTerraformToSchema")

	var diags diag.Diagnostics

	schemaRules := []models.RuntimeContainerPolicyRuleResourceModel{}

	if len(rules) == 0 {
		return schemaRules, diags
	}

	for _, rule := range rules {
		var (
			planRule    models.RuntimeContainerPolicyRuleResourceModel
			processes   *models.RuntimeContainerPolicyProcessesResourceModel
			networking  *models.RuntimeContainerPolicyNetworkingResourceModel
			fileSystem  *models.RuntimeContainerPolicyFileSystemResourceModel
			customRules *[]models.RuntimePolicyCustomRuleResourceModel
			notes       basetypes.StringValue
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
			return []models.RuntimeContainerPolicyRuleResourceModel{}, diags
		}

		if planRule.Processes == nil {
			processes = nil
		} else {
			processesValue, diags := containerRuntimeProcessesToSchema(ctx, rule.Processes, rule.SkipExecSessions)
			if diags.HasError() {
				return []models.RuntimeContainerPolicyRuleResourceModel{}, diags
			}

			processes = &processesValue
		}

		if planRule.Networking == nil {
			networking = nil
		} else {
			networkingValue, diags := containerRuntimeNetworkingToSchema(ctx, rule.Network, rule.Dns)
			if diags.HasError() {
				return []models.RuntimeContainerPolicyRuleResourceModel{}, diags
			}

			networking = &networkingValue
		}

		if planRule.FileSystem == nil {
			fileSystem = nil
		} else {
			fileSystemValue, diags := containerRuntimeFileSystemToSchema(ctx, rule.Filesystem)
			if diags.HasError() {
				return []models.RuntimeContainerPolicyRuleResourceModel{}, diags
			}

			fileSystem = &fileSystemValue
		}

		if planRule.CustomRules == nil {
			customRules = nil
		} else {
			customRulesValue, diags := models.CustomRuntimeRulesToSchema(ctx, rule.CustomRules, customRuleIdMap)
			if diags.HasError() {
				return []models.RuntimeContainerPolicyRuleResourceModel{}, diags
			}

			customRules = &customRulesValue
		}

		if planRule.Notes.IsNull() {
			notes = types.StringNull()
		} else {
			notes = types.StringValue(rule.Notes)
		}

		schemaRule := models.RuntimeContainerPolicyRuleResourceModel{
			AntiMalware: &models.RuntimeContainerPolicyAntiMalwareResourceModel{
				MalwareFromAdvancedThreatProtection: types.StringValue(rule.AdvancedProtectionEffect),
				KubernetesAttacks:                   types.StringValue(rule.KubernetesEnforcementEffect),
				SuspiciousCloudProviderApiQueries:   types.StringValue(rule.CloudMetadataEnforcementEffect),
				WildFireAnalysis:                    types.StringValue(rule.WildFireAnalysis),
			},
			Processes:    processes,
			Networking:   networking,
			FileSystem:   fileSystem,
			CustomRules:  customRules,
			Collections:  collections,
			Disabled:     types.BoolValue(rule.Disabled),
			Modified:     types.StringValue(""),
			Name:         types.StringValue(rule.Name),
			Notes:        notes,
			Order:        planRule.Order,
			Owner:        types.StringValue(rule.Owner),
			PreviousName: types.StringValue(rule.PreviousName),
		}

		schemaRules = append(schemaRules, schemaRule)
	}

	util.HCLogDebug(ctx, "Finishing RuntimeContainerPolicyRulesTerraformToSchema exection")

	return schemaRules, diags
}

func containerRuntimeProcessesToTerraform(ctx context.Context, schemaProcesses *models.RuntimeContainerPolicyProcessesResourceModel) (policyAPI.RuntimeContainerProcesses, bool, diag.Diagnostics) {
	var (
		diags                diag.Diagnostics
		allowedProcesses     []string
		deniedProcessesPaths []string
	)

	if schemaProcesses == nil {
		return policyAPI.RuntimeContainerProcesses{}, false, diags
	}

	if schemaProcesses.AllowedProcesses.IsNull() {
		allowedProcesses = []string{}
	} else {
		diags = schemaProcesses.AllowedProcesses.ElementsAs(ctx, &allowedProcesses, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerProcesses{}, false, diags
		}
	}

	if schemaProcesses.DeniedProcesses.Paths.IsNull() {
		deniedProcessesPaths = []string{}
	} else {
		diags = schemaProcesses.DeniedProcesses.Paths.ElementsAs(ctx, &deniedProcessesPaths, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerProcesses{}, false, diags
		}
	}

	deniedList := policyAPI.RuntimeContainerDeniedList{
		Paths:  deniedProcessesPaths,
		Effect: schemaProcesses.DeniedProcesses.Effect.ValueString(),
	}

	return policyAPI.RuntimeContainerProcesses{
		ModifiedProcessEffect: schemaProcesses.ProcessesFromModifiedBinaries.ValueString(),
		CryptoMinersEffect:    schemaProcesses.CryptoMiners.ValueString(),
		LateralMovementEffect: schemaProcesses.LateralMovementProcesses.ValueString(),
		ReverseShellEffect:    schemaProcesses.ReverseShell.ValueString(),
		SuidBinariesEffect:    schemaProcesses.ProcessesStartedWithSUID.ValueString(),
		DefaultEffect:         schemaProcesses.AllOtherProcessesEffect.ValueString(),
		CheckParentChild:      schemaProcesses.AllowOnlyLearnedProcessesFromParents.ValueBool(),
		AllowedList:           allowedProcesses,
		DeniedList:            deniedList,
		Disabled:              !schemaProcesses.Enabled.ValueBool(),
	}, schemaProcesses.AllowAllActivityInAttachedSessions.ValueBool(), diags
}

func containerRuntimeProcessesToSchema(ctx context.Context, tfProcesses policyAPI.RuntimeContainerProcesses, skipExecSessions bool) (models.RuntimeContainerPolicyProcessesResourceModel, diag.Diagnostics) {
	var (
		diags diag.Diagnostics
	)

	allowedProcesses, diags := types.ListValueFrom(ctx, types.StringType, tfProcesses.AllowedList)
	if diags.HasError() {
		return models.RuntimeContainerPolicyProcessesResourceModel{}, diags
	}

	deniedProcessesPaths, diags := types.ListValueFrom(ctx, types.StringType, tfProcesses.DeniedList.Paths)
	if diags.HasError() {
		return models.RuntimeContainerPolicyProcessesResourceModel{}, diags
	}

	deniedProcesses := models.RuntimePolicyDeniedProcessesResourceModel{
		Effect: types.StringValue(tfProcesses.DeniedList.Effect),
		Paths:  deniedProcessesPaths,
	}

	return models.RuntimeContainerPolicyProcessesResourceModel{
		Enabled:                              types.BoolValue(!tfProcesses.Disabled),
		AllowedProcesses:                     allowedProcesses,
		AllowOnlyLearnedProcessesFromParents: types.BoolValue(tfProcesses.CheckParentChild),
		AllowAllActivityInAttachedSessions:   types.BoolValue(!skipExecSessions),
		ProcessesFromModifiedBinaries:        types.StringValue(tfProcesses.ModifiedProcessEffect),
		CryptoMiners:                         types.StringValue(tfProcesses.CryptoMinersEffect),
		ReverseShell:                         types.StringValue(tfProcesses.ReverseShellEffect),
		LateralMovementProcesses:             types.StringValue(tfProcesses.LateralMovementEffect),
		ProcessesStartedWithSUID:             types.StringValue(tfProcesses.SuidBinariesEffect),
		DeniedProcesses:                      &deniedProcesses,
		AllOtherProcessesEffect:              types.StringValue(tfProcesses.DefaultEffect),
	}, diags
}

func containerRuntimeNetworkingToTerraform(ctx context.Context, schemaNetworking *models.RuntimeContainerPolicyNetworkingResourceModel) (policyAPI.RuntimeContainerNetwork, policyAPI.RuntimeContainerDns, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		allowedOutboundIPs []string
		deniedOutboundIPs  []string
		allowedDnsDomains  []string
		deniedDnsDomains   []string
	)

	if schemaNetworking == nil {
		return policyAPI.RuntimeContainerNetwork{}, policyAPI.RuntimeContainerDns{}, diags
	}

	// TODO: change it so that we only return on diags.HasError() at the end

	allowedListeningPorts, diags := models.PortRangesToTerraform(ctx, schemaNetworking.AllowedListeningPorts)
	allowedOutboundInternetPorts, diags := models.PortRangesToTerraform(ctx, schemaNetworking.AllowedOutboundInternetPorts)

	if schemaNetworking.AllowedOutboundIPs.IsNull() {
		allowedOutboundIPs = []string{}
	} else {
		diags = schemaNetworking.AllowedOutboundIPs.ElementsAs(ctx, &allowedOutboundIPs, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerNetwork{}, policyAPI.RuntimeContainerDns{}, diags
		}
	}

	deniedListeningPorts, diags := models.PortRangesToTerraform(ctx, schemaNetworking.DeniedListeningPorts)
	deniedOutboundInternetPorts, diags := models.PortRangesToTerraform(ctx, schemaNetworking.DeniedOutboundInternetPorts)

	if schemaNetworking.DeniedOutboundIPs.IsNull() {
		deniedOutboundIPs = []string{}
	} else {
		diags = schemaNetworking.DeniedOutboundIPs.ElementsAs(ctx, &deniedOutboundIPs, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerNetwork{}, policyAPI.RuntimeContainerDns{}, diags
		}
	}

	if schemaNetworking.AllowedDnsDomains.IsNull() {
		allowedDnsDomains = []string{}
	} else {
		diags = schemaNetworking.AllowedDnsDomains.ElementsAs(ctx, &allowedDnsDomains, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerNetwork{}, policyAPI.RuntimeContainerDns{}, diags
		}
	}

	if schemaNetworking.DeniedDnsDomains.IsNull() {
		deniedDnsDomains = []string{}
	} else {
		diags = schemaNetworking.DeniedDnsDomains.ElementsAs(ctx, &deniedDnsDomains, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerNetwork{}, policyAPI.RuntimeContainerDns{}, diags
		}
	}

	network := policyAPI.RuntimeContainerNetwork{
		AllowedIps:      allowedOutboundIPs,
		DefaultEffect:   schemaNetworking.AllOtherActivityEffect.ValueString(),
		DeniedIps:       deniedOutboundIPs,
		DeniedIpsEffect: schemaNetworking.DeniedOutboundIPsEffect.ValueString(),
		Disabled:        !schemaNetworking.IpConnectivityEnabled.ValueBool(),
		ListeningPorts: policyAPI.NetworkPorts{
			Allowed: allowedListeningPorts,
			Denied:  deniedListeningPorts,
			Effect:  schemaNetworking.DeniedListeningPortEffect.ValueString(),
		},
		ModifiedProcEffect: schemaNetworking.NetworkActivityFromModifiedBinaries.ValueString(),
		OutboundPorts: policyAPI.NetworkPorts{
			Allowed: allowedOutboundInternetPorts,
			Denied:  deniedOutboundInternetPorts,
			Effect:  schemaNetworking.DeniedOutboundInternetPortsEffect.ValueString(),
		},
		PortScanEffect:   schemaNetworking.PortScanning.ValueString(),
		RawSocketsEffect: schemaNetworking.RawSockets.ValueString(),
	}

	dns := policyAPI.RuntimeContainerDns{
		Disabled:      !schemaNetworking.DnsEnabled.ValueBool(),
		DefaultEffect: schemaNetworking.AllOtherDomainsEffect.ValueString(),
		DomainList: policyAPI.DnsDomainList{
			Allowed: allowedDnsDomains,
			Denied:  deniedDnsDomains,
			Effect:  schemaNetworking.DeniedDnsDomainsEffect.ValueString(),
		},
	}

	return network, dns, diags
}

func containerRuntimeNetworkingToSchema(ctx context.Context, tfNetwork policyAPI.RuntimeContainerNetwork, tfDns policyAPI.RuntimeContainerDns) (models.RuntimeContainerPolicyNetworkingResourceModel, diag.Diagnostics) {
	var (
		diags diag.Diagnostics
	)

	allowedListeningPortsValues := policy.PortRangeToStringSlice(tfNetwork.ListeningPorts.Allowed)
	allowedListeningPorts, diags := types.ListValueFrom(ctx, types.StringType, allowedListeningPortsValues)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	allowedOutboundInternetPortsValues := policy.PortRangeToStringSlice(tfNetwork.OutboundPorts.Allowed)
	allowedOutboundInternetPorts, diags := types.ListValueFrom(ctx, types.StringType, allowedOutboundInternetPortsValues)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	allowedOutboundIPs, diags := types.ListValueFrom(ctx, types.StringType, tfNetwork.AllowedIps)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	deniedListeningPortsValues := policy.PortRangeToStringSlice(tfNetwork.ListeningPorts.Denied)
	deniedListeningPorts, diags := types.ListValueFrom(ctx, types.StringType, deniedListeningPortsValues)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	deniedOutboundInternetPortsValues := policy.PortRangeToStringSlice(tfNetwork.OutboundPorts.Denied)
	deniedOutboundInternetPorts, diags := types.ListValueFrom(ctx, types.StringType, deniedOutboundInternetPortsValues)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	deniedOutboundIPs, diags := types.ListValueFrom(ctx, types.StringType, tfNetwork.DeniedIps)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	allowedDnsDomains, diags := types.ListValueFrom(ctx, types.StringType, tfDns.DomainList.Allowed)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	deniedDnsDomains, diags := types.ListValueFrom(ctx, types.StringType, tfDns.DomainList.Denied)
	if diags.HasError() {
		return models.RuntimeContainerPolicyNetworkingResourceModel{}, diags
	}

	return models.RuntimeContainerPolicyNetworkingResourceModel{
		IpConnectivityEnabled:               types.BoolValue(!tfNetwork.Disabled),
		AllowedListeningPorts:               allowedListeningPorts,
		AllowedOutboundInternetPorts:        allowedOutboundInternetPorts,
		AllowedOutboundIPs:                  allowedOutboundIPs,
		NetworkActivityFromModifiedBinaries: types.StringValue(tfNetwork.ModifiedProcEffect),
		PortScanning:                        types.StringValue(tfNetwork.PortScanEffect),
		RawSockets:                          types.StringValue(tfNetwork.RawSocketsEffect),
		DeniedListeningPorts:                deniedListeningPorts,
		DeniedListeningPortEffect:           types.StringValue(tfNetwork.ListeningPorts.Effect),
		DeniedOutboundInternetPorts:         deniedOutboundInternetPorts,
		DeniedOutboundInternetPortsEffect:   types.StringValue(tfNetwork.OutboundPorts.Effect),
		DeniedOutboundIPs:                   deniedOutboundIPs,
		DeniedOutboundIPsEffect:             types.StringValue(tfNetwork.DeniedIpsEffect),
		AllOtherActivityEffect:              types.StringValue(tfNetwork.DefaultEffect),
		DnsEnabled:                          types.BoolValue(!tfDns.Disabled),
		AllowedDnsDomains:                   allowedDnsDomains,
		DeniedDnsDomains:                    deniedDnsDomains,
		DeniedDnsDomainsEffect:              types.StringValue(tfDns.DomainList.Effect),
		AllOtherDomainsEffect:               types.StringValue(tfDns.DefaultEffect),
	}, diags
}

func containerRuntimeFileSystemToTerraform(ctx context.Context, schemaFileSystem *models.RuntimeContainerPolicyFileSystemResourceModel) (policyAPI.RuntimeContainerFilesystem, diag.Diagnostics) {
	var (
		diags        diag.Diagnostics
		allowedPaths []string
		deniedPaths  []string
	)

	if schemaFileSystem == nil {
		return policyAPI.RuntimeContainerFilesystem{}, diags
	}

	if schemaFileSystem.AllowedPaths.IsNull() {
		allowedPaths = []string{}
	} else {
		diags = schemaFileSystem.AllowedPaths.ElementsAs(ctx, &allowedPaths, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerFilesystem{}, diags
		}
	}

	if schemaFileSystem.DeniedPaths.IsNull() {
		deniedPaths = []string{}
	} else {
		diags = schemaFileSystem.DeniedPaths.ElementsAs(ctx, &deniedPaths, false)
		if diags.HasError() {
			return policyAPI.RuntimeContainerFilesystem{}, diags
		}
	}

	return policyAPI.RuntimeContainerFilesystem{
		AllowedList:         allowedPaths,
		BackdoorFilesEffect: schemaFileSystem.ChangesToBinaries.ValueString(),
		DefaultEffect:       schemaFileSystem.AllOtherPathsEffect.ValueString(),
		DeniedList: policyAPI.RuntimeContainerDeniedList{
			Effect: schemaFileSystem.DeniedPathsEffect.ValueString(),
			Paths:  deniedPaths,
		},
		Disabled:                   !schemaFileSystem.Enabled.ValueBool(),
		EncryptedBinariesEffect:    schemaFileSystem.DetectionOfEncryptedBinaries.ValueString(),
		NewFilesEffect:             schemaFileSystem.ChangesToSshAdminAccountConfigFiles.ValueString(),
		SuspiciousElfHeadersEffect: schemaFileSystem.SuspiciousELFHeaders.ValueString(),
	}, diags
}

func containerRuntimeFileSystemToSchema(ctx context.Context, tfFileSystem policyAPI.RuntimeContainerFilesystem) (models.RuntimeContainerPolicyFileSystemResourceModel, diag.Diagnostics) {
	var (
		diags diag.Diagnostics
	)

	allowedPaths, diags := types.ListValueFrom(ctx, types.StringType, tfFileSystem.AllowedList)
	if diags.HasError() {
		return models.RuntimeContainerPolicyFileSystemResourceModel{}, diags
	}

	deniedPaths, diags := types.ListValueFrom(ctx, types.StringType, tfFileSystem.DeniedList.Paths)
	if diags.HasError() {
		return models.RuntimeContainerPolicyFileSystemResourceModel{}, diags
	}

	return models.RuntimeContainerPolicyFileSystemResourceModel{
		Enabled:                             types.BoolValue(!tfFileSystem.Disabled),
		AllOtherPathsEffect:                 types.StringValue(tfFileSystem.DefaultEffect),
		AllowedPaths:                        allowedPaths,
		ChangesToBinaries:                   types.StringValue(tfFileSystem.BackdoorFilesEffect),
		ChangesToSshAdminAccountConfigFiles: types.StringValue(tfFileSystem.NewFilesEffect),
		DeniedPaths:                         deniedPaths,
		DeniedPathsEffect:                   types.StringValue(tfFileSystem.DeniedList.Effect),
		DetectionOfEncryptedBinaries:        types.StringValue(tfFileSystem.EncryptedBinariesEffect),
		SuspiciousELFHeaders:                types.StringValue(tfFileSystem.SuspiciousElfHeadersEffect),
	}, diags
}
