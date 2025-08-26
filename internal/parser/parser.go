package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/fdanctl/jsontypify/internal/utils"
)

type Lang string

const (
	GO Lang = "go"
	TS Lang = "ts"
)

var validLangs = []string{"go", "ts"}

func findClosingIdx(b []byte, op byte) (int, error) {
	if op != '[' && op != '{' {
		return -1, fmt.Errorf("op char must be '[' or '{'")
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
			re := regexp.MustCompile(`[^\\]"`)
			idx := re.FindIndex(b[i:])
			i += idx[1]
		}
		if b[i] == op+2 {
			openingCount--
			if openingCount <= 0 {
				return i, nil
			}
		}
		i++
	}
	return -1, fmt.Errorf("closing %c not found in %s", op, b)
}

func isDate(b []byte) bool {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	return re.Match(b)
}

func findType(b []byte, param string, allMaps map[string]map[string]string) (string, int, error) {
	var i int
	for i < len(b) {
		if b[i] == '"' {
			re := regexp.MustCompile(`,"\w+":|$`)
			endIdxs := re.FindIndex(b[i:])
			if isDate(b[i : i+endIdxs[0]]) {
				return "time.Time", i + endIdxs[0], nil
			} else {
				return "string", i + endIdxs[0], nil
			}
		}

		if b[i] == 't' || b[i] == 'f' {
			re := regexp.MustCompile(`,"\w+":|$`)
			endIdxs := re.FindIndex(b[i:])
			return "bool", i + endIdxs[0], nil
		}
		// 0 - 9 char have rune/byte values [48-57]
		if b[i] >= 48 && b[i] <= 57 {
			re := regexp.MustCompile(`,"\w+":|$`)
			endIdxs := re.FindIndex(b[i:])
			re = regexp.MustCompile(`.`)
			if re.Match(b[i : i+endIdxs[0]]) {
				return "float64", i + endIdxs[0], nil
			}
			return "int", i + endIdxs[0], nil
		}

		if b[i] == '{' {
			idx, err := findClosingIdx(b[i:], '{')
			if err != nil {
				log.Fatal(err)
			}
			makeTypeMap(b[i:idx+1], utils.SnakeToCamelCase(param), allMaps)
			return utils.Capitalize(param), i + idx + 1, nil
		}

		if b[i] == '[' {
			t, _, err := findType(b[i+1:], param, allMaps)
			if err != nil {
				break
			}
			idx, err := findClosingIdx(b[i:], '[')
			if err != nil {
				log.Fatal(err)
			}
			return "[]" + t, i + idx + 1, nil
		}

		if b[i] == 'n' {
			re := regexp.MustCompile(`,"\w+":|$`)
			endIdxs := re.FindIndex(b[i:])
			return "null", i + endIdxs[0], nil
		}
		i++
	}

	return "", -1, fmt.Errorf("Malformed json")
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

		paramType, end, err := findType(b[i+paramIdxs[1]:], param, allMaps)
		if err != nil {
			log.Fatal(err)
		}

		typeMap[param] = paramType
		i += paramIdxs[1] + end + 1

	}
	allMaps[name] = typeMap
}

func goStruct(allMaps *map[string]map[string]string, indent int) string {
	var indentStr string
	for range indent {
		indentStr += " "
	}

	var str string
	for name, m := range *allMaps {
		str += fmt.Sprintf("type %s struct {\n", utils.Capitalize(name))
		for p, t := range m {
			str += fmt.Sprintf("%s%s %s `json:\"%s\"`\n", indentStr, utils.Capitalize(utils.SnakeToCamelCase(p)), t, p)
		}
		str += "}\n\n"
	}
	return str
}

func tsInterface(allMaps *map[string]map[string]string, indent int) string {
	var indentStr string
	for range indent {
		indentStr += " "
	}

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
			re = regexp.MustCompile(`time.Time`)
			t = string(re.ReplaceAll([]byte(t), []byte("Date")))
			str += fmt.Sprintf("%s%s: %s;\n", indentStr, utils.SnakeToCamelCase(p), t)
		}
		str += "}\n\n"
	}
	return str
}

func IsValidLang(s string) bool {
	return slices.Contains(validLangs, s)
}

func GetValidLangs() string {
	return strings.Join(validLangs, ", ")
}

func ParseTypes(s []byte, lang Lang, indent int) string {
	if !json.Valid(s) {
		log.Fatal("Invalid json")
	}

	re := regexp.MustCompile(`(?m)^\s+|\n`)
	flaten := re.ReplaceAll([]byte(s), []byte(""))

	allMaps := make(map[string]map[string]string, 0)
	makeTypeMap(flaten, "main", allMaps)

	if lang == GO {
		return goStruct(&allMaps, indent)
	} else {
		return tsInterface(&allMaps, indent)
	}
}
