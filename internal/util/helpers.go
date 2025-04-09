package util

import (
    "context"
    "slices"
    "strconv"
    //"regexp"
    
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Converts basetypes.ListValue to string slice and populates it in response
func ListToStringSlice(ctx context.Context, l *basetypes.ListValue, response *[]string) diag.Diagnostics {
    if (l == nil || (*l).IsNull()) {
        return nil
    }

    diags := l.ElementsAs(ctx, response, false)

    return diags 
}

//func StringToInt(str string) (int, diag.Diagnostic) {
//    i, err := strconv.Atoi(str)
//    if err != nil {
//        diag.NewAttributeErrorDiagnostic(
//
//        )
//    }
//}

func StringToInt(str string) (int, error) {
    i, err := strconv.Atoi(str)
    if err != nil {
        return -1, err 
    }

    return i, nil
}

// Returns true if s1 shares any elements with s2
func SliceSharesOneOrMoreElements(s1 []string, s2 []string) bool {
    for _, elem := range s1 {
        if slices.Contains(s2, elem) {
            return true
        }
    }

    return false
}
