package diff

import "github.com/getkin/kin-openapi/openapi3"

// ResponseDiff is a diff between two OAS responses
type ResponseDiff struct {
	// ExtensionProps
	DescriptionDiff *ValueDiff   `json:"description,omitempty"` // diff of 'description' property
	HeadersDiff     *HeadersDiff `json:"headers,omitempty"`     // diff of 'headers' property
	ContentDiff     *ContentDiff `json:"content,omitempty"`     // diff of 'content' property
	// Links
}

func (responseDiff ResponseDiff) empty() bool {
	return responseDiff == ResponseDiff{}
}

func diffResponseValues(response1, response2 *openapi3.Response) ResponseDiff {
	result := ResponseDiff{}

	result.DescriptionDiff = getStringRefDiff(response1.Description, response2.Description)

	if headersDiff := getHeadersDiff(response1.Headers, response2.Headers); !headersDiff.empty() {
		result.HeadersDiff = headersDiff
	}

	if contentDiff := getContentDiff(response1.Content, response2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
