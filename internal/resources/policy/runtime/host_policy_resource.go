package policy

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	ruleAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/rule"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &HostRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &HostRuntimePolicyResource{}

func (r *HostRuntimePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host_runtime_policy"
}

func (r *HostRuntimePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.GetSchema(ctx)
}

func (r *HostRuntimePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostRuntimePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.RuntimeHostPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	data, diags := RuntimePolicySchemaToTerraform(ctx, &plan, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new host runtime policy
	err := policyAPI.UpsertRuntimeHostPolicyFiltered(*r.client, data, []string{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Host Runtime Policy resource",
			"Failed to create host runtime policy: "+err.Error(),
		)
		return
	}

	// Retrieve newly created host runtime policy
	response, err := policyAPI.GetRuntimeHostPolicyFiltered(*r.client, plan.GetRuleNames())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving created Host Runtime Policy resource",
			"Failed to retrieve created host runtime policy: "+err.Error(),
		)
		return
	}

	createdPolicy, diags := RuntimePolicyTerraformToSchema(ctx, response, plan, r.client)
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

func (r *HostRuntimePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.RuntimeHostPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get policy value from Prisma Cloud
	data, err := policyAPI.GetRuntimeHostPolicyFiltered(*r.client, state.GetRuleNames())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Host Runtime Policy resource",
			"Failed to read host runtime policy: "+err.Error(),
		)
		return
	}

	// Overwrite state values with Prisma Cloud data
	policySchema, diags := RuntimePolicyTerraformToSchema(ctx, data, state, r.client)
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

func (r *HostRuntimePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state models.RuntimeHostPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.RuntimeHostPolicyResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	planPolicy, diags := RuntimePolicySchemaToTerraform(ctx, &plan, r.client)
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
	err := policyAPI.UpsertRuntimeHostPolicyFiltered(*r.client, planPolicy, deletedRuleNames)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Host Runtime Policy resource",
			"Failed to update host runtime policy: "+err.Error(),
		)
		return
	}

	// Get updated policy value from Prisma Cloud
	updatedPolicy, err := policyAPI.GetRuntimeHostPolicyFiltered(*r.client, plan.GetRuleNames())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Host Runtime Policy resource",
			"Failed to read Host Runtime Policy: "+err.Error(),
		)
		return
	}

	// Convert updated policy into schema
	policySchema, diags := RuntimePolicyTerraformToSchema(ctx, updatedPolicy, plan, r.client)
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

func (r *HostRuntimePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.RuntimeHostPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleNames := state.GetRuleNames()

	// Clear policy rules
	state.Rules = &[]models.RuntimeHostPolicyRuleResourceModel{}

	// Generate API request body from plan
	updatedPlan, diags := RuntimePolicySchemaToTerraform(ctx, &state, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing policy
	err := policyAPI.UpsertRuntimeHostPolicyFiltered(*r.client, updatedPlan, ruleNames)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Host Runtime Policy resource",
			"Failed to delete host runtime policy: "+err.Error(),
		)
		return
	}
}

func (r *HostRuntimePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *HostRuntimePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	util.HCLogDebug(ctx, "entering ModifyPlan")

	//policy.ModifyPolicyResourcePlan(ctx, r.client, req.Plan, resp)

	util.HCLogDebug(ctx, "exiting ModifyPlan")
}

