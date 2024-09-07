package generator

import (
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

	getAll().generate(out)

	return nil
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

func generateId(hierarchy []string, noun, action string) string {
	if before, _, found := strings.Cut(noun, "/"); found {
		noun = before
	}

	return strcase.ToKebab(strings.Join(filterStrings([]string{concat(hierarchy), noun, conjugate(action)}, isEmpty), "-"))
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
	return s == "request body" ||
		s == "paths"
}

func isAtttibuted(s string) bool {
	return s == "request parameter"
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
	case "become":
		return "became"
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

func addAttribute(noun, attributiveAdjective, predicativeAdjective string) string {
	return strings.Join([]string{attributiveAdjective + " " + noun + " " + predicativeAdjective}, " ")
}
