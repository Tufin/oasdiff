package diff

import "github.com/getkin/kin-openapi/openapi3"

// ResponseDiff is a diff between response objects: https://swagger.io/specification/#response-object
type ResponseDiff struct {
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty"`
	HeadersDiff     *HeadersDiff    `json:"headers,omitempty"`
	ContentDiff     *ContentDiff    `json:"content,omitempty"`
	// Links
}

func (responseDiff ResponseDiff) empty() bool {
	return responseDiff == ResponseDiff{}
}

func diffResponseValues(config *Config, response1, response2 *openapi3.Response) ResponseDiff {
	result := ResponseDiff{}

	if diff := getExtensionsDiff(config, response1.ExtensionProps, response2.ExtensionProps); !diff.empty() {
		result.ExtensionsDiff = diff
	}

	result.DescriptionDiff = getStringRefDiff(response1.Description, response2.Description)

	if headersDiff := getHeadersDiff(config, response1.Headers, response2.Headers); !headersDiff.empty() {
		result.HeadersDiff = headersDiff
	}

	if contentDiff := getContentDiff(config, response1.Content, response2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
