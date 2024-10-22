package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type _RolePermissionResourceModel struct {
    Name        types.String    `tfsdk:"name"`
    ReadWrite   types.Bool      `tfsdk:"read_write"`
}

func UseEmptySetForUnknownRolePermissions() planmodifier.Set {
    return useEmptySetForUnknownRolePermissions{} 
}

type useEmptySetForUnknownRolePermissions struct {}

func (m useEmptySetForUnknownRolePermissions) Description(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknownRolePermissions) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknownRolePermissions) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    //fmt.Println(")()()()()()()()()()()(")
    //fmt.Println("entering UseEmptySetForUnknownRolePermissions")
    ////fmt.Println("planvalue:")
    ////fmt.Println(req.PlanValue)
    ////fmt.Println("configvalue:")
    ////fmt.Println(req.ConfigValue)
    //fmt.Println(resp.PlanValue)
    //fmt.Println(")()()()()()()()()()()(")
    //if req.PlanValue.IsNull() {
    //    fmt.Println(")()()()()()()()()()()(")
    //    fmt.Println("role plan value is null")
    //    fmt.Println(")()()()()()()()()()()(")
    //}
    //var plan []_RolePermissionResourceModel 
    //diags := req.PlanValue.ElementsAs(ctx, &plan, false)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //for _, planValue := range plan {
    //    fmt.Println(planValue)
    //    //if planValue.Name.IsNull(){
    //    //    fmt.Println(")()()()()()()()()()()(")
    //    //    fmt.Println("role plan value is null")
    //    //    fmt.Println(")()()()()()()()()()()(")
    //    //}
    //}

    
    if req.ConfigValue.IsNull() {
        resp.PlanValue = types.SetValueMust(
            types.ObjectType{
                AttrTypes: map[string]attr.Type{
                    "name": types.StringType,
                    "read_write": types.BoolType,
                },
            },
            []attr.Value{},
        )
        
        return
    }

    return
}
