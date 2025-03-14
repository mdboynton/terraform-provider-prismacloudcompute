package planmodifiers

import (
	"context"

    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/path"
)

func SetEmptyIfUnknown() planmodifier.List {
    return setEmptyIfUnknown{} 
}

type setEmptyIfUnknown struct {}

func (m setEmptyIfUnknown) Description(_ context.Context) string {
    return ""
}

func (m setEmptyIfUnknown) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setEmptyIfUnknown) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    var rules basetypes.ListValue
    diags := req.Plan.GetAttribute(ctx, path.Root("rules"), &rules)
    if diags.HasError() {
        return
    }

    if rules.IsUnknown() {
        resp.PlanValue = types.ListValueMust(rules.ElementType(ctx), []attr.Value{})
    }

    return
}
