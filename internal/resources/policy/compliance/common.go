package policy

import (
    "context"
    "sort"
	"fmt"
    "slices"
    "strings"
    "cmp"
    "time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GenerateCompliancePolicyRulesOrderMap(rules []models.CompliancePolicyRuleResourceModel) map[string]int {
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

func SortCompliancePolicyRules(rules *[]policyAPI.CompliancePolicyRule, planRules *[]models.CompliancePolicyRuleResourceModel) {
    rulesOrderMap := GenerateCompliancePolicyRulesOrderMap(*planRules) 
    r := *rules
    sort.Slice(r, func(i, j int) bool {
        return rulesOrderMap[r[i].Name] < rulesOrderMap[r[j].Name]
    })
    rules = &r 
}

func SortComplianceSchemaRules(ctx context.Context, schemaRules *[]models.CompliancePolicyRuleResourceModel, planRules *[]models.CompliancePolicyRuleResourceModel) {
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

    slices.SortFunc(*schemaRules, func(a, b models.CompliancePolicyRuleResourceModel) int {
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


    util.DLogf(ctx, *schemaRules)

    for i := 0; i < len(*schemaRules); i++ {
        (*schemaRules)[i].Order = (*planRules)[i].Order
    }
}

func GenerateConditionFromEffect(
    ctx context.Context, 
    client api.PrismaCloudComputeAPIClient, 
    policyType string,
    rule models.CompliancePolicyRuleResourceModel, 
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
            if effect != "ignore" {
                if effect == "default" {
                    complianceVulnerabilities = systemAPI.GetHighOrCriticalVulnerabilities(complianceVulnerabilities)
                }
                // If this block executes and effect is default, effect must be set to "alert"
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
            if effect == "default" {
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

func CompliancePolicySchemaToPolicy(ctx context.Context, plan *models.CompliancePolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (policyAPI.CompliancePolicy, diag.Diagnostics) {
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

func CompliancePolicyRuleSchemaToPolicy(ctx context.Context, planRules []models.CompliancePolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient) ([]policyAPI.CompliancePolicyRule, diag.Diagnostics) {
    util.DLog(ctx, "entering ComplianceRuleSchemaToPolicy")

    var diags diag.Diagnostics

    rules := []policyAPI.CompliancePolicyRule{}

    for _, planRule := range planRules {
        collectionNames := []string{}
        diags = planRule.Collections.ElementsAs(ctx, &collectionNames, false)
        if diags.HasError() {
            return rules, diags
        }

        collections, err := collectionAPI.GetCollections(*client, collectionNames)
        if err != nil {
            diags.AddError(
                "Value Conversion Error",
                fmt.Sprintf("Error retrieving collection names while converting compliance policy rules to schema: %s", err.Error()),
            )
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

    SortCompliancePolicyRules(&rules, &planRules)

    util.DLog(ctx, fmt.Sprintf("exiting ComplianceRuleSchemaToPolicy with return value rules:\n\n %+v", rules))
    
    return rules, diags
}

func CompliancePolicyToSchema(ctx context.Context, policy policyAPI.CompliancePolicy, plan models.CompliancePolicyResourceModel) (models.CompliancePolicyResourceModel, diag.Diagnostics) {
    var diags diag.Diagnostics

    schema := models.CompliancePolicyResourceModel{
        Id: types.StringValue(policy.Id),
        PolicyType: types.StringValue(policy.PolicyType),
    }

    var rules []models.CompliancePolicyRuleResourceModel

    if policy.Rules != nil {
        rules, diags = CompliancePolicyRulesToSchema(ctx, *policy.Rules, plan.Rules)
        if diags.HasError() {
            util.DLog(ctx, "CompliancePolicyRulesToSchema error")
            return schema, diags
        }
    } else {
        rules = []models.CompliancePolicyRuleResourceModel{}
    }

    schema.Rules = &rules

    return schema, diags
}

func CompliancePolicyRulesToSchema(ctx context.Context, rules []policyAPI.CompliancePolicyRule, planRules *[]models.CompliancePolicyRuleResourceModel) ([]models.CompliancePolicyRuleResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "entering CompliancePolicyRulesToSchema")

    var diags diag.Diagnostics

    schemaRules := []models.CompliancePolicyRuleResourceModel{}

    if len(rules) == 0 {
        return schemaRules, diags
    }

    for _, rule := range rules {
        schemaRule := models.CompliancePolicyRuleResourceModel{
            Disabled: types.BoolValue(rule.Disabled),
            Effect: types.StringValue(rule.Effect),
            Modified: types.StringValue(""),
            Name: types.StringValue(rule.Name),
            Owner: types.StringValue(rule.Owner),
            ReportAllPassedAndFailedChecks: types.BoolValue(rule.ReportAllPassedAndFailedChecks),
            Verbose: types.BoolValue(rule.Verbose),
        }

        collectionNames := []string{}
        for _, collection := range rule.Collections {
            collectionNames = append(collectionNames, collection.Name) 
        }

        collections, diags := types.SetValueFrom(ctx, types.StringType, collectionNames)
        if diags.HasError() {
            return schemaRules, diags
        }

        schemaRule.Collections = collections

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

func isWildcard(val []string) bool {
    return len(val) == 1 && val[0] == "*"
}

func validateRuleCollectionsByPolicyType(ctx context.Context, policyType string, ruleName string, ruleCollectionNames []string, collections []collectionAPI.Collection) diag.Diagnostics {
    var diags diag.Diagnostics

    for _, ruleCollectionName := range ruleCollectionNames {
        found := false
        invalidFields := make([]string, 0)

        for _, collection := range collections {
            if ruleCollectionName == collection.Name {
                found = true

                switch policyType {
                    case "containerCompliance":
                        if !isWildcard(collection.Functions) {
                            invalidFields = append(invalidFields, "Functions")
                        }
                    case "ciImagesCompliance":
                        if !isWildcard(collection.Containers) {
                            invalidFields = append(invalidFields, "Containers")
                        }
                        if !isWildcard(collection.Hosts) {
                            invalidFields = append(invalidFields, "Hosts")
                        }
                        if !isWildcard(collection.AppIDs) {
                            invalidFields = append(invalidFields, "App IDs")
                        }
                        if !isWildcard(collection.Functions) {
                            invalidFields = append(invalidFields, "Functions")
                        }
                        if !isWildcard(collection.Namespaces) {
                            invalidFields = append(invalidFields, "Namespaces")
                        }
                        if !isWildcard(collection.AccountIDs) {
                            invalidFields = append(invalidFields, "Account IDs")
                        }
                        if !isWildcard(collection.Clusters) {
                            invalidFields = append(invalidFields, "Clusters")
                        }
                    case "hostCompliance":
                        if !isWildcard(collection.Containers) {
                            invalidFields = append(invalidFields, "Containers")
                        }
                        if !isWildcard(collection.Images) {
                            invalidFields = append(invalidFields, "Images")
                        }
                        if !isWildcard(collection.AppIDs) {
                            invalidFields = append(invalidFields, "App IDs")
                        }
                        if !isWildcard(collection.Functions) {
                            invalidFields = append(invalidFields, "Functions")
                        }
                        if !isWildcard(collection.Namespaces) {
                            invalidFields = append(invalidFields, "Namespaces")
                        }
                    case "vmCompliance":
                        if !isWildcard(collection.Containers) {
                            invalidFields = append(invalidFields, "Containers")
                        }
                        if !isWildcard(collection.Hosts) {
                            invalidFields = append(invalidFields, "Hosts")
                        }
                        if !isWildcard(collection.AppIDs) {
                            invalidFields = append(invalidFields, "App IDs")
                        }
                        if !isWildcard(collection.Functions) {
                            invalidFields = append(invalidFields, "Functions")
                        }
                        if !isWildcard(collection.Namespaces) {
                            invalidFields = append(invalidFields, "Namespaces")
                        }
                        if !isWildcard(collection.Clusters) {
                            invalidFields = append(invalidFields, "Clusters")
                        }
                    case "serverlessCompliance":
                        if !isWildcard(collection.Containers) {
                            invalidFields = append(invalidFields, "Containers")
                        }
                        if !isWildcard(collection.Hosts) {
                            invalidFields = append(invalidFields, "Hosts")
                        }
                        if !isWildcard(collection.Images) {
                            invalidFields = append(invalidFields, "Images")
                        }
                        if !isWildcard(collection.AppIDs) {
                            invalidFields = append(invalidFields, "App IDs")
                        }
                        if !isWildcard(collection.Namespaces) {
                            invalidFields = append(invalidFields, "Namespaces")
                        }
                        if !isWildcard(collection.Clusters) {
                            invalidFields = append(invalidFields, "Clusters")
                        }
                    case "ciServerlessCompliance":
                        if !isWildcard(collection.Containers) {
                            invalidFields = append(invalidFields, "Containers")
                        }
                        if !isWildcard(collection.Hosts) {
                            invalidFields = append(invalidFields, "Hosts")
                        }
                        if !isWildcard(collection.Images) {
                            invalidFields = append(invalidFields, "Images")
                        }
                        if !isWildcard(collection.AppIDs) {
                            invalidFields = append(invalidFields, "App IDs")
                        }
                        if !isWildcard(collection.Namespaces) {
                            invalidFields = append(invalidFields, "Namespaces")
                        }
                        if !isWildcard(collection.AccountIDs) {
                            invalidFields = append(invalidFields, "Account IDs")
                        }
                        if !isWildcard(collection.Clusters) {
                            invalidFields = append(invalidFields, "Clusters")
                        }
                    // TODO: implement logic for trusted images
                    //case "trust":
                    default:
                        diags.AddError(
                            "Resource Validation Error",
                            fmt.Sprintf("Error occured during validation of collections for policy rule \"%s\": Invalid policy type \"%s\"", ruleName, policyType),
                        )
                }

                if !found {
                    diags.AddError(
                        "Resource Validation Error",
                        fmt.Sprintf("Error occured during validation of collections for policy rule \"%s\": Collection name \"%s\" not found", ruleName, policyType),
                    )
                } else if len(invalidFields) > 0 {
                    invalidFieldsString := strings.Join(invalidFields, ", ")
                    diags.AddError(
                        "Resource Validation Error",
                        fmt.Sprintf("%s policy rule \"%s\" is configured with collection \"%s\", which contains invalid values. The following fields in the collection must only contain the wildcard value (\"*\") for the collection to be compatible with this policy: %s", policyType, ruleName, ruleCollectionName, invalidFieldsString),
                    )
                }
            }
        }
    }

    return diags
}

func ModifyCompliancePolicyResourcePlan(ctx context.Context, client *api.PrismaCloudComputeAPIClient, plan *models.CompliancePolicyResourceModel, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyCompliancePolicyResourcePlan")
    
    var diags diag.Diagnostics

    if plan == nil {
        return
    }

    if plan.Rules == nil {
        emptyRules := []models.CompliancePolicyRuleResourceModel{}
        diags.Append(resp.Plan.SetAttribute(ctx, path.Root("rules"), &emptyRules)...)
        return
    }

    policyType := plan.PolicyType.ValueString()

    util.DLog(ctx, "Retrieving vulnerability data")
    complianceVulnerabilities, err := systemAPI.GetComplianceVulnerabilities(*client, policyType)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error modifying planned policy rules", 
            "Failed to retrieve compliance host vulnerabilities from Prisma Cloud while modifying plan rules: " + err.Error(),
        )
        return
	}

    util.DLog(ctx, "Retrieving collections")
    collections, err := collectionAPI.ListCollections(*client)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error modifying planned policy rules", 
            "Failed to retrieve collections from Prisma Cloud while modifying plan rules: " + err.Error(),
        )
        return
	}
    
    util.DLog(ctx, "Beginning loop over rules")

    for index, rule := range *plan.Rules {
        ruleName := rule.Name.ValueString()
        rulePath := path.Root("rules").AtListIndex(index)

        ruleCollectionNames := []string{}
        resp.Diagnostics.Append(rule.Collections.ElementsAs(ctx, &ruleCollectionNames, false)...)
        if resp.Diagnostics.HasError() {
            return
        }
       
        // Validate configured collections to ensure compatibility with policy type
        diags = validateRuleCollectionsByPolicyType(ctx, policyType, ruleName, ruleCollectionNames, collections) 
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return 
        }
        
        // Generate rule condition value from effect
        conditionObject, diags := GenerateConditionFromEffect(ctx, *client, plan.PolicyType.ValueString(), rule, complianceVulnerabilities)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return 
        }

        // Set condition to generated value
        resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rulePath.AtName("condition"), conditionObject)...)
        if resp.Diagnostics.HasError() {
            return 
        }

        // Set effect to "alert" if its currently set to the placeholder value ("default")
        if rule.Effect.ValueString() == "default" {
            resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rulePath.AtName("effect"), types.StringValue("alert"))...)
            if resp.Diagnostics.HasError() {
                return 
            }
        }
    }

    util.DLog(ctx, "exiting ModifyCompliancePolicyResourcePlan")
}
