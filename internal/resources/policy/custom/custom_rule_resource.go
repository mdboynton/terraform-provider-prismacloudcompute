package custom

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	ruleAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/rule"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"

	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	//"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"

	//"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &CustomRuleResource{}
var _ resource.ResourceWithImportState = &CustomRuleResource{}

func NewCustomRuleResource() resource.Resource {
	return &CustomRuleResource{}
}

type CustomRuleResource struct {
	client *api.PrismaCloudComputeAPIClient
}

func (r *CustomRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_rule"
}

func (r *CustomRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the resource",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"attack_techniques": schema.ListAttribute{
				// TODO: should probably be computed?
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of attack techniques",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Description of the resource",
			},
			"message": schema.StringAttribute{
				Required:    true,
				Description: "Message associated with the resource",
			},
			"min_version": schema.StringAttribute{
				Computed: true,
				// TODO: description (this is set through adding a condition for the version in the rule script, at least for WAAS)
				Description: "Minimum version required for the resource",
			},
			"modified": schema.Int64Attribute{
				Computed:    true,
				Description: "Timestamp of last modification",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owner": schema.StringAttribute{
				Computed:    true,
				Description: "Owner of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				//Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
		},
	}
}

func (r *CustomRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	//resource.ImportStatePassthroughWithIdentity(ctx, path.Root("name"), req, resp)
}

func (r *CustomRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.client.MutexMap[api.MutexMapKeyCustomRules].Lock()
	defer r.client.MutexMap[api.MutexMapKeyCustomRules].Unlock()

	var plan models.CustomRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	data := plan.ToTerraform(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// TESTING:
	rules, _ := ruleAPI.ListCustomRules(*r.client)
	for _, rule := range rules {
		if rule.Owner == "admin" && rule.Name == data.Name {
			tflog.Debug(ctx, fmt.Sprintf("\n\n\ndeleting rule %s\n\n\n", rule.Name))
			ruleAPI.DeleteCustomRule(*r.client, rule.Id)
		}
	}
	// TESTING:

	// Create new custom runtime rule
	createdRuleId, err := ruleAPI.CreateCustomRule(*r.client, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Custom Rule resource",
			"Failed to create custom rule: "+err.Error(),
		)
		return
	}

	// Retrieve newly created custom runtime rule
	response, err := ruleAPI.GetCustomRuleById(*r.client, createdRuleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Custom Rule resource",
			fmt.Sprintf("Failed to retrieve created custom rule: %s", err.Error()),
		)
		return
	}

	plan.UpdateValues(ctx, &resp.Diagnostics, *response)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.client.MutexMap[api.MutexMapKeyCustomRules].Lock()
	defer r.client.MutexMap[api.MutexMapKeyCustomRules].Unlock()

	// Get current state
	var state models.CustomRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get policy value from Prisma Cloud
	//data, err := ruleAPI.GetCustomRuleById(*r.client, int(state.Id.ValueInt64()))
	data, err := ruleAPI.GetCustomRuleByName(*r.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Custom Rule resource",
			fmt.Sprintf("Failed to read custom rule: %s", err.Error()),
		)
		return
	}

	// Overwrite state values with Prisma Cloud data
	state.UpdateValues(ctx, &resp.Diagnostics, *data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.client.MutexMap[api.MutexMapKeyCustomRules].Lock()
	defer r.client.MutexMap[api.MutexMapKeyCustomRules].Unlock()

	// Get current state
	var state models.CustomRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.CustomRuleResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	data := plan.ToTerraform(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing rule
	err := ruleAPI.UpdateCustomRule(*r.client, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Custom Rule resource",
			fmt.Sprintf("Failed to update custom rule: %s", err.Error()),
		)
		return
	}

	// Get updated rule value from Prisma Cloud
	updatedRule, err := ruleAPI.GetCustomRuleById(*r.client, int(plan.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Custom Rule resource",
			fmt.Sprintf("Failed to read updated custom rule: %s", err.Error()),
		)
		return
	}

	//// Convert updated rule into schema
	//policySchema, diags := terraformToSchema(ctx, *updatedRule)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//	// Set updated state
	//	diags = resp.State.Set(ctx, policySchema)
	//	resp.Diagnostics.Append(diags...)
	//	if resp.Diagnostics.HasError() {
	//		return
	//	}
	plan.UpdateValues(ctx, &resp.Diagnostics, *updatedRule)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.client.MutexMap[api.MutexMapKeyCustomRules].Lock()
	defer r.client.MutexMap[api.MutexMapKeyCustomRules].Unlock()

	// Retrieve values from state
	var state models.CustomRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing policy
	err := ruleAPI.DeleteCustomRule(*r.client, int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Custom Rule resource",
			fmt.Sprintf("Failed to delete custom rule: %s", err.Error()),
		)
		return
	}
}

