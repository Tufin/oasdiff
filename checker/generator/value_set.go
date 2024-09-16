package generator

import (
	"fmt"
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

func (vs ValueSets) generate() []string {
	result := []string{}
	for _, v := range vs {
		result = append(result, v.generate()...)
	}
	return result
}

type IValueSet interface {
	generate() []string
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

func (v ValueSetA) generate() []string {
	generateMessage := func(hierarchy []string, name, attributiveAdjective, predicativeAdjective, action, adverb string) string {
		prefix := addAttribute(name, attributiveAdjective, predicativeAdjective)
		if hierarchyMessage := getHierarchyMessage(hierarchy); hierarchyMessage != "" {
			prefix += " of " + hierarchyMessage
		}

		return standardizeSpaces(fmt.Sprintf("%s was %s %s %s", prefix, conjugate(action), getActionMessage(action), adverb))
	}

	result := []string{}
	for _, name := range v.Names {
		for _, action := range v.Actions {
			for _, adverb := range oneAtLeast(v.Adverbs) {
				id := generateId(v.Hierarchy, name, action, adverb)
				message := generateMessage(v.Hierarchy, name, v.AttributiveAdjective, v.PredicativeAdjective, action, adverb)
				result = append(result, fmt.Sprintf("%s: %s", id, message))
			}
		}
	}
	return result
}

// ValueSetB messages start with the action
// for example: "removed %s request parameter %s"
type ValueSetB ValueSet

func (v ValueSetB) setHierarchy(hierarchy []string) IValueSet {
	return ValueSetB(ValueSet(v).setHierarchy(hierarchy))
}

func (v ValueSetB) generate() []string {
	generateMessage := func(hierarchy []string, name, attributiveAdjective, predicativeAdjective, action string) string {
		return standardizeSpaces(strings.Join([]string{conjugate(action), addAttribute(name, attributiveAdjective, predicativeAdjective), getHierarchyPostfix(action, hierarchy)}, " "))
	}

	result := []string{}
	for _, name := range v.Names {
		for _, action := range v.Actions {
			result = append(result, fmt.Sprintf("%s: %s", generateId(v.Hierarchy, name, action, ""), generateMessage(v.Hierarchy, name, v.AttributiveAdjective, v.PredicativeAdjective, action)))
		}
	}

	return result
}

func oneAtLeast(list []string) []string {
	if len(list) == 0 {
		return []string{""}
	}
	return list
}
