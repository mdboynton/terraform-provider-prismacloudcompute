package models

import (
	"cmp"
	"context"
	"slices"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RuntimeContainerPolicyResourceModel struct {
	AutomaticRuntimeLearning types.Bool                                 `tfsdk:"automatic_runtime_learning"`
	Rules                    *[]RuntimeContainerPolicyRuleResourceModel `tfsdk:"rules"`
}

func (m *RuntimeContainerPolicyResourceModel) GetRuleNames() []string {
	ruleNames := []string{}

	if m.Rules == nil || len(*m.Rules) == 0 {
		return ruleNames
	}

	for _, rule := range *m.Rules {
		ruleNames = append(ruleNames, rule.Name.ValueString())
	}

	return ruleNames
}

type RuntimeContainerPolicyRuleResourceModel struct {
	AntiMalware  *RuntimeContainerPolicyAntiMalwareResourceModel `tfsdk:"anti_malware"`
	Processes    *RuntimeContainerPolicyProcessesResourceModel   `tfsdk:"processes"`
	Networking   *RuntimeContainerPolicyNetworkingResourceModel  `tfsdk:"networking"`
	FileSystem   *RuntimeContainerPolicyFileSystemResourceModel  `tfsdk:"file_system"`
	CustomRules  *[]RuntimePolicyCustomRuleResourceModel         `tfsdk:"custom_rules"`
	Collections  types.Set                                       `tfsdk:"collections"`
	Disabled     types.Bool                                      `tfsdk:"disabled"`
	Modified     types.String                                    `tfsdk:"modified"`
	Name         types.String                                    `tfsdk:"name"`
	Notes        types.String                                    `tfsdk:"notes"`
	Order        types.Int32                                     `tfsdk:"order"`
	Owner        types.String                                    `tfsdk:"owner"`
	PreviousName types.String                                    `tfsdk:"previous_name"`
}

type RuntimeContainerPolicyAntiMalwareResourceModel struct {
	MalwareFromAdvancedThreatProtection types.String `tfsdk:"malware_from_advanced_threat_protection"`
	KubernetesAttacks                   types.String `tfsdk:"kubernetes_attacks"`
	SuspiciousCloudProviderApiQueries   types.String `tfsdk:"suspicious_cloud_provider_api_queries"`
	WildFireAnalysis                    types.String `tfsdk:"wild_fire_analysis"`
}

type RuntimeContainerPolicyProcessesResourceModel struct {
	Enabled                              types.Bool                                 `tfsdk:"enabled"`
	AllowedProcesses                     types.List                                 `tfsdk:"allowed_processes"`
	AllowOnlyLearnedProcessesFromParents types.Bool                                 `tfsdk:"allow_only_learned_processes_from_parents"`
	AllowAllActivityInAttachedSessions   types.Bool                                 `tfsdk:"allow_all_activity_in_attached_sessions"`
	ProcessesFromModifiedBinaries        types.String                               `tfsdk:"processes_from_modified_binaries"`
	CryptoMiners                         types.String                               `tfsdk:"crypto_miners"`
	ReverseShell                         types.String                               `tfsdk:"reverse_shell"`
	LateralMovementProcesses             types.String                               `tfsdk:"lateral_movement_processes"`
	ProcessesStartedWithSUID             types.String                               `tfsdk:"processes_started_with_suid"`
	DeniedProcesses                      *RuntimePolicyDeniedProcessesResourceModel `tfsdk:"denied_processes"`
	AllOtherProcessesEffect              types.String                               `tfsdk:"all_other_processes_effect"`
}

type RuntimeContainerPolicyNetworkingResourceModel struct {
	IpConnectivityEnabled               types.Bool   `tfsdk:"ip_connectivity_enabled"`
	AllowedListeningPorts               types.List   `tfsdk:"allowed_listening_ports"`
	AllowedOutboundInternetPorts        types.List   `tfsdk:"allowed_outbound_internet_ports"`
	AllowedOutboundIPs                  types.List   `tfsdk:"allowed_outbound_ips"`
	NetworkActivityFromModifiedBinaries types.String `tfsdk:"network_activity_from_modified_binaries"`
	PortScanning                        types.String `tfsdk:"port_scanning"`
	RawSockets                          types.String `tfsdk:"raw_sockets"`
	DeniedListeningPorts                types.List   `tfsdk:"denied_listening_ports"`
	DeniedListeningPortEffect           types.String `tfsdk:"denied_listening_ports_effect"`
	DeniedOutboundInternetPorts         types.List   `tfsdk:"denied_outbound_internet_ports"`
	DeniedOutboundInternetPortsEffect   types.String `tfsdk:"denied_outbound_internet_ports_effect"`
	DeniedOutboundIPs                   types.List   `tfsdk:"denied_outbound_ips"`
	DeniedOutboundIPsEffect             types.String `tfsdk:"denied_outbound_ips_effect"`
	AllOtherActivityEffect              types.String `tfsdk:"all_other_activity_effect"`
	DnsEnabled                          types.Bool   `tfsdk:"dns_enabled"`
	AllowedDnsDomains                   types.List   `tfsdk:"allowed_dns_domains"`
	DeniedDnsDomains                    types.List   `tfsdk:"denied_dns_domains"`
	DeniedDnsDomainsEffect              types.String `tfsdk:"denied_dns_domains_effect"`
	AllOtherDomainsEffect               types.String `tfsdk:"all_other_domains_effect"`
}

type RuntimeContainerPolicyFileSystemResourceModel struct {
	Enabled                             types.Bool   `tfsdk:"enabled"`
	AllOtherPathsEffect                 types.String `tfsdk:"all_other_paths_effect"`
	AllowedPaths                        types.List   `tfsdk:"allowed_paths"`
	ChangesToBinaries                   types.String `tfsdk:"changes_to_binaries"`
	ChangesToSshAdminAccountConfigFiles types.String `tfsdk:"changes_to_ssh_admin_account_config_files"`
	DeniedPaths                         types.List   `tfsdk:"denied_paths"`
	DeniedPathsEffect                   types.String `tfsdk:"denied_paths_effect"`
	DetectionOfEncryptedBinaries        types.String `tfsdk:"detection_of_encrypted_binaries"`
	SuspiciousELFHeaders                types.String `tfsdk:"suspicious_elf_headers"`
}

func (m *RuntimeContainerPolicyResourceModel) SortRules(ctx context.Context, planRules *[]RuntimeContainerPolicyRuleResourceModel) {
	util.LogDebug(ctx, "Executing RuntimeContainerPolicyResourceModel.SortRules()")

	if m.Rules != nil && len(*m.Rules) > 0 {
		// TODO: check for mismatched lengths and return diags with error
		// (if there's additional rules added outside of terraform, they will be reflected in m.SortRules. therefore,
		// if this happens, return a diag error with message suggesting they check those rules, until we can put in logic
		// to be able to automatically handle that scenario)

		//if len(*m.Rules) != len(*planRules) {
		//}

		if len(*m.Rules) == 1 && len(*planRules) == 1 {
			(*m.Rules)[0].Order = (*planRules)[0].Order
			return
		}

		if planRules == nil {
			for i := 0; i < len(*m.Rules); i++ {
				(*m.Rules)[i].Order = types.Int32Value(int32(i + 1))
			}
			return
		}

		ruleOrderMap := make(map[string]int32)
		for index, planRule := range *planRules {
			ruleOrderMap[planRule.Name.ValueString()] = int32(index)
		}

		slices.SortFunc(*m.Rules, func(a, b RuntimeContainerPolicyRuleResourceModel) int {
			orderA, okA := ruleOrderMap[a.Name.ValueString()]
			if !okA {
				orderA = int32(len(ruleOrderMap) + 1)
			}
			orderB, okB := ruleOrderMap[b.Name.ValueString()]
			if !okB {
				orderB = int32(len(ruleOrderMap) + 1)
			}
			return cmp.Compare(orderA, orderB)
		})

		for i := 0; i < len(*m.Rules); i++ {
			(*m.Rules)[i].Order = (*planRules)[i].Order
		}
	}

	util.LogDebug(ctx, "Finishing RuntimeContainerPolicyResourceModel.SortRules() execution")
}
