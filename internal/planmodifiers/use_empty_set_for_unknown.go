package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseEmptySetForUnknown() planmodifier.Set{
    return useEmptySetForUnknown{} 
}

type useEmptySetForUnknown struct {}

func (m useEmptySetForUnknown) Description(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknown) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknown) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.SetValueMust(types.StringType, []attr.Value{})
    }

    return
}
