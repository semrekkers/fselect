package fselect

import "reflect"

func getFieldName(v *reflect.StructField) string {
	name := v.Name
	if tag, ok := v.Tag.Lookup(StructTagKey); ok {
		name = tag
	}
	return name
}

func sliceContains(v string, slice []string) bool {
	for _, str := range slice {
		if str == v {
			return true
		}
	}
	return false
}

func repeatString(v string, sep string, times int) string {
	stringLen := (len(v) * times) + (len(sep) * (times - 1))
	out := make([]byte, stringLen)

	outp := copy(out, v)
	for i := 0; i < times; i++ {
		outp += copy(out[outp:], sep)
		outp += copy(out[outp:], v)
	}

	return string(out)
}
