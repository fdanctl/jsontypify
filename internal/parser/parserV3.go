package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/fdanctl/jsontypify/internal/utils"
)

// change to handle ts and go (maybe map)
func assertType(v any) string {
	return ""
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
				parseObj(t, dec, allMaps, keys)
			}
			if val == json.Delim('[') {
				t = parseArr(name, dec, allMaps, keys)
				fmt.Println("Dude")
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
		fmt.Println(k)
		switch val := tk.(type) {
		case json.Delim:
			t = utils.Capitalize(k)
			if val == json.Delim('{') {
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
	(*keys)[name] = orderedFields

	(*keys)["root"] = append((*keys)["root"], name)
}

func ParseTypes(s []byte, lang Lang, indent int, name string) string {
	if !json.Valid(s) {
		panic("Invalid json")
	}

	dec := json.NewDecoder(bytes.NewReader(s))
	dec.UseNumber()

	allMaps := make(map[string]map[string]string, 0)
	keys := make(map[string][]string, 0)

	// dec.Token()
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			panic(t)
		}

		// keys := make([]string, 0)
		switch t.(type) {
		case json.Delim:
			parseObj(name, dec, &allMaps, &keys)
		default:
			fmt.Println("hello", t)
		}
	}

	return goStruct(&allMaps, 4, &keys)
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
			str += fmt.Sprintf("%s%s %s `json:\"%s\"`\n", indentStr, utils.Capitalize(utils.SnakeToCamelCase(param)), (*allMaps)[k][param], param)
		}
		str += "}\n\n"
	}
	return str
}