//func schemaToTerraform(ctx context.Context, schema models.CustomRuleResourceModel) (ruleAPI.CustomRule, diag.Diagnostics) {
//	util.HCLogDebug(ctx, "Executing CustomRuleResource.schemaToTerraform")
//
//	var (
//		diags            diag.Diagnostics
//		attackTechniques []string
//		vulnIDs          []string
//	)
//
//	diags = schema.AttackTechniques.ElementsAs(ctx, &attackTechniques, false)
//	if diags.HasError() {
//		return ruleAPI.CustomRule{}, diags
//	}
//
//	//diags = schema.VulnIDs.ElementsAs(ctx, &vulnIDs, false)
//	//if diags.HasError() {
//	//    return ruleAPI.CustomRule{}, diags
//	//}
//
//	resp := ruleAPI.CustomRule{
//		Id:               int(schema.Id.ValueInt64()),
//		AttackTechniques: attackTechniques,
//		Description:      schema.Description.ValueString(),
//		Message:          schema.Message.ValueString(),
//		MinVersion:       schema.MinVersion.ValueString(),
//		Modified:         int(schema.Modified.ValueInt64()),
//		Name:             schema.Name.ValueString(),
//		Owner:            schema.Owner.ValueString(),
//		Script:           schema.Script.ValueString(),
//		Type:             schema.Type.ValueString(),
//		VulnIDs:          vulnIDs,
//	}
//
//	util.HCLogDebug(ctx, "Finishing CustomRuleResource.schemaToTerraform execution")
//
//	return resp, diags
//}
//
//func terraformToSchema(ctx context.Context, tf ruleAPI.CustomRule) (models.CustomRuleResourceModel, diag.Diagnostics) {
//	util.HCLogDebug(ctx, "Executing CustomRuleResource.terraformToSchema")
//
//	var diags diag.Diagnostics
//
//	attackTechniques, diags := types.ListValueFrom(ctx, types.StringType, tf.AttackTechniques)
//	if diags.HasError() {
//		return models.CustomRuleResourceModel{}, diags
//	}
//
//	vulnIDs, diags := types.ListValueFrom(ctx, types.StringType, tf.VulnIDs)
//	if diags.HasError() {
//		return models.CustomRuleResourceModel{}, diags
//	}
//
//	resp := models.CustomRuleResourceModel{
//		Id:               types.Int64Value(int64(tf.Id)),
//		AttackTechniques: attackTechniques,
//		Description:      types.StringValue(tf.Description),
//		Message:          types.StringValue(tf.Message),
//		MinVersion:       types.StringValue(tf.MinVersion),
//		Modified:         types.Int64Value(int64(tf.Modified)),
//		Name:             types.StringValue(tf.Name),
//		Owner:            types.StringValue(tf.Owner),
//		Script:           types.StringValue(tf.Script),
//		Type:             types.StringValue(tf.Type),
//		VulnIDs:          vulnIDs,
//	}
//
//	util.HCLogDebug(ctx, "Finishing CustomRuleResource.terraformToSchema execution")
//
//	return resp, diags
//}
