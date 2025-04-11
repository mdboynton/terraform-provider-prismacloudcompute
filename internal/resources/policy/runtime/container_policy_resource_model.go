package policy

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ContainerRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &ContainerRuntimePolicyResource{}

func NewContainerRuntimePolicyResource() resource.Resource {
    return &ContainerRuntimePolicyResource{}
}

type ContainerRuntimePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *ContainerRuntimePolicyResource) GetSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
            "automatic_runtime_learning": schema.BoolAttribute{
                //MarkdownDescription: "TODO",
                Optional:            true,
                Computed:            true,
                Description: "Enables automatic behavioural learning.",
                Default:             booldefault.StaticBool(true),
            },
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional:            true,
                Computed:            true,
                PlanModifiers: []planmodifier.List{
                    planmodifiers.UseIndexForUnknownOrder(""),
                    planmodifiers.SetEmptyIfUnknown(),
                },
                Validators: []validator.List{
                    validators.PolicyRuleNameIsUnique(policyAPI.PolicyTypeRuntimeHostFormatted),
                    validators.PolicyRuleOrderIsPositiveNonZero(policyAPI.PolicyTypeRuntimeHostFormatted),
                },
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "order": schema.Int32Attribute{
                            //MarkdownDescription: "TODO",
                            Description: "Order in which the rule will be evaluated. Rules are evaluated from the lowest order value to the highest. When a match is found, the subsequent rules are skipped.",
                            Optional:            true,
                            Computed:            true,
                        },
    			        "anti_malware": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
                            Description: "Configuration for malware monitoring.",
    			        	Attributes: map[string]schema.Attribute{
    			        		"malware_from_advanced_threat_protection": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected files classified as malware by Prisma Cloud Advanced Threat Protection. Must one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["anti_malware.malware_from_advanced_threat_protection"]),
                                    },
    			        		},
    			        		"kubernetes_attacks": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected Kubernetes attacks. Must one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["anti_malware.kubernetes_attacks"]),
                                    },
                                    Default: stringdefault.StaticString("disable"),
    			        		},
    			        		"suspicious_cloud_provider_api_queries": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected suspicious queries to cloud service provider APIs. Must one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["anti_malware.suspicious_cloud_provider_api_queries"]),
                                    },
                                    Default: stringdefault.StaticString("disable"),
    			        		},
    			        		"wild_fire_analysis": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["anti_malware.wild_fire_analysis"]),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for detected files classified as malware by WildFire, Palo Alto Networks' malware analysis engine. Must be either \"disable\" or \"alert\". WildFire must be enabled for runtime protection under Manage > System > WildFire.",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(containerAntiMalwareDefault),
    			        },
    			        "collections": schema.SetAttribute{
    			        	Optional:    true,
    			        	Computed:    true,
                            Default: setdefault.StaticValue(
                                types.SetValueMust(types.StringType, []attr.Value{types.StringValue("All")}),
                            ),
    			        	ElementType: types.StringType,
                            Description: "List of collection names. Used to scope the rule. Note that in order for a collection to be attached to this type of policy rule, it must contain only the wildcard value (\"*\") for all of the following resource types: App IDs and Functions.",
    			        },
    			        "custom_rules": schema.ListNestedAttribute{
                            Optional: true,
                            Description: "List of custom runtime rules.",
                            Validators: []validator.List{
                                validators.CustomRulesAreValid(),
                            },
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "id": schema.Int64Attribute{
                                        Computed:   true,
                                        Description: "Custom rule ID.",
                                    },
                                    "name": schema.StringAttribute{
                                        Optional:   true,
                                        Description: "Name of the custom rule.",
                                    },
                                    "effect": schema.StringAttribute{
                                        Required:   true,
                                        Description: "Effect for the custom rule. Must be one of \"allow\", \"alert\" or \"prevent\". Note that if set to \"allow\", the value of log_as will not have any effect.",
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid([]string{"allow", "alert", "prevent", "block"}),
                                        },
                                    },
                                    "log_as": schema.StringAttribute{
                                        Optional:   true,
    			        			    Computed:   true,
                                        Description: "How violations of this custom rule will be logged. Must be \"audit\" or \"incident\".",
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid([]string{"audit", "incident"}),
                                        },
                                        Default:    stringdefault.StaticString(""),
                                    },
                                },
                            },
                        },
    			        "disabled": schema.BoolAttribute{
    			        	Optional:   true,
                            Computed:   true,
                            Default: booldefault.StaticBool(false),
    			        	Description: "Indicates whether to disable the rule.",
    			        },
    			        "file_system": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
                            Description: "Configuration for file system monitoring.",
    			        	Attributes: map[string]schema.Attribute{
                                "enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables file system activity collection/monitoring.",
                                },
    			        		"allowed_paths": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "List of file system paths which will not be monitored.",
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"changes_to_binaries": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected changes to binary files. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["file_system.changes_to_binaries"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"detection_of_encrypted_binaries": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected encrypted/packed binaries. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["file_system.detection_of_encrypted_binaries"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"changes_to_ssh_admin_account_config_files": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected changes to SSH or admin account configuration files. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["file_system.changes_to_ssh_admin_account_config_files"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"suspicious_elf_headers": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected binaries with suspicious ELF headers. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"denied_paths": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "File system paths to be alerted on or prevented/blocked.",
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"denied_paths_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected file system paths from the deny list. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["file_system.denied_paths_effect"]),
                                    },
                                    Default: stringdefault.StaticString("disable"),
    			        		},
    			        		"all_other_paths_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for all other file system paths not specified in the allow or deny lists. Must be one of \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["file_system.all_other_paths_effect"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(containerFileSystemDefault),
    			        },
    			        "processes": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
                            Description: "Configuration for process monitoring.",
                            Validators: []validator.Object{
                                validators.NoSharedValuesBetweenProcessLists(),
                            },
    			        	Attributes: map[string]schema.Attribute{
                                "enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Description: "Enables process monitoring.",
                                    Default: booldefault.StaticBool(true),
                                },
    			        		"allowed_processes": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
                                    Description: "List of process names to be whitelisted.",
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
                                "allow_only_learned_processes_from_parents": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(false),
                                    Description: "Enable execution of processes only if they are present in the automatic behavioural learning model.",
                                },
                                "allow_all_activity_in_attached_sessions": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Description: "Enable evaluation of processes triggered from attached sessions (docker, kubectl, etc).",
                                    Default: booldefault.StaticBool(false),
                                },
    			        		"processes_from_modified_binaries": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected execution of binaries that do not belong to the original image. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["processes.processes_from_modified_binaries"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"crypto_miners": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected crypto miners. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["processes.crypto_miners"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"reverse_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected reverse shell attacks. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"lateral_movement_processes": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected lateral movement processes. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["processes.lateral_movement_processes"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"processes_started_with_suid": schema.StringAttribute{
    			        			Optional:    true,
                                    Description: "Effect for detected processes started with superuser ID. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["processes.processes_started_with_suid"]),
                                    },
                                    Default: stringdefault.StaticString("disable"),
    			        		},
    			        		"denied_processes": schema.SingleNestedAttribute{
    			        			Optional:    true,
    			        			Description: "Processes to alert on or prevent execution of based on the process name or full path of the binary from which the process is executed.",
    			        			Attributes: map[string]schema.Attribute{
                                        // TODO: validation (cant have values in allowed_processes in paths)
    			        				"effect": schema.StringAttribute{
    			        					Optional:    true,
                                            Description: "Effect for denied processes. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                            Validators: []validator.String{
                                                validators.PolicyEffectIsValid(ValidEffects["container"]["processes.denied_processes.effect"]),
                                            },
    			        				},
                                        // TODO: rename to "processes"
    			        				"paths": schema.ListAttribute{
    			        					Optional:    true,
                                            Computed: true,
    			        					ElementType: types.StringType,
                                            Default: listdefault.StaticValue(
                                                types.ListValueMust(types.StringType, []attr.Value{}),
                                            ),
    			        			        Description: "List of names or full paths of denied processes.",
    			        				},
    			        			},
    			        		},
    			        		"all_other_processes_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["processes.all_other_processes_effect"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for all other processes not specified in the allow or deny lists. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(containerProcessesDefault),
    			        },
    			        "modified": schema.StringAttribute{
    			        	Optional:   true,
                            Computed:   true,
    			        	Description: "Modified timestamp.",
                            PlanModifiers: []planmodifier.String{
                                planmodifiers.UseEmptyStringForNull(),
                            },
    			        },
    			        "name": schema.StringAttribute{
    			        	Required:    true,
    			        	Description: "Name of the policy rule.",
    			        },
    			        "networking": schema.SingleNestedAttribute{
    			        	Optional:    true,
    			        	Computed:    true,
                            Description: "Configuration for network monitoring.",
    			        	Attributes: map[string]schema.Attribute{
                                "ip_connectivity_enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Description: "Enables IP connectivity monitoring.",
                                    Default: booldefault.StaticBool(true),
                                },
    			        		"allowed_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "List of listening ports which will not generate alerts or be blocked.",
    			        			ElementType: types.StringType,
                                    Validators: []validator.List{
                                        validators.PortRangesAreValid(),
                                    },
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"allowed_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "List of outbound internet ports which will not generate alerts or be blocked.",
    			        			ElementType: types.StringType,
                                    Validators: []validator.List{
                                        validators.PortRangesAreValid(),
                                    },
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"allowed_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "List of outbound IPs which will not generate alerts or be blocked.",
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"network_activity_from_modified_binaries": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for activity from modified binaries. Must be one of \"disable\", \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.network_activity_from_modified_binaries"]),
                                    },
    			        		},
    			        		"port_scanning": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for detected port scanning activity. Must be one of \"disable\", \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.port_scanning"]),
                                    },
    			        		},
    			        		"raw_sockets": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for detected raw sockets. Must be either \"disable\" or \"alert\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.raw_sockets"]),
                                    },
    			        		},
    			        		"denied_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
                                    Description: "List of listening ports for which access will be alerted on or blocked.",
    			        			ElementType: types.StringType,
                                    Validators: []validator.List{
                                        validators.PortRangesAreValid(),
                                    },
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"denied_listening_ports_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
    			        			Description: "Effect for denied listening ports. Must be one of \"disable\", \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("disable"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.denied_listening_ports_effect"]),
                                    },
    			        		},
    			        		"denied_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
                                    Description: "List of outbound internet ports for which access will be alerted on or blocked.",
    			        			ElementType: types.StringType,
                                    Validators: []validator.List{
                                        validators.PortRangesAreValid(),
                                    },
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"denied_outbound_internet_ports_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
    			        			Description: "Effect for denied outbound internet ports. Must be one of \"disable\", \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("disable"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.denied_outbound_internet_ports_effect"]),
                                    },
    			        		},
    			        		"denied_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
                                    Description: "List of outbound IPs for which access will be alerted on or blocked.",
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"denied_outbound_ips_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
    			        			Description: "Effect for denied outbound IPs. Must be one of \"disable\", \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("disable"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.denied_outbound_ips_effect"]),
                                    },
    			        		},
    			        		"all_other_activity_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for all other network activity for IPs/ports not specified in the allow or deny lists. Must be either \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.all_other_activity_effect"]),
                                    },
    			        		},
                                "dns_enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Description: "Enables DNS monitoring.",
                                    Default: booldefault.StaticBool(false),
                                },
    			        		"allowed_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "List of DNS domains which will not generate alerts or be prevented/blocked.",
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"denied_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "List of denied DNS domains.",
    			        		},
    			        		"denied_dns_domains_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
    			        			Description: "Effect for DNS domains specified in the deny list. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Default: stringdefault.StaticString("disable"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.denied_dns_domains_effect"]),
                                    },
    			        		},
    			        		"all_other_domains_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for all other network activity for domains not specified in the allow or deny lists. Must be one of \"alert\", \"prevent\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["container"]["networking.all_other_domains_effect"]),
                                    },
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(containerNetworkingDefault),
    			        },
    			        "notes": schema.StringAttribute{
    			        	Optional:    true,
    			        	Description: "Notes for the policy rule.",
    			        },
    			        "owner": schema.StringAttribute{
                            Computed: true,
    			        	Description: "Owner of the policy rule.",
    			        },
    			        "previous_name": schema.StringAttribute{
                            Computed: true,
    			        	Description: "Previous name of the policy rule.",
    			        },
                    },
                },
    		},
    	},
    }
}
