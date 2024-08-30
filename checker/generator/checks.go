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

func getPreposition(action string) string {
	switch action {
	case "add":
		return "to"
	}
	return "from"
}

func addAttribute(noun, adjective string, adjectiveType AdjectiveType) string {
	if adjectiveType == ATTRIBUTIVE {
		return adjective + " " + noun
	}
	return noun + " " + adjective
}
