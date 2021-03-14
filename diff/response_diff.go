package diff

import "github.com/getkin/kin-openapi/openapi3"

// ResponseDiff is a diff between response objects: https://swagger.io/specification/#response-object
type ResponseDiff struct {
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	HeadersDiff     *HeadersDiff    `json:"headers,omitempty" yaml:"headers,omitempty"`
	ContentDiff     *ContentDiff    `json:"content,omitempty" yaml:"content,omitempty"`
	// Links
}

// Empty return true if there is no diff
func (diff *ResponseDiff) Empty() bool {
	return diff == nil || *diff == ResponseDiff{}
}

func diffResponseValues(config *Config, response1, response2 *openapi3.Response) (*ResponseDiff, error) {
	diff, err := diffResponseValuesInternal(config, response1, response2)
	if err != nil {
		return nil, err
	}
	if diff.Empty() {
		return nil, nil
	}
	return diff, nil
}

func diffResponseValuesInternal(config *Config, response1, response2 *openapi3.Response) (*ResponseDiff, error) {
	result := ResponseDiff{}
	var err error

	result.ExtensionsDiff = getExtensionsDiff(config, response1.ExtensionProps, response2.ExtensionProps)
	result.DescriptionDiff = getStringRefDiff(response1.Description, response2.Description)
	result.HeadersDiff, err = getHeadersDiff(config, response1.Headers, response2.Headers)
	if err != nil {
		return nil, err
	}
	result.ContentDiff, err = getContentDiff(config, response1.Content, response2.Content)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
