package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownCondition() planmodifier.Object{
    return useDefaultForUnknownCondition{} 
}

type useDefaultForUnknownCondition struct {}

func (m useDefaultForUnknownCondition) Description(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownCondition) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownCondition) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    if req.PlanValue.IsUnknown() && !req.ConfigValue.IsNull() {
        //resp.PlanValue = types.SetValueMust(types.StringType, []attr.Value{})
        resp.PlanValue = req.ConfigValue
    }

    return
}
