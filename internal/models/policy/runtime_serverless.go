package models

import (
    "context"
    "slices"
    "cmp"
    
    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    "github.com/hashicorp/terraform-plugin-framework/types"
)

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
