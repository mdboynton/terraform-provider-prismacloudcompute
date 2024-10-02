package planmodifiers

import (
	"context"
    //"time"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseCurrentTimeForDefaultCollectionModified() planmodifier.String {
    return &useCurrentTimeForDefaultCollectionModified{}
}

type useCurrentTimeForDefaultCollectionModified struct {}

func (m *useCurrentTimeForDefaultCollectionModified) Description(_ context.Context) string {
    return ""
}

func (m *useCurrentTimeForDefaultCollectionModified) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useCurrentTimeForDefaultCollectionModified) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
        if !req.StateValue.IsUnknown() && !req.StateValue.IsNull() {
            resp.PlanValue = req.StateValue
            return
        }

        return
}
