package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

//func UseStateForNullRolePermissions() planmodifier.Set {
//    return useStateForNullRolePermissions{} 
//}
//
//type useStateForNullRolePermissions struct {}
//
//func (m useStateForNullRolePermissions) Description(_ context.Context) string {
//    return ""
//}
//
//func (m useStateForNullRolePermissions) MarkdownDescription(_ context.Context) string {
//    return ""
//}
//
//func (m useStateForNullRolePermissions) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
//    //if req.ConfigValue.IsNull() {
//    //    resp.PlanValue = types.SetValueMust(
//    //        types.ObjectType{
//    //            AttrTypes: map[string]attr.Type{
//    //                "name": types.StringType,
//    //                "read_write": types.BoolType,
//    //            },
//    //        },
//    //        []attr.Value{},
//    //    )
//    //    
//    //    return
//    //}
//    fmt.Println(")()()()()()()()()()()(")
//    fmt.Println("entering UseStateForNullRolePermissions")
//    fmt.Println(")()()()()()()()()()()(")
//
//    if req.PlanValue.IsNull() {
//        fmt.Println(")()()()()()()()()()()(")
//        fmt.Println("role plan value is null")
//        fmt.Println(")()()()()()()()()()()(")
//    }
//    
//    if req.PlanValue.IsUnknown() {
//        fmt.Println(")()()()()()()()()()()(")
//        fmt.Println("role plan value is Unknown")
//        fmt.Println(")()()()()()()()()()()(")
//    }
//
//    if !req.ConfigValue.IsUnknown() && req.ConfigValue.IsNull() {
//        fmt.Println(")()()()()()()()()()()(")
//        fmt.Println("role config value is known and not null")
//        fmt.Println(")()()()()()()()()()()(")
//    }
//
//    return
//}

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
    //fmt.Println(")()()()()()()()()()()(")
    //fmt.Println("entering UseStateForUnknownRolePermissionName")
    //fmt.Println("planvalue:")
    //fmt.Println(req.PlanValue)
    //fmt.Println("configvalue:")
    //fmt.Println(req.ConfigValue)
    //fmt.Println(")()()()()()()()()()()(")

    ////if req.PlanValue.IsUnknown() && (!req.StateValue.IsUnknown() && !req.StateValue.IsNull()) {
    ////if req.PlanValue.IsNull() && (!req.StateValue.IsUnknown() && !req.StateValue.IsNull()) {
    ////    resp.PlanValue = req.StateValue
    ////    return
    ////}

    //if req.PlanValue.IsNull() {
    //    fmt.Println(")()()()()()()()()()()(")
    //    fmt.Println("role plan value is null")
    //    fmt.Println(")()()()()()()()()()()(")
    //}
    //
    //if req.PlanValue.IsUnknown() {
    //    fmt.Println(")()()()()()()()()()()(")
    //    fmt.Println("role plan value is Unknown")
    //    fmt.Println(")()()()()()()()()()()(")
    //}

    //if !req.ConfigValue.IsUnknown() && !req.ConfigValue.IsNull() {
    //    fmt.Println(")()()()()()()()()()()(")
    //    fmt.Println("role config value is known and not null")
    //    fmt.Println(")()()()()()()()()()()(")
    //}
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
    //if req.PlanValue.IsNull() {
    //    resp.PlanValue = nil
    //    return
    //}


    //fmt.Println(")()()()()()()()()()()(")
    //fmt.Println("entering UseStateForUnknownRolePermissionReadWrite")
    //fmt.Println("planvalue:")
    //fmt.Println(req.PlanValue)
    //fmt.Println("configvalue:")
    //fmt.Println(req.ConfigValue)
    //fmt.Println(")()()()()()()()()()()(")
    //if req.PlanValue.IsUnknown() {
    //    resp.PlanValue = types.BoolValue(true)
    //    return
    //}

    return
}

