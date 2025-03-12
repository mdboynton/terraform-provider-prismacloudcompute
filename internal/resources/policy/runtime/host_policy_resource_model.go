package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"

    "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ resource.Resource = &HostRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &HostRuntimePolicyResource{}
//var _ resource.ResourceWithModifyPlan = &HostRuntimePolicyResource{}

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
                //PlanModifiers: []planmodifier.List{
                //    //planmodifiers.UseIndexForUnknownOrder(policyAPI.PolicyTypeComplianceHostFormatted),
                //    //planmodifiers.SetNullByModuleType(metaType),
                //    planmodifiers.UseIndexForUnknownOrder(policyTypeFormatted),
                //    //planmodifiers.GenerateConditionFromEffect(policyType, complianceVulnerabilities),
                //},
                //Validators: []validator.List{
                //    //validators.PolicyRuleNameIsUnique(policyAPI.PolicyTypeComplianceHostFormatted),
                //    //validators.PolicyRuleOrderIsPositiveNonZero(policyAPI.PolicyTypeComplianceHostFormatted),
                //    validators.PolicyRuleNameIsUnique(policyTypeFormatted),
                //    validators.PolicyRuleOrderIsPositiveNonZero(policyTypeFormatted),
                //},
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "order": schema.Int32Attribute{
                            MarkdownDescription: "TODO",
                            Optional:            true,
                            Computed:            true,
                        },
    			        "anti_malware": schema.SingleNestedAttribute{
    			        	Optional: true,
    			        	Attributes: map[string]schema.Attribute{
    			        		"allowed_processes": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of allowed processes",
    			        		},
    			        		"crypto_miners": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
    			        			Description: "Effect for crypto miner detection",
    			        		},
    			        		"denied_processes": schema.SingleNestedAttribute{
    			        			Optional:    true,
    			        			Description: "Denied Processes configuration",
    			        			Attributes: map[string]schema.Attribute{
    			        				"effect": schema.StringAttribute{
    			        					Optional:    true,
    			        					Description: "Effect for denied processes",
                                            // TODO: validation
    			        				},
    			        				"paths": schema.ListAttribute{
    			        					Optional:    true,
    			        					ElementType: types.StringType,
    			        					Description: "List of paths for denied processes",
    			        				},
    			        			},
    			        		},
    			        		"suppress_compiler_generated_binaries": schema.BoolAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Default:    booldefault.StaticBool(false),
    			        			Description: "Detect compiler generated binary.",
                                    // TODO: update Description
                                    // TODO: do not allow if non_packaged_binaries_run_by_service is set to "disable"
    			        		},
    			        		"encrypted_binaries": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for encrypted/packed binaries.",
    			        		},
    			        		"execution_flow_hijacking": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for execution flow hijacking.",
    			        		},
    			        		"malware_from_advanced_threat_protection": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for intelligence feed.",
    			        		},
    			        		"malware_from_custom_feed": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for malware based on custom feed",
    			        		},
    			        		"reverse_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for reverse shell attacks.",
    			        		},
    			        		"non_packaged_binaries_service": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for service unknown origin binary.",
                                    // TODO: description
    			        		},
    			        		"non_packaged_binaries_user": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for user unknown origin binary.",
    			        		},
    			        		"suspicious_elf_headers": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for binaries with suspicious ELF headers.",
    			        		},
    			        		"processes_temporary_storage": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for processes running from temporary storage.",
    			        		},
    			        		"web_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for web shell attacks. Note that when setting this effect to \"Prevent\", the \"Prevent\" effect will only apply to file execution. Alerts will be generated on file creation.",
    			        		},
    			        		"wild_fire_analysis": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for WildFire analysis. WildFire must be enabled for runtime protection under Manage > System > WildFire.",
    			        		},
    			        	},
    			        },
    			        "collections": schema.SetAttribute{
    			        	Optional:    true,
    			        	ElementType: types.StringType,
    			        	Description: "List of collections.",
    			        },
    			        "custom_rules": schema.ListNestedAttribute{
                            Optional: true,
                            Description: "List of custom rules.",
                            NestedObject: schema.NestedAttributeObject{
                            // TODO: validation (cant have prevent effect with incident log_as)
                            // TODO: validation (log_as cant be specified with allow effect)
                                Attributes: map[string]schema.Attribute{
                                    "id": schema.Int64Attribute{
                                        Computed:   true,
                                        Description: "",
                                    },
                                    "name": schema.StringAttribute{
                                        Optional:   true,
                                        Description: "",
                                    },
                                    "effect": schema.StringAttribute{
                                        Required:   true,
                                        Description: "",
                                    },
                                    "log_as": schema.StringAttribute{
                                        Optional:   true,
    			        			    Computed:   true,
                                        Default:    stringdefault.StaticString(""),
                                        Description: "",
                                    },
                                },
                            },
                        },
    			        "disabled": schema.BoolAttribute{
    			        	Optional:   true,
                            Computed:   true,
                            Default: booldefault.StaticBool(false),
    			        	Description: "Disable the rule.",
    			        },
    			        "file_integrity_rules": schema.ListNestedAttribute{
    			        	Optional:    true,
                            Description: "List of file integrity rules",
                            // TODO: validation (see below)
                            // TODO: ensure at least one monitor argument provided
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "file_path": schema.StringAttribute{
                                        Optional:   true,
                                        Description: "File path.",
                                    },
                                    "allowed_processes": schema.ListAttribute{
                                        Optional:   true,
                                        ElementType: types.StringType,
                                        Description: "Process names to be whitelisted.",
                                    },
                                    "excluded_file_patterns": schema.ListAttribute{
                                        Optional:   true,
                                        ElementType: types.StringType,
                                        Description: "File patterns to exclude.",
                                    },
                                    /*
                                        when monitor_subdirectories is enabled:
                                            monitor_write_ops must be enabled?
                                            monitor_read_ops is disabled
                                            monitor_metadata_changes is disabled
                                        when monitor_read_ops is enabled:
                                            monitor_subdirectories is disabled
                                        when monitor_metadata_changes is enabled:
                                            monitor_subdirectories is disabled
                                    */
                                    "monitor_subdirectories": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "TODO",
                                    },
                                    "monitor_write_ops": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "TODO",
                                    },
                                    "monitor_read_ops": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "TODO",
                                    },
                                    "monitor_metadata_changes": schema.BoolAttribute{
                                        Optional:   true,
                                        Computed:   true,
                                        Default: booldefault.StaticBool(false),
                                        Description: "TODO",
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
                                            Description: "Enables host activity monitoring.",
                                        },
                                        "docker_commands": schema.SingleNestedAttribute{
                                            Optional: true,
                                            Attributes: map[string]schema.Attribute{
                                                "enabled": schema.BoolAttribute{
                                                    Optional: true,
                                                    Description: "Collect data from docker commands.",
                                                },
                                                "include_read_only_events": schema.BoolAttribute{
                                                    Optional: true,
                                                    Description: "Toggle collection of read-only Docker events.",
                                                },
                                            },
                                        },
                                        "sshd_sessions": schema.BoolAttribute{
                                            Optional: true,
                                            Description: "Collect activity data from new sessions spawned by sshd.",
                                        },
                                        "sudo_commands": schema.BoolAttribute{
                                            Optional: true,
                                            Description: "Collect activity data from commands executed with sudo or su.",
                                        },
                                        "log_background_apps": schema.BoolAttribute{
                                            Optional: true,
                                            Description: "Collect activity data from background applications. Note that this will result in additional performance overhead and significant data being logged.",
                                        },
                                    },
                                },
                                "track_ssh_events": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(false),
                                    Description: "Enables host activity monitoring.",
                                },
    			        	},
    			        },
    			        "log_inspection_rules": schema.ListAttribute{
    			        	Optional:    true,
    			        	ElementType: types.ObjectType{
    			        		AttrTypes: map[string]attr.Type{
    			        			"path": types.StringType,
    			        			"regex": types.ListType{
    			        				ElemType: types.StringType,
    			        			},
    			        		},
    			        	},
    			        	Description: "List of log inspection rules",
                            // TODO: validate path (must be absolute, only use asterisks in filename)
                            // TODO: enforce unique paths
                            // TODO: enfoce 
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
    			        	Description: "Name of the resource",
    			        },
    			        "networking": schema.SingleNestedAttribute{
    			        	Optional:    true,
    			        	Attributes: map[string]schema.Attribute{
    			        		"allowed_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of allowed outbound IPs.",
    			        		},
    			        		"suspicious_ips_custom_feed": schema.StringAttribute{
    			        			Optional:    true,
    			        			Description: "Effect for custom feed",
                                    // TODO: description
    			        		},
    			        		"denied_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of denied listening ports",
                                    // TODO: change to list of ints
    			        		},
    			        		"denied_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of denied outbound IPs",
    			        		},
    			        		"denied_outbound_ports": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of denied outbound ports",
                                    // TODO: change to list of ints
    			        		},
    			        		"denied_ips_ports_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Description: "Effect for the IP/ports deny list",
                                    // TODO: description
    			        		},
    			        		"suspicious_ips_advanced_threat_protection_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Description: "Effect for intelligence feed",
                                    // TODO: description
    			        		},
    			        		"allowed_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of allowed DNS domains.",
                                    // TODO: validation
    			        		},
    			        		"denied_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of denied DNS domains.",
                                    // TODO: validation
    			        		},
    			        		"denied_dns_domains_effect": schema.StringAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Default:    stringdefault.StaticString("disable"),
    			        			Description: "Effect for the DNS deny list.",
                                    // TODO: description
                                    // TODO: default
    			        		},
    			        		"suspicious_domains_advanced_threat_protection_effect": schema.StringAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Default:    stringdefault.StaticString("disable"),
    			        			Description: "Effect for the intelligence feed",
                                    // TODO: description
                                    // TODO: validation
    			        		},
    			        	},
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
