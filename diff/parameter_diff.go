package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ParameterDiff is a diff between parameter objects: https://swagger.io/specification/#parameter-object
type ParameterDiff struct {
	ExtensionProps      *ExtensionsDiff `json:"extensions,omitempty"`
	DescriptionDiff     *ValueDiff      `json:"description,omitempty"`
	StyleDiff           *ValueDiff      `json:"style,omitempty"`
	ExplodeDiff         *ValueDiff      `json:"explode,omitempty"`
	AllowEmptyValueDiff *ValueDiff      `json:"allowEmptyValue,omitempty"`
	AllowReservedDiff   *ValueDiff      `json:"allowReserved,omitempty"`
	DeprecatedDiff      *ValueDiff      `json:"deprecated,omitempty"`
	RequiredDiff        *ValueDiff      `json:"required,omitempty"`
	SchemaDiff          *SchemaDiff     `json:"schema,omitempty"`
	ExampleDiff         *ValueDiff      `json:"example,omitempty"`
	ContentDiff         *ContentDiff    `json:"content,omitempty"`
}

func (parameterDiff ParameterDiff) empty() bool {
	return parameterDiff == ParameterDiff{}
}

func getParameterDiff(param1, param2 *openapi3.Parameter) ParameterDiff {

	result := ParameterDiff{}

	if diff := getExtensionsDiff(param1.ExtensionProps, param2.ExtensionProps); !diff.empty() {
		result.ExtensionProps = diff
	}

	result.DescriptionDiff = getValueDiff(param1.Description, param2.Description)
	result.StyleDiff = getValueDiff(param1.Style, param2.Style)
	result.ExplodeDiff = getBoolRefDiff(param1.Explode, param2.Explode)
	result.AllowEmptyValueDiff = getValueDiff(param1.AllowEmptyValue, param2.AllowEmptyValue)
	result.AllowReservedDiff = getValueDiff(param1.AllowReserved, param2.AllowReserved)
	result.DeprecatedDiff = getValueDiff(param1.Deprecated, param2.Deprecated)
	result.RequiredDiff = getValueDiff(param1.Required, param2.Required)

	if schemaDiff := getSchemaDiff(param1.Schema, param2.Schema); !schemaDiff.empty() {
		result.SchemaDiff = &schemaDiff
	}

	result.ExampleDiff = getValueDiff(param1.Example, param2.Example)
	// Examples

	if contentDiff := getContentDiff(param1.Content, param2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
