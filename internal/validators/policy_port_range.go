package validators

import (
	"context"
	"fmt"

	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func PortRangesAreValid() portRangesAreValid {
	return portRangesAreValid{}
}

type portRangesAreValid struct{}

func (v portRangesAreValid) Description(ctx context.Context) string {
	return ""
}

func (v portRangesAreValid) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (v portRangesAreValid) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	util.LogDebug(ctx, "Executing PortRangesAreValid")

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	values := make([]string, 0, len(req.ConfigValue.Elements()))
	resp.Diagnostics.Append(util.ListToStringSlice(ctx, &req.ConfigValue, &values)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for idx, value := range values {
		start, end, diags := models.GetPortRangeValues(value)

		if diags.HasError() {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(idx),
				diags[0].Summary(),
				diags[0].Detail(),
			)
			continue

		}

		if start == end {
			if !isValidTCPPort(start) {
				resp.Diagnostics.AddAttributeError(
					req.Path.AtListIndex(idx),
					"Invalid Resource Configuration",
					fmt.Sprintf("Invalid value specified: %s\nValues must be between greater than or equal to 0 and less than or equal to 65535.", value),
				)
				continue
			}
		} else {
			if !isValidTCPPort(start) {
				resp.Diagnostics.AddAttributeError(
					req.Path.AtListIndex(idx),
					"Invalid Resource Configuration",
					fmt.Sprintf("Invalid port range start value specified: %d\nValues must be between greater than or equal to 0 and less than or equal to 65535.", start),
				)
				continue
			}

			if !isValidTCPPort(end) {
				resp.Diagnostics.AddAttributeError(
					req.Path.AtListIndex(idx),
					"Invalid Resource Configuration",
					fmt.Sprintf("Invalid port range end value specified: %d\nValues must be between greater than or equal to 0 and less than or equal to 65535.", end),
				)
				continue
			}

			if start > end {
				resp.Diagnostics.AddAttributeError(
					req.Path.AtListIndex(idx),
					"Invalid Resource Configuration",
					fmt.Sprintf("Invalid port range specified: %s\nThe starting value in the port range must be greater than the ending value.", value),
				)
				continue
			}
		}
	}

	util.LogDebug(ctx, "Finishing PortRangesAreValid execution")

	return
}

func isValidTCPPort(port int) bool {
	return (port >= 0 && port <= 65535)
}
