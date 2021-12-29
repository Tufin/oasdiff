package diff

import "github.com/getkin/kin-openapi/openapi3"

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

func (diff *VariableDiff) removeNonBreaking() {

	if diff.Empty() {
		return
	}

	diff.ExtensionsDiff = nil
	diff.DescriptionDiff = nil
}

func getVariableDiff(config *Config, var1, var2 *openapi3.ServerVariable) *VariableDiff {
	diff := getVariableDiffInternal(config, var1, var2)

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil
	}

	return diff
}

func getVariableDiffInternal(config *Config, var1, var2 *openapi3.ServerVariable) *VariableDiff {
	result := VariableDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, var1.ExtensionProps, var2.ExtensionProps)
	result.EnumDiff = getStringsDiff(var1.Enum, var2.Enum)
	result.DefaultDiff = getValueDiff(var1.Default, var2.Default)
	result.DescriptionDiff = getValueDiffConditional(config.ExcludeDescription, var1.Description, var2.Description)

	return &result
}
