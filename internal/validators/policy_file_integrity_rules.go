package validators

import (
    "fmt"
    "context"
    "slices"

    "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FileIntegrityRulesAreValid() fileIntegrityRulesAreValid {
    return fileIntegrityRulesAreValid{}
}

type fileIntegrityRulesAreValid struct {}

func (v fileIntegrityRulesAreValid) Description(ctx context.Context) string {
    return ""
}

func (v fileIntegrityRulesAreValid) MarkdownDescription(ctx context.Context) string {
    return ""
}

func (v fileIntegrityRulesAreValid) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
    if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
        return
    }

    fileIntegrityRules := []models.RuntimeHostPolicyFileIntegrityRuleResourceModel{}
    resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &fileIntegrityRules, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    filePaths := []string{}
    for idx, fileIntegrityRules := range fileIntegrityRules {
        filePath := fileIntegrityRules.Path.ValueString()
        monitorSubdirectories := fileIntegrityRules.MonitorSubdirectories.ValueBool()
        monitorWriteOps := fileIntegrityRules.MonitorWriteOps.ValueBool()
        monitorReadOps := fileIntegrityRules.MonitorReadOps.ValueBool()
        monitorMetadata := fileIntegrityRules.MonitorMetadataChanges.ValueBool()

        if (!monitorSubdirectories && !monitorWriteOps && !monitorReadOps && !monitorMetadata) {
            resp.Diagnostics.AddAttributeError(
                req.Path, 
                "Invalid Resource Configuration",
                fmt.Sprintf("file_integrity_rule[%d] does not have any monitoring settings enabled. At least one monitoring setting must be set to true.", idx), 
            )
            continue
        }

        if monitorSubdirectories {
            if !monitorWriteOps{
                resp.Diagnostics.AddAttributeError(
                    req.Path, 
                    "Invalid Resource Configuration", 
                    fmt.Sprintf("file_integrity_rule[%d] has monitor_subdirectories enabled and monitor_write_ops disabled. Both arguments must be enabled when using monitor_subdirectories.", idx), 
                )
            }

            if (monitorReadOps && monitorMetadata) {
                resp.Diagnostics.AddAttributeError(
                    req.Path, 
                    "Invalid Resource Configuration", 
                    fmt.Sprintf("file_integrity_rule[%d] has monitor_subdirectories enabled along with monitor_read_ops and monitor_metadata_changes. Only monitor_write_ops can be enabled when using monitor_subdirectories.", idx), 
                )
            } else if monitorReadOps {
                resp.Diagnostics.AddAttributeError(
                    req.Path, 
                    "Invalid Resource Configuration", 
                    fmt.Sprintf("file_integrity_rule[%d] has monitor_subdirectories enabled along with monitor_read_ops. Only monitor_write_ops can be enabled when using monitor_subdirectories.", idx), 
                )
            } else if monitorMetadata {
                resp.Diagnostics.AddAttributeError(
                    req.Path, 
                    "Invalid Resource Configuration", 
                    fmt.Sprintf("file_integrity_rule[%d] has monitor_subdirectories enabled along with monitor_metadata. Only monitor_write_ops can be enabled when using monitor_subdirectories.", idx), 
                )
            }
        }

        if slices.Contains(filePaths, filePath) {
            resp.Diagnostics.AddAttributeError(
                req.Path, 
                "Invalid Resource Configuration", 
                fmt.Sprintf("file_integrity_rule[%d] has a duplicate file/directory path. Paths values must be unique.", idx), 
            )
        } else {
            filePaths = append(filePaths, filePath)
        }
    }

    return
}
