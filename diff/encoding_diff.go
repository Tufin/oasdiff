package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// EncodingDiff is a diff between encoding objects: https://swagger.io/specification/#encoding-object
type EncodingDiff struct {
	ExtensionProps    *ExtensionsDiff `json:"extensions,omitempty"`
	ContentTypeDiff   *ValueDiff      `json:"contentType,omitempty"`
	HeadersDiff       *HeadersDiff    `json:"headers,omitempty"`
	StyleDiff         *ValueDiff      `json:"styleDiff,omitempty"`
	ExplodeDiff       *ValueDiff      `json:"explode,omitempty"`
	AllowReservedDiff *ValueDiff      `json:"allowReservedDiff,omitempty"`
}

func (diff EncodingDiff) empty() bool {
	return diff == EncodingDiff{}
}

func getEncodingDiff(value1, value2 *openapi3.Encoding) EncodingDiff {
	result := EncodingDiff{}

	if diff := getExtensionsDiff(value1.ExtensionProps, value2.ExtensionProps); !diff.empty() {
		result.ExtensionProps = diff
	}

	result.ContentTypeDiff = getValueDiff(value1.ContentType, value2.ContentType)

	if headersDiff := getHeadersDiff(value1.Headers, value2.Headers); !headersDiff.empty() {
		result.HeadersDiff = headersDiff
	}

	result.StyleDiff = getValueDiff(value1.Style, value2.Style)
	result.ExplodeDiff = getBoolRefDiff(value1.Explode, value2.Explode)
	result.AllowReservedDiff = getValueDiff(value1.AllowReserved, value2.AllowReserved)

	return result
}
