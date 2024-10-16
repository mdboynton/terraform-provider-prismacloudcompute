package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func AllowUnknownEffect() planmodifier.String {
    return &allowUnknownEffect{}
}

type allowUnknownEffect struct {}

func (m *allowUnknownEffect) Description(_ context.Context) string {
    return ""
}

func (m *allowUnknownEffect) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *allowUnknownEffect) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.ConfigValue.IsNull() {
        if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
            resp.PlanValue = basetypes.NewStringUnknown()
            return
        }
    }

    return
}
