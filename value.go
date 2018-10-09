package ltsvdoc

import (
	"strconv"
	"strings"
)

func parseValue(s string) interface{} {
	if strings.IndexAny(s, "\\") < 0 {
		return s
	}
	t, _ := strconv.Unquote(s)
	return t
}
