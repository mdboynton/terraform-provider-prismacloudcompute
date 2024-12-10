package policy

import (
	"context"
	"fmt"
	"sort"

	//"slices"
	"strings"
	//"cmp"
	"time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	//models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	systemAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/system"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/planmodifiers"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GetPolicySchema(ctx context.Context, policyType string, policyTypeFormatted string, policyContext string, metaType string) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "TODO",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Optional:            true,
				Computed:            true,
				//Default: stringdefault.StaticString(policyAPI.PolicyTypeComplianceHost),
				Default: stringdefault.StaticString(policyType),
			},
			"policy_context": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				//Default: stringdefault.StaticString(policyAPI.PolicyContextVulnerabilityDeployedImage),
				Default: stringdefault.StaticString(policyContext),
			},
			"policy_type": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Optional:            true,
				Computed:            true,
				//Default: stringdefault.StaticString(policyAPI.PolicyTypeComplianceHost),
				Default: stringdefault.StaticString(policyType),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				//Default: stringdefault.StaticString(policyAPI.TypeVulnerability),
				Default: stringdefault.StaticString(metaType),
			},
			"rules": schema.ListNestedAttribute{
				MarkdownDescription: "TODO",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					//planmodifiers.UseIndexForUnknownOrder(policyAPI.PolicyTypeComplianceHostFormatted),
					//planmodifiers.SetNullByModuleType(metaType),
					planmodifiers.UseIndexForUnknownOrder(policyTypeFormatted),
				},
				Validators: []validator.List{
					//validators.PolicyRuleNameIsUnique(policyAPI.PolicyTypeComplianceHostFormatted),
					//validators.PolicyRuleOrderIsPositiveNonZero(policyAPI.PolicyTypeComplianceHostFormatted),
					validators.PolicyRuleNameIsUnique(policyTypeFormatted),
					validators.PolicyRuleOrderIsPositiveNonZero(policyTypeFormatted),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"order": schema.Int32Attribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
						},
						"apply_only_when_fix_available": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//Default: booldefault.StaticBool(false),
						},
						"alert_threshold": schema.SingleNestedAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							//Default: objectdefault.StaticValue(
							//    types.ObjectValueMust(
							//        map[string]attr.Type{
							//            "threshold": types.StringType,
							//            "risk_factors": types.SetType{ ElemType: types.StringType },
							//        },
							//        map[string]attr.Value{
							//            "threshold": types.StringValue("off"),
							//            "risk_factors": types.SetValueMust(
							//                types.StringType,
							//                []attr.Value{},
							//            ),
							//        },
							//    ),
							//),
							PlanModifiers: []planmodifier.Object{
								planmodifiers.SetThresholdByModuleType(metaType),
							},
							Attributes: map[string]schema.Attribute{
								"threshold": schema.StringAttribute{
									MarkdownDescription: "TODO",
									Optional:            true,
									Computed:            true,
								},
								"risk_factors": schema.SetAttribute{
									ElementType:         types.StringType,
									MarkdownDescription: "TODO",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.Set{
										planmodifiers.UseEmptySetForUnknownRiskFactors(),
									},
								},
							},
						},
						"block_threshold": schema.SingleNestedAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							//Default: objectdefault.StaticValue(
							//    types.ObjectValueMust(
							//        map[string]attr.Type{
							//            "threshold": types.StringType,
							//            "risk_factors": types.SetType{ ElemType: types.StringType },
							//        },
							//        map[string]attr.Value{
							//            "threshold": types.StringValue("off"),
							//            "risk_factors": types.SetValueMust(
							//                types.StringType,
							//                []attr.Value{},
							//            ),
							//        },
							//    ),
							//),
							PlanModifiers: []planmodifier.Object{
								planmodifiers.SetThresholdByModuleType(metaType),
							},
							Attributes: map[string]schema.Attribute{
								"threshold": schema.StringAttribute{
									MarkdownDescription: "TODO",
									Optional:            true,
									Computed:            true,
								},
								"risk_factors": schema.SetAttribute{
									ElementType:         types.StringType,
									MarkdownDescription: "TODO",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.Set{
										planmodifiers.UseEmptySetForUnknownRiskFactors(),
									},
								},
							},
						},
						"report_passed_and_failed_checks": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//Default: booldefault.StaticBool(false),
						},
						"block_message": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
						},
						"collections": schema.SetAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							Default: setdefault.StaticValue(
								types.SetValueMust(types.StringType, []attr.Value{types.StringValue("All")}),
							),
							PlanModifiers: []planmodifier.Set{
								planmodifiers.UseEmptySetForUnknown(),
							},
						},
						"condition": schema.SingleNestedAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"vulnerabilities": schema.ListNestedAttribute{
									MarkdownDescription: "TODO",
									Optional:            true,
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.Int32Attribute{
												MarkdownDescription: "TODO",
												Optional:            true,
												Computed:            true,
											},
											"block": schema.BoolAttribute{
												MarkdownDescription: "TODO",
												Optional:            true,
												Computed:            true,
											},
										},
									},
									//PlanModifiers: []planmodifier.List{
									//    planmodifiers.UseEmptyListForUnknownConditionVulnerabilities(),
									//},
								},
							},
							Default: objectdefault.StaticValue(
								types.ObjectValueMust(
									map[string]attr.Type{
										"vulnerabilities": types.ListType{
											ElemType: types.ObjectType{
												AttrTypes: map[string]attr.Type{
													"id":    types.Int32Type,
													"block": types.BoolType,
												},
											},
										},
									},
									map[string]attr.Value{
										"vulnerabilities": types.ListValueMust(
											types.ObjectType{
												AttrTypes: map[string]attr.Type{
													"id":    types.Int32Type,
													"block": types.BoolType,
												},
											},
											[]attr.Value{},
										),
									},
								),
							),
							PlanModifiers: []planmodifier.Object{
								planmodifiers.UseDefaultForUnknownCondition(),
							},
						},
						"cve_rules": schema.SetNestedAttribute{
							// TODO: validation
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//Default: setdefault.StaticValue(
							//    types.SetValueMust(
							//        types.ObjectType{
							//            AttrTypes: map[string]attr.Type{
							//                "name": types.StringType,
							//                "effect": types.StringType,
							//                "id": types.StringType,
							//                "description": types.StringType,
							//                "type": types.StringType,
							//                "expiration": types.ObjectType{
							//                    AttrTypes: map[string]attr.Type{
							//                        "enabled": types.BoolType,
							//                        "date": types.StringType,
							//                    },
							//                },
							//            },
							//        },
							//        []attr.Value{},
							//    ),
							//),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									//"name": schema.StringAttribute{
									//    MarkdownDescription: "TODO",
									//    Optional: true,
									//    //Computed: true,
									//    //Default: stringdefault.StaticString(""),
									//},
									"effect": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
										// TODO: validator
									},
									"id": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
									},
									//"type": schema.StringAttribute{
									//    MarkdownDescription: "TODO",
									//    Optional: true,
									//    //Computed: true,
									//    //Default: stringdefault.StaticString("cve"),
									//},
									"expiration": schema.SingleNestedAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												MarkdownDescription: "TODO",
												Optional:            true,
												//Computed: true,
											},
											"date": schema.StringAttribute{
												// TODO: validator (check format, ensure alert effect != block)
												MarkdownDescription: "TODO",
												Optional:            true,
												//Computed: true,
											},
										},
									},
								},
							},
						},
						"disabled": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"effect": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								planmodifiers.UseDefaultForUnknownEffect(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("ignore", "alert", "block", "alert, block"),
							},
						},
						"exclude_base_image_vulns": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//Default: booldefault.StaticBool(false),
						},
						"grace_days_all_severities": schema.Int32Attribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
						},
						"grace_days_by_severity": schema.SingleNestedAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//Default: objectdefault.StaticValue(
							//    types.ObjectValueMust(
							//        map[string]attr.Type{
							//            "low": types.Int32Type,
							//            "medium": types.Int32Type,
							//            "high": types.Int32Type,
							//            "critical": types.Int32Type,
							//        },
							//        map[string]attr.Value{
							//            "low": types.Int32Value(int32(0)),
							//            "medium": types.Int32Value(int32(0)),
							//            "high": types.Int32Value(int32(0)),
							//            "critical": types.Int32Value(int32(0)),
							//        },
							//   ),
							//),
							// TODO: add validator to make sure no values are specified for
							// severities that are below the block threshold
							Attributes: map[string]schema.Attribute{
								"low": schema.Int32Attribute{
									MarkdownDescription: "TODO",
									Optional:            true,
									//Computed: true,
								},
								"medium": schema.Int32Attribute{
									MarkdownDescription: "TODO",
									Optional:            true,
									//Computed: true,
								},
								"high": schema.Int32Attribute{
									MarkdownDescription: "TODO",
									Optional:            true,
									//Computed: true,
								},
								"critical": schema.Int32Attribute{
									MarkdownDescription: "TODO",
									Optional:            true,
									//Computed: true,
								},
							},
						},
						"modified": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							//Default: stringdefault.StaticString(time.Now().Format("2006-01-02T15:04:05.000Z")),
							PlanModifiers: []planmodifier.String{
								//UseStateForUnknown(),
								//UsePlanForUnknownString(),
								planmodifiers.UseEmptyStringForNull(),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
							//Optional: true,
							//Computed: true,
						},
						"notes": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//PlanModifiers: []planmodifier.String{
							//    stringplanmodifier.UseStateForUnknown(),
							//},
						},
						"owner": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"package_types_thresholds": schema.SetNestedAttribute{
							// TODO: validation (check type, check thresholds)
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//Default: setdefault.StaticValue(
							//    types.SetValueMust(
							//        types.ObjectType{
							//            AttrTypes: map[string]attr.Type{
							//                "type": types.StringType,
							//                "alert_threshold": types.StringType,
							//                "block_threshold": types.StringType,
							//            },
							//        },
							//        []attr.Value{},
							//    ),
							//),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
										// TODO: validator
									},
									"alert_threshold": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
									},
									"block_threshold": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
									},
								},
							},
						},
						"tags": schema.SetNestedAttribute{
							// TODO: validation
							MarkdownDescription: "TODO",
							Optional:            true,
							//Computed: true,
							//Default: setdefault.StaticValue(
							//    types.SetValueMust(
							//        types.ObjectType{
							//            AttrTypes: map[string]attr.Type{
							//                "effect": types.StringType,
							//                "id": types.StringType,
							//                "description": types.StringType,
							//                "type": types.StringType,
							//                "name": types.StringType,
							//                "expiration": types.ObjectType{
							//                    AttrTypes: map[string]attr.Type{
							//                        "enabled": types.BoolType,
							//                        "date": types.StringType,
							//                    },
							//                },
							//            },
							//        },
							//        []attr.Value{},
							//    ),
							//),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									//"name": schema.StringAttribute{
									//    MarkdownDescription: "TODO",
									//    Required: true,
									//    //Optional: true,
									//    //Computed: true,
									//},
									"effect": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
										// TODO: validator
									},
									"id": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										//Computed: true,
									},
									//"type": schema.StringAttribute{
									//    MarkdownDescription: "TODO",
									//    Optional: true,
									//    //Computed: true,
									//    //Default: stringdefault.StaticString("tag"),
									//},
									"expiration": schema.SingleNestedAttribute{
										MarkdownDescription: "TODO",
										Optional:            true,
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												MarkdownDescription: "TODO",
												Optional:            true,
												//Computed: true,
											},
											"date": schema.StringAttribute{
												// TODO: validator (check format, ensure alert effect != block)
												MarkdownDescription: "TODO",
												Optional:            true,
												//Computed: true,
											},
										},
									},
								},
							},
						},
						"verbose": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func PolicySchemaToTerraform(ctx context.Context, plan *models.PolicyResourceModel, client *api.PrismaCloudComputeAPIClient) (policyAPI.Policy, diag.Diagnostics) {
	util.DLog(ctx, "Executing PolicySchemaToTerraform")

	var (
		diags diag.Diagnostics
		rules []policyAPI.PolicyRule
	)

	if plan.Rules != nil {
		rules, diags = PolicyRulesSchemaToTerraform(ctx, plan.Type.ValueString(), *plan.Rules, client)
		if diags.HasError() {
			return policyAPI.Policy{}, diags
		}
	} else {
		rules = []policyAPI.PolicyRule{}
	}

	tfPolicy := policyAPI.Policy{
		Id:            plan.Id.ValueString(),
		PolicyType:    plan.PolicyType.ValueString(),
		PolicyContext: plan.PolicyContext.ValueString(),
		Rules:         &rules,
		Type:          plan.Type.ValueString(),
	}

	tfPolicy.SortRules(ctx, plan.Rules)

	util.DLog(ctx, "Finishing PolicySchemaToTerraform execution")

	return tfPolicy, diags
}