func RuntimePolicySchemaToTerraform(ctx context.Context, plan *models.RuntimeHostPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (policyAPI.RuntimeHostPolicy, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing RuntimePolicySchemaToTerraform")

	var (
		diags diag.Diagnostics
		rules []policyAPI.RuntimeHostPolicyRule
	)

	if plan.Rules != nil {
		customRuleIdMap, err := ruleAPI.GetCustomRuleIdToNameMappings(*client)
		if err != nil {
			diags.AddError(
				"API Error",
				"Error during retrieval of custom rules data: "+err.Error(),
			)
			return policyAPI.RuntimeHostPolicy{}, diags
		}

		rules, diags = RuntimePolicyRulesSchemaToTerraform(ctx, *plan.Rules, client, customRuleIdMap)
		if diags.HasError() {
			return policyAPI.RuntimeHostPolicy{}, diags
		}
	} else {
		rules = []policyAPI.RuntimeHostPolicyRule{}
	}

	tfPolicy := policyAPI.RuntimeHostPolicy{
		Id:    "hostRuntime",
		Owner: "system",
		Rules: &rules,
	}

	//tfPolicy.SortRules(ctx, plan.Rules)

	util.HCLogDebug(ctx, "Finishing RuntimePolicySchemaToTerraform execution")

	return tfPolicy, diags
}

func RuntimePolicyRulesSchemaToTerraform(ctx context.Context, schemaRules []models.RuntimeHostPolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient, customRuleIdMap map[string]int) ([]policyAPI.RuntimeHostPolicyRule, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing RuntimePolicyRulesSchemaToTerraform")

	var (
		diags           diag.Diagnostics
		collectionNames []string
	)

	rules := []policyAPI.RuntimeHostPolicyRule{}

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

		activities, trackSshEvents := activitiesToTerraform(ctx, schemaRule.Activities)
		if diags.HasError() {
			return rules, diags
		}

		antiMalware, diags := antiMalwareToTerraform(ctx, schemaRule.AntiMalware, trackSshEvents)
		if diags.HasError() {
			return rules, diags
		}

		network, dns, diags := runtimeHostNetworkingToTerraform(ctx, schemaRule.Networking)
		if diags.HasError() {
			return rules, diags
		}

		logInspectionRules, diags := logInspectionRulesToTerraform(ctx, schemaRule.LogInspectionRules)
		if diags.HasError() {
			return rules, diags
		}

		fileIntegrityRules, diags := fileIntegrityRulesToTerraform(ctx, schemaRule.FileIntegrityRules)
		if diags.HasError() {
			return rules, diags
		}

		customRules, diags := models.CustomRuntimeRulesToTerraform(ctx, schemaRule.CustomRules, customRuleIdMap)
		if diags.HasError() {
			return rules, diags
		}

		rule := policyAPI.RuntimeHostPolicyRule{
			AntiMalware:        antiMalware,
			Collections:        collections,
			CustomRules:        customRules,
			DNS:                dns,
			Network:            network,
			LogInspectionRules: logInspectionRules,
			FileIntegrityRules: fileIntegrityRules,
			Forensic:           activities,
			Disabled:           schemaRule.Disabled.ValueBool(),
			Modified:           time.Now().Format("2006-01-02T15:04:05.000Z"),
			Name:               schemaRule.Name.ValueString(),
			Notes:              schemaRule.Notes.ValueString(),
			Owner:              schemaRule.Owner.ValueString(),
			PreviousName:       schemaRule.PreviousName.ValueString(),
		}

		rules = append(rules, rule)
	}

	util.HCLogDebug(ctx, "Finishing RuntimePolicyRulesSchemaToTerraform execution")

	return rules, diags
}

