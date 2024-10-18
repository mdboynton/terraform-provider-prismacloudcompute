package planmodifiers

import (
	"context"
    //"fmt"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseStateRuleIdForUnknownIfExists() planmodifier.Int32 {
    return &useStateRuleIdForUnknownIfExists{}
}

type useStateRuleIdForUnknownIfExists struct {}

func (m *useStateRuleIdForUnknownIfExists) Description(_ context.Context) string {
    return ""
}

func (m *useStateRuleIdForUnknownIfExists) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useStateRuleIdForUnknownIfExists) PlanModifyInt32(_ context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
    if req.PlanValue.IsUnknown() && (!req.StateValue.IsUnknown() && !req.StateValue.IsNull()) {
        resp.PlanValue = req.StateValue
        return
    }

    return
}
