package policy

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &HostRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &HostRuntimePolicyResource{}

func NewHostRuntimePolicyResource() resource.Resource {
    return &HostRuntimePolicyResource{}
}

type HostRuntimePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *HostRuntimePolicyResource) GetSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional:            true,
                Computed:            true,
                PlanModifiers: []planmodifier.List{
                    planmodifiers.UseIndexForUnknownOrder(""),
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
    			        		"allowed_processes": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "Processes marked as safe to use based on the process name or full path of the binary from which the process is executed. Processes added to this list will not be alerted on or prevented by any of the malware runtime capabilities.",
    			        		},
    			        		"crypto_miners": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for detected crypto miners. Must be one of \"disable\", \"alert\", or \"prevent\". Note that when setting to \"prevent\", only some detected use cases will be prevented. Others will only generate an alert.",
    			        		},
    			        		"denied_processes": schema.SingleNestedAttribute{
    			        			Optional:    true,
    			        			Description: "Processes to alert on or prevent execution of based on the process name or full path of the binary from which the process is executed.",
    			        			Attributes: map[string]schema.Attribute{
    			        				"effect": schema.StringAttribute{
    			        					Optional:    true,
    			        					Description: "Effect for denied processes. Must be either \"alert\" or \"prevent\".",
                                            Validators: []validator.String{
                                                validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
                                            },
    			        				},
    			        				"paths": schema.ListAttribute{
    			        					Optional:    true,
    			        					ElementType: types.StringType,
    			        			        Description: "List of names or full paths of denied processes.",
    			        				},
    			        			},
    			        		},
    			        		"encrypted_binaries": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected encrypted/packed binaries. Must be either \"disable\" or \"alert\".",
    			        		},
    			        		"execution_flow_hijacking": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected execution flow hijack attempts. Must be either \"disable\" or \"alert\".",
    			        		},
    			        		"malware_from_advanced_threat_protection": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for detected files classified as malware by Prisma Cloud Advanced Threat Protection. Must be either \"disable\" or \"alert\".",
    			        		},
    			        		"malware_from_custom_feed": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected files classified as malware by the presence of their MD5 hash in the custom malware signature feed. Must be either \"disable\" or \"alert\". Custom malware signatures can be configured in the console by navigating to Manage > System > Custom feeds > Malware signatures.",
    			        		},
    			        		"non_packaged_binaries_service": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected binaries created or executed by a service without a package manager. Must be one of \"disable\", \"alert\", or \"prevent\". Defender must be running when a file is written to detect its source. Note that when setting to \"prevent\", only file execution will be prevented, while alerts will be generated on file creation.",
    			        		},
    			        		"non_packaged_binaries_user": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected binaries created or executed by a user without a package manager. Must be one of \"disable\", \"alert\", or \"prevent\". Defender must be running when a file is written to detect its source. Note that when setting to \"prevent\", only file execution will be prevented, while alerts will be generated on file creation.",
    			        		},
    			        		"processes_temporary_storage": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected processes executed from temporary storage. Must be one of \"disable\", \"alert\", or \"prevent\".",
    			        		},
    			        		"reverse_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected reverse shell attacks. Must be either \"disable\" or \"alert\".",
    			        		},
    			        		"suppress_compiler_generated_binaries": schema.BoolAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Default:    booldefault.StaticBool(false),
                                    Description: "Toggle suppression of alerts created by detected non-packaged binaries executed by compiler services. This setting has no effect if non_packaged_binaries_service is set to \"disable\".",
                                    // TODO: do not allow if non_packaged_binaries_run_by_service is set to "disable"
    			        		},
    			        		"suspicious_elf_headers": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
                                    Description: "Effect for detected suspicious ELF headers. Must be either \"disable\" or \"alert\".",
    			        		},
    			        		"web_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for detected web shell attacks. Must be one of \"disable\", \"alert\" or \"prevent\". Note that when setting to \"prevent\", only Linux command line tool execution will be prevented. Alerts will be generated on web shell creation.",
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
                            Default: objectdefault.StaticValue(hostAntiMalwareDefault),
    			        },
    			        "collections": schema.SetAttribute{
    			        	Optional:    true,
    			        	ElementType: types.StringType,
                            Description: "List of collection names. Used to scope the rule. Note that in order for a collection to be attached to this type of policy rule, it must contain only the wildcard value (\"*\") for all of the following resource types: Containers, Images, App IDs, Functions, Namespaces and Clusters.",
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
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid([]string{"allow", "alert", "prevent"}),
                                        },
                                        Description: "Effect for the custom rule. Must be one of \"allow\", \"alert\" or \"prevent\". Note that if set to \"allow\", the value of log_as will not have any effect.",
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
    			        "file_integrity_rules": schema.ListNestedAttribute{
    			        	Optional:    true,
                            Description: "List of file integrity rules. Each rule must have at least one of the monitoring options enabled.",
                            Validators: []validator.List{
                                validators.FileIntegrityRulesAreValid(),
                            },
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    // TODO: add logic to set "dir" to "true" in API call if the provided file_path value is a directory
                                    // TODO: update FileIntegrityRulesAreValid to only apply the validation rules if the specified path is a directory
                                    "file_path": schema.StringAttribute{
                                        Optional:   true,
                                        Description: "File/directory path to monitor. Values must be unique.",
                                    },
                                    "allowed_processes": schema.ListAttribute{
                                        Optional:   true,
                                        ElementType: types.StringType,
                                        Description: "List of process names to be whitelisted.",
                                    },
                                    "excluded_file_patterns": schema.ListAttribute{
                                        Optional:   true,
                                        ElementType: types.StringType,
                                        Description: "Filename patterns to exclude from monitoring. The wildcard character \"*\" can be used in the filename only (e.g. \"foo*.log\", \"*.cache\").",
                                    },
                                    "monitor_subdirectories": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "Enables monitoring of subdirectories. If enabled for a directory, monitor_write_ops must also be set to true and both monitor_read_ops and monitor_metadata_changes must be set to false.",
                                    },
                                    "monitor_write_ops": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "Enables monitoring of write operations.",
                                    },
                                    "monitor_read_ops": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "Enables monitoring of read operations.",
                                    },
                                    "monitor_metadata_changes": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "Enables monitoring of metadata changes (e.g. chmod, chown).",
                                    },
                                },
                            },
    			        },
    			        "activities": schema.SingleNestedAttribute{
    			        	Optional: true,
    			        	Attributes: map[string]schema.Attribute{
                                "host_activity_monitoring": schema.SingleNestedAttribute{
                                    Optional: true,
                                    Attributes: map[string]schema.Attribute{
                                        "enabled": schema.BoolAttribute{
                                            Optional: true,
                                            Description: "Enables host activity collection/monitoring.",
                                        },
                                        "docker_commands": schema.SingleNestedAttribute{
                                            Optional: true,
                                            Attributes: map[string]schema.Attribute{
                                                "enabled": schema.BoolAttribute{
                                                    Optional: true,
                                                    Description: "Enables monitoring of docker commands.",
                                                },
                                                "include_read_only_events": schema.BoolAttribute{
                                                    Optional: true,
                                                    Description: "Enables monitoring of read-only Docker events.",
                                                },
                                            },
                                        },
                                        "sshd_sessions": schema.BoolAttribute{
                                            Optional: true,
                                            Description: "Enables monitoring of activity data from new sessions spawned by sshd.",
                                        },
                                        "sudo_commands": schema.BoolAttribute{
                                            Optional: true,
                                            Description: "Enables monitoring of activity data from commands executed with sudo or su.",
                                        },
                                        "log_background_apps": schema.BoolAttribute{
                                            Optional: true,
                                            Description: "Enables monitoring of activity data from background applications. Note that this will result in additional performance overhead and significant data being logged.",
                                        },
                                    },
                                },
                                "track_ssh_events": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(false),
                                    Description: "Enables monitoring of SSH events.",
                                },
    			        	},
    			        },
    			        "log_inspection_rules": schema.ListAttribute{
    			        	Optional:    true,
                            Validators: []validator.List{
                                validators.LogInspectionPathIsValid(),
                            },
    			        	ElementType: types.ObjectType{
    			        		AttrTypes: map[string]attr.Type{
    			        			"path": types.StringType,
    			        			"regex": types.ListType{
    			        				ElemType: types.StringType,
    			        			},
    			        		},
    			        	},
    			        	Description: "List of log inspection rules. Path value must be non-empty, absolute (begins with \"/\"), and unique. Asterisks may only be used in the file name portion of the path.",
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
    			        		"allowed_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of allowed outbound IPs which will not generate alerts.",
    			        		},
    			        		"suspicious_ips_custom_feed": schema.StringAttribute{
    			        			Optional:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
    			        			Description: "Effect for detected IPs from the list of suspicious or high risk IPs under Manage > System > Custom feeds > IP Reputation lists. Must be either \"disable\" or \"alert\".",
    			        		},
    			        		"denied_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of listening ports for which accessing will generate alerts.",
    			        		},
    			        		"denied_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of outbound IPs for which accessing will generate alerts.",
    			        		},
    			        		"denied_outbound_ports": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of outbound ports for which accessing will generate alerts.",
    			        		},
    			        		"denied_ips_ports_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
    			        			Description: "Effect for denied IP/ports. Must be either \"disable\" or \"alert\".",
    			        		},
    			        		"suspicious_ips_advanced_threat_protection_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert"}),
                                    },
    			        			Description: "Effect for detected malicious IPs based on the Prisma Cloud advanced threat protection intelligence stream.",
                                    // TODO: description
    			        		},
    			        		"allowed_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of DNS domains which will not generate alerts.",
                                    // TODO: validation (no duplicates)
    			        		},
    			        		"denied_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
                                    Description: "List of DNS domains for which access will generate an alert or be prevented.",
                                    // TODO: validation (no duplicates)
    			        		},
    			        		"denied_dns_domains_effect": schema.StringAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("disable"),
    			        			Description: "Effect for denied DNS domains. Must be one of \"disable\", \"alert\" or \"prevent\".",
    			        		},
    			        		"suspicious_domains_advanced_threat_protection_effect": schema.StringAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("disable"),
    			        			Description: "Effect for detected malicious domains based on the Prisma Cloud advanced threat protection intelligence stream.",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(hostNetworkingDefault),
                            //Default: objectdefault.StaticValue(
                            //    types.ObjectValueMust(
                            //        map[string]attr.Type{
                            //            "allowed_outbound_ips": types.ListType{
                            //                ElemType: types.StringType,
                            //            },
                            //            "suspicious_ips_custom_feed": types.StringType,
                            //            "denied_listening_ports": types.ListType{
                            //                ElemType: types.StringType,
                            //            },
                            //            "denied_outbound_ips": types.ListType{
                            //                ElemType: types.StringType,
                            //            },
                            //            "denied_outbound_ports": types.ListType{
                            //                ElemType: types.StringType,
                            //            },
                            //            "denied_ips_ports_effect": types.StringType,
                            //            "suspicious_ips_advanced_threat_protection_effect": types.StringType,
                            //            "allowed_dns_domains": types.ListType{
                            //                ElemType: types.StringType,
                            //            },
                            //            "denied_dns_domains": types.ListType{
                            //                ElemType: types.StringType,
                            //            },
                            //            "denied_dns_domains_effect": types.StringType,
                            //            "suspicious_domains_advanced_threat_protection_effect": types.StringType,
                            //        }, 
                            //        map[string]attr.Value{
                            //            "allowed_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
                            //            "suspicious_ips_custom_feed": types.StringValue("alert"),
                            //            "denied_listening_ports": types.ListValueMust(types.StringType, []attr.Value{}),
                            //            "denied_outbound_ips": types.ListValueMust(types.StringType, []attr.Value{}),
                            //            "denied_outbound_ports": types.ListValueMust(types.StringType, []attr.Value{}),
                            //            "denied_ips_ports_effect": types.StringValue("alert"),
                            //            "suspicious_ips_advanced_threat_protection_effect": types.StringValue("alert"),
                            //            "allowed_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
                            //            "denied_dns_domains": types.ListValueMust(types.StringType, []attr.Value{}),
                            //            "denied_dns_domains_effect": types.StringValue("disable"),
                            //            "suspicious_domains_advanced_threat_protection_effect": types.StringValue("disable"),
                            //        },
                            //    ),
                            //),
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