func RuntimePolicyTerraformToSchema(ctx context.Context, policy policyAPI.RuntimeHostPolicy, plan models.RuntimeHostPolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (models.RuntimeHostPolicyResourceModel, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing RuntimePolicyTerraformToSchema")

	var (
		diags diag.Diagnostics
		rules []models.RuntimeHostPolicyRuleResourceModel
	)

	if policy.Rules != nil {
		customRuleIdMap, err := ruleAPI.GetCustomRuleNameToIdMappings(*client)
		if err != nil {
			diags.AddError(
				"API Error",
				"Error during retrieval of custom rules data: "+err.Error(),
			)
			return models.RuntimeHostPolicyResourceModel{}, diags
		}

		rules, diags = RuntimePolicyRulesTerraformToSchema(ctx, *policy.Rules, plan.Rules, customRuleIdMap)
		if diags.HasError() {
			return models.RuntimeHostPolicyResourceModel{}, diags
		}
	} else {
		rules = []models.RuntimeHostPolicyRuleResourceModel{}
	}

	schema := models.RuntimeHostPolicyResourceModel{
		Rules: &rules,
	}

	schema.SortRules(ctx, plan.Rules)

	util.HCLogDebug(ctx, "Finishing RuntimeHostPolicyTerraformToSchema execution")

	return schema, diags
}

func RuntimePolicyRulesTerraformToSchema(ctx context.Context, rules []policyAPI.RuntimeHostPolicyRule, planRules *[]models.RuntimeHostPolicyRuleResourceModel, customRuleIdMap map[int]string) ([]models.RuntimeHostPolicyRuleResourceModel, diag.Diagnostics) {
	util.HCLogDebug(ctx, "Executing RuntimeHostPolicyRulesTerraformToSchema")

	var diags diag.Diagnostics

	schemaRules := []models.RuntimeHostPolicyRuleResourceModel{}

	if len(rules) == 0 {
		return schemaRules, diags
	}

	for _, rule := range rules {
		var (
			planRule           models.RuntimeHostPolicyRuleResourceModel
			antiMalware        *models.RuntimeHostPolicyAntiMalwareResourceModel
			skipSshTracking    bool
			customRules        *[]models.RuntimePolicyCustomRuleResourceModel
			activities         *models.RuntimeHostPolicyActivitiesResourceModel
			fileIntegrityRules *[]models.RuntimeHostPolicyFileIntegrityRuleResourceModel
			logInspectionRules *[]models.RuntimeHostPolicyLogInspectionRuleResourceModel
			networking         *models.RuntimeHostPolicyNetworkingResourceModel
			notes              basetypes.StringValue
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
			return []models.RuntimeHostPolicyRuleResourceModel{}, diags
		}

		if planRule.AntiMalware == nil {
			antiMalware = nil
			skipSshTracking = true
		} else {
			antiMalwareValue, skipSshTrackingValue, diags := antiMalwareToSchema(ctx, rule.AntiMalware)
			if diags.HasError() {
				return []models.RuntimeHostPolicyRuleResourceModel{}, diags
			}

			antiMalware = &antiMalwareValue
			skipSshTracking = skipSshTrackingValue
		}

		if planRule.Activities == nil {
			activities = nil
		} else {
			activitesValue, diags := activitiesToSchema(ctx, rule.Forensic, skipSshTracking)
			if diags.HasError() {
				return []models.RuntimeHostPolicyRuleResourceModel{}, diags
			}

			activities = &activitesValue
		}

		if planRule.FileIntegrityRules == nil {
			fileIntegrityRules = nil
		} else {
			fileIntegrityRulesValue, diags := fileIntegrityRulesToSchema(ctx, rule.FileIntegrityRules)
			if diags.HasError() {
				return []models.RuntimeHostPolicyRuleResourceModel{}, diags
			}

			fileIntegrityRules = &fileIntegrityRulesValue
		}

		if planRule.LogInspectionRules == nil {
			logInspectionRules = nil
		} else {
			logInspectionRulesValue, diags := logInspectionRulesToSchema(ctx, rule.LogInspectionRules)
			if diags.HasError() {
				return []models.RuntimeHostPolicyRuleResourceModel{}, diags
			}

			logInspectionRules = &logInspectionRulesValue
		}

		if planRule.Networking == nil {
			networking = nil
		} else {
			networkingValue, diags := runtimeHostNetworkingToSchema(ctx, rule.Network, rule.DNS)
			if diags.HasError() {
				return []models.RuntimeHostPolicyRuleResourceModel{}, diags
			}

			networking = &networkingValue
		}

		if planRule.CustomRules == nil {
			customRules = nil
		} else {
			customRulesValue, diags := models.CustomRuntimeRulesToSchema(ctx, rule.CustomRules, customRuleIdMap)
			if diags.HasError() {
				return []models.RuntimeHostPolicyRuleResourceModel{}, diags
			}

			customRules = &customRulesValue
		}

		if planRule.Notes.IsNull() {
			notes = types.StringNull()
		} else {
			notes = types.StringValue(rule.Notes)
		}

		schemaRule := models.RuntimeHostPolicyRuleResourceModel{
			AntiMalware:        antiMalware,
			FileIntegrityRules: fileIntegrityRules,
			Activities:         activities,
			LogInspectionRules: logInspectionRules,
			Networking:         networking,
			Collections:        collections,
			CustomRules:        customRules,
			Disabled:           types.BoolValue(rule.Disabled),
			Modified:           types.StringValue(""),
			Name:               types.StringValue(rule.Name),
			Notes:              notes,
			Order:              planRule.Order,
			Owner:              types.StringValue(rule.Owner),
			PreviousName:       types.StringValue(rule.PreviousName),
		}

		schemaRules = append(schemaRules, schemaRule)
	}

	util.HCLogDebug(ctx, "Finishing RuntimeHostPolicyRulesTerraformToSchema exection")

	return schemaRules, diags
}

func antiMalwareToTerraform(ctx context.Context, schemaAntiMalware *models.RuntimeHostPolicyAntiMalwareResourceModel, trackSshEvents bool) (policyAPI.RuntimeHostAntiMalware, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		allowedProcesses   []string
		deniedProcesses    policyAPI.DeniedProcesses
		deniedProcessPaths []string
	)

	if schemaAntiMalware == nil {
		return policyAPI.RuntimeHostAntiMalware{}, diags
	}

	if schemaAntiMalware.AllowedProcesses.IsNull() {
		allowedProcesses = []string{}
	} else {
		diags = schemaAntiMalware.AllowedProcesses.ElementsAs(ctx, &allowedProcesses, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostAntiMalware{}, diags
		}
	}

	deniedProcesses = policyAPI.DeniedProcesses{}
	if schemaAntiMalware.DeniedProcesses != nil {
		deniedProcesses.Effect = schemaAntiMalware.DeniedProcesses.Effect.ValueString()

		diags = schemaAntiMalware.DeniedProcesses.Paths.ElementsAs(ctx, &deniedProcessPaths, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostAntiMalware{}, diags
		}

		deniedProcesses.Paths = deniedProcessPaths
	}

	return policyAPI.RuntimeHostAntiMalware{
		AllowedProcesses:              allowedProcesses,
		CryptoMiner:                   schemaAntiMalware.CryptoMiners.ValueString(),
		DeniedProcesses:               deniedProcesses,
		DetectCompilerGeneratedBinary: schemaAntiMalware.SuppressCompilerGeneratedBinaries.ValueBool(),
		EncryptedBinaries:             schemaAntiMalware.EncryptedBinaries.ValueString(),
		ExecutionFlowHijack:           schemaAntiMalware.ExecutionFlowHijacking.ValueString(),
		IntelligenceFeed:              schemaAntiMalware.MalwareFromAdvancedThreatProtection.ValueString(),
		CustomFeed:                    schemaAntiMalware.MalwareFromCustomFeed.ValueString(),
		ReverseShell:                  schemaAntiMalware.ReverseShell.ValueString(),
		ServiceUnknownOriginBinary:    schemaAntiMalware.NonPackagedBinariesService.ValueString(),
		UserUnknownOriginBinary:       schemaAntiMalware.NonPackagedBinariesUser.ValueString(),
		SkipSSHTracking:               !trackSshEvents,
		SuspiciousELFHeaders:          schemaAntiMalware.SuspiciousELFHeaders.ValueString(),
		TempFSProc:                    schemaAntiMalware.ProcessesTemporaryStorage.ValueString(),
		WebShell:                      schemaAntiMalware.WebShell.ValueString(),
		WildFireAnalysis:              schemaAntiMalware.WildFireAnalysis.ValueString(),
	}, diags
}

