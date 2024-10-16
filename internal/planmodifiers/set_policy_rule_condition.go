package planmodifiers

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func conditionObjectTypeMap() map[string]attr.Type {
    return map[string]attr.Type{
        "vulnerabilities": types.ListType{
            ElemType: conditionVulnerabilitiesObjectType(),
        },
    }
}

func conditionVulnerabilitiesObjectType() types.ObjectType {
    return types.ObjectType{
        AttrTypes: conditionVulnerabilitiesObjectAttrTypeMap(),
    }
}

func conditionVulnerabilitiesObjectAttrTypeMap() map[string]attr.Type {
    return map[string]attr.Type{
        "id": types.Int32Type,
        "block": types.BoolType,
    }
}

func AllowUnknownCondition() planmodifier.Object{
    return allowUnknownCondition{} 
}

type allowUnknownCondition struct {}

func (m allowUnknownCondition) Description(_ context.Context) string {
    return ""
}

func (m allowUnknownCondition) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m allowUnknownCondition) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    if req.ConfigValue.IsNull() {
        if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
            resp.PlanValue = basetypes.NewObjectUnknown(conditionObjectTypeMap())
        }
        return
    }

    return
}
