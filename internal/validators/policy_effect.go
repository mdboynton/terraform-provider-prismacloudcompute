package validators

import (
    "context"
    "fmt"
    "slices"
    "strings"
	
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func PolicyEffectIsValid(attributeName string, validEffects []string) policyEffectIsValid {
    return policyEffectIsValid{
        AttributeName: attributeName,
        ValidEffects: validEffects,
    }
}

type policyEffectIsValid struct {
    AttributeName string
    ValidEffects []string
}

func (v policyEffectIsValid) Description(ctx context.Context) string {
    return ""
}

func (v policyEffectIsValid) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v policyEffectIsValid) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
    if (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
        return
    }

    if !slices.Contains(v.ValidEffects, req.ConfigValue.ValueString()) {
        resp.Diagnostics.AddAttributeError(
            req.Path, 
            "Invalid Argument Value", 
            fmt.Sprintf("Invalid value \"%s\" provided for argument \"%s\". Must be one of the following: %s", req.ConfigValue.ValueString(), v.AttributeName, strings.Join(v.ValidEffects[:], ", ")),
        )
    }

    return
}