func antiMalwareToSchema(ctx context.Context, tfAntiMalware policyAPI.RuntimeHostAntiMalware) (models.RuntimeHostPolicyAntiMalwareResourceModel, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	allowedProcesses, diags := types.ListValueFrom(ctx, types.StringType, tfAntiMalware.AllowedProcesses)
	if diags.HasError() {
		return models.RuntimeHostPolicyAntiMalwareResourceModel{}, false, diags
	}

	deniedProcessPaths, diags := types.ListValueFrom(ctx, types.StringType, tfAntiMalware.DeniedProcesses.Paths)
	if diags.HasError() {
		return models.RuntimeHostPolicyAntiMalwareResourceModel{}, false, diags
	}

	deniedProcesses := models.RuntimePolicyDeniedProcessesResourceModel{
		Effect: types.StringValue(tfAntiMalware.DeniedProcesses.Effect),
		Paths:  deniedProcessPaths,
	}

	return models.RuntimeHostPolicyAntiMalwareResourceModel{
		AllowedProcesses:                    allowedProcesses,
		CryptoMiners:                        types.StringValue(tfAntiMalware.CryptoMiner),
		DeniedProcesses:                     &deniedProcesses,
		SuppressCompilerGeneratedBinaries:   types.BoolValue(tfAntiMalware.DetectCompilerGeneratedBinary),
		EncryptedBinaries:                   types.StringValue(tfAntiMalware.EncryptedBinaries),
		ExecutionFlowHijacking:              types.StringValue(tfAntiMalware.ExecutionFlowHijack),
		MalwareFromAdvancedThreatProtection: types.StringValue(tfAntiMalware.IntelligenceFeed),
		MalwareFromCustomFeed:               types.StringValue(tfAntiMalware.CustomFeed),
		ReverseShell:                        types.StringValue(tfAntiMalware.ReverseShell),
		NonPackagedBinariesService:          types.StringValue(tfAntiMalware.ServiceUnknownOriginBinary),
		NonPackagedBinariesUser:             types.StringValue(tfAntiMalware.UserUnknownOriginBinary),
		SuspiciousELFHeaders:                types.StringValue(tfAntiMalware.SuspiciousELFHeaders),
		ProcessesTemporaryStorage:           types.StringValue(tfAntiMalware.TempFSProc),
		WebShell:                            types.StringValue(tfAntiMalware.WebShell),
		WildFireAnalysis:                    types.StringValue(tfAntiMalware.WildFireAnalysis),
	}, tfAntiMalware.SkipSSHTracking, diags
}

