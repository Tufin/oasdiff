package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// EncodingDiff is a diff between encoding objects: https://swagger.io/specification/#encoding-object
type EncodingDiff struct {
	ContentTypeDiff   *ValueDiff   `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	HeadersDiff       *HeadersDiff `json:"headers,omitempty" yaml:"headers,omitempty"`
	StyleDiff         *ValueDiff   `json:"styleDiff,omitempty" yaml:"styleDiff,omitempty"`
	ExplodeDiff       *ValueDiff   `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReservedDiff *ValueDiff   `json:"allowReservedDiff,omitempty" yaml:"allowReservedDiff,omitempty"`
}

// Empty return true if there is no diff
func (diff *EncodingDiff) Empty() bool {
	return diff == nil || *diff == EncodingDiff{}
}

func getEncodingDiff(config *Config, value1, value2 *openapi3.Encoding) *EncodingDiff {
	diff := getEncodingDiffInternal(config, value1, value2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getEncodingDiffInternal(config *Config, value1, value2 *openapi3.Encoding) *EncodingDiff {
	result := EncodingDiff{}

	result.ContentTypeDiff = getValueDiff(value1.ContentType, value2.ContentType)
	result.HeadersDiff = getHeadersDiff(config, value1.Headers, value2.Headers)
	result.StyleDiff = getValueDiff(value1.Style, value2.Style)
	result.ExplodeDiff = getBoolRefDiff(value1.Explode, value2.Explode)
	result.AllowReservedDiff = getValueDiff(value1.AllowReserved, value2.AllowReserved)

	return &result
}
