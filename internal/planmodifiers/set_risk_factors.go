package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseEmptySetForUnknownRiskFactors() planmodifier.Set{
    return useEmptySetForUnknownRiskFactors{} 
}

type useEmptySetForUnknownRiskFactors struct {}

func (m useEmptySetForUnknownRiskFactors) Description(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknownRiskFactors) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknownRiskFactors) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.SetValueMust(types.StringType, []attr.Value{})
    }

    return
}