func runtimeHostNetworkingToTerraform(ctx context.Context, schemaNetworking *models.RuntimeHostPolicyNetworkingResourceModel) (policyAPI.RuntimeHostNetwork, policyAPI.RuntimeHostDns, diag.Diagnostics) {
	var (
		diags                     diag.Diagnostics
		allowedOutboundIPs        []string
		deniedOutboundIPs         []string
		deniedOutboundPorts       []string
		deniedOutboundPortRanges  []policyAPI.PortRange
		deniedListeningPorts      []string
		deniedListeningPortRanges []policyAPI.PortRange
		allowedDnsDomains         []string
		deniedDnsDomains          []string
	)

	if schemaNetworking == nil {
		return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
	}

	if schemaNetworking.AllowedOutboundIPs.IsNull() {
		allowedOutboundIPs = []string{}
	} else {
		diags = schemaNetworking.AllowedOutboundIPs.ElementsAs(ctx, &allowedOutboundIPs, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
		}
	}

	if schemaNetworking.DeniedListeningPorts.IsNull() {
		deniedListeningPorts = []string{}
	} else {
		diags = schemaNetworking.DeniedListeningPorts.ElementsAs(ctx, &deniedListeningPorts, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
		}
	}

	if schemaNetworking.DeniedOutboundIPs.IsNull() {
		deniedOutboundIPs = []string{}
	} else {
		diags = schemaNetworking.DeniedOutboundIPs.ElementsAs(ctx, &deniedOutboundIPs, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
		}
	}

	if schemaNetworking.DeniedOutboundPorts.IsNull() {
		deniedOutboundPorts = []string{}
	} else {
		diags = schemaNetworking.DeniedOutboundPorts.ElementsAs(ctx, &deniedOutboundPorts, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
		}
	}

	if schemaNetworking.AllowedDnsDomains.IsNull() {
		allowedDnsDomains = []string{}
	} else {
		diags = schemaNetworking.AllowedDnsDomains.ElementsAs(ctx, &allowedDnsDomains, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
		}
	}

	if schemaNetworking.DeniedDnsDomains.IsNull() {
		deniedDnsDomains = []string{}
	} else {
		diags = schemaNetworking.DeniedDnsDomains.ElementsAs(ctx, &deniedDnsDomains, false)
		if diags.HasError() {
			return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
		}
	}

	deniedListeningPortRanges = []policyAPI.PortRange{}
	for _, port := range deniedListeningPorts {
		if strings.ContainsAny(port, "-") {
			splitPort := strings.Split(port, "-")

			start, err := strconv.Atoi(splitPort[0])
			if err != nil {
				diags.AddError(
					"Value Conversion Error",
					err.Error(),
				)
				return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
			}

			end, err := strconv.Atoi(splitPort[1])
			if err != nil {
				diags.AddError(
					"Value Conversion Error",
					err.Error(),
				)
				return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
			}

			deniedListeningPortRanges = append(deniedListeningPortRanges, policyAPI.PortRange{
				Start: start,
				End:   end,
			})
		} else {
			portInt, err := strconv.Atoi(port)
			if err != nil {
				diags.AddError(
					"Value Conversion Error",
					err.Error(),
				)
				return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
			}

			deniedListeningPortRanges = append(deniedListeningPortRanges, policyAPI.PortRange{
				Start: portInt,
				End:   portInt,
			})
		}
	}

	deniedOutboundPortRanges = []policyAPI.PortRange{}
	for _, port := range deniedOutboundPorts {
		if strings.ContainsAny(port, "-") {
			splitPort := strings.Split(port, "-")

			start, err := strconv.Atoi(splitPort[0])
			if err != nil {
				diags.AddError(
					"Value Conversion Error",
					err.Error(),
				)
				return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
			}

			end, err := strconv.Atoi(splitPort[1])
			if err != nil {
				diags.AddError(
					"Value Conversion Error",
					err.Error(),
				)
				return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
			}

			deniedOutboundPortRanges = append(deniedOutboundPortRanges, policyAPI.PortRange{
				Start: start,
				End:   end,
			})
		} else {
			portInt, err := strconv.Atoi(port)
			if err != nil {
				diags.AddError(
					"Value Conversion Error",
					err.Error(),
				)
				return policyAPI.RuntimeHostNetwork{}, policyAPI.RuntimeHostDns{}, diags
			}

			deniedOutboundPortRanges = append(deniedOutboundPortRanges, policyAPI.PortRange{
				Start: portInt,
				End:   portInt,
			})
		}
	}

	network := policyAPI.RuntimeHostNetwork{
		AllowedOutboundIPs:   allowedOutboundIPs,
		CustomFeed:           schemaNetworking.SuspiciousIPsCustomFeed.ValueString(),
		DeniedListeningPorts: deniedListeningPortRanges,
		DeniedOutboundIPs:    deniedOutboundIPs,
		DeniedOutboundPorts:  deniedOutboundPortRanges,
		DenyListEffect:       schemaNetworking.DeniedIPsPortsEffect.ValueString(),
		IntelligenceFeed:     schemaNetworking.SuspiciousIPsAdvancedThreatProtectionEffect.ValueString(),
	}

	dns := policyAPI.RuntimeHostDns{
		Allow:            allowedDnsDomains,
		Deny:             deniedDnsDomains,
		DenyListEffect:   schemaNetworking.DeniedDnsDomainsEffect.ValueString(),
		IntelligenceFeed: schemaNetworking.SuspiciousDomainsAdvancedThreatProtectionEffect.ValueString(),
	}

	return network, dns, diags
}