func PolicyRulesSchemaToTerraform(ctx context.Context, moduleType string, schemaRules []models.PolicyRuleResourceModel, client *api.PrismaCloudComputeAPIClient) ([]policyAPI.PolicyRule, diag.Diagnostics) {
	util.DLog(ctx, "Executing PolicyRulesSchemaToTerraform")

	var diags diag.Diagnostics

	rules := []policyAPI.PolicyRule{}

	for _, schemaRule := range schemaRules {
		collectionNames := []string{}
		diags = schemaRule.Collections.ElementsAs(ctx, &collectionNames, false)
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

		for idx := range collections {
			collections[idx].Modified = time.Now().Format("2006-01-02T15:04:05.000Z")
		}

		// TODO: put this in a validator
		//if schemaRule.Effect.ValueString() == "alert, block" && schemaRule.Condition.IsUnknown() {
		//    diags.AddError(
		//        "Missing condition from \"alert, block\" effect rule",
		//        "Condition attribute must be defined for rules with effect \"alert, block\".",
		//    )
		//    return rules, diags
		//}

		alertThreshold := policyAPI.AlertThreshold{}
		blockThreshold := policyAPI.BlockThreshold{}
		condition := policyAPI.Condition{}
		graceDaysPolicy := policyAPI.GraceDaysPolicy{}
		pkgTypesThresholds := []policyAPI.PkgTypesThreshold{}
		cveRules := []policyAPI.Exception{}
		tags := []policyAPI.Exception{}
		riskFactorsEffects := []policyAPI.RiskFactorsEffect{}

		if moduleType == "vulnerability" {
			tfAlertThreshold, tfAlertRiskFactorsEffects, diags := alertThresholdSchemaToTerraform(ctx, schemaRule.AlertThreshold)
			if diags.HasError() {
				return rules, diags
			}
			alertThreshold = tfAlertThreshold
			riskFactorsEffects = append(riskFactorsEffects, tfAlertRiskFactorsEffects...)

			tfBlockThreshold, tfBlockRiskFactorsEffects, diags := blockThresholdSchemaToTerraform(ctx, schemaRule.BlockThreshold)
			if diags.HasError() {
				return rules, diags
			}
			blockThreshold = tfBlockThreshold
			riskFactorsEffects = append(riskFactorsEffects, tfBlockRiskFactorsEffects...)
		}

		if schemaRule.Condition != nil {
			condition, diags = conditionToTerraform(ctx, schemaRule.Condition)
			if diags.HasError() {
				return rules, diags
			}
		}

		if schemaRule.GraceDaysPolicy != nil {
			graceDaysPolicy, diags = graceDaysPolicyToTerraform(ctx, *schemaRule.GraceDaysPolicy)
			if diags.HasError() {
				return rules, diags
			}
		}

		if schemaRule.PkgTypesThresholds != nil {
			pkgTypesThresholds, diags = pkgTypesThresholdsToTerraform(ctx, *schemaRule.PkgTypesThresholds)
			if diags.HasError() {
				return rules, diags
			}
		}

		if schemaRule.CVERules != nil {
			cveRules, diags = exceptionsToTerraform(ctx, *schemaRule.CVERules, false)
			if diags.HasError() {
				return rules, diags
			}
		}

		if schemaRule.Tags != nil {
			tags, diags = exceptionsToTerraform(ctx, *schemaRule.Tags, true)
			if diags.HasError() {
				return rules, diags
			}
		}

		rule := policyAPI.PolicyRule{
			AlertThreshold:                 alertThreshold,
			BlockMessage:                   schemaRule.BlockMessage.ValueString(),
			BlockThreshold:                 blockThreshold,
			Collections:                    collections,
			Condition:                      &condition,
			CVERules:                       cveRules,
			Disabled:                       schemaRule.Disabled.ValueBool(),
			Effect:                         schemaRule.Effect.ValueString(),
			ExcludeBaseImageVulns:          schemaRule.ExcludeBaseImageVulns.ValueBool(),
			GraceDays:                      int(schemaRule.GraceDays.ValueInt32()),
			GraceDaysPolicy:                graceDaysPolicy,
			Modified:                       time.Now().Format("2006-01-02T15:04:05.000Z"),
			Name:                           schemaRule.Name.ValueString(),
			Notes:                          schemaRule.Notes.ValueString(),
			Order:                          int(schemaRule.Order.ValueInt32()),
			OnlyFixed:                      schemaRule.OnlyFixed.ValueBool(),
			PkgTypesThresholds:             pkgTypesThresholds,
			ReportAllPassedAndFailedChecks: schemaRule.ReportAllPassedAndFailedChecks.ValueBool(),
			RiskFactorEffects:              riskFactorsEffects,
			Tags:                           tags,
			Verbose:                        schemaRule.Verbose.ValueBool(),
		}

		rules = append(rules, rule)
	}

	util.DLog(ctx, "Finishing PolicyRulesSchemaToTerraform execution")

	return rules, diags
}

