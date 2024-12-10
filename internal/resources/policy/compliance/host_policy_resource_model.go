package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	
    "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &HostCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &HostCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &HostCompliancePolicyResource{}

func NewHostCompliancePolicyResource() resource.Resource {
    return &HostCompliancePolicyResource{}
}

type HostCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *HostCompliancePolicyResource) GetSchema(ctx context.Context) schema.Schema {
    return models.GetPolicySchema(
        ctx,
        policyAPI.PolicyTypeComplianceHost, 
        policyAPI.PolicyTypeComplianceHostFormatted, 
        policyAPI.PolicyContextHost,
        policyAPI.TypeCompliance,
    )
}
