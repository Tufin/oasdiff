package generator

import (
	"fmt"
	"io"
	"strings"
)

type ValueSets []IValueSet

func NewValueSets(hierarchy []string, attributed []bool, valueSets ValueSets) ValueSets {

	result := make(ValueSets, len(valueSets))

	if attributed == nil {
		attributed = make([]bool, len(hierarchy))
	}

	for i, vs := range valueSets {
		result[i] = vs.setHierarchy(hierarchy, attributed)
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
	setHierarchy(hierarchy []string, attributed []bool) IValueSet
}

type AdjectiveType bool

const (
	PREDICATIVE AdjectiveType = false // PREDICATIVE adjectives are added after the noun (default)
	ATTRIBUTIVE AdjectiveType = true  // ATTRIBUTIVE adjectives are added before the noun
)

type ValueSet struct {
	adjective     string // adjective is added to the noun
	adjectiveType AdjectiveType
	hierarchy     []string
	attributed    []bool // attributed levels in the hierarchy are preceded by a name (%s)
	nouns         []string
	actions       []string
}

func (v ValueSet) setHierarchy(hierarchy []string, attributed []bool) ValueSet {
	if len(hierarchy) == 0 {
		return v
	}

	v.hierarchy = append(v.hierarchy, hierarchy...)
	v.attributed = append(v.attributed, attributed...)

	return v
}

// ValueSetA messages start with the noun
type ValueSetA ValueSet

func (v ValueSetA) setHierarchy(hierarchy []string, attributed []bool) IValueSet {
	return ValueSetA(ValueSet(v).setHierarchy(hierarchy, attributed))
}

func (v ValueSetA) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, atttibuted []bool, noun, adjective, action string) string {
		return standardizeSpaces(fmt.Sprintf("%s of %s was %s", addAttribute(noun, adjective, v.adjectiveType), getHierarchyMessage(hierarchy, atttibuted), getActionMessage(action)))
	}

	for _, noun := range v.nouns {
		for _, action := range v.actions {
			fmt.Fprintln(out, fmt.Sprintf("%s: %s", generateId(v.hierarchy, noun, action), generateMessage(v.hierarchy, v.attributed, noun, v.adjective, action)))
		}
	}
}

// ValueSetB messages start with the action
type ValueSetB ValueSet

func (v ValueSetB) setHierarchy(hierarchy []string, attributed []bool) IValueSet {
	return ValueSetB(ValueSet(v).setHierarchy(hierarchy, attributed))
}

func (v ValueSetB) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, atttibuted []bool, noun, adjective, action string) string {
		return standardizeSpaces(strings.Join([]string{conjugate(action), addAttribute(noun, adjective, v.adjectiveType), getHierarchyPostfix(action, hierarchy, atttibuted)}, " "))
	}

	for _, noun := range v.nouns {
		for _, action := range v.actions {
			fmt.Fprintln(out, fmt.Sprintf("%s: %s", generateId(v.hierarchy, noun, action), generateMessage(v.hierarchy, v.attributed, noun, v.adjective, action)))
		}
	}
}
