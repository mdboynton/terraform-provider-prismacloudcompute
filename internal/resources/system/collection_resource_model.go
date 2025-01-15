package system

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

var _ resource.Resource = &CollectionResource{}
var _ resource.ResourceWithImportState = &CollectionResource{}

func NewCollectionResource() resource.Resource {
    return &CollectionResource{}
}

type CollectionResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *CollectionResource) GetSchema() schema.Schema {
    return schema.Schema{
        //MarkdownDescription: "TODO",
        Description: "Collections are predefined filters that let you group related resources together. They can be used to scope policy rules and segment data/views in the Console UI and the Prisma Cloud API.",
        Attributes: map[string]schema.Attribute{
            "account_ids": schema.SetAttribute{
                Description: "List of account IDs.",
                //MarkdownDescription: "TODO",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "app_ids": schema.SetAttribute{
                Description: "List of application IDs.",
                //MarkdownDescription: "TODO",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Validators: []validator.Set{
                    validators.AppIDsEndWithWildcard(),
                },
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "clusters": schema.SetAttribute{
                Description: "List of Kubernetes cluster names.",
                //MarkdownDescription: "TODO",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "color": schema.StringAttribute{
                Description: "Hexadecimal representation of the collection's color value.",
                //MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("#3FA2F7"),
            },
            "containers": schema.SetAttribute{
                Description: "List of containers.",
                //MarkdownDescription: "TODO",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "description": schema.StringAttribute{
                //MarkdownDescription: "TODO",
                Description: "Description of the collection.",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString(""),
            },
            "functions": schema.SetAttribute{
                //MarkdownDescription: "TODO",
                Description: "List of serverless functions.",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "hosts": schema.SetAttribute{
                //MarkdownDescription: "TODO",
                Description: "List of hosts.",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "images": schema.SetAttribute{
                //MarkdownDescription: "TODO",
                Description: "List of images.",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "labels": schema.SetAttribute{
                //MarkdownDescription: "TODO",
                Description: "List of labels.",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "modified": schema.StringAttribute{
                //MarkdownDescription: "TODO",
                Description: "Date and time that the collection was last modified.",
                //Optional: true, // TODO: get rid of this and make it just Computed
                Computed: true,
            },
            "name": schema.StringAttribute{
                //MarkdownDescription: "TODO",
                Description: "Collection name. Must be unique.",
                Required: true,
            },
            "namespaces": schema.SetAttribute{
                //MarkdownDescription: "TODO",
                Description: "List of Kubernetes namespaces.",
                ElementType: types.StringType,
                Optional: true,
                Computed: true,
                Default: setdefault.StaticValue(
                    types.SetValueMust(
                        types.StringType,
                        []attr.Value{
                            types.StringValue("*"),
                        },
                    ),
                ),
            },
            "owner": schema.StringAttribute{
                //MarkdownDescription: "TODO",
                Description: "User who created or last modified the collection.",
                Computed: true,
            },
            "prisma": schema.BoolAttribute{
                //MarkdownDescription: "TODO",
                Description: "Indicates whether this collection originated from Prisma Cloud.",
                Computed: true,
                Default: booldefault.StaticBool(false),
            },
            "system": schema.BoolAttribute{
                //MarkdownDescription: "TODO",
                Description: "Indicates whether this collection was created by a user (true) or the system (false).",
                Computed: true,
                Default: booldefault.StaticBool(false),
            },
        },
    }
}
