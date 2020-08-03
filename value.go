package ltsvdoc

import (
	"strconv"
	"strings"
)

func parseValue(s string) interface{} {
	if !strings.ContainsAny(s, "\\") {
		return s
	}
	t, _ := strconv.Unquote("\"" + s + "\"")
	return t
}
