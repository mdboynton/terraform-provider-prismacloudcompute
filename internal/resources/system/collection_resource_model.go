package system

import (
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
)

var _ resource.Resource = &CollectionResource{}
var _ resource.ResourceWithImportState = &CollectionResource{}

func NewCollectionResource() resource.Resource {
    return &CollectionResource{}
}

type CollectionResource struct {
    client *api.PrismaCloudComputeAPIClient
}

type CollectionResourceModel struct {
    AccountIDs types.Set `tfsdk:"account_ids"`
    AppIDs types.Set `tfsdk:"app_ids"`
    Clusters types.Set `tfsdk:"clusters"`
    Color types.String `tfsdk:"color"`
    Containers types.Set `tfsdk:"containers"`
    Description types.String `tfsdk:"description"`
    Functions types.Set `tfsdk:"functions"`
    Hosts types.Set `tfsdk:"hosts"`
    Images types.Set `tfsdk:"images"`
    Labels types.Set `tfsdk:"labels"`
    Modified types.String `tfsdk:"modified"`
    Name types.String `tfsdk:"name"`
    Namespaces types.Set `tfsdk:"namespaces"`
    Owner types.String `tfsdk:"owner"`
    Prisma types.Bool `tfsdk:"prisma"`
    System types.Bool `tfsdk:"system"`
}

func (r *CollectionResource) GetSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "TODO",
        Attributes: map[string]schema.Attribute{
            "account_ids": schema.SetAttribute{
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
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
            "clusters": schema.SetAttribute{
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString("#3FA2F7"),
            },
            "containers": schema.SetAttribute{
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
                Optional: true,
                Computed: true,
                Default: stringdefault.StaticString(""),
            },
            "functions": schema.SetAttribute{
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "TODO",
                Required: true,
            },
            "namespaces": schema.SetAttribute{
                MarkdownDescription: "TODO",
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
                MarkdownDescription: "TODO",
                Computed: true,
            },
            "prisma": schema.BoolAttribute{
                MarkdownDescription: "TODO",
                Computed: true,
                Default: booldefault.StaticBool(false),
            },
            "system": schema.BoolAttribute{
                MarkdownDescription: "TODO",
                Computed: true,
                Default: booldefault.StaticBool(false),
            },
        },
    }
}
