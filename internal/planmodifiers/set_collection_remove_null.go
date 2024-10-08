package planmodifiers

import (
	"context"
    "fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CollectionObjectType() types.ObjectType {
    return types.ObjectType{
        AttrTypes: CollectionObjectAttrTypeMap(),
    }
}

func CollectionObjectAttrTypeMap() map[string]attr.Type {
    return map[string]attr.Type{
        "account_ids":  types.SetType{ElemType: types.StringType},
        "app_ids":  types.SetType{ElemType: types.StringType},
        "clusters":  types.SetType{ElemType: types.StringType},
        "color":        types.StringType,
        "containers":  types.SetType{ElemType: types.StringType},
        "description":  types.StringType,
        "functions":  types.SetType{ElemType: types.StringType},
        "hosts":  types.SetType{ElemType: types.StringType},
        "images":  types.SetType{ElemType: types.StringType},
        "labels":  types.SetType{ElemType: types.StringType},
        "modified":         types.StringType,
        "name":         types.StringType,
        "namespaces":  types.SetType{ElemType: types.StringType},
        "owner":        types.StringType,
        "prisma":       types.BoolType,
        "system":       types.BoolType,
    }
}

func RemoveNullObjects() planmodifier.List {
    return removeNullObjects{} 
}

type removeNullObjects struct {}

func (m removeNullObjects) Description(_ context.Context) string {
    return ""
}

func (m removeNullObjects) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m removeNullObjects) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    //fmt.Println("!!!!!!!!!!!!!!!!!!!!")
    //fmt.Println("entering RemoveNullObjects")
    //fmt.Println("planValue:")
    //fmt.Println(req.PlanValue)
    //fmt.Println("stateValue:")
    //fmt.Println(req.StateValue)
    //fmt.Println("!!!!!!!!!!!!!!!!!!!!")
    
    //if !req.PlanValue.IsNull() && 
    //if req.PlanValue != req.StateValue {
    //    
    //}

    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    fmt.Println(req.PlanValue.Elements())
    fmt.Println(len(req.ConfigValue.Elements()))
    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")

    if len(req.ConfigValue.Elements()) == 0 {
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
        fmt.Println("RemoveNullObjects: setting PlanValue to NewListUnknown")
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")

        resp.PlanValue = basetypes.NewListUnknown(CollectionObjectType())
        return
    }

    return
}
