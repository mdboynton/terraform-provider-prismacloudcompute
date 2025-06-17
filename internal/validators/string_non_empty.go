package validators

import (
    "context"
    "fmt"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type isNonEmpty struct {
    AttributeName string
}

func IsNonEmpty(attributeName string) isNonEmpty {
    return isNonEmpty{
        AttributeName: attributeName,
    }
}

func (v isNonEmpty) Description(ctx context.Context) string {
    return ""
}

func (v isNonEmpty) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v isNonEmpty) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
    if (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
        return
    }

    if req.ConfigValue.ValueString() == "" {
        resp.Diagnostics.AddAttributeError(
            req.Path,
            "Invalid Resource Configuration",
            fmt.Sprintf("%s must be a non-empty value.", v.AttributeName),
        )
    }

    return
}
