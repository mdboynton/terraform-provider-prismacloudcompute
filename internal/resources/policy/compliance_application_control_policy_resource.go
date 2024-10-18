package policy 

import (
    "context"
	"fmt"
    "time"
    "reflect"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    //"github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
    policyAPI.DeleteApplicationControlPolicyRule(*r.client, policyAPI.ApplicationControlPolicyRule{Id: 11000})
    _ = reflect.TypeOf(ctx)

    // Retrieve values from plan
    util.DLog(ctx, "retrieving plan and serializing into ApplicationControlPolicyResourceModel")
    var plan ApplicationControlPolicyResourceModel 
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    policy, diags := schemaToPolicy(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new application control policy 
    util.DLog(ctx, fmt.Sprintf("creating application control policy resource with payload:\n\n %+v", policy))
    for _, policyRule := range policy {
        // TODO: if there's a failed insert, continue adding the other rules and return diags at the end
        // (check to see if this is the desired behaviour)
        err := policyAPI.UpsertApplicationControlPolicyRule(*r.client, policyRule)
	    if err != nil {
	    	resp.Diagnostics.AddError(
                "Error creating Application Control Policy resource", 
                "Failed to create application control policy rule: " + err.Error(),
            )
            return
	    }
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

    createdPolicy, diags := policyToSchema(ctx, *response)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, fmt.Sprintf("created policy with rules:\n\n %+v", *createdPolicy.Rules))
    
    // Set state to collection data
    diags = resp.State.Set(ctx, createdPolicy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *ApplicationControlPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    util.DLog(ctx, "starting Read() execution")

    // Get current state
    var state ApplicationControlPolicyResourceModel 
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

    //util.DLog(ctx, fmt.Sprintf("retrieved application control policy with rules:\n\n %+v", *policy.Rules))
  
    // Convert policy to terraform schema
    policySchema, diags := policyToSchema(ctx, *policy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    //util.DLog(ctx, fmt.Sprintf("policy schema rules:\n\n %+v", policySchema.Rules))

    // Set refreshed state
    diags = resp.State.Set(ctx, &policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, "ending Read() execution")
}

func (r *ApplicationControlPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state ApplicationControlPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan ApplicationControlPolicyResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    //for index, stateRule := range *state.Rules {
    //    if !stateRule.Equals(ctx, (*plan.Rules)[index]) {
    //        util.DLog(ctx, "rule mismatch")
    //    }
    //}

    // Generate API request body from plan
    planPolicy, diags := schemaToPolicy(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, fmt.Sprintf("updating policy with payload: \n\n%v", planPolicy))

    // Update existing policy
    for _, policyRule := range planPolicy {
        // TODO: if there's a failed insert, continue adding the other rules and return diags at the end
        // (check to see if this is the desired behaviour)
        err := policyAPI.UpsertApplicationControlPolicyRule(*r.client, policyRule)
	    if err != nil {
	    	resp.Diagnostics.AddError(
                "Error creating Application Control Policy resource", 
                "Failed to create application control policy rule: " + err.Error(),
            )
            return
	    }
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
    createdPolicy, diags := policyToSchema(ctx, *response)
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
	var state ApplicationControlPolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Convert state to API struct
    updatedPlan, diags := schemaToPolicy(ctx, &state, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete all rules
    for _, policyRule := range updatedPlan {
        err := policyAPI.DeleteApplicationControlPolicyRule(*r.client, policyRule)
	    if err != nil {
	    	resp.Diagnostics.AddError(
                "Error deleting Application Control Policy resource", 
                "Failed to delete application control policy rule: " + err.Error(),
            )
            return
	    }
    }
    
    // Clear policy rules
    state.Rules = &[]ApplicationControlPolicyRuleResourceModel{}
}

func (r *ApplicationControlPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    util.DLog(ctx, "executing ImportState")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//func (r *ApplicationControlPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
//    util.DLog(ctx, "entering ModifyPlan")
//
//    var plan *ApplicationControlPolicyResourceModel
//    diags := req.Plan.Get(ctx, &plan)
//    resp.Diagnostics.Append(diags...)
//    if resp.Diagnostics.HasError() {
//        return
//    }
//
//    ModifyCompliancePolicyResourcePlan(ctx, r.client, plan, resp)
//
//    util.DLog(ctx, "exiting ModifyPlan")
//}


func schemaToPolicy(ctx context.Context, plan *ApplicationControlPolicyResourceModel, client *api.PrismaCloudComputeAPIClient,/*, username types.String*/) ([]policyAPI.ApplicationControlPolicyRule, diag.Diagnostics) {
    util.DLog(ctx, "entering schemaToPolicy")
    var diags diag.Diagnostics

    policy := []policyAPI.ApplicationControlPolicyRule{}

    if plan.Rules == nil {
        return policy, diags
    }

    for _, planRule := range *plan.Rules {
        util.DLogf(ctx, planRule)
        rule := policyAPI.ApplicationControlPolicyRule{
            Id: int(planRule.Id.ValueInt32()),
            Name: planRule.Name.ValueString(),
            Description: planRule.Description.ValueString(),
            Notes: planRule.Notes.ValueString(),
            Modified: time.Now().Format("2006-01-02T15:04:05.000Z"),
            //Owner: planRule.Notes.
            PreviousName: planRule.PreviousName.ValueString(),
            Severity: planRule.Severity.ValueString(),
        }

        planApplications := []policyAPI.ApplicationControlPolicyRuleApplication{}
        diags = planRule.Applications.ElementsAs(ctx, &planApplications, false)
        if diags.HasError() {
            return policy, diags
        }

        rule.Applications = planApplications

        policy = append(policy, rule)
    }

    util.DLog(ctx, "exiting schemaToPolicy")
    return policy, diags
}

func policyToSchema(ctx context.Context, rules []policyAPI.ApplicationControlPolicyRule) (ApplicationControlPolicyResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "entering policyToSchema")
    var diags diag.Diagnostics
    
    schemaPolicy := ApplicationControlPolicyResourceModel{}
    schemaRules := []ApplicationControlPolicyRuleResourceModel{}
    
    for _, rule := range rules {
        schemaRule, diags := policyRuleToSchema(ctx, rule)
        if diags.HasError() {
            return schemaPolicy, diags
        }

        schemaRules = append(schemaRules, schemaRule) 
    }

    schemaPolicy.Rules = &schemaRules
    util.DLogf(ctx, schemaPolicy)
    
    util.DLog(ctx, "exiting policyToSchema")
    return schemaPolicy, diags
}

func policyRuleToSchema(ctx context.Context, rule policyAPI.ApplicationControlPolicyRule) (ApplicationControlPolicyRuleResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "entering policyRuleToSchema")
    //util.DLogf(ctx, rule)
    var diags diag.Diagnostics

    schemaRule := ApplicationControlPolicyRuleResourceModel{
        Id: types.Int32Value(int32(rule.Id)),
        Name: types.StringValue(rule.Name),
        Description: types.StringValue(rule.Description),
        Modified: types.StringValue(""),
        Notes: types.StringValue(rule.Notes),
        Owner: types.StringValue(rule.Owner),
        PreviousName: types.StringValue(rule.PreviousName),
        Severity: types.StringValue(rule.Severity),
    }

    allowedVersionsConditionsSetType := types.SetType{
        ElemType: types.StringType,
    }

    applicationsSetType := types.ObjectType{
        AttrTypes: map[string]attr.Type{
            "name": types.StringType,
            "allowed_versions": types.SetType{
                ElemType: types.SetType{
                    ElemType: types.StringType,
                },
            },
        },
    }

    applicationsTypeMap := map[string]attr.Type{
        "name": types.StringType,
        "allowed_versions": types.SetType{
            ElemType: types.SetType{
                ElemType: types.StringType,
            },
        },
    }

    schemaApplications := []attr.Value{}
    for _, application := range rule.Applications {
        schemaAllowedVersions := []attr.Value{}
        for _, allowedVersion := range application.AllowedVersions {
            allowedVersionSet, diags := types.SetValueFrom(ctx, types.StringType, allowedVersion)
            if diags.HasError() {
                return schemaRule, diags
            }

            schemaAllowedVersions = append(schemaAllowedVersions, allowedVersionSet)
        }
    
        allowedVersionsSet, diags := types.SetValueFrom(ctx, allowedVersionsConditionsSetType, schemaAllowedVersions)
        if diags.HasError() {
            return schemaRule, diags
        }

        schemaApplication := types.ObjectValueMust(
            applicationsTypeMap, 
            map[string]attr.Value{
                "name": types.StringValue(application.Name),
                "allowed_versions": allowedVersionsSet,
            },
        )

        schemaApplications = append(schemaApplications, schemaApplication)
    }

    util.DLogf(ctx, schemaApplications)

    schemaApplicationsSet, diags := types.SetValueFrom(ctx, applicationsSetType, schemaApplications)

    if diags.HasError() {
        util.DLogf(ctx, diags)
        return schemaRule, diags
    }

    schemaRule.Applications = schemaApplicationsSet

    util.DLog(ctx, "exiting policyRuleToSchema")
    return schemaRule, diags
}


func (r *ApplicationControlPolicyRuleResourceModel) Equals(ctx context.Context, compare ApplicationControlPolicyRuleResourceModel) bool {
    if (*r).Description.ValueString() != compare.Description.ValueString() {
        return false
    }

    return true
}
