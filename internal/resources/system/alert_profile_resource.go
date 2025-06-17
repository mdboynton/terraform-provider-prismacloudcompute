package system

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	alertProfileAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/alertprofile"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &AlertProfileResource{}
var _ resource.ResourceWithImportState = &AlertProfileResource{}

type AlertProfileResource struct {
	client *api.PrismaCloudComputeAPIClient
}

func NewAlertProfileResource() resource.Resource {
	return &AlertProfileResource{}
}

func (r *AlertProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_profile"
}

var (
	policyModuleAttributes map[string]schema.Attribute = map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Description: "TODO",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"all_rules": schema.BoolAttribute{
			Description: "TODO",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
		"rules": schema.ListAttribute{
			Description: "TODO",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			Default:     listdefault.StaticValue(types.ListNull(types.StringType)),
		},
	}

	policyModuleAttributeTypes map[string]attr.Type = map[string]attr.Type{
		"enabled":   types.BoolType,
		"all_rules": types.BoolType,
		"rules": types.ListType{
			ElemType: types.StringType,
		},
	}

	policyModuleAttributeDefaultValues map[string]attr.Value = map[string]attr.Value{
		"enabled":   types.BoolValue(false),
		"all_rules": types.BoolValue(true),
		"rules":     types.ListValueMust(types.StringType, []attr.Value{}),
		//"rules": types.ListNull(types.StringType),
	}

	policyModuleDefaultValue basetypes.ObjectValue = types.ObjectValueMust(policyModuleAttributeTypes, policyModuleAttributeDefaultValues)

	policyDefaultValue basetypes.ObjectValue = types.ObjectValueMust(
		map[string]attr.Type{
			"admission": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"agentless_app_firewall": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"app_embedded_app_firewall": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"app_embedded_runtime": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"cloud_discovery": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"container_app_firewall": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"container_compliance": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"container_compliance_scan": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"container_runtime": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"container_vulnerability": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"defender": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"host_app_firewall": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"host_compliance": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"host_compliance_scan": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"host_runtime": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"host_vulnerability": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"incident": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"kubernetes_audit": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"network_firewall": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"registry_vulnerability": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"serverless_app_firewall": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"serverless_runtime": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"vm_compliance": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"vm_vulnerability": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
			"waas_health": types.ObjectType{
				AttrTypes: policyModuleAttributeTypes,
			},
		},
		map[string]attr.Value{
			"admission":                 policyModuleDefaultValue,
			"agentless_app_firewall":    policyModuleDefaultValue,
			"app_embedded_app_firewall": policyModuleDefaultValue,
			"app_embedded_runtime":      policyModuleDefaultValue,
			"cloud_discovery":           policyModuleDefaultValue,
			"container_app_firewall":    policyModuleDefaultValue,
			"container_compliance":      policyModuleDefaultValue,
			"container_compliance_scan": policyModuleDefaultValue,
			"container_runtime":         policyModuleDefaultValue,
			"container_vulnerability":   policyModuleDefaultValue,
			"defender":                  policyModuleDefaultValue,
			"host_app_firewall":         policyModuleDefaultValue,
			"host_compliance":           policyModuleDefaultValue,
			"host_compliance_scan":      policyModuleDefaultValue,
			"host_runtime":              policyModuleDefaultValue,
			"host_vulnerability":        policyModuleDefaultValue,
			"incident":                  policyModuleDefaultValue,
			"kubernetes_audit":          policyModuleDefaultValue,
			"network_firewall":          policyModuleDefaultValue,
			"registry_vulnerability":    policyModuleDefaultValue,
			"serverless_app_firewall":   policyModuleDefaultValue,
			"serverless_runtime":        policyModuleDefaultValue,
			"vm_compliance":             policyModuleDefaultValue,
			"vm_vulnerability":          policyModuleDefaultValue,
			"waas_health":               policyModuleDefaultValue,
		},
	)

	webhookTypes map[string]attr.Type = map[string]attr.Type{
		"ca_cert": types.StringType,
		"enabled": types.BoolType,
		"json":    types.StringType,
		"url":     types.StringType,
	}

	webhookJsonDefault string = `{
	"type": "#type",
	"time": "#time",
	"container": "#container",
	"containerID": "#containerID",
	"image": "#image",
	"imageID": "#imageID",
	"tags": "#tags",
	"host": "#host",
	"fqdn": "#fqdn",
	"function": "#function",
	"region": "#region",
	"provider": "#provider",
	"osRelease": "#osRelease",
	"osDistro": "#osDistro",
	"runtime": "#runtime",
	"appID": "#appID",
	"rule": "#rule",
	"message": "#message",
	"aggregatedAlerts": #aggregatedAlerts,
	"dropped": #dropped,
	"forensics": "#forensics",
	"category": "#category",
	"command": "#command",
	"startupProcess": "#startupProcess",
	"labels": #labels,
	"collections": #collections,
	"complianceIssues": #complianceIssues,
	"vulnerabilities": #vulnerabilities,
	"clusters": #clusters,
	"namespaces": #namespaces,
	"accountIDs": #accountIDs,
	"user": "#user",
	"incidentTime": "#incidentTime"
}`
)