func PolicyTerraformToSchema(ctx context.Context, policy policyAPI.Policy, plan models.PolicyResourceModel) (models.PolicyResourceModel, diag.Diagnostics) {
	util.DLog(ctx, "Executing PolicyTerraformToSchema")

	var diags diag.Diagnostics

	var (
		rules         []models.PolicyRuleResourceModel
		policyContext basetypes.StringValue
		moduleType    basetypes.StringValue
	)

	policyContext, moduleType, diags = getPolicyContextAndModuleType(policy.PolicyType)
	if diags.HasError() {
		return models.PolicyResourceModel{}, diags
	}

	if policy.Rules != nil {
		rules, diags = PolicyRulesTerraformToSchema(ctx, moduleType.ValueString(), *policy.Rules, plan.Rules)
		if diags.HasError() {
			return models.PolicyResourceModel{}, diags
		}
	} else {
		rules = []models.PolicyRuleResourceModel{}
	}

	schema := models.PolicyResourceModel{
		Id:            types.StringValue(policy.Id),
		PolicyType:    types.StringValue(policy.PolicyType),
		PolicyContext: policyContext,
		Rules:         &rules,
		Type:          moduleType,
	}

	schema.SortRules(ctx, plan.Rules)

	util.DLog(ctx, "Finishing PolicyTerraformToSchema execution")

	return schema, diags
}

