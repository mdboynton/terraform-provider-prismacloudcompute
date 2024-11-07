package planmodifiers

import (
	"context"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseEmptyStringForUnknown() planmodifier.String {
    return &useEmptyStringForUnknown{}
}

type useEmptyStringForUnknown struct {}

func (m *useEmptyStringForUnknown) Description(_ context.Context) string {
    return ""
}

func (m *useEmptyStringForUnknown) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useEmptyStringForUnknown) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.StringValue("")
    }

    return
}
