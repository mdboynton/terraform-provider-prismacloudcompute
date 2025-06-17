package validators

import (
    "fmt"
    "context"
    "strings"
    "slices"

    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func LogInspectionPathIsValid() logInspectionPathIsValid {
    return logInspectionPathIsValid{}
}

type logInspectionPathIsValid struct {}

func (v logInspectionPathIsValid) Description(ctx context.Context) string {
    return ""
}

func (v logInspectionPathIsValid) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v logInspectionPathIsValid) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
    if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
        return
    }

    logInspectionRules := []models.RuntimeHostPolicyLogInspectionRuleResourceModel{}
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &logInspectionRules, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    pathValues := []string{}
    for idx, logInspectionRule := range logInspectionRules {
        pathValue := logInspectionRule.Path.ValueString()

        if len(pathValue) == 0 && pathValue != "/" {
            resp.Diagnostics.AddAttributeError(
                req.Path, 
                "Invalid Argument Value", 
                fmt.Sprintf("log_inspection_rule[%d] contains invalid file path. File path must have a valid, non-empty path to log file.", idx), 
            )
            continue
        }

        if !strings.HasPrefix(pathValue, "/") {
            resp.Diagnostics.AddAttributeError(
                req.Path, 
                "Invalid Argument Value", 
                fmt.Sprintf("log_inspection_rule[%d] contains invalid file path. File path must specify the absolute path to the log file (must begin with \"/\").", idx), 
            )
            continue
        }

        if strings.Contains(pathValue, "*") {
            splitPath := strings.Split(pathValue, "/")
            
            for splitIdx, value := range splitPath {
                if (strings.Contains(value, "*") && splitIdx != (len(splitPath) - 1)) {
                    resp.Diagnostics.AddAttributeError(
                        req.Path, 
                        "Invalid Argument Value", 
                        fmt.Sprintf("log_inspection_rule[%d] contains invalid file path. File path can only use asterisks in the file name portion.", idx),
                    )
                    continue
                }
            }
        }

        if slices.Contains(pathValues, pathValue) {
            resp.Diagnostics.AddAttributeError(
                req.Path, 
                "Invalid Argument Value", 
                fmt.Sprintf("log_inspection_rule[%d] contains duplicate file path. File paths must be unique for each rule.", idx), 
            )
        } else {
            pathValues = append(pathValues, pathValue)
        }
    }

    return
}