func PolicyRulesTerraformToSchema(ctx context.Context, moduleType string, rules []policyAPI.PolicyRule, planRules *[]models.PolicyRuleResourceModel) ([]models.PolicyRuleResourceModel, diag.Diagnostics) {
	util.DLog(ctx, "Executing PolicyRulesTerraformToSchema")

	var diags diag.Diagnostics

	schemaRules := []models.PolicyRuleResourceModel{}

	if len(rules) == 0 {
		return schemaRules, diags
	}

	for _, rule := range rules {
		// Find the matching plan rule
		var planRule models.PolicyRuleResourceModel
		for idx, pRule := range *planRules {
			if pRule.Name.ValueString() == rule.Name {
				planRule = (*planRules)[idx]
				break
			}
		}

		collectionNames := []string{}
		for _, collection := range rule.Collections {
			collectionNames = append(collectionNames, collection.Name)
		}

		collections, diags := types.SetValueFrom(ctx, types.StringType, collectionNames)
		if diags.HasError() {
			return []models.PolicyRuleResourceModel{}, diags
		}

		var (
			alertThreshold *models.PolicyRuleThresholdResourceModel
			blockThreshold *models.PolicyRuleThresholdResourceModel
		)

		if moduleType == "vulnerability" {
			alertThreshold, diags = alertThresholdTerraformToSchema(ctx, rule.AlertThreshold, planRule.AlertThreshold, rule.RiskFactorEffects)
			if diags.HasError() {
				return []models.PolicyRuleResourceModel{}, diags
			}

			blockThreshold, diags = blockThresholdTerraformToSchema(ctx, rule.BlockThreshold, planRule.BlockThreshold, rule.RiskFactorEffects)
			if diags.HasError() {
				return []models.PolicyRuleResourceModel{}, diags
			}
		} else {
			alertThreshold = nil
			blockThreshold = nil
		}

		condition, diags := conditionToSchema(ctx, *rule.Condition)
		if diags.HasError() {
			return []models.PolicyRuleResourceModel{}, diags
		}

		var graceDaysPolicy *models.PolicyRuleGraceDaysPolicyResourceModel
		if planRule.GraceDaysPolicy == nil {
			graceDaysPolicy = nil
		} else {
			schemaGraceDaysPolicy, diags := graceDaysPolicyToSchema(ctx, rule.GraceDaysPolicy, planRule.GraceDaysPolicy)
			if diags.HasError() {
				return []models.PolicyRuleResourceModel{}, diags
			}
			graceDaysPolicy = &schemaGraceDaysPolicy
		}

		var pkgTypesThresholds *[]models.PolicyRulePkgTypesThresholdResourceModel
		if planRule.PkgTypesThresholds == nil {
			pkgTypesThresholds = nil
		} else {
			schemaPkgTypesThresholds, diags := pkgTypesThresholdsToSchema(ctx, rule.PkgTypesThresholds)
			if diags.HasError() {
				return []models.PolicyRuleResourceModel{}, diags
			}
			pkgTypesThresholds = &schemaPkgTypesThresholds
		}

		cveRules, diags := exceptionsToSchema(ctx, rule.CVERules, false)
		if diags.HasError() {
			return []models.PolicyRuleResourceModel{}, diags
		}

		tags, diags := exceptionsToSchema(ctx, rule.Tags, true)
		if diags.HasError() {
			return []models.PolicyRuleResourceModel{}, diags
		}

		//if rule.Effect == "alert, block" {
		//    rule.Effect = "block"
		//}

		var excludeBaseImageVulns basetypes.BoolValue
		if planRule.ExcludeBaseImageVulns.IsNull() {
			excludeBaseImageVulns = types.BoolNull()
		} else {
			excludeBaseImageVulns = types.BoolValue(rule.ExcludeBaseImageVulns)
		}

		var graceDays basetypes.Int32Value
		if planRule.GraceDays.IsNull() {
			graceDays = types.Int32Null()
		} else {
			graceDays = types.Int32Value(int32(rule.GraceDays))
		}

		var notes basetypes.StringValue
		if planRule.Notes.IsNull() {
			notes = types.StringNull()
		} else {
			notes = types.StringValue(rule.Notes)
		}

		var onlyFixed basetypes.BoolValue
		if planRule.OnlyFixed.IsNull() {
			onlyFixed = types.BoolNull()
		} else {
			onlyFixed = types.BoolValue(rule.OnlyFixed)
		}

		var reportAllPassedAndFailedChecks basetypes.BoolValue
		if planRule.ReportAllPassedAndFailedChecks.IsNull() {
			reportAllPassedAndFailedChecks = types.BoolNull()
		} else {
			reportAllPassedAndFailedChecks = types.BoolValue(rule.ReportAllPassedAndFailedChecks)
		}

		schemaRule := models.PolicyRuleResourceModel{
			AlertThreshold:                 alertThreshold,
			BlockMessage:                   types.StringValue(rule.BlockMessage),
			BlockThreshold:                 blockThreshold,
			Collections:                    collections,
			Condition:                      condition,
			CVERules:                       cveRules,
			Disabled:                       types.BoolValue(rule.Disabled),
			Effect:                         types.StringValue(rule.Effect),
			ExcludeBaseImageVulns:          excludeBaseImageVulns,
			GraceDays:                      graceDays,
			GraceDaysPolicy:                graceDaysPolicy,
			Modified:                       types.StringValue(""),
			Name:                           types.StringValue(rule.Name),
			Notes:                          notes,
			OnlyFixed:                      onlyFixed,
			Owner:                          types.StringValue(rule.Owner),
			PkgTypesThresholds:             pkgTypesThresholds,
			ReportAllPassedAndFailedChecks: reportAllPassedAndFailedChecks,
			Tags:                           tags,
			Verbose:                        types.BoolValue(rule.Verbose),
		}

		schemaRules = append(schemaRules, schemaRule)
	}

	util.DLog(ctx, "Finishing PolicyRulesTerraformToSchema exection")

	return schemaRules, diags
}

func getPolicyContextAndModuleType(policyType string) (basetypes.StringValue, basetypes.StringValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch policyType {
	case policyAPI.PolicyTypeComplianceContainer:
		return types.StringValue(policyAPI.PolicyContextComplianceContainer), types.StringValue(policyAPI.TypeCompliance), diags
	case policyAPI.PolicyTypeComplianceCiImage:
		return types.StringValue(policyAPI.PolicyContextComplianceCiImage), types.StringValue(policyAPI.TypeCompliance), diags
	case policyAPI.PolicyTypeComplianceHost:
		return types.StringValue(policyAPI.PolicyContextComplianceHost), types.StringValue(policyAPI.TypeCompliance), diags
	case policyAPI.PolicyTypeComplianceVmImage:
		return types.StringValue(policyAPI.PolicyContextComplianceVmImage), types.StringValue(policyAPI.TypeCompliance), diags
	case policyAPI.PolicyTypeComplianceFunction:
		return types.StringValue(policyAPI.PolicyContextComplianceFunction), types.StringValue(policyAPI.TypeCompliance), diags
	case policyAPI.PolicyTypeComplianceCiFunction:
		return types.StringValue(policyAPI.PolicyContextComplianceCiFunction), types.StringValue(policyAPI.TypeCompliance), diags
	case policyAPI.PolicyTypeVulnerabilityDeployedImage:
		return types.StringValue(policyAPI.PolicyContextVulnerabilityDeployedImage), types.StringValue(policyAPI.TypeVulnerability), diags
	case policyAPI.PolicyTypeVulnerabilityCiImage:
		return types.StringValue(policyAPI.PolicyContextVulnerabilityCiImage), types.StringValue(policyAPI.TypeVulnerability), diags
	default:
		diags.AddError(
			"Value Conversion Error", // TODO: probably need a better error type here
			fmt.Sprintf("unknown policy type value \"%s\"", policyType),
		)
		return types.StringNull(), types.StringNull(), diags
	}
}