func (r *AlertProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TODO",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "TODO",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vuln_immediate_alerts_enabled": schema.BoolAttribute{
				Description: "TODO",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"policy": schema.SingleNestedAttribute{
				Description: "TODO",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"admission": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"agentless_app_firewall": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"app_embedded_app_firewall": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"app_embedded_runtime": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"cloud_discovery": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"container_app_firewall": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"container_compliance": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"container_compliance_scan": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"container_runtime": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"container_vulnerability": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"defender": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"host_app_firewall": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"host_compliance": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"host_compliance_scan": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"host_runtime": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"host_vulnerability": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"incident": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"kubernetes_audit": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"network_firewall": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"registry_vulnerability": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"serverless_app_firewall": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"serverless_runtime": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"vm_compliance": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"vm_vulnerability": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
					"waas_health": schema.SingleNestedAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Attributes:  policyModuleAttributes,
						Default:     objectdefault.StaticValue(policyModuleDefaultValue),
						Validators: []validator.Object{
							validators.RulesNotConfiguredWithAllRulesEnabled(),
						},
					},
				},
				Default: objectdefault.StaticValue(policyDefaultValue),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"webhook": schema.SingleNestedAttribute{
				Description: "TODO",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"ca_cert": schema.StringAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"enabled": schema.BoolAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"json": schema.StringAttribute{
						Description: "TODO",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(webhookJsonDefault),
					},
					"url": schema.StringAttribute{
						// TODO: validation
						Description: "TODO",
						Required:    true,
					},
				},
				Default: objectdefault.StaticValue(types.ObjectValueMust(webhookTypes, map[string]attr.Value{
					"ca_cert": types.StringValue(""),
					"enabled": types.BoolValue(false),
					"json":    types.StringValue(""),
					"url":     types.StringValue(""),
				})),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AlertProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *AlertProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.AlertProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	request := plan.ToCreateOrUpdateRequest(ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create resource
	err := alertProfileAPI.CreateOrUpdateAlertprofile(*r.client, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Alert Profile resource",
			fmt.Sprintf("Failed to create alert profile: %s", err.Error()),
		)
		return
	}

	// Retrieve newly created resource
	createdResource, err := alertProfileAPI.GetAlertprofile(*r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Alert Profile resource",
			fmt.Sprintf("Failed to retrieve created alert profile: %s", err.Error()),
		)
		return
	}

	plan.RefreshPropertyValues(ctx, &resp.Diagnostics, *createdResource)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *AlertProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.AlertProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve data from Prisma Cloud
	data, err := alertProfileAPI.GetAlertprofile(*r.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Alert Profile resource",
			fmt.Sprintf("Failed to retrieve alert profile: %s", err.Error()),
		)
		return
	}

	state.RefreshPropertyValues(ctx, &resp.Diagnostics, *data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *AlertProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state models.AlertProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.AlertProfileResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	request := plan.ToCreateOrUpdateRequest(ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update resource
	err := alertProfileAPI.CreateOrUpdateAlertprofile(*r.client, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Alert Profile resource",
			fmt.Sprintf("Failed to update alert profile: %s", err.Error()),
		)
		return
	}

	// Retrieve newly created resource
	updatedResource, err := alertProfileAPI.GetAlertprofile(*r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Alert Profile resource",
			fmt.Sprintf("Failed to retrieve updated alert profile: %s", err.Error()),
		)
		return
	}

	plan.RefreshPropertyValues(ctx, &resp.Diagnostics, *updatedResource)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *AlertProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.AlertProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete alert profile
	err := alertProfileAPI.DeleteAlertprofile(*r.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Alert Profile resource",
			"Failed to delete alert profile: "+err.Error(),
		)
	}
}

// TODO: Ensure this works properly
func (r *AlertProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
