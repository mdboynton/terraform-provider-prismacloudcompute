package policy 

import (
    "context"
	"fmt"
    "reflect"
    "slices"
    "cmp"
    "sort"
    "time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
    //"github.com/hashicorp/terraform-plugin-log/tflog"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

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
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("retrieving plan and serializing into HostCompliancePolicyResourceModel")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    var plan HostCompliancePolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    //fmt.Println(*plan.Rules)
    //fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

    // Generate API request body from plan
    policy, diags := schemaToPolicy(ctx, &plan, r.client)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new host compliance policy 
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("creating policy resource with payload:")
    //fmt.Printf("%+v\n", policy)
    fmt.Printf("%+v\n", *policy.Rules)
    //r1 := *policy.Rules
    //fmt.Printf("%+v\n", r1)
    //fmt.Printf("%+v\n", *r1[0].Condition)
    //fmt.Println("number of vulns:")
    //fmt.Println(len(r1[0].Condition.Vulnerabilities))
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
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


    // TODO: explore passing in the CreateRequest to policyToSchema in order to be
    // able to reference configured order values that arent returned from the API

    //fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    ////fmt.Println(req.Plan)
    //fmt.Println(req.Plan.Raw)
    ////fmt.Println(reflect.TypeOf(req.Plan))
    //fmt.Println(reflect.TypeOf(req.Plan.Raw))
    //fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    //createdPolicy, diags := policyToSchema(ctx, *response)
    createdPolicy, diags := policyToSchema(ctx, *response, plan)
    if diags.HasError() {
        return
    }

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("createdPolicy in Create():")
    fmt.Println(reflect.TypeOf(createdPolicy))
    fmt.Printf("%+v\n", createdPolicy)
    //fmt.Printf("%+v\n", *createdPolicy.Rules)
    fmt.Printf("%+v\n", *createdPolicy.Rules)
    //fmt.Println("number of vulns:")
    //fmt.Println(len(r1[0].Condition.Vulnerabilities))
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    
    // Set state to collection data
    diags = resp.State.Set(ctx, createdPolicy)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%")
        fmt.Println("error in resp.State.Set")
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%")
        return
    }

}

