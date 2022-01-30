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
	return t
}
