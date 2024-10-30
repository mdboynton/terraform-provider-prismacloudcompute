package policy 

import (
    "context"
	"fmt"
    "time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    //"github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *CustomComplianceCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_custom_compliance_check"
}

func (r *CustomComplianceCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema()
}

func (r *CustomComplianceCheckResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomComplianceCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    util.DLog(ctx, "retrieving plan and serializing into CustomComplianceCheckResourceModel")
    var plan CustomComplianceCheckResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    check, diags := schemaToCheck(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Upsert custom compliance checks with new values
    util.DLog(ctx, fmt.Sprintf("creating custom check resource with payload:\n\n %+v", check))
    createdCheck, err := policyAPI.UpsertCustomComplianceCheck(*r.client, check)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Custom Compliance Check resource", 
            fmt.Sprintf("Failed to create custom compliance check: %s", err.Error()),
        )
        return
	}

    // Convert created check to schema
    createdCheckSchema, diags := checkToSchema(ctx, *createdCheck)
    if diags.HasError() {
        return
    }

    // Set state to created check data
    diags = resp.State.Set(ctx, createdCheckSchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *CustomComplianceCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    util.DLog(ctx, "starting Read() execution")

    // Get current state
    var state CustomComplianceCheckResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get policy value from Prisma Cloud
    check, err := policyAPI.GetCustomComplianceCheckById(*r.client, int(state.Id.ValueInt32()))
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Custom Compliance Check resource", 
            fmt.Sprintf("Failed to read custom compliance check: %s", err.Error()),
        )
        return
    }

    //util.DLog(ctx, fmt.Sprintf("retrieved container compliance policy with rules:\n\n %+v", *policy.Rules))
  
    // Overwrite state values with Prisma Cloud data
    checkSchema, diags := checkToSchema(ctx, *check)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    //util.DLog(ctx, fmt.Sprintf("policy schema rules:\n\n %+v", policySchema.Rules))

    // Set refreshed state
    diags = resp.State.Set(ctx, &checkSchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, "ending Read() execution")
}

func (r *CustomComplianceCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state CustomComplianceCheckResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan CustomComplianceCheckResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    check, diags := schemaToCheck(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Update existing policy
    updatedCheck, err := policyAPI.UpsertCustomComplianceCheck(*r.client, check)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating Custom Compliance Check resource", 
            fmt.Sprintf("Failed to update custom compliance check: %s", err.Error()),
        )
        return
	}

    // Convert updated policy into schema
    checkSchema, diags := checkToSchema(ctx, *updatedCheck)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    //util.DLog(ctx, fmt.Sprintf("setting state from Update() with rules:\n\n %+v", policySchema.Rules))

    // Set updated state
    diags = resp.State.Set(ctx, checkSchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *CustomComplianceCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve value from state
	var state CustomComplianceCheckResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    check, diags := schemaToCheck(ctx, &state, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing check 
    err := policyAPI.DeleteCustomComplianceCheck(*r.client, check.Id)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Container Compliance Policy resource", 
            "Failed to delete container compliance policy: " + err.Error(),
        )
        return
	}
}

func (r *CustomComplianceCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    util.DLog(ctx, "executing ImportState")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CustomComplianceCheckResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyPlan")

    //var plan *CustomComplianceCheckResourceModel
    //diags := req.Plan.Get(ctx, &plan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //ModifyCustomComplianceCheckResourcePlan(ctx, r.client, plan, resp)

    util.DLog(ctx, "exiting ModifyPlan")
}

func schemaToCheck(ctx context.Context, plan *CustomComplianceCheckResourceModel, client *api.PrismaCloudComputeAPIClient) (policyAPI.CustomComplianceCheck, diag.Diagnostics) {
    util.DLog(ctx, "entering schemaToCheck")

    var diags diag.Diagnostics

    check := policyAPI.CustomComplianceCheck{
        Owner: plan.Owner.ValueString(),
        Modified: time.Now().Format("2006-01-02T15:04:05.000Z"),
        Name: plan.Name.ValueString(),
        Script: plan.Script.ValueString(),
        Severity: plan.Severity.ValueString(),
        Title: plan.Description.ValueString(),
    }

    if !plan.Id.IsUnknown() && !plan.Id.IsNull() {
        check.Id = int(plan.Id.ValueInt32())
    }

    util.DLog(ctx, "exiting schemaToCheck")

    return check, diags
}

func checkToSchema(ctx context.Context, check policyAPI.CustomComplianceCheck) (CustomComplianceCheckResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "entering checkToSchema")

    var diags diag.Diagnostics
    
    schemaCheck := CustomComplianceCheckResourceModel{
        Id: types.Int32Value(int32(check.Id)),
        Modified: types.StringValue(""),
        Description: types.StringValue(check.Title),
        Owner: types.StringValue(check.Owner),
        Name: types.StringValue(check.Name),
        PreviousName: types.StringValue(check.PreviousName),
        Script: types.StringValue(check.Script),
        Severity: types.StringValue(check.Severity),
    }

    util.DLogf(ctx, schemaCheck)

    util.DLog(ctx, "exiting checkToSchema")

    return schemaCheck, diags
}
