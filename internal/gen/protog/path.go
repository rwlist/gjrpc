package protog

import "strings"

func PathToString(path []string) string {
	return strings.Join(path, ".")
}

func StringToPath(path string) []string {
	return strings.Split(path, ".")
}

func PathAppend(path []string, elems ...string) []string {
	var res []string
	res = append(res, path...)
	res = append(res, elems...)
	return res
}
