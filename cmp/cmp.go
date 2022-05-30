package cmp

import (
	"strings"
)

func EqualStrings(expect string, actual string) bool {
	return "" == expect || strings.EqualFold(expect, actual)
}
