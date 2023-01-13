package diff

import "github.com/getkin/kin-openapi/openapi3"

// ResponseDiff describes the changes between a pair of response objects: https://swagger.io/specification/#response-object
type ResponseDiff struct {
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	HeadersDiff     *HeadersDiff    `json:"headers,omitempty" yaml:"headers,omitempty"`
	ContentDiff     *ContentDiff    `json:"content,omitempty" yaml:"content,omitempty"`
	LinksDiff       *LinksDiff      `json:"links,omitempty" yaml:"links,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *ResponseDiff) Empty() bool {
	return diff == nil || *diff == ResponseDiff{}
}

func (diff *ResponseDiff) removeNonBreaking() {

	if diff.Empty() {
		return
	}

	diff.ExtensionsDiff = nil
	diff.DescriptionDiff = nil
}

func diffResponseValues(config *Config, state *state, response1, response2 *openapi3.Response) (*ResponseDiff, error) {
	diff, err := diffResponseValuesInternal(config, state, response1, response2)
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

func diffResponseValuesInternal(config *Config, state *state, response1, response2 *openapi3.Response) (*ResponseDiff, error) {
	result := ResponseDiff{}
	var err error

	result.ExtensionsDiff = getExtensionsDiff(config, state, response1.Extensions, response2.Extensions)
	result.DescriptionDiff = getStringRefDiffConditional(config.ExcludeDescription, response1.Description, response2.Description)
	result.HeadersDiff, err = getHeadersDiff(config, state, response1.Headers, response2.Headers)
	if err != nil {
		return nil, err
	}

	result.ContentDiff, err = getContentDiff(config, state, response1.Content, response2.Content)
	if err != nil {
		return nil, err
	}

	result.LinksDiff, err = getLinksDiff(config, state, response1.Links, response2.Links)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
