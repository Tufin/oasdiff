package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// RequestBodyDiff describes the changes between a pair of request body objects: https://swagger.io/specification/#request-body-object
type RequestBodyDiff struct {
	Added           bool         `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted         bool         `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	DescriptionDiff *ValueDiff   `json:"description,omitempty" yaml:"description,omitempty"`
	ContentDiff     *ContentDiff `json:"content,omitempty" yaml:"content,omitempty"`
}

// Empty indicates whether a change was found in this element
func (requestBodyDiff *RequestBodyDiff) Empty() bool {
	if requestBodyDiff == nil {
		return true
	}
	return *requestBodyDiff == RequestBodyDiff{}
}

func newRequestBodyDiff() *RequestBodyDiff {
	return &RequestBodyDiff{}
}

func getRequestBodyDiff(config *Config, requestBodyRef1, requestBodyRef2 *openapi3.RequestBodyRef) (*RequestBodyDiff, error) {
	diff, err := getRequestBodyDiffInternal(config, requestBodyRef1, requestBodyRef2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}
	return diff, nil
}

func getRequestBodyDiffInternal(config *Config, requestBodyRef1, requestBodyRef2 *openapi3.RequestBodyRef) (*RequestBodyDiff, error) {

	if requestBodyRef1 == nil && requestBodyRef2 == nil {
		return nil, nil
	}

	result := newRequestBodyDiff()

	if requestBodyRef1 == nil && requestBodyRef2 != nil {
		result.Added = true
		return result, nil
	}

	if requestBodyRef1 != nil && requestBodyRef2 == nil {
		result.Deleted = true
		return result, nil
	}

	requestBody1, err := derefRequestBody(requestBodyRef1)
	if err != nil {
		return nil, err
	}
	requestBody2, err := derefRequestBody(requestBodyRef2)
	if err != nil {
		return nil, err
	}

	result.DescriptionDiff = getValueDiff(requestBody1.Description, requestBody2.Description)
	result.ContentDiff, err = getContentDiff(config, requestBody1.Content, requestBody2.Content)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func derefRequestBody(ref *openapi3.RequestBodyRef) (*openapi3.RequestBody, error) {

	if ref == nil || ref.Value == nil {
		return nil, fmt.Errorf("request body reference is nil")
	}

	return ref.Value, nil
}
