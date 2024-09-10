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
	Object               string   `yaml:"object"`
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

	for _, change := range changeMap {
		for action, objects := range change.Actions {
			for _, object := range objects {
				getValueSet(object, action).generate(out)
			}
		}
		change.NextLevel.generate(out)
	}
}

func fillHierarchy(changeMap ChangeMap, hierarchy []string) {
	for container, change := range changeMap {
		if !change.ExcludeFromHierarchy {
			hierarchy = append([]string{container}, hierarchy...)
		}
		for _, objects := range change.Actions {
			for _, object := range objects {
				object.Hierarchy = hierarchy
			}
		}
		fillHierarchy(change.NextLevel, hierarchy)
	}
}

func getValueSet(object *Object, action string) IValueSet {
	valueSet := ValueSet{
		AttributiveAdjective: object.AttributiveAdjective,
		PredicativeAdjective: object.PredicativeAdjective,
		Hierarchy:            object.Hierarchy,
		Names:                []string{object.Object},
		Actions:              parseAction(action),
		Adverbs:              object.Adverbs,
	}

	if object.StartWith == "object" {
		return ValueSetA(valueSet)
	}
	return ValueSetB(valueSet)
}

func parseAction(action string) []string {
	if before, after, found := strings.Cut(action, "/"); found {
		return []string{before, after}
	}
	return []string{action}
}
