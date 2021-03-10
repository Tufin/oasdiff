package diff

import "github.com/getkin/kin-openapi/openapi3"

// ComponentsDiff is the diff between two component objects: https://swagger.io/specification/#components-object
type ComponentsDiff struct {
	SchemasDiff         *SchemasDiff         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	ParametersDiff      *ParametersDiff      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	HeadersDiff         *HeadersDiff         `json:"headers,omitempty" yaml:"headers,omitempty"`
	RequestBodiesDiff   *RequestBodiesDiff   `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	ResponsesDiff       *ResponsesDiff       `json:"responses,omitempty" yaml:"responses,omitempty"`
	SecuritySchemesDiff *SecuritySchemesDiff `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	CallbacksDiff       *CallbacksDiff       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

func getComponentsDiff(config *Config, s1, s2 openapi3.Components) ComponentsDiff {

	result := ComponentsDiff{}

	result.SchemasDiff = getSchemasDiff(config, s1.Schemas, s2.Schemas)
	result.ParametersDiff = getParametersDiff(config, toParameters(s1.Parameters), toParameters(s2.Parameters))
	result.HeadersDiff = getHeadersDiff(config, s1.Headers, s2.Headers)
	result.RequestBodiesDiff = getRequestBodiesDiff(config, s1.RequestBodies, s2.RequestBodies)
	result.ResponsesDiff = getResponsesDiff(config, s1.Responses, s2.Responses)
	result.SecuritySchemesDiff = getSecuritySchemesDiff(config, s1.SecuritySchemes, s2.SecuritySchemes)
	// Examples
	// Links
	result.CallbacksDiff = getCallbacksDiff(config, s1.Callbacks, s2.Callbacks)

	return result
}
