package planmodifiers

import (
	"context"
    "fmt"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownThreshold() planmodifier.Object {
    return &useDefaultForUnknownThreshold{}
}

type useDefaultForUnknownThreshold struct {}

func (m *useDefaultForUnknownThreshold) Description(_ context.Context) string {
    return ""
}

func (m *useDefaultForUnknownThreshold) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useDefaultForUnknownThreshold) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    fmt.Println("")

    if req.PlanValue.IsUnknown() {
        return
    }

    return
}
