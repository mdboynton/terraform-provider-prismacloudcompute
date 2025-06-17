package planmodifiers

import (
	"context"
    "log"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func SetMaxScannersIfUsingHubAccount() planmodifier.Object{
    return setMaxScannersIfUsingHubAccount{} 
}

type setMaxScannersIfUsingHubAccount struct {}

func (m setMaxScannersIfUsingHubAccount) Description(_ context.Context) string {
    return ""
}

func (m setMaxScannersIfUsingHubAccount) MarkdownDescription(_ context.Context) string {
    return ""
}

func (m setMaxScannersIfUsingHubAccount) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
    log.Println("running SetMaxScannersIfUsingHubAccount")

    var hubAccountId basetypes.StringValue
    resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("agentless_scanning").AtName("hub_account_id"), &hubAccountId)...)
    if resp.Diagnostics.HasError() {
        return
    }
    log.Println("got hub account id")

    if (!hubAccountId.IsNull() && !hubAccountId.IsUnknown() && hubAccountId.ValueString() != "") {
        log.Println("modifying plan")
        resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("agentless_scanning").AtName("max_number_of_scanners"), types.Int32Value(0))...)
	    resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("agentless_scanning"), &resp.PlanValue)...)
    } else {
        log.Println("NO MODIFICATION")
    }

    return
}
