package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// MethodDiff is the diff between two operation objects: https://swagger.io/specification/#operation-object
type MethodDiff struct {
	ExtensionsDiff  *ExtensionsDiff           `json:"extensions,omitempty"`
	TagsDiff        *StringsDiff              `json:"tags,omitempty"`
	SummaryDiff     *ValueDiff                `json:"summary,omitempty"`
	DescriptionDiff *ValueDiff                `json:"description,omitempty"`
	OperationIDDiff *ValueDiff                `json:"operationID,omitempty"`
	ParametersDiff  *ParametersDiff           `json:"parameters,omitempty"`
	RequestBodyDiff *RequestBodyDiff          `json:"requestBody,omitempty"`
	ResponsesDiff   *ResponsesDiff            `json:"responses,omitempty"`
	CallbacksDiff   *CallbacksDiff            `json:"callbacks,omitempty"`
	DeprecatedDiff  *ValueDiff                `json:"deprecated,omitempty"`
	SecurityDiff    *SecurityRequirementsDiff `json:"securityRequirements,omitempty"`
	ServersDiff     *ServersDiff              `json:"servers,omitempty"`
	// ExternalDocs
}

func newMethodDiff() *MethodDiff {
	return &MethodDiff{}
}

func (methodDiff *MethodDiff) empty() bool {
	if methodDiff == nil {
		return true
	}

	return *methodDiff == MethodDiff{}
}

func getMethodDiff(config *Config, pathItem1, pathItem2 *openapi3.Operation) *MethodDiff {
	diff := getMethodDiffInternal(config, pathItem1, pathItem2)
	if diff.empty() {
		return nil
	}
	return diff
}

func getMethodDiffInternal(config *Config, pathItem1, pathItem2 *openapi3.Operation) *MethodDiff {

	result := newMethodDiff()

	result.ExtensionsDiff = getExtensionsDiff(config, pathItem1.ExtensionProps, pathItem2.ExtensionProps)
	result.TagsDiff = getStringsDiff(pathItem1.Tags, pathItem2.Tags)
	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiff(pathItem1.Description, pathItem2.Description)
	result.OperationIDDiff = getValueDiff(pathItem1.OperationID, pathItem2.OperationID)
	result.ParametersDiff = getParametersDiff(config, pathItem1.Parameters, pathItem2.Parameters)
	result.RequestBodyDiff = getRequestBodyDiff(config, pathItem1.RequestBody, pathItem2.RequestBody)
	result.ResponsesDiff = getResponsesDiff(config, pathItem1.Responses, pathItem2.Responses)
	result.CallbacksDiff = getCallbacksDiff(config, pathItem1.Callbacks, pathItem2.Callbacks)
	result.DeprecatedDiff = getValueDiff(pathItem1.Deprecated, pathItem2.Deprecated)
	result.SecurityDiff = getSecurityRequirementsDiff(config, pathItem1.Security, pathItem2.Security)
	result.ServersDiff = getServersDiff(config, pathItem1.Servers, pathItem2.Servers)

	return result
}
