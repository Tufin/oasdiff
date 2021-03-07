package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ParameterDiff is a diff between parameter objects: https://swagger.io/specification/#parameter-object
type ParameterDiff struct {
	ExtensionsDiff      *ExtensionsDiff `json:"extensions,omitempty"`
	DescriptionDiff     *ValueDiff      `json:"description,omitempty"`
	StyleDiff           *ValueDiff      `json:"style,omitempty"`
	ExplodeDiff         *ValueDiff      `json:"explode,omitempty"`
	AllowEmptyValueDiff *ValueDiff      `json:"allowEmptyValue,omitempty"`
	AllowReservedDiff   *ValueDiff      `json:"allowReserved,omitempty"`
	DeprecatedDiff      *ValueDiff      `json:"deprecated,omitempty"`
	RequiredDiff        *ValueDiff      `json:"required,omitempty"`
	SchemaDiff          *SchemaDiff     `json:"schema,omitempty"`
	ExampleDiff         *ValueDiff      `json:"example,omitempty"`
	// Examples
	ContentDiff *ContentDiff `json:"content,omitempty"`
}

func (diff *ParameterDiff) empty() bool {
	return diff == nil || *diff == ParameterDiff{}
}

func getParameterDiff(config *Config, param1, param2 *openapi3.Parameter) *ParameterDiff {
	diff := getParameterDiffInternal(config, param1, param2)
	if diff.empty() {
		return nil
	}
	return diff
}

func getParameterDiffInternal(config *Config, param1, param2 *openapi3.Parameter) *ParameterDiff {

	result := ParameterDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, param1.ExtensionProps, param2.ExtensionProps)
	result.DescriptionDiff = getValueDiff(param1.Description, param2.Description)
	result.StyleDiff = getValueDiff(param1.Style, param2.Style)
	result.ExplodeDiff = getBoolRefDiff(param1.Explode, param2.Explode)
	result.AllowEmptyValueDiff = getValueDiff(param1.AllowEmptyValue, param2.AllowEmptyValue)
	result.AllowReservedDiff = getValueDiff(param1.AllowReserved, param2.AllowReserved)
	result.DeprecatedDiff = getValueDiff(param1.Deprecated, param2.Deprecated)
	result.RequiredDiff = getValueDiff(param1.Required, param2.Required)
	result.SchemaDiff = getSchemaDiff(config, param1.Schema, param2.Schema)

	if config.IncludeExamples {
		result.ExampleDiff = getValueDiff(param1.Example, param2.Example)
	}

	result.ContentDiff = getContentDiff(config, param1.Content, param2.Content)

	return &result
}
