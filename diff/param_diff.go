package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ParamDiff struct {
	DescriptionDiff     *ValueDiff   `json:"descriptionDiff,omitempty"`
	StyleDiff           *ValueDiff   `json:"styleDiff,omitempty"`
	ExplodeDiff         *ValueDiff   `json:"explodeDiff,omitempty"`
	AllowEmptyValueDiff *ValueDiff   `json:"allow_empty_valueDiff,omitempty"`
	AllowReservedDiff   *ValueDiff   `json:"allow_reservedDiff,omitempty"`
	DeprecatedDiff      *ValueDiff   `json:"deprecatedDiff,omitempty"`
	RequiredDiff        *ValueDiff   `json:"requiredDiff,omitempty"`
	SchemaDiff          *SchemaDiff  `json:"schemaDiff,omitempty"`
	ContentDiff         *ContentDiff `json:"contentDiff,omitempty"`
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
		result.SchemaDiff = &schemaDiff
	}

	if contentDiff := diffContent(param1.Content, param2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
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
