package planmodifiers

import (
	"context"

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
    if len(req.ConfigValue.Elements()) == 0 {
        resp.PlanValue = basetypes.NewListUnknown(CollectionObjectType())
        return
    }

    return
}
