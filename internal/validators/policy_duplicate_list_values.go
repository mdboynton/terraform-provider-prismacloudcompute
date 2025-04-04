package validators

import (
	"context"
	"fmt"
	"slices"
    "strings"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type noSharedValuesBetweenProcessLists struct {
}

func NoSharedValuesBetweenProcessLists() noSharedValuesBetweenProcessLists {
    return noSharedValuesBetweenProcessLists{}
}

func (v noSharedValuesBetweenProcessLists) Description(ctx context.Context) string {
    return ""
}

func (v noSharedValuesBetweenProcessLists) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v noSharedValuesBetweenProcessLists) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
    util.DLog(ctx, "Executing NoSharedValuesBetweenProcessLists")

    var (
        allowedProcesses basetypes.ListValue 
        deniedProcesses basetypes.ListValue 
        sharedProcessNames []string = []string{}
    )

    resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName("allowed_processes"), &allowedProcesses)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.AtName("denied_processes").AtName("paths"), &deniedProcesses)...)
    if resp.Diagnostics.HasError() {
        return
    }

    if (len(allowedProcesses.Elements()) == 0 && len(deniedProcesses.Elements()) == 0) {
        return
    }

    allowedProcessNames := make([]string, 0, len(allowedProcesses.Elements()))
    resp.Diagnostics.Append(allowedProcesses.ElementsAs(ctx, &allowedProcessNames, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    deniedProcessNames := make([]string, 0, len(deniedProcesses.Elements()))
    resp.Diagnostics.Append(deniedProcesses.ElementsAs(ctx, &deniedProcessNames, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    for _, allowedProcessName := range allowedProcessNames {
        if slices.Contains(deniedProcessNames, allowedProcessName) {
            sharedProcessNames = append(sharedProcessNames, allowedProcessName)
        }
    }

    if len(sharedProcessNames) > 0 {
        resp.Diagnostics.AddError(
			"Invalid Resource Configuration",
            fmt.Sprintf("Invalid configuration at %s\nValues in allowed_processes and denied_processes.paths cannot be shared.\nThe following values must be removed from one of either argument: %s", req.Path.String(), strings.Join(sharedProcessNames, ", ")),
        )
    }

    return 
}
