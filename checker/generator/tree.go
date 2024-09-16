package generator

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ChangeTree struct {
	Changes    ChangeMap `yaml:"changes"`
	Components ChangeMap `yaml:"components"`
}

type ChangeMap map[string]*Changes

type Changes struct {
	Ref                  string              `yaml:"$ref"`
	ExcludeFromHierarchy bool                `yaml:"excludeFromHierarchy"`
	Actions              map[string]*Objects `yaml:"actions"`
	NextLevel            ChangeMap           `yaml:"nextLevel"`
}

type Objects []*Object

type Object struct {
	Hierarchy            []string `yaml:"hierarchy"`
	Names                []string `yaml:"names"`
	Adverbs              []string `yaml:"adverbs"`
	StartWithName        bool     `yaml:"startWith"`
	PredicativeAdjective string   `yaml:"predicativeAdjective"`
	AttributiveAdjective string   `yaml:"attributiveAdjective"`
}

func GetTree(file string) func() (MessageGenerator, error) {
	return func() (MessageGenerator, error) {
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
		}

		var changeMap ChangeTree
		err = yaml.Unmarshal(yamlFile, &changeMap)
		if err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		return changeMap, nil
	}
}

func (changeMap ChangeTree) generate(out io.Writer) {
	resolveRefs(changeMap.Changes, changeMap.Components)
	fillHierarchy(changeMap.Changes, nil)
	generateRecursive(changeMap.Changes, out)
}

func resolveRefs(changes ChangeMap, components ChangeMap) {
	for container, change := range changes {
		if change.Ref != "" {
			changes[container] = components[change.Ref]
		}
		resolveRefs(changes[container].NextLevel, components)
	}
}

func generateRecursive(changes ChangeMap, out io.Writer) {
	for _, change := range changes {
		for action, objects := range change.Actions {
			for _, object := range *objects {
				getValueSet(object, action).generate(out)
			}
		}
		generateRecursive(change.NextLevel, out)
	}
}

func fillHierarchy(changes ChangeMap, hierarchy []string) {
	for container, change := range changes {
		containerHierarchy := getContainerHierarchy(container, change, hierarchy)
		for _, objects := range change.Actions {
			for _, object := range *objects {
				object.Hierarchy = containerHierarchy
			}
		}
		fillHierarchy(change.NextLevel, containerHierarchy)
	}
}

func getContainerHierarchy(container string, change *Changes, hierarchy []string) []string {
	if change.ExcludeFromHierarchy {
		return hierarchy
	}
	return append([]string{container}, hierarchy...)
}

func getValueSet(object *Object, action string) IValueSet {
	valueSet := ValueSet{
		AttributiveAdjective: object.AttributiveAdjective,
		PredicativeAdjective: object.PredicativeAdjective,
		Hierarchy:            object.Hierarchy,
		Names:                object.Names,
		Actions:              parseAction(action),
		Adverbs:              object.Adverbs,
	}

	if object.StartWithName {
		return ValueSetA(valueSet)
	}
	return ValueSetB(valueSet)
}

func parseAction(action string) []string {
	return strings.Split(action, "/")
}