func runtimeHostNetworkingToSchema(ctx context.Context, tfNetwork policyAPI.RuntimeHostNetwork, tfDns policyAPI.RuntimeHostDns) (models.RuntimeHostPolicyNetworkingResourceModel, diag.Diagnostics) {
	var (
		diags                      diag.Diagnostics
		deniedListeningPortsValues = []string{}
		deniedOutboundPortsValues  = []string{}
	)

	allowedOutboundIPs, diags := types.ListValueFrom(ctx, types.StringType, tfNetwork.AllowedOutboundIPs)
	if diags.HasError() {
		return models.RuntimeHostPolicyNetworkingResourceModel{}, diags
	}

	deniedOutboundIPs, diags := types.ListValueFrom(ctx, types.StringType, tfNetwork.DeniedOutboundIPs)
	if diags.HasError() {
		return models.RuntimeHostPolicyNetworkingResourceModel{}, diags
	}

	allowedDnsDomains, diags := types.ListValueFrom(ctx, types.StringType, tfDns.Allow)
	if diags.HasError() {
		return models.RuntimeHostPolicyNetworkingResourceModel{}, diags
	}

	deniedDnsDomains, diags := types.ListValueFrom(ctx, types.StringType, tfDns.Deny)
	if diags.HasError() {
		return models.RuntimeHostPolicyNetworkingResourceModel{}, diags
	}

	for _, portRange := range tfNetwork.DeniedListeningPorts {
		start := strconv.Itoa(portRange.Start)
		end := strconv.Itoa(portRange.End)

		if start == end {
			deniedListeningPortsValues = append(deniedListeningPortsValues, start)
		} else {
			deniedListeningPortsValues = append(deniedListeningPortsValues, start+"-"+end)
		}
	}

	deniedListeningPorts, diags := types.ListValueFrom(ctx, types.StringType, deniedListeningPortsValues)
	if diags.HasError() {
		return models.RuntimeHostPolicyNetworkingResourceModel{}, diags
	}

	for _, portRange := range tfNetwork.DeniedOutboundPorts {
		start := strconv.Itoa(portRange.Start)
		end := strconv.Itoa(portRange.End)

		if start == end {
			deniedOutboundPortsValues = append(deniedOutboundPortsValues, start)
		} else {
			deniedOutboundPortsValues = append(deniedOutboundPortsValues, start+"-"+end)
		}
	}

	deniedOutboundPorts, diags := types.ListValueFrom(ctx, types.StringType, deniedOutboundPortsValues)
	if diags.HasError() {
		return models.RuntimeHostPolicyNetworkingResourceModel{}, diags
	}

	return models.RuntimeHostPolicyNetworkingResourceModel{
		AllowedOutboundIPs:                              allowedOutboundIPs,
		DeniedOutboundIPs:                               deniedOutboundIPs,
		DeniedOutboundPorts:                             deniedOutboundPorts,
		DeniedListeningPorts:                            deniedListeningPorts,
		DeniedIPsPortsEffect:                            types.StringValue(tfNetwork.DenyListEffect),
		SuspiciousIPsAdvancedThreatProtectionEffect:     types.StringValue(tfNetwork.IntelligenceFeed),
		SuspiciousIPsCustomFeed:                         types.StringValue(tfNetwork.CustomFeed),
		AllowedDnsDomains:                               allowedDnsDomains,
		DeniedDnsDomains:                                deniedDnsDomains,
		DeniedDnsDomainsEffect:                          types.StringValue(tfDns.DenyListEffect),
		SuspiciousDomainsAdvancedThreatProtectionEffect: types.StringValue(tfDns.IntelligenceFeed),
	}, diags
}

