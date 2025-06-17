package policy

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ApplicationControlPolicyResource{}
var _ resource.ResourceWithImportState = &ApplicationControlPolicyResource{}

func NewApplicationControlPolicyResource() resource.Resource {
    return &ApplicationControlPolicyResource{}
}

type ApplicationControlPolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *ApplicationControlPolicyResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "rules": schema.ListNestedAttribute{
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.Int32Attribute{
                            MarkdownDescription: "TODO",
                            Computed: true,
                            PlanModifiers: []planmodifier.Int32{
                                int32planmodifier.UseStateForUnknown(),
                            },
                        },
                        "applications": schema.ListNestedAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "name": schema.StringAttribute{
                                        MarkdownDescription: "TODO",
                                        Optional: true,
                                        Computed: true,
                                        Validators: []validator.String{
                                            validators.IsNonEmpty("Application name"),
                                        },
                                    },
                                    "allowed_versions": schema.ListAttribute{
                                        MarkdownDescription: "TODO",
                                        Required: true,
                                        ElementType: types.ListType{
                                            ElemType: types.StringType,
                                        },
                                    },
                                },
                            },
                        },
                        "description": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Optional: true,
                            Computed: true,
                            Default: stringdefault.StaticString(""),
                        },
                        "modified": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Computed: true,
                        },
                        "name": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Required: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.RequiresReplace(),
                            },
                            Validators: []validator.String{
                                validators.IsNonEmpty("Rule name"),
                            },
                        },
                        "owner": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.UseStateForUnknown(),
                            },
                        },
                        "previous_name": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Computed: true,
                            PlanModifiers: []planmodifier.String{
                                stringplanmodifier.UseStateForUnknown(),
                            },
                        },
                        "severity": schema.StringAttribute{
                            MarkdownDescription: "TODO",
                            Required: true,
                            Validators: []validator.String{
                                stringvalidator.OneOf("low", "medium", "high", "critical"),
                            },
                        },
                    },
                },
            },
        },
    }
}
