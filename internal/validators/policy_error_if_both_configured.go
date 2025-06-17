package validators

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	//"github.com/hashicorp/terraform-plugin-framework/path"
)

type errorIfBothListsConfigured struct {
	Attribute1Name string
	Attribute2Name string
}

func ErrorIfBothListsConfigured(attribute1Name string, attribute2Name string) errorIfBothListsConfigured {
	return errorIfBothListsConfigured{
		Attribute1Name: attribute1Name,
		Attribute2Name: attribute2Name,
	}
}

func (v errorIfBothListsConfigured) Description(ctx context.Context) string {
	return ""
}

func (v errorIfBothListsConfigured) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (v errorIfBothListsConfigured) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	util.HCLogDebug(ctx, "Executing ErrorIfBothListsConfigured")

	var (
		attribute1Value basetypes.ListValue
		attribute2Value basetypes.ListValue
	)

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName(v.Attribute1Name), &attribute1Value)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName(v.Attribute2Name), &attribute2Value)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !attribute1Value.IsNull() && !attribute1Value.IsUnknown() && !attribute2Value.IsNull() && !attribute2Value.IsUnknown() {
		resp.Diagnostics.AddError(
			"Invalid Resource Configuration",
			fmt.Sprintf("Invalid configuration at %s \nThe following attributes cannot be defined in the same policy rule: %s, %s", req.Path.String(), v.Attribute1Name, v.Attribute2Name),
		)
	}

	return
}
