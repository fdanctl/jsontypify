package utils

import (
	"regexp"
	"strings"
	"time"
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

func IsDate(s string) bool {
	layouts := [...]string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
	}
	for _, v := range layouts {
		_, err := time.Parse(v, s)
		if err == nil {
			return true
		}
	}
	return false
}