func logInspectionRulesToTerraform(ctx context.Context, schemaLogInspectionRules *[]models.RuntimeHostPolicyLogInspectionRuleResourceModel) ([]policyAPI.LogInspectionRule, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		logInspectionRules = []policyAPI.LogInspectionRule{}
	)

	if schemaLogInspectionRules == nil {
		return logInspectionRules, diags
	}

	for _, rule := range *schemaLogInspectionRules {
		regex := []string{}

		diags = rule.Regex.ElementsAs(ctx, &regex, false)
		if diags.HasError() {
			return logInspectionRules, diags
		}

		logInspectionRules = append(logInspectionRules, policyAPI.LogInspectionRule{
			Path:  rule.Path.ValueString(),
			Regex: regex,
		})
	}

	return logInspectionRules, diags
}

func logInspectionRulesToSchema(ctx context.Context, tfLogInspectionRules []policyAPI.LogInspectionRule) ([]models.RuntimeHostPolicyLogInspectionRuleResourceModel, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		logInspectionRules = []models.RuntimeHostPolicyLogInspectionRuleResourceModel{}
	)

	for _, rule := range tfLogInspectionRules {
		regex, diags := types.ListValueFrom(ctx, types.StringType, rule.Regex)
		if diags.HasError() {
			return logInspectionRules, diags
		}

		logInspectionRules = append(logInspectionRules, models.RuntimeHostPolicyLogInspectionRuleResourceModel{
			Path:  types.StringValue(rule.Path),
			Regex: regex,
		})
	}

	return logInspectionRules, diags
}

func fileIntegrityRulesToTerraform(ctx context.Context, schemaFileIntegrityRules *[]models.RuntimeHostPolicyFileIntegrityRuleResourceModel) ([]policyAPI.FileIntegrityRule, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		fileIntegrityRules = []policyAPI.FileIntegrityRule{}
	)

	if schemaFileIntegrityRules == nil {
		return fileIntegrityRules, diags
	}

	for _, rule := range *schemaFileIntegrityRules {
		exclusions := []string{}
		diags = rule.ExcludedFilePatterns.ElementsAs(ctx, &exclusions, false)
		if diags.HasError() {
			return fileIntegrityRules, diags
		}

		procWhitelist := []string{}
		diags = rule.AllowedProcesses.ElementsAs(ctx, &procWhitelist, false)
		if diags.HasError() {
			return fileIntegrityRules, diags
		}

		fileIntegrityRules = append(fileIntegrityRules, policyAPI.FileIntegrityRule{
			Exclusions:    exclusions,
			Path:          rule.Path.ValueString(),
			ProcWhitelist: procWhitelist,
			Recursive:     rule.MonitorSubdirectories.ValueBool(),
			Read:          rule.MonitorReadOps.ValueBool(),
			Write:         rule.MonitorWriteOps.ValueBool(),
			Metadata:      rule.MonitorMetadataChanges.ValueBool(),
		})
	}

	return fileIntegrityRules, diags
}

func fileIntegrityRulesToSchema(ctx context.Context, tfFileIntegrityRules []policyAPI.FileIntegrityRule) ([]models.RuntimeHostPolicyFileIntegrityRuleResourceModel, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		fileIntegrityRules = []models.RuntimeHostPolicyFileIntegrityRuleResourceModel{}
	)

	for _, rule := range tfFileIntegrityRules {
		excludedFilePatterns, diags := types.ListValueFrom(ctx, types.StringType, rule.Exclusions)
		if diags.HasError() {
			return fileIntegrityRules, diags
		}

		// NOTE: hotfix until this can be re-written
		if excludedFilePatterns.IsNull() {
			excludedFilePatterns = types.ListValueMust(types.StringType, []attr.Value{})
		}

		allowedProcesses, diags := types.ListValueFrom(ctx, types.StringType, rule.ProcWhitelist)
		if diags.HasError() {
			return fileIntegrityRules, diags
		}
		// NOTE: hotfix until this can be re-written
		if allowedProcesses.IsNull() {
			allowedProcesses = types.ListValueMust(types.StringType, []attr.Value{})
		}

		fileIntegrityRules = append(fileIntegrityRules, models.RuntimeHostPolicyFileIntegrityRuleResourceModel{
			Path:                   types.StringValue(rule.Path),
			AllowedProcesses:       allowedProcesses,
			ExcludedFilePatterns:   excludedFilePatterns,
			MonitorSubdirectories:  types.BoolValue(rule.Recursive),
			MonitorWriteOps:        types.BoolValue(rule.Write),
			MonitorReadOps:         types.BoolValue(rule.Read),
			MonitorMetadataChanges: types.BoolValue(rule.Metadata),
		})
	}

	return fileIntegrityRules, diags
}

