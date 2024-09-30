package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownCollectionSets() planmodifier.Set {
    return useDefaultForUnknownCollectionSets{} 
}

type useDefaultForUnknownCollectionSets struct {}

func (m useDefaultForUnknownCollectionSets) Description(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownCollectionSets) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownCollectionSets) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    if req.PlanValue.IsUnknown() {
        defaultSet := types.SetValueMust(
            types.StringType,
            []attr.Value{
                types.StringValue("*"),
            },
        )
    
        resp.PlanValue = defaultSet
        return
    }

    return
}
