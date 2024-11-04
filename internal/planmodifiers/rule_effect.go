package planmodifiers

import (
	"context"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownEffect() planmodifier.String {
    return &useDefaultForUnknownEffect{}
}

type useDefaultForUnknownEffect struct {}

func (m *useDefaultForUnknownEffect) Description(_ context.Context) string {
    return ""
}

func (m *useDefaultForUnknownEffect) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useDefaultForUnknownEffect) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.StringValue("default")
        return
    }
    return
}
