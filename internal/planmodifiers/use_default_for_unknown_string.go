
package planmodifiers

import (
	"context"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownString() planmodifier.String {
    return &useDefaultForUnknownString{}
}

type useDefaultForUnknownString struct {}

func (m *useDefaultForUnknownString) Description(_ context.Context) string {
    return ""
}

func (m *useDefaultForUnknownString) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useDefaultForUnknownString) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    //if req.PlanValue.IsUnknown() {
    //    resp.PlanValue = types.StringValue("All")
    //    return
    //}

    return
}
