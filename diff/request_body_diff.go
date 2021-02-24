package diff

import "github.com/getkin/kin-openapi/openapi3"

// RequestBodyDiff is a diff between request body objects: https://swagger.io/specification/#request-body-object
type RequestBodyDiff struct {
	DescriptionDiff *ValueDiff   `json:"description,omitempty"`
	ContentDiff     *ContentDiff `json:"content,omitempty"`
}

func (requestBodyDiff RequestBodyDiff) empty() bool {
	return requestBodyDiff == RequestBodyDiff{}
}

func newRequestBodyDiff() *RequestBodyDiff {
	return &RequestBodyDiff{}
}

func getRequestBodyDiff(requestBodyRef1, requestBodyRef2 *openapi3.RequestBodyRef) *RequestBodyDiff {
	result := newRequestBodyDiff()

	if requestBodyRef1 == nil || requestBodyRef1.Value == nil {
		// TODO: handle added, deleted etc.
		return result
	}

	requestBody1 := requestBodyRef1.Value
	requestBody2 := requestBodyRef2.Value

	result.DescriptionDiff = getValueDiff(requestBody1.Description, requestBody2.Description)

	if contentDiff := getContentDiff(requestBody1.Content, requestBody2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
