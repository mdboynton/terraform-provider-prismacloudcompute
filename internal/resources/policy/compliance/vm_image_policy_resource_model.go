package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"

    "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &VmImageCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &VmImageCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &VmImageCompliancePolicyResource{}

func NewVmImageCompliancePolicyResource() resource.Resource {
    return &VmImageCompliancePolicyResource{}
}

type VmImageCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *VmImageCompliancePolicyResource) GetSchema(ctx context.Context) schema.Schema {
    return policy.GetPolicySchema(
        ctx,
        policyAPI.PolicyTypeComplianceVmImage, 
        policyAPI.PolicyTypeComplianceVmImageFormatted, 
        policyAPI.PolicyContextComplianceVmImage,
        policyAPI.TypeCompliance,
    )
}
