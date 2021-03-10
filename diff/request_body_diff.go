package diff

import "github.com/getkin/kin-openapi/openapi3"

// RequestBodyDiff is a diff between request body objects: https://swagger.io/specification/#request-body-object
type RequestBodyDiff struct {
	DescriptionDiff *ValueDiff   `json:"description,omitempty" yaml:"description,omitempty"`
	ContentDiff     *ContentDiff `json:"content,omitempty" yaml:"content,omitempty"`
}

// Empty return true if there is no diff
func (requestBodyDiff *RequestBodyDiff) Empty() bool {
	if requestBodyDiff == nil {
		return true
	}
	return *requestBodyDiff == RequestBodyDiff{}
}

func newRequestBodyDiff() *RequestBodyDiff {
	return &RequestBodyDiff{}
}

func getRequestBodyDiff(config *Config, requestBodyRef1, requestBodyRef2 *openapi3.RequestBodyRef) *RequestBodyDiff {
	diff := getRequestBodyDiffInternal(config, requestBodyRef1, requestBodyRef2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getRequestBodyDiffInternal(config *Config, requestBodyRef1, requestBodyRef2 *openapi3.RequestBodyRef) *RequestBodyDiff {
	result := newRequestBodyDiff()

	if requestBodyRef1 == nil || requestBodyRef1.Value == nil {
		// TODO: handle added, deleted etc.
		return nil
	}

	requestBody1 := requestBodyRef1.Value
	requestBody2 := requestBodyRef2.Value

	result.DescriptionDiff = getValueDiff(requestBody1.Description, requestBody2.Description)
	result.ContentDiff = getContentDiff(config, requestBody1.Content, requestBody2.Content)

	return result
}
