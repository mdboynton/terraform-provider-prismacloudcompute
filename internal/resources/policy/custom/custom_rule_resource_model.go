package custom 

import (
    "context"

    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

var _ resource.Resource = &CustomRuntimeRuleResource{}
var _ resource.ResourceWithImportState = &CustomRuntimeRuleResource{}

func NewCustomRuntimeRuleResource() resource.Resource {
    return &CustomRuntimeRuleResource{}
}

type CustomRuntimeRuleResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *CustomRuntimeRuleResource) GetSchema(ctx context.Context) schema.Schema {
    return schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.Int64Attribute{
                Computed:    true,
                Description: "ID of the resource",
            },
            "attack_techniques": schema.ListAttribute{
                Optional:    true,
                ElementType: types.StringType,
                Description: "List of attack techniques",
            },
            "description": schema.StringAttribute{
                Optional:    true,
                Computed:   true,
                Default:    stringdefault.StaticString(""),
                Description: "Description of the resource",
            },
            "message": schema.StringAttribute{
                Required:    true,
                Description: "Message associated with the resource",
            },
            "min_version": schema.StringAttribute{
                Optional:    true,
                Computed:   true,
                Default:    stringdefault.StaticString(""),
                // TODO: description (you never really set this, seemingly)
                Description: "Minimum version required for the resource",
            },
            "modified": schema.Int64Attribute{
                Computed:    true,
                Description: "Timestamp of last modification",
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "Name of the resource",
            },
            "owner": schema.StringAttribute{
                Computed:    true,
                Description: "Owner of the resource",
            },
            "script": schema.StringAttribute{
                Required:    true,
                Description: "Script associated with the resource",
            },
            "type": schema.StringAttribute{
                Required:    true,
                Description: "Type of the resource",
            },
            "vuln_ids": schema.ListAttribute{
                Computed:    true,
                ElementType: types.StringType,
                Description: "List of vulnerability IDs",
                // TODO: description (you never really set this, seemingly)
            },
        },
    }
}
