package parser

import (
	"slices"
	"strings"
)

type Lang string

const (
	GO Lang = "go"
	TS Lang = "ts"
)

var validLangs = []string{"go", "ts"}

func IsValidLang(s string) bool {
	return slices.Contains(validLangs, s)
}

func GetValidLangs() string {
	return strings.Join(validLangs, ", ")
}
