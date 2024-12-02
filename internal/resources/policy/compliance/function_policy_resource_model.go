package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &FunctionCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &FunctionCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &FunctionCompliancePolicyResource{}

func NewFunctionCompliancePolicyResource() resource.Resource {
    return &FunctionCompliancePolicyResource{}
}

type FunctionCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *FunctionCompliancePolicyResource) GetSchema(ctx context.Context) schema.Schema {
    return policy.GetPolicySchema(
        ctx,
        policyAPI.PolicyTypeComplianceFunction, 
        policyAPI.PolicyTypeComplianceFunctionFormatted, 
        policyAPI.PolicyContextComplianceFunction,
        policyAPI.TypeCompliance,
    )
}