func thresholdValueToSchema(value int) (basetypes.StringValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch value {
	case 0:
		return types.StringValue("off"), diags
	case 1:
		return types.StringValue("low"), diags
	case 4:
		return types.StringValue("medium"), diags
	case 7:
		return types.StringValue("high"), diags
	case 9:
		return types.StringValue("critical"), diags
	default:
		diags.AddError(
			"Value Conversion Error",
			fmt.Sprintf("unknown threshold value \"%d\"", value),
		)
		return types.StringValue(""), diags
	}
}

func thresholdValueToTerraform(value string) (int, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch value {
	case "off":
		return 0, true, diags
	case "low":
		return 1, false, diags
	case "medium":
		return 4, false, diags
	case "high":
		return 7, false, diags
	case "critical":
		return 9, false, diags
	default:
		diags.AddError(
			"Value Conversion Error",
			fmt.Sprintf("unknown alert threshold value \"%s\"", value),
		)
		return 0, false, diags
	}
}

func alertThresholdTerraformToSchema(ctx context.Context, tfAlertThreshold policyAPI.AlertThreshold, planAlertThreshold *models.PolicyRuleThresholdResourceModel, tfRiskFactorsEffects []policyAPI.RiskFactorsEffect) (*models.PolicyRuleThresholdResourceModel, diag.Diagnostics) {
	var (
		diags          diag.Diagnostics
		riskFactorsSet basetypes.SetValue
	)

	threshold, diags := thresholdValueToSchema(tfAlertThreshold.Value)
	if diags.HasError() {
		return nil, diags
	}

	if planAlertThreshold != nil && planAlertThreshold.RiskFactors.IsNull() {
		riskFactorsSet = types.SetNull(types.StringType)
	} else {
		riskFactors := []attr.Value{}
		for _, tfRiskFactorEffect := range tfRiskFactorsEffects {
			if tfRiskFactorEffect.Effect == "alert" {
				riskFactors = append(riskFactors, types.StringValue(tfRiskFactorEffect.RiskFactor))
			}
		}
		riskFactorsSet = types.SetValueMust(types.StringType, riskFactors)
	}

	return &models.PolicyRuleThresholdResourceModel{
		Threshold:   threshold,
		RiskFactors: riskFactorsSet,
	}, diags
}

func blockThresholdTerraformToSchema(ctx context.Context, tfBlockThreshold policyAPI.BlockThreshold, planBlockThreshold *models.PolicyRuleThresholdResourceModel, tfRiskFactorsEffects []policyAPI.RiskFactorsEffect) (*models.PolicyRuleThresholdResourceModel, diag.Diagnostics) {
	var (
		diags          diag.Diagnostics
		riskFactorsSet basetypes.SetValue
	)

	threshold, diags := thresholdValueToSchema(tfBlockThreshold.Value)
	if diags.HasError() {
		return nil, diags
	}

	if planBlockThreshold != nil && planBlockThreshold.RiskFactors.IsNull() {
		riskFactorsSet = types.SetNull(types.StringType)
	} else {
		riskFactors := []attr.Value{}
		for _, tfRiskFactorEffect := range tfRiskFactorsEffects {
			if tfRiskFactorEffect.Effect == "block" {
				riskFactors = append(riskFactors, types.StringValue(tfRiskFactorEffect.RiskFactor))
			}
		}
		riskFactorsSet = types.SetValueMust(types.StringType, riskFactors)
	}

	return &models.PolicyRuleThresholdResourceModel{
		Threshold:   threshold,
		RiskFactors: riskFactorsSet,
	}, diags
}

func alertThresholdSchemaToTerraform(ctx context.Context, schemaAlertThreshold *models.PolicyRuleThresholdResourceModel) (policyAPI.AlertThreshold, []policyAPI.RiskFactorsEffect, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		alertThreshold     policyAPI.AlertThreshold
		riskFactorsEffects []policyAPI.RiskFactorsEffect
	)

	if schemaAlertThreshold == nil {
		return policyAPI.AlertThreshold{
				Value:    0,
				Disabled: true,
			},
			[]policyAPI.RiskFactorsEffect{},
			diags
	}

	value, disabled, diags := thresholdValueToTerraform(schemaAlertThreshold.Threshold.ValueString())
	if diags.HasError() {
		return policyAPI.AlertThreshold{}, []policyAPI.RiskFactorsEffect{}, diags
	}
	alertThreshold = policyAPI.AlertThreshold{
		Value:    value,
		Disabled: disabled,
	}

	var riskFactors []string
	if schemaAlertThreshold.RiskFactors.IsNull() {
		riskFactors = []string{}
	} else {
		diags = schemaAlertThreshold.RiskFactors.ElementsAs(ctx, &riskFactors, false)
		if diags.HasError() {
			return policyAPI.AlertThreshold{}, []policyAPI.RiskFactorsEffect{}, diags
		}
	}

	riskFactorsEffects = []policyAPI.RiskFactorsEffect{}
	for _, riskFactor := range riskFactors {
		riskFactorsEffects = append(riskFactorsEffects, policyAPI.RiskFactorsEffect{
			Effect:     "alert",
			RiskFactor: riskFactor,
		})
	}

	return alertThreshold, riskFactorsEffects, diags
}

func blockThresholdSchemaToTerraform(ctx context.Context, schemaBlockThreshold *models.PolicyRuleThresholdResourceModel) (policyAPI.BlockThreshold, []policyAPI.RiskFactorsEffect, diag.Diagnostics) {
	var (
		diags              diag.Diagnostics
		blockThreshold     policyAPI.BlockThreshold
		riskFactorsEffects []policyAPI.RiskFactorsEffect
	)

	if schemaBlockThreshold == nil {
		return policyAPI.BlockThreshold{
				Value:   0,
				Enabled: true,
			},
			[]policyAPI.RiskFactorsEffect{},
			diags
	}

	value, disabled, diags := thresholdValueToTerraform(schemaBlockThreshold.Threshold.ValueString())
	if diags.HasError() {
		return policyAPI.BlockThreshold{}, []policyAPI.RiskFactorsEffect{}, diags
	}
	blockThreshold = policyAPI.BlockThreshold{
		Value:   value,
		Enabled: !disabled,
	}

	var riskFactors []string
	if schemaBlockThreshold.RiskFactors.IsNull() {
		riskFactors = []string{}
	} else {
		diags = schemaBlockThreshold.RiskFactors.ElementsAs(ctx, &riskFactors, false)
		if diags.HasError() {
			return policyAPI.BlockThreshold{}, []policyAPI.RiskFactorsEffect{}, diags
		}
	}

	riskFactorsEffects = []policyAPI.RiskFactorsEffect{}
	for _, riskFactor := range riskFactors {
		riskFactorsEffects = append(riskFactorsEffects, policyAPI.RiskFactorsEffect{
			Effect:     "block",
			RiskFactor: riskFactor,
		})
	}

	return blockThreshold, riskFactorsEffects, diags
}

