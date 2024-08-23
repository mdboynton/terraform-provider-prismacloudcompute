package provider

//import (
//	"fmt"
//
//	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
//	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
//	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/convert"
//    //"github.com/hashicorp/terraform-plugin-framework/diag"
//    //"github.com/hashicorp/terraform-plugin-framework/path"
//	//"github.com/hashicorp/terraform-plugin-framework/resource"
//	//"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
//	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
//	//"github.com/hashicorp/terraform-plugin-framework/types"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//
////var _ resource.Resource = &HostCompliancePolicy{}
////var _ resource.ResourceWithImportState = &HostCompliancePolicy{}
////
////func NewHostCompliancePolicyResource() resource.Resource {
////    return &HostCompliancePolicy{}
////}
////
////type HostCompliancePolicyResource struct {
////    client *api.PrismaCloudComputeAPIClient
////}
////
////type HostCompliancePolicyResourceModel struct {
////    //Id types.String `tfsdk:"id"`
////    //PolicyType types.String `tfsdk:"policy_type"`
////    Rules *[]HostCompliancePolicyRuleResourceModel `tfsdk:"rules"`
////}
////
////type HostCompliancePolicyRuleResourceModel struct {
////    Action []types.String `tfsdk:"action"`
////    AlertThreshold HostCompliancePolicyRuleAlertThresholdResourceModel `tfsdk:"alert_threshold"`
////    ReportAllCompliance type.Bool `tfsdk:"report_all_compliance"`
////    AuditAllowed type.Bool `tfsdk:"audit_allowed"`
////    BlockMessage type.String `tfsdk:"block_message"`
////    BlockThreshold  HostCompliancePolicyRuleBlockThresholdResourceModel `tfsdk:"block_threshold"`
////    Collections 
////}
//
//func resourcePoliciesComplianceHost() *schema.Resource {
//	return &schema.Resource{
//		Create: createPolicyComplianceHost,
//		Read:   readPolicyComplianceHost,
//		Update: updatePolicyComplianceHost,
//		Delete: deletePolicyComplianceHost,
//
//		Importer: &schema.ResourceImporter{
//			StateContext: schema.ImportStatePassthroughContext,
//		},
//
//		Schema: map[string]*schema.Schema{
//			"id": {
//				Description: "The ID of the policy.",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"rule": {
//				Type:        schema.TypeList,
//				Optional:    true,
//				Description: "Rules that make up the policy.",
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"block_message": {
//							Type:        schema.TypeString,
//							Optional:    true,
//							Description: "Message to display for blocked requests.",
//						},
//						"collections": {
//							Type:        schema.TypeList,
//							Optional:    true,
//							Description: "Collections used to scope the rule.",
//							Elem: &schema.Schema{
//								Type: schema.TypeString,
//							},
//						},
//						"compliance_check": {
//							Type:        schema.TypeList,
//							Optional:    true,
//							Description: "Compliance checks. Omitted checks are ignored.",
//							Elem: &schema.Resource{
//								Schema: map[string]*schema.Schema{
//									"block": {
//										Type:        schema.TypeBool,
//										Optional:    true,
//										Description: "Whether or not to block if this check is failed. Setting to 'false' will only alert if the check is failed.",
//									},
//									"id": {
//										Type:        schema.TypeInt,
//										Optional:    true,
//										Description: "Compliance check number.",
//									},
//								},
//							},
//						},
//						"disabled": {
//							Type:        schema.TypeBool,
//							Optional:    true,
//							Description: "Whether or not to disable the rule.",
//						},
//						"effect": {
//							Type:        schema.TypeString,
//							Optional:    true,
//							Description: "The effect of the rule. Can be set to 'ignore', 'alert', 'block', or 'alert, block'.",
//						},
//						"name": {
//							Type:        schema.TypeString,
//							Optional:    true,
//							Description: "Unique name of the rule.",
//						},
//						"notes": {
//							Type:        schema.TypeString,
//							Optional:    true,
//							Description: "Free-form text field.",
//						},
//						"show_passed_checks": {
//							Type:        schema.TypeBool,
//							Optional:    true,
//							Description: "Whether or not to report both failed and passed compliance checks.",
//						},
//						"verbose": {
//							Type:        schema.TypeBool,
//							Optional:    true,
//							Description: "Whether or not to provide verbose output for blocked requests.",
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func createPolicyComplianceHost(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*api.PrismaCloudComputeAPIClient)
//	parsedRules, err := convert.SchemaToComplianceDeployedRules(d)
//	if err != nil {
//		return fmt.Errorf("error creating %s policy: %s", policyTypeComplianceHost, err)
//	}
//
//	parsedPolicy := policy.CompliancePolicy{
//		Type:  policyTypeComplianceHost,
//		Rules: parsedRules,
//	}
//
//	if err := policy.UpdateComplianceHost(*client, parsedPolicy); err != nil {
//		return fmt.Errorf("error creating %s policy: %s", policyTypeComplianceHost, err)
//	}
//
//	d.SetId(policyTypeComplianceHost)
//	return readPolicyComplianceHost(d, meta)
//}
//
//func readPolicyComplianceHost(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*api.PrismaCloudComputeAPIClient)
//	retrievedPolicy, err := policy.GetComplianceHost(*client)
//	if err != nil {
//		return fmt.Errorf("error reading %s policy: %s", policyTypeComplianceHost, err)
//	}
//
//	if err := d.Set("rule", convert.ComplianceDeployedRulesToSchema(retrievedPolicy.Rules)); err != nil {
//		return fmt.Errorf("error reading %s policy: %s", policyTypeComplianceHost, err)
//	}
//	return nil
//}
//
//func updatePolicyComplianceHost(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*api.PrismaCloudComputeAPIClient)
//	parsedRules, err := convert.SchemaToComplianceDeployedRules(d)
//	if err != nil {
//		return fmt.Errorf("error updating %s policy: %s", policyTypeComplianceHost, err)
//	}
//
//	parsedPolicy := policy.CompliancePolicy{
//		Type:  policyTypeComplianceHost,
//		Rules: parsedRules,
//	}
//
//	if err := policy.UpdateComplianceHost(*client, parsedPolicy); err != nil {
//		return fmt.Errorf("error updating %s policy: %s", policyTypeComplianceHost, err)
//	}
//
//	return readPolicyComplianceHost(d, meta)
//}
//
//func deletePolicyComplianceHost(d *schema.ResourceData, meta interface{}) error {
//	// TODO: reset to default policy
//	return nil
//}
