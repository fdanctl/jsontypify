package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/fdanctl/jsontypify/internal/utils"
)

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

func parseArr(name string, dec *json.Decoder, allMaps *map[string]map[string]string, keys *map[string][]string) string {
	tArr := make([]string, 0)
	isAny := false

	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			panic(err)
		}

		var t string
		switch val := tk.(type) {
		case json.Delim:
			if val == json.Delim('{') {
				t = utils.Capitalize(name)
				// t = safeName(t, allMaps)
				parseObj(t, dec, allMaps, keys)
			}
			if val == json.Delim('[') {
				t = parseArr(name, dec, allMaps, keys)
			}
		case bool:
			t = "bool"
		case json.Number:
			re := regexp.MustCompile(`\.`)
			if re.Match([]byte(val)) {
				t = "float64"
			} else {
				t = "int"
			}
		case string:
			if utils.IsDate(t) {
				t = "time.Time"
			}
			t = "string"
		case nil:
			t = "any"
		}

		if !isAny && len(tArr) > 0 && tArr[0] != t {
			isAny = true
		}
		tArr = append(tArr, t)
	}

	// Consume '}'
	_, err := dec.Token()
	if err != nil {
		panic(err)
	}

	if len(tArr) == 0 {
		return "[]any"
	} else if isAny {
		return "[" + strings.Join(tArr, ", ") + "]"
	} else {
		return "[]" + tArr[0]
	}
}

func parseObj(name string, dec *json.Decoder, allMaps *map[string]map[string]string, keys *map[string][]string) {
	typeMap := make(map[string]string)
	orderedFields := make([]string, 0)

	// Consume '{'
	// _, err := dec.Token()
	// if err != nil {
	// 	panic(err)
	// }

	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			panic(err)
		}
		k := tk.(string)

		tk, err = dec.Token()
		if err != nil {
			panic(err)
		}

		var t string
		switch val := tk.(type) {
		case json.Delim:
			t = utils.Capitalize(k)
			if val == json.Delim('{') {
				t = safeName(t, allMaps)
				parseObj(t, dec, allMaps, keys)
			}
			if val == json.Delim('[') {
				t = parseArr(t, dec, allMaps, keys)
			}
		case bool:
			t = "bool"
		case json.Number:
			re := regexp.MustCompile(`\.`)
			if re.Match([]byte(val)) {
				t = "float64"
			} else {
				t = "int"
			}
		case string:
			if utils.IsDate(t) {
				t = "time.Time"
			}
			t = "string"
		case nil:
			t = "any"
		}

		typeMap[k] = t
		orderedFields = append(orderedFields, k)
	}

	// Consume '}'
	_, err := dec.Token()
	if err != nil {
		panic(err)
	}
	(*allMaps)[name] = typeMap

	if _, ok := (*keys)[name]; !ok {
		(*keys)["root"] = append((*keys)["root"], name)
	}

	(*keys)[name] = orderedFields
}

func ParseTypes(s io.Reader, lang Lang, indent int, name string) string {
	// if !json.Valid(s) {
	// 	panic("Invalid json")
	// }
	//
	dec := json.NewDecoder(s)
	dec.UseNumber()

	allMaps := make(map[string]map[string]string, 0)
	keys := make(map[string][]string, 0)

	loop:
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			panic(t)
		}

		switch val := t.(type) {
		case json.Delim:
			_, ok := allMaps[name]
			if val == json.Delim('{') && !ok {
				parseObj(name, dec, &allMaps, &keys)
			} else if ok {
				break loop
			}
		default: 
			panic(fmt.Sprint("unexpected char:", val))
		}
	}

	if lang == GO {
		return goStruct(&allMaps, indent, &keys)
	} else {
		return tsInterface(&allMaps, indent, &keys)
	}
}

func goStruct(allMaps *map[string]map[string]string, indent int, keys *map[string][]string) string {
	var indentStr string
	for range indent {
		indentStr += " "
	}

	var str string
	for _, k := range (*keys)["root"] {
		str += fmt.Sprintf("type %s struct {\n", utils.Capitalize(k))
		for _, param := range (*keys)[k] {
			t := (*allMaps)[k][param]
			re := regexp.MustCompile(`\[[^\]]{1}.*\]`)
			t = re.ReplaceAllStringFunc(t, func(s string) string {
				re := regexp.MustCompile(`[A-Z]\w*`)
				unusedMaps := re.FindAllString(t, -1)
				for _, val := range unusedMaps {
					delete(*keys, val)
				}
				return "[]any"
			})
			t = re.ReplaceAllString(t, "[]any")
			str += fmt.Sprintf("%s%s %s `json:\"%s\"`\n", indentStr, utils.Capitalize(utils.SnakeToCamelCase(param)), t, param)
		}
		str += "}\n\n"
	}
	return str
}

func tsInterface(allMaps *map[string]map[string]string, indent int, keys *map[string][]string) string {
	var indentStr string
	for range indent {
		indentStr += " "
	}

	var str string
	for _, k := range (*keys)["root"] {
		str += fmt.Sprintf("type %s struct {\n", utils.Capitalize(k))
		for _, param := range (*keys)[k] {
			t := (*allMaps)[k][param]
			re := regexp.MustCompile(`int|float64`)
			t = re.ReplaceAllString(t, "number")

			re = regexp.MustCompile(`bool`)
			t = re.ReplaceAllString(t, "boolean")

			re = regexp.MustCompile(`time\.Time`)
			t = re.ReplaceAllString(t, "Date")

			re = regexp.MustCompile(`((\[\])+)(\w+)(,|$)`)
			t = re.ReplaceAllStringFunc(t, func(s string) string {
				groups := re.FindStringSubmatch(s)
				return groups[3] + groups[1] + groups[4]
			})

			str += fmt.Sprintf("%s%s: %s;\n", indentStr, utils.SnakeToCamelCase(param), t)
		}
		str += "}\n\n"
	}
	return str
}
