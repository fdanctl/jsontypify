package utils

import (
	"regexp"
	"strings"
)

func SnakeToCamelCase(s string) string {
	re := regexp.MustCompile(`_.`)
	bytes := re.ReplaceAllFunc([]byte(s), func(b []byte) []byte {
		return []byte(strings.ToUpper(string(b[1])))
	})
	return string(bytes)
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + string(s[1:])
}


