package generator

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ChangeMap struct {
	Changes    map[string]Changes `yaml:"changes"`
	Components map[string]Changes `yaml:"components"`
}

type Changes struct {
	Ref                  string             `yaml:"$ref"`
	ExcludeFromHierarchy bool               `yaml:"excludeFromHierarchy"`
	Actions              map[string]Objects `yaml:"actions"`
	NextLevel            map[string]Changes `yaml:"nextLevel"`
}

type Objects []*Object

type Object struct {
	Hierarchy            []string `yaml:"hierarchy"`
	Objects              []string `yaml:"objects"`
	Adverbs              []string `yaml:"adverbs"`
	StartWith            string   `yaml:"startWith"`
	PredicativeAdjective string   `yaml:"predicativeAdjective"`
	AttributiveAdjective string   `yaml:"attributiveAdjective"`
}

func GetTree() (MessageGenerator, error) {
	yamlFile, err := os.ReadFile("tree.yaml")
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}

	var changeMap ChangeMap
	err = yaml.Unmarshal(yamlFile, &changeMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return changeMap, nil
}

func (changeMap ChangeMap) generate(out io.Writer) {
	resolveRefs(changeMap.Changes, changeMap.Components)
	fillHierarchy(changeMap.Changes, nil)
	generateRecursive(changeMap.Changes, out)
}

func resolveRefs(changes map[string]Changes, components map[string]Changes) {
	for _, change := range changes {
		if change.Ref != "" {
			changes[change.Ref] = components[change.Ref]
		}
		resolveRefs(change.NextLevel, components)
	}
}

func generateRecursive(changes map[string]Changes, out io.Writer) {
	for _, change := range changes {
		for action, objects := range change.Actions {
			for _, object := range objects {
				getValueSet(object, action).generate(out)
			}
		}
		generateRecursive(change.NextLevel, out)
	}
}

func fillHierarchy(changes map[string]Changes, hierarchy []string) {
	for container, change := range changes {
		containerHierarchy := getContainerHierarchy(container, change, hierarchy)
		for _, objects := range change.Actions {
			for _, object := range objects {
				object.Hierarchy = containerHierarchy
			}
		}
		fillHierarchy(change.NextLevel, containerHierarchy)
	}
}

func getContainerHierarchy(container string, change Changes, hierarchy []string) []string {
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
		Names:                object.Objects,
		Actions:              parseAction(action),
		Adverbs:              object.Adverbs,
	}

	if object.StartWith == "object" {
		return ValueSetA(valueSet)
	}
	return ValueSetB(valueSet)
}

func parseAction(action string) []string {
	return strings.Split(action, "/")
}
