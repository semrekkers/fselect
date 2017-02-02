package fselect

import (
	"reflect"
	"strings"
)

// getFieldName returns the name of a struct field.
func getFieldName(v *reflect.StructField) string {
	name := v.Name

	if tag, ok := v.Tag.Lookup(StructTagKey); ok {
		// The field has a tag like `col:"<name>"`, use this <name> instead.
		name = tag
	}

	return name
}

// sliceContains returns whether slice contains v. If ignoreCase is true then differing cases will be ignored.
func sliceContains(v string, ignoreCase bool, slice []string) bool {
	for _, str := range slice {
		if !ignoreCase && str == v {
			return true
		}
		if ignoreCase && strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

// repeatString repeats string v with seperator sep, n times.
func repeatString(v string, sep string, n int) string {
	stringLen := (len(v) * n) + (len(sep) * (n - 1))
	out := make([]byte, stringLen)

	outp := copy(out, v)
	for i := 0; i < n; i++ {
		outp += copy(out[outp:], sep)
		outp += copy(out[outp:], v)
	}

	return string(out)
}
