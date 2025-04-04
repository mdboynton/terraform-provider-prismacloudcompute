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

var _ resource.Resource = &AppEmbeddedRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &AppEmbeddedRuntimePolicyResource{}

func NewAppEmbeddedRuntimePolicyResource() resource.Resource {
    return &AppEmbeddedRuntimePolicyResource{}
}

type AppEmbeddedRuntimePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *AppEmbeddedRuntimePolicyResource) GetSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
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
    			        "collections": schema.SetAttribute{
    			        	Optional:    true,
    			        	Computed:    true,
                            Default: setdefault.StaticValue(
                                types.SetValueMust(types.StringType, []attr.Value{types.StringValue("All")}),
                            ),
    			        	ElementType: types.StringType,
                            Description: "List of collection names. Used to scope the rule. Note that in order for a collection to be attached to this type of policy rule, it must contain only the wildcard value (\"*\") for all of the following resource types: Hosts, Labels, Functions and Namespaces.",
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
                                            validators.PolicyEffectIsValid([]string{"allow", "alert", "prevent", "block"}),
                                        },
                                        Description: "Effect for the custom rule. Must be one of \"allow\", \"alert\" or \"prevent\". Note that if set to \"allow\", the value of log_as will not have any effect.",
                                    },
                                    "log_as": schema.StringAttribute{
                                        Optional:   true,
    			        			    Computed:   true,
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid([]string{"audit", "incident"}),
                                        },
                                        Default:    stringdefault.StaticString(""),
                                        Description: "How violations of this custom rule will be logged. Must be \"audit\" or \"incident\".",
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
                            Validators: []validator.Object{
                                validators.ErrorIfBothListsConfigured("allowed_paths", "denied_paths"),
                            },
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
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "List of file system paths which will not be monitored.",
    			        		},
    			        		"denied_paths": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "File system paths to be alerted on or suppressed.",
    			        		},
    			        		"denied_paths_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["app_embedded"]["file_system.denied_paths_effect"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "Effect for detected file system paths from the deny list. Must be either \"alert\" or \"prevent\"",
    			        		},
                                "changes_to_binaries_and_certs": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables monitoring of changes to binary or certificate files.",
                                },
                                "detection_of_encrypted_binaries": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables detection of encrypted/packed binaries.",
                                },
                                "changes_to_ssh_admin_account_config_files": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables monitoring of changes to SSH and admin account configuration files.",
                                },
                                "suspicious_elf_headers": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables monitoring of binaries with suspicious ELF headers.",
                                },
                                "malware_from_custom_feed": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables monitoring of files classified as malware by the presence of their MD5 hash in the custom malware signature feed. Custom malware signatures can be configured in the console by navigating to Manage > System > Custom feeds > Malware signatures.",
                                },
    			        		"wild_fire_analysis": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["app_embedded"]["file_system.wild_fire_analysis"]),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        			Description: "Effect for detected files classified as malware by WildFire, Palo Alto Networks' malware analysis engine. Must be either \"disable\" or \"alert\". WildFire must be enabled for runtime protection under Manage > System > WildFire.",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(appEmbeddedFileSystemDefault),
    			        },
    			        "processes": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
                            Description: "Configuration for process monitoring.",
                            Validators: []validator.Object{
                                validators.ErrorIfBothListsConfigured("allowed_processes", "denied_processes"),
                            },
    			        	Attributes: map[string]schema.Attribute{
                                "enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables process monitoring.",
                                },
    			        		"allowed_processes": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
                                    Description: "List of process names to be whitelisted.",
    			        		},
    			        		"denied_processes": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
                                    Description: "List of processes to deny or alert on.",
    			        		},
    			        		"denied_processes_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["app_embedded"]["processes.denied_processes_effect"]),
                                    },
                                    Description: "Effect for detected denied processes. Must be either \"alert\" or \"prevent\".",
    			        		},
    			        		"crypto_miners": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enable crypto miner detection.",
    			        		},
    			        		"processes_from_modified_binaries": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enable execution of binaries that do not belong to the original image.",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(appEmbeddedProcessesDefault),
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
                            Validators: []validator.Object{
                                validators.ErrorIfBothListsConfigured("allowed_listening_ports", "denied_listening_ports"),
                                validators.ErrorIfBothListsConfigured("allowed_outbound_internet_ports", "denied_outbound_internet_ports"),
                                validators.ErrorIfBothListsConfigured("allowed_outbound_ips", "denied_outbound_ips"),
                            },
    			        	Attributes: map[string]schema.Attribute{
                                "ip_connectivity_enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables IP connectivity monitoring.",
                                },
    			        		"allowed_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "List of listening ports which will not generate alerts or be prevented.",
    			        		},
    			        		"allowed_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "List of outbound internet ports which will not generate alerts or be prevented.",
    			        		},
    			        		"allowed_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "List of outbound IPs which will not generate alerts or be prevented.",
    			        		},
    			        		"denied_ips_ports_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["app_embedded"]["networking.denied_ips_ports_effect"]),
                                    },
    			        			Description: "Effect for denied IPs and ports. Must be either \"alert\" or \"prevent\".",
    			        		},
    			        		"denied_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
                                    Description: "List of listening ports for which access will be alerted on or prevented.",
    			        		},
    			        		"denied_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
                                    Description: "List of outbound internet ports for which access will be alerted on or prevented.",
    			        		},
    			        		"denied_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
                                    Description: "List of outbound IPs for which access will be alerted on or prevented.",
    			        		},
                                "dns_enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables DNS monitoring.",
                                },
    			        		"allowed_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "List of DNS domains which will not generate alerts or be prevented.",
    			        		},
    			        		"denied_dns_domains_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid(ValidEffects["app_embedded"]["networking.denied_dns_domains_effect"]),
                                    },
    			        			Description: "Effect for DNS domains not specified in the allow list. Must be either \"alert\" or \"prevent\".",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(appEmbeddedNetworkingDefault),
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
