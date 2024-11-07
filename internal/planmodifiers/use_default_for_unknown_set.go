package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownSet() planmodifier.Set{
    return useDefaultForUnknownSet{} 
}

type useDefaultForUnknownSet struct {}

func (m useDefaultForUnknownSet) Description(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownSet) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownSet) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    if req.PlanValue.IsUnknown() && !req.ConfigValue.IsNull() {
        //resp.PlanValue = types.SetValueMust(types.StringType, []attr.Value{})
        resp.PlanValue = req.ConfigValue
    }

    return
}
