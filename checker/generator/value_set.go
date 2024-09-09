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

type ValueSetList []ValueSet

type ValueSet struct {
	AttributiveAdjective string // attributive adjectives are added before the object
	PredicativeAdjective string // predicative adjectives are added after the object
	Hierarchy            []string
	Names                []string
	Actions              []string
	Adverbs              []string
}

func (v ValueSet) setHierarchy(hierarchy []string) ValueSet {
	if len(hierarchy) > 0 {
		v.Hierarchy = append(v.Hierarchy, hierarchy...)
	}

	return v
}

// ValueSetA messages start with the object
// for example: "api was removed without deprecation"
type ValueSetA ValueSet

func (v ValueSetA) setHierarchy(hierarchy []string) IValueSet {
	return ValueSetA(ValueSet(v).setHierarchy(hierarchy))
}

func (v ValueSetA) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, object, attributiveAdjective, predicativeAdjective, action, adverb string) string {
		prefix := addAttribute(object, attributiveAdjective, predicativeAdjective)
		if hierarchyMessage := getHierarchyMessage(hierarchy); hierarchyMessage != "" {
			prefix += " of " + hierarchyMessage
		}

		return standardizeSpaces(fmt.Sprintf("%s was %s %s %s", prefix, conjugate(action), getActionMessage(action), adverb))
	}

	for _, object := range v.Names {
		for _, action := range v.Actions {
			id := generateId(v.Hierarchy, object, action)

			adverbs := v.Adverbs
			if v.Adverbs == nil {
				adverbs = []string{""}
			}
			for _, adverb := range adverbs {
				message := generateMessage(v.Hierarchy, object, v.AttributiveAdjective, v.PredicativeAdjective, action, adverb)
				fmt.Fprintf(out, "%s: %s\n", id, message)
			}
		}
	}
}

// ValueSetB messages start with the action
// for example: "removed %s request parameter %s"
type ValueSetB ValueSet

func (v ValueSetB) setHierarchy(hierarchy []string) IValueSet {
	return ValueSetB(ValueSet(v).setHierarchy(hierarchy))
}

func (v ValueSetB) generate(out io.Writer) {
	generateMessage := func(hierarchy []string, object, attributiveAdjective, predicativeAdjective, action string) string {
		return standardizeSpaces(strings.Join([]string{conjugate(action), addAttribute(object, attributiveAdjective, predicativeAdjective), getHierarchyPostfix(action, hierarchy)}, " "))
	}

	for _, object := range v.Names {
		for _, action := range v.Actions {
			fmt.Fprintf(out, "%s: %s\n", generateId(v.Hierarchy, object, action), generateMessage(v.Hierarchy, object, v.AttributiveAdjective, v.PredicativeAdjective, action))
		}
	}
}
