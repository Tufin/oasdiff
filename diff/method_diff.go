package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// MethodDiff is the diff between two operation objects: https://swagger.io/specification/#operation-object
type MethodDiff struct {
	ExtensionsDiff  *ExtensionsDiff           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	TagsDiff        *StringsDiff              `json:"tags,omitempty" yaml:"tags,omitempty"`
	SummaryDiff     *ValueDiff                `json:"summary,omitempty" yaml:"summary,omitempty"`
	DescriptionDiff *ValueDiff                `json:"description,omitempty" yaml:"description,omitempty"`
	OperationIDDiff *ValueDiff                `json:"operationID,omitempty" yaml:"operationID,omitempty"`
	ParametersDiff  *ParametersDiff           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBodyDiff *RequestBodyDiff          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	ResponsesDiff   *ResponsesDiff            `json:"responses,omitempty" yaml:"responses,omitempty"`
	CallbacksDiff   *CallbacksDiff            `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	DeprecatedDiff  *ValueDiff                `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	SecurityDiff    *SecurityRequirementsDiff `json:"securityRequirements,omitempty" yaml:"securityRequirements,omitempty"`
	ServersDiff     *ServersDiff              `json:"servers,omitempty" yaml:"servers,omitempty"`
	// ExternalDocs
}

func newMethodDiff() *MethodDiff {
	return &MethodDiff{}
}

// Empty return true if there is no diff
func (methodDiff *MethodDiff) Empty() bool {
	if methodDiff == nil {
		return true
	}

	return *methodDiff == MethodDiff{}
}

func getMethodDiff(config *Config, pathItem1, pathItem2 *openapi3.Operation) (*MethodDiff, error) {
	diff, err := getMethodDiffInternal(config, pathItem1, pathItem2)

	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getMethodDiffInternal(config *Config, pathItem1, pathItem2 *openapi3.Operation) (*MethodDiff, error) {

	result := newMethodDiff()
	var err error

	result.ExtensionsDiff = getExtensionsDiff(config, pathItem1.ExtensionProps, pathItem2.ExtensionProps)
	result.TagsDiff = getStringsDiff(pathItem1.Tags, pathItem2.Tags)
	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiff(pathItem1.Description, pathItem2.Description)
	result.OperationIDDiff = getValueDiff(pathItem1.OperationID, pathItem2.OperationID)
	result.ParametersDiff, err = getParametersDiff(config, pathItem1.Parameters, pathItem2.Parameters)
	if err != nil {
		return nil, err
	}

	result.RequestBodyDiff, err = getRequestBodyDiff(config, pathItem1.RequestBody, pathItem2.RequestBody)
	if err != nil {
		return nil, err
	}

	result.ResponsesDiff, err = getResponsesDiff(config, pathItem1.Responses, pathItem2.Responses)
	if err != nil {
		return nil, err
	}

	result.CallbacksDiff, err = getCallbacksDiff(config, pathItem1.Callbacks, pathItem2.Callbacks)
	if err != nil {
		return nil, err
	}
	result.DeprecatedDiff = getValueDiff(pathItem1.Deprecated, pathItem2.Deprecated)
	result.SecurityDiff = getSecurityRequirementsDiff(config, pathItem1.Security, pathItem2.Security)
	result.ServersDiff = getServersDiff(config, pathItem1.Servers, pathItem2.Servers)

	return result, nil
}
