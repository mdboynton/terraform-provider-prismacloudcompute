package policy

import (
	"context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"
	//policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
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
                            Computed: true,
    			        	Attributes: map[string]schema.Attribute{
    			        		"allowed_processes": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of allowed processes",
    			        		},
    			        		"crypto_miners": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.crypto_miners", []string{"disable", "alert", "prevent"}),
                                    },
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
                                            Validators: []validator.String{
                                                validators.PolicyEffectIsValid("anti_malware.denied_processes.effect", []string{"alert", "prevent"}),
                                            },
    			        				},
    			        				"paths": schema.ListAttribute{
    			        					Optional:    true,
    			        					ElementType: types.StringType,
    			        					Description: "List of paths for denied processes",
                                            // TODO: validation
    			        				},
    			        			},
    			        		},
    			        		"encrypted_binaries": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.encrypted_binaries", []string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for encrypted/packed binaries.",
    			        		},
    			        		"execution_flow_hijacking": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.execution_flow_hijacking", []string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for execution flow hijacking.",
    			        		},
    			        		"malware_from_advanced_threat_protection": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.malware_from_advanced_threat_protection", []string{"disable", "alert"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for intelligence feed.",
    			        		},
    			        		"malware_from_custom_feed": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.malware_from_custom_feed", []string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for malware based on custom feed",
    			        		},
    			        		"non_packaged_binaries_service": schema.StringAttribute{
    			        			Optional:   true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.non_packaged_binaries_service", []string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for service unknown origin binary.",
                                    // TODO: description
    			        		},
    			        		"non_packaged_binaries_user": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.non_packaged_binaries_user", []string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for user unknown origin binary.",
    			        		},
    			        		"processes_temporary_storage": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.processes_temporary_storage", []string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for processes running from temporary storage.",
    			        		},
    			        		"reverse_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.reverse_shell", []string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for reverse shell attacks.",
    			        		},
    			        		"suppress_compiler_generated_binaries": schema.BoolAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Default:    booldefault.StaticBool(false),
    			        			Description: "Detect compiler generated binary.",
                                    // TODO: update Description
                                    // TODO: do not allow if non_packaged_binaries_run_by_service is set to "disable"
    			        		},
    			        		"suspicious_elf_headers": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.suspicious_elf_headers", []string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for binaries with suspicious ELF headers.",
    			        		},
    			        		"web_shell": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.web_shell", []string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for web shell attacks. Note that when setting this effect to \"Prevent\", the \"Prevent\" effect will only apply to file execution. Alerts will be generated on file creation.",
    			        		},
    			        		"wild_fire_analysis": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("anti_malware.wild_fire_analysis", []string{"disable", "alert"}),
                                    },
                                    Default:    stringdefault.StaticString("alert"),
    			        			Description: "Effect for WildFire analysis. WildFire must be enabled for runtime protection under Manage > System > WildFire.",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
                                ),
                            ),
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
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid("custom_rules.effect", []string{"allow", "alert", "prevent"}),
                                        },
                                        Description: "",
                                    },
                                    "log_as": schema.StringAttribute{
                                        Optional:   true,
    			        			    Computed:   true,
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid("custom_rules.log_as", []string{"audit", "incident"}),
                                        },
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
                            Validators: []validator.List{
                                validators.FileIntegrityRulesAreValid(),
                            },
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
    			        	Description: "Name of the resource",
    			        },
    			        "networking": schema.SingleNestedAttribute{
    			        	Optional:    true,
    			        	Computed:    true,
    			        	Attributes: map[string]schema.Attribute{
    			        		"allowed_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
    			        			ElementType: types.StringType,
    			        			Description: "List of allowed outbound IPs.",
    			        		},
    			        		"suspicious_ips_custom_feed": schema.StringAttribute{
    			        			Optional:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("networking.suspicious_ips_custom_feed", []string{"disable", "alert"}),
                                    },
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
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("networking.denied_ips_ports_effect", []string{"disable", "alert"}),
                                    },
    			        			Description: "Effect for the IP/ports deny list",
                                    // TODO: description
    			        		},
    			        		"suspicious_ips_advanced_threat_protection_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("networking.suspicious_ips_advanced_threat_protection_effect", []string{"disable", "alert"}),
                                    },
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
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("networking.denied_dns_domains_effect", []string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("disable"),
    			        			Description: "Effect for the DNS deny list.",
                                    // TODO: description
    			        		},
    			        		"suspicious_domains_advanced_threat_protection_effect": schema.StringAttribute{
    			        			Optional:   true,
                                    Computed:   true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("networking.suspicious_domains_advanced_threat_protection_effect", []string{"disable", "alert", "prevent"}),
                                    },
                                    Default:    stringdefault.StaticString("disable"),
    			        			Description: "Effect for the intelligence feed",
                                    // TODO: description
                                    // TODO: validation
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
