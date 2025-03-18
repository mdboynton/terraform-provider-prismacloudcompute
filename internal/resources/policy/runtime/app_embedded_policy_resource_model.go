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
                            MarkdownDescription: "TODO",
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
    			        	Description: "List of collections.",
    			        },
    			        "custom_rules": schema.ListNestedAttribute{
                            Optional: true,
                            Description: "List of custom rules.",
                            Validators: []validator.List{
                                validators.CustomRulesAreValid(),
                            },
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "id": schema.Int64Attribute{
                                        Computed:   true,
                                        Description: "TODO",
                                    },
                                    "name": schema.StringAttribute{
                                        Optional:   true,
                                        Description: "TODO",
                                    },
                                    "effect": schema.StringAttribute{
                                        Required:   true,
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid("custom_rules.effect", []string{"allow", "alert", "prevent", "block"}),
                                        },
                                        Description: "TODO",
                                    },
                                    "log_as": schema.StringAttribute{
                                        Optional:   true,
    			        			    Computed:   true,
                                        Validators: []validator.String{
                                            validators.PolicyEffectIsValid("custom_rules.log_as", []string{"audit", "incident"}),
                                        },
                                        Default:    stringdefault.StaticString(""),
                                        Description: "TODO",
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
    			        "file_system": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
    			        	Attributes: map[string]schema.Attribute{
                                "enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO.",
                                },
                                // TODO: allowed_paths cannot be set if denied_paths is set
    			        		"allowed_paths": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
                                // TODO: denied_paths cannot be set if allowed_paths is set
    			        		"denied_paths": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
    			        		"denied_paths_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("file_system.denied_paths_effect", []string{"alert", "prevent"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "TODO",
    			        		},
                                "changes_to_binaries_and_certs": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO.",
                                },
                                "detection_of_encrypted_binaries": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO.",
                                },
                                "changes_to_ssh_admin_account_config_files": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO.",
                                },
                                "suspicious_elf_headers": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO.",
                                },
                                "malware_from_custom_feed": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO.",
                                },
    			        		"wild_fire_analysis": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("file_system.wild_fire_analysis", []string{"disable", "alert"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
                                    Description: "TODO",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
                                ),
                            ),
    			        },
    			        "processes": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
    			        	Attributes: map[string]schema.Attribute{
                                "enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO.",
                                },
                                // TODO: allowed_processes cannot be set if denied_processes is set and vis versa
    			        		"allowed_processes": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
    			        		"denied_processes": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
    			        		"denied_processes_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("processes.denied_processes_effect", []string{"alert", "prevent"}),
                                    },
    			        			Description: "",
    			        		},
    			        		"crypto_miners": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO",
    			        		},
    			        		"processes_from_modified_binaries": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(
                                types.ObjectValueMust(
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
    			        	Description: "Name of the resource",
    			        },
    			        "networking": schema.SingleNestedAttribute{
    			        	Optional:    true,
    			        	Computed:    true,
    			        	Attributes: map[string]schema.Attribute{
                                "ip_connectivity_enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO",
                                },
                                // TODO: allowed_listening_ports cannot be set if denied_listening_ports is set and vis versa
    			        		"allowed_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
                                // TODO: allowed_outbound_internet_ports cannot be set if denied_outbound_internet_ports is set and vis versa
    			        		"allowed_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
                                // TODO: allowed_outbound_ips cannot be set if denied_outbound_internet_ports is set and vis versa
    			        		"allowed_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
    			        		"denied_ips_ports_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("networking.denied_ips_ports_effect", []string{"alert", "prevent"}),
                                    },
    			        			Description: "",
    			        		},
    			        		"denied_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
    			        		"denied_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
    			        		"denied_outbound_ips": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
                                "dns_enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO",
                                },
    			        		"allowed_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "",
    			        		},
    			        		"denied_dns_domains_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid("networking.denied_dns_domains_effect", []string{"alert", "prevent"}),
                                    },
    			        			Description: "",
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
