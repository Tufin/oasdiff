package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ParameterDiff is a diff between two OAS parameters
type ParameterDiff struct {
	DescriptionDiff     *ValueDiff   `json:"description,omitempty"`     // diff of 'description' property
	StyleDiff           *ValueDiff   `json:"style,omitempty"`           // diff of 'style' property
	ExplodeDiff         *ValueDiff   `json:"explode,omitempty"`         // diff of 'explode' property
	AllowEmptyValueDiff *ValueDiff   `json:"allowEmptyValue,omitempty"` // diff of 'allowEmptyValue' property
	AllowReservedDiff   *ValueDiff   `json:"allowReserved,omitempty"`   // diff of 'allowReserved' property
	DeprecatedDiff      *ValueDiff   `json:"deprecated,omitempty"`      // diff of 'deprecated' property
	RequiredDiff        *ValueDiff   `json:"required,omitempty"`        // diff of 'required' property
	SchemaDiff          *SchemaDiff  `json:"schema,omitempty"`          // diff of 'schema' property
	ExampleDiff         *ValueDiff   `json:"example,omitempty"`         // diff of 'example' property
	ContentDiff         *ContentDiff `json:"content,omitempty"`         // diff of 'content' property
}

func (parameterDiff ParameterDiff) empty() bool {
	return parameterDiff == ParameterDiff{}
}

func getParameterDiff(param1, param2 *openapi3.Parameter) ParameterDiff {

	result := ParameterDiff{}

	// TODO: ExtensionProps

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

	if contentDiff := getContentDiff(param1.Content, param2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