func activitiesToTerraform(ctx context.Context, schemaActivities *models.RuntimeHostPolicyActivitiesResourceModel) (policyAPI.Forensic, bool) {
	if schemaActivities == nil {
		return policyAPI.Forensic{}, false
	}

	return policyAPI.Forensic{
		ActivitiesDisabled:       !schemaActivities.HostActivityMonitoring.Enabled.ValueBool(),
		DockerEnabled:            schemaActivities.HostActivityMonitoring.DockerCommands.Enabled.ValueBool(),
		ReadonlyDockerEnabled:    schemaActivities.HostActivityMonitoring.DockerCommands.IncludeReadOnlyEvents.ValueBool(),
		ServiceActivitiesEnabled: schemaActivities.HostActivityMonitoring.LogBackgroundApps.ValueBool(),
		SshdEnabled:              schemaActivities.HostActivityMonitoring.SshdSessions.ValueBool(),
		SudoEnabled:              schemaActivities.HostActivityMonitoring.SudoCommands.ValueBool(),
	}, schemaActivities.TrackSshEvents.ValueBool()
}

func activitiesToSchema(ctx context.Context, tfActivities policyAPI.Forensic, skipSshTracking bool) (models.RuntimeHostPolicyActivitiesResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	return models.RuntimeHostPolicyActivitiesResourceModel{
		HostActivityMonitoring: models.RuntimeHostPolicyActivityMonitoringResourceModel{
			Enabled: types.BoolValue(!tfActivities.ActivitiesDisabled),
			DockerCommands: models.RuntimeHostPolicyDockerActivityMonitoringResourceModel{
				Enabled:               types.BoolValue(tfActivities.DockerEnabled),
				IncludeReadOnlyEvents: types.BoolValue(tfActivities.ReadonlyDockerEnabled),
			},
			SshdSessions:      types.BoolValue(tfActivities.SshdEnabled),
			SudoCommands:      types.BoolValue(tfActivities.SudoEnabled),
			LogBackgroundApps: types.BoolValue(tfActivities.ServiceActivitiesEnabled),
		},
		TrackSshEvents: types.BoolValue(!skipSshTracking),
	}, diags
}

//func customRulesToTerraform(ctx context.Context, schemaCustomRules *[]models.RuntimeHostPolicyCustomRuleResourceModel, customRuleIdMap map[string]int) ([]policyAPI.CustomRule, diag.Diagnostics) {
//    // TODO: find a way to have the error message denote which rule has the error
//
//    var diags diag.Diagnostics
//
//    if schemaCustomRules == nil {
//        return []policyAPI.CustomRule{}, diags
//    }
//
//    tfCustomRules := []policyAPI.CustomRule{}
//
//    for _, schemaCustomRule := range *schemaCustomRules {
//        customRuleId, ok := customRuleIdMap[schemaCustomRule.Name.ValueString()]
//        if !ok {
//            diags.AddError(
//                "Value Conversion Error",
//                fmt.Sprintf("No matching custom rule found for specified rule name \"%s\"", schemaCustomRule.Name.ValueString()),
//            )
//
//            return []policyAPI.CustomRule{}, diags
//        }
//
//        tfCustomRules = append(tfCustomRules, policyAPI.CustomRule{
//            ID: customRuleId,
//            Action: schemaCustomRule.LogAs.ValueString(),
//            Effect: schemaCustomRule.Effect.ValueString(),
//        })
//    }
//
//    return tfCustomRules, diags
//}
//
//func customRulesToSchema(ctx context.Context, tfCustomRules []policyAPI.CustomRule, customRuleIdMap map[int]string) ([]models.RuntimeHostPolicyCustomRuleResourceModel, diag.Diagnostics) {
//    var diags diag.Diagnostics
//
//    schemaCustomRules := []models.RuntimeHostPolicyCustomRuleResourceModel{}
//
//    for _, tfCustomRule := range tfCustomRules {
//        customRuleName, ok := customRuleIdMap[tfCustomRule.ID]
//        if !ok {
//            diags.AddError(
//                "Value Conversion Error",
//                fmt.Sprintf("No matching custom rule found for specified rule ID %d", tfCustomRule.ID),
//            )
//
//            return []models.RuntimeHostPolicyCustomRuleResourceModel{}, diags
//        }
//
//        schemaCustomRules = append(schemaCustomRules, models.RuntimeHostPolicyCustomRuleResourceModel{
//            ID: types.Int64Value(int64(tfCustomRule.ID)),
//            Name: types.StringValue(customRuleName),
//            LogAs: types.StringValue(tfCustomRule.Action),
//            Effect: types.StringValue(tfCustomRule.Effect),
//        })
//    }
//
//    return schemaCustomRules, diags
//}
