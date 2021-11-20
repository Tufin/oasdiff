package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// EncodingDiff describes the changes betweena pair of encoding objects: https://swagger.io/specification/#encoding-object
type EncodingDiff struct {
	ExtensionsDiff    *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	ContentTypeDiff   *ValueDiff      `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	HeadersDiff       *HeadersDiff    `json:"headers,omitempty" yaml:"headers,omitempty"`
	StyleDiff         *ValueDiff      `json:"styleDiff,omitempty" yaml:"styleDiff,omitempty"`
	ExplodeDiff       *ValueDiff      `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReservedDiff *ValueDiff      `json:"allowReservedDiff,omitempty" yaml:"allowReservedDiff,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *EncodingDiff) Empty() bool {
	return diff == nil || *diff == EncodingDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff EncodingDiff) Breaking() bool {
	if diff.Empty() {
		return false
	}

	return diff.ContentTypeDiff.Breaking() ||
		diff.HeadersDiff.Breaking() ||
		diff.StyleDiff.Breaking() || // TODO: check only if request body media type is application/x-www-form-urlencoded
		diff.ExplodeDiff.Breaking() || // TODO: be more specific
		diff.AllowReservedDiff.Breaking() // TODO: be more specific
}

func getEncodingDiff(config *Config, value1, value2 *openapi3.Encoding) (*EncodingDiff, error) {
	diff, err := getEncodingDiffInternal(config, value1, value2)
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

func getEncodingDiffInternal(config *Config, value1, value2 *openapi3.Encoding) (*EncodingDiff, error) {
	result := EncodingDiff{}
	var err error

	if value1 == nil || value2 == nil {
		return nil, fmt.Errorf("encoding is nil")
	}

	result.ExtensionsDiff = getExtensionsDiff(config, value1.ExtensionProps, value2.ExtensionProps)
	result.ContentTypeDiff = getValueDiff(config, true, value1.ContentType, value2.ContentType)
	result.HeadersDiff, err = getHeadersDiff(config, value1.Headers, value2.Headers)
	if err != nil {
		return nil, err
	}
	result.StyleDiff = getValueDiff(config, true, value1.Style, value2.Style)
	result.ExplodeDiff = getBoolRefDiff(config, true, value1.Explode, value2.Explode)
	result.AllowReservedDiff = getValueDiff(config, true, value1.AllowReserved, value2.AllowReserved)

	return &result, nil
}
