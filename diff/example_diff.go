package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ExampleDiff describes the changes between a pair of example objects: https://swagger.io/specification/#example-object
type ExampleDiff struct {
	ExtensionsDiff    *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	SummaryDiff       *ValueDiff      `json:"summary,omitempty" yaml:"summary,omitempty"`
	DescriptionDiff   *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	ValueDiff         *ValueDiff      `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValueDiff *ValueDiff      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *ExampleDiff) Empty() bool {
	return diff == nil || *diff == ExampleDiff{}
}

func getExampleDiff(config *Config, state *state, value1, value2 *openapi3.Example) *ExampleDiff {
	diff := getExampleDiffInternal(config, state, value1, value2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getExampleDiffInternal(config *Config, state *state, value1, value2 *openapi3.Example) *ExampleDiff {
	result := ExampleDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, state, value1.Extensions, value2.Extensions)
	result.SummaryDiff = getValueDiff(value1.Summary, value2.Summary)
	result.DescriptionDiff = getValueDiffConditional(config.ExcludeDescription, value1.Description, value2.Description)
	result.ValueDiff = getValueDiff(value1.Value, value2.Value)
	result.ExternalValueDiff = getValueDiff(value1.ExternalValue, value2.ExternalValue)

	return &result
}
