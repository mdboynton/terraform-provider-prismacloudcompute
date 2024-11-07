package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseStateForUnknownRolePermissionName() planmodifier.String {
    return &useStateForUnknownRolePermissionName{}
}

type useStateForUnknownRolePermissionName struct {}

func (m *useStateForUnknownRolePermissionName) Description(_ context.Context) string {
    return ""
}

func (m *useStateForUnknownRolePermissionName) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useStateForUnknownRolePermissionName) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    return
}


func UseStateForUnknownRolePermissionReadWrite() planmodifier.Bool {
    return &useStateForUnknownRolePermissionReadWrite{}
}

type useStateForUnknownRolePermissionReadWrite struct {}

func (m *useStateForUnknownRolePermissionReadWrite) Description(_ context.Context) string {
    return ""
}

func (m *useStateForUnknownRolePermissionReadWrite) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useStateForUnknownRolePermissionReadWrite) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    return
}

