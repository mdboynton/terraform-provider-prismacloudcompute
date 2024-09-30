package planmodifiers

import (
	"context"
	"fmt"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForDefaultCollectionDescription() planmodifier.String {
    return &useDefaultForDefaultCollectionDescription{}
}

type useDefaultForDefaultCollectionDescription struct {}

func (m *useDefaultForDefaultCollectionDescription) Description(_ context.Context) string {
    return ""
}

func (m *useDefaultForDefaultCollectionDescription) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useDefaultForDefaultCollectionDescription) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.StringValue("System - all resources collection")
        return
    }

    return
}
