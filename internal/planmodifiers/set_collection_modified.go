package planmodifiers

import (
	"context"
	"fmt"
    //"time"

	//"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/types"
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
    //if req.PlanValue.IsUnknown() {
    //    fmt.Println("setting resp.PlanValue to empty string")
    //    resp.PlanValue = types.StringValue("")
    //}

    //if resp.PlanValue.IsUnknown() {
    //    //req.PlanValue = req.ConfigValue
    //    resp.PlanValue = req.ConfigValue
    //    return
    //}

    //if req.PlanValue.IsUnknown() {
    //    resp.PlanValue = types.StringValue("")
    //    return
    //}

    //if resp.PlanValue.IsUnknown() {
    //    return
    //}

    //resp.PlanValue = req.StateValue


    //if req.PlanValue.IsUnknown() {
    //    resp.PlanValue = types.StringValue(time.Now().Format("2006-01-02T15:04:05.000Z"))
    //} else {
    //    fmt.Println("req planvalue:")
    //    fmt.Println(req.PlanValue)
    //    fmt.Println("req configvluae:")
    //    fmt.Println(req.ConfigValue)
    //    fmt.Println("assigning req planvalue")
    //    //resp.PlanValue = req.PlanValue
    //}
    //return
}
