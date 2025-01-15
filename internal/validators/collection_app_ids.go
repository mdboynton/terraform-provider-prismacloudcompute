package validators

import (
    "context"
    "strings"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func AppIDsEndWithWildcard() appIDsEndWithWildcard {
    return appIDsEndWithWildcard{}
}

type appIDsEndWithWildcard struct {}

func (v appIDsEndWithWildcard) Description(ctx context.Context) string {
    return ""
}

func (v appIDsEndWithWildcard) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v appIDsEndWithWildcard) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
    if (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
        return
    }

    var appIDs []string
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &appIDs, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    for _, val := range appIDs {
        if !strings.HasSuffix(val, "*") {
            resp.Diagnostics.AddError(
                "Invalid Resource Configuration",
                "All App ID resources must end with a wildcard (\"*\")",
            )
        }
    }

    return
}
