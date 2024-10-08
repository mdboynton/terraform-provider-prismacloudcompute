package planmodifiers

import (
	"context"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseTrueForDefaultCollectionBools() planmodifier.Bool {
    return &useTrueForDefaultCollectionBools{}
}

type useTrueForDefaultCollectionBools struct {}

func (m *useTrueForDefaultCollectionBools) Description(_ context.Context) string {
    return ""
}

func (m *useTrueForDefaultCollectionBools) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useTrueForDefaultCollectionBools) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    //if req.PlanValue.IsNull() {
    //    resp.PlanValue = nil
    //    return
    //}


    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.BoolValue(true)
        return
    }

    return
}

func UseFalseForDefaultCollectionBools() planmodifier.Bool {
    return &useFalseForDefaultCollectionBools{}
}

type useFalseForDefaultCollectionBools struct {}

func (m *useFalseForDefaultCollectionBools) Description(_ context.Context) string {
    return ""
}

func (m *useFalseForDefaultCollectionBools) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useFalseForDefaultCollectionBools) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    if req.PlanValue.IsUnknown() {
        resp.PlanValue = types.BoolValue(false)
        return
    }

    return
}
