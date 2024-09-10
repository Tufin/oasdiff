package generator

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type ChangeMap map[string]Changes

type Changes struct {
	ExcludeFromHierarchy bool               `yaml:"excludeFromHierarchy"`
	Actions              map[string]Objects `yaml:"actions"`
	NextLevel            ChangeMap          `yaml:"nextLevel"`
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
	fillHierarchy(changeMap, nil)
	generateRecursive(changeMap, out)
}

func generateRecursive(changeMap ChangeMap, out io.Writer) {
	for _, change := range changeMap {
		for action, objects := range change.Actions {
			for _, object := range objects {
				getValueSet(object, action).generate(out)
			}
		}
		generateRecursive(change.NextLevel, out)
	}
}

func fillHierarchy(changeMap ChangeMap, hierarchy []string) {
	for container, change := range changeMap {
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
