package planmodifiers

import (
	"context"
    "fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/attr"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultForUnknownCollectionSets() planmodifier.List {
    return useDefaultForUnknownCollectionSets{} 
}

type useDefaultForUnknownCollectionSets struct {}

func (m useDefaultForUnknownCollectionSets) Description(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownCollectionSets) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m useDefaultForUnknownCollectionSets) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
    if req.PlanValue.IsUnknown() {
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
        fmt.Println("planvalue is unknown")
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
        return
    }
    if req.PlanValue.IsNull() {
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
        fmt.Println("planvalue is null")
        fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
        return
    }
        
    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")
    fmt.Println("in set_collection_sets")
    fmt.Println(req.PlanValue)
    fmt.Println("%%%%%%%%%%%%%%%%%%%%%%")

        
    //if req.PlanValue.IsUnknown() {
    //    defaultSet := types.SetValueMust(
    //        types.StringType,
    //        []attr.Value{
    //            types.StringValue("*"),
    //        },
    //    )
    //
    //    resp.PlanValue = defaultSet
    //    return
    //}

    return
}

//func UseDefaultForUnknownCollectionSets() planmodifier.Set {
//    return useDefaultForUnknownCollectionSets{} 
//}
//
//type useDefaultForUnknownCollectionSets struct {}
//
//func (m useDefaultForUnknownCollectionSets) Description(_ context.Context) string {
//    return ""
//}
//
//func (m useDefaultForUnknownCollectionSets) MarkdownDescription(_ context.Context) string {
//    return ""
//}
//
//func (m useDefaultForUnknownCollectionSets) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
//    if req.PlanValue.IsUnknown() {
//        defaultSet := types.SetValueMust(
//            types.StringType,
//            []attr.Value{
//                types.StringValue("*"),
//            },
//        )
//    
//        resp.PlanValue = defaultSet
//        return
//    }
//
//    return
//}
