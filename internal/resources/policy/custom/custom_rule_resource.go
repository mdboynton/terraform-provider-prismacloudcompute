package custom

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	ruleAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/rule"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &CustomRuntimeRuleResource{}
var _ resource.ResourceWithImportState = &CustomRuntimeRuleResource{}

func (r *CustomRuntimeRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_runtime_rule"
}

func (r *CustomRuntimeRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.GetSchema(ctx)
}

func (r *CustomRuntimeRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomRuntimeRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.CustomRuntimeRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	data, diags := schemaToTerraform(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new custom runtime rule
	createdRuleId, err := ruleAPI.CreateCustomRule(*r.client, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Custom Runtime Rule resource",
			"Failed to create custom runtime rule: "+err.Error(),
		)
		return
	}

	// Retrieve newly created custom runtime rule
	response, err := ruleAPI.GetCustomRuleById(*r.client, createdRuleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving created Custom Runtime Rule resource",
			"Failed to retrieve created custom runtime rule: "+err.Error(),
		)
		return
	}

	createdRule, diags := terraformToSchema(ctx, *response)
	if diags.HasError() {
		return
	}

	// Set state to collection data
	diags = resp.State.Set(ctx, createdRule)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomRuntimeRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.CustomRuntimeRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get policy value from Prisma Cloud
	data, err := ruleAPI.GetCustomRuleById(*r.client, int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Custom Runtime Rule resource",
			"Failed to read custom runtime rule: "+err.Error(),
		)
		return
	}

	// Overwrite state values with Prisma Cloud data
	ruleSchema, diags := terraformToSchema(ctx, *data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &ruleSchema)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomRuntimeRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state models.CustomRuntimeRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.CustomRuntimeRuleResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	planRule, diags := schemaToTerraform(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing rule
	err := ruleAPI.UpdateCustomRule(*r.client, planRule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Custom Runtime Rule resource",
			"Failed to update custom runtime rule: "+err.Error(),
		)
		return
	}

	// Get updated rule value from Prisma Cloud
	updatedRule, err := ruleAPI.GetCustomRuleById(*r.client, int(plan.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Custom Runtime Rule resource",
			"Failed to read Custom Runtime Rule: "+err.Error(),
		)
		return
	}

	// Convert updated rule into schema
	policySchema, diags := terraformToSchema(ctx, *updatedRule)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set updated state
	diags = resp.State.Set(ctx, policySchema)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomRuntimeRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.CustomRuntimeRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing policy
	err := ruleAPI.DeleteCustomRule(*r.client, int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Custom Runtime Rule resource",
			"Failed to delete custom runtime rule: "+err.Error(),
		)
		return
	}
}

func (r *CustomRuntimeRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CustomRuntimeRuleResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	util.LogDebug(ctx, "entering ModifyPlan")

	//policy.ModifyPolicyResourcePlan(ctx, r.client, req.Plan, resp)

	util.LogDebug(ctx, "exiting ModifyPlan")
}

func schemaToTerraform(ctx context.Context, schema models.CustomRuntimeRuleResourceModel) (ruleAPI.CustomRule, diag.Diagnostics) {
	util.LogDebug(ctx, "Executing CustomRuleResource.schemaToTerraform")

	var (
		diags            diag.Diagnostics
		attackTechniques []string
		vulnIDs          []string
	)

	diags = schema.AttackTechniques.ElementsAs(ctx, &attackTechniques, false)
	if diags.HasError() {
		return ruleAPI.CustomRule{}, diags
	}

	//diags = schema.VulnIDs.ElementsAs(ctx, &vulnIDs, false)
	//if diags.HasError() {
	//    return ruleAPI.CustomRule{}, diags
	//}

	resp := ruleAPI.CustomRule{
		Id:               int(schema.Id.ValueInt64()),
		AttackTechniques: attackTechniques,
		Description:      schema.Description.ValueString(),
		Message:          schema.Message.ValueString(),
		MinVersion:       schema.MinVersion.ValueString(),
		Modified:         int(schema.Modified.ValueInt64()),
		Name:             schema.Name.ValueString(),
		Owner:            schema.Owner.ValueString(),
		Script:           schema.Script.ValueString(),
		Type:             schema.Type.ValueString(),
		VulnIDs:          vulnIDs,
	}

	util.LogDebug(ctx, "Finishing CustomRuleResource.schemaToTerraform execution")

	return resp, diags
}

func terraformToSchema(ctx context.Context, tf ruleAPI.CustomRule) (models.CustomRuntimeRuleResourceModel, diag.Diagnostics) {
	util.LogDebug(ctx, "Executing CustomRuleResource.terraformToSchema")

	var diags diag.Diagnostics

	attackTechniques, diags := types.ListValueFrom(ctx, types.StringType, tf.AttackTechniques)
	if diags.HasError() {
		return models.CustomRuntimeRuleResourceModel{}, diags
	}

	vulnIDs, diags := types.ListValueFrom(ctx, types.StringType, tf.VulnIDs)
	if diags.HasError() {
		return models.CustomRuntimeRuleResourceModel{}, diags
	}

	resp := models.CustomRuntimeRuleResourceModel{
		Id:               types.Int64Value(int64(tf.Id)),
		AttackTechniques: attackTechniques,
		Description:      types.StringValue(tf.Description),
		Message:          types.StringValue(tf.Message),
		MinVersion:       types.StringValue(tf.MinVersion),
		Modified:         types.Int64Value(int64(tf.Modified)),
		Name:             types.StringValue(tf.Name),
		Owner:            types.StringValue(tf.Owner),
		Script:           types.StringValue(tf.Script),
		Type:             types.StringValue(tf.Type),
		VulnIDs:          vulnIDs,
	}

	util.LogDebug(ctx, "Finishing CustomRuleResource.terraformToSchema execution")

	return resp, diags
}
