package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type MethodDiff struct {
	ParamDiff *Params `json:"parameters,omitempty"`
}

func newMethodDiff() *MethodDiff {
	return &MethodDiff{}
}

func (methodDiff *MethodDiff) empty() bool {
	return methodDiff.ParamDiff == nil
}

func getMethodDiff(pathItem1, pathItem2 *openapi3.Operation) *MethodDiff {

	result := newMethodDiff()

	if diff := getParamDiff(pathItem1.Parameters, pathItem2.Parameters); !diff.empty() {
		result.ParamDiff = diff
	}

	return result
}
