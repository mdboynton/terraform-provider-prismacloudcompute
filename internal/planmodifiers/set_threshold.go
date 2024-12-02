package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SetThresholdByModuleType(moduleType string) planmodifier.Object{
    return setThresholdByModuleType{
        ModuleType: moduleType,
    } 
}

type setThresholdByModuleType struct {
    ModuleType string
}

func (m setThresholdByModuleType) Description(_ context.Context) string {
    return ""
}

func (m setThresholdByModuleType) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setThresholdByModuleType) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    if m.ModuleType == "vulnerability" {
        if req.ConfigValue.IsNull() {
            resp.PlanValue = types.ObjectValueMust(
                map[string]attr.Type{
                    "threshold": types.StringType,
                    "risk_factors": types.SetType{ ElemType: types.StringType },
                },
                map[string]attr.Value{
                    "threshold": types.StringValue("off"),
                    "risk_factors": types.SetValueMust(
                        types.StringType,
                        []attr.Value{},
                    ),
                },
            )
        } else {
            resp.PlanValue = req.ConfigValue
        }
        return
    }

    if m.ModuleType == "compliance" {
        attrTypes := map[string]attr.Type{
            "threshold": types.StringType,
            "risk_factors": types.SetType{ ElemType: types.StringType },
        }
        resp.PlanValue = basetypes.NewObjectNull(attrTypes)
        return
    }

    return
}
