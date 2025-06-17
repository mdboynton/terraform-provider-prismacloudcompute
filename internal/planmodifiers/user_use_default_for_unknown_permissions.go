package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownUserPermissions() planmodifier.Set {
	return useDefaultForUnknownUserPermissions{}
}

type useDefaultForUnknownUserPermissions struct{}

func (m useDefaultForUnknownUserPermissions) Description(_ context.Context) string {
	return ""
}

func (m useDefaultForUnknownUserPermissions) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m useDefaultForUnknownUserPermissions) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	var role string
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("role"), &role)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrMap := map[string]attr.Type{
		"project":     types.StringType,
		"collections": types.ListType{ElemType: types.StringType},
	}

	if role == "admin" || role == "operator" {
		resp.PlanValue = types.SetValueMust(
			types.ObjectType{
				AttrTypes: attrMap,
			},
			[]attr.Value{},
		)
		return
	}

	if req.PlanValue.IsUnknown() {
		if !req.StateValue.IsUnknown() && !req.StateValue.IsNull() {
			resp.PlanValue = req.StateValue
			return
		}

		defaultValue := types.ObjectValueMust(attrMap, map[string]attr.Value{
			"project":     types.StringValue("Central Console"),
			"collections": types.ListNull(types.StringType),
		})

		resp.PlanValue = types.SetValueMust(
			types.ObjectType{
				AttrTypes: attrMap,
			},
			[]attr.Value{
				defaultValue,
			},
		)
	}

	return
}
