package diff

import "github.com/getkin/kin-openapi/openapi3"

// ResponseDiff is a diff between two OAS responses
type ResponseDiff struct {
	// ExtensionProps
	DescriptionDiff *ValueDiff `json:"description,omitempty"`
	// Headers
	ContentDiff *ContentDiff `json:"content,omitempty"`
	// Links
}

func (responseDiff ResponseDiff) empty() bool {
	return responseDiff == ResponseDiff{}
}

func diffResponseValues(response1, response2 *openapi3.Response) ResponseDiff {
	result := ResponseDiff{}

	result.DescriptionDiff = getStringRefDiff(response1.Description, response2.Description)

	if contentDiff := getContentDiff(response1.Content, response2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
