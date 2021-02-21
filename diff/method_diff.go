package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// MethodDiff is the diff between two methods
type MethodDiff struct {
	ParamDiff    *ParametersDiff `json:"parameters,omitempty"`
	ResponseDiff *ResponsesDiff  `json:"responses,omitempty"`
}

func newMethodDiff() *MethodDiff {
	return &MethodDiff{}
}

func (methodDiff *MethodDiff) empty() bool {
	return methodDiff.ParamDiff == nil &&
		methodDiff.ResponseDiff == nil
}

func getMethodDiff(pathItem1, pathItem2 *openapi3.Operation) *MethodDiff {

	result := newMethodDiff()

	if diff := getParamDiff(pathItem1.Parameters, pathItem2.Parameters); !diff.empty() {
		result.ParamDiff = diff
	}

	if diff := getResponseDiff(pathItem1.Responses, pathItem2.Responses); !diff.empty() {
		result.ResponseDiff = diff
	}

	return result
}
