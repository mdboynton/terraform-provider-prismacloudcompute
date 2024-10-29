package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownTrustedImageRules() planmodifier.Set {
    return useDefaultForUnknownTrustedImageRules{} 
}

type useDefaultForUnknownTrustedImageRules struct {}

func (m useDefaultForUnknownTrustedImageRules) Description(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownTrustedImageRules) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownTrustedImageRules) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    if req.PlanValue.IsUnknown() {
        defaultSet := types.SetValueMust(
            RulesObjectType(),
            []attr.Value{
            },
        )
    
        resp.PlanValue = defaultSet
        return
    }

    return
}

// TODO: move these to a common location

func RulesObjectType() types.ObjectType {
    return types.ObjectType{
        AttrTypes: RulesObjectAttrTypeMap(),
    }
}

func RulesObjectAttrTypeMap() map[string]attr.Type {
    return map[string]attr.Type{
        "action":  types.SetType{ElemType: types.StringType},
        "allowed_groups":  types.SetType{ElemType: types.StringType},
        "collections": types.ListType{ElemType: collectionObjectType()},
        "denied_groups":  types.SetType{ElemType: types.StringType},
        "effect":        types.StringType,
        "modified": types.StringType,
        "name":         types.StringType,
        "notes":         types.StringType,
        "owner":        types.StringType,
        "previous_name":       types.StringType,
    }
}

func collectionObjectType() types.ObjectType {
    return types.ObjectType{
        AttrTypes: collectionObjectAttrTypeMap(),
    }
}

func collectionObjectAttrTypeMap() map[string]attr.Type {
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
        "modified": types.StringType,
        "name":         types.StringType,
        "namespaces":  types.SetType{ElemType: types.StringType},
        "owner":        types.StringType,
        "prisma":       types.BoolType,
        "system":       types.BoolType,
    }
}
