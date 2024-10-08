package util

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

func DLog(ctx context.Context, message string) {
    tflog.Debug(ctx, fmt.Sprintf("\n\n%s\n\n", message))
}
