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

var _ resource.Resource = &ServerlessRuntimePolicyResource{}
var _ resource.ResourceWithImportState = &ServerlessRuntimePolicyResource{}

func NewServerlessRuntimePolicyResource() resource.Resource {
    return &ServerlessRuntimePolicyResource{}
}

type ServerlessRuntimePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *ServerlessRuntimePolicyResource) GetSchema(ctx context.Context) schema.Schema {
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
                            Description: "List of collection names. Used to scope the rule. Note that in order for a collection to be attached to this type of policy rule, it must contain only the wildcard value (\"*\") for all of the following resource types: Containers, Hosts, Images, Labels, App IDs, Namespaces and Clusters.",
                            Default: setdefault.StaticValue(
                                types.SetValueMust(types.StringType, []attr.Value{types.StringValue("All")}),
                            ),
    			        	ElementType: types.StringType,
    			        },
    			        "disabled": schema.BoolAttribute{
    			        	Optional:   true,
                            Computed:   true,
                            Default: booldefault.StaticBool(false),
    			        	Description: "Indicates whether to disable the rule.",
    			        },
    			        "advanced_threat_protection": schema.BoolAttribute{
    			        	Optional:   true,
                            Computed:   true,
                            Default: booldefault.StaticBool(true),
    			        	Description: "Toggles Serverless Advanced Threat Protection. Serverless Advanced Threat Protection (ATP) is a collection of paths (researched by Prisma Cloud Labs) that define which file system or process activity is allowed within the function. Activities that do not match these paths will raise a security audit (Note: filesystem monitoring must be enabled for this to work). When enabled, it creates an automatic hardening for the function in runtime, without the need to manually configure the runtime policy.",
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
                                // TODO: allowed_paths cannot be set if denied_paths is set
    			        		"allowed_paths": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "List of file system paths which will not be monitored. Cannot be configured if denied_paths is configured.",
    			        		},
                                // TODO: denied_paths cannot be set if allowed_paths is set
    			        		"denied_paths": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "File system paths to be alerted on or suppressed. Cannot be configured if allowed_paths is configured.",
    			        		},
    			        		"denied_paths_effect": schema.StringAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Description: "Effect for detected file system paths from the deny list. Must be either \"alert\" or \"prevent\"",
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
                                    },
                                    Default: stringdefault.StaticString("alert"),
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(serverlessFileSystemDefault),
    			        },
    			        "processes": schema.SingleNestedAttribute{
    			        	Optional: true,
                            Computed: true,
                            Description: "Configuration for process monitoring.",
    			        	Attributes: map[string]schema.Attribute{
                                "enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables process monitoring.",
                                },
                                // TODO: if block_all_processes_except_main is set to true, this cannot be set
    			        		"allowed_processes": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
                                    Description: "List of process names to be whitelisted.",
    			        		},
    			        		"denied_processes_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
                                    Description: "Effect for detected denied processes. Must be either \"alert\" or \"prevent\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
                                    },
    			        		},
    			        		"crypto_miners": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enable crypto miner detection.",
    			        		},
    			        		"block_all_processes_except_main": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "Enables blocking of all processes except the main process.",
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(serverlessProcessesDefault),
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
    			        			Description: "Effect for denied IPs and ports. Must be either \"alert\" or \"prevent\".",
                                    Default: stringdefault.StaticString("disable"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
                                    },
    			        		},
                                "dns_enabled": schema.BoolAttribute{
                                    Optional: true,
                                    Computed:   true,
                                    Description: "Enables DNS monitoring.",
                                    Default: booldefault.StaticBool(true),
                                },
    			        		"allowed_dns_domains": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			Description: "List of DNS domains which will not generate alerts or be prevented.",
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        		},
    			        		"denied_dns_domains_effect": schema.StringAttribute{
    			        			Optional:    true,
                                    Computed:    true,
    			        			Description: "Effect for DNS domains not specified in the allow list. Must be either \"alert\" or \"prevent\".",
                                    Default: stringdefault.StaticString("alert"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
                                    },
    			        		},
    			        	},
                            Default: objectdefault.StaticValue(serverlessNetworkingDefault),
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
