package planmodifiers

import (
	"context"
    "fmt"
    //"time"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseCurrentTimeForDefaultCollectionModified() planmodifier.String {
    return &useCurrentTimeForDefaultCollectionModified{}
}

type useCurrentTimeForDefaultCollectionModified struct {}

func (m *useCurrentTimeForDefaultCollectionModified) Description(_ context.Context) string {
    return ""
}

func (m *useCurrentTimeForDefaultCollectionModified) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useCurrentTimeForDefaultCollectionModified) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
        if req.PlanValue.IsUnknown() && (!req.StateValue.IsUnknown() && !req.StateValue.IsNull()) {
            fmt.Println("%%%%%%%%%%%%%%%%%%%%")
            fmt.Println("modifying modified value")
            fmt.Println("%%%%%%%%%%%%%%%%%%%%")
            resp.PlanValue = req.StateValue
            return
        }

        return
}



func UseEmptyStringForNull() planmodifier.String {
    return &useEmptyStringForNull{}
}

type useEmptyStringForNull struct {}

func (m *useEmptyStringForNull) Description(_ context.Context) string {
    return ""
}

func (m *useEmptyStringForNull) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m *useEmptyStringForNull) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
        fmt.Println("%%%%%%%%%%%%%%%%%%%%")
        fmt.Println("entering UseEmptyStringForNull")
        fmt.Println("planValue:")
        fmt.Println(req.PlanValue)
        fmt.Println("%%%%%%%%%%%%%%%%%%%%")

        if req.PlanValue.IsUnknown() {
        //if req.PlanValue.IsNull() {
            fmt.Println("%%%%%%%%%%%%%%%%%%%%")
            fmt.Println("setting null modified to empty string")
            fmt.Println("%%%%%%%%%%%%%%%%%%%%")
            resp.PlanValue = types.StringValue("")
            return
        }

        return
}
