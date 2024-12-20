package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &CiImageCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &CiImageCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &CiImageCompliancePolicyResource{}

func NewCiImageCompliancePolicyResource() resource.Resource {
    return &CiImageCompliancePolicyResource{}
}

type CiImageCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *CiImageCompliancePolicyResource) GetSchema(ctx context.Context) schema.Schema {
    policySchema, _ := policy.GetPolicySchema(
        ctx,
        policyAPI.PolicyTypeComplianceCiImage, 
        policyAPI.PolicyTypeComplianceCiImageFormatted, 
        policyAPI.PolicyContextCiImage,
        policyAPI.TypeCompliance,
    )

    return policySchema
}
