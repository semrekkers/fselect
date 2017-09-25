package qbuilder

import "bytes"

// sliceContains returns whether the string slice contains v.
func sliceContains(v string, slice []string) bool {
	for _, str := range slice {
		if str == v {
			return true
		}
	}
	return false
}

// repeatString repeats string v, with separator sep, n times.
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

// joinStringsWithSuffix joins all strings in v concatenated with suffix and separated by sep.
func joinStringsWithSuffix(v []string, suffix string, sep string) string {
	// calculate length of all strings
	var length int
	for _, str := range v {
		length += len(str)
	}
	length += (len(suffix) * len(v)) + (len(sep) * (len(v) - 1))

	out := make([]byte, length)
	cursor := copy(out, v[0])
	cursor += copy(out[cursor:], suffix)
	for i := 1; i < len(v); i++ {
		cursor += copy(out[cursor:], sep)
		cursor += copy(out[cursor:], v[i])
		cursor += copy(out[cursor:], suffix)
	}

	return string(out)
}

func joinTwoSlices(a, b []string, suffix, sep string) string {
	var buf bytes.Buffer
	lastI := len(a) - 1
	for i := 0; i < len(a); i++ {
		buf.WriteString(a[i])
		buf.WriteString(suffix)
		buf.WriteString(b[i])
		if i < lastI {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}
