package tsproto

import "strings"

var funcMap = map[string]interface{}{
	"convertGoType": convertGoType,
}

func convertGoType(t string) string {
	if strings.HasPrefix(t, "*") {
		// TODO: handle pointers
		return "unknown"
	}
	if strings.HasPrefix(t, "[]") {
		// TODO: handle slices
		return t[2:] + "[]"
	}
	if t == "uint" {
		// TODO: handle numeric types
		return "number"
	}
	return t
}
