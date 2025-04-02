package models

import (
    "context"
    "slices"
    "cmp"
    
    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Generic Policy Resource Models

type PolicyResourceModel struct {
    Id              types.String                            `tfsdk:"id"`
    PolicyType      types.String                            `tfsdk:"policy_type"`
    PolicyContext   types.String                            `tfsdk:"policy_context"`
    Rules           *[]PolicyRuleResourceModel              `tfsdk:"rules"`
    Type            types.String                            `tfsdk:"type"`
}

type PolicyRuleResourceModel struct {
    AlertThreshold                  *PolicyRuleThresholdResourceModel `tfsdk:"alert_threshold" json:"alertThreshold"`
    BlockMessage                    types.String    `tfsdk:"block_message"`
    BlockThreshold                  *PolicyRuleThresholdResourceModel `tfsdk:"block_threshold" json:"blockThreshold"`
    Collections                     types.Set `tfsdk:"collections"`
    CVERules                        *[]PolicyRuleExceptionResourceModel  `tfsdk:"cve_rules"`
    ComplianceActions *PolicyRuleComplianceActionsResourceModel    `tfsdk:"compliance_actions"`
    Disabled                        types.Bool      `tfsdk:"disabled"`
    Effect                          types.String    `tfsdk:"effect"`
    ExcludeBaseImageVulns           types.Bool      `tfsdk:"exclude_base_image_vulns"`
    GraceDays                       types.Int32     `tfsdk:"grace_days_all_severities"`
    GraceDaysPolicy                 *PolicyRuleGraceDaysPolicyResourceModel `tfsdk:"grace_days_by_severity"` 
    Modified                        types.String    `tfsdk:"modified"`
    Name                            types.String    `tfsdk:"name"`
    Notes                           types.String    `tfsdk:"notes"`
    OnlyFixed                       types.Bool      `tfsdk:"apply_only_when_fix_available"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner                           types.String    `tfsdk:"owner"`
    PkgTypesThresholds              *[]PolicyRulePkgTypesThresholdResourceModel    `tfsdk:"package_types_thresholds"`
    ReportAllPassedAndFailedChecks  types.Bool      `tfsdk:"report_passed_and_failed_checks"`
    Tags                            *[]PolicyRuleExceptionResourceModel  `tfsdk:"tags"`
    Verbose                         types.Bool      `tfsdk:"verbose"`
}

type PolicyRuleThresholdResourceModel struct {
    Threshold   types.String    `tfsdk:"threshold"`
    RiskFactors types.Set       `tfsdk:"risk_factors"`
}

type PolicyRuleGraceDaysPolicyResourceModel struct {
    Low         types.Int32     `tfsdk:"low"`
    Medium      types.Int32     `tfsdk:"medium"`
    High        types.Int32     `tfsdk:"high"`
    Critical    types.Int32     `tfsdk:"critical"`
}

type PolicyRulePkgTypesThresholdResourceModel struct {
    Type            types.String    `tfsdk:"type"`
    AlertThreshold  types.String    `tfsdk:"alert_threshold"`
    BlockThreshold  types.String    `tfsdk:"block_threshold"`
}

type PolicyRuleExceptionResourceModel struct {
    //Name            types.String                                            `tfsdk:"name"`
    Effect          types.String                                            `tfsdk:"effect"`
    Id              types.String                                            `tfsdk:"id"`
    Description     types.String                                            `tfsdk:"description"`
    //Type            types.String                                            `tfsdk:"type"`
    Expiration      PolicyRuleExceptionExpirationResourceModel              `tfsdk:"expiration"`
}

type PolicyRuleExceptionExpirationResourceModel struct {
    Enabled types.Bool      `tfsdk:"enabled"`
    Date    types.String    `tfsdk:"date"`
}

type PolicyRuleComplianceActionsResourceModel struct {
    Template    types.String        `tfsdk:"template"`
    Types       types.Set           `tfsdk:"types"`
    Severities  types.Set           `tfsdk:"severities"`
    Checks      []PolicyRuleComplianceActionCheckResourceModel `tfsdk:"checks"`
}

type PolicyRuleComplianceActionCheckResourceModel struct {
    Id types.Int32      `tfsdk:"id"`
    Action types.String `tfsdk:"action"`
}

func (m *PolicyResourceModel) SortRules(ctx context.Context, planRules *[]PolicyRuleResourceModel) {
    util.DLog(ctx, "Executing PolicyResourceModel.SortRules()")

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

        slices.SortFunc(*m.Rules, func(a, b PolicyRuleResourceModel) int {
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
    
    util.DLog(ctx, "Finishing PolicyResourceModel.SortRules() execution")
}


// Application Control Policy Models

type ApplicationControlPolicyResourceModel struct {
    Rules       *[]ApplicationControlPolicyRuleResourceModel    `tfsdk:"rules"`
}

type ApplicationControlPolicyRuleResourceModel struct {
    Id              types.Int32     `tfsdk:"id"`
    Applications    types.Set       `tfsdk:"applications"`
    Description     types.String    `tfsdk:"description"`
    //Disabled        types.Bool      `tfsdk:"disabled"` // TODO: is this even supported? doesnt appear in UI
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Notes           types.String    `tfsdk:"notes"` // TODO: is this even supported? doesnt appear in UI
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
    Severity        types.String    `tfsdk:"severity"`
}

type ApplicationControlPolicyRuleAppliactionResourceModel struct {
    Name            types.String    `tfsdk:"name"`
    AllowedVersions types.List      `tfsdk:"allowed_versions"`
}


// Trusted Images Policy Models

type TrustedImagesPolicyResourceModel struct {
    Groups      *[]TrustGroupResourceModel           `tfsdk:"groups"`
    Policy      TrustedImagesPolicyRulesResourceModel    `tfsdk:"policy"`
}

type TrustGroupResourceModel struct {
    Id              types.String    `tfsdk:"id"`
    Images          types.Set       `tfsdk:"images"`
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
}

type TrustedImagesPolicyRulesResourceModel struct {
    Id      types.String                            `tfsdk:"id"`
    Enabled types.Bool                              `tfsdk:"enabled"`
    Rules   *[]TrustedImagesPolicyRuleResourceModel `tfsdk:"rules"`
}

type TrustedImagesPolicyRuleResourceModel struct {
    //Action          types.Set       `tfsdk:"action"`
    AllowedGroups   types.Set       `tfsdk:"allowed_groups"`
    Collections     types.Set       `tfsdk:"collections"`
    DeniedGroups    types.Set       `tfsdk:"denied_groups"`
    Effect          types.String    `tfsdk:"effect"`
    Modified        types.String    `tfsdk:"modified"`
    Name            types.String    `tfsdk:"name"`
    Notes           types.String    `tfsdk:"notes"`
    Owner           types.String    `tfsdk:"owner"`
    PreviousName    types.String    `tfsdk:"previous_name"`
}

// Host Runtime Policy Model

type RuntimeHostPolicyResourceModel struct {
    Rules  *[]RuntimeHostPolicyRuleResourceModel `tfsdk:"rules"`
}

func (m *RuntimeHostPolicyResourceModel) GetRuleNames() []string {
    ruleNames := []string{}

    if (m.Rules == nil || len(*m.Rules) == 0) {
        return ruleNames
    }

    for _, rule := range *m.Rules {
        ruleNames = append(ruleNames, rule.Name.ValueString())
    }

    return ruleNames
}

type RuntimeHostPolicyRuleResourceModel struct {
    AntiMalware       *RuntimeHostPolicyAntiMalwareResourceModel `tfsdk:"anti_malware"`
    Collections       types.Set `tfsdk:"collections"`
    CustomRules       *[]RuntimeHostPolicyCustomRuleResourceModel `tfsdk:"custom_rules"`
    Disabled          types.Bool   `tfsdk:"disabled"`
    FileIntegrityRules *[]RuntimeHostPolicyFileIntegrityRuleResourceModel   `tfsdk:"file_integrity_rules"`
    Activities        *RuntimeHostPolicyActivitiesResourceModel `tfsdk:"activities"`
    LogInspectionRules *[]RuntimeHostPolicyLogInspectionRuleResourceModel `tfsdk:"log_inspection_rules"`
    Modified          types.String `tfsdk:"modified"`
    Name              types.String `tfsdk:"name"`
    Networking        *RuntimeHostPolicyNetworkingResourceModel `tfsdk:"networking"`
    Notes             types.String `tfsdk:"notes"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner             types.String `tfsdk:"owner"`
    PreviousName      types.String `tfsdk:"previous_name"`
}

type RuntimeHostPolicyAntiMalwareResourceModel struct {
    AllowedProcesses          types.List   `tfsdk:"allowed_processes"`
    CryptoMiners              types.String `tfsdk:"crypto_miners"`
    DeniedProcesses           *RuntimePolicyDeniedProcessesResourceModel `tfsdk:"denied_processes"`
    SuppressCompilerGeneratedBinaries types.Bool   `tfsdk:"suppress_compiler_generated_binaries"`
    EncryptedBinaries         types.String `tfsdk:"encrypted_binaries"`
    ExecutionFlowHijacking       types.String `tfsdk:"execution_flow_hijacking"`
    MalwareFromAdvancedThreatProtection types.String `tfsdk:"malware_from_advanced_threat_protection"`
    MalwareFromCustomFeed     types.String `tfsdk:"malware_from_custom_feed"`
    ReverseShell              types.String `tfsdk:"reverse_shell"`
    NonPackagedBinariesService types.String `tfsdk:"non_packaged_binaries_service"`
    NonPackagedBinariesUser   types.String `tfsdk:"non_packaged_binaries_user"`
    SuspiciousELFHeaders      types.String `tfsdk:"suspicious_elf_headers"`
    ProcessesTemporaryStorage types.String `tfsdk:"processes_temporary_storage"`
    WebShell                  types.String `tfsdk:"web_shell"`
    WildFireAnalysis          types.String `tfsdk:"wild_fire_analysis"`
}

type RuntimePolicyDeniedProcessesResourceModel struct {
    Effect types.String `tfsdk:"effect"`
    Paths  types.List `tfsdk:"paths"`
}

type RuntimeHostPolicyCustomRuleResourceModel struct {
    ID     types.Int64  `tfsdk:"id"`
    Name     types.String `tfsdk:"name"`
    LogAs  types.String `tfsdk:"log_as"`
    Effect types.String `tfsdk:"effect"`
}

type RuntimeHostPolicyFileIntegrityRuleResourceModel struct {
    Path          types.String `tfsdk:"file_path"`
    AllowedProcesses types.List `tfsdk:"allowed_processes"`
    ExcludedFilePatterns types.List `tfsdk:"excluded_file_patterns"`
    MonitorSubdirectories types.Bool   `tfsdk:"monitor_subdirectories"`
    MonitorWriteOps   types.Bool   `tfsdk:"monitor_write_ops"`
    MonitorReadOps   types.Bool   `tfsdk:"monitor_read_ops"`
    MonitorMetadataChanges   types.Bool   `tfsdk:"monitor_metadata_changes"`
}

type RuntimeHostPolicyActivitiesResourceModel struct {
    HostActivityMonitoring  RuntimeHostPolicyActivityMonitoringResourceModel `tfsdk:"host_activity_monitoring"`
    TrackSshEvents          types.Bool                                       `tfsdk:"track_ssh_events"`
}

type RuntimeHostPolicyActivityMonitoringResourceModel struct {
    Enabled             types.Bool `tfsdk:"enabled"`
    DockerCommands      RuntimeHostPolicyDockerActivityMonitoringResourceModel `tfsdk:"docker_commands"`
    SshdSessions        types.Bool `tfsdk:"sshd_sessions"`
    SudoCommands        types.Bool `tfsdk:"sudo_commands"`
    LogBackgroundApps   types.Bool `tfsdk:"log_background_apps"`
}

type RuntimeHostPolicyDockerActivityMonitoringResourceModel struct {
    Enabled                 types.Bool `tfsdk:"enabled"`
    IncludeReadOnlyEvents   types.Bool `tfsdk:"include_read_only_events"`
}

type RuntimeHostPolicyLogInspectionRuleResourceModel struct {
    Path  types.String `tfsdk:"path"`
    Regex types.List   `tfsdk:"regex"`
}

type RuntimeHostPolicyNetworkingResourceModel struct {
    AllowedOutboundIPs types.List   `tfsdk:"allowed_outbound_ips"`
    DeniedListeningPorts types.List   `tfsdk:"denied_listening_ports"`
    DeniedOutboundIPs  types.List   `tfsdk:"denied_outbound_ips"`
    DeniedOutboundPorts  types.List   `tfsdk:"denied_outbound_ports"`
    DeniedIPsPortsEffect     types.String `tfsdk:"denied_ips_ports_effect"`
    SuspiciousIPsCustomFeed types.String `tfsdk:"suspicious_ips_custom_feed"`
    SuspiciousIPsAdvancedThreatProtectionEffect types.String `tfsdk:"suspicious_ips_advanced_threat_protection_effect"`
    AllowedDnsDomains types.List   `tfsdk:"allowed_dns_domains"`
    DeniedDnsDomains types.List   `tfsdk:"denied_dns_domains"`
    DeniedDnsDomainsEffect types.String `tfsdk:"denied_dns_domains_effect"`
    SuspiciousDomainsAdvancedThreatProtectionEffect types.String`tfsdk:"suspicious_domains_advanced_threat_protection_effect"`
}

type RuntimeHostPolicyPortRangeResourceModel struct {
    Deny  types.Bool  `tfsdk:"deny"`
    End   types.Int64 `tfsdk:"end"`
    Start types.Int64 `tfsdk:"start"`
}

func (m *RuntimeHostPolicyResourceModel) SortRules(ctx context.Context, planRules *[]RuntimeHostPolicyRuleResourceModel) {
    util.DLog(ctx, "Executing RuntimeHostPolicyResourceModel.SortRules()")

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
    
    util.DLog(ctx, "Finishing RuntimeHostPolicyResourceModel.SortRules() execution")
}

// Container Runtime Policy Model

type RuntimeContainerPolicyResourceModel struct {
    AutomaticRuntimeLearning types.Bool `tfsdk:"automatic_runtime_learning"`
    Rules  *[]RuntimeContainerPolicyRuleResourceModel `tfsdk:"rules"`
}

func (m *RuntimeContainerPolicyResourceModel) GetRuleNames() []string {
    ruleNames := []string{}

    if (m.Rules == nil || len(*m.Rules) == 0) {
        return ruleNames
    }

    for _, rule := range *m.Rules {
        ruleNames = append(ruleNames, rule.Name.ValueString())
    }

    return ruleNames
}

type RuntimeContainerPolicyRuleResourceModel struct {
    AntiMalware       *RuntimeContainerPolicyAntiMalwareResourceModel `tfsdk:"anti_malware"`
    Processes         *RuntimeContainerPolicyProcessesResourceModel `tfsdk:"processes"`
    Networking        *RuntimeContainerPolicyNetworkingResourceModel `tfsdk:"networking"`
    FileSystem        *RuntimeContainerPolicyFileSystemResourceModel `tfsdk:"file_system"`
    CustomRules       *[]RuntimeHostPolicyCustomRuleResourceModel `tfsdk:"custom_rules"`
    Collections       types.Set `tfsdk:"collections"`
    Disabled          types.Bool   `tfsdk:"disabled"`
    Modified          types.String `tfsdk:"modified"`
    Name              types.String `tfsdk:"name"`
    Notes             types.String `tfsdk:"notes"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner             types.String `tfsdk:"owner"`
    PreviousName      types.String `tfsdk:"previous_name"`
}

type RuntimeContainerPolicyAntiMalwareResourceModel struct {
    MalwareFromAdvancedThreatProtection types.String `tfsdk:"malware_from_advanced_threat_protection"`
    KubernetesAttacks                   types.String `tfsdk:"kubernetes_attacks"`
    SuspiciousCloudProviderApiQueries   types.String `tfsdk:"suspicious_cloud_provider_api_queries"`
    WildFireAnalysis                    types.String `tfsdk:"wild_fire_analysis"`
}

type RuntimeContainerPolicyProcessesResourceModel struct {
    Enabled                                 types.Bool `tfsdk:"enabled"`
    AllowedProcesses                        types.List   `tfsdk:"allowed_processes"`
    AllowOnlyLearnedProcessesFromParents    types.Bool `tfsdk:"allow_only_learned_processes_from_parents"`
    AllowAllActivityInAttachedSessions      types.Bool `tfsdk:"allow_all_activity_in_attached_sessions"`
    ProcessesFromModifiedBinaries           types.String `tfsdk:"processes_from_modified_binaries"`
    CryptoMiners                            types.String `tfsdk:"crypto_miners"`
    ReverseShell                            types.String `tfsdk:"reverse_shell"`
    LateralMovementProcesses                types.String `tfsdk:"lateral_movement_processes"`
    ProcessesStartedWithSUID                types.String `tfsdk:"processes_started_with_suid"`
    DeniedProcesses                         *RuntimePolicyDeniedProcessesResourceModel `tfsdk:"denied_processes"`
    AllOtherProcessesEffect                 types.String `tfsdk:"all_other_processes_effect"`
}

type RuntimeContainerPolicyNetworkingResourceModel struct {
    IpConnectivityEnabled types.Bool   `tfsdk:"ip_connectivity_enabled"`
    AllowedListeningPorts types.List   `tfsdk:"allowed_listening_ports"`
    AllowedOutboundInternetPorts types.List   `tfsdk:"allowed_outbound_internet_ports"`
    AllowedOutboundIPs types.List   `tfsdk:"allowed_outbound_ips"`
    NetworkActivityFromModifiedBinaries types.String `tfsdk:"network_activity_from_modified_binaries"`
    PortScanning types.String `tfsdk:"port_scanning"`
    RawSockets types.String `tfsdk:"raw_sockets"`
    DeniedListeningPorts types.List   `tfsdk:"denied_listening_ports"`
    DeniedListeningPortEffect types.String `tfsdk:"denied_listening_ports_effect"`
    DeniedOutboundInternetPorts types.List   `tfsdk:"denied_outbound_internet_ports"`
    DeniedOutboundInternetPortsEffect types.String `tfsdk:"denied_outbound_internet_ports_effect"`
    DeniedOutboundIPs  types.List   `tfsdk:"denied_outbound_ips"`
    DeniedOutboundIPsEffect  types.String `tfsdk:"denied_outbound_ips_effect"`
    AllOtherActivityEffect types.String `tfsdk:"all_other_activity_effect"`
    DnsEnabled types.Bool   `tfsdk:"dns_enabled"`
    AllowedDnsDomains types.List   `tfsdk:"allowed_dns_domains"`
    DeniedDnsDomains types.List   `tfsdk:"denied_dns_domains"`
    DeniedDnsDomainsEffect types.String `tfsdk:"denied_dns_domains_effect"`
    AllOtherDomainsEffect types.String `tfsdk:"all_other_domains_effect"`
}

type RuntimeContainerPolicyFileSystemResourceModel struct {
    Enabled types.Bool `tfsdk:"enabled"`
    AllOtherPathsEffect types.String `tfsdk:"all_other_paths_effect"`
    AllowedPaths types.List   `tfsdk:"allowed_paths"`
    ChangesToBinaries types.String `tfsdk:"changes_to_binaries"`
    ChangesToSshAdminAccountConfigFiles types.String `tfsdk:"changes_to_ssh_admin_account_config_files"`
    DeniedPaths types.List   `tfsdk:"denied_paths"`
    DeniedPathsEffect types.String `tfsdk:"denied_paths_effect"`
    DetectionOfEncryptedBinaries types.String `tfsdk:"detection_of_encrypted_binaries"`
    SuspiciousELFHeaders      types.String `tfsdk:"suspicious_elf_headers"`
}

func (m *RuntimeContainerPolicyResourceModel) SortRules(ctx context.Context, planRules *[]RuntimeContainerPolicyRuleResourceModel) {
    util.DLog(ctx, "Executing RuntimeContainerPolicyResourceModel.SortRules()")

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
    
    util.DLog(ctx, "Finishing RuntimeContainerPolicyResourceModel.SortRules() execution")
}

// Serverless Runtime Policy Model

type RuntimeServerlessPolicyResourceModel struct {
    Rules  *[]RuntimeServerlessPolicyRuleResourceModel `tfsdk:"rules"`
}

func (m *RuntimeServerlessPolicyResourceModel) GetRuleNames() []string {
    ruleNames := []string{}

    if (m.Rules == nil || len(*m.Rules) == 0) {
        return ruleNames
    }

    for _, rule := range *m.Rules {
        ruleNames = append(ruleNames, rule.Name.ValueString())
    }

    return ruleNames
}

type RuntimeServerlessPolicyRuleResourceModel struct {
    Processes         *RuntimeServerlessPolicyProcessesResourceModel `tfsdk:"processes"`
    Networking        *RuntimeServerlessPolicyNetworkingResourceModel `tfsdk:"networking"`
    FileSystem        *RuntimeServerlessPolicyFileSystemResourceModel `tfsdk:"file_system"`
    AdvancedThreatProtection types.Bool `tfsdk:"advanced_threat_protection"`
    Collections       types.Set `tfsdk:"collections"`
    Disabled          types.Bool   `tfsdk:"disabled"`
    Modified          types.String `tfsdk:"modified"`
    Name              types.String `tfsdk:"name"`
    Notes             types.String `tfsdk:"notes"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner             types.String `tfsdk:"owner"`
    PreviousName      types.String `tfsdk:"previous_name"`
}

type RuntimeServerlessPolicyProcessesResourceModel struct {
    Enabled                                 types.Bool `tfsdk:"enabled"`
    AllowedProcesses                        types.List   `tfsdk:"allowed_processes"`
    DeniedProcessesEffect                   types.String  `tfsdk:"denied_processes_effect"`
    CryptoMiners                            types.Bool `tfsdk:"crypto_miners"`
    BlockAllProcessesExceptMain             types.Bool `tfsdk:"block_all_processes_except_main"`
}

type RuntimeServerlessPolicyNetworkingResourceModel struct {
    IpConnectivityEnabled types.Bool   `tfsdk:"ip_connectivity_enabled"`
    AllowedListeningPorts types.List   `tfsdk:"allowed_listening_ports"`
    AllowedOutboundInternetPorts types.List   `tfsdk:"allowed_outbound_internet_ports"`
    AllowedOutboundIPs types.List   `tfsdk:"allowed_outbound_ips"`
    DeniedIPsPortsEffect types.String `tfsdk:"denied_ips_ports_effect"`
    DnsEnabled types.Bool   `tfsdk:"dns_enabled"`
    AllowedDnsDomains types.List   `tfsdk:"allowed_dns_domains"`
    DeniedDnsDomainsEffect types.String `tfsdk:"denied_dns_domains_effect"`
}

type RuntimeServerlessPolicyFileSystemResourceModel struct {
    Enabled types.Bool `tfsdk:"enabled"`
    AllowedPaths types.List   `tfsdk:"allowed_paths"`
    DeniedPaths types.List   `tfsdk:"denied_paths"`
    DeniedPathsEffect types.String `tfsdk:"denied_paths_effect"`
}

func (m *RuntimeServerlessPolicyResourceModel) SortRules(ctx context.Context, planRules *[]RuntimeServerlessPolicyRuleResourceModel) {
    util.DLog(ctx, "Executing RuntimeServerlessPolicyResourceModel.SortRules()")

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

        slices.SortFunc(*m.Rules, func(a, b RuntimeServerlessPolicyRuleResourceModel) int {
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
    
    util.DLog(ctx, "Finishing RuntimeServerlessPolicyResourceModel.SortRules() execution")
}

// App-Embedded Runtime Policy Model

type RuntimeAppEmbeddedPolicyResourceModel struct {
    Rules  *[]RuntimeAppEmbeddedPolicyRuleResourceModel `tfsdk:"rules"`
}

func (m *RuntimeAppEmbeddedPolicyResourceModel) GetRuleNames() []string {
    ruleNames := []string{}

    if (m.Rules == nil || len(*m.Rules) == 0) {
        return ruleNames
    }

    for _, rule := range *m.Rules {
        ruleNames = append(ruleNames, rule.Name.ValueString())
    }

    return ruleNames
}

type RuntimeAppEmbeddedPolicyRuleResourceModel struct {
    Processes         *RuntimeAppEmbeddedPolicyProcessesResourceModel `tfsdk:"processes"`
    Networking        *RuntimeAppEmbeddedPolicyNetworkingResourceModel `tfsdk:"networking"`
    FileSystem        *RuntimeAppEmbeddedPolicyFileSystemResourceModel `tfsdk:"file_system"`
    CustomRules       *[]RuntimeHostPolicyCustomRuleResourceModel `tfsdk:"custom_rules"`
    Collections       types.Set `tfsdk:"collections"`
    Disabled          types.Bool   `tfsdk:"disabled"`
    Modified          types.String `tfsdk:"modified"`
    Name              types.String `tfsdk:"name"`
    Notes             types.String `tfsdk:"notes"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner             types.String `tfsdk:"owner"`
    PreviousName      types.String `tfsdk:"previous_name"`
}

type RuntimeAppEmbeddedPolicyProcessesResourceModel struct {
    Enabled                                 types.Bool `tfsdk:"enabled"`
    AllowedProcesses                        types.List   `tfsdk:"allowed_processes"`
    DeniedProcesses                         types.List   `tfsdk:"denied_processes"`
    DeniedProcessesEffect                   types.String  `tfsdk:"denied_processes_effect"`
    CryptoMiners                            types.Bool `tfsdk:"crypto_miners"`
    ProcessesFromModifiedBinaries           types.Bool `tfsdk:"processes_from_modified_binaries"`
}

type RuntimeAppEmbeddedPolicyNetworkingResourceModel struct {
    IpConnectivityEnabled types.Bool   `tfsdk:"ip_connectivity_enabled"`
    AllowedListeningPorts types.List   `tfsdk:"allowed_listening_ports"`
    AllowedOutboundInternetPorts types.List   `tfsdk:"allowed_outbound_internet_ports"`
    AllowedOutboundIPs types.List   `tfsdk:"allowed_outbound_ips"`
    DeniedIPsPortsEffect types.String `tfsdk:"denied_ips_ports_effect"`
    DeniedListeningPorts types.List   `tfsdk:"denied_listening_ports"`
    DeniedOutboundInternetPorts types.List   `tfsdk:"denied_outbound_internet_ports"`
    DeniedOutboundIPs types.List   `tfsdk:"denied_outbound_ips"`
    DnsEnabled types.Bool   `tfsdk:"dns_enabled"`
    AllowedDnsDomains types.List   `tfsdk:"allowed_dns_domains"`
    DeniedDnsDomainsEffect types.String `tfsdk:"denied_dns_domains_effect"`
}

type RuntimeAppEmbeddedPolicyFileSystemResourceModel struct {
    Enabled types.Bool `tfsdk:"enabled"`
    AllowedPaths types.List   `tfsdk:"allowed_paths"`
    DeniedPaths types.List   `tfsdk:"denied_paths"`
    DeniedPathsEffect types.String `tfsdk:"denied_paths_effect"`
    ChangesToBinariesAndCerts types.Bool `tfsdk:"changes_to_binaries_and_certs"`
    DetectionOfEncryptedBinaries types.Bool `tfsdk:"detection_of_encrypted_binaries"`
    ChangesToSshAdminAccountConfigFiles types.Bool `tfsdk:"changes_to_ssh_admin_account_config_files"`
    SuspiciousELFHeaders      types.Bool `tfsdk:"suspicious_elf_headers"`
    MalwareFromCustomFeed     types.Bool `tfsdk:"malware_from_custom_feed"`
    WildFireAnalysis          types.String `tfsdk:"wild_fire_analysis"`
}

func (m *RuntimeAppEmbeddedPolicyResourceModel) SortRules(ctx context.Context, planRules *[]RuntimeAppEmbeddedPolicyRuleResourceModel) {
    util.DLog(ctx, "Executing RuntimeAppEmbeddedPolicyResourceModel.SortRules()")

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

        slices.SortFunc(*m.Rules, func(a, b RuntimeAppEmbeddedPolicyRuleResourceModel) int {
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
    
    util.DLog(ctx, "Finishing RuntimeAppEmbeddedPolicyResourceModel.SortRules() execution")
}

// Custom Runtime Rule Model

type CustomRuntimeRuleResourceModel struct {
    Id              types.Int64  `tfsdk:"id"`
    AttackTechniques types.List   `tfsdk:"attack_techniques"`
    Description     types.String `tfsdk:"description"`
    Message         types.String `tfsdk:"message"`
    MinVersion      types.String `tfsdk:"min_version"`
    Modified        types.Int64  `tfsdk:"modified"`
    Name            types.String `tfsdk:"name"`
    Owner           types.String `tfsdk:"owner"`
    Script          types.String `tfsdk:"script"`
    Type            types.String `tfsdk:"type"`
    VulnIDs         types.List   `tfsdk:"vuln_ids"`
}
