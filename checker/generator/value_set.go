package generator

import (
	"fmt"
	"io"
)

type ValueSets []ValueSet

func (vs ValueSets) generate(out io.Writer) {
	for _, v := range vs {
		v.generate(out)
	}
}

type ValueSet interface {
	generate(out io.Writer)
}

type ValueSetA struct {
	adjective  string
	hierarchy  []string
	attributed []bool
	nouns      []string
	actions    []string
}

// ValueSetA messages start with the noun
func (v ValueSetA) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, atttibuted []bool, noun, adjective, action string) string {
		return standardizeSpaces(fmt.Sprintf("%s %s of %s was %s", noun, adjective, getHierarchyMessage(hierarchy, atttibuted), getActionMessage(action)))
	}

	for _, noun := range v.nouns {
		for _, action := range v.actions {
			fmt.Fprintln(out, fmt.Sprintf("%s: %s", generateId(v.hierarchy, noun, action), generateMessage(v.hierarchy, v.attributed, noun, v.adjective, action)))
		}
	}
}

// ValueSetB messages start with the action
type ValueSetB struct {
	adjective  string
	hierarchy  []string
	attributed []bool
	nouns      []string
	actions    []string
}

func (v ValueSetB) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, atttibuted []bool, noun, _, action string) string {
		return standardizeSpaces(fmt.Sprintf("%s %s %s %s", conjugate(action), noun, getPreposition(action), getHierarchyMessage(hierarchy, atttibuted)))
	}

	for _, noun := range v.nouns {
		for _, action := range v.actions {
			fmt.Fprintln(out, fmt.Sprintf("%s: %s", generateId(v.hierarchy, noun, action), generateMessage(v.hierarchy, v.attributed, noun, v.adjective, action)))
		}
	}
}
