package generator

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
)

type MessageGenerator interface {
	generate() []string
}

type Getter func() (MessageGenerator, error)

func Generate(getter Getter) ([]string, error) {
	data, err := getter()
	if err != nil {
		return nil, err
	}

	return data.generate(), nil
}

func isEmpty(s string) bool {
	return s == ""
}

func filterStrings(list []string, f func(string) bool) []string {
	var result []string
	for _, s := range list {
		if !f(s) {
			result = append(result, s)
		}
	}
	return result
}

func generateId(hierarchy []string, object, action, adverb string) string {
	if prefix, _, found := strings.Cut(object, "/"); found {
		object = prefix
	}

	return strcase.ToKebab(strings.Join(filterStrings([]string{concat(hierarchy), object, conjugate(action), adverb}, isEmpty), "-"))
}

func concat(list []string) string {
	if len(list) == 0 {
		return ""
	}

	copy := slices.Clone(list)
	slices.Reverse(copy)
	return strings.Join(copy, "-")
}

func getHierarchyPostfix(action string, hierarchy []string) string {
	if len(hierarchy) == 0 {
		return ""
	}

	return getPreposition(action) + " " + getHierarchyMessage(hierarchy)
}

func getHierarchyMessage(hierarchy []string) string {

	copy := slices.Clone(hierarchy)

	for i, s := range hierarchy {
		if isAtttibuted(s) {
			copy[i] = "%s " + s
		}
	}
	result := strings.Join(copy, " %s of ")

	if hierarchy != nil && !isTopLevel(hierarchy[len(hierarchy)-1]) {
		result += " %s"
	}

	return result
}

func isTopLevel(s string) bool {
	return s == "request body"
}

func isAtttibuted(s string) bool {
	return s == "request parameter"
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func getActionMessage(action string) string {
	switch getArity(action) {
	case 0:
		return ""
	case 1:
		return " to %s"
	case 2:
		return " from %s to %s"
	default:
		return ""
	}
}

func getArity(action string) int {
	switch action {
	case "add", "remove":
		return 0
	case "set":
		return 1
	}
	return 2
}

func conjugate(verb string) string {
	switch verb {
	case "set":
		return "set"
	case "add":
		return "added"
	case "fail to parse":
		return "failed to parse"
	}
	return verb + "d"
}

func getPreposition(action string) string {
	switch action {
	case "add":
		return "to"
	}
	return "from"
}

func addAttribute(name, attributiveAdjective, predicativeAdjective string) string {
	return strings.Join([]string{attributiveAdjective + " " + name + " " + predicativeAdjective}, " ")
}
