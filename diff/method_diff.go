package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// MethodDiff describes the changes between a pair of operation objects: https://swagger.io/specification/#operation-object
type MethodDiff struct {
	ExtensionsDiff   *ExtensionsDiff           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	TagsDiff         *StringsDiff              `json:"tags,omitempty" yaml:"tags,omitempty"`
	SummaryDiff      *ValueDiff                `json:"summary,omitempty" yaml:"summary,omitempty"`
	DescriptionDiff  *ValueDiff                `json:"description,omitempty" yaml:"description,omitempty"`
	OperationIDDiff  *ValueDiff                `json:"operationID,omitempty" yaml:"operationID,omitempty"`
	ParametersDiff   *ParametersDiff           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBodyDiff  *RequestBodyDiff          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	ResponsesDiff    *ResponsesDiff            `json:"responses,omitempty" yaml:"responses,omitempty"`
	CallbacksDiff    *CallbacksDiff            `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	DeprecatedDiff   *ValueDiff                `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	SecurityDiff     *SecurityRequirementsDiff `json:"securityRequirements,omitempty" yaml:"securityRequirements,omitempty"`
	ServersDiff      *ServersDiff              `json:"servers,omitempty" yaml:"servers,omitempty"`
	ExternalDocsDiff *ExternalDocsDiff         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Base             *openapi3.Operation       `json:"-" yaml:"-"`
	Revision         *openapi3.Operation       `json:"-" yaml:"-"`
}

func newMethodDiff() *MethodDiff {
	return &MethodDiff{}
}

// Empty indicates whether a change was found in this element
func (methodDiff *MethodDiff) Empty() bool {
	if methodDiff == nil {
		return true
	}

	return *methodDiff == MethodDiff{Base: methodDiff.Base, Revision: methodDiff.Revision}
}

func (methodDiff *MethodDiff) removeNonBreaking(config *Config, pathItem2 *openapi3.Operation) {

	if methodDiff.Empty() {
		return
	}

	methodDiff.ExtensionsDiff = nil
	methodDiff.TagsDiff = nil
	methodDiff.SummaryDiff = nil
	methodDiff.DescriptionDiff = nil
	methodDiff.OperationIDDiff = nil
	if deprecationPeriodSufficient(config.DeprecationDays, pathItem2.ExtensionProps) {
		methodDiff.DeprecatedDiff = nil
	}
	methodDiff.ServersDiff = nil
	methodDiff.ExternalDocsDiff = nil
}

func getMethodDiff(config *Config, state *state, pathItem1, pathItem2 *openapi3.Operation) (*MethodDiff, error) {
	diff, err := getMethodDiffInternal(config, state, pathItem1, pathItem2)

	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking(config, pathItem2)
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getMethodDiffInternal(config *Config, state *state, pathItem1, pathItem2 *openapi3.Operation) (*MethodDiff, error) {

	result := newMethodDiff()
	var err error

	result.ExtensionsDiff = getExtensionsDiff(config, state, pathItem1.ExtensionProps, pathItem2.ExtensionProps)
	result.TagsDiff = getStringsDiff(pathItem1.Tags, pathItem2.Tags)
	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiffConditional(config.ExcludeDescription, pathItem1.Description, pathItem2.Description)
	result.OperationIDDiff = getValueDiff(pathItem1.OperationID, pathItem2.OperationID)
	result.ParametersDiff, err = getParametersDiff(config, state, pathItem1.Parameters, pathItem2.Parameters)
	if err != nil {
		return nil, err
	}

	result.RequestBodyDiff, err = getRequestBodyDiff(config, state, pathItem1.RequestBody, pathItem2.RequestBody)
	if err != nil {
		return nil, err
	}

	result.ResponsesDiff, err = getResponsesDiff(config, state, pathItem1.Responses, pathItem2.Responses)
	if err != nil {
		return nil, err
	}

	result.CallbacksDiff, err = getCallbacksDiff(config, state, pathItem1.Callbacks, pathItem2.Callbacks)
	if err != nil {
		return nil, err
	}
	result.DeprecatedDiff = getValueDiff(pathItem1.Deprecated, pathItem2.Deprecated)
	result.SecurityDiff = getSecurityRequirementsDiff(config, state, pathItem1.Security, pathItem2.Security)
	result.ServersDiff = getServersDiff(config, state, pathItem1.Servers, pathItem2.Servers)
	result.ExternalDocsDiff = getExternalDocsDiff(config, state, pathItem1.ExternalDocs, pathItem2.ExternalDocs)
	result.Base = pathItem1
	result.Revision = pathItem2

	return result, nil
}

// Patch applies the patch to a method
func (methodDiff *MethodDiff) Patch(operation *openapi3.Operation) error {

	if methodDiff.Empty() {
		return nil
	}

	methodDiff.DescriptionDiff.patchString(&operation.Description)
	err := methodDiff.ParametersDiff.Patch(operation.Parameters)
	if err != nil {
		return err
	}
	return nil
}
