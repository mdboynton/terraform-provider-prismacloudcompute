package util

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var logger *log.Logger

func init() {
	logger = log.Default()
	logger.SetPrefix("")
}

func HCLogDebug(ctx context.Context, message string) {
	tflog.Debug(ctx, fmt.Sprintf("\n\n%s\n\n", message))
}

func HCLogfDebug(ctx context.Context, object interface{}) {
	tflog.Debug(ctx, fmt.Sprintf("\n\n%v\n\n", object))
}
