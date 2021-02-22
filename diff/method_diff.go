package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// MethodDiff is the diff between two OAS operations (methods)
type MethodDiff struct {

	// ExtensionProps
	// Tags
	SummaryDiff     *ValueDiff      `json:"summary,omitempty"`     // diff of 'summary' property
	DescriptionDiff *ValueDiff      `json:"description,omitempty"` // diff of 'description' property
	OperationIDDiff *ValueDiff      `json:"operationID,omitempty"` // diff of 'operationID' property
	ParamDiff       *ParametersDiff `json:"parameters,omitempty"`  // diff of 'parameters' property
	// RequestBody
	ResponseDiff *ResponsesDiff `json:"responses,omitempty"` // diff of 'responses' property
	// Callbacks
	DeprecatedDiff *ValueDiff `json:"deprecated,omitempty"` // diff of 'deprecated' property
	// Security
	// Servers
	// ExternalDocs
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

	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiff(pathItem1.Description, pathItem2.Description)
	result.OperationIDDiff = getValueDiff(pathItem1.OperationID, pathItem2.OperationID)

	if diff := getParametersDiff(pathItem1.Parameters, pathItem2.Parameters); !diff.empty() {
		result.ParamDiff = diff
	}

	if diff := getResponsesDiff(pathItem1.Responses, pathItem2.Responses); !diff.empty() {
		result.ResponseDiff = diff
	}
	result.DeprecatedDiff = getValueDiff(pathItem1.Deprecated, pathItem2.Deprecated)

	return result
}
