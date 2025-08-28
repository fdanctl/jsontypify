package parser
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"regexp"
// 	"strconv"
// 	"strings"
//
// 	"github.com/fdanctl/jsontypify/internal/utils"
// )
//
// // change to handle ts and go (maybe map)
// func assertType(v any) string {
// 	switch t := v.(type) {
// 	case float64:
// 		if math.Trunc(t) == t {
// 			return "int"
// 		}
// 		return "float64"
// 	case string:
// 		if utils.IsDate(t) {
// 			return "time.Time"
// 		}
// 		return "string"
// 	case []any:
// 		return "array"
// 	case map[string]any:
// 		return "map"
// 	case bool:
// 		return "bool"
// 	default:
// 		return "any"
//
// 	}
// }
//
// func safeName(name string, allMaps *map[string]map[string]string) string {
// 	var n int
// 	_, ok := (*allMaps)[name]
// 	for ok {
// 		n++
// 		k := name + strconv.Itoa(n)
// 		_, ok = (*allMaps)[k]
// 	}
// 	if n == 0 {
// 		return name
// 	} else {
// 		name += strconv.Itoa(n)
// 		return name
// 	}
// }
//
// func makeTypeMap(name string, obj *map[string]any, allMaps *map[string]map[string]string, keys *map[string][]string) {
// 	typeMap := make(map[string]string)
// 	orderedFields := make([]string, 0)
//
// 	for k, v := range *obj {
// 		t := assertType(v)
//
// 		if t == "map" {
// 			m := v.(map[string]any)
// 			t = utils.Capitalize(utils.SnakeToCamelCase(k))
// 			t = safeName(t, allMaps)
// 			makeTypeMap(t, &m, allMaps, keys)
// 		}
//
// 		if t == "array" {
// 			arr := v.([]any)
// 			if len(arr) == 0 {
// 				t = "[]any"
// 			} else {
// 				tArr := make([]string, len(arr))
// 				isAny := false
// 				for i, val := range arr {
// 					tipo := assertType(val)
//
// 					if tipo == "map" {
// 						m := arr[i].(map[string]any)
// 						tipo = utils.Capitalize(utils.SnakeToCamelCase(k))
//
// 						if i == 0 {
// 							tipo = safeName(tipo, allMaps)
// 							makeTypeMap(tipo, &m, allMaps, keys)
// 						} else if !isAny && tArr[0] != tipo {
// 							isAny = true
// 							tipo = safeName(tipo, allMaps)
// 							makeTypeMap(tipo, &m, allMaps, keys)
// 						}
// 					} else {
// 						if !isAny && i > 0 && tArr[0] != tipo {
// 							isAny = true
// 						}
// 					}
// 					tArr[i] = tipo
// 				}
// 				if isAny {
// 					t = "[" + strings.Join(tArr, ", ") + "]"
// 				} else {
// 					t = "[]" + tArr[0]
// 				}
// 				// of := assertType(arr[0])
// 				// if of == "map" {
// 				// 	m := arr[0].(map[string]any)
// 				// 	of = utils.Capitalize(utils.SnakeToCamelCase(k))
// 				// 	of = safeName(of, allMaps)
// 				// 	makeTypeMap(of, &m, allMaps)
// 				// }
// 				// t = "[]" + of
// 			}
// 		}
//
// 		typeMap[k] = t
// 		orderedFields = append(orderedFields, k)
//
// 	}
// 	name = safeName(name, allMaps)
// 	(*allMaps)[name] = typeMap
// 	(*keys)[name] = orderedFields
//
// 	(*keys)["root"] = append((*keys)["root"], name)
// }
//
// func ParseTypes(s []byte, lang Lang, indent int, name string) string {
// 	if !json.Valid(s) {
// 		panic("Invalid json")
// 	}
//
// 	var data any
// 	err := json.Unmarshal(s, &data)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	allMaps := make(map[string]map[string]string, 0)
// 	keys := make(map[string][]string, 0)
//
// 	switch val := data.(type) {
// 	case map[string]any:
// 		makeTypeMap(name, &val, &allMaps, &keys)
// 	case []any:
// 		if len(val) > 0 {
// 			switch first := val[0].(type) {
// 			case map[string]any:
// 				makeTypeMap(name, &first, &allMaps, &keys)
// 			default:
// 				panic("Invalid json")
// 			}
// 		}
// 	default:
// 		panic("Invalid json")
// 	}
//
// 	if lang == GO {
// 		return goStruct(&allMaps, indent, &keys)
// 	} else {
// 		return tsInterface(&allMaps, indent)
// 	}
// }
//
// func goStruct(allMaps *map[string]map[string]string, indent int, keys *map[string][]string) string {
// 	var indentStr string
// 	for range indent {
// 		indentStr += " "
// 	}
//
// 	var str string
// 	for _, k := range (*keys)["root"] {
// 		// str += fmt.Sprintf("type %s struct {\n", utils.Capitalize(k))
// 		for _, param := range (*keys)[k] {
// 			println(k, param)
// 			// str += fmt.Sprintf("%s%s %s `json:\"%s\"`\n", indentStr, utils.Capitalize(utils.SnakeToCamelCase(param)), (*allMaps)[k][param], param)
// 		}
// 		// str += "}\n\n"
// 	}
// 	return str
// }
//
// func tsInterface(allMaps *map[string]map[string]string, indent int) string {
// 	var indentStr string
// 	for range indent {
// 		indentStr += " "
// 	}
//
// 	var str string
// 	for name, m := range *allMaps {
// 		str += fmt.Sprintf("interface %s {\n", utils.Capitalize(name))
// 		for p, t := range m {
// 			re := regexp.MustCompile(`(^|\[\])(int|float64)$`)
// 			t = re.ReplaceAllStringFunc(t, func(s string) string {
// 				groups := re.FindStringSubmatch(s)
// 				groups[2] = "number"
// 				return groups[1] + groups[2]
// 			})
//
// 			re = regexp.MustCompile(`(^|\[\])(bool)`)
// 			t = re.ReplaceAllStringFunc(t, func(s string) string {
// 				groups := re.FindStringSubmatch(s)
// 				groups[2] = "bool"
// 				return groups[1] + groups[2]
// 			})
//
// 			re = regexp.MustCompile(`(^|\[\])(time\.Time)`)
// 			t = re.ReplaceAllStringFunc(t, func(s string) string {
// 				groups := re.FindStringSubmatch(s)
// 				groups[2] = "Date"
// 				return groups[1] + groups[2]
// 			})
//
// 			if t[0] == '[' {
// 				t = t[2:] + t[0:2]
// 			}
//
// 			str += fmt.Sprintf("%s%s: %s;\n", indentStr, utils.SnakeToCamelCase(p), t)
// 		}
// 		str += "}\n\n"
// 	}
// 	return str
// }