func conditionToTerraform(ctx context.Context, schemaCondition *models.PolicyRuleConditionResourceModel) (policyAPI.Condition, diag.Diagnostics) {
	var diags diag.Diagnostics

	if schemaCondition == nil {
		return policyAPI.Condition{
			Vulnerabilities: []policyAPI.Vulnerability{},
		}, diags
	}

	tfConditionVulnerabilities := []policyAPI.Vulnerability{}
	for _, schemaConditionVulnerability := range (*schemaCondition).Vulnerabilities {
		tfConditionVulnerabilities = append(tfConditionVulnerabilities, policyAPI.Vulnerability{
			Id:    int(schemaConditionVulnerability.Id.ValueInt32()),
			Block: schemaConditionVulnerability.Block.ValueBool(),
		})
	}

	tfCondition := policyAPI.Condition{
		Vulnerabilities: tfConditionVulnerabilities,
	}

	return tfCondition, diags
}

func conditionToSchema(ctx context.Context, tfCondition policyAPI.Condition) (*models.PolicyRuleConditionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	schemaConditionVulnerabilities := []models.PolicyRuleConditionVulnerabilityResourceModel{}
	for _, tfConditionVulnerability := range tfCondition.Vulnerabilities {
		schemaConditionVulnerabilities = append(schemaConditionVulnerabilities, models.PolicyRuleConditionVulnerabilityResourceModel{
			Id:    types.Int32Value(int32(tfConditionVulnerability.Id)),
			Block: types.BoolValue(tfConditionVulnerability.Block),
		})
	}

	schemaCondition := models.PolicyRuleConditionResourceModel{
		Vulnerabilities: schemaConditionVulnerabilities,
	}

	return &schemaCondition, diags
}

func graceDaysPolicyToTerraform(ctx context.Context, schemaGraceDaysPolicy models.PolicyRuleGraceDaysPolicyResourceModel) (policyAPI.GraceDaysPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	// TODO: set default value for each count to 0
	low := int(schemaGraceDaysPolicy.Low.ValueInt32())
	medium := int(schemaGraceDaysPolicy.Medium.ValueInt32())
	high := int(schemaGraceDaysPolicy.High.ValueInt32())
	critical := int(schemaGraceDaysPolicy.Critical.ValueInt32())
	enabled := (low != 0 || medium != 0 || high != 0 || critical != 0)

	tfGraceDaysPolicy := policyAPI.GraceDaysPolicy{
		Enabled:  enabled,
		Low:      low,
		Medium:   medium,
		High:     high,
		Critical: critical,
	}

	return tfGraceDaysPolicy, diags
}

func graceDaysPolicyToSchema(ctx context.Context, tfGraceDaysPolicy policyAPI.GraceDaysPolicy, planGraceDaysPolicy *models.PolicyRuleGraceDaysPolicyResourceModel) (models.PolicyRuleGraceDaysPolicyResourceModel, diag.Diagnostics) {
	var (
		diags    diag.Diagnostics
		low      basetypes.Int32Value
		medium   basetypes.Int32Value
		high     basetypes.Int32Value
		critical basetypes.Int32Value
	)

	planPolicyIsNil := (planGraceDaysPolicy == nil)
	if !planPolicyIsNil && (*planGraceDaysPolicy).Low.IsNull() {
		low = types.Int32Null()
	} else {
		low = types.Int32Value(int32(tfGraceDaysPolicy.Low))
	}

	if !planPolicyIsNil && (*planGraceDaysPolicy).Medium.IsNull() {
		medium = types.Int32Null()
	} else {
		medium = types.Int32Value(int32(tfGraceDaysPolicy.Medium))
	}

	if !planPolicyIsNil && (*planGraceDaysPolicy).High.IsNull() {
		high = types.Int32Null()
	} else {
		high = types.Int32Value(int32(tfGraceDaysPolicy.High))
	}

	if !planPolicyIsNil && (*planGraceDaysPolicy).Critical.IsNull() {
		critical = types.Int32Null()
	} else {
		critical = types.Int32Value(int32(tfGraceDaysPolicy.Critical))
	}

	return models.PolicyRuleGraceDaysPolicyResourceModel{
		Low:      low,
		Medium:   medium,
		High:     high,
		Critical: critical,
	}, diags
}

func pkgTypesThresholdsToTerraform(ctx context.Context, schemaPkgTypesThresholds []models.PolicyRulePkgTypesThresholdResourceModel) ([]policyAPI.PkgTypesThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfPkgTypesThresholds := []policyAPI.PkgTypesThreshold{}
	for _, schemaPkgTypesThreshold := range schemaPkgTypesThresholds {
		alertThresholdDisabled, blockThresholdEnabled := false, true
		alertThresholdValue, blockThresholdValue := 0, 0

		switch schemaPkgTypesThreshold.AlertThreshold.ValueString() {
		case "off":
			alertThresholdDisabled = true
		case "low":
			alertThresholdValue = 1
		case "medium":
			alertThresholdValue = 4
		case "high":
			alertThresholdValue = 7
		case "critical":
			alertThresholdValue = 9
		}

		switch schemaPkgTypesThreshold.BlockThreshold.ValueString() {
		case "off":
			blockThresholdEnabled = false
		case "low":
			blockThresholdValue = 1
		case "medium":
			blockThresholdValue = 4
		case "high":
			blockThresholdValue = 7
		case "critical":
			blockThresholdValue = 9
		}

		tfPkgTypesThresholds = append(tfPkgTypesThresholds, policyAPI.PkgTypesThreshold{
			Type: schemaPkgTypesThreshold.Type.ValueString(),
			AlertThreshold: policyAPI.AlertThreshold{
				Disabled: alertThresholdDisabled,
				Value:    alertThresholdValue,
			},
			BlockThreshold: policyAPI.BlockThreshold{
				Enabled: blockThresholdEnabled,
				Value:   blockThresholdValue,
			},
		})
	}

	return tfPkgTypesThresholds, diags
}

func pkgTypesThresholdsToSchema(ctx context.Context, tfPkgTypesThresholds []policyAPI.PkgTypesThreshold) ([]models.PolicyRulePkgTypesThresholdResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	schemaPkgTypesThresholds := []models.PolicyRulePkgTypesThresholdResourceModel{}
	for _, tfPkgTypesThreshold := range tfPkgTypesThresholds {
		alertThreshold, blockThreshold := "off", "off"

		switch tfPkgTypesThreshold.AlertThreshold.Value {
		case 1:
			alertThreshold = "low"
		case 4:
			alertThreshold = "medium"
		case 7:
			alertThreshold = "high"
		case 9:
			alertThreshold = "critical"
		}

		switch tfPkgTypesThreshold.BlockThreshold.Value {
		case 1:
			blockThreshold = "low"
		case 4:
			blockThreshold = "medium"
		case 7:
			blockThreshold = "high"
		case 9:
			blockThreshold = "critical"
		}

		schemaPkgTypesThresholds = append(schemaPkgTypesThresholds, models.PolicyRulePkgTypesThresholdResourceModel{
			Type:           types.StringValue(tfPkgTypesThreshold.Type),
			AlertThreshold: types.StringValue(alertThreshold),
			BlockThreshold: types.StringValue(blockThreshold),
		})
	}

	return schemaPkgTypesThresholds, diags
}

