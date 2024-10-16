package policy

import (
    "context"
    "sort"
	"fmt"
    "slices"
    "cmp"
    "time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type CompliancePolicyResourceModel struct {
    Id          types.String                                `tfsdk:"id"`
    PolicyType  types.String                                `tfsdk:"policy_type"`
    Rules       *[]CompliancePolicyRuleResourceModel    `tfsdk:"rules"`
}

type CompliancePolicyRuleResourceModel struct {
    BlockMessage                    types.String    `tfsdk:"block_message"`
    Collections                     types.List      `tfsdk:"collections"`
    Condition                       types.Object    `tfsdk:"condition"`
    Disabled                        types.Bool      `tfsdk:"disabled"`
    Effect                          types.String    `tfsdk:"effect"`
    Modified                        types.String    `tfsdk:"modified"`
    Name                            types.String    `tfsdk:"name"`
    Notes                           types.String    `tfsdk:"notes"`
    Order                           types.Int32     `tfsdk:"order"`
    Owner                           types.String    `tfsdk:"owner"`
    ReportAllPassedAndFailedChecks  types.Bool      `tfsdk:"report_passed_and_failed_checks"`
    Verbose                         types.Bool      `tfsdk:"verbose"`
}

func GenerateCompliancePolicyRulesOrderMap(rules []CompliancePolicyRuleResourceModel) map[string]int {
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

func SortCompliancePolicyRules(rules *[]policyAPI.CompliancePolicyRule, planRules *[]CompliancePolicyRuleResourceModel) {
    rulesOrderMap := GenerateCompliancePolicyRulesOrderMap(*planRules) 
    r := *rules
    sort.Slice(r, func(i, j int) bool {
        return rulesOrderMap[r[i].Name] < rulesOrderMap[r[j].Name]
    })
    rules = &r 
}

func SortComplianceSchemaRules(ctx context.Context, schemaRules *[]CompliancePolicyRuleResourceModel, planRules *[]CompliancePolicyRuleResourceModel) {
    util.DLog(ctx, "entering SortComplianceSchemaRules")

    if planRules == nil {
        for i := 0; i < len(*schemaRules); i++ {
            (*schemaRules)[i].Order = types.Int32Value(int32(i + 1))
        }
        return
    }

    ruleOrderMap := make(map[string]int32)
    for index, planRule := range *planRules {
        ruleOrderMap[planRule.Name.ValueString()] = int32(index)
    }

    slices.SortFunc(*schemaRules, func(a, b CompliancePolicyRuleResourceModel) int {
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


    //util.DLog(ctx, fmt.Sprintf("entering SortComplianceSchemaRules")
    util.DLogf(ctx, *schemaRules)

    for i := 0; i < len(*schemaRules); i++ {
        (*schemaRules)[i].Order = (*planRules)[i].Order
    }
}

func GenerateConditionFromEffect(
    ctx context.Context, 
    client api.PrismaCloudComputeAPIClient, 
    policyType string,
    rule CompliancePolicyRuleResourceModel, 
    complianceVulnerabilities []systemAPI.Vulnerability) (basetypes.ObjectValue, diag.Diagnostics) {
    util.DLog(ctx, "entering GenerateConditionFromEffect")
    // TODO: fix modification from "effect = alert" to no effect not creating the right values (doesnt think any
    // changes are needed since effect gets set to "alert" when initially creating a rule with no effect value)

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

    util.DLogf(ctx, rule)
    
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

    if policyType == "serverlessCompliance" || policyType == "ciServerlessCompliance" {
        if rule.Condition.IsUnknown() {
            util.DLog(ctx, "condition is unknown")
            if effect != "ignore" {
                if effect == "unknown" {
                    complianceVulnerabilities = systemAPI.GetHighOrCriticalVulnerabilities(complianceVulnerabilities)
                }
                // If this block executes and effect is unknown, effect must be set to "alert"
                for _, vuln := range complianceVulnerabilities {
                    vulnerabilityObjectValue := types.ObjectValueMust(
                        vulnerabilitiesAttributeTypes,
                        map[string]attr.Value{
                            "id": types.Int32Value(int32(vuln.Id)),
                            "block": types.BoolValue(false),
                        },
                    )
                    vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
                }
            }
        } else {
            util.DLog(ctx, "condition is not unknown")
            // Get rule condition configuration
            ruleCondition := policyAPI.CompliancePolicyRuleCondition{} 
            diags = rule.Condition.As(ctx, &ruleCondition, basetypes.ObjectAsOptions{})
            if diags.HasError() {
                return conditionObject, diags
            }
            ruleConditionVulns := ruleCondition.Vulnerabilities

            // Loop over condition configuration and populate specified values to vulnerabilityObjectValues
            for _, vuln := range ruleConditionVulns {
                vulnerabilityObjectValue := types.ObjectValueMust(
                    vulnerabilitiesAttributeTypes,
                    map[string]attr.Value{
                        "id": types.Int32Value(int32(vuln.Id)),
                        //"block": types.BoolValue(vuln.Block),
                        // TODO: change this value in resp to "false" if not already set to false
                        "block": types.BoolValue(false),
                    },
                )
                vulnerabilityObjectValues = append(vulnerabilityObjectValues, vulnerabilityObjectValue)
            }
        }
    } else {
        // If the effect is "alert, block", create condition vulnerabilities object from plan
        if effect == "alert, block" {
            var ruleConditionVulns []policyAPI.CompliancePolicyRuleVulnerability

            if rule.Condition.IsUnknown() {
                diags.AddError(
                    "Missing condition from \"alert, block\" effect rule",
                    "Condition attribute must be defined for rules with effect \"alert, block\".",
                )
                return conditionObject, diags
            }

            ruleCondition := policyAPI.CompliancePolicyRuleCondition{} 
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

            var block bool
            for _, vuln := range complianceVulnerabilities {
                if policyType == "hostCompliance" {
                    block = (isBlockEffect && !(vuln.Type == "windows"))
                } else if policyType == "containerCompliance" {
                    block = (isBlockEffect && !(vuln.Type == "istio" || vuln.Id == 58 || vuln.Id == 596 || vuln.Id == 598))
                } else if policyType == "ciImagesCompliance" {
                    block = isBlockEffect
                } else if policyType == "vmCompliance" {
                    block = (isBlockEffect && !(vuln.Type == "istio" || vuln.Id == 58 || vuln.Id == 596 || vuln.Id == 598))
                } else {
                    // TODO: append error here
                    return conditionObject, diags
                }

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
    
    util.DLog(ctx, fmt.Sprintf("created conditionObject:\n\n%v\n", conditionObject))

    return conditionObject, diags
}

func CompliancePolicySchemaToPolicy(ctx context.Context, plan *CompliancePolicyResourceModel, client *api.PrismaCloudComputeAPIClient,/*, username types.String*/) (policyAPI.CompliancePolicy, diag.Diagnostics) {
    var diags diag.Diagnostics

    policy := policyAPI.CompliancePolicy{
        Id: plan.Id.ValueString(),
        PolicyType: plan.PolicyType.ValueString(),
    }

    if plan.Rules == nil {
        rules := []policyAPI.CompliancePolicyRule{}
        policy.Rules = &rules
        return policy, diags
    }

    rules, diags := CompliancePolicyRuleSchemaToPolicy(ctx, *plan.Rules, client)
    if diags.HasError() {
        return policy, diags
    }

    policy.Rules = &rules

    return policy, diags
}

func CompliancePolicyRuleSchemaToPolicy(ctx context.Context, planRules []CompliancePolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient, /*, username types.String*/) ([]policyAPI.CompliancePolicyRule, diag.Diagnostics) {
    util.DLog(ctx, "entering ComplianceRuleSchemaToPolicy")

    var diags diag.Diagnostics

    rules := []policyAPI.CompliancePolicyRule{}

    for _, planRule := range planRules {
        collections := []collectionAPI.Collection{}
        diags = planRule.Collections.ElementsAs(ctx, &collections, false)
        if diags.HasError() {
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

        condition := policyAPI.CompliancePolicyRuleCondition{} 
        diags = planRule.Condition.As(ctx, &condition, basetypes.ObjectAsOptions{})
        if diags.HasError() {
            return rules, diags
        }

        rule := policyAPI.CompliancePolicyRule{
            BlockMessage: planRule.BlockMessage.ValueString(),
            Collections: collections,
            Condition: &condition,
            Disabled: planRule.Disabled.ValueBool(),
            Effect: planRule.Effect.ValueString(),
            Modified: time.Now().Format("2006-01-02T15:04:05.000Z"),
            Name: planRule.Name.ValueString(), 
            Order: int(planRule.Order.ValueInt32()),
            ReportAllPassedAndFailedChecks: planRule.ReportAllPassedAndFailedChecks.ValueBool(),
            Verbose: planRule.Verbose.ValueBool(),
        }
        
        if !planRule.Notes.IsUnknown() && !planRule.Notes.IsNull() {
            rule.Notes = planRule.Notes.ValueString()
        }

        rules = append(rules, rule)
    }

    //sortContainerCompliancePolicyRules(&rules, &planRules)
    SortCompliancePolicyRules(&rules, &planRules)

    util.DLog(ctx, fmt.Sprintf("exiting ComplianceRuleSchemaToPolicy with return value rules:\n\n %+v", rules))
    
    return rules, diags
}

func CollectionsToSchema(ctx context.Context, collections []collectionAPI.Collection) (types.List, diag.Diagnostics) {
    var diags diag.Diagnostics

    collectionList := types.ListNull(system.CollectionObjectType())
    collectionObjectValues := []attr.Value{}
    for _, collection := range(collections) {
        accountIDs, diags := types.SetValueFrom(ctx, types.StringType, collection.AccountIDs)
        if diags.HasError() {
            return collectionList, diags
        }

        appIDs, diags := types.SetValueFrom(ctx, types.StringType, collection.AppIDs)
        if diags.HasError() {
            return collectionList, diags
        }

        clusters, diags := types.SetValueFrom(ctx, types.StringType, collection.Clusters)
        if diags.HasError() {
            return collectionList, diags
        }

        containers, diags := types.SetValueFrom(ctx, types.StringType, collection.Containers)
        if diags.HasError() {
            return collectionList, diags
        }

        functions, diags := types.SetValueFrom(ctx, types.StringType, collection.Functions)
        if diags.HasError() {
            return collectionList, diags
        }

        hosts, diags := types.SetValueFrom(ctx, types.StringType, collection.Hosts)
        if diags.HasError() {
            return collectionList, diags
        }

        images, diags := types.SetValueFrom(ctx, types.StringType, collection.Images)
        if diags.HasError() {
            return collectionList, diags
        }

        labels, diags := types.SetValueFrom(ctx, types.StringType, collection.Labels)
        if diags.HasError() {
            return collectionList, diags
        }
        
        namespaces, diags := types.SetValueFrom(ctx, types.StringType, collection.Namespaces)
        if diags.HasError() {
            return collectionList, diags
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

    collectionList, diags = types.ListValueFrom(ctx, system.CollectionObjectType(), collectionObjectValues)

    util.DLog(ctx, "exiting collectionsToContainerSchema")

    return collectionList, diags
}


func CompliancePolicyToSchema(ctx context.Context, policy policyAPI.CompliancePolicy, plan CompliancePolicyResourceModel) (CompliancePolicyResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    schema := CompliancePolicyResourceModel{
        Id: types.StringValue(policy.Id),
        PolicyType: types.StringValue(policy.PolicyType),
    }

    var rules []CompliancePolicyRuleResourceModel

    if policy.Rules != nil {
        rules, diags = CompliancePolicyRulesToSchema(ctx, *policy.Rules, plan.Rules)
        if diags.HasError() {
            util.DLog(ctx, "CompliancePolicyRulesToSchema error")
            return schema, diags
        }
    } else {
        rules = []CompliancePolicyRuleResourceModel{}
    }

    schema.Rules = &rules

    return schema, diags
}

func CompliancePolicyRulesToSchema(ctx context.Context, rules []policyAPI.CompliancePolicyRule, planRules *[]CompliancePolicyRuleResourceModel) ([]CompliancePolicyRuleResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "entering CompliancePolicyRulesToSchema")

    var diags diag.Diagnostics

    schemaRules := []CompliancePolicyRuleResourceModel{}

    if len(rules) == 0 {
        return schemaRules, diags
    }

    for _, rule := range rules {
        schemaRule := CompliancePolicyRuleResourceModel{
            Disabled: types.BoolValue(rule.Disabled),
            Effect: types.StringValue(rule.Effect),
            Modified: types.StringValue(""),
            Name: types.StringValue(rule.Name),
            Owner: types.StringValue(rule.Owner),
            ReportAllPassedAndFailedChecks: types.BoolValue(rule.ReportAllPassedAndFailedChecks),
            Verbose: types.BoolValue(rule.Verbose),
        }

        if rule.Collections != nil {
            collections, diags := CollectionsToSchema(ctx, rule.Collections)
            if diags.HasError() {
                return schemaRules, diags
            }

            schemaRule.Collections = collections
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

    SortComplianceSchemaRules(ctx, &schemaRules, planRules)

    return schemaRules, diags
}


func ModifyCompliancePolicyResourcePlan(ctx context.Context, client *api.PrismaCloudComputeAPIClient, plan *CompliancePolicyResourceModel, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyCompliancePolicyResourcePlan")
    
    var diags diag.Diagnostics

    if plan == nil {
        return
    }

    if plan.Rules == nil {
        emptyRules := []CompliancePolicyRuleResourceModel{}
        diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules"), &emptyRules)...)
        return
    }

    //fmt.Printf("%v\n", *plan.Rules)

    util.DLog(ctx, "getting vulns")
    complianceVulnerabilities, err := systemAPI.GetComplianceVulnerabilities(*client, plan.PolicyType.ValueString())
	if err != nil {
		//diags.AddError(
		resp.Diagnostics.AddError(
            "Error modifying planned policy rules", 
            "Failed to retrieve compliance host vulnerabilities from Prisma Cloud while modifying plan rules: " + err.Error(),
        )
        return
	}

    util.DLog(ctx, "starting loop over rules")

    for index, rule := range *plan.Rules {
        // Set default collection if one is not specified in rule configuration
        if len(rule.Collections.Elements()) == 0 {
            diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("collections"), system.GetDefaultCollectionObject())...)
        }

        // Set unknown/null rule order to the value of the rule's index in the configuration
        if rule.Order.IsUnknown() || rule.Order.IsNull() {
            rule.Order = types.Int32Value(int32(index + 1))
            diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("order"), types.Int32Value(int32(index + 1)))...)
            if diags.HasError() {
                return
            }
        // Raise error if specified rule order value is less than 1
        // TODO: can this logic live in a validator for the attribute?
        } else if int(rule.Order.ValueInt32()) < 1 {
            resp.Diagnostics.AddError(
		    	"Invalid Resource Configuration",
		    	fmt.Sprintf("Host Compliance Policy Rule specified an invalid order (%d). Order values must be positive non-zero integers.", int(rule.Order.ValueInt32())),
		    )
            return
        }
        
        // Set unknown rule effect to "alert"
        if rule.Effect.IsUnknown() {
            rule.Effect = types.StringValue("unknown")
            diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("effect"), types.StringValue("alert"))...)
            if diags.HasError() {
                return
            }
        } 

        // Generate rule condition value from effect and set in response
        conditionObject, diags := GenerateConditionFromEffect(ctx, *client, plan.PolicyType.ValueString(), rule, complianceVulnerabilities)
        if diags.HasError() {
            return 
        }

        resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), conditionObject)...)
    }

    util.DLog(ctx, "exiting ModifyCompliancePolicyResourcePlan")
}
