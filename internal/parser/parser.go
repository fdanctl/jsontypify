package parser

import (
	"fmt"
	"regexp"

	"github.com/fdanctl/jsontypify/internal/utils"
)

func findClosingIdx(b []byte, op byte) int {
	if op != '[' && op != '{' {
		fmt.Println("error")
		return -1
	}

	var openingCount int
	var i int
	for i < len(b) {
		if b[i] == op {
			openingCount++
			i++
			continue
		}
		if b[i] == '"' {
			i++
			re := regexp.MustCompile(`"`)
			idx := re.FindIndex(b[i:])
			i += idx[0]
		}
		if b[i] == op+2 {
			openingCount--
			if openingCount <= 0 {
				return i
			}
		}
		i++
	}
	return -1
}

func findType(b []byte, param string, allMaps map[string]map[string]string) (string, int) {
	for i := range b {
		if b[i] == '"' {
			re := regexp.MustCompile(`,"\w+":|$`)
			endIdxs := re.FindIndex(b[i:])
			return "string", i + endIdxs[0]
		}

		if b[i] == 't' || b[i] == 'f' {
			re := regexp.MustCompile(`,"\w+":|$`)
			endIdxs := re.FindIndex(b[i:])
			return "bool", i + endIdxs[0]
		}
		// 0 - 9 char have rune/byte values [48-57]
		if b[i] >= 48 && b[i] <= 57 {
			re := regexp.MustCompile(`,"\w+":|$`)
			endIdxs := re.FindIndex(b[i:])
			re = regexp.MustCompile(`.`)
			if re.Match(b[i : i+endIdxs[0]]) {
				return "float64", i + endIdxs[0]
			}
			return "int", i + endIdxs[0]
		}

		if b[i] == '{' {
			idx := findClosingIdx(b[i:], '{')
			makeTypeMap(b[i:idx], utils.SnakeToCamelCase(param), allMaps)
			return utils.Capitalize(param), i + idx + 1
		}

		if b[i] == '[' {
			t, _ := findType(b[i+1:], param, allMaps)
			return "[]" + t, i + findClosingIdx(b[i:], '[')
		}
		i++
	}

	return "error", -1
}

func makeTypeMap(b []byte, name string, allMaps map[string]map[string]string) {
	typeMap := make(map[string]string)
	var i int
	for i < len(b) {

		re := regexp.MustCompile(`"[^"]*":`)

		if !re.Match(b[i:]) {
			break
		}

		paramIdxs := re.FindIndex(b[i:])
		param := string(b[i+paramIdxs[0]+1 : i+paramIdxs[1]-2])

		paramType, end := findType(b[i+paramIdxs[1]:], param, allMaps)

		typeMap[param] = paramType
		i += paramIdxs[1] + end + 1

	}
	allMaps[name] = typeMap
}

type Lang int

const (
	GO Lang = iota
	TS
)

func goStruct(allMaps *map[string]map[string]string) string {
	var str string
	for name, m := range *allMaps {
		str += fmt.Sprintf("type %s struct {\n", utils.Capitalize(name))
		for p, t := range m {
			str += fmt.Sprintf("  %s %s `json:\"%s\"`\n", utils.Capitalize(utils.SnakeToCamelCase(p)), t, p)
		}
		str += "}\n\n"
	}
	return str
}

func tsInterface(allMaps *map[string]map[string]string) string {
	var str string
	for name, m := range *allMaps {
		str += fmt.Sprintf("type interface %s {\n", utils.Capitalize(name))
		for p, t := range m {
			if t[0] == '[' {
				t = t[2:] + t[0:2]
			}
			re := regexp.MustCompile(`(int|float64)`)
			t = string(re.ReplaceAll([]byte(t), []byte("number")))
			re = regexp.MustCompile(`bool`)
			t = string(re.ReplaceAll([]byte(t), []byte("boolean")))
			str += fmt.Sprintf("  %s: %s;\n", p, t)
		}
		str += "}\n\n"
	}
	return str
}

func ParseTypes(s string, lang Lang) string {
	re := regexp.MustCompile(`(?m)^\s+|\n`)
	flaten := re.ReplaceAll([]byte(s), []byte(""))

	allMaps := make(map[string]map[string]string, 0)
	makeTypeMap(flaten, "main", allMaps)

	if lang == GO {
		return goStruct(&allMaps)
	} else {
		return tsInterface(&allMaps)
	}
}
