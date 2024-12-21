package policy

import (
    "context"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"
	
    "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &ContainerCompliancePolicyResource{}
var _ resource.ResourceWithImportState = &ContainerCompliancePolicyResource{}
var _ resource.ResourceWithModifyPlan = &ContainerCompliancePolicyResource{}

func NewContainerCompliancePolicyResource() resource.Resource {
    return &ContainerCompliancePolicyResource{}
}

type ContainerCompliancePolicyResource struct {
    client *api.PrismaCloudComputeAPIClient
}

func (r *ContainerCompliancePolicyResource) GetSchema(ctx context.Context) schema.Schema {
    util.DLog(ctx, "starting GetSchema for container compliance policy")
    if r.client == nil {
        util.DLog(ctx, "client is nil in GetSchema")
    }

    policySchema, _ := models.GetPolicySchema(
        ctx,
        policyAPI.PolicyTypeComplianceContainer, 
        policyAPI.PolicyTypeComplianceContainerFormatted, 
        policyAPI.PolicyContextContainer,
        policyAPI.TypeCompliance,
    )

    return policySchema
}
