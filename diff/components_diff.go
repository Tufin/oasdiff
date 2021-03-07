package diff

import "github.com/getkin/kin-openapi/openapi3"

// ComponentsDiff is the diff between two component objects: https://swagger.io/specification/#components-object
type ComponentsDiff struct {
	SchemasDiff         *SchemasDiff         `json:"schemas,omitempty"`
	ParametersDiff      *ParametersDiff      `json:"parameters,omitempty"`
	HeadersDiff         *HeadersDiff         `json:"headers,omitempty"`
	RequestBodiesDiff   *RequestBodiesDiff   `json:"requestBodies,omitempty"`
	ResponsesDiff       *ResponsesDiff       `json:"responses,omitempty"`
	SecuritySchemesDiff *SecuritySchemesDiff `json:"securitySchemes,omitempty"`
	CallbacksDiff       *CallbacksDiff       `json:"callbacks,omitempty"`
}

func getComponentsDiff(config *Config, s1, s2 openapi3.Components) ComponentsDiff {

	result := ComponentsDiff{}

	result.setSchemasDiff(getSchemasDiff(config, s1.Schemas, s2.Schemas))
	result.setParametersDiff(getParametersDiff(config, toParameters(s1.Parameters), toParameters(s2.Parameters)))
	result.setHeadersDiff(getHeadersDiff(config, s1.Headers, s2.Headers))
	result.setRequestBodiesDiff(getRequestBodiesDiff(config, s1.RequestBodies, s2.RequestBodies))
	result.setResponsesDiff(getResponsesDiff(config, s1.Responses, s2.Responses))
	result.setSecuritySchemesDiff(getSecuritySchemesDiff(config, s1.SecuritySchemes, s2.SecuritySchemes))
	// Examples
	// Links
	result.setCallbacksDiff(getCallbacksDiff(config, s1.Callbacks, s2.Callbacks))

	return result
}

func (componentsDiff *ComponentsDiff) setSchemasDiff(diff *SchemasDiff) {
	componentsDiff.SchemasDiff = nil

	if !diff.empty() {
		componentsDiff.SchemasDiff = diff
	}
}

func (componentsDiff *ComponentsDiff) setParametersDiff(diff *ParametersDiff) {
	componentsDiff.ParametersDiff = nil

	if !diff.empty() {
		componentsDiff.ParametersDiff = diff
	}
}

func (componentsDiff *ComponentsDiff) setHeadersDiff(diff *HeadersDiff) {
	componentsDiff.HeadersDiff = nil

	if !diff.empty() {
		componentsDiff.HeadersDiff = diff
	}
}

func (componentsDiff *ComponentsDiff) setRequestBodiesDiff(diff *RequestBodiesDiff) {
	componentsDiff.RequestBodiesDiff = nil

	if !diff.empty() {
		componentsDiff.RequestBodiesDiff = diff
	}
}

func (componentsDiff *ComponentsDiff) setResponsesDiff(diff *ResponsesDiff) {
	componentsDiff.ResponsesDiff = nil

	if !diff.empty() {
		componentsDiff.ResponsesDiff = diff
	}
}

func (componentsDiff *ComponentsDiff) setSecuritySchemesDiff(diff *SecuritySchemesDiff) {
	componentsDiff.SecuritySchemesDiff = nil

	if !diff.empty() {
		componentsDiff.SecuritySchemesDiff = diff
	}
}

func (componentsDiff *ComponentsDiff) setCallbacksDiff(diff *CallbacksDiff) {
	componentsDiff.CallbacksDiff = nil

	if !diff.empty() {
		componentsDiff.CallbacksDiff = diff
	}
}
