package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseEmptyListForUnknownConditionVulnerabilities() planmodifier.List{
    return useEmptyListForUnknownConditionVulnerabilities{} 
}

type useEmptyListForUnknownConditionVulnerabilities struct {}

func (m useEmptyListForUnknownConditionVulnerabilities) Description(_ context.Context) string {
    return ""
}

func (m useEmptyListForUnknownConditionVulnerabilities) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useEmptyListForUnknownConditionVulnerabilities) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.ListValueMust(
            types.ObjectType{
                AttrTypes: map[string]attr.Type{
                    "id": types.Int32Type,
                    "block": types.BoolType,
                },
            }, 
            []attr.Value{},
        )
    }

    return
}
