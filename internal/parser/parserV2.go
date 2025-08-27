package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/fdanctl/jsontypify/internal/utils"
)

// change to handle ts and go (maybe map)
func assertType(v any) string {
	switch t := v.(type) {
	case float64:
		// verify if it's int
		return "float64"
	case string:
		if utils.IsDate(t) {
			return "time.Time"
		}
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "map"
	case bool:
		return "bool"
	default:
		return "any"

	}
}

func safeName(name string, allMaps *map[string]map[string]string) string {
	var n int
	_, ok := (*allMaps)[name]
	for ok {
		n++
		k := name + strconv.Itoa(n)
		_, ok = (*allMaps)[k]
	}
	if n == 0 {
		return name
	} else {
		name += strconv.Itoa(n)
		return name
	}
}

func makeTypeMap(name string, obj *map[string]any, allMaps *map[string]map[string]string) {
	typeMap := make(map[string]string)
	for k, v := range *obj {
		t := assertType(v)

		if t == "map" {
			m := v.(map[string]any)
			t = utils.Capitalize(utils.SnakeToCamelCase(k))
			t = safeName(t, allMaps)
			makeTypeMap(t, &m, allMaps)
		}

		if t == "array" {
			arr := v.([]any)
			if len(arr) == 0 {
				t = "[]any"
			} else {
				of := assertType(arr[0])
				if of == "map" {
					m := arr[0].(map[string]any)
					of = utils.Capitalize(utils.SnakeToCamelCase(k))
					of = safeName(of, allMaps)
					makeTypeMap(of, &m, allMaps)
				}
				t = "[]" + of
			}
		}

		typeMap[k] = t
	}
	name = safeName(name, allMaps)
	(*allMaps)[name] = typeMap
}

func ParseTypes(s []byte, lang Lang, indent int, name string) string {
	if !json.Valid(s) {
		panic("Invalid json")
	}

	var data any
	err := json.Unmarshal(s, &data)
	if err != nil {
		panic(err)
	}

	allMaps := make(map[string]map[string]string, 0)

	switch val := data.(type) {
	case map[string]any:
		makeTypeMap(name, &val, &allMaps)
	case []any:
		if len(val) > 0 {
			switch first := val[0].(type) {
			case map[string]any:
				makeTypeMap(name, &first, &allMaps)
			default:
				panic("Invalid json")
			}
		}
	default:
		panic("Invalid json")
	}

	if lang == GO {
		return goStruct(&allMaps, indent)
	} else {
		return tsInterface(&allMaps, indent)
	}
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
		str += fmt.Sprintf("interface %s {\n", utils.Capitalize(name))
		for p, t := range m {
			re := regexp.MustCompile(`(^|\[\])(int|float64)$`)
			t = re.ReplaceAllStringFunc(t, func(s string) string {
				groups := re.FindStringSubmatch(s)
				groups[2] = "number"
				return groups[1] + groups[2]
			})

			re = regexp.MustCompile(`(^|\[\])(bool)`)
			t = re.ReplaceAllStringFunc(t, func(s string) string {
				groups := re.FindStringSubmatch(s)
				groups[2] = "bool"
				return groups[1] + groups[2]
			})

			re = regexp.MustCompile(`(^|\[\])(time\.Time)`)
			t = re.ReplaceAllStringFunc(t, func(s string) string {
				groups := re.FindStringSubmatch(s)
				groups[2] = "Date"
				return groups[1] + groups[2]
			})

			if t[0] == '[' {
				t = t[2:] + t[0:2]
			}

			str += fmt.Sprintf("%s%s: %s;\n", indentStr, utils.SnakeToCamelCase(p), t)
		}
		str += "}\n\n"
	}
	return str
}
