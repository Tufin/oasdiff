package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ParamDiff struct {
	DescriptionDiff     *ValueDiff  `json:"descriptionDiff,omitempty"`
	StyleDiff           *ValueDiff  `json:"styleDiff,omitempty"`
	ExplodeDiff         *ValueDiff  `json:"explodeDiff,omitempty"`
	AllowEmptyValueDiff *ValueDiff  `json:"allow_empty_valueDiff,omitempty"`
	AllowReservedDiff   *ValueDiff  `json:"allow_reservedDiff,omitempty"`
	DeprecatedDiff      *ValueDiff  `json:"deprecatedDiff,omitempty"`
	RequiredDiff        *ValueDiff  `json:"requiredDiff,omitempty"`
	ShcemaDiff          *SchemaDiff `json:"schemaDiff,omitempty"`
	ExampleDiff         *ValueDiff  `json:"exampleDiff,omitempty"`
	ExamplesDiff        *ValueDiff  `json:"examplesDiff,omitempty"`
	ContentDiff         *ValueDiff  `json:"contentDiff,omitempty"`
}

func (paramDiff ParamDiff) empty() bool {
	return paramDiff == ParamDiff{}
}

func diffParamValues(param1 *openapi3.Parameter, param2 *openapi3.Parameter) ParamDiff {

	result := ParamDiff{}

	// TODO: ExtensionProps

	result.DescriptionDiff = getValueDiff(param1.Description, param2.Description)
	result.StyleDiff = getValueDiff(param1.Style, param2.Style)

	if diffExplode(param1.Explode, param2.Explode) {
		result.ExplodeDiff = getValueDiff(param1.Explode, param2.Explode)
	}

	result.AllowEmptyValueDiff = getValueDiff(param1.AllowEmptyValue, param2.AllowEmptyValue)
	result.AllowReservedDiff = getValueDiff(param1.AllowReserved, param2.AllowReserved)
	result.DeprecatedDiff = getValueDiff(param1.Deprecated, param2.Deprecated)
	result.RequiredDiff = getValueDiff(param1.Required, param2.Required)

	if schemaDiff := diffSchema(param1.Schema, param2.Schema); !schemaDiff.empty() {
		result.ShcemaDiff = &schemaDiff
	}

	if diffExample(param1.Example, param2.Example) {
		result.ExampleDiff = getValueDiff(param1.Example, param2.Example)
	}

	if diffExamples(param1.Examples, param2.Examples) {
		result.ExamplesDiff = getValueDiff(param1.Examples, param2.Examples)
	}

	if diffContent(param1.Content, param2.Content) {
		result.ContentDiff = getValueDiff(param1.Content, param2.Content)
	}

	return result
}

func diffExplode(pExplode1 *bool, pExplode2 *bool) bool {
	explode1 := derefExplode(pExplode1)
	explode2 := derefExplode(pExplode2)

	return explode1 != explode2
}

func derefExplode(pExplode *bool) bool {
	if pExplode == nil {
		return false // this is the default value for explode
	}

	return *pExplode
}

func diffExample(example1 interface{}, example2 interface{}) bool {
	return false
}

func diffExamples(examples1 openapi3.Examples, examples2 openapi3.Examples) bool {
	return false
}

func diffContent(content1 openapi3.Content, content2 openapi3.Content) bool {
	return false
}
