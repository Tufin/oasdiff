package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// VariableDiff describes the changes between a pair of server variable objects: https://swagger.io/specification/#server-variable-object
type VariableDiff struct {
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	EnumDiff        *StringsDiff    `json:"enum,omitempty" yaml:"enum,omitempty"`
	DefaultDiff     *ValueDiff      `json:"default,omitempty" yaml:"default,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *VariableDiff) Empty() bool {
	return diff == nil || *diff == VariableDiff{}
}

func getVariableDiff(config *Config, state *state, var1, var2 *openapi3.ServerVariable) *VariableDiff {
	diff := getVariableDiffInternal(config, state, var1, var2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getVariableDiffInternal(config *Config, state *state, var1, var2 *openapi3.ServerVariable) *VariableDiff {
	result := VariableDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, state, var1.Extensions, var2.Extensions)
	result.EnumDiff = getStringsDiff(var1.Enum, var2.Enum)
	result.DefaultDiff = getValueDiff(var1.Default, var2.Default)
	result.DescriptionDiff = getValueDiffConditional(config.IsExcludeDescription(), var1.Description, var2.Description)

	return &result
}
