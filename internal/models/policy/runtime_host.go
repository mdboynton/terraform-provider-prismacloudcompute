package models

import (
	"cmp"
	"context"
	"slices"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RuntimeHostPolicyResourceModel struct {
	Rules *[]RuntimeHostPolicyRuleResourceModel `tfsdk:"rules"`
}

func (m *RuntimeHostPolicyResourceModel) GetRuleNames() []string {
	ruleNames := []string{}

	if m.Rules == nil || len(*m.Rules) == 0 {
		return ruleNames
	}

	for _, rule := range *m.Rules {
		ruleNames = append(ruleNames, rule.Name.ValueString())
	}

	return ruleNames
}

type RuntimeHostPolicyRuleResourceModel struct {
	AntiMalware        *RuntimeHostPolicyAntiMalwareResourceModel         `tfsdk:"anti_malware"`
	Collections        types.Set                                          `tfsdk:"collections"`
	CustomRules        *[]RuntimePolicyCustomRuleResourceModel            `tfsdk:"custom_rules"`
	Disabled           types.Bool                                         `tfsdk:"disabled"`
	FileIntegrityRules *[]RuntimeHostPolicyFileIntegrityRuleResourceModel `tfsdk:"file_integrity_rules"`
	Activities         *RuntimeHostPolicyActivitiesResourceModel          `tfsdk:"activities"`
	LogInspectionRules *[]RuntimeHostPolicyLogInspectionRuleResourceModel `tfsdk:"log_inspection_rules"`
	Modified           types.String                                       `tfsdk:"modified"`
	Name               types.String                                       `tfsdk:"name"`
	Networking         *RuntimeHostPolicyNetworkingResourceModel          `tfsdk:"networking"`
	Notes              types.String                                       `tfsdk:"notes"`
	Order              types.Int32                                        `tfsdk:"order"`
	Owner              types.String                                       `tfsdk:"owner"`
	PreviousName       types.String                                       `tfsdk:"previous_name"`
}

type RuntimeHostPolicyAntiMalwareResourceModel struct {
	AllowedProcesses                    types.List                                 `tfsdk:"allowed_processes"`
	CryptoMiners                        types.String                               `tfsdk:"crypto_miners"`
	DeniedProcesses                     *RuntimePolicyDeniedProcessesResourceModel `tfsdk:"denied_processes"`
	SuppressCompilerGeneratedBinaries   types.Bool                                 `tfsdk:"suppress_compiler_generated_binaries"`
	EncryptedBinaries                   types.String                               `tfsdk:"encrypted_binaries"`
	ExecutionFlowHijacking              types.String                               `tfsdk:"execution_flow_hijacking"`
	MalwareFromAdvancedThreatProtection types.String                               `tfsdk:"malware_from_advanced_threat_protection"`
	MalwareFromCustomFeed               types.String                               `tfsdk:"malware_from_custom_feed"`
	ReverseShell                        types.String                               `tfsdk:"reverse_shell"`
	NonPackagedBinariesService          types.String                               `tfsdk:"non_packaged_binaries_service"`
	NonPackagedBinariesUser             types.String                               `tfsdk:"non_packaged_binaries_user"`
	SuspiciousELFHeaders                types.String                               `tfsdk:"suspicious_elf_headers"`
	ProcessesTemporaryStorage           types.String                               `tfsdk:"processes_temporary_storage"`
	WebShell                            types.String                               `tfsdk:"web_shell"`
	WildFireAnalysis                    types.String                               `tfsdk:"wild_fire_analysis"`
}

type RuntimeHostPolicyFileIntegrityRuleResourceModel struct {
	Path                   types.String `tfsdk:"file_path"`
	AllowedProcesses       types.List   `tfsdk:"allowed_processes"`
	ExcludedFilePatterns   types.List   `tfsdk:"excluded_file_patterns"`
	MonitorSubdirectories  types.Bool   `tfsdk:"monitor_subdirectories"`
	MonitorWriteOps        types.Bool   `tfsdk:"monitor_write_ops"`
	MonitorReadOps         types.Bool   `tfsdk:"monitor_read_ops"`
	MonitorMetadataChanges types.Bool   `tfsdk:"monitor_metadata_changes"`
}

type RuntimeHostPolicyActivitiesResourceModel struct {
	HostActivityMonitoring RuntimeHostPolicyActivityMonitoringResourceModel `tfsdk:"host_activity_monitoring"`
	TrackSshEvents         types.Bool                                       `tfsdk:"track_ssh_events"`
}

type RuntimeHostPolicyActivityMonitoringResourceModel struct {
	Enabled           types.Bool                                             `tfsdk:"enabled"`
	DockerCommands    RuntimeHostPolicyDockerActivityMonitoringResourceModel `tfsdk:"docker_commands"`
	SshdSessions      types.Bool                                             `tfsdk:"sshd_sessions"`
	SudoCommands      types.Bool                                             `tfsdk:"sudo_commands"`
	LogBackgroundApps types.Bool                                             `tfsdk:"log_background_apps"`
}

type RuntimeHostPolicyDockerActivityMonitoringResourceModel struct {
	Enabled               types.Bool `tfsdk:"enabled"`
	IncludeReadOnlyEvents types.Bool `tfsdk:"include_read_only_events"`
}

type RuntimeHostPolicyLogInspectionRuleResourceModel struct {
	Path  types.String `tfsdk:"path"`
	Regex types.List   `tfsdk:"regex"`
}

type RuntimeHostPolicyNetworkingResourceModel struct {
	AllowedOutboundIPs                              types.List   `tfsdk:"allowed_outbound_ips"`
	DeniedListeningPorts                            types.List   `tfsdk:"denied_listening_ports"`
	DeniedOutboundIPs                               types.List   `tfsdk:"denied_outbound_ips"`
	DeniedOutboundPorts                             types.List   `tfsdk:"denied_outbound_ports"`
	DeniedIPsPortsEffect                            types.String `tfsdk:"denied_ips_ports_effect"`
	SuspiciousIPsCustomFeed                         types.String `tfsdk:"suspicious_ips_custom_feed"`
	SuspiciousIPsAdvancedThreatProtectionEffect     types.String `tfsdk:"suspicious_ips_advanced_threat_protection_effect"`
	AllowedDnsDomains                               types.List   `tfsdk:"allowed_dns_domains"`
	DeniedDnsDomains                                types.List   `tfsdk:"denied_dns_domains"`
	DeniedDnsDomainsEffect                          types.String `tfsdk:"denied_dns_domains_effect"`
	SuspiciousDomainsAdvancedThreatProtectionEffect types.String `tfsdk:"suspicious_domains_advanced_threat_protection_effect"`
}

func (m *RuntimeHostPolicyResourceModel) SortRules(ctx context.Context, planRules *[]RuntimeHostPolicyRuleResourceModel) {
	util.LogDebug(ctx, "Executing RuntimeHostPolicyResourceModel.SortRules()")

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

		slices.SortFunc(*m.Rules, func(a, b RuntimeHostPolicyRuleResourceModel) int {
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

	util.LogDebug(ctx, "Finishing RuntimeHostPolicyResourceModel.SortRules() execution")
}
