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
            // TODO: probably dont need this
            "automatic_runtime_learning": schema.BoolAttribute{
                MarkdownDescription: "TODO",
                Optional:            true,
                Computed:            true,
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
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
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
                                    }, 
                                    map[string]attr.Value{
                                        "enabled": types.BoolValue(true),
                                        "allowed_paths": types.ListValueMust(types.StringType, []attr.Value{}),
                                        "denied_paths": types.ListValueMust(types.StringType, []attr.Value{}),
                                        "denied_paths_effect": types.StringValue("alert"),
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
                                // TODO: if block_all_processes_except_main is set to true, this cannot be set
    			        		"allowed_processes": schema.ListAttribute{
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
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
                                    },
    			        			Description: "",
    			        		},
    			        		"crypto_miners": schema.BoolAttribute{
    			        			Optional:    true,
    			        			Computed:    true,
                                    Default: booldefault.StaticBool(true),
                                    Description: "TODO",
    			        		},
    			        		"block_all_processes_except_main": schema.BoolAttribute{
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
    			        		"allowed_listening_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
    			        		"allowed_outbound_internet_ports": schema.ListAttribute{
    			        			Optional:    true,
                                    Computed: true,
    			        			ElementType: types.StringType,
                                    Default: listdefault.StaticValue(
                                        types.ListValueMust(types.StringType, []attr.Value{}),
                                    ),
    			        			Description: "TODO",
    			        		},
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
                                    Default: stringdefault.StaticString("disable"),
                                    Validators: []validator.String{
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
                                    },
    			        			Description: "",
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
                                        validators.PolicyEffectIsValid([]string{"alert", "prevent"}),
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
