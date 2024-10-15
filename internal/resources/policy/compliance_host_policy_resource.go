package policy 

import (
    "context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	//collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    //"github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *HostCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_host_compliance_policy"
}

func (r *HostCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema()
}

func (r *HostCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *HostCompliancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // TODO: refine this logic to populate Owner with the value in config, if it exists
    //var username types.String
    //diags := req.Config.GetAttribute(ctx, path.Root("username"), &username)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    // Retrieve values from plan
    util.DLog(ctx, "retrieving plan and serializing into CompliancePolicyResourceModel")
    var plan CompliancePolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    policy, diags := CompliancePolicySchemaToPolicy(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new host compliance policy 
    util.DLog(ctx, fmt.Sprintf("creating policy resource with payload:\n\n %+v", *policy.Rules))
    err := policyAPI.UpsertHostCompliancePolicy(*r.client, policy)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Host Compliance Policy resource", 
            "Failed to create host compliance policy: " + err.Error(),
        )
        return
	}

    // Retrieve newly created host compliance policy 
    response, err := policyAPI.GetHostCompliancePolicy(*r.client)
    if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created Host Compliance Policy resource", 
            "Failed to retrieve created host compliance policy: " + err.Error(),
        )
        return
    }


    // TODO: explore passing in the CreateRequest to hostCompliancePolicyToSchema in order to be
    // able to reference configured order values that arent returned from the API

    createdPolicy, diags := CompliancePolicyToSchema(ctx, *response, plan)
    if diags.HasError() {
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

func (r *HostCompliancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    util.DLog(ctx, "starting Read() execution")

    // Get current state
    var state CompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get policy value from Prisma Cloud
    policy, err := policyAPI.GetHostCompliancePolicy(*r.client)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Host Compliance Policy resource", 
            "Failed to read host compliance Policy: " + err.Error(),
        )
        return
    }

    util.DLog(ctx, fmt.Sprintf("retrieved host compliance policy with rules:\n\n %+v", *policy.Rules))
  
    // Overwrite state values with Prisma Cloud data
    policySchema, diags := CompliancePolicyToSchema(ctx, *policy, state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, fmt.Sprintf("policy schema rules:\n\n %+v", policySchema.Rules))

    // Set refreshed state
    diags = resp.State.Set(ctx, &policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, "ending Read() execution")
}

func (r *HostCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state CompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan CompliancePolicyResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    planPolicy, diags := CompliancePolicySchemaToPolicy(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Update existing policy
    err := policyAPI.UpsertHostCompliancePolicy(*r.client, planPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating Host Compliance Policy resource", 
            "Failed to update host compliance policy: " + err.Error(),
        )
        return
	}

    // Get updated policy value from Prisma Cloud
    policy, err := policyAPI.GetHostCompliancePolicy(*r.client)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Host Compliance Policy resource", 
            "Failed to read Host Compliance Policy: " + err.Error(),
        )
        return
    }

    util.DLog(ctx, fmt.Sprintf("retrieved host compliance policy during Update() execution with rules:\n\n %+v", *policy.Rules))
  
    // Convert updated policy into schema
    policySchema, diags := CompliancePolicyToSchema(ctx, *policy, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, fmt.Sprintf("setting state from Update() with rules:\n\n %+v", policySchema.Rules))

    // Set updated state
    diags = resp.State.Set(ctx, policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *HostCompliancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state CompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Clear policy rules
    state.Rules = &[]CompliancePolicyRuleResourceModel{}

    // Generate API request body from plan
    updatedPlan, diags := CompliancePolicySchemaToPolicy(ctx, &state, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing policy 
    err := policyAPI.UpsertHostCompliancePolicy(*r.client, updatedPlan)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Host Compliance Policy resource", 
            "Failed to delete host compliance policy: " + err.Error(),
        )
        return
	}
}

func (r *HostCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    util.DLog(ctx, "executing ImportState")
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *HostCompliancePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyPlan")
    //util.DLog(ctx, fmt.Sprintf("%v+", resp))
    //util.DLog(ctx, fmt.Sprintf("%v+", req))

    var plan *CompliancePolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    if plan == nil {
        return
    }

    if plan.Rules == nil {
        emptyRules := []CompliancePolicyRuleResourceModel{}
        diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules"), &emptyRules)...)
        return
    }

    complianceVulnerabilities, err := systemAPI.GetComplianceVulnerabilities(*r.client, plan.PolicyType.ValueString())
	if err != nil {
		diags.AddError(
            "Error modifying planned policy rules", 
            "Failed to retrieve compliance host vulnerabilities from Prisma Cloud while modifying plan rules: " + err.Error(),
        )
        return
	}

    util.DLog(ctx, "starting loop over rules")

    for index, rule := range *plan.Rules {
        if len(rule.Collections.Elements()) == 0 {
            collections := types.ListValueMust(
                system.CollectionObjectType(),
                []attr.Value{
                    types.ObjectValueMust(
                        system.CollectionObjectAttrTypeMap(),
                        system.CollectionObjectDefaultAttrValueMap(),
                    ),
                },
            )
            diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("collections"), collections)...)
        }

        if rule.Order.IsUnknown() || rule.Order.IsNull() {
            rule.Order = types.Int32Value(int32(index + 1))
            diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("order"), types.Int32Value(int32(index + 1)))...)
            if diags.HasError() {
                return
            }
        } else if int(rule.Order.ValueInt32()) < 1 {
            resp.Diagnostics.AddError(
		    	"Invalid Resource Configuration",
		    	fmt.Sprintf("Host Compliance Policy Rule specified an invalid order (%d). Order values must be positive non-zero integers.", int(rule.Order.ValueInt32())),
		    )
            return
        }

        if rule.Effect.IsUnknown() {
            fmt.Printf("rule %s has unknown effect\n", rule.Name.ValueString())
            rule.Effect = types.StringValue("unknown")
            diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("effect"), types.StringValue("alert"))...)
            if diags.HasError() {
                return
            }
        } 

        conditionObject, diags := GenerateConditionFromEffect(ctx, *r.client, plan.PolicyType.ValueString(), rule, complianceVulnerabilities)
        if diags.HasError() {
            return 
        }

        resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)
    }

    util.DLog(ctx, "exiting ModifyPlan")
}
