package util

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
    Log *log.Logger
    debugLogPaddingEnvVar string = "PRISMACLOUDCOMPUTE_DEBUG_LOG_PADDING"
    debugLogPadding string = ""
    debugLogPaddingEnabled bool
)

func init() {
    // Assign address of default logger to Log
	Log = log.Default()

    // Remove log message datetime prefix by clearing the Ldate and Ltime flags using a bit clear operaton
    Log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

    // Add newlines around log messages for easier troubleshooting if PRISMACLOUDCOMPUTE_DEBUG_LOG_PADDING env var is set to true
    err := GetEnvironmentVariable(debugLogPaddingEnvVar, &debugLogPaddingEnabled)
    if err != nil {
        Log.Printf("Failed to parse %s value: %s", debugLogPaddingEnvVar, err.Error())
    } else if (debugLogPaddingEnabled == true) {
        debugLogPadding = "\n\n" 
    }
}

func LogDebug(message string) {
    Log.Printf(fmt.Sprintf("%s%s%s\n", debugLogPadding, message, debugLogPadding))
}

func LogfDebug(object interface{}) {
    Log.Printf(fmt.Sprintf("%s%s%s\n", debugLogPadding, object, debugLogPadding))
}

func HCLogInfo(ctx context.Context, message string) {
	tflog.Info(ctx, fmt.Sprintf("%s%s%s", debugLogPadding, message, debugLogPadding))
}

func HCLogWarn(ctx context.Context, message string) {
	tflog.Warn(ctx, fmt.Sprintf("%s%s%s", debugLogPadding, message, debugLogPadding))
}

func HCLogDebug(ctx context.Context, message string) {
	tflog.Debug(ctx, fmt.Sprintf("%s%s%s", debugLogPadding, message, debugLogPadding))
}

func HCLogfDebug(ctx context.Context, object interface{}) {
	tflog.Debug(ctx, fmt.Sprintf("%s%v%s", debugLogPadding, object, debugLogPadding))
}
