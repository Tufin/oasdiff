package generator

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
)

func Generate() error {
	out, err := os.Create("messages.yaml")
	if err != nil {
		return err
	}
	defer out.Close()

	// request
	getSchemaValues([]string{"media-type", "request body"}, nil).generate(out)
	getSchemaValues([]string{"property", "media-type", "request body"}, nil).generate(out)
	getSchemaValues([]string{"request parameter"}, []bool{true}).generate(out)

	// response
	getSchemaValues([]string{"media-type", "response"}, nil).generate(out)
	getSchemaValues([]string{"property", "media-type", "response"}, nil).generate(out)
	// getSchemaValues([]string{"request parameter"}, []bool{true}).generate(out)

	return nil
}

func getSchemaValues(hierarchy []string, attributed []bool) valueSets {
	return []valueSet{
		{
			adjective:  "value",
			hierarchy:  hierarchy,
			attributed: attributed,
			nouns:      []string{"max", "maxLength", "min", "minLength", "minItems", "maxItems"},
			actions:    []string{"set", "increase", "decrease"},
		},
		{
			adjective: "",
			hierarchy: hierarchy,
			nouns:     []string{"type/format"},
			actions:   []string{"change", "generalize"},
		},
		{
			adjective: "",
			hierarchy: hierarchy,
			nouns:     []string{"anyOf", "oneOf", "allOf"},
			actions:   []string{"add", "remove"},
		},
	}
}

func generateId(hierarchy []string, noun, action string) string {
	if before, _, found := strings.Cut(noun, "/"); found {
		noun = before
	}
	return strcase.ToKebab(concat(hierarchy) + "-" + noun + "-" + conjugate(action))
}

func concat(list []string) string {
	copy := slices.Clone(list)
	slices.Reverse(copy)
	return strings.Join(copy, "-")
}

func generateMessage(hierarchy []string, atttibuted []bool, noun, adjective, action string) string {
	return standardizeSpaces(fmt.Sprintf("%s %s of %s was %s", noun, adjective, getHierarchyMessage(hierarchy, atttibuted), getActionMessage(action)))
}

func getHierarchyMessage(hierarchy []string, atttibuted []bool) string {

	copy := slices.Clone(hierarchy)

	if atttibuted != nil {
		for i, s := range hierarchy {
			if atttibuted[i] {
				copy[i] = "%s " + s
			}
		}
	}

	result := strings.Join(copy, " %s of ")

	if !isTopLevel(hierarchy[len(hierarchy)-1]) {
		result += " %s"
	}

	return result
}

func isTopLevel(s string) bool {
	return s == "request body"
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func getActionMessage(action string) string {
	if isUnary(action) {
		return conjugate(action) + " to %s"
	}
	return conjugate(action) + " from %s to %s"
}

func isUnary(action string) bool {
	switch action {
	case "set":
		return true
	}
	return false
}

func conjugate(verb string) string {
	switch verb {
	case "set":
		return "set"
	case "add":
		return "added"
	}
	return verb + "d"
}
