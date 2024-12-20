package policy 

import (
    "context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    //"github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ContainerCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_container_compliance_policy"
}

func (r *ContainerCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    util.DLog(ctx, "ContainerCompliancePolicyResource.Schema() called")
    resp.Schema = r.GetSchema(ctx)
}

func (r *ContainerCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    util.DLog(ctx, "ContainerCompliancePolicyResource.Configure() called")
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

func (r *ContainerCompliancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

    // Create new container compliance policy 
    err := policyAPI.UpsertPolicy(*r.client, data)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Container Compliance Policy resource", 
            "Failed to create container compliance policy: " + err.Error(),
        )
        return
	}

    // Retrieve newly created container compliance policy 
    response, err := policyAPI.GetPolicy(*r.client, policyAPI.PolicyTypeComplianceContainer)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created Container Compliance Policy resource", 
            "Failed to retrieve created container compliance policy: " + err.Error(),
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

func (r *ContainerCompliancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state models.PolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get policy value from Prisma Cloud
    data, err := policyAPI.GetPolicy(*r.client, policyAPI.PolicyTypeComplianceContainer)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Container Compliance Policy resource", 
            "Failed to read container compliance policy: " + err.Error(),
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

func (r *ContainerCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
            "Error updating Container Compliance Policy resource", 
            "Failed to update container compliance policy: " + err.Error(),
        )
        return
	}

    // Get updated policy value from Prisma Cloud
    updatedPolicy, err := policyAPI.GetPolicy(*r.client, policyAPI.PolicyTypeComplianceContainer)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Container Compliance Policy resource", 
            "Failed to read container compliance policy: " + err.Error(),
        )
        return
    }

    // Convert updated policy into schema
    policySchema, diags := policy.PolicyTerraformToSchema(ctx, *updatedPolicy, plan)
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

func (r *ContainerCompliancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
            "Error deleting Container Compliance Policy resource", 
            "Failed to delete container compliance policy: " + err.Error(),
        )
        return
	}
}

func (r *ContainerCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    util.DLog(ctx, "executing ImportState")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ContainerCompliancePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyPlan")

    policy.ModifyPolicyResourcePlan(ctx, r.client, req.Plan, resp)

    util.DLog(ctx, "exiting ModifyPlan")
}
