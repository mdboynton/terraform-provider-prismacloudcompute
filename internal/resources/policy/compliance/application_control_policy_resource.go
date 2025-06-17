package policy 

import (
    "context"
	"fmt"
    "slices"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"

    //"github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *ApplicationControlPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_application_control_policy"
}

func (r *ApplicationControlPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema()
}

func (r *ApplicationControlPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationControlPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan models.ApplicationControlPolicyResourceModel 
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    policy, diags := plan.ToTerraform(ctx)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new application control policy 
    for _, policyRule := range policy {
        err := policyAPI.UpsertApplicationControlPolicyRule(*r.client, policyRule)
	    if err != nil {
	    	resp.Diagnostics.AddError(
                "Error creating Application Control Policy resource", 
                "Failed to create application control policy rule: " + err.Error(),
            )
	    }
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve newly created application control policy 
    response, err := policyAPI.GetApplicationControlPolicy(*r.client)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created Application Control Policy resource", 
            "Failed to retrieve created application control policy: " + err.Error(),
        )
        return
    }

    //createdPolicy, diags := policyToSchema(ctx, *response)
    createdPolicy := models.ApplicationControlPolicyResourceModel{}
    diags = createdPolicy.FromTerraform(ctx, response)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set state to collection data
    diags = resp.State.Set(ctx, createdPolicy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *ApplicationControlPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state models.ApplicationControlPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get policy value from Prisma Cloud
    policy, err := policyAPI.GetApplicationControlPolicy(*r.client)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Application Control Policy resource", 
            "Failed to read application control policy: " + err.Error(),
        )
        return
    }

    // Convert policy to schema
    policySchema := models.ApplicationControlPolicyResourceModel{}
    diags = policySchema.FromTerraform(ctx, policy)
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

func (r *ApplicationControlPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state models.ApplicationControlPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan models.ApplicationControlPolicyResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    planPolicy, diags := plan.ToTerraform(ctx)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Update existing policy
    stateRuleIDs := state.GetRuleIDs(ctx)
   
    // If there's no rules in the plan and >0 rules in the state, rules are only being deleted
    if (len(planPolicy) == 0 && len(stateRuleIDs) > 0) {
        for _, stateRule := range *state.Rules {
            err := policyAPI.DeleteApplicationControlPolicyRule(*r.client, int(stateRule.Id.ValueInt32()))
	        if err != nil {
	        	resp.Diagnostics.AddError(
                    "Error updating Application Control Policy resource", 
                    "Failed to delete application control policy rule during update: " + err.Error(),
                )
	        }
        }
    // Otherwise, update and delete rules according to the difference between plan and state
    } else {
        for i := 0; i < len(planPolicy); i++ {
            // If the rule ID is in the state but not in the plan, delete the rule
            if !slices.Contains(stateRuleIDs, planPolicy[i].Name) {
                err := policyAPI.DeleteApplicationControlPolicyRule(*r.client, planPolicy[i].Id)
	            if err != nil {
	            	resp.Diagnostics.AddError(
                        "Error updating Application Control Policy resource", 
                        "Failed to delete application control policy rule during update: " + err.Error(),
                    )
	            }
            // Otherwise, upsert the rule
            } else {
                err := policyAPI.UpsertApplicationControlPolicyRule(*r.client, planPolicy[i])
	            if err != nil {
	            	resp.Diagnostics.AddError(
                        "Error updating Application Control Policy resource", 
                        "Failed to update application control policy rule: " + err.Error(),
                    )
	            }
            }
        }
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Get updated policy value from Prisma Cloud
    response, err := policyAPI.GetApplicationControlPolicy(*r.client)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving updated Application Control Policy resource", 
            "Failed to retrieve updated application control policy: " + err.Error(),
        )
        return
    }

    // Convert updated policy to schema
    createdPolicy := models.ApplicationControlPolicyResourceModel{}
    diags = createdPolicy.FromTerraform(ctx, response)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set updated state
    diags = resp.State.Set(ctx, createdPolicy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *ApplicationControlPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state models.ApplicationControlPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Convert state to API struct
    updatedPlan, diags := state.ToTerraform(ctx)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete all rules
    for _, policyRule := range updatedPlan {
        err := policyAPI.DeleteApplicationControlPolicyRule(*r.client, policyRule.Id)
	    if err != nil {
	    	resp.Diagnostics.AddError(
                "Error deleting Application Control Policy resource", 
                "Failed to delete application control policy rule: " + err.Error(),
            )
            return
	    }
    }
    
    // Clear policy rules
    state.Rules = &[]models.ApplicationControlPolicyRuleResourceModel{}
}

func (r *ApplicationControlPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
