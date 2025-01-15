package util

import (
    "slices"
    
    //"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Returns true if s1 shares any elements with s2
func SliceSharesOneOrMoreElements(s1 []string, s2 []string) bool {
    for _, elem := range s1 {
        if slices.Contains(s2, elem) {
            return true
        }
    }

    return false
}
