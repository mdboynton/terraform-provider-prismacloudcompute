package data_sources

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/system"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &CollectionDataSource{}

func NewCollectionDataSource() datasource.DataSource {
	return &CollectionDataSource{}
}

type CollectionDataSource struct {
	client *api.PrismaCloudComputeAPIClient
}

func (d *CollectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (d *CollectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		//MarkdownDescription: "TODO",
		Description: "Collections are predefined filters that let you group related resources together. They can be used to scope policy rules and segment data/views in the Console UI and the Prisma Cloud API.",
		Attributes: map[string]schema.Attribute{
			"account_ids": schema.SetAttribute{
				Description: "List of account IDs.",
				//MarkdownDescription: "TODO",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"application_ids": schema.SetAttribute{
				Description: "List of application IDs.",
				//MarkdownDescription: "TODO",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"clusters": schema.SetAttribute{
				Description: "List of Kubernetes cluster names.",
				//MarkdownDescription: "TODO",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"color": schema.StringAttribute{
				Description: "Hexadecimal representation of the collection's color value.",
				//MarkdownDescription: "TODO",
				Optional: true,
				Computed: true,
			},
			"containers": schema.SetAttribute{
				Description: "List of containers.",
				//MarkdownDescription: "TODO",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				//MarkdownDescription: "TODO",
				Description: "Description of the collection.",
				Optional:    true,
				Computed:    true,
			},
			"functions": schema.SetAttribute{
				//MarkdownDescription: "TODO",
				Description: "List of serverless functions.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"hosts": schema.SetAttribute{
				//MarkdownDescription: "TODO",
				Description: "List of hosts.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"images": schema.SetAttribute{
				//MarkdownDescription: "TODO",
				Description: "List of images.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"labels": schema.SetAttribute{
				//MarkdownDescription: "TODO",
				Description: "List of labels.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"modified": schema.StringAttribute{
				//MarkdownDescription: "TODO",
				Description: "Date and time that the collection was last modified.",
				//Optional: true, // TODO: get rid of thi)s and make it just Computed
				Computed: true,
			},
			"name": schema.StringAttribute{
				//MarkdownDescription: "TODO",
				Description: "Collection name. Must be unique.",
				Required:    true,
			},
			"namespaces": schema.SetAttribute{
				//MarkdownDescription: "TODO",
				Description: "List of Kubernetes namespaces.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"owner": schema.StringAttribute{
				//MarkdownDescription: "TODO",
				Description: "User who created or last modified the collection.",
				Computed:    true,
			},
			"prisma": schema.BoolAttribute{
				//MarkdownDescription: "TODO",
				Description: "Indicates whether this collection originated from Prisma Cloud.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				//MarkdownDescription: "TODO",
				Description: "Indicates whether this collection was created by a user (true) or the system (false).",
				Computed:    true,
			},
		},
	}
}

func (d *CollectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PrismaCloudComputeAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *CollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data  models.CollectionDataSourceModel
		diags diag.Diagnostics
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the data from the API
	collection, err := collectionAPI.GetCollection(*d.client, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error fetching Collection data source with name: %s", data.Name.ValueString()),
			fmt.Sprintf("Error message: %s", err.Error()),
		)
		return
	}

	data, diags = data.RefreshPropertyValues(ctx, collection)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
