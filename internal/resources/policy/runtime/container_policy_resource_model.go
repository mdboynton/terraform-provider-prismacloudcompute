package policy

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
    			        		},
    			        		"kubernetes_attacks": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected Kubernetes attacks. Must one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("disable"),
    			        		},
    			        		"suspicious_cloud_provider_api_queries": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected suspicious queries to cloud service provider APIs. Must one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("disable"),
    			        		},
    			        		"wild_fire_analysis": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for detected files classified as malware by WildFire, Palo Alto Networks' malware analysis engine. Must be either \"disable\" or \"alert\". WildFire must be enabled for runtime protection under Manage > System > WildFire.",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
                                ),
                            ),
    			        },
    			        "collections": schema.SetAttribute{
    			        	Optional:    true,
    			        	Computed:    true,
                            Default: setdefault.StaticValue(
                                types.SetValueMust(types.StringType, []attr.Value{types.StringValue("All")}),
                            ),
    			        	ElementType: types.StringType,
                            // TODO: update resource types
                            Description: "List of collection names. Used to scope the rule. Note that in order for a collection to be attached to this type of policy rule, it must contain only the wildcard value (\"*\") for all of the following resource types: Containers, Images, App IDs, Functions, Namespaces, and Clusters.",
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"detection_of_encrypted_binaries": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected encrypted/packed binaries. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"changes_to_ssh_admin_account_config_files": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected changes to SSH or admin account configuration files. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("disable"),
    			        		},
    			        		"all_other_paths_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for all other file system paths not specified in the allow or deny lists. Must be one of \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
                                ),
                            ),
    			        },
    			        "processes": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
                            Description: "Configuration for process monitoring.",
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"crypto_miners": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected crypto miners. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        		"processes_started_with_suid": schema.StringAttribute{
    			        			Optional:    true,
                                    Description: "Effect for detected processes started with superuser ID. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
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
                                                validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                            },
    			        				},
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for all other processes not specified in the allow or deny lists. Must be one of \"disable\", \"alert\", \"prevent\" or \"block\".",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
                                                "effect": types.StringValue("alert"),
                                                "paths": types.ListValueMust(types.StringType, []attr.Value{}),
                                            },
                                        ),
                                        "all_other_processes_effect": types.StringValue("alert"),
                                    },
                                ),
                            ),
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
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"allowed_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "List of outbound internet ports which will not generate alerts or be blocked.",
    			        			ElementType: types.StringType,
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "block"}),
                                    },
    			        		},
    			        		"port_scanning": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for detected port scanning activity. Must be one of \"disable\", \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "block"}),
                                    },
    			        		},
    			        		"raw_sockets": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for detected raw sockets. Must be either \"disable\" or \"alert\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
    			        		},
    			        		"denied_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
                                    Description: "List of listening ports for which access will be alerted on or blocked.",
    			        			ElementType: types.StringType,
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "block"}),
                                    },
    			        		},
    			        		"denied_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
                                    Description: "List of outbound internet ports for which access will be alerted on or blocked.",
    			        			ElementType: types.StringType,
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "block"}),
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "block"}),
                                    },
    			        		},
    			        		"all_other_activity_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for all other network activity for IPs/ports not specified in the allow or deny lists. Must be either \"alert\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "block"}),
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
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent", "block"}),
                                    },
    			        		},
    			        		"all_other_domains_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for all other network activity for domains not specified in the allow or deny lists. Must be one of \"alert\", \"prevent\" or \"block\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent", "block"}),
                                    },
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
                                ),
                            ),
    			        },
    			        "notes": schema.StringAttribute{
    			        	Optional:    true,
    			        	Description: "Notes for the resource",
    			        },
    			        "owner": schema.StringAttribute{
                            Computed: true,
    			        	Description: "Owner of the resource",
    			        },
    			        "previous_name": schema.StringAttribute{
                            Computed: true,
    			        	Description: "Previous name of the resource",
    			        },
                    },
                },
    		},
    	},
    }
}
