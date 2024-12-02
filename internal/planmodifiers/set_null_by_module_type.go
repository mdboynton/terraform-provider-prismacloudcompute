package planmodifiers

import (
	"context"

    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/path"
)

func SetNullByModuleType(moduleType string) planmodifier.List {
    return setNullByModuleType{
        ModuleType: moduleType,
    } 
}

type setNullByModuleType struct {
    ModuleType string
}

func (m setNullByModuleType) Description(_ context.Context) string {
    return ""
}

func (m setNullByModuleType) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setNullByModuleType) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    var rules basetypes.ListValue
    diags := req.Plan.GetAttribute(ctx, path.Root("rules"), &rules)
    if diags.HasError() {
        return
    }

    //var order basetypes.Int32Value
    //var (
    //    cveRules basetypes.SetValue
    //)

    for index := range rules.Elements() {
        if m.ModuleType == "compliance" {
            cveRulesAttrTypes := map[string]attr.Type{
                "name": types.StringType,
                "effect": types.StringType,
                "id": types.StringType,
                "description": types.StringType,
                "type": types.StringType,
                "expiration": types.ObjectType{
                    AttrTypes: map[string]attr.Type{
                        "enabled": types.BoolType,
                        "date": types.StringType,
                    },
                },
            }
            cveRulesNull := basetypes.NewSetNull(types.ObjectType{
                AttrTypes: cveRulesAttrTypes,
            })
            diags = req.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("cve_rules"), &cveRulesNull)
            if diags.HasError() {
                return
            }
        } else if m.ModuleType == "vulnerability" {
            var (
                cveRules basetypes.SetValue
            )

            diags = req.Plan.GetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("cve_rules"), &cveRules)
            if diags.HasError() {
                return
            }

            if cveRules.IsUnknown() {
                defaultCveRules := types.SetValueMust(
                    types.ObjectType{
                        AttrTypes: map[string]attr.Type{
                            "name": types.StringType,
                            "effect": types.StringType,
                            "id": types.StringType,
                            "description": types.StringType,
                            "type": types.StringType,
                            "expiration": types.ObjectType{
                                AttrTypes: map[string]attr.Type{
                                    "enabled": types.BoolType,
                                    "date": types.StringType,
                                },
                            },
                        },
                    },
                    []attr.Value{},
                )

                diags = req.Plan.SetAttribute(ctx, path.Root("rules").AtListIndex(index).AtName("cve_rules"), &defaultCveRules)
                if diags.HasError() {
                    return
                }
            }
        }
    }

    diags = req.Plan.GetAttribute(ctx, path.Root("rules"), &resp.PlanValue)

    return
}
