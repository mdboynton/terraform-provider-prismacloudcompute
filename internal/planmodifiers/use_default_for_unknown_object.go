package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownObject() planmodifier.Object{
    return useDefaultForUnknownObject{} 
}

type useDefaultForUnknownObject struct {}

func (m useDefaultForUnknownObject) Description(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownObject) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownObject) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    if req.PlanValue.IsUnknown() && !req.ConfigValue.IsNull() {
        //resp.PlanValue = types.SetValueMust(types.StringType, []attr.Value{})
        resp.PlanValue = req.ConfigValue
    }

    return
}
