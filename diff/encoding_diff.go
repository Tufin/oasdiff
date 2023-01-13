package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// EncodingDiff describes the changes between a pair of encoding objects: https://swagger.io/specification/#encoding-object
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

func (diff *EncodingDiff) removeNonBreaking() {

	if diff.Empty() {
		return
	}

	diff.ExtensionsDiff = nil

	// TODO: if request body media type isn't application/x-www-form-urlencoded => diff.StyleDiff = nil
	// TODO: diff.ExplodeDiff is non breaking in specific cases
	// TODO: diff.AllowReservedDiff is non breaking in specific cases
}

func getEncodingDiff(config *Config, state *state, value1, value2 *openapi3.Encoding) (*EncodingDiff, error) {
	diff, err := getEncodingDiffInternal(config, state, value1, value2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getEncodingDiffInternal(config *Config, state *state, value1, value2 *openapi3.Encoding) (*EncodingDiff, error) {
	result := EncodingDiff{}
	var err error

	if value1 == nil || value2 == nil {
		return nil, fmt.Errorf("encoding is nil")
	}

	result.ExtensionsDiff = getExtensionsDiff(config, state, value1.Extensions, value2.Extensions)
	result.ContentTypeDiff = getValueDiff(value1.ContentType, value2.ContentType)
	result.HeadersDiff, err = getHeadersDiff(config, state, value1.Headers, value2.Headers)
	if err != nil {
		return nil, err
	}
	result.StyleDiff = getValueDiff(value1.Style, value2.Style)
	result.ExplodeDiff = getBoolRefDiff(value1.Explode, value2.Explode)
	result.AllowReservedDiff = getValueDiff(value1.AllowReserved, value2.AllowReserved)

	return &result, nil
}
