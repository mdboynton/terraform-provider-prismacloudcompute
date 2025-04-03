package policy 

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
    // Host policy

    hostAntiMalwareDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "allowed_processes": types.ListType{
                ElemType: types.StringType,
            },
            "crypto_miners": types.StringType,
            "denied_processes": types.ObjectType{
                AttrTypes: map[string]attr.Type{
                    "effect": types.StringType,
                    "paths": types.ListType{
                        ElemType: types.StringType,
                    },
                },
            },
            "encrypted_binaries": types.StringType,
            "execution_flow_hijacking": types.StringType,
            "malware_from_advanced_threat_protection": types.StringType,
            "malware_from_custom_feed": types.StringType,
            "non_packaged_binaries_service": types.StringType,
            "non_packaged_binaries_user": types.StringType,
            "processes_temporary_storage": types.StringType,
            "reverse_shell": types.StringType,
            "suppress_compiler_generated_binaries": types.BoolType,
            "suspicious_elf_headers": types.StringType,
            "web_shell": types.StringType,
            "wild_fire_analysis": types.StringType,
        }, 
        map[string]attr.Value{
            "allowed_processes": types.ListValueMust(types.StringType, []attr.Value{}),
            "crypto_miners": types.StringValue("alert"),
            "denied_processes": types.ObjectValueMust(
                map[string]attr.Type{
                    "effect": types.StringType,
                    "paths": types.ListType{
                        ElemType: types.StringType,
                    },
                },
                map[string]attr.Value{
                    "effect": types.StringValue("alert"),
                    "paths": types.ListValueMust(types.StringType, []attr.Value{}),
                },
            ),
            "encrypted_binaries": types.StringValue("alert"),
            "execution_flow_hijacking": types.StringValue("alert"),
            "malware_from_advanced_threat_protection": types.StringValue("alert"),
            "malware_from_custom_feed": types.StringValue("alert"),
            "non_packaged_binaries_service": types.StringValue("alert"),
            "non_packaged_binaries_user": types.StringValue("alert"),
            "processes_temporary_storage": types.StringValue("alert"),
            "reverse_shell": types.StringValue("alert"),
            "suppress_compiler_generated_binaries": types.BoolValue(false),
            "suspicious_elf_headers": types.StringValue("alert"),
            "web_shell": types.StringValue("alert"),
            "wild_fire_analysis": types.StringValue("alert"),
        },
    )

    hostNetworkingDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "allowed_outbound_ips": types.ListType{
                ElemType: types.StringType,
            },
            "suspicious_ips_custom_feed": types.StringType,
            "denied_listening_ports": types.ListType{
                ElemType: types.StringType,
            },
            "denied_outbound_ips": types.ListType{
                ElemType: types.StringType,
            },
            "denied_outbound_ports": types.ListType{
                ElemType: types.StringType,
            },
            "denied_ips_ports_effect": types.StringType,
            "suspicious_ips_advanced_threat_protection_effect": types.StringType,
            "allowed_dns_domains": types.ListType{
                ElemType: types.StringType,
            },
            "denied_dns_domains": types.ListType{
                ElemType: types.StringType,
            },
            "denied_dns_domains_effect": types.StringType,
            "suspicious_domains_advanced_threat_protection_effect": types.StringType,
        }, 
        map[string]attr.Value{
            "allowed_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
            "suspicious_ips_custom_feed": types.StringValue("alert"),
            "denied_listening_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_outbound_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_ips_ports_effect": types.StringValue("alert"),
            "suspicious_ips_advanced_threat_protection_effect": types.StringValue("alert"),
            "allowed_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_dns_domains_effect": types.StringValue("disable"),
            "suspicious_domains_advanced_threat_protection_effect": types.StringValue("disable"),
        },
    )

    // Container policy

    containerAntiMalwareDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "malware_from_advanced_threat_protection": types.StringType,
            "kubernetes_attacks": types.StringType,
            "suspicious_cloud_provider_api_queries": types.StringType,
            "wild_fire_analysis": types.StringType,
        }, 
        map[string]attr.Value{
            "malware_from_advanced_threat_protection": types.StringValue("alert"),
            "kubernetes_attacks": types.StringValue("disable"),
            "suspicious_cloud_provider_api_queries": types.StringValue("disable"),
            "wild_fire_analysis": types.StringValue("alert"),
        },
    )


    containerFileSystemDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "enabled": types.BoolType,
            "allowed_paths": types.ListType{
                ElemType: types.StringType,
            },
            "changes_to_binaries": types.StringType,
            "detection_of_encrypted_binaries": types.StringType,
            "changes_to_ssh_admin_account_config_files": types.StringType,
            "suspicious_elf_headers": types.StringType,
            "denied_paths": types.ListType{
                ElemType: types.StringType,
            },
            "denied_paths_effect": types.StringType,
            "all_other_paths_effect": types.StringType,
        }, 
        map[string]attr.Value{
            "enabled": types.BoolValue(true),
            "allowed_paths": types.ListValueMust(types.StringType, []attr.Value{}),
            "changes_to_binaries": types.StringValue("alert"),
            "detection_of_encrypted_binaries": types.StringValue("alert"),
            "changes_to_ssh_admin_account_config_files": types.StringValue("alert"),
            "suspicious_elf_headers": types.StringValue("alert"),
            "denied_paths": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_paths_effect": types.StringValue("disable"),
            "all_other_paths_effect": types.StringValue("alert"),
        },
    )

    containerProcessesDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "enabled": types.BoolType,
            "allowed_processes": types.ListType{
                ElemType: types.StringType,
            },
            "allow_only_learned_processes_from_parents": types.BoolType,
            "allow_all_activity_in_attached_sessions": types.BoolType,
            "processes_from_modified_binaries": types.StringType,
            "crypto_miners": types.StringType,
            "reverse_shell": types.StringType,
            "lateral_movement_processes": types.StringType,
            "processes_started_with_suid": types.StringType,
            "denied_processes": types.ObjectType{
                AttrTypes: map[string]attr.Type{
                    "effect": types.StringType,
                    "paths": types.ListType{
                        ElemType: types.StringType,
                    },
                },
            },
            "all_other_processes_effect": types.StringType,

        }, 
        map[string]attr.Value{
            "enabled": types.BoolValue(true),
            "allowed_processes": types.ListValueMust(types.StringType, []attr.Value{}),
            "allow_only_learned_processes_from_parents": types.BoolValue(false),
            "allow_all_activity_in_attached_sessions": types.BoolValue(false),
            "processes_from_modified_binaries": types.StringValue("alert"),
            "crypto_miners": types.StringValue("alert"),
            "reverse_shell": types.StringValue("alert"),
            "lateral_movement_processes": types.StringValue("alert"),
            "processes_started_with_suid": types.StringValue("disable"),
            "denied_processes": types.ObjectValueMust(
                map[string]attr.Type{
                    "effect": types.StringType,
                    "paths": types.ListType{
                        ElemType: types.StringType,
                    },
                },
                map[string]attr.Value{
                    "effect": types.StringValue("disable"),
                    "paths": types.ListValueMust(types.StringType, []attr.Value{}),
                },
            ),
            "all_other_processes_effect": types.StringValue("alert"),
        },
    )

    containerNetworkingDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "ip_connectivity_enabled": types.BoolType,
            "allowed_listening_ports": types.ListType{
                ElemType: types.StringType,
            },
            "allowed_outbound_internet_ports": types.ListType{
                ElemType: types.StringType,
            },
            "allowed_outbound_ips": types.ListType{
                ElemType: types.StringType,
            },
            "network_activity_from_modified_binaries": types.StringType,
            "port_scanning": types.StringType,
            "raw_sockets": types.StringType,
            "denied_listening_ports": types.ListType{
                ElemType: types.StringType,
            },
            "denied_listening_ports_effect": types.StringType,
            "denied_outbound_internet_ports": types.ListType{
                ElemType: types.StringType,
            },
            "denied_outbound_internet_ports_effect": types.StringType,
            "denied_outbound_ips": types.ListType{
                ElemType: types.StringType,
            },
            "denied_outbound_ips_effect": types.StringType,
            "all_other_activity_effect": types.StringType,
            "dns_enabled": types.BoolType,
            "allowed_dns_domains": types.ListType{
                ElemType: types.StringType,
            },
            "denied_dns_domains": types.ListType{
                ElemType: types.StringType,
            },
            "denied_dns_domains_effect": types.StringType,
            "all_other_domains_effect": types.StringType,
        }, 
        map[string]attr.Value{
            "ip_connectivity_enabled": types.BoolValue(true),
            "allowed_listening_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "allowed_outbound_internet_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "allowed_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
            "network_activity_from_modified_binaries": types.StringValue("alert"),
            "port_scanning": types.StringValue("alert"),
            "raw_sockets": types.StringValue("alert"),
            "denied_listening_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_listening_ports_effect": types.StringValue("disable"),
            "denied_outbound_internet_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_outbound_internet_ports_effect": types.StringValue("disable"),
            "denied_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_outbound_ips_effect": types.StringValue("disable"),
            "all_other_activity_effect": types.StringValue("alert"),
            "dns_enabled": types.BoolValue(true),
            "allowed_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_dns_domains_effect": types.StringValue("disable"),
            "all_other_domains_effect": types.StringValue("alert"),
        },
    )

    // Serverless policy

    serverlessFileSystemDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "enabled": types.BoolType,
            "allowed_paths": types.ListType{
                ElemType: types.StringType,
            },
            "denied_paths": types.ListType{
                ElemType: types.StringType,
            },
            "denied_paths_effect": types.StringType,
        }, 
        map[string]attr.Value{
            "enabled": types.BoolValue(true),
            "allowed_paths": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_paths": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_paths_effect": types.StringValue("alert"),
        },
    )

    serverlessProcessesDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "enabled": types.BoolType,
            "allowed_processes": types.ListType{
                ElemType: types.StringType,
            },
            "denied_processes_effect": types.StringType,
            "crypto_miners": types.BoolType,
            "block_all_processes_except_main": types.BoolType,
        }, 
        map[string]attr.Value{
            "enabled": types.BoolValue(true),
            "allowed_processes": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_processes_effect": types.StringValue("alert"),
            "crypto_miners": types.BoolValue(true),
            "block_all_processes_except_main": types.BoolValue(true),
        },
    )

    serverlessNetworkingDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "ip_connectivity_enabled": types.BoolType,
            "allowed_listening_ports": types.ListType{
                ElemType: types.StringType,
            },
            "allowed_outbound_internet_ports": types.ListType{
                ElemType: types.StringType,
            },
            "allowed_outbound_ips": types.ListType{
                ElemType: types.StringType,
            },
            "denied_ips_ports_effect": types.StringType,
            "dns_enabled": types.BoolType,
            "allowed_dns_domains": types.ListType{
                ElemType: types.StringType,
            },
            "denied_dns_domains_effect": types.StringType,
        }, 
        map[string]attr.Value{
            "ip_connectivity_enabled": types.BoolValue(true),
            "allowed_listening_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "allowed_outbound_internet_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "allowed_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_ips_ports_effect": types.StringValue("alert"),
            "dns_enabled": types.BoolValue(true),
            "allowed_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_dns_domains_effect": types.StringValue("alert"),
        },
    )

    // App-embedded policy

    appEmbeddedFileSystemDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "enabled": types.BoolType,
            "allowed_paths": types.ListType{
                ElemType: types.StringType,
            },
            "denied_paths": types.ListType{
                ElemType: types.StringType,
            },
            "denied_paths_effect": types.StringType,
            "changes_to_binaries_and_certs": types.BoolType,
            "detection_of_encrypted_binaries": types.BoolType,
            "changes_to_ssh_admin_account_config_files": types.BoolType,
            "suspicious_elf_headers": types.BoolType,
            "malware_from_custom_feed": types.BoolType,
            "wild_fire_analysis": types.StringType,
        }, 
        map[string]attr.Value{
            "enabled": types.BoolValue(true),
            "allowed_paths": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_paths": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_paths_effect": types.StringValue("alert"),
            "changes_to_binaries_and_certs": types.BoolValue(true),
            "detection_of_encrypted_binaries": types.BoolValue(true),
            "changes_to_ssh_admin_account_config_files": types.BoolValue(true),
            "suspicious_elf_headers": types.BoolValue(true),
            "malware_from_custom_feed": types.BoolValue(true),
            "wild_fire_analysis": types.StringValue("alert"),
        },
    )
    
    appEmbeddedProcessesDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "enabled": types.BoolType,
            "allowed_processes": types.ListType{
                ElemType: types.StringType,
            },
            "denied_processes": types.ListType{
                ElemType: types.StringType,
            },
            "denied_processes_effect": types.StringType,
            "crypto_miners": types.BoolType,
            "processes_from_modified_binaries": types.BoolType,
        }, 
        map[string]attr.Value{
            "enabled": types.BoolValue(true),
            "allowed_processes": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_processes": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_processes_effect": types.StringValue("alert"),
            "crypto_miners": types.BoolValue(true),
            "processes_from_modified_binaries": types.BoolValue(true),
        },
    )

    appEmbeddedNetworkingDefault = types.ObjectValueMust(
        map[string]attr.Type{
            "ip_connectivity_enabled": types.BoolType,
            "allowed_listening_ports": types.ListType{
                ElemType: types.StringType,
            },
            "allowed_outbound_internet_ports": types.ListType{
                ElemType: types.StringType,
            },
            "allowed_outbound_ips": types.ListType{
                ElemType: types.StringType,
            },
            "denied_ips_ports_effect": types.StringType,
            "denied_listening_ports": types.ListType{
                ElemType: types.StringType,
            },
            "denied_outbound_internet_ports": types.ListType{
                ElemType: types.StringType,
            },
            "denied_outbound_ips": types.ListType{
                ElemType: types.StringType,
            },
            "dns_enabled": types.BoolType,
            "allowed_dns_domains": types.ListType{
                ElemType: types.StringType,
            },
            "denied_dns_domains_effect": types.StringType,
        }, 
        map[string]attr.Value{
            "ip_connectivity_enabled": types.BoolValue(true),
            "allowed_listening_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "allowed_outbound_internet_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "allowed_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_ips_ports_effect": types.StringValue("alert"),
            "denied_listening_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_outbound_internet_ports": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
            "dns_enabled": types.BoolValue(true),
            "allowed_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
            "denied_dns_domains_effect": types.StringValue("alert"),
        },
    )
)
