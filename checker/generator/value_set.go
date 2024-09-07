package generator

import (
	"fmt"
	"io"
	"strings"
)

type ValueSets []IValueSet

// NewValueSets creates a new ValueSets object
func NewValueSets(hierarchy []string, valueSets ValueSets) ValueSets {

	result := make(ValueSets, len(valueSets))

	for i, vs := range valueSets {
		result[i] = vs.setHierarchy(hierarchy)
	}

	return result
}

func (vs ValueSets) generate(out io.Writer) {
	for _, v := range vs {
		v.generate(out)
	}
}

type IValueSet interface {
	generate(out io.Writer)
	setHierarchy(hierarchy []string) IValueSet
}

type ValueSet struct {
	attributiveAdjective string // attributive adjectives are added before the noun
	predicativeAdjective string // predicative adjectives are added after the noun
	hierarchy            []string
	nouns                []string
	actions              []string
}

func (v ValueSet) setHierarchy(hierarchy []string) ValueSet {
	if len(hierarchy) == 0 {
		return v
	}

	v.hierarchy = append(v.hierarchy, hierarchy...)

	return v
}

// ValueSetA messages start with the noun
type ValueSetA ValueSet

func (v ValueSetA) setHierarchy(hierarchy []string) IValueSet {
	return ValueSetA(ValueSet(v).setHierarchy(hierarchy))
}

func (v ValueSetA) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, noun, attributiveAdjective, predicativeAdjective, action string) string {
		prefix := addAttribute(noun, attributiveAdjective, predicativeAdjective)
		if hierarchyMessage := getHierarchyMessage(hierarchy); hierarchyMessage != "" {
			prefix += " of " + hierarchyMessage
		}

		return standardizeSpaces(fmt.Sprintf("%s was %s", prefix, getActionMessage(action)))
	}

	for _, noun := range v.nouns {
		for _, action := range v.actions {
			id := generateId(v.hierarchy, noun, action)
			message := generateMessage(v.hierarchy, noun, v.attributiveAdjective, v.predicativeAdjective, action)
			fmt.Fprintln(out, fmt.Sprintf("%s: %s", id, message))
		}
	}
}

// ValueSetB messages start with the action
type ValueSetB ValueSet

func (v ValueSetB) setHierarchy(hierarchy []string) IValueSet {
	return ValueSetB(ValueSet(v).setHierarchy(hierarchy))
}

func (v ValueSetB) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, noun, attributiveAdjective, predicativeAdjective, action string) string {
		return standardizeSpaces(strings.Join([]string{conjugate(action), addAttribute(noun, attributiveAdjective, predicativeAdjective), getHierarchyPostfix(action, hierarchy)}, " "))
	}

	for _, noun := range v.nouns {
		for _, action := range v.actions {
			fmt.Fprintln(out, fmt.Sprintf("%s: %s", generateId(v.hierarchy, noun, action), generateMessage(v.hierarchy, noun, v.attributiveAdjective, v.predicativeAdjective, action)))
		}
	}
}
