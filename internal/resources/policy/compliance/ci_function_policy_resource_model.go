package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &CiFunctionCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &CiFunctionCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &CiFunctionCompliancePolicyResource{}

func NewCiFunctionCompliancePolicyResource() resource.Resource {
    return &CiFunctionCompliancePolicyResource{}
}

type CiFunctionCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *CiFunctionCompliancePolicyResource) GetSchema(ctx context.Context) schema.Schema {
    return policy.GetPolicySchema(
        ctx,
        policyAPI.PolicyTypeComplianceCiFunction, 
        policyAPI.PolicyTypeComplianceCiFunctionFormatted, 
        policyAPI.PolicyContextCiFunction,
        policyAPI.TypeCompliance,
    )
}
