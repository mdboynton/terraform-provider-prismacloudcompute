package planmodifiers

import (
	"context"
	"fmt"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseAllForDefaultCollectionName() planmodifier.String {
    return &useAllForDefaultCollectionName{}
}

type useAllForDefaultCollectionName struct {}

func (m *useAllForDefaultCollectionName) Description(_ context.Context) string {
    return ""
}

func (m *useAllForDefaultCollectionName) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useAllForDefaultCollectionName) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.StringValue("All")
        return
    }

    return
}
