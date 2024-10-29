package planmodifiers

import (
	"context"
    //"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseEmptySetForUnknown() planmodifier.Set{
    return useEmptySetForUnknown{} 
}

type useEmptySetForUnknown struct {}

func (m useEmptySetForUnknown) Description(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknown) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useEmptySetForUnknown) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
    //if req.PlanValue.IsUnknown() {
    //    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    //    fmt.Println("planvalue is unknown")
    //    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    //    return
    //}
    //if req.PlanValue.IsNull() {
    //    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    //    fmt.Println("planvalue is null")
    //    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    //    return
    //}
    //    
    //fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    //fmt.Println("in set_collection_sets")
    //fmt.Println(req.PlanValue)
    //fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")

        
    if req.PlanValue.IsUnknown() {
        defaultSet := types.SetValueMust(
            types.StringType,
            []attr.Value{
                //types.StringValue("*"),
            },
        )
    
        resp.PlanValue = defaultSet
        return
    }

    return
}
