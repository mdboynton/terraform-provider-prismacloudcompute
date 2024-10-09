package system

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
)


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