func (r *HostCompliancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("we're in Read")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")

    // Get current state
    var state HostCompliancePolicyResourceModel 
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
            "Failed to read Host Compliance Policy: " + err.Error(),
        )
        return
    }

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("retrieved host compliance policy: ") 
    fmt.Printf("%+v\n", policy)
    fmt.Printf("%+v\n", *policy.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
  
    // Overwrite state values with Prisma Cloud data
    //createdPolicy, diags := policyToSchema(ctx, *response)
    policySchema, diags := policyToSchema(ctx, *policy, state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("policySchema: ") 
    fmt.Printf("%+v\n", policySchema)
    fmt.Printf("%+v\n", policySchema.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")

    // Set refreshed state
    diags = resp.State.Set(ctx, &policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *HostCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get current state
    var state HostCompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from plan
    var plan HostCompliancePolicyResourceModel 
    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    planPolicy, diags := schemaToPolicy(ctx, &plan, r.client)
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

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("retrieved host compliance policy during Update() execution: ") 
    fmt.Printf("%+v\n", policy)
    fmt.Printf("%+v\n", *policy.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
  
    // Convert updated policy into schema
    //createdPolicy, diags := policyToSchema(ctx, *response)
    policySchema, diags := policyToSchema(ctx, *policy, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    fmt.Println("(Update) setting state") 
    fmt.Printf("%+v\n", policySchema.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    // Set updated state
    diags = resp.State.Set(ctx, policySchema)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *HostCompliancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state HostCompliancePolicyResourceModel 
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Clear policy rules
    emptyRules := []HostCompliancePolicyRuleResourceModel{}
    state.Rules = &emptyRules

    // Generate API request body from plan
    updatedPlan, diags := schemaToPolicy(ctx, &state, r.client)
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

// TODO: Define ImportState to work properly with this resource
func (r *HostCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("executing ImportState")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func generateRulesOrderMap(rules []HostCompliancePolicyRuleResourceModel) map[string]int {
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("starting rule ordering")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    
    orderedRulesMap := make(map[int][]string)

    for _, rule := range rules {
        order := int(rule.Order.ValueInt32())
        if _, exists := orderedRulesMap[order]; exists {
            orderedRulesMap[order] = append(orderedRulesMap[order], rule.Name.ValueString())
        } else {
            orderedRulesMap[order] = []string{rule.Name.ValueString()}
        }
    }
        
    sortedKeys := make([]int, 0, len(orderedRulesMap))
    for key, _ := range orderedRulesMap {
        sortedKeys = append(sortedKeys, key)
    }
    sort.Ints(sortedKeys)

    ruleOrders := make(map[string]int)
    lastOrderValue := -1
    for _, key := range sortedKeys {
        offset := 0
        if lastOrderValue != -1 && lastOrderValue >= key {
            offset = lastOrderValue - key + 1
        }

        for sliceIndex, ruleName := range orderedRulesMap[key] {
            orderValue := key + sliceIndex + offset
            ruleOrders[ruleName] = orderValue
            lastOrderValue = orderValue
        }
    }

    return ruleOrders
}

func sortPolicyRules(rules *[]policyAPI.HostCompliancePolicyRule, planRules *[]HostCompliancePolicyRuleResourceModel) {
    rulesOrderMap := generateRulesOrderMap(*planRules) 
    r := *rules
    sort.Slice(r, func(i, j int) bool {
        return rulesOrderMap[r[i].Name] < rulesOrderMap[r[j].Name]
    })
    rules = &r 
}

func sortSchemaRules(schemaRules *[]HostCompliancePolicyRuleResourceModel, planRules *[]HostCompliancePolicyRuleResourceModel) {
    ruleOrderMap := make(map[string]int32)
    pr := *planRules
    sr := *schemaRules

    for index, planRule := range pr {
        ruleOrderMap[planRule.Name.ValueString()] = int32(index)
    }

    slices.SortFunc(sr, func(a, b HostCompliancePolicyRuleResourceModel) int {
        orderA, okA := ruleOrderMap[a.Name.ValueString()]
        if !okA {
            orderA = int32(len(ruleOrderMap) + 1)
        }
        orderB, okB := ruleOrderMap[b.Name.ValueString()]
        if !okB {
            orderB = int32(len(ruleOrderMap) + 1)
        }
        return cmp.Compare(orderA, orderB)
    })

    for i := 0; i < len(sr); i++ {
        sr[i].Order = pr[i].Order
    }

    schemaRules = &sr
}

func createConditionFromEffect(ctx context.Context, client api.PrismaCloudComputeAPIClient, rule HostCompliancePolicyRuleResourceModel, complianceVulnerabilities []systemAPI.Vulnerability) (basetypes.ObjectValue, diag.Diagnostics) {
    // TODO: finish implementing this more compact implementation of this function
    //      currently, the issue is that we're dealing with two different types of
    //      vulnerability objects depending on whether we're taking the vulnerability
    //      data from the TF resource configuration or from Prisma Cloud
    //if effect != "ignore" {
    //    if effect == "alert, block" {
    //        //var ruleConditionVulns []policyAPI.HostCompliancePolicyRuleVulnerability

    //        if rule.Condition.IsUnknown() {
    //            diags.AddError(
    //                "Missing condition from \"alert, block\" effect rule",
    //                "Condition attribute must be defined for rules with effect \"alert, block\".",
    //            )
    //            return conditionObject, diags
    //        }

    //        ruleCondition := policyAPI.HostCompliancePolicyRuleCondition{} 
    //        diags = rule.Condition.As(ctx, &ruleCondition, basetypes.ObjectAsOptions{})
    //        if diags.HasError() {
    //            return conditionObject, diags
    //        }

    //        //ruleConditionVulns = ruleCondition.Vulnerabilities
    //        //complianceVulnerabilities = ruleCondition.Vulnerabilities
    //        vulnerabilities := ruleCondition.Vulnerabilities
    //    } else if effect == "unknown" {
    //        complianceVulnerabilities = systemAPI.GetHighOrCriticalVulnerabilities(complianceVulnerabilities)
    //        vulnerabilities := complianceVulnerabilities
    //    } else {
    //        vulnerabilities := complianceVulnerabilities
    //    }
    //   
    //    //var block func(string, HostCompliancePolicyRuleVulnerabilityResourceModel) bool
    //    var block func(string, interface{}) bool
    //    //block = func(effect string, vuln HostCompliancePolicyRuleVulnerability) bool {
    //    block = func(effect string, vuln interface{}) bool {
    //        if effect == "alert, block" {
    //            return vuln.Block
    //        } else {
    //            return (effect == "block") && !(vuln.Type == "windows")
    //        }
    //    }

    //    //isBlockEffect := (effect == "block")
    //    for _, vuln := range complianceVulnerabilities {
    //        //if effect == "alert, block" {
    //        //    block := vuln.Block
    //        //} else {
    //        //    block := (isBlockEffect && !(vuln.Type == "windows"))
    //        //}

    //        vulnerabilityObjectValue := types.ObjectValueMust(
    //            vulnerabilitiesAttributeTypes,
    //            map[string]attr.Value{
    //                "id": types.Int32Value(int32(vuln.Id)),
    //                //"block": types.BoolValue(block),
    //                "block": types.BoolValue(block(effect, vuln)),
    //            },
    //        )
    //        
    //        vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
    //    }
    //}
    var diags diag.Diagnostics

    // Create static values
    effect := rule.Effect.ValueString()
    vulnerabilityObjectValues := []attr.Value{}
    vulnerabilitiesAttributeTypes := map[string]attr.Type{
        "id": types.Int32Type,
        "block": types.BoolType,
    }
    conditionObjectValueTypes := map[string]attr.Type{
        //"device": types.StringType,
        //"read_only": types.BoolType,
        "vulnerabilities": types.ListType{
            ElemType: types.ObjectType{
                AttrTypes: vulnerabilitiesAttributeTypes,
            },
        },
    }
    conditionObject := types.ObjectNull(conditionObjectValueTypes)

    // If the effect is "alert, block", create condition vulnerabilities object from plan
    if effect == "alert, block" {
        var ruleConditionVulns []policyAPI.HostCompliancePolicyRuleVulnerability

        if rule.Condition.IsUnknown() {
            diags.AddError(
                "Missing condition from \"alert, block\" effect rule",
                "Condition attribute must be defined for rules with effect \"alert, block\".",
            )
            return conditionObject, diags
        }

        ruleCondition := policyAPI.HostCompliancePolicyRuleCondition{} 
        diags = rule.Condition.As(ctx, &ruleCondition, basetypes.ObjectAsOptions{})
        if diags.HasError() {
            return conditionObject, diags
        }
        ruleConditionVulns = ruleCondition.Vulnerabilities

        for _, vuln := range ruleConditionVulns {
            vulnerabilityObjectValue := types.ObjectValueMust(
                vulnerabilitiesAttributeTypes,
                map[string]attr.Value{
                    "id": types.Int32Value(int32(vuln.Id)),
                    "block": types.BoolValue(vuln.Block),
                },
            )
            vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
        }
    // Otherwise, if the rule effect is not "ignore", create condition vulnerabilities using Prisma Cloud vulnerability data
    } else if effect != "ignore" {
        if effect == "unknown" {
           complianceVulnerabilities = systemAPI.GetHighOrCriticalVulnerabilities(complianceVulnerabilities)
        }

        isBlockEffect := (effect == "block")

        for _, vuln := range complianceVulnerabilities {
            block := (isBlockEffect && !(vuln.Type == "windows"))

            vulnerabilityObjectValue := types.ObjectValueMust(
                vulnerabilitiesAttributeTypes,
                map[string]attr.Value{
                    "id": types.Int32Value(int32(vuln.Id)),
                    "block": types.BoolValue(block),
                },
            )
            
            vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
        }
    }

    // Create vulnerability list value
    vulnerabilityObject, diags := types.ListValueFrom(
        ctx,
        types.ObjectType{
            AttrTypes: vulnerabilitiesAttributeTypes,
        },
        vulnerabilityObjectValues,
    )

    if diags.HasError() {
        return conditionObject, diags
    }

    // Create condition object value
    conditionObject = types.ObjectValueMust(
        conditionObjectValueTypes,
        map[string]attr.Value{
            //"device": types.StringValue(rule.Condition.Device),
            //"read_only": types.BoolValue(rule.Condition.ReadOnly),
            "vulnerabilities": vulnerabilityObject,
        },
    )

    return conditionObject, diags
}

func schemaToPolicy(ctx context.Context, plan *HostCompliancePolicyResourceModel, client *api.PrismaCloudComputeAPIClient,/*, username types.String*/) (policyAPI.HostCompliancePolicy, diag.Diagnostics) {
    var diags diag.Diagnostics

    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("entering schemaToPolicy")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    policy := policyAPI.HostCompliancePolicy{
        Id: plan.Id.ValueString(),
        PolicyType: plan.PolicyType.ValueString(),
    }

    if plan.Rules == nil {
        rules := []policyAPI.HostCompliancePolicyRule{}
        policy.Rules = &rules
        return policy, diags
    } else {
        rules, diags := ruleSchemaToPolicy(ctx, *plan.Rules, client)
        if diags.HasError() {
            return policy, diags
        }

        policy.Rules = &rules
    }

    return policy, diags
}

func ruleSchemaToPolicy(ctx context.Context, planRules []HostCompliancePolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient, /*, username types.String*/) ([]policyAPI.HostCompliancePolicyRule, diag.Diagnostics) {
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("entering ruleSchemaToPolicy")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    _ = reflect.TypeOf(ctx)

    var diags diag.Diagnostics

    rules := []policyAPI.HostCompliancePolicyRule{}

    for _, planRule := range planRules {
        collections := []collectionAPI.Collection{}
        diags = planRule.Collections.ElementsAs(ctx, &collections, false)
        if diags.HasError() {
            fmt.Println(diags)
            return rules, diags
        }

        for i, _ := range collections {
            collections[i].Modified = time.Now().Format("2006-01-02T15:04:05.000Z")
        }

        if planRule.Effect.ValueString() == "alert, block" && planRule.Condition.IsUnknown() {
            diags.AddError(
                "Missing condition from \"alert, block\" effect rule",
                "Condition attribute must be defined for rules with effect \"alert, block\".",
            )
            return rules, diags
        }

        condition := policyAPI.HostCompliancePolicyRuleCondition{} 
        diags = planRule.Condition.As(ctx, &condition, basetypes.ObjectAsOptions{})
        if diags.HasError() {
            return rules, diags
        }

        rule := policyAPI.HostCompliancePolicyRule{
            Order: int(planRule.Order.ValueInt32()),
            Name: planRule.Name.ValueString(), 
            Collections: collections,
            BlockMessage: planRule.BlockMessage.ValueString(),
            //Collections: col,
            Condition: &condition,
            Effect: planRule.Effect.ValueString(),
            Verbose: planRule.Verbose.ValueBool(),
            ReportAllPassedAndFailedChecks: planRule.ReportAllPassedAndFailedChecks.ValueBool(),
            //Owner: planRule.Owner.ValueString(),
            Disabled: planRule.Disabled.ValueBool(),
            Modified: time.Now().Format("2006-01-02T15:04:05.000Z"),
        }
        
        if !planRule.Notes.IsUnknown() && !planRule.Notes.IsNull() {
            rule.Notes = planRule.Notes.ValueString()
        }

        rules = append(rules, rule)
    }

    sortPolicyRules(&rules, &planRules)

    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("exiting ruleSchemaToPolicy")
    fmt.Printf("%+v\n", rules)
    //fmt.Printf("%+v\n", *rules[0].Condition)
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    
    return rules, diags
}

func policyToSchema(ctx context.Context, policy policyAPI.HostCompliancePolicy, plan HostCompliancePolicyResourceModel) (HostCompliancePolicyResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    schema := HostCompliancePolicyResourceModel{
        Id: types.StringValue(policy.Id),
        PolicyType: types.StringValue(policy.PolicyType),
    }

    rules, diags := policyRulesToSchema(ctx, *policy.Rules, *plan.Rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    //fmt.Println(plan)
    fmt.Println("!!! returned from policyRulesToSchema: ") 
    //fmt.Printf("%+v\n", rules)
    fmt.Println("#$%#$%#$%#$%#$%#$%#$%")
    if diags.HasError() {
        return schema, diags
    }

    schema.Rules = &rules

    return schema, diags
}

func policyRulesToSchema(ctx context.Context, rules []policyAPI.HostCompliancePolicyRule, planRules []HostCompliancePolicyRuleResourceModel) ([]HostCompliancePolicyRuleResourceModel, diag.Diagnostics) {
    fmt.Println("***********************")
    fmt.Println("entering policyRulesToSchema")
    fmt.Println("***********************")

    var diags diag.Diagnostics

    schemaRules := []HostCompliancePolicyRuleResourceModel{}

    if len(rules) == 0 {
        return nil, diags
    }

    for _, rule := range rules {
        schemaRule := HostCompliancePolicyRuleResourceModel{
            Name: types.StringValue(rule.Name),
            Disabled: types.BoolValue(rule.Disabled),
            Effect: types.StringValue(rule.Effect),
            Verbose: types.BoolValue(rule.Verbose),
            Owner: types.StringValue(rule.Owner),
            ReportAllPassedAndFailedChecks: types.BoolValue(rule.ReportAllPassedAndFailedChecks),
            Modified: types.StringValue(""),
        }

        if rule.Collections != nil {
            fmt.Println("***********************")
            fmt.Printf("entering collectionsToSchema (rule name: %s)\n", rule.Name)
            fmt.Println("***********************")
            collectionSet, diags := collectionsToSchema(ctx, rule.Collections)
            if diags.HasError() {
                return schemaRules, diags
            }

            schemaRule.Collections = collectionSet
        }

        if rule.Effect == "alert, block" {
            rule.Effect = "block" 
        }

        if rule.Condition != nil {
            vulnerabilityObjectValues := []attr.Value{}
            for _, vulnerability := range(rule.Condition.Vulnerabilities) {
                vulnerabilityObjectValue := types.ObjectValueMust(
                    map[string]attr.Type{
                        "id":        types.Int32Type,
                        "block":       types.BoolType,
                    },
                    map[string]attr.Value{
                        "id": types.Int32Value(int32(vulnerability.Id)),
                        "block": types.BoolValue(vulnerability.Block),
                    },
                )
               
                vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
            }

            vulnerabilityObject, diags := types.ListValueFrom(
                ctx,
                types.ObjectType{
                    AttrTypes: map[string]attr.Type{
                        "id": types.Int32Type,
                        "block": types.BoolType,
                    },
                },
                vulnerabilityObjectValues,
            )

            if diags.HasError() {
                return schemaRules, diags
            }

            conditionObject := types.ObjectValueMust(
                map[string]attr.Type{
                    //"device": types.StringType,
                    //"read_only": types.BoolType,
                    "vulnerabilities": types.ListType{
                        ElemType: types.ObjectType{
                            AttrTypes: map[string]attr.Type{
                                "id": types.Int32Type,
                                "block": types.BoolType,
                            },
                        },
                    },
                },
                map[string]attr.Value{
                    //"device": types.StringValue(rule.Condition.Device),
                    //"read_only": types.BoolValue(rule.Condition.ReadOnly),
                    "vulnerabilities": vulnerabilityObject,
                },
            )
            
            schemaRule.Condition = conditionObject
        }

        schemaRule.Notes = types.StringValue(rule.Notes)
        schemaRule.BlockMessage = types.StringValue(rule.BlockMessage) 
            
        schemaRules = append(schemaRules, schemaRule)
    }

    sortSchemaRules(&schemaRules, &planRules)

    return schemaRules, diags
}

//func collectionsToSchema(ctx context.Context, collections []collectionAPI.Collection) (types.Set, diag.Diagnostics) {
func collectionsToSchema(ctx context.Context, collections []collectionAPI.Collection) (types.List, diag.Diagnostics) {
    var diags diag.Diagnostics

    //collectionSet := types.SetNull(system.CollectionObjectType())
    collectionSet := types.ListNull(system.CollectionObjectType())
    collectionObjectValues := []attr.Value{}
    //for _, collection := range(rule.Collections) {
    for _, collection := range(collections) {
        accountIDs, diags := types.SetValueFrom(ctx, types.StringType, collection.AccountIDs)
        if diags.HasError() {
            return collectionSet, diags
        }

        appIDs, diags := types.SetValueFrom(ctx, types.StringType, collection.AppIDs)
        if diags.HasError() {
            return collectionSet, diags
        }

        clusters, diags := types.SetValueFrom(ctx, types.StringType, collection.Clusters)
        if diags.HasError() {
            return collectionSet, diags
        }

        containers, diags := types.SetValueFrom(ctx, types.StringType, collection.Containers)
        if diags.HasError() {
            return collectionSet, diags
        }

        functions, diags := types.SetValueFrom(ctx, types.StringType, collection.Functions)
        if diags.HasError() {
            return collectionSet, diags
        }

        hosts, diags := types.SetValueFrom(ctx, types.StringType, collection.Hosts)
        if diags.HasError() {
            return collectionSet, diags
        }

        images, diags := types.SetValueFrom(ctx, types.StringType, collection.Images)
        if diags.HasError() {
            return collectionSet, diags
        }

        labels, diags := types.SetValueFrom(ctx, types.StringType, collection.Labels)
        if diags.HasError() {
            return collectionSet, diags
        }
        
        namespaces, diags := types.SetValueFrom(ctx, types.StringType, collection.Namespaces)
        if diags.HasError() {
            return collectionSet, diags
        }

        collectionObjectValue := types.ObjectValueMust(
            system.CollectionObjectAttrTypeMap(),
            map[string]attr.Value{
                "account_ids": accountIDs,
                "app_ids": appIDs,
                "clusters": clusters,
                "color": types.StringValue(collection.Color),
                "containers": containers,
                "description": types.StringValue(collection.Description),
                "functions": functions,
                "hosts": hosts,
                "images": images,
                "labels": labels,
                "modified": types.StringValue(""),
                "name": types.StringValue(collection.Name),
                "namespaces": namespaces,
                "owner": types.StringValue(collection.Owner),
                "prisma": types.BoolValue(collection.Prisma),
                "system": types.BoolValue(collection.System),
            },
        )

        collectionObjectValues = append(collectionObjectValues, collectionObjectValue)
    }

    //collectionSet, diags = types.SetValueFrom(
    //    ctx,
    //    system.CollectionObjectType(),
    //    collectionObjectValues,
    //)
    collectionSet, diags = types.ListValueFrom(
        ctx,
        system.CollectionObjectType(),
        collectionObjectValues,
    )

    fmt.Println("***********************")
    fmt.Println("exiting collectionsToSchema")
    fmt.Println("***********************")

    return collectionSet, diags
}

func (r *HostCompliancePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("entering ModifyPlan")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    //fmt.Printf("%v\n", req.State.Raw)

    var plan HostCompliancePolicyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    fmt.Printf("%v\n", *plan.Rules)

    complianceVulnerabilities, err := systemAPI.GetComplianceHostVulnerabilities(*r.client)
	if err != nil {
		diags.AddError(
            "Error modifying planned policy rules", 
            "Failed to retrieve compliance host vulnerabilities from Prisma Cloud while modifying plan rules: " + err.Error(),
        )
        return
	}

    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("starting loop over rules")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    for index, rule := range *plan.Rules {
        if rule.Name.ValueString() == "r1" {
            sv := []system.CollectionResourceModel{}
            dgs := req.State.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("collections"), &sv)
            if dgs.HasError() {
                fmt.Println(dgs)
                return
            }
            fmt.Println("planValue:")
            fmt.Println(rule.Collections)
            fmt.Println("stateValue:") 
            fmt.Println(sv)
        }

        //if rule.Collections.IsUnknown() {
        if len(rule.Collections.Elements()) == 0 {
            fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
            fmt.Println("collections is unknown")
            fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
            //collections := types.SetValueMust(
            collections := types.ListValueMust(
                system.CollectionObjectType(),
                []attr.Value{
                    types.ObjectValueMust(
                        system.CollectionObjectAttrTypeMap(),
                        map[string]attr.Value{
                            "account_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "app_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "clusters": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "color": types.StringValue("#3FA2F7"),
                            "containers": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "description": types.StringValue("System - all resources collection"),
                            "functions": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "hosts": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "images": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "labels": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "modified": types.StringValue(""),
                            "name": types.StringValue("All"),
                            "namespaces": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
                            "owner": types.StringValue("system"),
                            "prisma": types.BoolValue(false),
                            "system": types.BoolValue(true),
                        },
                    ),
                },
            )
            diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("collections"), collections)...)
        }

        //for _, c := range rule.Collections.Elements() {
        //    //fmt.Println("!!!!!!!!!!!!!!!!!!!!")
        //    //fmt.Println(c)
        //    //fmt.Println("!!!!!!!!!!!!!!!!!!!!")
        //    if c.IsNull() {
        //        fmt.Println("!!!!!!!!!!!!!!!!!!!!")
        //        fmt.Println("we got a null here")
        //        fmt.Println("!!!!!!!!!!!!!!!!!!!!")
        //    }
        //}

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

        conditionObject, diags := createConditionFromEffect(ctx, *r.client, rule, complianceVulnerabilities)
        if diags.HasError() {
            return 
        }

        resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)
    }

    //var respPlan HostCompliancePolicyResourceModel
    //diags = resp.Plan.Get(ctx, &respPlan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //fmt.Printf("%+v\n", respPlan)
    //fmt.Printf("%+v\n", *respPlan.Rules)
    ////fmt.Printf("%+v\n", &respPlan.Rules.Elements()[0].Condition)
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
    fmt.Println("exiting ModifyPlan")
    fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
}
