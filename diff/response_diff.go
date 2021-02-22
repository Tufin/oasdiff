package diff

import "github.com/getkin/kin-openapi/openapi3"

// ResponseDiff is a diff between two OAS responses
type ResponseDiff struct {
	// ExtensionProps
	DescriptionDiff *ValueDiff `json:"description,omitempty"`
	// Headers
	// Content
	// Links
}

func (responseDiff ResponseDiff) empty() bool {
	return responseDiff == ResponseDiff{}
}

func diffResponseValues(response1, response2 *openapi3.Response) ResponseDiff {
	result := ResponseDiff{}

	result.DescriptionDiff = getStringRefDiff(response1.Description, response2.Description)

	return result
}
