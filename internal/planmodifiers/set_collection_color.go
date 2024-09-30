package planmodifiers

import (
	"context"
	"fmt"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultColorForDefaultCollectionColor() planmodifier.String {
    return &useDefaultColorForDefaultCollectionColor{}
}

type useDefaultColorForDefaultCollectionColor struct {}

func (m *useDefaultColorForDefaultCollectionColor) Description(_ context.Context) string {
    return ""
}

func (m *useDefaultColorForDefaultCollectionColor) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useDefaultColorForDefaultCollectionColor) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.StringValue("#3FA2F7")
        return
    }
    return
}
