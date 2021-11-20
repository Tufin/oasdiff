package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ParameterDiff describes the changes between a pair of parameter objects: https://swagger.io/specification/#parameter-object
type ParameterDiff struct {
	ExtensionsDiff      *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	DescriptionDiff     *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	StyleDiff           *ValueDiff      `json:"style,omitempty" yaml:"style,omitempty"`
	ExplodeDiff         *ValueDiff      `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowEmptyValueDiff *ValueDiff      `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	AllowReservedDiff   *ValueDiff      `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	DeprecatedDiff      *ValueDiff      `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	RequiredDiff        *ValueDiff      `json:"required,omitempty" yaml:"required,omitempty"`
	SchemaDiff          *SchemaDiff     `json:"schema,omitempty" yaml:"schema,omitempty"`
	ExampleDiff         *ValueDiff      `json:"example,omitempty" yaml:"example,omitempty"`
	ExamplesDiff        *ExamplesDiff   `json:"examples,omitempty" yaml:"examples,omitempty"`
	ContentDiff         *ContentDiff    `json:"content,omitempty" yaml:"content,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *ParameterDiff) Empty() bool {
	return diff == nil || *diff == ParameterDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff *ParameterDiff) Breaking() bool {
	if diff.Empty() {
		return false
	}

	if !diff.AllowEmptyValueDiff.Empty() {
		if diff.AllowEmptyValueDiff.To == true {
			// TODO: return true only if this is a query param
			return true
		}
	}

	if !diff.RequiredDiff.Empty() {
		if diff.RequiredDiff.To == true {
			return true
		}
	}

	return diff.SchemaDiff.Breaking() ||
		diff.ContentDiff.Breaking()
}

func getParameterDiff(config *Config, param1, param2 *openapi3.Parameter) (*ParameterDiff, error) {
	diff, err := getParameterDiffInternal(config, param1, param2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil, nil
	}

	return diff, nil
}

func getParameterDiffInternal(config *Config, param1, param2 *openapi3.Parameter) (*ParameterDiff, error) {

	result := ParameterDiff{}
	var err error

	result.ExtensionsDiff = getExtensionsDiff(config, param1.ExtensionProps, param2.ExtensionProps)
	result.DescriptionDiff = getValueDiffConditional(config, false, config.ExcludeDescription, param1.Description, param2.Description)
	result.StyleDiff = getValueDiff(config, false, param1.Style, param2.Style)
	result.ExplodeDiff = getBoolRefDiff(config, false, param1.Explode, param2.Explode)
	result.AllowEmptyValueDiff = getValueDiff(config, false, param1.AllowEmptyValue, param2.AllowEmptyValue)
	result.AllowReservedDiff = getValueDiff(config, false, param1.AllowReserved, param2.AllowReserved)
	result.DeprecatedDiff = getValueDiff(config, false, param1.Deprecated, param2.Deprecated)
	result.RequiredDiff = getValueDiff(config, false, param1.Required, param2.Required)
	result.SchemaDiff, err = getSchemaDiff(config, param1.Schema, param2.Schema)
	if err != nil {
		return nil, err
	}

	result.ExampleDiff = getValueDiffConditional(config, false, config.ExcludeExamples, param1.Example, param2.Example)

	result.ExamplesDiff, err = getExamplesDiff(config, param1.Examples, param2.Examples)
	if err != nil {
		return nil, err
	}

	result.ContentDiff, err = getContentDiff(config, param1.Content, param2.Content)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Patch applies the patch to a parameter
func (diff *ParameterDiff) Patch(parameter *openapi3.Parameter) error {
	diff.DescriptionDiff.patchString(&parameter.Description)

	schema, err := derefSchema(parameter.Schema)
	if err != nil {
		// no schema to patch, continue.
		return nil
	}

	return diff.SchemaDiff.Patch(schema)
}
