package planmodifiers

import (
	"context"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseAlertForUnknownEffect() planmodifier.String {
    return &useAlertForUnknownEffect{}
}

type useAlertForUnknownEffect struct {}

func (m *useAlertForUnknownEffect) Description(_ context.Context) string {
    return ""
}

func (m *useAlertForUnknownEffect) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useAlertForUnknownEffect) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.PlanValue.IsUnknown() {
        //resp.PlanValue = types.StringValue("unknown")
        resp.PlanValue = types.StringValue("default")
        return
    }
    return
}
