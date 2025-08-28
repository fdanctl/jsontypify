package parser

// import (
// 	// "encoding/json"
// 	// "fmt"
// 	// "log"
// 	// "regexp"
// 	"slices"
// 	// "strconv"
// 	"strings"
//
// 	// "github.com/fdanctl/jsontypify/internal/utils"
// )
// func findClosingIdx(b []byte, op byte) (int, error) {
// 	if op != '[' && op != '{' {
// 		return -1, fmt.Errorf("op char must be '[' or '{'")
// 	}
//
// 	var openingCount int
// 	var i int
// 	for i < len(b) {
// 		if b[i] == op {
// 			openingCount++
// 			i++
// 			continue
// 		}
// 		if b[i] == '"' {
// 			re := regexp.MustCompile(`[^\\]"`)
// 			idx := re.FindIndex(b[i:])
// 			if idx == nil {
// 				return -1, fmt.Errorf("unmatched \"")
// 			}
// 			i += idx[1]
// 			continue
// 		}
// 		if b[i] == op+2 {
// 			openingCount--
// 			if openingCount <= 0 {
// 				return i, nil
// 			}
// 		}
// 		i++
// 	}
// 	return -1, fmt.Errorf("closing %c not found in %s", op, b)
// }
//
// func isDate(b []byte) bool {
// 	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
// 	return re.Match(b)
// }
//
// func safelyAddMap(name string, v *map[string]string, allMaps *map[string]map[string]string) {
// 	var n int
// 	_, ok := (*allMaps)[name]
// 	for ok {
// 		n++
// 		k := name + strconv.Itoa(n)
// 		_, ok = (*allMaps)[k]
// 	}
// 	if n == 0 {
// 		(*allMaps)[name] = *v
// 	} else {
// 		k := name + strconv.Itoa(n)
// 		(*allMaps)[k] = *v
// 	}
// }
//
// func findType(b []byte, param string, allMaps *map[string]map[string]string) (string, int, error) {
// 	var i int
// 	for i < len(b) {
// 		if b[i] == '"' {
// 			re := regexp.MustCompile(`,"\w+":|$`)
// 			endIdxs := re.FindIndex(b[i:])
// 			if isDate(b[i : i+endIdxs[0]]) {
// 				return "time.Time", i + endIdxs[0], nil
// 			} else {
// 				return "string", i + endIdxs[0], nil
// 			}
// 		}
//
// 		if b[i] == 't' || b[i] == 'f' {
// 			re := regexp.MustCompile(`,"\w+":|$`)
// 			endIdxs := re.FindIndex(b[i:])
// 			return "bool", i + endIdxs[0], nil
// 		}
// 		// 0 - 9 char have rune/byte values [48-57]
// 		if b[i] >= 48 && b[i] <= 57 {
// 			re := regexp.MustCompile(`,"\w+":|$`)
// 			endIdxs := re.FindIndex(b[i:])
// 			re = regexp.MustCompile(`.`)
// 			if re.Match(b[i : i+endIdxs[0]]) {
// 				return "float64", i + endIdxs[0], nil
// 			}
// 			return "int", i + endIdxs[0], nil
// 		}
//
// 		if b[i] == '{' {
// 			if b[i+1] == '}' {
// 				v := make(map[string]string, 0) 
// 				safelyAddMap(utils.SnakeToCamelCase(param), &v, allMaps)
// 				return utils.Capitalize(param), i + 2, nil
// 			}
//
// 			idx, err := findClosingIdx(b[i:], '{')
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			makeTypeMap(b[i:idx+1], utils.SnakeToCamelCase(param), allMaps)
// 			return utils.Capitalize(utils.SnakeToCamelCase(param)), i + idx + 1, nil
// 		}
//
// 		if b[i] == '[' {
// 			if b[i+1] == ']' {
// 				return "any", i + 2, nil
// 			}
// 			t, _, err := findType(b[i+1:], param, allMaps)
// 			if err != nil {
// 				break
// 			}
// 			idx, err := findClosingIdx(b[i:], '[')
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			return "[]" + t, i + idx + 1, nil
// 		}
//
// 		if b[i] == 'n' {
// 			re := regexp.MustCompile(`,"\w+":|$`)
// 			endIdxs := re.FindIndex(b[i:])
// 			return "any", i + endIdxs[0], nil
// 		}
// 		i++
// 	}
//
// 	return "", -1, fmt.Errorf("Malformed json")
// }
//
// func makeTypeMap(b []byte, name string, allMaps *map[string]map[string]string) {
// 	typeMap := make(map[string]string)
// 	var i int
// 	for i < len(b) {
// 		re := regexp.MustCompile(`"[^"]*":`)
//
// 		if !re.Match(b[i:]) {
// 			break
// 		}
//
// 		paramIdxs := re.FindIndex(b[i:])
// 		param := string(b[i+paramIdxs[0]+1 : i+paramIdxs[1]-2])
//
// 		paramType, end, err := findType(b[i+paramIdxs[1]:], param, allMaps)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		typeMap[param] = paramType
// 		i += paramIdxs[1] + end + 1
//
// 	}
// 	safelyAddMap(name, &typeMap, allMaps)
// }
//
// func goStruct(allMaps *map[string]map[string]string, indent int) string {
// 	var indentStr string
// 	for range indent {
// 		indentStr += " "
// 	}
//
// 	var str string
// 	for name, m := range *allMaps {
// 		str += fmt.Sprintf("type %s struct {\n", utils.Capitalize(name))
// 		for p, t := range m {
// 			str += fmt.Sprintf("%s%s %s `json:\"%s\"`\n", indentStr, utils.Capitalize(utils.SnakeToCamelCase(p)), t, p)
// 		}
// 		str += "}\n\n"
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
//
// 			str += fmt.Sprintf("%s%s: %s;\n", indentStr, utils.SnakeToCamelCase(p), t)
// 		}
// 		str += "}\n\n"
// 	}
// 	return str
// }
// func ParseTypes(s []byte, lang Lang, indent int, name string) string {
// 	if !json.Valid(s) {
// 		log.Fatal("Invalid json")
// 	}
//
// 	re := regexp.MustCompile(`(?m)^\s+|\n`)
// 	flaten := re.ReplaceAll([]byte(s), []byte(""))
//
// 	allMaps := make(map[string]map[string]string, 0)
// 	if flaten[0] == '[' {
// 		idx, _:= findClosingIdx(flaten[1:], '{')
// 		makeTypeMap(flaten[1:idx], name, &allMaps)
// 	} else {
// 		makeTypeMap(flaten, name, &allMaps)
// 	}
//
// 	if lang == GO {
// 		return goStruct(&allMaps, indent)
// 	} else {
// 		return tsInterface(&allMaps, indent)
// 	}
// }
