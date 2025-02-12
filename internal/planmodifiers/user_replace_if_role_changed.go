package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ReplaceIfRoleNameChanged() planmodifier.String {
    return stringplanmodifier.RequiresReplaceIf(
        func(
            ctx context.Context,
            sr planmodifier.StringRequest,
            rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse,
        ) {
            rrifr.RequiresReplace = (sr.PlanValue.ValueString() != sr.StateValue.ValueString()) && (sr.PlanValue.ValueString() == "admin" || sr.PlanValue.ValueString() == "operator")
        },
        "TODO",
        "TODO",
    )
}