func exceptionsToTerraform(ctx context.Context, schemaExceptions []models.PolicyRuleExceptionResourceModel, isTag bool) ([]policyAPI.Exception, diag.Diagnostics) {
	var diags diag.Diagnostics

	var exceptionType string
	if isTag {
		exceptionType = "tag"
	} else {
		exceptionType = "cve"
	}

	tfExceptions := []policyAPI.Exception{}

	for _, schemaException := range schemaExceptions {
		tfExceptions = append(tfExceptions, policyAPI.Exception{
			//Name: schemaException.Name.ValueString(),
			Name:        schemaException.Id.ValueString(),
			Effect:      schemaException.Effect.ValueString(),
			Id:          schemaException.Id.ValueString(),
			Description: schemaException.Description.ValueString(),
			Type:        exceptionType,
			Expiration: policyAPI.ExceptionExpiration{
				Enabled: schemaException.Expiration.Enabled.ValueBool(),
				Date:    schemaException.Expiration.Date.ValueString(),
			},
		})
	}

	return tfExceptions, diags
}

func exceptionsToSchema(ctx context.Context, tfExceptions []policyAPI.Exception, isTag bool) (*[]models.PolicyRuleExceptionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(tfExceptions) == 0 {
		return nil, diags
	}

	var exceptionId basetypes.StringValue

	schemaExceptions := []models.PolicyRuleExceptionResourceModel{}

	for _, tfException := range tfExceptions {
		if isTag {
			exceptionId = types.StringValue(tfException.Name)
		} else {
			exceptionId = types.StringValue(tfException.Id)
		}

		schemaExceptions = append(schemaExceptions, models.PolicyRuleExceptionResourceModel{
			//Name: types.StringValue(tfException.Name),
			Effect:      types.StringValue(tfException.Effect),
			Id:          exceptionId,
			Description: types.StringValue(tfException.Description),
			//Type: exceptionType,
			Expiration: models.PolicyRuleExceptionExpirationResourceModel{
				Enabled: types.BoolValue(tfException.Expiration.Enabled),
				Date:    types.StringValue(tfException.Expiration.Date),
			},
		})
	}

	return &schemaExceptions, diags
}

