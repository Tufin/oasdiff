package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// MethodDiff is the diff between two operation objects: https://swagger.io/specification/#operation-object
type MethodDiff struct {
	ExtensionProps  *ExtensionsDiff  `json:"extensions,omitempty"`
	TagsDiff        *StringsDiff     `json:"tags,omitempty"`
	SummaryDiff     *ValueDiff       `json:"summary,omitempty"`
	DescriptionDiff *ValueDiff       `json:"description,omitempty"`
	OperationIDDiff *ValueDiff       `json:"operationID,omitempty"`
	ParametersDiff  *ParametersDiff  `json:"parameters,omitempty"`
	RequestBodyDiff *RequestBodyDiff `json:"requestBody,omitempty"`
	ResponsesDiff   *ResponsesDiff   `json:"responses,omitempty"`
	CallbacksDiff   *CallbacksDiff   `json:"callbacks,omitempty"`
	DeprecatedDiff  *ValueDiff       `json:"deprecated,omitempty"`
	// Security
	ServersDiff *ServersDiff `json:"servers,omitempty"`
	// ExternalDocs
}

func newMethodDiff() *MethodDiff {
	return &MethodDiff{}
}

func (methodDiff *MethodDiff) empty() bool {
	return *methodDiff == MethodDiff{}
}

func getMethodDiff(config *Config, pathItem1, pathItem2 *openapi3.Operation) *MethodDiff {

	result := newMethodDiff()

	if diff := getExtensionsDiff(pathItem1.ExtensionProps, pathItem2.ExtensionProps); !diff.empty() {
		result.ExtensionProps = diff
	}

	result.TagsDiff = getStringsDiff(pathItem1.Tags, pathItem2.Tags)
	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiff(pathItem1.Description, pathItem2.Description)
	result.OperationIDDiff = getValueDiff(pathItem1.OperationID, pathItem2.OperationID)

	if diff := getParametersDiff(config, pathItem1.Parameters, pathItem2.Parameters); !diff.empty() {
		result.ParametersDiff = diff
	}

	if diff := getRequestBodyDiff(config, pathItem1.RequestBody, pathItem2.RequestBody); !diff.empty() {
		result.RequestBodyDiff = diff
	}

	if diff := getResponsesDiff(config, pathItem1.Responses, pathItem2.Responses); !diff.empty() {
		result.ResponsesDiff = diff
	}

	if diff := getCallbacksDiff(config, pathItem1.Callbacks, pathItem2.Callbacks); !diff.empty() {
		result.CallbacksDiff = diff
	}

	result.DeprecatedDiff = getValueDiff(pathItem1.Deprecated, pathItem2.Deprecated)

	if diff := getServersDiff(pathItem1.Servers, pathItem2.Servers); !diff.empty() {
		result.ServersDiff = diff
	}

	return result
}
