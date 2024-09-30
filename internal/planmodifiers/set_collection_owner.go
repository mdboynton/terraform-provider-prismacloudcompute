package planmodifiers

import (
	"context"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseSystemForDefaultCollectionOwner() planmodifier.String {
    return &useSystemForDefaultCollectionOwner{}
}

type useSystemForDefaultCollectionOwner struct {}

func (m *useSystemForDefaultCollectionOwner) Description(_ context.Context) string {
    return ""
}

func (m *useSystemForDefaultCollectionOwner) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useSystemForDefaultCollectionOwner) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.StringValue("system")
        return
    }

    return
}