func GeneratePolicyRulesOrderMap(rules []models.PolicyRuleResourceModel) map[string]int {
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
	for key := range orderedRulesMap {
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

func SchemaRuleOrderIsRestored(ctx context.Context, planRules *[]models.PolicyRuleResourceModel, createdRules *[]models.PolicyRuleResourceModel) bool {
	orderIsRestored := true

	for index := range *planRules {
		planRuleName, createdRuleName := (*planRules)[index].Name.ValueString(), (*createdRules)[index].Name.ValueString()
		if planRuleName != createdRuleName {
			orderIsRestored = false
		}
	}

	return orderIsRestored
}

func isWildcard(val []string) bool {
	return len(val) == 1 && val[0] == "*"
}

func validateRuleCollectionsByPolicyType(ctx context.Context, policyType string, ruleName string, ruleCollectionNames []string, collections []collectionAPI.Collection) diag.Diagnostics {
	var diags diag.Diagnostics

	// Define what collection fields must be wildcards according to policy type
	policyValidations := map[string][]string{
		"containerCompliance":    {"Functions"},
		"ciImagesCompliance":     {"Containers", "Hosts", "AppIDs", "Functions", "Namespaces", "AccountIDs", "Clusters"},
		"hostCompliance":         {"Containers", "Images", "AppIDs", "Functions", "Namespaces"},
		"vmCompliance":           {"Containers", "Hosts", "AppIDs", "Functions", "Namespaces", "Clusters"},
		"serverlessCompliance":   {"Containers", "Hosts", "Images", "AppIDs", "Namespaces", "Clusters"},
		"ciServerlessCompliance": {"Containers", "Hosts", "Images", "AppIDs", "Namespaces", "AccountIDs", "Clusters"},
		// "trust":  add rules for "trust" when implemented
		"containerVulnerability": {"Functions"},
		"ciImagesVulnerability": {"Containers", "Hosts", "AppIDs", "Functions", "Namespaces", "AccountIDs", "Clusters"},
	}

	// Retrieve the validation fields for the given policyType
	fieldsToCheck, ok := policyValidations[policyType]
	if !ok {
		diags.AddError(
			"Resource Validation Error",
			fmt.Sprintf("While validating collections for policy rule \"%s\", encountered unknown policy type \"%s\"", ruleName, policyType),
		)
		return diags
	}

	// Iterate over the collection names in rule config
	for _, ruleCollectionName := range ruleCollectionNames {
		found := false
		invalidFields := make([]string, 0)

		// Iterate over the collections from the console to find the collection with the current ruleCollectionName
		for _, collection := range collections {
			if ruleCollectionName == collection.Name {
				found = true

				// Check each field specified in fieldsToCheck for non-wildcard values, adding the field name to invalidFields if
				// the value is not a slice with a single wildcard string
				for _, field := range fieldsToCheck {
					switch field {
					case "Containers":
						if !isWildcard(collection.Containers) {
							invalidFields = append(invalidFields, "Containers")
						}
					case "Hosts":
						if !isWildcard(collection.Hosts) {
							invalidFields = append(invalidFields, "Hosts")
						}
					case "Images":
						if !isWildcard(collection.Images) {
							invalidFields = append(invalidFields, "Images")
						}
					case "Labels":
						if !isWildcard(collection.Labels) {
							invalidFields = append(invalidFields, "Labels")
						}
					case "AppIDs":
						if !isWildcard(collection.AppIDs) {
							invalidFields = append(invalidFields, "AppIDs")
						}
					case "Functions":
						if !isWildcard(collection.Functions) {
							invalidFields = append(invalidFields, "Functions")
						}
					case "Namespaces":
						if !isWildcard(collection.Namespaces) {
							invalidFields = append(invalidFields, "Namespaces")
						}
					case "AccountIDs":
						if !isWildcard(collection.AccountIDs) {
							invalidFields = append(invalidFields, "AccountIDs")
						}
					case "Clusters":
						if !isWildcard(collection.Clusters) {
							invalidFields = append(invalidFields, "Clusters")
						}
					}
				}

				// Append error if the collection is not found in the console
				if !found {
					diags.AddError(
						"Resource Validation Error",
						fmt.Sprintf("Error occured during validation of collections for policy rule \"%s\": Collection name \"%s\" not found", ruleName, policyType),
					)
					// Otherwise, if there's a non-zero amount of invalid fields, append an error with the field names
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

func generateConditionFromEffect(ctx context.Context, policyType string, name string, effect string, index int, plan tfsdk.Plan, complianceVulnerabilities []systemAPI.Vulnerability) (models.PolicyRuleConditionResourceModel, diag.Diagnostics) {
	util.DLog(ctx, "Executing generateConditionFromEffect")
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

	generatedCondition := models.PolicyRuleConditionResourceModel{}
	conditionVulnerabilities := []models.PolicyRuleConditionVulnerabilityResourceModel{}

	// TODO: add a check to see if we're configuring the condition instead of
	// just checking to see if the effect is "alert, block"

	if policyType == "serverlessCompliance" || policyType == "ciServerlessCompliance" {
		var configCondition basetypes.ObjectValue
		diags.Append(plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), &configCondition)...)
		//if configCondition.IsUnknown() || configCondition.IsNull() {
		//if effect == "alert, block"
		if effect != "ignore" {
			if effect == "default" {
				complianceVulnerabilities = systemAPI.GetHighOrCriticalVulnerabilities(complianceVulnerabilities)
			}
			// If this block executes and effect is default, effect must be set to "alert"
			for _, vuln := range complianceVulnerabilities {
				conditionVulnerability := models.PolicyRuleConditionVulnerabilityResourceModel{
					Id:    types.Int32Value(int32(vuln.Id)),
					Block: types.BoolValue(false),
				}
				conditionVulnerabilities = append(conditionVulnerabilities, conditionVulnerability)
			}
		}
		//}
	} else {
		// If the effect is "alert, block", create condition vulnerabilities object from plan
		if effect == "alert, block" {
			// TODO: add validator that ensures a condition is configured if the effect is set to "alert, block"
			var configVulnerabilities []models.PolicyRuleConditionVulnerabilityResourceModel
			diags.Append(plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition").AtName("vulnerabilities"), &configVulnerabilities)...)
			if diags.HasError() {
				return generatedCondition, diags
			}
			conditionVulnerabilities = configVulnerabilities
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
					//return conditionObject, diags
					diags.AddError(
						"Policy Rule Condition Conversion Error",
						//fmt.Sprintf("policy rule \"%s\" configured with unknown policy type \"%s\"", rule.Name.ValueString(), policyType),
						fmt.Sprintf("policy rule \"%s\" configured with unknown policy type \"%s\"", name, policyType),
					)
					return generatedCondition, diags
				}

				conditionVulnerability := models.PolicyRuleConditionVulnerabilityResourceModel{
					Id:    types.Int32Value(int32(vuln.Id)),
					Block: types.BoolValue(block),
				}
				conditionVulnerabilities = append(conditionVulnerabilities, conditionVulnerability)
			}
		}
	}

	generatedCondition.Vulnerabilities = conditionVulnerabilities

	util.DLog(ctx, "Finishing generateConditionFromEffect execution")

	return generatedCondition, diags
}

func ModifyPolicyResourcePlan(ctx context.Context, client *api.PrismaCloudComputeAPIClient, plan tfsdk.Plan, resp *resource.ModifyPlanResponse) {
	util.DLog(ctx, "Executing ModifyPolicyResourcePlan")

	var (
		rules                     basetypes.ListValue
		policyTypeValue           basetypes.StringValue
		policyType                string
		moduleTypeValue           basetypes.StringValue
		moduleType                string
		nameValue                 basetypes.StringValue
		name                      string
		collectionNamesValue      basetypes.SetValue
		collectionNames           []string
		effectValue               basetypes.StringValue
		effect                    string
		complianceVulnerabilities []systemAPI.Vulnerability
		err                       error
	)

	resp.Diagnostics.Append(plan.GetAttribute(ctx, path.Root("rules"), &rules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rules.IsNull() {
		return
	}

	// Retrieve policy type
	resp.Diagnostics.Append(plan.GetAttribute(ctx, path.Root("policy_type"), &policyTypeValue)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policyType = policyTypeValue.ValueString()

	// Retrieve module type
	resp.Diagnostics.Append(plan.GetAttribute(ctx, path.Root("type"), &moduleTypeValue)...)
	if resp.Diagnostics.HasError() {
		return
	}
	moduleType = moduleTypeValue.ValueString()

	util.DLog(ctx, "Retrieving vulnerability data")
	if moduleType == "compliance" {
		complianceVulnerabilities, err = systemAPI.GetComplianceVulnerabilities(*client, policyType)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error modifying planned policy rules",
				"Failed to retrieve compliance host vulnerabilities from Prisma Cloud while modifying plan rules: "+err.Error(),
			)
			return
		}
	} else {
		complianceVulnerabilities = []systemAPI.Vulnerability{}
	}

	util.DLog(ctx, "Retrieving collections")
	collections, err := collectionAPI.ListCollections(*client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error modifying planned policy rules",
			"Failed to retrieve collections from Prisma Cloud while modifying plan rules: "+err.Error(),
		)
		return
	}

	util.DLog(ctx, "Beginning loop over rules")

	for index := range rules.Elements() {
		resp.Diagnostics.Append(plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("name"), &nameValue)...)
		name = nameValue.ValueString()
		resp.Diagnostics.Append(plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("effect"), &effectValue)...)
		effect = effectValue.ValueString()

		if resp.Diagnostics.HasError() {
			return
		}

		rulePath := path.Root("rules").AtListIndex(index)

		// Retrieve collection names
		resp.Diagnostics.Append(plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("collections"), &collectionNamesValue)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert collection names to string slice
		resp.Diagnostics.Append(collectionNamesValue.ElementsAs(ctx, &collectionNames, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Validate configured collections to ensure compatibility with policy type
		resp.Diagnostics.Append(validateRuleCollectionsByPolicyType(ctx, policyType, name, collectionNames, collections)...)
		if resp.Diagnostics.HasError() {
			return
		}

		condition, diags := generateConditionFromEffect(ctx, policyType, name, effect, index, plan, complianceVulnerabilities)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if effect == "default" {
			resp.Diagnostics.Append(plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("effect"), "alert")...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		// Set condition to generated value
		resp.Diagnostics.Append(plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("condition"), condition)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Set effect to "alert" if its currently set to the placeholder value ("default")
		var ruleEffectValue basetypes.StringValue
		resp.Diagnostics.Append(plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("effect"), &ruleEffectValue)...)
		effect := ruleEffectValue.ValueString()
		if effect == "default" {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, rulePath.AtName("effect"), types.StringValue("alert"))...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	resp.Plan = plan

	util.DLog(ctx, "Finishing ModifyPolicyResourcePlan execution")
}
