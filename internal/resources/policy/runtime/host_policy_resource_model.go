package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"

    "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
    			        		//"crypto_miner": schema.ListAttribute{
    			        		"crypto_miners": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
    			        			Description: "Effect for crypto miner detection",
    			        		},
    			        		//"alert_processes": schema.SingleNestedAttribute{
    			        		"denied_processes": schema.SingleNestedAttribute{
    			        			Optional:    true,
    			        			Description: "Denied Processes configuration",
    			        			Attributes: map[string]schema.Attribute{
    			        				//"effect": schema.ListAttribute{
    			        				//	Optional:    true,
    			        				//	ElementType: types.StringType,
    			        				//	Description: "Effect for denied processes",
    			        				//},
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
    			        		//"detect_compiler_generated_binary": schema.BoolAttribute{
    			        		"suppress_compiler_generated_binaries": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Description: "Detect compiler generated binary.",
                                    // TODO: default to "true"
                                    // TODO: update Description
                                    // TODO: do not allow if non_packaged_binaries_run_by_service is set to "disable"
    			        		},
    			        		"encrypted_binaries": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for encrypted/packed binaries.",
    			        		},
    			        		"execution_flow_hijack": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for execution flow hijacking.",
    			        		},
    			        		//"intelligence_feed": schema.ListAttribute{
    			        		"malware_from_advanced_threat_protection": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for intelligence feed.",
    			        		},
    			        		//"custom_feed": schema.ListAttribute{
    			        		//	Optional:    true,
    			        		//	ElementType: types.StringType,
    			        		//	Description: "Effect for custom feed",
    			        		//},
    			        		"malware_from_custom_feed": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for malware based on custom feed",
                                    // TODO: default to Alert
    			        		},
    			        		//"reverse_shell": schema.ListAttribute{
    			        		"reverse_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for reverse shell attacks.",
    			        		},
    			        		//"service_unknown_origin_binary": schema.ListAttribute{
    			        		"non_packaged_binaries_service": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for service unknown origin binary.",
                                    // TODO: description
    			        		},
    			        		//"user_unknown_origin_binary": schema.ListAttribute{
    			        		"non_packaged_binaries_user": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for user unknown origin binary.",
    			        		},
    			        		//"skip_ssh_tracking": schema.BoolAttribute{
    			        		//	Optional:    true,
    			        		//	Description: "Skip SSH tracking",
    			        		//},
    			        		//"suspicious_elf_headers": schema.ListAttribute{
    			        		"suspicious_elf_headers": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for binaries with suspicious ELF headers.",
    			        		},
    			        		//"temp_fs_proc": schema.ListAttribute{
    			        		"processes_temporary_storage": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for processes running from temporary storage.",
    			        		},
    			        		//"web_shell": schema.ListAttribute{
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
    			        "custom_rules": schema.ListAttribute{
    			        	Optional:    true,
    			        	ElementType: types.ObjectType{
    			        		AttrTypes: map[string]attr.Type{
    			        			"_id": types.Int32Type,
    			        			"action": types.ListType{
    			        				ElemType: types.StringType,
    			        			},
    			        			"effect": types.ListType{
    			        				ElemType: types.StringType,
    			        			},
    			        		},
    			        	},
    			        	Description: "List of custom rules",
    			        },
    			        "disabled": schema.BoolAttribute{
    			        	Optional:    true,
    			        	Description: "Disable the resource",
    			        },
    			        //"file_integrity_rules": schema.ListAttribute{
    			        //	Optional:    true,
    			        //	ElementType: types.ObjectType{
    			        //		AttrTypes: map[string]attr.Type{
    			        //			"dir": types.BoolType,
    			        //			"exclusions": types.ListType{
    			        //				ElemType: types.StringType,
    			        //			},
    			        //			"metadata": types.BoolType,
    			        //			"path": types.StringType,
    			        //			"proc_whitelist": types.ListType{
    			        //				ElemType: types.StringType,
    			        //			},
    			        //			"read": types.BoolType,
    			        //			"recursive": types.BoolType,
    			        //			"write": types.BoolType,
    			        //		},
    			        //	},
    			        //	Description: "List of file integrity rules",
    			        //},
    			        "file_integrity_rules": schema.ListNestedAttribute{
    			        	Optional:    true,
                            Description: "List of file integrity rules",
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
                                            monitor_read_ops is disabled
                                            monitor_metadata_changes is disabled
                                        when monitor_read_ops is enabled:
                                            monitor_subdirectories is disabled
                                        when monitor_metadata_changes is enabled:
                                            monitor_subdirectories is disabled
                                    */
                                    "monitor_subdirectories": schema.BoolAttribute{
                                        Optional:   true,
                                        Description: "TODO",
                                    },
                                    "monitor_write_ops": schema.BoolAttribute{
                                        Optional:   true,
                                        Description: "TODO",
                                    },
                                    "monitor_read_ops": schema.BoolAttribute{
                                        Optional:   true,
                                        Description: "TODO",
                                    },
                                    "monitor_metadata_changes": schema.BoolAttribute{
                                        Optional:   true,
                                        Description: "TODO",
                                    },
                                },
                            },
    			        },
    			        //"forensic": schema.SingleNestedAttribute{
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
    			        },
    			        "modified": schema.StringAttribute{
    			        	Optional:    true,
    			        	Description: "Modified timestamp",
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
    			        			//ElementType: types.ObjectType{
    			        			//	AttrTypes: map[string]attr.Type{
    			        			//		"deny": types.BoolType,
    			        			//		"end":  types.Int32Type,
    			        			//		"start": types.Int32Type,
    			        			//	},
    			        			//},
    			        			Description: "List of denied listening ports",
    			        		},
    			        		"denied_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of denied outbound IPs",
    			        		},
    			        		"denied_outbound_ports": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			//ElementType: types.ObjectType{
    			        			//	AttrTypes: map[string]attr.Type{
    			        			//		"deny": types.BoolType,
    			        			//		"end":  types.Int32Type,
    			        			//		"start": types.Int32Type,
    			        			//	},
    			        			//},
    			        			Description: "List of denied outbound ports",
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
    			        		//"intelligence_feed": schema.ListAttribute{
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
    			        //"dns": schema.SingleNestedAttribute{
    			        //	Optional:    true,
    			        //	Attributes: map[string]schema.Attribute{
    			        //		"allow": schema.ListAttribute{
    			        //			Optional:    true,
    			        //			ElementType: types.StringType,
    			        //			Description: "List of allowed DNS",
    			        //		},
    			        //		"deny": schema.ListAttribute{
    			        //			Optional:    true,
    			        //			ElementType: types.StringType,
    			        //			Description: "List of denied DNS",
    			        //		},
    			        //		"deny_list_effect": schema.ListAttribute{
    			        //			Optional:    true,
    			        //			ElementType: types.StringType,
    			        //			Description: "Effect for the deny list",
    			        //		},
    			        //		"intelligence_feed": schema.ListAttribute{
    			        //			Optional:    true,
    			        //			ElementType: types.StringType,
    			        //			Description: "Effect for the intelligence feed",
    			        //		},
    			        //	},
    			        //},
    			        "notes": schema.StringAttribute{
    			        	Optional:    true,
    			        	Description: "Notes for the resource",
    			        },
    			        "owner": schema.StringAttribute{
    			        	Optional:    true,
    			        	Description: "Owner of the resource",
    			        },
    			        "previous_name": schema.StringAttribute{
    			        	Optional:    true,
    			        	Description: "Previous name of the resource",
    			        },
                    },
                },
    		},
    	},
    }
}
