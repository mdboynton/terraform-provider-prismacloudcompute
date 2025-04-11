package policy

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *CiImageCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ci_image_compliance_policy"
}

func (r *CiImageCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.GetSchema(ctx)
}

func (r *CiImageCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CiImageCompliancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.PolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	data, diags := policy.PolicySchemaToTerraform(ctx, &plan, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new CI image compliance policy
	err := policyAPI.UpsertPolicy(*r.client, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating CI Image Compliance Policy resource",
			"Failed to create CI image compliance policy: "+err.Error(),
		)
		return
	}

	// Retrieve newly created container compliance policy
	response, err := policyAPI.GetPolicy(*r.client, policyAPI.PolicyTypeComplianceCiImage)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving created CI Image Compliance Policy resource",
			"Failed to retrieve created CI image compliance policy: "+err.Error(),
		)
		return
	}

	createdPolicy, diags := policy.PolicyTerraformToSchema(ctx, *response, plan)
	if diags.HasError() {
		return
	}

	// Set state to collection data
	diags = resp.State.Set(ctx, createdPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CiImageCompliancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.PolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get policy value from Prisma Cloud
	data, err := policyAPI.GetPolicy(*r.client, policyAPI.PolicyTypeComplianceCiImage)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading CI Image Compliance Policy resource",
			"Failed to read CI image compliance policy: "+err.Error(),
		)
		return
	}

	// Overwrite state values with Prisma Cloud data
	policySchema, diags := policy.PolicyTerraformToSchema(ctx, *data, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &policySchema)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CiImageCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state models.PolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.PolicyResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	planPolicy, diags := policy.PolicySchemaToTerraform(ctx, &plan, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing policy
	err := policyAPI.UpsertPolicy(*r.client, planPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating CI Image Compliance Policy resource",
			"Failed to update CI image compliance policy: "+err.Error(),
		)
		return
	}

	// Get updated policy value from Prisma Cloud
	data, err := policyAPI.GetPolicy(*r.client, policyAPI.PolicyTypeComplianceCiImage)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading CI Image Compliance Policy resource",
			"Failed to read CI image compliance policy: "+err.Error(),
		)
		return
	}

	// Convert updated policy into schema
	policySchema, diags := policy.PolicyTerraformToSchema(ctx, *data, plan)
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

func (r *CiImageCompliancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.PolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Clear policy rules
	state.Rules = &[]models.PolicyRuleResourceModel{}

	// Generate API request body from plan
	updatedPlan, diags := policy.PolicySchemaToTerraform(ctx, &state, r.client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing policy
	err := policyAPI.UpsertPolicy(*r.client, updatedPlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting CI Image Compliance Policy resource",
			"Failed to delete CI image compliance policy: "+err.Error(),
		)
		return
	}
}

func (r *CiImageCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	util.LogDebug(ctx, "executing ImportState")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CiImageCompliancePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	util.LogDebug(ctx, "entering ModifyPlan")

	policy.ModifyPolicyResourcePlan(ctx, r.client, req.Plan, resp)

	util.LogDebug(ctx, "exiting ModifyPlan")
}
