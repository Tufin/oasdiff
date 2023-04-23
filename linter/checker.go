package linter

import (
	"sort"

	"github.com/tufin/oasdiff/load"
)

const (
	LEVEL_ERROR = 0
	LEVEL_WARN  = 1
)

type Check func(*load.OpenAPISpecInfo) []Error

type Error struct {
	Id      string `json:"id,omitempty" yaml:"id,omitempty"`
	Text    string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level   int    `json:"level" yaml:"level"`
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`
}

type Errors []Error

func (e Errors) Len() int {
	return len(e)
}

func (e Errors) Less(i, j int) bool {
	iv, jv := e[i], e[j]

	switch {
	case iv.Level != jv.Level:
		return iv.Level < jv.Level
	case iv.Path != jv.Path:
		return iv.Path < jv.Path
	case iv.Id != jv.Id:
		return iv.Id < jv.Id
	case iv.Text != jv.Text:
		return iv.Text < jv.Text
	default:
		return iv.Comment < jv.Comment
	}
}

func (e Errors) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func Run(config Config, spec *load.OpenAPISpecInfo) Errors {
	result := make(Errors, 0)

	if spec == nil {
		return result
	}

	for _, check := range config.Checks {
		errs := check(spec)
		result = append(result, errs...)
	}

	sort.Sort(result)
	return result
}
